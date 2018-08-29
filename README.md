
## Goparse 

[![CircleCI](https://circleci.com/gh/MizukiSonoko/goparse.svg?style=shield)](https://circleci.com/gh/MizukiSonoko/goparse)
[![codecov](https://codecov.io/gh/MizukiSonoko/goparse/branch/master/graph/badge.svg)](https://codecov.io/gh/MizukiSonoko/goparse)
[![Go Report Card](https://goreportcard.com/badge/github.com/MizukiSonoko/goparse)](https://goreportcard.com/report/github.com/MizukiSonoko/goparse)
[![MIT licensed](https://img.shields.io/packagist/l/doctrine/orm.svg)](https://github.com/MizukiSonoko/goparse/blob/master/LICENSE)
[![Godoc](https://godoc.org/github.com/MizukiSonoko/goparse/parse?status.svg)](https://godoc.org/github.com/MizukiSonoko/goparse/parse)
  
  
This library is inspired by [r1chardj0n3s/parse](https://github.com/r1chardj0n3s/parse) in Python

> Parse() is the opposite of fmt.Sprintf()

```go
res, err := goparse.Parse("Hello %s", "Hello World")
fmt.Println(res[0].Value())
// Output:
// World
```

```go
format := "Hello %s"
expected := "World"
res, _ := goparse.Parse(format,fmt.Sprintf(format,expected))
fmt.Println(res[0].Value())
// Output:
// World
```

```go
format := "Hello my number is %d"
expected := 100
res, _ := goparse.Parse(format, fmt.Sprintf(format, expected))
fmt.Println(res[0].Value())
// Output:
// 100
```

```go
format := "水樹素子「%s」。秋穂伊織「%s」"
str := "水樹素子「今日は天気が悪いね」。秋穂伊織「そうだね」"
expected1 := "今日は天気が悪いね"
expected2 := "そうだね"
res, _ := goparse.Parse(format,str)
fmt.Println(res[0].Value())
fmt.Println(res[1].Value())
// Output:
// 今日は天気が悪いね
// そうだね
```

## Installation

```sh
go get github.com/MizukiSonoko/goparse
```

## The format 'verbs'
Cite by https://golang.org/pkg/fmt/

- `[ ]` blank means i should decide goparse supports or not
- `[o]` already implemented
- `[x]` not supported
- `[A]` will be supported

### General:
```
[ ] %v	the value in a default format
	when printing structs, the plus flag (%+v) adds field names
[ ] %#v	a Go-syntax representation of the value
[ ] %T	a Go-syntax representation of the type of the value
[ ] %%	a literal percent sign; consumes no value
```

### Boolean:
```
[A] %t	the word true or false
```

### Integer:
```
[ ] %b	base 2
[ ] %c	the character represented by the corresponding Unicode code point
[o] %d	base 10
[ ] %o	base 8
[ ] %q	a single-quoted character literal safely escaped with Go syntax.
[ ] %x	base 16, with lower-case letters for a-f
[ ] %X	base 16, with upper-case letters for A-F
[ ] %U	Unicode format: U+1234; same as "U+%04X"
```

### Floating-point and complex constituents:
```
[ ] %b	decimalless scientific notation with exponent a power of two,
	in the manner of strconv.FormatFloat with the 'b' format,
	e.g. -123456p-78
[ ] %e	scientific notation, e.g. -1.234456e+78
[ ] %E	scientific notation, e.g. -1.234456E+78
[ ] %f	decimal point but no exponent, e.g. 123.456
[ ] %F	synonym for %f
[ ] %g	%e for large exponents, %f otherwise. Precision is discussed below.
[ ] %G	%E for large exponents, %F otherwise
```

### String and slice of bytes (treated equivalently with these verbs):
```
[o] %s	the uninterpreted bytes of the string or slice
[ ] %q	a double-quoted string safely escaped with Go syntax
[ ] %x	base 16, lower-case, two characters per byte
[ ] %X	base 16, upper-case, two characters per byte
```