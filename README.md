
## Goparse 

[![CircleCI](https://circleci.com/gh/MizukiSonoko/goparse.svg?style=shield)](https://circleci.com/gh/MizukiSonoko/goparse)
[![codecov](https://codecov.io/gh/MizukiSonoko/goparse/branch/master/graph/badge.svg)](https://codecov.io/gh/MizukiSonoko/goparse)
[![Go Report Card](https://goreportcard.com/badge/github.com/MizukiSonoko/goparse)](https://goreportcard.com/report/github.com/MizukiSonoko/goparse)
[![MIT licensed](https://img.shields.io/packagist/l/doctrine/orm.svg)](https://github.com/MizukiSonoko/goparse/blob/master/LICENSE)
[![Godoc](https://godoc.org/github.com/MizukiSonoko/goparse?status.svg)](https://godoc.org/github.com/MizukiSonoko/goparse)
  
  
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

## ToDo

- [x] support string `%s`
- [ ] support integer `%b,%d,%o,...`
- [ ] support float `%f,%e,...`
- [ ] other flag `+`, `#`

 
