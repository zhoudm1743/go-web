<script setup lang="ts">
import { useBoolean } from '@/hooks'
import { ref } from 'vue'
import { fetchAllRoutes, getRoleMenus, updateRoleMenus } from '@/service'

const { bool: modalVisible, setTrue: showModal, setFalse: hiddenModal } = useBoolean(false)
const { bool: submitLoading, setTrue: startLoading, setFalse: endLoading } = useBoolean(false)
const { bool: treeLoading, setTrue: startTreeLoading, setFalse: endTreeLoading } = useBoolean(false)

// 当前编辑的角色
const currentRole = ref<Entity.Role | null>(null)
// 选中的菜单ID
const checkedKeys = ref<number[]>([])
// 菜单树数据
const treeData = ref<any[]>([])

// 打开弹窗
async function openModal(role: Entity.Role) {
  currentRole.value = role
  showModal()
  await loadMenuTree()
  await loadRoleMenus(role.id)
}

// 加载菜单树
async function loadMenuTree() {
  startTreeLoading()
  try {
    const { data, isSuccess } = await fetchAllRoutes()
    if (isSuccess) {
      treeData.value = data
    }
  } 
  finally {
    endTreeLoading()
  }
}

// 加载角色已有菜单
async function loadRoleMenus(roleId: number) {
  try {
    const { data, isSuccess } = await getRoleMenus(roleId)
    if (isSuccess) {
      checkedKeys.value = data
    }
  } catch (error) {
    window.$message.error('获取角色菜单失败')
    console.error('获取角色菜单失败', error)
  }
}

// 关闭弹窗
function closeModal() {
  hiddenModal()
  endLoading()
  currentRole.value = null
  checkedKeys.value = []
}

// 提交分配
async function submitAssign() {
  if (!currentRole.value) return
  
  startLoading()
  try {
    // 调用API保存角色菜单关联
    const { isSuccess } = await updateRoleMenus({
      roleId: currentRole.value.id,
      menuIds: checkedKeys.value
    })
    
    if (isSuccess) {
      window.$message.success('菜单分配成功')
      closeModal()
    }
  }
  catch (error) {
    console.error('菜单分配失败', error)
    window.$message.error('操作失败，请重试')
  }
  finally {
    endLoading()
  }
}

// 暴露方法给父组件调用
defineExpose({
  openModal,
})
</script>

<template>
  <n-modal
    v-model:show="modalVisible" 
    :mask-closable="false" 
    preset="card" 
    :title="`角色菜单分配 - ${currentRole?.name || ''}`"
    class="w-600px"
    style="max-height: 80vh"
  >
    <n-spin :show="treeLoading">
      <n-tree
        checkable
        cascade
        v-model:checked-keys="checkedKeys"
        :data="treeData"
        key-field="id"
        label-field="title"
        children-field="children"
        :render-label="({ option }) => {
          return option.title + ' [' + option.path + ']'
        }"
      />
    </n-spin>
    
    <template #footer>
      <n-space justify="end">
        <n-button @click="closeModal">取消</n-button>
        <n-button type="primary" :loading="submitLoading" @click="submitAssign">
          确定
        </n-button>
      </n-space>
    </template>
  </n-modal>
</template> 