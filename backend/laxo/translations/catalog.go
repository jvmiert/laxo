// Code generated by running "go generate" in golang.org/x/text. DO NOT EDIT.

package translations

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
)

type dictionary struct {
	index []uint32
	data  string
}

func (d *dictionary) Lookup(key string) (data string, ok bool) {
	p, ok := messageKeyToIndex[key]
	if !ok {
		return "", false
	}
	start, end := d.index[p], d.index[p+1]
	if start == end {
		return "", false
	}
	return d.data[start:end], true
}

func init() {
	dict := map[string]catalog.Dictionary{
		"en": &dictionary{index: enIndex, data: enData},
		"vi": &dictionary{index: viIndex, data: viData},
	}
	fallback := language.MustParse("en")
	cat, err := catalog.NewFromMap(dict, catalog.Fallback(fallback))
	if err != nil {
		panic(err)
	}
	message.DefaultCatalog = cat
}

var messageKeyToIndex = map[string]int{
	"Email not found":                      0,
	"Password is incorrect":                1,
	"cannot be blank":                      2,
	"must be a valid email address":        4,
	"must contain a digit":                 5,
	"must contain a letter":                6,
	"the length must be between %v and %v": 3,
	"unknown validation error":             7,
	"user already exists":                  8,
}

var enIndex = []uint32{ // 10 elements
	0x00000000, 0x00000010, 0x00000026, 0x00000036,
	0x00000061, 0x0000007f, 0x00000094, 0x000000aa,
	0x000000c3, 0x000000d7,
} // Size: 64 bytes

const enData string = "" + // Size: 215 bytes
	"\x02Email not found\x02Password is incorrect\x02cannot be blank\x02the l" +
	"ength must be between %[1]v and %[2]v\x02must be a valid email address" +
	"\x02must contain a digit\x02must contain a letter\x02unknown validation " +
	"error\x02user already exists"

var viIndex = []uint32{ // 10 elements
	0x00000000, 0x00000019, 0x00000034, 0x00000034,
	0x00000034, 0x00000065, 0x00000065, 0x00000065,
	0x00000065, 0x00000065,
} // Size: 64 bytes

const viData string = "" + // Size: 101 bytes
	"\x02Không tìm thấy email\x02Mật khẩu không đúng\x02Phải la một địa ch" +
	"ỉ email hợp lệ"

	// Total table size 444 bytes (0KiB); checksum: 2A4B5670