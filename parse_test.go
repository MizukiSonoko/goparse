// Copyright (C) 2018 MizukiSonoko. All rights reserved.

package goparse_test

import (
	"fmt"
	"github.com/MizukiSonoko/goparse"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strings"
	"testing"
)

func checkTestCase(t *testing.T, expected, format string, str ...interface{}) {
	if expected != fmt.Sprintf(format, str...) {
		panic(
			"Invalid test case. expected by test = '" + expected + "' want '" + fmt.Sprintf(format, str...) + "'")
	}
}

func TestParseString(t *testing.T) {

	t.Run("Normal case", func(t *testing.T) {
		format := "Dayo"
		str := "SonokoDayo"
		expected := "Sonoko"
		res, err := goparse.ParseStringForTest(format, str)
		assert.NoErrorf(t, err, "parseString(%s,%s) not failed want fail",
			format, str)
		assert.Equal(t, expected, res)
	})

	t.Run("last string", func(t *testing.T) {
		format := ""
		str := "Sonoko"
		expected := "Sonoko"
		res, err := goparse.ParseStringForTest(format, str)
		assert.NoErrorf(t, err, "parseString(%s,%s) not failed want fail",
			format, str)
		assert.Equal(t, expected, res)
	})

	t.Run("contains duplicated text", func(t *testing.T) {
		format := "OkOkOk"
		str := "OkOkOkOk"
		expected := "Ok"
		res, err := goparse.ParseStringForTest(format, str)
		assert.NoErrorf(t, err, "parseString(%s,%s) not failed want fail",
			format, str)
		assert.Equal(t, expected, res)
	})

	t.Run("format is same as str", func(t *testing.T) {
		format := "Ok"
		str := "Ok"
		expected := ""
		res, err := goparse.ParseStringForTest(format, str)
		assert.NoErrorf(t, err, "parseString(%s,%s) not failed want fail",
			format, str)
		assert.Equal(t, expected, res)
	})

	t.Run("format contains %s", func(t *testing.T) {
		format := "Ok%sOk"
		str := "OkOk"
		expected := "Ok"
		res, err := goparse.ParseStringForTest(format, str)
		assert.NoErrorf(t, err, "parseString(%s,%s) not failed want fail",
			format, str)
		assert.Equal(t, expected, res)
	})

	t.Run("format contains %s at the end", func(t *testing.T) {
		format := "Ok%s"
		str := "OkOk"
		expected := "Ok"
		res, err := goparse.ParseStringForTest(format, str)
		assert.NoErrorf(t, err, "parseString(%s,%s) not failed want fail",
			format, str)
		assert.Equal(t, expected, res)
	})

	t.Run("format contains three %s", func(t *testing.T) {
		format := "_%s_%s_%s"
		str := "han_maru_gin"
		expected := "han"
		res, err := goparse.ParseStringForTest(format, str)
		assert.NoErrorf(t, err, "parseString(%s,%s) not failed want fail",
			format, str)
		assert.Equal(t, expected, res)
	})

}

func TestParse(t *testing.T) {

	t.Run("the mix of %s and %d", func(t *testing.T) {
		format := "Hello %s, my number is %d"
		str := "Hello iorin, my number is 9753"
		expected1 := "iorin"
		expected2 := 9753
		res, _ := goparse.Parse(format, str)
		assert.Equal(t, expected1, res[0].Value())
		assert.Equal(t, expected2, res[1].Value())
	})
}

