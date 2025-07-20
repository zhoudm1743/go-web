<script setup lang="tsx">
import type { DataTableColumns } from 'naive-ui'
import CopyText from '@/components/custom/CopyText.vue'
import { useBoolean } from '@/hooks'
import { fetchDictList, deleteDict } from '@/service'
import { useDictStore } from '@/store'
import { NButton, NFlex, NPopconfirm } from 'naive-ui'
import DictModal from './components/DictModal.vue'

const { bool: dictLoading, setTrue: startDictLoading, setFalse: endDictLoading } = useBoolean(false)
const { bool: contentLoading, setTrue: startContentLoading, setFalse: endContentLoading } = useBoolean(false)

const { getDictByNet } = useDictStore()

const dictRef = ref<InstanceType<typeof DictModal>>()
const dictContentRef = ref<InstanceType<typeof DictModal>>()

onMounted(() => {
  getDictList()
})

const dictData = ref<Entity.Dict[]>([])
const dictContentData = ref<Entity.Dict[]>([])

async function getDictList() {
  startDictLoading()
  const { data, isSuccess } = await fetchDictList()
  if (isSuccess) {
    dictData.value = data
  }
  endDictLoading()
}

const lastDictCode = ref('')
async function getDictContent(code: string) {
  startContentLoading()
  dictContentData.value = await getDictByNet(code)
  lastDictCode.value = code
  endContentLoading()
}

// 处理字典项成功事件
function handleDictSuccess() {
  getDictList()
}

// 处理字典值成功事件
function handleDictContentSuccess() {
  if (lastDictCode.value) {
    getDictContent(lastDictCode.value)
  }
}

async function deleteDictItem(id: number) {
  try {
    const { isSuccess } = await deleteDict(id)
    if (isSuccess) {
      window.$message.success('字典删除成功')
      getDictList()
      
      // 如果当前正在查看的字典内容被删除，清空字典内容
      if (lastDictCode.value) {
        getDictContent(lastDictCode.value)
      }
    }
  } catch (error) {
    console.error('删除字典失败:', error)
    window.$message.error('删除字典失败')
  }
}

const dictColumns: DataTableColumns<Entity.Dict> = [
  {
    title: '字典项',
    key: 'label',
  },
  {
    title: '字典码',
    key: 'code',
    render: (row) => {
      return (
        <CopyText value={row.code} />
      )
    },
  },
  {
    title: '操作',
    key: 'actions',
    align: 'center',
    render: (row) => {
      return (
        <NFlex justify="center">
          <NButton
            size="small"
            onClick={() => getDictContent(row.code)}
          >
            查看字典
          </NButton>
          <NButton
            size="small"
            onClick={() => dictRef.value!.openModal('edit', row)}
          >
            编辑
          </NButton>
          <NPopconfirm onPositiveClick={() => deleteDictItem(row.id!)}>
            {{
              default: () => (
                <span>
                  确认删除字典
                  <b>{row.label}</b>
                  {' '}
                  ？
                </span>
              ),
              trigger: () => <NButton size="small" type="error">删除</NButton>,
            }}
          </NPopconfirm>
        </NFlex>
      )
    },
  },
]

const contentColumns: DataTableColumns<Entity.Dict> = [
  {
    title: '字典名称',
    key: 'label',
  },
  {
    title: '字典码',
    key: 'code',
  },
  {
    title: '字典值',
    key: 'value',
  },
  {
    title: '操作',
    key: 'actions',
    align: 'center',
    width: '15em',
    render: (row) => {
      return (
        <NFlex justify="center">
          <NButton
            size="small"
            onClick={() => dictContentRef.value!.openModal('edit', row)}
          >
            编辑
          </NButton>
          <NPopconfirm onPositiveClick={() => deleteDictItem(row.id!)}>
            {{
              default: () => (
                <span>
                  确认删除字典值
                  <b>{row.label}</b>
                  {' '}
                  ？
                </span>
              ),
              trigger: () => <NButton size="small" type="error">删除</NButton>,
            }}
          </NPopconfirm>
        </NFlex>
      )
    },
  },
]
</script>

<template>
  <NFlex>
    <div class="basis-2/5">
      <n-card>
        <template #header>
          <NButton type="primary" @click="dictRef!.openModal('add')">
            <template #icon>
              <icon-park-outline-add-one />
            </template>
            新建
          </NButton>
        </template>
        <template #header-extra>
          <NFlex>
            <NButton type="primary" secondary @click="getDictList">
              <template #icon>
                <icon-park-outline-refresh />
              </template>
              刷新
            </NButton>
          </NFlex>
        </template>
        <n-data-table
          :columns="dictColumns" :data="dictData" :loading="dictLoading" :pagination="false"
          :bordered="false"
        />
      </n-card>
    </div>
    <div class="flex-1">
      <n-card>
        <template #header>
          <NButton type="primary" :disabled="!lastDictCode" @click="dictContentRef!.openModal('add')">
            <template #icon>
              <icon-park-outline-add-one />
            </template>
            新建
          </NButton>
        </template>
        <template #header-extra>
          <NFlex>
            <NButton type="primary" :disabled="!lastDictCode" secondary @click="getDictContent(lastDictCode)">
              <template #icon>
                <icon-park-outline-refresh />
              </template>
              刷新
            </NButton>
          </NFlex>
        </template>
        <n-data-table
          :columns="contentColumns" :data="dictContentData" :loading="contentLoading" :pagination="false"
          :bordered="false"
        />
      </n-card>
    </div>

    <DictModal ref="dictRef" modal-name="字典项" is-root @success="handleDictSuccess" />
    <DictModal ref="dictContentRef" modal-name="字典值" :dict-code="lastDictCode" @success="handleDictContentSuccess" />
  </NFlex>
</template>

<style scoped></style>
