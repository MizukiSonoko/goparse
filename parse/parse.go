// Copyright (C) 2018 MizukiSonoko. All rights reserved.

package goparse

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"strconv"
	"strings"
)

// A Result is returned by Parse
// It has kind like string, int,..., and value
type Result struct {
	kind  reflect.Kind
	value interface{}
}

// Kind returns a type of value. string, int, ...
func (r Result) Kind() reflect.Kind {
	return r.kind
}

// Value returns a value as interface
func (r Result) Value() interface{} {
	return r.value
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

func parseInteger(format, str string) (int, error) {
	s, err := parseString(format, str)
	if err != nil {
		return 0, errors.Wrapf(err, "parseString(%s,%s) failed",
			format, str)
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, errors.Wrapf(err, "Atoi(%s) failed", s)
	}
	return i, nil
}

// Parse parse str uses format
func Parse(format, str string) ([]Result, error) {
	var res []Result
	end := len(format)
	strOffset := 0

	for i := 0; i < end; {
		if format[i] != '%' && format[i] != str[strOffset+i] {
			return res, fmt.Errorf("invalid string (%s) with (%s). expect %c but it is %c",
				str, format, format[i], str[strOffset+i])
		}

		if format[i] == '%' {
			for ; i < end; i++ {
				c := format[i]
				switch c {
				case 's':
					// first arguments except format
					s, err := parseString(format[i+1:], str[strOffset+i-1:])
					if err != nil {
						return res, errors.Wrapf(err, "parseString(%s,%s) failed",
							format[i:], str[strOffset+i-1:])
					}
					strOffset += len(s) - 2
					res = append(res, Result{reflect.String, s})
					i += 2
					goto formatLoop
				case 'd':
					// first arguments except format
					n, err := parseInteger(format[i+1:], str[strOffset+i-1:])
					if err != nil {
						return res, errors.Wrapf(err, "parseInteger(%s,%s) failed",
							format[i:], str[strOffset+i-1:])
					}
					strOffset += len(strconv.Itoa(n)) - 2
					res = append(res, Result{reflect.Int, n})
					i += 2
					goto formatLoop
				}
			}
		}
		i++
	formatLoop:
	}
	return res, nil
}
