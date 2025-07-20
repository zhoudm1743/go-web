<script setup lang="ts">
import type { FormItemRule } from 'naive-ui'
import { useBoolean } from '@/hooks'
import { computed, ref } from 'vue'
import { createRole, updateRole } from '@/service'

interface Props {
  modalName?: string
}

const {
  modalName = '角色',
} = defineProps<Props>()

const emit = defineEmits<{
  success: []
}>()

const { bool: modalVisible, setTrue: showModal, setFalse: hiddenModal } = useBoolean(false)
const { bool: submitLoading, setTrue: startLoading, setFalse: endLoading } = useBoolean(false)

const formDefault = {
  id: 0,
  name: '',
  code: '',
  sort: 0,
  status: 1,
  remark: ''
}

const formModel = ref({ ...formDefault })

type ModalType = 'add' | 'view' | 'edit'
const modalType = ref<ModalType>('add')
const modalTitle = computed(() => {
  const titleMap: Record<ModalType, string> = {
    add: '添加',
    view: '查看',
    edit: '编辑',
  }
  return `${titleMap[modalType.value]}${modalName}`
})

function openModal(type: ModalType = 'add', data?: Entity.Role) {
  modalType.value = type
  showModal()
  
  if (type === 'add') {
    formModel.value = { ...formDefault }
  } 
  else if (data) {
    formModel.value = { ...data }
  }
}

function closeModal() {
  hiddenModal()
  endLoading()
}

defineExpose({
  openModal,
})

const formRef = ref()

async function submitModal() {
  await formRef.value?.validate()
  startLoading()
  
  try {
    // 调用API保存角色
    const api = modalType.value === 'add' ? createRole : updateRole
    const { isSuccess } = await api(formModel.value)
    
    if (isSuccess) {
      window.$message.success(`${modalType.value === 'add' ? '新增' : '更新'}角色成功`)
      closeModal()
      emit('success')
    }
  }
  catch (error) {
    console.error('角色保存失败', error)
    window.$message.error('操作失败，请重试')
  }
  finally {
    endLoading()
  }
}

// 表单验证规则
const rules = {
  name: {
    required: true,
    message: '请输入角色名称',
    trigger: 'blur',
  },
  code: {
    required: true,
    message: '请输入角色编码',
    trigger: 'blur',
  },
}
</script>

<template>
  <n-modal
    v-model:show="modalVisible" 
    :mask-closable="false" 
    preset="card" 
    :title="modalTitle" 
    class="w-500px"
  >
    <n-form
      ref="formRef"
      :rules="rules"
      label-placement="left" 
      :label-width="80"
      :model="formModel"
      :disabled="modalType === 'view'"
    >
      <n-form-item label="角色名称" path="name">
        <n-input v-model:value="formModel.name" placeholder="请输入角色名称" />
      </n-form-item>
      
      <n-form-item label="角色编码" path="code">
        <n-input v-model:value="formModel.code" placeholder="请输入角色编码" />
      </n-form-item>
      
      <n-form-item label="排序" path="sort">
        <n-input-number v-model:value="formModel.sort" />
      </n-form-item>
      
      <n-form-item label="状态" path="status">
        <n-radio-group v-model:value="formModel.status" name="status">
          <n-space>
            <n-radio :value="1">启用</n-radio>
            <n-radio :value="2">禁用</n-radio>
          </n-space>
        </n-radio-group>
      </n-form-item>
      
      <n-form-item label="备注" path="remark">
        <n-input v-model:value="formModel.remark" type="textarea" placeholder="请输入备注" />
      </n-form-item>
    </n-form>
    
    <template #footer>
      <n-space justify="end">
        <n-button @click="closeModal">取消</n-button>
        <n-button v-if="modalType !== 'view'" type="primary" :loading="submitLoading" @click="submitModal">
          确定
        </n-button>
      </n-space>
    </template>
  </n-modal>
</template> 