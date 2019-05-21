
## Goparse 

[![CircleCI](https://circleci.com/gh/MizukiSonoko/goparse.svg?style=shield)](https://circleci.com/gh/MizukiSonoko/goparse)
[![codecov](https://codecov.io/gh/MizukiSonoko/goparse/branch/master/graph/badge.svg)](https://codecov.io/gh/MizukiSonoko/goparse)
[![Go Report Card](https://goreportcard.com/badge/github.com/MizukiSonoko/goparse)](https://goreportcard.com/report/github.com/MizukiSonoko/goparse)
[![MIT licensed](https://img.shields.io/packagist/l/doctrine/orm.svg)](https://github.com/MizukiSonoko/goparse/blob/master/LICENSE)
[![Godoc](https://godoc.org/github.com/MizukiSonoko/goparse/parse?status.svg)](https://godoc.org/github.com/MizukiSonoko/goparse/parse)
  
  
This library is inspired by [r1chardj0n3s/parse](https://github.com/r1chardj0n3s/parse) in Python

> Parse() is the opposite of fmt.Sprintf()

```go
var s string
err := goparse.Parse("Hello %s", "Hello World").Insert(&s)
fmt.Println(s)
// Output:
// World
```

## Example

### Single string
```go
format := "Hello %s"
expected := "World"
var res string
_ = goparse.Parse(format,fmt.Sprintf(format,expected)).Insert(&res)
fmt.Println(s)
// Output:
// World
```

### Multiple string
```go
format := "水樹素子「%s」。秋穂伊織「%s」"
str := "水樹素子「今日は天気が悪いね」。秋穂伊織「そうだね」"
var mizukiMsg, ioriMsg string
_ = goparse.Parse(format,str).Insert(&mizukiMsg, &ioriMsg)
fmt.Println(mizukiMsg)
fmt.Println(ioriMsg)
// Output:
// 今日は天気が悪いね
// そうだね
```

### Base2 integer
```go
format := "Robot says '%d'"
expected := 12345
var num int
_ = goparse.Parse(format, fmt.Sprintf(format, expected)).Insert(&num)
fmt.Println(num)
// Output:
// 12345
```

### Base10 integer
```go
format := "Hello my number is %d"
expected := 100
var num int
_ = goparse.Parse(format, fmt.Sprintf(format, expected)).Insert(&num)
fmt.Println(num)
// Output:
// 100
```

### Base8 integer
```go
format := "Hello my number is %o"
expected := 123
var numOct int
_ = goparse.Parse(format, fmt.Sprintf(format, expected)).Insert(&numOct)
fmt.Println(numOct)
// Output:
// 123
```

### Boolean
```go
format := "I can't tell whether it is %t or %t"
str := "I can't tell whether it is false or true"
var boolRes1, boolRes2 bool
_ = goparse.Parse(format, str).Insert(&boolRes1,&boolRes2)
fmt.Println(boolRes1)
fmt.Println(boolRes2)
// Output:
// false
// true
```

## Installation

```sh
go get github.com/MizukiSonoko/goparse
```

## The format 'verbs'
Cite by https://golang.org/pkg/fmt/

I support this verbs as follows:

### Boolean:
```
[o] %t	the word true or false
```

### Integer:
```
[o] %b	base 2
[o] %d	base 10
[o] %o	base 8
```

### String and slice of bytes (treated equivalently with these verbs):
```
[o] %s	the uninterpreted bytes of the string or slice
```