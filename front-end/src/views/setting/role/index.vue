<script setup lang="tsx">
import type { DataTableColumns } from 'naive-ui'
import { useBoolean } from '@/hooks'
import { NButton, NPopconfirm, NSpace, NTag } from 'naive-ui'
import { fetchRoles, deleteRole } from '@/service'
import RoleModal from './components/RoleModal.vue'
import RoleMenuModal from './components/RoleMenuModal.vue'

const { bool: loading, setTrue: startLoading, setFalse: endLoading } = useBoolean(false)

const roleModalRef = ref()
const roleMenuModalRef = ref()

onMounted(() => {
  getRoleList()
})

const roles = ref<Entity.Role[]>([])

async function getRoleList() {
  startLoading()
  try {
    const { data, isSuccess } = await fetchRoles()
    if (isSuccess)
      roles.value = data
  }
  finally {
    endLoading()
  }
}

async function handleDeleteRole(id: number) {
  window.$dialog.warning({
    title: '确认删除',
    content: '确定要删除这个角色吗？删除后不可恢复，且会影响已分配该角色的用户。',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        const { isSuccess } = await deleteRole(id)
        if (isSuccess) {
          window.$message.success('删除角色成功')
          getRoleList()
        }
      } catch (error) {
        console.error('删除角色失败', error)
        window.$message.error('删除失败，请重试')
      }
    },
  })
}

// 分配菜单
function assignMenus(role: Entity.Role) {
  roleMenuModalRef.value.openModal(role)
}

const columns: DataTableColumns<Entity.Role> = [
  {
    title: '角色名称',
    key: 'name',
  },
  {
    title: '角色编码',
    key: 'code',
  },
  {
    title: '状态',
    key: 'status',
    align: 'center',
    width: 100,
    render: (row) => {
      // 处理字符串或数字类型的状态值
      const statusValue = Number(row.status)
      const status = statusValue === 1
      return (
        <NTag type={status ? 'success' : 'error'}>
          { status ? '启用' : '禁用' }
        </NTag>
      )
    },
  },
  {
    title: '排序',
    key: 'sort',
    align: 'center',
    width: 100,
  },
  {
    title: '备注',
    key: 'remark',
  },
  {
    title: '操作',
    key: 'actions',
    align: 'center',
    width: 280,
    render: (row) => {
      return (
        <NSpace justify="center">
          <NButton
            size="small"
            type="primary"
            onClick={() => assignMenus(row)}
          >
            分配菜单
          </NButton>
          <NButton
            size="small"
            onClick={() => roleModalRef.value.openModal('edit', row)}
          >
            编辑
          </NButton>
          <NPopconfirm onPositiveClick={() => handleDeleteRole(row.id)}>
            {{
              default: () => '确认删除',
              trigger: () => <NButton size="small" type="error">删除</NButton>,
            }}
          </NPopconfirm>
        </NSpace>
      )
    },
  },
]

const checkedRowKeys = ref<number[]>([])
</script>

<template>
  <n-card>
    <template #header>
      <n-button type="primary" @click="roleModalRef.openModal('add')">
        <template #icon>
          <icon-park-outline-add-one />
        </template>
        新建角色
      </n-button>
    </template>

    <template #header-extra>
      <n-flex>
        <n-button type="primary" secondary @click="getRoleList">
          <template #icon>
            <icon-park-outline-refresh />
          </template>
          刷新
        </n-button>
      </n-flex>
    </template>

    <n-data-table
      v-model:checked-row-keys="checkedRowKeys"
      :row-key="(row:Entity.Role) => row.id"
      :columns="columns"
      :data="roles"
      :loading="loading"
      size="small"
      :scroll-x="1200"
    />
    
    <role-modal ref="roleModalRef" @success="getRoleList" />
    <role-menu-modal ref="roleMenuModalRef" />
  </n-card>
</template> 