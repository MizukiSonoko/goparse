// Copyright (C) 2018,2019 MizukiSonoko. All rights reserved.

package goparse

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type Result interface {
	// Insert insets format values to dest
	Insert(dest ...interface{}) error

	// Insert inserts a selected format value to dest
	InsertOnly(index int, dest interface{}) error
}

// parseString returns string before format
//  ( format=" %s ", str= "a b c") => "a"
//  ( format="%s", str= "nnnn") => "n"
//  ( format="(%s)", str= "(yes)(no)") => "(yes)"
//  ( format="or", str= "(yes)or(no)") => "(yes)"
func parseString(format, str string) (string, error) {

	// This case is happened by %s is in end of a text.
	if len(format) == 0 {
		return str, nil
	}

	terms := strings.Split(format, "%")
	i := strings.Index(str, terms[0])
	if i == -1 {
		return "", fmt.Errorf("LastIndex return -1, [%s] not contains [%s]",
			str, terms[0])
	} else if i == 0 {
		ni := strings.Index(str[1:], terms[0])
		if ni == -1 {
			return str[:i], nil
		}
		i = ni + 1
	}
	return str[:i], nil
}

// parseInteger returns number
// (format="yen",str="123yen",base=10) => 123
// (format="yen",str="10101yen",base=2) => 21
func parseInteger(format, str string, base int) (int, error) {
	s, err := parseString(format, str)
	if err != nil {
		return 0, errors.Wrapf(err, "parseString(%s,%s) failed",
			format, str)
	}
	// ToDo We should think about int64
	i, err := strconv.ParseInt(s, base, 0)
	if err != nil {
		return 0, errors.Wrapf(err, "ParseInt(\"%s\",%d) failed", s, base)
	}
	return int(i), nil
}

func parseBool(format, str string) (bool, error) {
	s, err := parseString(format, str)
	if err != nil {
		return false, errors.Wrapf(err, "parseString(%s,%s) failed",
			format, str)
	}
	b, err := strconv.ParseBool(s)
	if err != nil {
		return false, errors.Wrapf(err, "ParseBool(%s) failed", s)
	}
	return b, nil
}

func parseFloat(format, str string) (float64, error) {
	s, err := parseString(format, str)
	if err != nil {
		return 0, errors.Wrapf(err, "parseString(%s,%s) failed",
			format, str)
	}
	f, err := strconv.ParseFloat(s, 0)
	if err != nil {
		return 0, errors.Wrapf(err, "ParseFloat(%s) failed", s)
	}
	return f, nil
}

type value struct {
	kind  reflect.Kind
	value interface{}
}

type result struct {
	err    error
	values []value
}

func assign(dest interface{}, src value) error {
	switch src.kind {
	case reflect.String:
		switch d := dest.(type) {
		case *string:
			if d == nil {
				return fmt.Errorf("destination pointer should not be nil")
			}
			*d = src.value.(string)
			return nil
		case *[]byte:
			if d == nil {
				return fmt.Errorf("destination pointer should not be nil")
			}
			*d = []byte(src.value.(string))
			return nil
		default:
			return fmt.Errorf("type mismatch: expected *string,*[]byte, actual %T", d)
		}
	case reflect.Int:
		switch d := dest.(type) {
		case *int:
			*d = src.value.(int)
			return nil
		case *int8:
			if src.value.(int) > math.MaxInt8 {
				return fmt.Errorf("overflow: %d is greater than MaxInt8(%d)",
					src.value.(int), math.MaxInt8)
			}
			*d = int8(src.value.(int))
			return nil
		case *int32:
			if src.value.(int) > math.MaxInt32 {
				return fmt.Errorf("overflow: %d is greater than MaxInt32(%d)",
					src.value.(int), math.MaxInt32)
			}
			*d = int32(src.value.(int))
			return nil
		case *int64:
			// Note: if value is over MaxInt64, strconv.ParseInt failed
			*d = int64(src.value.(int))
			return nil
		default:
			return fmt.Errorf("type mismatch: expected *int{8,32,64}, actual %T", d)
		}
	case reflect.Bool:
		switch d := dest.(type) {
		case *bool:
			*d = src.value.(bool)
			return nil
		}
	case reflect.Float64:
		switch d := dest.(type) {
		case *float64:
			*d = src.value.(float64)
			return nil
		case *float32:
			*d = float32(src.value.(float64))
			return nil
		}
	}
	return fmt.Errorf("unsupported type %s into type %s",
		src.kind.String(), reflect.TypeOf(dest).Kind().String())
}

func (r result) Insert(dest ...interface{}) error {
	// r.err is happened by parser
	if r.err != nil {
		return r.err
	}

	if len(dest) != len(r.values) {
		return fmt.Errorf(
			"expected %d destination arguments in Insert, not %d",
			len(r.values), len(dest))
	}
	for i, sv := range r.values {
		err := assign(dest[i], sv)
		if err != nil {
			return fmt.Errorf(`assign(src{kind:%s,%v} => dest[%d]) failed err:%s`,
				sv.kind.String(), sv.value, i, err)
		}
	}
	return nil
}

