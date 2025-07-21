package generator

import (
	"testing"
)

// DebugModel 用于调试模型生成
func DebugModel(t *testing.T) {
	// 模拟外键字段名生成
	relationFieldName := "Category"
	foreignKeyFieldName := relationFieldName + "ID" // 默认外键字段名，如"CategoryID"

	// 手动生成正确格式的外键列名（下划线分隔的小写形式）
	var foreignKeyColumnName string
	for i, c := range foreignKeyFieldName {
		if i > 0 && c >= 'A' && c <= 'Z' {
			foreignKeyColumnName += "_"
		}
		if c >= 'A' && c <= 'Z' {
			foreignKeyColumnName += string(c - 'A' + 'a')
		} else {
			foreignKeyColumnName += string(c)
		}
	}

	// 打印结果
	t.Logf("关系字段名: %s", relationFieldName)
	t.Logf("外键字段名: %s", foreignKeyFieldName)
	t.Logf("外键列名: %s", foreignKeyColumnName)

	// 检查结果
	if foreignKeyColumnName != "category_id" {
		t.Errorf("外键列名生成错误，期望: category_id, 实际: %s", foreignKeyColumnName)
	}
}