func TestParse_string(t *testing.T) {

	checkString := func(t *testing.T, expected string, actual goparse.Result) {
		if reflect.String != actual.Kind() {
			t.Errorf("Kind = <%d> want <%d>", actual.Kind(), reflect.String)
			return
		}
		if reflect.TypeOf(string("")) != reflect.TypeOf(actual.Value()) {
			t.Errorf("type(res[0].Value) = <%v> want <%v>",
				reflect.TypeOf(actual.Value()), reflect.TypeOf(string("")))
			return
		}
		if expected != actual.Value().(string) {
			t.Errorf("res[0].Value = <%s> want <%s>",
				actual.Value().(string), expected)
		}
	}

	t.Run("The opposite of Sprintf", func(t *testing.T) {
		format := "Hello %s"
		expected := "World"
		res, _ := goparse.Parse(format, fmt.Sprintf(format, expected))
		assert.Equal(t, expected, res[0].Value())
	})

	t.Run("format is splitted by blank", func(t *testing.T) {
		format := "Hello %s"
		str := "Hello World"
		expected := "World"
		checkTestCase(t, str, format, expected)
		res, err := goparse.Parse(format, str)
		assert.NoErrorf(t, err, "Parse(%s,%s) failed")
		if len(res) != 1 {
			t.Errorf("len(res) = <%d> want <1>", len(res))
		}
		checkString(t, expected, res[0])
	})

	t.Run("format has one %s", func(t *testing.T) {
		for _, tt := range []struct {
			format   string
			str      string
			expected string
		}{
			{
				format:   "Good_%s_",
				str:      "Good_Morning_",
				expected: "Morning",
			},
			{
				format:   "Hello%sMosaic",
				str:      "HelloGoldenMosaic",
				expected: "Golden",
			},
			{
				format:   "%ssonoko",
				str:      "ssssonoko",
				expected: "sss",
			},
		} {
			t.Logf("test case: Parse(%s,%s) failed", tt.format, tt.str)
			checkTestCase(t, tt.str, tt.format, tt.expected)
			res, err := goparse.Parse(tt.format, tt.str)
			assert.NoErrorf(t, err, "Parse(%s,%s) failed", tt.format, tt.str)
			if len(res) != 1 {
				t.Errorf("len(res) = <%d> want <1>", len(res))
			}
			checkString(t, tt.expected, res[0])
		}
	})

	t.Run("format contains Number", func(t *testing.T) {
		format := "12%s90"
		str := "1234567890"
		expected := "345678"
		checkTestCase(t, str, format, expected)
		res, err := goparse.Parse(format, str)
		assert.NoErrorf(t, err, "Parse(%s,%s) failed")
		if len(res) != 1 {
			t.Errorf("len(res) = <%d> want <1>", len(res))
		}
		checkString(t, expected, res[0])
	})

	t.Run("format contains 日本語", func(t *testing.T) {
		format := "Hello %s"
		str := "Hello こんにちは"
		expected := "こんにちは"
		checkTestCase(t, str, format, expected)
		res, err := goparse.Parse(format, str)
		assert.NoErrorf(t, err, "Parse(%s,%s) failed")
		if len(res) != 1 {
			t.Errorf("len(res) = <%d> want <1>", len(res))
		}
		checkString(t, expected, res[0])
	})

	t.Run("format contains 日本語 part 2", func(t *testing.T) {
		format := "み%sっ%sのこ"
		str := "みかんとずっきーにときのこ"
		expected1 := "かんとず"
		expected2 := "きーにとき"
		checkTestCase(t, str, format, expected1, expected2)
		res, err := goparse.Parse(format, str)
		assert.NoErrorf(t, err, "Parse(%s,%s) failed")
		if len(res) != 2 {
			t.Errorf("len(res) = <%d> want <2>", len(res))
		}
		checkString(t, expected1, res[0])
		checkString(t, expected2, res[1])
	})

	t.Run("format contains 日本語 part 3", func(t *testing.T) {
		format := "水樹素子「%s」。秋穂伊織「%s」"
		str := "水樹素子「今日は天気が悪いね」。秋穂伊織「そうだね」"
		expected1 := "今日は天気が悪いね"
		expected2 := "そうだね"
		checkTestCase(t, str, format, expected1, expected2)
		res, err := goparse.Parse(format, str)
		assert.NoErrorf(t, err, "Parse(%s,%s) failed")
		if len(res) != 2 {
			t.Errorf("len(res) = <%d> want <1>", len(res))
		}
		checkString(t, expected1, res[0])
		checkString(t, expected2, res[1])
	})

	t.Run("text contains multiple %s", func(t *testing.T) {
		format := "Hello %s!, How are you? %s?"
		str := "Hello Mizuki!, How are you? Ok?"
		expected1 := "Mizuki"
		expected2 := "Ok"
		checkTestCase(t, str, format, expected1, expected2)
		res, err := goparse.Parse(format, str)
		assert.NoErrorf(t, err, "Parse(%s,%s) failed")
		if len(res) != 2 {
			t.Fatalf("len(res) = <%d> want <2>", len(res))
		}
		checkString(t, expected1, res[0])
		checkString(t, expected2, res[1])
	})

	asI := func(ss []string) []interface{} {
		res := make([]interface{}, len(ss))
		for i, str := range ss {
			res[i] = str
		}
		return res
	}

	t.Run("text contains many %s", func(t *testing.T) {
		names := []string{
			"chiyoda",
			"chuo",
			"shinagawa",
			"shinjuku",
			"shibuya",
			"taito",
			"edogawa",
			"nakano",
			"suginami",
			"katsushika",
			"kita",
			"minato",
			"itabashi",
			"toshima",
			"adachi",
			"oota",
			"sumida",
			"bunkyo",
			"koto",
			"setagaya",
			"nerima",
			"arakawa",
			"meguro",
		}
		str := ""
		format := ""
		for i := 1; i < len(names); i++ {
			str = strings.Join(names[:i], "_")
			format = "%s" + strings.Repeat("_%s", i-1)
			checkTestCase(t, str, format, asI(names[:i])...)

			res, err := goparse.Parse(format, str)
			assert.NoErrorf(t, err, "Parse(%s,%s) failed")
			if len(res) != i {
				t.Fatalf("len(res) = <%d> want <6>", len(res))
			}
			for j := range names[:i] {
				checkString(t, names[j], res[j])
			}
		}
	})

	t.Run("Invalid argumetns", func(t *testing.T) {
		t.Run("No match", func(t *testing.T) {
			format := "Hello"
			str := "noHello"
			_, err := goparse.Parse(format, str)
			assert.Errorf(t, err, "Parse(%s,%s) not failed want fail")
			assert.Contains(t, err.Error(), "invalid string")
		})

		t.Run("a cuple of %s", func(t *testing.T) {
			format := "%s%s%s"
			str := "Hello"
			_, err := goparse.Parse(format, str)
			assert.Errorf(t, err, "Parse(%s,%s) not failed want fail")
			assert.Contains(t, err.Error(), "invalid string")
		})

		t.Run("different number of %s from str", func(t *testing.T) {
			format := "%s_%s_%s"
			str := "H_He"
			_, err := goparse.Parse(format, str)
			assert.Errorf(t, err, "Parse(%s,%s) not failed want fail")
			assert.Contains(t, err.Error(), "parseString")
		})
	})
}

