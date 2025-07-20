package response

import "github.com/jinzhu/copier"

// Copy 拷贝结构体
func Copy(toValue interface{}, fromValue interface{}) interface{} {
	if err := copier.Copy(toValue, fromValue); err != nil {
		panic(SystemError)
	}
	return toValue
}
