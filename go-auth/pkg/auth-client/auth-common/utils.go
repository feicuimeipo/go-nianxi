package auth_common

import (
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"reflect"
)

// 结构体转为json
func Struct2Json(obj interface{}) string {
	str, err := json.Marshal(obj)
	if err != nil {
		panic(fmt.Sprintf("[Struct2Json]转换异常: %v", err))
	}
	return string(str)
}

// json转为结构体
func Json2Struct(str string, obj interface{}) {
	// 将json转为结构体
	err := json.Unmarshal([]byte(str), obj)
	if err != nil {
		panic(fmt.Sprintf("[Json2Struct]转换异常: %v", err))
	}
}

// json interface转为结构体
func JsonI2Struct(str interface{}, obj interface{}) {
	JsonStr := str.(string)
	Json2Struct(JsonStr, obj)
}

// CopyFields 用b的所有字段覆盖a的
// 如果fields不为空, 表示用b的特定字段覆盖a的
// a应该为结构体指针
func CopyFields(target interface{}, source interface{}, fields ...string) (err error) {
	at := reflect.TypeOf(target)
	av := reflect.ValueOf(target)
	bt := reflect.TypeOf(source)
	bv := reflect.ValueOf(source)

	// 简单判断下
	if at.Kind() != reflect.Ptr {
		err = errors.New("target 必须是指针")
		return
	}
	av = reflect.ValueOf(av.Interface())

	// 要复制哪些字段
	_fields := make([]string, 0)
	if len(fields) > 0 {
		_fields = fields
	} else {
		for i := 0; i < bv.NumField(); i++ {
			//if bt.Field(i).Name == "Id" {
			//    continue
			//}
			_fields = append(_fields, bt.Field(i).Name)
		}
	}

	if len(_fields) == 0 {
		err = errors.New("没有字段被copy")
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
			err = errors.New(fmt.Sprintf("%s字段类型不一致\n", name))
		}
	}
	return
}
