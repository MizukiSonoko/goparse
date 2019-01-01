// Copyright (C) 2018 MizukiSonoko. All rights reserved.

package goparse

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type Result interface {
	Insert(dest ...interface{}) error
}

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

func parseInteger(format, str string, base int) (int, error) {
	s, err := parseString(format, str)
	if err != nil {
		return 0, errors.Wrapf(err, "parseString(%s,%s) failed",
			format, str)
	}
	// ToDo We should think about int64
	i, err := strconv.ParseInt(s, base, 0)
	if err != nil {
		return 0, errors.Wrapf(err, "ParseInt(%s,%d) failed", s, base)
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
			*d = src.value.(int8)
			return nil
		case *int32:
			*d = src.value.(int32)
			return nil
		case *int64:
			*d = src.value.(int64)
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
				case 'o':
					// first arguments except format
					n, err := parseInteger(format[i+1:], str[strOffset+i-1:], 8)
					if err != nil {
						return result{
							err: errors.Wrapf(err, "parseInteger(%s,%s,8) failed",
								format[i:], str[strOffset+i-1:]),
						}
					}
					strOffset += len(strconv.Itoa(n)) - 2
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
				}
			}
		}
		i++
	formatLoop:
	}
	return res
}
