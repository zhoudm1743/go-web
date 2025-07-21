<template>
  <CommonWrapper>
    <template #title>
      文章管理
    </template>

    <template #buttons>
      <n-button type="primary" @click="handleAdd">
        <template #icon>
          <icon-mdi:plus-circle />
        </template>
        新增文章
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
import { createArticle, deleteArticle, getArticles, updateArticle } from '@/service/api/article';

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
const title = '文章管理';
const description = '文章';

// 表格列
const columns = [
  {
    title: '标题',
    key: 'title',
    width: 120,
  },
  {
    title: '分类',
    key: 'categoryID',
    width: 120,
  },
  {
    title: '作者',
    key: 'author',
    width: 120,
  },

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
    const res = await getArticles({
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
  modalTitle.value = '新增文章';
  modalMode.value = 'add';
  modalRef.value.openModal();
}

// 处理编辑
function handleEdit(row) {
  modalTitle.value = '编辑文章';
  modalMode.value = 'edit';
  modalRef.value.openModal(row);
}

// 处理删除
async function handleDelete(row) {
  try {
    await deleteArticle(row.id);
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
      await createArticle(formData);
      message.success('创建成功');
    } else {
      await updateArticle(formData);
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
