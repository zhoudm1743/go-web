<template>
  <div>
    <n-card :bordered="false">
      <n-data-table
        remote
        :loading="historyLoading"
        :columns="historyColumns"
        :data="historyData"
        :pagination="historyPagination"
        @update:page="handleHistoryPageChange"
        @update:page-size="handleHistoryPageSizeChange"
      />
    </n-card>

    <!-- 回滚确认对话框 -->
    <n-modal
      v-model:show="showRollbackModal"
      preset="dialog"
      title="确认回滚"
      positive-text="确认"
      negative-text="取消"
      :loading="rollbackLoading"
      @positive-click="confirmRollback"
      @negative-click="cancelRollback"
    >
      <div class="py-4">
        <p class="mb-4">确定要回滚此次生成的代码吗？请选择回滚选项：</p>
        <div class="flex flex-col gap-2">
          <n-checkbox v-model:checked="rollbackOptions.deleteFiles">
            删除生成的文件
          </n-checkbox>
          <n-checkbox v-model:checked="rollbackOptions.deleteApi">
            删除生成的API
          </n-checkbox>
          <n-checkbox v-model:checked="rollbackOptions.deleteMenu">
            删除生成的菜单
          </n-checkbox>
          <n-checkbox v-model:checked="rollbackOptions.deleteTable">
            删除生成的数据表
          </n-checkbox>
        </div>
      </div>
    </n-modal>

    <!-- 文件预览对话框 -->
    <n-modal
      v-model:show="showPreviewModal"
      preset="card"
      title="文件预览"
      style="width: 80%; max-width: 900px;"
      :bordered="false"
    >
      <n-tabs type="segment">
        <n-tab-pane name="model" tab="模型">
          <n-code
            :code="previewData.modelCode"
            language="go"
            show-line-numbers
          />
        </n-tab-pane>
        <n-tab-pane name="dto" tab="DTO">
          <n-code
            :code="previewData.dtoCode"
            language="go"
            show-line-numbers
          />
        </n-tab-pane>
        <n-tab-pane name="controller" tab="控制器">
          <n-code
            :code="previewData.controllerCode"
            language="go"
            show-line-numbers
          />
        </n-tab-pane>
        <n-tab-pane name="router" tab="路由">
          <n-code
            :code="previewData.routeCode"
            language="go"
            show-line-numbers
          />
        </n-tab-pane>
        <n-tab-pane name="frontend" tab="前端">
          <n-code
            :code="previewData.frontendCode"
            language="html"
            show-line-numbers
          />
        </n-tab-pane>
      </n-tabs>
    </n-modal>
  </div>
</template>

<script lang="ts" setup>
import { h, reactive, ref } from 'vue';
import { useMessage } from 'naive-ui';
import type { DataTableColumns } from 'naive-ui';
import { deleteHistory, getHistoryList, rollbackCode } from '@/service/api/codegen';
import type { HistoryRecord } from '@/service/api/codegen';

// Props 定义
const props = defineProps({
  onRefresh: {
    type: Function,
    default: () => {}
  }
});

// 消息组件
const message = useMessage();

// 历史记录数据
const historyLoading = ref(false);
const historyData = ref<HistoryRecord[]>([]);
const historyPagination = reactive({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 30, 50],
});

