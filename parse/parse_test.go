// Copyright (C) 2018 MizukiSonoko. All rights reserved.

package goparse_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/MizukiSonoko/goparse/parse"
	"github.com/stretchr/testify/assert"
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
		var res1 string
		var res2 int
		err := goparse.Parse(format, str).Insert(&res1, &res2)
		assert.NoError(t, err)
		assert.Equal(t, expected1, res1)
		assert.Equal(t, expected2, res2)
	})

	t.Run("format contains two %s, but the argument is one", func(t *testing.T) {
		format := "Hello %s, i'm %s"
		str := "Hello Iori, i'm sonoko"
		var res string
		err := goparse.Parse(format, str).Insert(&res)
		assert.Error(t, err)
	})
}

func TestParse_string(t *testing.T) {

	t.Run("The opposite of Sprintf", func(t *testing.T) {
		format := "Hello %s"
		expected := "World"
		var res string

		err := goparse.Parse(format, fmt.Sprintf(format, expected)).Insert(&res)
		assert.NoError(t, err)

		assert.Equal(t, expected, res)
	})

	t.Run("the argument is []byte", func(t *testing.T) {
		format := "Hello %s"
		str := "Hello iorin"
		expected := []byte("iorin")
		var res []byte

		err := goparse.Parse(format, str).Insert(&res)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("the argument is not string", func(t *testing.T) {
		format := "Hello %s"
		str := "Hello sodiu"
		var res int

		err := goparse.Parse(format, str).Insert(&res)
		assert.Error(t, err)
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
			t.Logf("test case: Parse(%s,%s)", tt.format, tt.str)
			checkTestCase(t, tt.str, tt.format, tt.expected)
			var res string
			err := goparse.Parse(tt.format, tt.str).Insert(&res)
			assert.NoErrorf(t, err, "Parse(%s,%s).Insert failed", tt.format, tt.str)
			assert.Equal(t, tt.expected, res)
		}
	})

	t.Run("format contains Number", func(t *testing.T) {
		for _, tt := range []struct {
			format   string
			str      string
			expected string
		}{
			{
				format:   "12%s90",
				str:      "1234567890",
				expected: "345678",
			},
			{
				format:   "12_%s_90",
				str:      "12______90",
				expected: "____",
			},
			{
				format:   "%s4567",
				str:      "1234567",
				expected: "123",
			},
		} {
			t.Logf("test case: Parse(%s,%s)", tt.format, tt.str)
			checkTestCase(t, tt.str, tt.format, tt.expected)
			var res string
			err := goparse.Parse(tt.format, tt.str).Insert(&res)
			assert.NoErrorf(t, err, "Parse(%s,%s).Insert failed", tt.format, tt.str)
			assert.Equal(t, tt.expected, res)
		}
	})

	t.Run("format contains 日本語", func(t *testing.T) {
		format := "Hello %s"
		str := "Hello こんにちは"
		expected := "こんにちは"

		var res string
		err := goparse.Parse(format, str).Insert(&res)
		assert.NoErrorf(t, err, "Parse(%s,%s) failed")

		assert.Equal(t, expected, res)
	})

	t.Run("format contains 日本語 part 2", func(t *testing.T) {
		format := "み%sっ%sのこ"
		str := "みかんとずっきーにときのこ"
		expected1 := "かんとず"
		expected2 := "きーにとき"

		var res1, res2 string

		err := goparse.Parse(format, str).Insert(&res1, &res2)
		assert.NoErrorf(t, err, "Parse(%s,%s) failed", format, str)

		assert.Equal(t, expected1, res1)
		assert.Equal(t, expected2, res2)
	})

	t.Run("format contains 日本語 part 3", func(t *testing.T) {
		format := "水樹素子「%s」。秋穂伊織「%s」"
		str := "水樹素子「今日は天気が悪いね」。秋穂伊織「そうだね」"
		expected1 := "今日は天気が悪いね"
		expected2 := "そうだね"
		var res1, res2 string

		err := goparse.Parse(format, str).Insert(&res1, &res2)
		assert.NoErrorf(t, err, "Parse(%s,%s) failed", format, str)

		assert.Equal(t, expected1, res1)
		assert.Equal(t, expected2, res2)
	})

	asI := func(ss []string) []interface{} {
		res := make([]interface{}, len(ss))
		for i, str := range ss {
			res[i] = str
		}
		return res
	}

	t.Run("text contains many %s", func(t *testing.T) {
		t.SkipNow()
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

		}
	})

	t.Run("Invalid argumetns", func(t *testing.T) {
		t.Run("No match", func(t *testing.T) {
			format := "Hello"
			str := "noHello"
			var res string
			err := goparse.Parse(format, str).Insert(&res)
			assert.Errorf(t, err, "Parse(%s,%s) not failed want fail")
		})

		t.Run("a cuple of %s", func(t *testing.T) {
			format := "%s%s%s"
			str := "Hello"
			var res string
			err := goparse.Parse(format, str).Insert(&res)
			assert.Errorf(t, err, "Parse(%s,%s) not failed want fail")
		})

		t.Run("different number of %s from str", func(t *testing.T) {
			format := "%s_%s_%s"
			str := "H_He"
			var res string
			err := goparse.Parse(format, str).Insert(&res)
			assert.Errorf(t, err, "Parse(%s,%s) not failed want fail")
		})
	})
}