func (r result) InsertOnly(index int, dest interface{}) error {
	// r.err is happened by parser
	if r.err != nil {
		return r.err
	}

	if index >= len(r.values) {
		return fmt.Errorf(
			"invalid index:%d, format is only %d",
			index, len(r.values))
	}

	for i, sv := range r.values {
		if i == index {
			err := assign(dest, sv)
			if err != nil {
				return fmt.Errorf(`assign(src{kind:%s,%v} => dest[%d]) failed err:%s`,
					sv.kind.String(), sv.value, i, err)
			}
		}
	}
	return nil
}

// Parse parse str uses format
func Parse(format, str string) Result {
	var res result
	end := len(format)
	strOffset := 0

	for i := 0; i < end; {
		if format[i] != '%' && format[i] != str[strOffset+i] {
			return result{
				err: fmt.Errorf("invalid string (%s) with (%s). expect %c but it is %c",
					str, format, format[i], str[strOffset+i]),
			}
		}

		if format[i] == '%' {
			for ; i < end; i++ {
				c := format[i]
				switch c {
				case 's':
					// first arguments except format
					s, err := parseString(format[i+1:], str[strOffset+i-1:])
					if err != nil {
						return result{
							err: errors.Wrapf(err, "parseString(%s,%s) failed",
								format[i:], str[strOffset+i-1:]),
						}
					}
					strOffset += len(s) - 2
					res.values = append(res.values, value{
						reflect.String, s})
					i += 2
					goto formatLoop
				case 'd':
					// first arguments except format
					n, err := parseInteger(format[i+1:], str[strOffset+i-1:], 10)
					if err != nil {
						return result{
							err: errors.Wrapf(err, "parseInteger(%%%s,\"%s\",10) failed",
								format[i:], str[strOffset+i-1:]),
						}
					}
					strOffset += len(strconv.Itoa(n)) - 2
					res.values = append(res.values, value{
						reflect.Int, n})
					i += 2
					goto formatLoop
				case 'b':
					// first arguments except format
					n, err := parseInteger(format[i+1:], str[strOffset+i-1:], 2)
					if err != nil {
						return result{
							err: errors.Wrapf(err, "parseInteger(%%%s,\"%s\",2) failed",
								format[i:], str[strOffset+i-1:]),
						}
					}
					offset, err := parseString(format[i+1:], str[strOffset+i-1:])
					if err != nil {
						return result{
							err: errors.Wrapf(err, "parseString(%%%s,\"%s\",2) failed",
								format[i+1:], str[strOffset+i-1:]),
						}
					}
					strOffset += len(offset) - 2
					res.values = append(res.values, value{
						reflect.Int, n})
					i += 2
					goto formatLoop
				case 'o':
					// first arguments except format
					n, err := parseInteger(format[i+1:], str[strOffset+i-1:], 8)
					if err != nil {
						return result{
							err: errors.Wrapf(err, "parseInteger(%s,%s,8) failed",
								format[i:], str[strOffset+i-1:]),
						}
					}
					offset, err := parseString(format[i+1:], str[strOffset+i-1:])
					if err != nil {
						return result{
							err: errors.Wrapf(err, "parseString(%%%s,\"%s\",2) failed",
								format[i+1:], str[strOffset+i-1:]),
						}
					}
					strOffset += len(offset) - 2
					res.values = append(res.values, value{
						reflect.Int, n})
					i += 2
					goto formatLoop
				case 't':
					// first arguments except format
					b, err := parseBool(format[i+1:], str[strOffset+i-1:])
					if err != nil {
						return result{
							err: errors.Wrapf(err, "parseInteger(%s,%s) failed",
								format[i:], str[strOffset+i-1:]),
						}
					}
					strOffset += len(strconv.FormatBool(b)) - 2
					res.values = append(res.values, value{
						reflect.Bool, b})
					i += 2
					goto formatLoop
				case 'f':
					// first arguments except format
					f, err := parseFloat(format[i+1:], str[strOffset+i-1:])
					if err != nil {
						return result{
							err: errors.Wrapf(err, "parseFloat(%s,%s) failed",
								format[i:], str[strOffset+i-1:]),
						}
					}
					offset, err := parseString(format[i+1:], str[strOffset+i-1:])
					if err != nil {
						return result{
							err: errors.Wrapf(err, "parseString(\"%s\",\"%s\") failed",
								format[i+1:], str[strOffset+i-1:]),
						}
					}
					strOffset += len(offset) - 2
					res.values = append(res.values, value{
						reflect.Float64, f})
					i += 2
					goto formatLoop
				}
			}
		}
		i++
	formatLoop:
	}
	return res
}