// 定义历史记录表格列
const historyColumns: DataTableColumns = [
  { 
    title: 'ID',
    key: 'id',
    width: 80
  },
  { 
    title: '结构体名称',
    key: 'structName',
    width: 120
  },
  { 
    title: '表名',
    key: 'table',
    width: 120
  },
  {
    title: '包名',
    key: 'packageName',
    width: 120
  },
  {
    title: '描述',
    key: 'description',
    width: 150
  },
  {
    title: '创建时间',
    key: 'createdAt',
    width: 160
  },
  {
    title: '状态',
    key: 'flag',
    width: 100,
    render(row: HistoryRecord) {
      return row.flag === 0 ? '正常' : '已回滚';
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 200,
    render(row: HistoryRecord) {
      if (row.flag === 0) {
        return h('div', { class: 'flex gap-2' }, [
          h(
            'button',
            {
              class: 'text-blue-500 hover:text-blue-700',
              onClick: () => handlePreview(row)
            },
            '预览'
          ),
          h(
            'button',
            {
              class: 'text-blue-500 hover:text-blue-700',
              onClick: () => handleRollback(row)
            },
            '回滚'
          ),
          h(
            'button',
            {
              class: 'text-red-500 hover:text-red-700',
              onClick: () => handleDelete(row.id)
            },
            '删除'
          )
        ]);
      } else {
        return h('div', { class: 'flex gap-2' }, [
          h(
            'button',
            {
              class: 'text-blue-500 hover:text-blue-700',
              onClick: () => handlePreview(row)
            },
            '预览'
          ),
          h(
            'button',
            {
              class: 'text-red-500 hover:text-red-700',
              onClick: () => handleDelete(row.id)
            },
            '删除'
          )
        ]);
      }
    }
  }
];

// 加载历史记录数据
async function loadHistoryData() {
  try {
    historyLoading.value = true;
    const { page, pageSize } = historyPagination;
    const res = await getHistoryList(page, pageSize);
    
    if (res.data) {
      historyData.value = res.data.list;
      historyPagination.itemCount = res.data.total;
    }
  } catch (error) {
    message.error('获取历史记录失败');
  } finally {
    historyLoading.value = false;
  }
}

// 处理页码变化
function handleHistoryPageChange(page: number) {
  historyPagination.page = page;
  loadHistoryData();
}

// 处理每页条数变化
function handleHistoryPageSizeChange(pageSize: number) {
  historyPagination.pageSize = pageSize;
  historyPagination.page = 1;
  loadHistoryData();
}

// 回滚相关
const showRollbackModal = ref(false);
const rollbackLoading = ref(false);
const currentRollbackId = ref<number>(0);
const rollbackOptions = reactive({
  deleteFiles: true,
  deleteApi: false,
  deleteMenu: false,
  deleteTable: false
});

// 处理回滚
function handleRollback(row: HistoryRecord) {
  currentRollbackId.value = row.id;
  showRollbackModal.value = true;
}

// 确认回滚
async function confirmRollback() {
  if (currentRollbackId.value === 0) {
    showRollbackModal.value = false;
    return;
  }

  try {
    rollbackLoading.value = true;
    await rollbackCode({
      id: currentRollbackId.value,
      deleteFiles: rollbackOptions.deleteFiles,
      deleteApi: rollbackOptions.deleteApi,
      deleteMenu: rollbackOptions.deleteMenu,
      deleteTable: rollbackOptions.deleteTable
    });

    message.success('回滚成功');
    showRollbackModal.value = false;
    loadHistoryData();
    props.onRefresh();
  } catch (error) {
    message.error('回滚失败');
  } finally {
    rollbackLoading.value = false;
  }
}

// 取消回滚
function cancelRollback() {
  showRollbackModal.value = false;
}

// 处理删除
async function handleDelete(id: number) {
  try {
    await deleteHistory(id);
    message.success('删除成功');
    loadHistoryData();
  } catch (error) {
    message.error('删除失败');
  }
}

// 文件预览功能
const showPreviewModal = ref(false);
const previewData = reactive({
  modelCode: '',
  dtoCode: '',
  controllerCode: '',
  routeCode: '',
  frontendCode: '',
});

// 处理预览
function handlePreview(row: HistoryRecord) {
  // 这里应该调用后端API获取生成的文件内容
  // 暂时使用模拟数据
  previewData.modelCode = `package models

import (
  "time"
  "gorm.io/gorm"
)

// ${row.structName} ${row.description}
type ${row.structName} struct {
  ID        uint           \`gorm:"primarykey" json:"id"\`
  CreatedAt time.Time      \`json:"createdAt"\`
  UpdatedAt time.Time      \`json:"updatedAt"\`
  DeletedAt gorm.DeletedAt \`gorm:"index" json:"-"\`
  // ...其他字段
}

// TableName 设置表名
func (${row.structName}) TableName() string {
  return "${row.table}"
}`;

  previewData.dtoCode = `package dto

// ${row.structName}CreateRequest 创建${row.description}请求
type ${row.structName}CreateRequest struct {
  // ...字段
}

// ${row.structName}Response ${row.description}响应
type ${row.structName}Response struct {
  ID        uint   \`json:"id"\`
  CreatedAt string \`json:"createdAt"\`
  UpdatedAt string \`json:"updatedAt"\`
  // ...其他字段
}`;

  previewData.controllerCode = `package controllers

import (
  "github.com/gin-gonic/gin"
)

// ${row.structName}Controller ${row.description}控制器
type ${row.structName}Controller struct{}

// New${row.structName}Controller 创建${row.description}控制器
func New${row.structName}Controller() *${row.structName}Controller {
  return &${row.structName}Controller{}
}`;

  previewData.routeCode = `package routes

import (
  "github.com/gin-gonic/gin"
)

// 路由注册
${row.structName}Controller := controllers.New${row.structName}Controller()
privateRoutes.GET("/${row.packageName}/${row.table}", ${row.structName}Controller.Get${row.structName}s)
privateRoutes.GET("/${row.packageName}/${row.table}/:id", ${row.structName}Controller.Get${row.structName})
privateRoutes.POST("/${row.packageName}/${row.table}", ${row.structName}Controller.Create${row.structName})
privateRoutes.PUT("/${row.packageName}/${row.table}", ${row.structName}Controller.Update${row.structName})
privateRoutes.DELETE("/${row.packageName}/${row.table}/:id", ${row.structName}Controller.Delete${row.structName})`;

  previewData.frontendCode = `<template>
  <CommonWrapper>
    <template #title>${row.description}管理</template>
    
    <template #buttons>
      <n-button type="primary" @click="openModal()">
        <template #icon><icon-mdi:plus-circle /></template>
        添加${row.description}
      </n-button>
    </template>
    
    <template #content>
      <!-- 表格和表单 -->
    </template>
  </CommonWrapper>
</template>`;

  showPreviewModal.value = true;
}

// 初始加载历史记录
loadHistoryData();

defineExpose({
  loadHistoryData
});
</script> 