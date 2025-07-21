<template>
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
      <n-form-item label="产品分类" path="categoryID">
        <n-input v-model:value="formData.categoryID" placeholder="请输入产品分类" />
      </n-form-item>
      <n-form-item label="名称" path="name">
        <n-input v-model:value="formData.name" placeholder="请输入名称" />
      </n-form-item>
      <n-form-item label="排序" path="sort">
        <n-input v-model:value="formData.sort" placeholder="请输入排序" />
      </n-form-item>

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
  iD: 0,
  categoryID: '',
  name: '',
  sort: null,

});

// 初始表单数据
const initialFormData = {
  iD: 0,
  categoryID: '',
  name: '',
  sort: null,

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
  name: {
    required: true,
    message: '请输入名称',
    trigger: ['blur', 'input'],
  },

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
