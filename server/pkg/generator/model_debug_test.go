package generator

import (
	"strings"
	"testing"
)

// TestDebugModel 测试列名生成逻辑
func TestDebugModel(t *testing.T) {
	// 不同的字段名测试
	testCases := []struct {
		input    string
		expected string
	}{
		{"CategoryID", "category_id"},
		{"AuthorID", "author_id"},
		{"UserRoleID", "user_role_id"},
		{"ID", "id"},
		{"ProductCategoryID", "product_category_id"},
	}

	for _, tc := range testCases {
		// 修复算法，特殊处理ID后缀
		var columnName string

		// 检查是否以ID结尾，特殊处理
		if strings.HasSuffix(tc.input, "ID") && tc.input != "ID" {
			// 移除ID后缀
			base := tc.input[:len(tc.input)-2]

			// 处理基础部分
			for i, c := range base {
				if i > 0 && c >= 'A' && c <= 'Z' {
					columnName += "_"
				}
				if c >= 'A' && c <= 'Z' {
					columnName += string(c - 'A' + 'a')
				} else {
					columnName += string(c)
				}
			}

			// 添加_id后缀
			columnName += "_id"
		} else if tc.input == "ID" {
			// ID直接转换为id
			columnName = "id"
		} else {
			// 常规处理
			for i, c := range tc.input {
				if i > 0 && c >= 'A' && c <= 'Z' {
					columnName += "_"
				}
				if c >= 'A' && c <= 'Z' {
					columnName += string(c - 'A' + 'a')
				} else {
					columnName += string(c)
				}
			}
		}

		t.Logf("输入: %s, 输出: %s, 期望: %s", tc.input, columnName, tc.expected)
		if columnName != tc.expected {
			t.Errorf("列名生成错误，输入: %s, 期望: %s, 实际: %s", tc.input, tc.expected, columnName)
		}
	}
}