func TestParse_integer(t *testing.T) {

	checkInt := func(t *testing.T, expected int, actual goparse.Result) {
		if reflect.Int != actual.Kind() {
			t.Errorf("Kind = <%d> want <%d>", actual.Kind(), reflect.Int)
			return
		}
		if reflect.TypeOf(int(0)) != reflect.TypeOf(actual.Value()) {
			t.Errorf("type(res[0].Value) = <%v> want <%v>",
				reflect.TypeOf(actual.Value()), reflect.TypeOf(int(0)))
			return
		}
		if expected != actual.Value().(int) {
			t.Errorf("res[0].Value = <%d> want <%d>",
				actual.Value().(int), expected)
		}
	}

	t.Run("The opposite of Sprintf", func(t *testing.T) {
		format := "Hello my number is %d"
		expected := 100
		res, err := goparse.Parse(format, fmt.Sprintf(format, expected))
		assert.NoError(t, err)
		assert.Equal(t, expected, res[0].Value())
	})

	t.Run("text contains multiple %d", func(t *testing.T) {
		format := "1%d456%d89"
		str := "123456789"
		expected1 := 23
		expected2 := 7
		checkTestCase(t, str, format, expected1, expected2)
		res, err := goparse.Parse(format, str)
		assert.NoErrorf(t, err, "Parse(%s,%s) failed")
		if len(res) != 2 {
			t.Fatalf("len(res) = <%d> want <2>", len(res))
		}
		checkInt(t, expected1, res[0])
		checkInt(t, expected2, res[1])
	})

}

func ExampleParse() {
	res, _ := goparse.Parse("Hello %s", "Hello World")
	fmt.Println(res[0].Value())
	// Output:
	// World
}

func ExampleParse_ja() {
	format := "水樹素子「%s」。秋穂伊織「%s」"
	str := "水樹素子「今日は天気が悪いね」。秋穂伊織「そうだね」"
	res, _ := goparse.Parse(format, str)
	fmt.Println(res[0].Value())
	fmt.Println(res[1].Value())
	// Output:
	// 今日は天気が悪いね
	// そうだね
}