func TestParse_integer(t *testing.T) {

	t.Run("The opposite of Sprintf", func(t *testing.T) {
		format := "Hello my number is %d"
		expected := 100
		var res int
		err := goparse.Parse(format, fmt.Sprintf(format, expected)).Insert(&res)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("the argument is int friends?", func(t *testing.T) {

		t.Run("the argument is int8", func(t *testing.T) {
			format := "-%d-"
			str := "-123-"
			expected := int8(123)
			checkTestCase(t, str, format, expected)

			var res int8
			err := goparse.Parse(format, str).Insert(&res)
			assert.NoErrorf(t, err, "Parse(%s,%s) failed")

			assert.Equal(t, expected, res)
		})

		t.Run("the argument is int32", func(t *testing.T) {
			format := "-%d-"
			str := "-123-"
			expected := int32(123)
			checkTestCase(t, str, format, expected)

			var res int32
			err := goparse.Parse(format, str).Insert(&res)
			assert.NoErrorf(t, err, "Parse(%s,%s) failed")

			assert.Equal(t, expected, res)
		})

		t.Run("the argument is int64", func(t *testing.T) {
			format := "-%d-"
			str := "-123-"
			expected := int64(123)
			checkTestCase(t, str, format, expected)

			var res int64
			err := goparse.Parse(format, str).Insert(&res)
			assert.NoErrorf(t, err, "Parse(%s,%s) failed")

			assert.Equal(t, expected, res)
		})

	})

	t.Run("text contains multiple %d", func(t *testing.T) {
		format := "1%d456%d89"
		str := "123456789"
		expected1 := 23
		expected2 := 7
		checkTestCase(t, str, format, expected1, expected2)

		var res1, res2 int
		err := goparse.Parse(format, str).Insert(&res1, &res2)
		assert.NoErrorf(t, err, "Parse(%s,%s) failed")

		assert.Equal(t, expected1, res1)
		assert.Equal(t, expected2, res2)
	})

	t.Run("invalid argument", func(t *testing.T) {
		t.Run("not number", func(t *testing.T) {
			format := "%d"
			str := "ss"
			var res int
			err := goparse.Parse(format, str).Insert(&res)
			assert.Errorf(t, err, "Parse(%s,%s) not failed want fail", format, str)
		})

		t.Run("empty format", func(t *testing.T) {
			format := "%d"
			str := ""
			var res int
			err := goparse.Parse(format, str).Insert(&res)
			assert.Errorf(t, err, "Parse(%s,%s) not failed want fail", format, str)
		})

	})

}

func TestParse_integer_base8(t *testing.T) {

	t.Run("The opposite of Sprintf", func(t *testing.T) {
		t.Run("case1", func(t *testing.T) {
			format := "Hello my number is %o"
			expected := 123
			var res int
			err := goparse.Parse(format, fmt.Sprintf(format, expected)).Insert(&res)
			assert.NoError(t, err)
			assert.Equal(t, expected, res)
		})

		t.Run("case2 multiple", func(t *testing.T) {
			format := "%o %o %o"
			expected1 := 123
			expected2 := 456
			expected3 := 135
			var res1, res2, res3 int
			err := goparse.Parse(format, fmt.Sprintf(format, expected1, expected2, expected3)).
				Insert(&res1, &res2, &res3)
			assert.NoError(t, err)
			assert.Equal(t, expected1, res1)
			assert.Equal(t, expected2, res2)
			assert.Equal(t, expected3, res3)
		})

	})

	t.Run("invalid argument", func(t *testing.T) {
		t.Run("not number", func(t *testing.T) {
			format := "%o"
			str := "ss"
			var res int

			err := goparse.Parse(format, str).Insert(&res)
			assert.Errorf(t, err, "Parse(%s,%s) not failed want fail")
		})

		t.Run("empty format", func(t *testing.T) {
			format := "%o"
			str := ""
			var res int

			err := goparse.Parse(format, str).Insert(&res)
			assert.Errorf(t, err, "Parse(%s,%s) not failed want fail")
		})

	})

}

func TestParse_boolean(t *testing.T) {

	t.Run("The opposite of Sprintf", func(t *testing.T) {
		format := "Hello my number is %t"
		expected := true
		var res bool

		err := goparse.Parse(format, fmt.Sprintf(format, expected)).Insert(&res)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("text contains multiple %d", func(t *testing.T) {
		format := "Ah%t_Oh%t_Uu%t"
		str := "Ahtrue_Ohfalse_Uutrue"
		checkTestCase(t, str, format, true, false, true)
		var res1, res2, res3 bool

		err := goparse.Parse(format, str).Insert(&res1, &res2, &res3)
		assert.NoErrorf(t, err, "Parse(%s,%s) failed")

		assert.Equal(t, true, res1)
		assert.Equal(t, false, res2)
		assert.Equal(t, true, res3)
	})

	t.Run("invalid argument", func(t *testing.T) {
		t.Run("not boolean", func(t *testing.T) {
			format := "%t"
			str := "ss"
			var res bool
			err := goparse.Parse(format, str).Insert(&res)
			assert.Errorf(t, err, "Parse(%s,%s) not failed want fail")
		})

		t.Run("empty format", func(t *testing.T) {
			format := "%t"
			str := ""
			var res bool

			err := goparse.Parse(format, str).Insert(&res)
			assert.Errorf(t, err, "Parse(%s,%s) not failed want fail")
		})

	})

}

func ExampleParse() {
	var str string
	_ = goparse.Parse("Hello %s", "Hello World").Insert(&str)
	fmt.Println(str)
	// Output:
	// World
}

func ExampleParse_ja() {
	format := "水樹素子「%s」。秋穂伊織「%s」"
	str := "水樹素子「今日は天気が悪いね」。秋穂伊織「そうだね」"
	var mizukiMsg, ioriMsg string
	_ = goparse.Parse(format, str).Insert(&mizukiMsg, &ioriMsg)
	fmt.Println(mizukiMsg)
	fmt.Println(ioriMsg)
	// Output:
	// 今日は天気が悪いね
	// そうだね
}

func ExampleParse_number() {
	format := "Room %d"
	str := "Room 101"
	var num int
	_ = goparse.Parse(format, str).Insert(&num)
	fmt.Println(num)
	// Output:
	// 101
}

func ExampleParse_boolean() {
	format := "I can't tell whether it is %t or %t"
	str := "I can't tell whether it is false or true"
	var res1, res2 bool
	_ = goparse.Parse(format, str).Insert(&res1, &res2)
	fmt.Println(res1)
	fmt.Println(res2)
	// Output:
	// false
	// true
}

func ExampleParse_ja_number() {
	format := "塩ラーメン ￥%d円"
	str := "塩ラーメン ￥409円"
	var num int
	_ = goparse.Parse(format, str).Insert(&num)
	fmt.Println(num)
	// Output:
	// 409
}

func ExampleParse_ja_number8() {
	format := "Hello my number is %o"
	expected := 123

	var num int
	fmt.Println(fmt.Sprintf(format, expected))
	_ = goparse.Parse(format, fmt.Sprintf(format, expected)).Insert(&num)
	fmt.Println(num)
	// Output:
	// Hello my number is 173
	// 123
}

func ExampleParse_failed_insert_int_to_string() {
	format := "Hello!! my number is %d"
	str := "Hello!! my number is 1"

	var resInvalidType string
	err := goparse.Parse(format, str).Insert(&resInvalidType)
	fmt.Println(err.Error())
	// Output:
	// assign(src{kind:int,1} => dest[0]) failed err:type mismatch: expected *int{8,32,64}, actual *string
}

func ExampleParse_failed() {
	format := "Hello!! my number is %d"
	str := "Hello!! my number is One"

	var resInvalidType int
	err := goparse.Parse(format, str).Insert(&resInvalidType)
	fmt.Println(err.Error())
	// Output:
	// parseInteger(%d,"One",10) failed: ParseInt(One,10) failed: strconv.ParseInt: parsing "One": invalid syntax
}
