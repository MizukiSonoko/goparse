
## Goparse 

[![Go Report Card](https://goreportcard.com/badge/github.com/MizukiSonoko/goparse)](https://goreportcard.com/report/github.com/MizukiSonoko/goparse)
[![MIT licensed](https://img.shields.io/packagist/l/doctrine/orm.svg)](https://github.com/MizukiSonoko/goparse/blob/master/LICENSE)
[![Godoc](https://godoc.org/github.com/MizukiSonoko/goparse?status.svg)](https://godoc.org/github.com/MizukiSonoko/goparse)
  
  
This library is inspired by [r1chardj0n3s/parse](https://github.com/r1chardj0n3s/parse) in Python

> Parse() is the opposite of fmt.Sprintf()

```go
res, err := goparse.Parse("Hello %s", "Hello World")
fmt.Println(res[0].Value())
// 'World'
```

```go
format := "Hello %s"
expected := "World"
res, _ := goparse.Parse(format,fmt.Sprintf(format,expected))
assert.Equal(t,expected,res[0].Value())
```

```go
format := "水樹素子「%s」。秋穂伊織「%s」"
str := "水樹素子「今日は天気が悪いね」。秋穂伊織「そうだね」"
expected1 := "今日は天気が悪いね"
expected2 := "そうだね"
res, _ := goparse.Parse(format,str)
assert.Equal(t,expected1,res[0].Value())
assert.Equal(t,expected2,res[1].Value())
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

 