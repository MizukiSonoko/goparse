// Copyright (C) 2018,2019 MizukiSonoko. All rights reserved.

package goparse_test

import (
	"fmt"
	"log"
	"math"
	"strings"
	"testing"

	goparse "github.com/MizukiSonoko/goparse/parse"
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

	t.Run("format contains an unsupported type", func(t *testing.T) {
		format := "Hello I want a coffee %g gram"
		str := "Hello I want a coffee 123.456 gram"
		var res float32
		err := goparse.Parse(format, str).Insert(&res)
		assert.Error(t, err)
	})

	t.Run("format contains %s, but the argument is struct", func(t *testing.T) {
		format := "Hello %s"
		str := "Hello Iori"
		var res struct {
			dummy string
		}
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

	t.Run("invalid string like %s%s%s", func(t *testing.T) {
		format := "%s%s%s"
		str := "abc"
		var res string

		err := goparse.Parse(format, str).Insert(&res)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ambiguous")
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

}

func TestParse_string_ja(t *testing.T) {

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

}

func TestParse_string_invalid(t *testing.T) {

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

	t.Run("Invalid target pointer (nil)", func(t *testing.T) {
		format := "Hello"
		str := "noHello"
		var res *string
		err := goparse.Parse(format, str).Insert(res)
		assert.Errorf(t, err, "Parse(%s,%s) not failed want fail")
	})

	t.Run("Invalid target pointer (func)", func(t *testing.T) {
		format := "Hello"
		str := "noHello"
		res := func() {}
		err := goparse.Parse(format, str).Insert(&res)
		assert.Errorf(t, err, "Parse(%s,%s) not failed want fail")
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

		t.Run("the argument is int8, and str happens overflow", func(t *testing.T) {
			format := "-%d-"
			str := fmt.Sprintf("-%d-", math.MaxInt8+1)

			var res int8
			err := goparse.Parse(format, str).Insert(&res)
			assert.Errorf(t, err, "Parse(%s,%s) failed")
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

		t.Run("the argument is int32, and str happens overflow", func(t *testing.T) {
			format := "-%d-"
			str := fmt.Sprintf("-%d-", math.MaxInt32+1)

			var res int32
			err := goparse.Parse(format, str).Insert(&res)
			assert.Errorf(t, err, "Parse(%s,%s) failed")
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

		t.Run("the argument is int64, and str happens overflow", func(t *testing.T) {
			format := "-%d-"
			// 9223372036854775801=(2^63-1)+1
			str := fmt.Sprintf("-9223372036854775808-")

			var res int64
			err := goparse.Parse(format, str).Insert(&res)
			print(err.Error())
			assert.Errorf(t, err, "Parse(%s,%s) failed")
		})

	})

	t.Run("case single", func(t *testing.T) {
		format := "%d"
		expected1 := 123
		var res1 int
		err := goparse.Parse(format, fmt.Sprintf(format, expected1)).
			Insert(&res1)
		assert.NoError(t, err)
		assert.Equal(t, expected1, res1)
	})

	t.Run("case2 multiple 2", func(t *testing.T) {
		format := "%d %d"
		expected1 := 123
		expected2 := 456
		var res1, res2 int
		err := goparse.Parse(format, fmt.Sprintf(format, expected1, expected2)).
			Insert(&res1, &res2)
		assert.NoError(t, err)
		assert.Equal(t, expected1, res1)
		assert.Equal(t, expected2, res2)
	})

	t.Run("case2 multiple 3", func(t *testing.T) {
		format := "%d %d %d"
		expected1 := 123
		expected2 := 456
		expected3 := 985
		var res1, res2, res3 int
		err := goparse.Parse(format, fmt.Sprintf(format, expected1, expected2, expected3)).
			Insert(&res1, &res2, &res3)
		assert.NoError(t, err)
		assert.Equal(t, expected1, res1)
		assert.Equal(t, expected2, res2)
		assert.Equal(t, expected3, res3)
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

	t.Run("invalid int like %d%d", func(t *testing.T) {
		format := "%d%d"
		str := "123"
		var res int

		err := goparse.Parse(format, str).Insert(&res)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ambiguous")
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

	t.Run("case single", func(t *testing.T) {
		format := "%o"
		expected1 := 123
		var res1 int
		err := goparse.Parse(format, fmt.Sprintf(format, expected1)).
			Insert(&res1)
		assert.NoError(t, err)
		assert.Equal(t, expected1, res1)
	})

	t.Run("case multiple 2", func(t *testing.T) {
		format := "%o %o"
		expected1 := 123
		expected2 := 456
		var res1, res2 int
		err := goparse.Parse(format, fmt.Sprintf(format, expected1, expected2)).
			Insert(&res1, &res2)
		assert.NoError(t, err)
		assert.Equal(t, expected1, res1)
		assert.Equal(t, expected2, res2)
	})

	t.Run("case multiple 3", func(t *testing.T) {
		format := "%o %o %o"
		expected1 := 123
		expected2 := 456
		expected3 := 985
		var res1, res2, res3 int
		err := goparse.Parse(format, fmt.Sprintf(format, expected1, expected2, expected3)).
			Insert(&res1, &res2, &res3)
		assert.NoError(t, err)
		assert.Equal(t, expected1, res1)
		assert.Equal(t, expected2, res2)
		assert.Equal(t, expected3, res3)
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

func TestParse_integer_base2(t *testing.T) {

	t.Run("The opposite of Sprintf", func(t *testing.T) {
		t.Run("case1", func(t *testing.T) {
			format := "Hello my number is %b"
			expected := 123
			var res int
			err := goparse.Parse(format, fmt.Sprintf(format, expected)).Insert(&res)
			assert.NoError(t, err)
			assert.Equal(t, expected, res)
		})

		t.Run("case2 single", func(t *testing.T) {
			format := "%b"
			expected1 := 123
			var res1 int
			err := goparse.Parse(format, fmt.Sprintf(format, expected1)).
				Insert(&res1)
			assert.NoError(t, err)
			assert.Equal(t, expected1, res1)
		})

		t.Run("case2 multiple 2", func(t *testing.T) {
			format := "%b %b"
			expected1 := 123
			expected2 := 456
			var res1, res2 int
			err := goparse.Parse(format, fmt.Sprintf(format, expected1, expected2)).
				Insert(&res1, &res2)
			assert.NoError(t, err)
			assert.Equal(t, expected1, res1)
			assert.Equal(t, expected2, res2)
		})

		t.Run("case2 multiple", func(t *testing.T) {
			format := "%b %b %b"
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
			format := "%b"
			str := "ss"
			var res int

			err := goparse.Parse(format, str).Insert(&res)
			assert.Errorf(t, err, "Parse(%s,%s) not failed want fail")
		})

		t.Run("empty format", func(t *testing.T) {
			format := "%b"
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

			if !strings.Contains(err.Error(), "parseBool") {
				t.Errorf("Parse error should contain %s, but it's %s",
					"parseBool", err.Error())

			}
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

func TestParse_float(t *testing.T) {

	t.Run("The opposite of Sprintf", func(t *testing.T) {
		format := "Hello my number is %f"
		expected := 123.456
		var res float64

		err := goparse.Parse(format, fmt.Sprintf(format, expected)).Insert(&res)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("The opposite of Sprintf (float32)", func(t *testing.T) {
		format := "Hello my number is %f"
		expected := float32(123.456)
		var res float32

		err := goparse.Parse(format, fmt.Sprintf(format, expected)).Insert(&res)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("text contains multiple %d", func(t *testing.T) {
		format := "Wow%f_Wow%f_Wow%f"
		str := "Wow12.340000_Wow45.670000_Wow78.900000"
		checkTestCase(t, str, format, 12.34, 45.67, 78.9)
		expected1 := 12.34
		expected2 := 45.67
		expected3 := 78.9
		var res1, res2, res3 float64

		err := goparse.Parse(format, str).Insert(&res1, &res2, &res3)
		assert.NoErrorf(t, err, "Parse(%s,%s) failed")

		assert.Equal(t, expected1, res1)
		assert.Equal(t, expected2, res2)
		assert.Equal(t, expected3, res3)
	})

	t.Run("invalid argument", func(t *testing.T) {
		t.Run("not float", func(t *testing.T) {
			format := "%f"
			str := "ss"
			var res float64
			err := goparse.Parse(format, str).Insert(&res)
			assert.Errorf(t, err, "Parse(%s,%s) not failed want fail")
		})

		t.Run("empty format", func(t *testing.T) {
			format := "%f"
			str := ""
			var res float64

			err := goparse.Parse(format, str).Insert(&res)
			assert.Errorf(t, err, "Parse(%s,%s) not failed want fail")
		})

	})

}

func TestParse_value_struct(t *testing.T) {

	t.Run("format contains %v", func(t *testing.T) {
		type sample struct {
			Name    string
			Value   int
			Boolean bool
			Point   float32
		}
		format := "sample %v"
		str := "sample {Hello 123 true 123.456}"
		var res sample
		err := goparse.Parse(format, str).Insert(&res)
		if err != nil {
			t.Fatalf("Parse failed err:%s", err)
		}
		if res.Name != "Hello" {
			t.Errorf("Name is %s, but want %s", res.Name, "Hello")
		}
		if res.Value != 123 {
			t.Errorf("Value is %d, but want %d", res.Value, 123)
		}
		if !res.Boolean {
			t.Errorf("Boolean is False, but want True")
		}
		if res.Point != 123.456 {
			t.Errorf("Point is %f, but want %f", res.Point, 123.456)
		}
	})

	t.Run("format contains %v, but struct has not exposed attribute", func(t *testing.T) {
		type sample struct {
			name    string
			value   int
			boolean bool
			point   float32
		}
		format := "sample %v"
		str := "sample {Hello 123}"
		var res sample
		err := goparse.Parse(format, str).Insert(&res)
		if err == nil {
			t.Fatalf("Parse returns nil err")
		}
		if !strings.Contains(err.Error(), "not exposed") {
			t.Errorf("err should contain %s, but it's %s",
				"not exposed", err.Error())
		}
	})

}

func TestParse_value_struct_invalid_struct_attribute(t *testing.T) {
	t.Run("format contains %v, but struct has invalid type int to bool", func(t *testing.T) {
		type sample struct {
			Value bool
		}
		format := "sample %v"
		str := "sample {123}"
		var res sample
		err := goparse.Parse(format, str).Insert(&res)
		if err == nil {
			t.Fatalf("Parse returns nil err")
		}
		if !strings.Contains(err.Error(), "invalid type") {
			t.Errorf("err should contain %s, but it's %s",
				"invalid type", err.Error())
		}
	})

	t.Run("format contains %v, but struct has invalid type string to bool", func(t *testing.T) {
		type sample struct {
			Value bool
		}
		format := "sample %v"
		str := "sample {Hello}"
		var res sample
		err := goparse.Parse(format, str).Insert(&res)
		if err == nil {
			t.Fatalf("Parse returns nil err")
		}
		if !strings.Contains(err.Error(), "invalid type") {
			t.Errorf("err should contain %s, but it's %s",
				"invalid type", err.Error())
		}
	})

	t.Run("format contains %v, but struct has invalid type string to float", func(t *testing.T) {
		type sample struct {
			Value float32
		}
		format := "sample %v"
		str := "sample {Hello}"
		var res sample
		err := goparse.Parse(format, str).Insert(&res)
		if err == nil {
			t.Fatalf("Parse returns nil err")
		}
		if !strings.Contains(err.Error(), "invalid type") {
			t.Errorf("err should contain %s, but it's %s",
				"invalid type", err.Error())
		}
	})

	t.Run("format contains %v, but struct has invalid type float to bool", func(t *testing.T) {
		type sample struct {
			Value bool
		}
		format := "sample %v"
		str := "sample {123.45}"
		var res sample
		err := goparse.Parse(format, str).Insert(&res)
		if err == nil {
			t.Fatalf("Parse returns nil err")
		}
		if !strings.Contains(err.Error(), "invalid type") {
			t.Errorf("err should contain %s, but it's %s",
				"invalid type", err.Error())
		}
	})

	t.Run("format contains %v, but struct has invalid type bool to int", func(t *testing.T) {
		type sample struct {
			Value int
		}
		format := "sample %v"
		str := "sample {false}"
		var res sample
		err := goparse.Parse(format, str).Insert(&res)
		if err == nil {
			t.Fatalf("Parse returns nil err")
		}
		if !strings.Contains(err.Error(), "invalid type") {
			t.Errorf("err should contain %s, but it's %s",
				"invalid type", err.Error())
		}
	})

}

func TestParse_value_primitive(t *testing.T) {

	t.Run("format contains %v, it has string", func(t *testing.T) {
		format := "sample %v"
		str := "sample Hello"
		var res string
		err := goparse.Parse(format, str).Insert(&res)
		if err != nil {
			t.Fatalf("Parse failed err:%s", err)
		}
		if res != "Hello" {
			t.Errorf("res is %s, but wants %s",
				res, "Hello")
		}
	})

	t.Run("format contains %v, it has int", func(t *testing.T) {
		format := "sample %v"
		str := "sample 123"
		var res int
		err := goparse.Parse(format, str).Insert(&res)
		if err != nil {
			t.Fatalf("Parse failed err:%s", err)
		}
		if res != 123 {
			t.Errorf("res is %d, but wants %d",
				res, 123)
		}
	})

	t.Run("format contains %v, it has float", func(t *testing.T) {
		format := "sample %v"
		str := "sample 123.456"
		var res float64
		err := goparse.Parse(format, str).Insert(&res)
		if err != nil {
			t.Fatalf("Parse failed err:%s", err)
		}
		if res != 123.456 {
			t.Errorf("res is %f, but wants %f",
				res, 123.456)
		}
	})

	t.Run("format contains %v, it has bool", func(t *testing.T) {
		format := "sample %v"
		str := "sample true"
		var res bool
		err := goparse.Parse(format, str).Insert(&res)
		if err != nil {
			t.Fatalf("Parse failed err:%s", err)
		}
		if !res {
			t.Errorf("res is false, but wants true")
		}
	})

}

func TestInsertOnly_normal(t *testing.T) {

	t.Run("working case", func(t *testing.T) {
		format := "Hello %s! I'm %s,%s"
		actual := "Hello Yukari! I'm SonokoMizuki, my favorite language is golang"
		var res string

		t.Run("First one", func(t *testing.T) {
			expected := "Yukari"
			err := goparse.Parse(format, actual).InsertOnly(0, &res)
			if err != nil {
				log.Fatalf("InsertOnly failed err:%s", err)
			}
			assert.Equal(t, expected, res)
		})

		t.Run("Second one", func(t *testing.T) {
			expected := "SonokoMizuki"
			err := goparse.Parse(format, actual).InsertOnly(1, &res)
			if err != nil {
				log.Fatalf("InsertOnly failed err:%s", err)
			}
			assert.Equal(t, expected, res)
		})
	})

	t.Run("invalid index case", func(t *testing.T) {
		format := "Hello %s!"
		actual := "Hello Yukari!"
		var res string

		t.Run("over index", func(t *testing.T) {
			err := goparse.Parse(format, actual).InsertOnly(99999999, &res)
			assert.Error(t, err)
		})

		t.Run("not number case (format expects number, but string)", func(t *testing.T) {
			format := "%d"
			str := "string"
			var res int
			err := goparse.Parse(format, str).InsertOnly(0, &res)
			assert.Errorf(t, err, "Parse(%s,%s) not failed want fail", format, str)
		})

		t.Run("type mismatch case (func expects target arguments is string, but int)", func(t *testing.T) {
			format := "%s"
			str := "string"
			var res int
			err := goparse.Parse(format, str).InsertOnly(0, &res)
			assert.Errorf(t, err, "Parse(%s,%s) not failed want fail", format, str)
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

func ExampleParse_insertOnly() {
	var greeting, name string
	result := goparse.Parse("%s, I'm %s.%s", "Hello, I'm MizukiSonoko.")
	result.InsertOnly(0, &greeting)
	result.InsertOnly(1, &name)
	fmt.Println(greeting)
	fmt.Println(name)
	// Output:
	// Hello
	// MizukiSonoko
}

func ExampleParse_struct() {
	type sample struct {
		Name  string
		Value int
	}
	format := "sample %v"
	str := "sample {Hello 123}"
	var res sample
	_ = goparse.Parse(format, str).Insert(&res)
	fmt.Println(res.Name)
	fmt.Println(res.Value)
	// Output:
	// Hello
	// 123
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

func ExampleParse_jaNumber() {
	format := "塩ラーメン ￥%d円"
	str := "塩ラーメン ￥409円"
	var num int
	_ = goparse.Parse(format, str).Insert(&num)
	fmt.Println(num)
	// Output:
	// 409
}

func ExampleParse_numberBase8() {
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

func ExampleParse_failedInsertIntToString() {
	format := "Hello!! my number is %d"
	str := "Hello!! my number is 1"

	var resInvalidType string
	err := goparse.Parse(format, str).Insert(&resInvalidType)
	fmt.Println(err.Error())
	// Output:
	// assign(src{kind:int,1} => dest[0]) failed err:type mismatch: expected *int{8,32,64}, actual *string
}

func ExampleParse_failedAmbigurousFormat() {
	format := "%s%s%s"
	str := "abc"
	var res string

	err := goparse.Parse(format, str).Insert(&res)
	fmt.Println(err.Error())
	// Output:
	// invalid format("%s%s%s"). too ambiguous to invese format
}

func ExampleParse_failed() {
	format := "Hello!! my number is %d"
	str := "Hello!! my number is One"

	var resInvalidType int
	err := goparse.Parse(format, str).Insert(&resInvalidType)
	fmt.Println(err.Error())
	// Output:
	// parseInteger(%d,"One",10) failed: ParseInt("One",10) failed: strconv.ParseInt: parsing "One": invalid syntax
}
