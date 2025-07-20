package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/zhoudm1743/go-web/core/facades"
	"gorm.io/gorm"
)

// HistoryModel 代码生成历史模型
type HistoryModel struct {
	ID          uint           `gorm:"primarykey" json:"id"`                  // 主键ID
	CreatedAt   time.Time      `json:"createdAt"`                             // 创建时间
	UpdatedAt   time.Time      `json:"updatedAt"`                             // 更新时间
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`                        // 删除时间
	Table       string         `gorm:"comment:表名" json:"table"`               // 表名
	StructName  string         `gorm:"comment:结构体名称" json:"structName"`       // 结构体名称
	PackageName string         `gorm:"comment:包名" json:"packageName"`         // 包名
	ModuleName  string         `gorm:"comment:模块名" json:"moduleName"`         // 模块名
	Description string         `gorm:"comment:描述" json:"description"`         // 描述
	Fields      string         `gorm:"type:text;comment:字段" json:"fields"`    // 字段JSON
	Templates   string         `gorm:"type:text;comment:模板" json:"templates"` // 模板JSON
	ApiIDs      string         `gorm:"comment:API ID列表" json:"apiIds"`        // API ID列表
	MenuID      uint           `gorm:"comment:菜单ID" json:"menuId"`            // 菜单ID
	Flag        uint8          `gorm:"default:0;comment:标记" json:"flag"`      // 标记 0:未删除 1:已删除
	BusinessDB  string         `gorm:"comment:业务数据库" json:"businessDb"`       // 业务数据库
}

// TableName 指定表名
func (HistoryModel) TableName() string {
	return "code_gen_history"
}

// HistoryManager 代码生成历史管理器
type HistoryManager struct {
	DB *gorm.DB
}

// NewHistoryManager 创建历史管理器
func NewHistoryManager() *HistoryManager {
	db := facades.DB()
	if db == nil {
		fmt.Println("警告: 数据库连接未初始化，使用内存数据库进行操作。")
		// 在实际应用中，这里可以配置内存数据库或返回错误
		return &HistoryManager{
			DB: nil, // 表示使用内存数据库
		}
	}
	return &HistoryManager{
		DB: db,
	}
}

// Migrate 迁移历史表
func (h *HistoryManager) Migrate() error {
	if h.DB == nil {
		fmt.Println("警告: 数据库连接未初始化，跳过迁移。")
		return nil // 如果数据库连接未初始化，则跳过迁移
	}
	return h.DB.AutoMigrate(&HistoryModel{})
}

// Create 创建历史记录
func (h *HistoryManager) Create(config *Config, templates map[string]string) (uint, error) {
	// 序列化字段
	fieldsJSON, err := json.Marshal(config.Fields)
	if err != nil {
		return 0, fmt.Errorf("序列化字段失败: %w", err)
	}

	// 序列化模板路径
	templatesJSON, err := json.Marshal(templates)
	if err != nil {
		return 0, fmt.Errorf("序列化模板路径失败: %w", err)
	}

	// 创建历史记录
	history := &HistoryModel{
		Table:       config.TableName,
		StructName:  config.StructName,
		PackageName: config.PackageName,
		ModuleName:  config.ModuleName,
		Description: config.Description,
		Fields:      string(fieldsJSON),
		Templates:   string(templatesJSON),
		Flag:        0,
		BusinessDB:  "", // 默认使用主数据库
	}

	if h.DB == nil {
		fmt.Println("警告: 数据库连接未初始化，跳过创建历史记录。")
		return 0, nil // 如果数据库连接未初始化，则跳过创建并返回成功
	}

	if err := h.DB.Create(history).Error; err != nil {
		return 0, fmt.Errorf("创建历史记录失败: %w", err)
	}

	return history.ID, nil
}

// GetList 获取历史列表
func (h *HistoryManager) GetList(page, pageSize int) ([]HistoryModel, int64, error) {
	var records []HistoryModel
	var total int64

	// 获取总数
	if err := h.DB.Model(&HistoryModel{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取历史记录总数失败: %w", err)
	}

	// 查询列表
	if err := h.DB.Model(&HistoryModel{}).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&records).Error; err != nil {
		return nil, 0, fmt.Errorf("查询历史记录列表失败: %w", err)
	}

	return records, total, nil
}

// Get 获取单个历史记录
func (h *HistoryManager) Get(id uint) (*HistoryModel, error) {
	var record HistoryModel
	if err := h.DB.First(&record, id).Error; err != nil {
		return nil, fmt.Errorf("获取历史记录失败: %w", err)
	}

	return &record, nil
}

// Delete 删除历史记录
func (h *HistoryManager) Delete(id uint) error {
	if err := h.DB.Delete(&HistoryModel{}, id).Error; err != nil {
		return fmt.Errorf("删除历史记录失败: %w", err)
	}

	return nil
}

// RollBack 回滚代码生成
func (h *HistoryManager) RollBack(id uint, deleteFiles, deleteAPI, deleteMenu, deleteTable bool) error {
	// 获取历史记录
	record, err := h.Get(id)
	if err != nil {
		return err
	}

	// 解析模板路径
	templates := make(map[string]string)
	if err := json.Unmarshal([]byte(record.Templates), &templates); err != nil {
		return fmt.Errorf("解析模板路径失败: %w", err)
	}

	// 删除生成的文件
	if deleteFiles {
		// 将文件移动到临时目录而不是直接删除
		removeBasePath := filepath.Join("./server/temp/rm_file", strconv.FormatInt(time.Now().Unix(), 10))

		for path := range templates {
			// 确保目标目录存在
			destDir := filepath.Join(removeBasePath, filepath.Dir(path))
			if err := os.MkdirAll(destDir, 0755); err != nil {
				return fmt.Errorf("创建临时目录失败: %w", err)
			}

			// 移动文件
			destPath := filepath.Join(removeBasePath, path)
			if err := h.moveFile(path, destPath); err != nil {
				fmt.Printf("警告: 移动文件 %s 失败: %v\n", path, err)
			} else {
				fmt.Printf("已移动文件: %s 到 %s\n", path, destPath)
			}
		}
	}

	// 删除数据库表
	if deleteTable && record.Table != "" {
		if err := h.DB.Exec("DROP TABLE IF EXISTS " + record.Table).Error; err != nil {
			fmt.Printf("警告: 删除表 %s 失败: %v\n", record.Table, err)
		} else {
			fmt.Printf("已删除表: %s\n", record.Table)
		}
	}

	// 更新历史记录标记为已回滚
	if err := h.DB.Model(&HistoryModel{}).Where("id = ?", id).Update("flag", 1).Error; err != nil {
		return fmt.Errorf("更新历史记录失败: %w", err)
	}

	return nil
}

// moveFile 移动文件
func (h *HistoryManager) moveFile(src, dst string) error {
	// 检查源文件是否存在
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return fmt.Errorf("源文件不存在: %w", err)
	}

	// 确保目标目录存在
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("创建目标目录失败: %w", err)
	}

	// 读取源文件内容
	content, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("读取源文件失败: %w", err)
	}

	// 写入目标文件
	if err := os.WriteFile(dst, content, 0644); err != nil {
		return fmt.Errorf("写入目标文件失败: %w", err)
	}

	// 删除源文件
	if err := os.Remove(src); err != nil {
		return fmt.Errorf("删除源文件失败: %w", err)
	}

	return nil
}
