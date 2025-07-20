<script setup lang="ts">
import { useBoolean } from '@/hooks'
import { fetchRoleList, createAdmin, updateAdmin } from '@/service'

interface Props {
  modalName?: string
}

const {
  modalName = '',
} = defineProps<Props>()

const emit = defineEmits<{
  open: []
  close: []
  success: []
}>()

const { bool: modalVisible, setTrue: showModal, setFalse: hiddenModal } = useBoolean(false)

const { bool: submitLoading, setTrue: startLoading, setFalse: endLoading } = useBoolean(false)

const formDefault: Entity.Admin = {
  username: '',
  nickname: '',
  email: '',
  mobile: '',
  roleId: undefined,
  status: 1,
}
const formModel = ref<Entity.Admin>({ ...formDefault })

type ModalType = 'add' | 'view' | 'edit'
const modalType = shallowRef<ModalType>('add')
const modalTitle = computed(() => {
  const titleMap: Record<ModalType, string> = {
    add: '添加',
    view: '查看',
    edit: '编辑',
  }
  return `${titleMap[modalType.value]}${modalName}`
})

async function openModal(type: ModalType = 'add', data: any) {
  emit('open')
  modalType.value = type
  showModal()
  getRoleList()
  const handlers = {
    async add() {
      formModel.value = { ...formDefault }
    },
    async view() {
      if (!data)
        return
      formModel.value = { ...data }
    },
    async edit() {
      if (!data)
        return
      formModel.value = { ...data }
    },
  }
  await handlers[type]()
}

function closeModal() {
  hiddenModal()
  endLoading()
  emit('close')
}

defineExpose({
  openModal,
})

const formRef = ref()
async function submitModal() {
  try {
    await formRef.value?.validate()
    startLoading()
    
    const handlers = {
      async add() {
        const { isSuccess } = await createAdmin(formModel.value)
        if (isSuccess) {
          window.$message.success('管理员创建成功')
          emit('success')
          return true
        }
        return false
      },
      async edit() {
        const { isSuccess } = await updateAdmin(formModel.value)
        if (isSuccess) {
          window.$message.success('管理员更新成功')
          emit('success')
          return true
        }
        return false
      },
      async view() {
        return true
      },
    }
    
    const result = await handlers[modalType.value]()
    if (result) {
      closeModal()
    }
  } catch (error) {
    console.error('提交表单出错:', error)
  } finally {
    endLoading()
  }
}

const rules = {
  username: {
    required: true,
    message: '请输入用户名',
    trigger: 'blur',
  },
}

const options = ref()
async function getRoleList() {
  const { data } = await fetchRoleList()
  options.value = data
}
</script>

<template>
  <n-modal
    v-model:show="modalVisible"
    :mask-closable="false"
    preset="card"
    :title="modalTitle"
    class="w-700px"
    :segmented="{
      content: true,
      action: true,
    }"
  >
    <n-form ref="formRef" :rules="rules" label-placement="left" :model="formModel" :label-width="100" :disabled="modalType === 'view'">
      <n-grid :cols="2" :x-gap="18">
        <n-form-item-grid-item :span="1" label="用户名" path="username">
          <n-input v-model:value="formModel.username" />
        </n-form-item-grid-item>
        <n-form-item-grid-item :span="1" label="昵称" path="nickname">
          <n-input v-model:value="formModel.nickname" />
        </n-form-item-grid-item>
        <n-form-item-grid-item :span="1" label="真实姓名" path="realName">
          <n-input v-model:value="formModel.realName" />
        </n-form-item-grid-item>
        <n-form-item-grid-item :span="1" label="邮箱" path="email">
          <n-input v-model:value="formModel.email" />
        </n-form-item-grid-item>
        <n-form-item-grid-item :span="1" label="联系方式" path="mobile">
          <n-input v-model:value="formModel.mobile" />
        </n-form-item-grid-item>
        <n-form-item-grid-item :span="1" label="角色" path="roleId">
          <n-select
            v-model:value="formModel.roleId" filterable
            label-field="name"
            value-field="id"
            :options="options"
          />
        </n-form-item-grid-item>
        <n-form-item-grid-item :span="1" label="用户状态" path="status">
          <n-switch
            v-model:value="formModel.status"
            :checked-value="1" :unchecked-value="2"
          >
            <template #checked>
              启用
            </template>
            <template #unchecked>
              禁用
            </template>
          </n-switch>
        </n-form-item-grid-item>
      </n-grid>
    </n-form>
    <template #action>
      <n-space justify="center">
        <n-button @click="closeModal">
          取消
        </n-button>
        <n-button type="primary" :loading="submitLoading" @click="submitModal">
          提交
        </n-button>
      </n-space>
    </template>
  </n-modal>
</template>
