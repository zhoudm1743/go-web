package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// generateFrontend 生成前端代码
func (g *Generator) generateFrontend() error {
	// 创建前端目录
	frontendDir := filepath.Join(g.RootPath, "front-end/src/views", strings.ToLower(g.Config.StructName))
	if err := os.MkdirAll(frontendDir, 0755); err != nil {
		return fmt.Errorf("创建前端目录失败: %w", err)
	}

	// 创建组件目录
	componentsDir := filepath.Join(frontendDir, "components")
	if err := os.MkdirAll(componentsDir, 0755); err != nil {
		return fmt.Errorf("创建组件目录失败: %w", err)
	}

	// 生成索引页面
	if err := g.generateIndexPage(frontendDir); err != nil {
		return err
	}

	// 生成表格模态框组件
	if err := g.generateTableModalComponent(componentsDir); err != nil {
		return err
	}

	// 生成API文件
	if err := g.generateApiFile(); err != nil {
		return err
	}

	return nil
}

// generateIndexPage 生成索引页面
func (g *Generator) generateIndexPage(dir string) error {
	const indexTemplate = `<template>
  <CommonWrapper>
    <template #title>
      {{ .Title }}
    </template>

    <template #buttons>
      <n-button type="primary" @click="handleAdd">
        <template #icon>
          <icon-ic-round-plus />
        </template>
        新增{{ .Description }}
      </n-button>
    </template>

    <template #content>
      <n-card :bordered="false">
        <n-data-table
          remote
          :loading="tableLoading"
          :columns="columns"
          :data="tableData"
          :pagination="pagination"
          @update:page="handlePageChange"
          @update:page-size="handlePageSizeChange"
        />
      </n-card>
    </template>

    <!-- 表格模态框组件 -->
    <TableModal
      ref="modalRef"
      :title="modalTitle"
      :loading="modalLoading"
      :mode="modalMode"
      @submit="handleModalSubmit"
    />
  </CommonWrapper>
</template>

<script lang="ts" setup>
import { reactive, ref } from 'vue';
import { useMessage } from 'naive-ui';
import { CommonWrapper } from '@/components/common';
import TableModal from './components/TableModal.vue';
import { create{{.StructName}}, delete{{.StructName}}, get{{.PluralName}}, update{{.StructName}} } from '@/service/api/{{.ApiFile}}';

// 表格设置
const tableLoading = ref(false);
const tableData = ref([]);
const pagination = reactive({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 30, 50],
  onChange: (page: number) => {
    pagination.page = page;
  },
  onUpdatePageSize: (pageSize: number) => {
    pagination.pageSize = pageSize;
    pagination.page = 1;
  },
});

// 模态框设置
const modalRef = ref();
const modalLoading = ref(false);
const modalTitle = ref('');
const modalMode = ref<'add' | 'edit'>('add');

// 消息组件
const message = useMessage();

// 常量
const title = '{{.Description}}管理';
const description = '{{.Description}}';

// 表格列
const columns = [
{{.ColumnDefs}}
  {
    title: '操作',
    key: 'actions',
    width: 150,
    fixed: 'right',
    render(row) {
      return [
        <n-button
          key="edit"
          type="primary"
          text
          size="small"
          onClick={() => handleEdit(row)}
        >
          编辑
        </n-button>,
        <n-button
          key="delete"
          type="error"
          text
          size="small"
          onClick={() => handleDelete(row)}
        >
          删除
        </n-button>
      ];
    },
  },
];

// 初始化加载数据
loadTableData();

// 加载表格数据
async function loadTableData() {
  try {
    tableLoading.value = true;
    const res = await get{{.PluralName}}({
      page: pagination.page,
      pageSize: pagination.pageSize
    });
    
    tableData.value = res.data.list;
    pagination.itemCount = res.data.total;
  } catch (error) {
    message.error('获取数据失败');
  } finally {
    tableLoading.value = false;
  }
}

// 处理页码变化
function handlePageChange(page: number) {
  pagination.page = page;
  loadTableData();
}

// 处理每页条数变化
function handlePageSizeChange(pageSize: number) {
  pagination.pageSize = pageSize;
  pagination.page = 1;
  loadTableData();
}

// 处理新增
function handleAdd() {
  modalTitle.value = '新增{{.Description}}';
  modalMode.value = 'add';
  modalRef.value.openModal();
}

// 处理编辑
function handleEdit(row) {
  modalTitle.value = '编辑{{.Description}}';
  modalMode.value = 'edit';
  modalRef.value.openModal(row);
}

// 处理删除
async function handleDelete(row) {
  try {
    await delete{{.StructName}}(row.id);
    message.success('删除成功');
    loadTableData();
  } catch (error) {
    message.error('删除失败');
  }
}

// 处理模态框提交
async function handleModalSubmit(formData) {
  try {
    modalLoading.value = true;
    
    if (modalMode.value === 'add') {
      await create{{.StructName}}(formData);
      message.success('创建成功');
    } else {
      await update{{.StructName}}(formData);
      message.success('更新成功');
    }
    
    modalRef.value.closeModal();
    loadTableData();
  } catch (error) {
    message.error(modalMode.value === 'add' ? '创建失败' : '更新失败');
  } finally {
    modalLoading.value = false;
  }
}
</script>
`

	// 准备模板数据
	type TemplateData struct {
		Title       string
		Description string
		StructName  string
		ApiFile     string
		PluralName  string
		ColumnDefs  string
	}

	// 生成列定义
	var columnDefs strings.Builder
	for _, field := range g.Config.Fields {
		if field.IsPrimaryKey {
			continue
		}

		jsonName := ToLowerCamel(field.FieldName)
		columnDefs.WriteString(fmt.Sprintf(`  {
    title: '%s',
    key: '%s',
    width: 120,
  },
`, field.FieldDesc, jsonName))
	}

	data := TemplateData{
		Title:       g.Config.Description + "管理",
		Description: g.Config.Description,
		StructName:  g.Config.StructName,
		ApiFile:     strings.ToLower(g.Config.StructName),
		PluralName:  ToPlural(g.Config.StructName),
		ColumnDefs:  columnDefs.String(),
	}

	// 解析和渲染模板
	t, err := template.New("indexPage").Parse(indexTemplate)
	if err != nil {
		return fmt.Errorf("解析索引页面模板失败: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return fmt.Errorf("渲染索引页面模板失败: %w", err)
	}

	// 写入文件
	indexFilePath := filepath.Join(dir, "index.vue")
	if err := os.WriteFile(indexFilePath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("写入索引页面文件失败: %w", err)
	}

	// 记录生成的文件
	g.AddGeneratedFile(indexFilePath, "frontend_index")

	fmt.Printf("生成索引页面文件: %s\n", indexFilePath)
	return nil
}

// generateTableModalComponent 生成表格模态框组件
func (g *Generator) generateTableModalComponent(dir string) error {
	const modalTemplate = `<template>
  <n-modal
    v-model:show="isShow"
    :title="title"
    preset="card"
    :mask-closable="false"
    @close="handleClose"
    class="modal-card"
  >
    <n-form
      ref="formRef"
      :model="formData"
      :rules="rules"
      label-placement="left"
      :label-width="80"
      :disabled="loading"
    >
{{.FormItems}}
    </n-form>
    <div class="modal-footer">
      <n-button class="mr-2" @click="closeModal" :disabled="loading">取消</n-button>
      <n-button type="primary" @click="handleSubmit" :loading="loading">确定</n-button>
    </div>
  </n-modal>
</template>

<script lang="ts" setup>
import { ref, reactive } from 'vue';
import type { FormInst, FormRules } from 'naive-ui';

// 表单实例
const formRef = ref<FormInst | null>(null);

// 基本状态
const isShow = ref(false);
const formData = reactive({
{{.FormDataFields}}
});

// 初始表单数据
const initialFormData = {
{{.FormDataFields}}
};

// 接收的Props
const props = defineProps({
  title: {
    type: String,
    default: '表单',
  },
  loading: {
    type: Boolean,
    default: false,
  },
  mode: {
    type: String as () => 'add' | 'edit',
    default: 'add',
  },
});

// 表单验证规则
const rules = {
{{.FormRules}}
} as FormRules;

// 事件
const emit = defineEmits(['submit']);

// 打开模态框
function openModal(data?: any) {
  // 重置表单
  resetForm();
  
  // 如果有数据，填充表单
  if (data) {
    Object.keys(formData).forEach((key) => {
      if (data[key] !== undefined) {
        formData[key] = data[key];
      }
    });
  }
  
  isShow.value = true;
}

// 关闭模态框
function closeModal() {
  isShow.value = false;
}

// 处理关闭
function handleClose() {
  resetForm();
}

// 重置表单
function resetForm() {
  if (formRef.value) {
    formRef.value.restoreValidation();
  }
  
  Object.keys(formData).forEach((key) => {
    formData[key] = initialFormData[key];
  });
}

// 处理提交
function handleSubmit() {
  formRef.value?.validate((errors) => {
    if (!errors) {
      emit('submit', { ...formData });
    }
  });
}

// 对外暴露的方法
defineExpose({
  openModal,
  closeModal,
});
</script>

<style scoped>
.modal-card {
  width: 550px;
}

.modal-footer {
  margin-top: 18px;
  text-align: right;
}
</style>
`

	// 准备表单项
	var formItems strings.Builder
	var formDataFields strings.Builder
	var formRules strings.Builder

	for _, field := range g.Config.Fields {
		if field.IsPrimaryKey {
			formDataFields.WriteString(fmt.Sprintf("  %s: 0,\n", ToLowerCamel(field.FieldName)))
			continue
		}

		jsonName := ToLowerCamel(field.FieldName)

		// 表单项
		formItems.WriteString(fmt.Sprintf(`      <n-form-item label="%s" path="%s">
        <n-input v-model:value="formData.%s" placeholder="请输入%s" />
      </n-form-item>
`, field.FieldDesc, jsonName, jsonName, field.FieldDesc))

		// 表单数据
		defaultValue := "null"
		if field.FieldType == "string" {
			defaultValue = "''"
		} else if field.FieldType == "bool" {
			defaultValue = "false"
		} else if field.FieldType == "uint" || field.FieldType == "int" {
			defaultValue = "0"
		}
		formDataFields.WriteString(fmt.Sprintf("  %s: %s,\n", jsonName, defaultValue))

		// 表单验证规则
		if field.Required {
			formRules.WriteString(fmt.Sprintf(`  %s: {
    required: true,
    message: '请输入%s',
    trigger: ['blur', 'input'],
  },
`, jsonName, field.FieldDesc))
		}
	}

	// 解析和渲染模板
	type TemplateData struct {
		FormItems      string
		FormDataFields string
		FormRules      string
	}

	data := TemplateData{
		FormItems:      formItems.String(),
		FormDataFields: formDataFields.String(),
		FormRules:      formRules.String(),
	}

	t, err := template.New("modalComponent").Parse(modalTemplate)
	if err != nil {
		return fmt.Errorf("解析模态框组件模板失败: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return fmt.Errorf("渲染模态框组件模板失败: %w", err)
	}

	// 写入文件
	modalFilePath := filepath.Join(dir, "TableModal.vue")
	if err := os.WriteFile(modalFilePath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("写入模态框组件文件失败: %w", err)
	}

	// 记录生成的文件
	g.AddGeneratedFile(modalFilePath, "frontend_component")

	fmt.Printf("生成模态框组件文件: %s\n", modalFilePath)
	return nil
}

// generateApiFile 生成API文件
func (g *Generator) generateApiFile() error {
	const apiTemplate = `import { http } from '../http';

// {{.Description}}列表查询参数
export interface {{.StructName}}QueryParams {
  page?: number;
  pageSize?: number;
{{.QueryParams}}
}

// {{.Description}}创建请求
export interface {{.StructName}}CreateRequest {
{{.CreateFields}}
}

// {{.Description}}更新请求
export interface {{.StructName}}UpdateRequest {
  id: number;
{{.UpdateFields}}
}

// {{.Description}}响应
export interface {{.StructName}}Response {
  id: number;
  createdAt: string;
  updatedAt: string;
{{.ResponseFields}}
}

// {{.Description}}列表响应
export interface {{.StructName}}ListResponse {
  total: number;
  list: {{.StructName}}Response[];
}

// 获取{{.Description}}列表
export const get{{.PluralName}} = (params: {{.StructName}}QueryParams) => {
  return http.request<{{.StructName}}ListResponse>({
    url: '/{{.ApiPrefix}}/list',
    method: 'GET',
    params,
  });
};

// 获取{{.Description}}详情
export const get{{.StructName}} = (id: number) => {
  return http.request<{{.StructName}}Response>({
    url: '/{{.ApiPrefix}}/detail/' + id,
    method: 'GET',
  });
};

// 创建{{.Description}}
export const create{{.StructName}} = (data: {{.StructName}}CreateRequest) => {
  return http.request<void>({
    url: '/{{.ApiPrefix}}/create',
    method: 'POST',
    data,
  });
};

// 更新{{.Description}}
export const update{{.StructName}} = (data: {{.StructName}}UpdateRequest) => {
  return http.request<void>({
    url: '/{{.ApiPrefix}}/update',
    method: 'PUT',
    data,
  });
};

// 删除{{.Description}}
export const delete{{.StructName}} = (id: number) => {
  return http.request<void>({
    url: '/{{.ApiPrefix}}/delete/' + id,
    method: 'DELETE',
  });
};
`

	// 准备模板数据
	var queryParams strings.Builder
	var createFields strings.Builder
	var updateFields strings.Builder
	var responseFields strings.Builder

	for _, field := range g.Config.Fields {
		if field.IsPrimaryKey {
			continue
		}

		jsonName := ToLowerCamel(field.FieldName)

		// TypeScript类型
		tsType := "string"
		if field.FieldType == "uint" || field.FieldType == "int" || field.FieldType == "int64" || field.FieldType == "uint64" {
			tsType = "number"
		} else if field.FieldType == "bool" {
			tsType = "boolean"
		}

		// 查询参数
		if field.IsSearchable || field.IsFilterable {
			queryParams.WriteString(fmt.Sprintf("  %s?: %s;\n", jsonName, tsType))
		}

		// 创建字段
		required := "?"
		if field.Required {
			required = ""
		}
		createFields.WriteString(fmt.Sprintf("  %s%s: %s;\n", jsonName, required, tsType))

		// 更新字段
		updateFields.WriteString(fmt.Sprintf("  %s?: %s;\n", jsonName, tsType))

		// 响应字段
		responseFields.WriteString(fmt.Sprintf("  %s: %s;\n", jsonName, tsType))
	}

	// 解析和渲染模板
	type TemplateData struct {
		Description    string
		StructName     string
		PluralName     string
		ApiPrefix      string
		QueryParams    string
		CreateFields   string
		UpdateFields   string
		ResponseFields string
	}

	data := TemplateData{
		Description:    g.Config.Description,
		StructName:     g.Config.StructName,
		PluralName:     ToPlural(g.Config.StructName),
		ApiPrefix:      g.Config.ApiPrefix,
		QueryParams:    queryParams.String(),
		CreateFields:   createFields.String(),
		UpdateFields:   updateFields.String(),
		ResponseFields: responseFields.String(),
	}

	t, err := template.New("apiFile").Parse(apiTemplate)
	if err != nil {
		return fmt.Errorf("解析API文件模板失败: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return fmt.Errorf("渲染API文件模板失败: %w", err)
	}

	// 确保目录存在
	apiDir := filepath.Join(g.RootPath, "front-end/src/service/api")
	if err := os.MkdirAll(apiDir, 0755); err != nil {
		return fmt.Errorf("创建API目录失败: %w", err)
	}

	// 写入文件
	apiFilePath := filepath.Join(apiDir, strings.ToLower(g.Config.StructName)+".ts")
	if err := os.WriteFile(apiFilePath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("写入API文件失败: %w", err)
	}

	// 记录生成的文件
	g.AddGeneratedFile(apiFilePath, "frontend_api")

	fmt.Printf("生成API文件: %s\n", apiFilePath)
	return nil
}
