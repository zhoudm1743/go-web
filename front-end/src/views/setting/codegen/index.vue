<template>
  <n-card>
    <template #header>代码生成器</template>
    <template #header-extra>
      <n-button type="primary" @click="openGeneratorForm">
        <template #icon>
          <icon-mdi:plus-circle />
        </template>
        创建新代码
      </n-button>
    </template>
      
      <!-- 选项卡 -->
      <n-tabs type="line" animated v-model:value="activeTab">
        <n-tab-pane name="generator" tab="代码生成器">
          <GeneratorForm
            v-if="showGeneratorForm"
            ref="generatorFormRef"
            :is-visible="showGeneratorForm"
            @cancel="closeGeneratorForm"
            @success="onGenerateSuccess"
          />

          <n-card v-else :bordered="false">
            <n-result
              status="info"
              title="代码生成器"
              description="使用代码生成器可以快速创建模型、控制器和视图，支持关系字段和预加载，提高开发效率。"
            >
              <template #footer>
                <n-button type="primary" @click="openGeneratorForm">
                  创建新代码
                </n-button>
              </template>
            </n-result>
          </n-card>
        </n-tab-pane>

        <n-tab-pane name="history" tab="历史记录">
          <HistoryList ref="historyListRef" :on-refresh="refreshForm" />
        </n-tab-pane>
      </n-tabs>
  </n-card>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import { useMessage } from 'naive-ui';
import GeneratorForm from './components/GeneratorForm.vue';
import HistoryList from './components/HistoryList.vue';

// 消息组件
const message = useMessage();


// 当前选中的标签
const activeTab = ref('generator');

// 代码生成器表单
const showGeneratorForm = ref(false);
const generatorFormRef = ref();

// 历史记录组件引用
const historyListRef = ref();

// 打开生成器表单
function openGeneratorForm() {
  showGeneratorForm.value = true;
}

// 关闭生成器表单
function closeGeneratorForm() {
  showGeneratorForm.value = false;
}

// 生成成功回调
function onGenerateSuccess() {
  message.success('代码生成成功');
  closeGeneratorForm();
  
  // 切换到历史记录标签
  activeTab.value = 'history';
  
  // 刷新历史记录
  if (historyListRef.value) {
    historyListRef.value.loadHistoryData();
  }
}

// 刷新表单
function refreshForm() {
  // 用于历史记录组件调用，刷新表单或其他操作
  showGeneratorForm.value = false;
}
</script> 