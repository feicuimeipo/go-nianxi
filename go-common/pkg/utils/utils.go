package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"unsafe"
)

func Int2Byte(data int) (ret []byte) {
	var len uintptr = unsafe.Sizeof(data)
	ret = make([]byte, len)
	var tmp int = 0xff
	var index uint = 0
	for index = 0; index < uint(len); index++ {
		ret[index] = byte((tmp << (index * 8) & data) >> (index * 8))
	}
	return ret
}

func Byte2Int(data []byte) int {
	var ret int = 0
	var len int = len(data)
	var i uint = 0
	for i = 0; i < uint(len); i++ {
		ret = ret | (int(data[i]) << (i * 8))
	}
	return ret
}

func Substr(str string, start, length int) string {
	if len(str) == 0 {
		return str
	}
	rs := []rune(str)
	rl := len(rs)
	end := 0
	if start < 0 {
		start = rl - 1 + start
	}

	end = start + length
	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}

	if start > rl {
		start = rl
	}

	if end < 0 {
		end = 0
	}

	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}

// 按名字和类型复制（&a,b）b->a
func CopyFields(target interface{}, source interface{}, fields ...string) (err error) {
	at := reflect.TypeOf(target)
	av := reflect.ValueOf(target)
	bt := reflect.TypeOf(source)
	bv := reflect.ValueOf(source)
	// 简单判断下
	if at.Kind() != reflect.Ptr {
		err = fmt.Errorf("a must be a struct pointer")
		return
	}
	av = reflect.ValueOf(av.Interface())
	// 要复制哪些字段
	_fields := make([]string, 0)
	if len(fields) > 0 {
		_fields = fields
	} else {
		for i := 0; i < bv.NumField(); i++ {
			_fields = append(_fields, bt.Field(i).Name)
		}
	}
	if len(_fields) == 0 {
		fmt.Println("no fields to copy")
		return
	}
	// 复制
	for i := 0; i < len(_fields); i++ {
		name := _fields[i]
		f := av.Elem().FieldByName(name)
		bValue := bv.FieldByName(name)
		// a中有同名的字段并且类型一致才复制
		if f.IsValid() && f.Kind() == bValue.Kind() {
			f.Set(bValue)
		} else {
			fmt.Printf("no such field or different kind, fieldName: %s\n", name)
		}
	}
	return
}

func StringStrip(input string) string {
	if input == "" {
		return ""
	}
	reg := regexp.MustCompile(`[\s\p{Zs}]{1,}`)
	return reg.ReplaceAllString(input, "")
}
