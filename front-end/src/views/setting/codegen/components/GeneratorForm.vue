<template>
  <n-drawer
    v-model:show="localVisible"
    :width="800"
    :height="'90%'"
    placement="right"
    :mask-closable="false"
    @close="handleClose"
  >
    <n-drawer-content :title="title" :native-scrollbar="false">
      <n-spin :show="loading">
        <n-steps :current="currentStep" :status="currentStatus">
          <n-step title="基本设置" description="设置代码生成的基本信息" />
          <n-step title="字段设置" description="设置字段信息" />
          <n-step title="生成选项" description="设置生成选项" />
        </n-steps>

        <div class="mt-6">
          <!-- 步骤1: 基本设置 -->
          <div v-show="currentStep === 0">
            <n-form
              ref="baseFormRef"
              :model="formData"
              :rules="baseRules"
              label-placement="left"
              label-width="100"
              require-mark-placement="right-hanging"
            >
              <!-- 选择应用 -->
              <n-form-item label="所属应用" path="appName">
                <n-input-group>
                  <n-select
                    v-model:value="formData.appName"
                    placeholder="选择应用"
                    :options="appOptions"
                    style="width: 80%"
                    @update:value="onAppChange"
                  />
                  <n-button type="primary" style="width: 20%" @click="showNewAppField = true">
                    新建应用
                  </n-button>
                </n-input-group>
              </n-form-item>

              <!-- 新建应用字段 -->
              <n-form-item v-if="showNewAppField" label="新应用名称" path="newAppName">
                <n-input
                  v-model:value="formData.newAppName"
                  placeholder="请输入新应用名称，英文字母，如：order"
                />
              </n-form-item>

              <n-form-item label="包名" path="packageName">
                <n-input
                  v-model:value="formData.packageName"
                  placeholder="请输入包名，如：admin"
                />
              </n-form-item>

              <n-form-item label="结构体名称" path="structName">
                <n-input
                  v-model:value="formData.structName"
                  placeholder="请输入结构体名称，如：Product"
                />
              </n-form-item>

              <n-form-item label="表名" path="tableName">
                <n-input-group>
                  <n-input
                    v-model:value="formData.tableName"
                    placeholder="请输入表名，如：products"
                    style="width: calc(100% - 80px)"
                  />
                  <n-button type="primary" style="width: 80px" @click="generateTableName">自动生成</n-button>
                </n-input-group>
              </n-form-item>

              <n-form-item label="描述" path="description">
                <n-input-group>
                  <n-input
                    v-model:value="formData.description"
                    placeholder="请输入描述，如：产品"
                    style="width: calc(100% - 80px)"
                  />
                  <n-button type="primary" style="width: 80px" @click="generateDescription">自动生成</n-button>
                </n-input-group>
              </n-form-item>

              <n-form-item label="API前缀" path="apiPrefix">
                <n-input-group>
                  <n-input
                    v-model:value="formData.apiPrefix"
                    placeholder="请输入API前缀，如：product"
                    style="width: calc(100% - 80px)"
                  />
                  <n-button type="primary" style="width: 80px" @click="generateApiPrefix">自动生成</n-button>
                </n-input-group>
              </n-form-item>

              <n-form-item label="从数据表导入">
                <n-button @click="importFromTable">从数据表导入字段</n-button>
              </n-form-item>
            </n-form>
          </div>

          <!-- 步骤2: 字段设置 -->
          <div v-show="currentStep === 1">
            <div class="flex justify-end mb-2">
              <n-button @click="addField">添加字段</n-button>
            </div>
            <n-data-table
              :columns="fieldColumns"
              :data="formData.fields"
              :pagination="false"
              :bordered="true"
              :row-key="(row) => row._id"
            />
          </div>

          <!-- 步骤3: 生成选项 -->
          <div v-show="currentStep === 2">
            <n-form
              ref="optionsFormRef"
              :model="formData"
              label-placement="left"
              label-width="160"
            >
              <n-grid :cols="3" :x-gap="12">
                <n-grid-item>
                  <n-form-item label="生成列表功能">
                    <n-switch v-model:value="formData.hasList" />
                  </n-form-item>
                </n-grid-item>

                <n-grid-item>
                  <n-form-item label="生成创建功能">
                    <n-switch v-model:value="formData.hasCreate" />
                  </n-form-item>
                </n-grid-item>

                <n-grid-item>
                  <n-form-item label="生成更新功能">
                    <n-switch v-model:value="formData.hasUpdate" />
                  </n-form-item>
                </n-grid-item>

                <n-grid-item>
                  <n-form-item label="生成删除功能">
                    <n-switch v-model:value="formData.hasDelete" />
                  </n-form-item>
                </n-grid-item>

                <n-grid-item>
                  <n-form-item label="生成详情功能">
                    <n-switch v-model:value="formData.hasDetail" />
                  </n-form-item>
                </n-grid-item>

                <n-grid-item>
                  <n-form-item label="启用分页">
                    <n-switch v-model:value="formData.hasPagination" />
                  </n-form-item>
                </n-grid-item>
              </n-grid>
            </n-form>

            <div class="mt-4">
              <n-alert title="即将生成的文件" type="info">
                <template #icon>
                  <n-icon><icon-ant-design:info-circle-outlined /></n-icon>
                </template>
                <div class="text-sm">
                  <p>模型文件: server/apps/{{ formData.appName || formData.newAppName }}/models/{{ formData.tableName }}.go</p>
                  <p>DTO文件: server/apps/{{ formData.appName || formData.newAppName }}/dto/{{ formData.tableName }}.go</p>
                  <p>控制器文件: server/apps/{{ formData.appName || formData.newAppName }}/controllers/{{ formData.tableName }}_controller.go</p>
                  <p>前端页面: front-end/src/views/{{ formData.appName || formData.newAppName }}/{{ formData.tableName }}/index.vue</p>
                  <p>前端API: front-end/src/service/api/{{ formData.tableName }}.ts</p>
                </div>
              </n-alert>
            </div>
          </div>

          <!-- 步骤操作按钮 -->
          <div class="flex justify-between mt-4">
            <n-button 
              v-if="currentStep > 0"
              @click="prevStep"
            >
              上一步
            </n-button>
            <div></div>
            <n-button 
              type="primary" 
              @click="currentStep < 2 ? nextStep() : generateCode()"
              :loading="loading"
            >
              {{ currentStep < 2 ? '下一步' : '生成代码' }}
            </n-button>
          </div>
        </div>

        <!-- 导入表字段对话框 -->
        <n-modal
          v-model:show="showImportTableModal"
          title="从数据表导入"
          preset="dialog"
          positive-text="导入"
          negative-text="取消"
          @positive-click="confirmImportTable"
          @negative-click="cancelImportTable"
        >
          <div class="py-2">
            <n-form>
              <n-form-item label="选择数据表">
                <n-select
                  v-model:value="selectedTable"
                  placeholder="请选择表"
                  :options="tableOptions"
                  :loading="tablesLoading"
                  @update:value="onTableChange"
                />
              </n-form-item>
            </n-form>

            <div v-if="selectedTable && columnsData.length > 0" class="mt-4">
              <h4 class="mb-2">表字段列表</h4>
              <n-data-table
                :columns="columnColumns"
                :data="columnsData"
                :pagination="false"
                :bordered="true"
              />
            </div>
          </div>
        </n-modal>

        <!-- 字段编辑对话框 -->
        <n-modal
          v-model:show="showFieldModal"
          title="编辑字段"
          preset="dialog"
          positive-text="确定"
          negative-text="取消"
          @positive-click="confirmField"
          @negative-click="cancelField"
        >
          <n-form
            ref="fieldFormRef"
            :model="currentField"
            label-placement="left"
            label-width="100"
            require-mark-placement="right-hanging"
            :rules="fieldRules"
          >
            <n-form-item label="字段名称" path="fieldName">
              <n-input v-model:value="currentField.fieldName" placeholder="请输入字段名称，如：Name" />
            </n-form-item>

            <n-form-item label="是否关系字段">
              <n-switch v-model:value="currentField.isRelation" />
            </n-form-item>

            <!-- 非关系字段 -->
            <template v-if="!currentField.isRelation">
              <n-form-item label="字段类型" path="fieldType">
                <n-select
                  v-model:value="currentField.fieldType"
                  :options="fieldTypeOptions"
                  placeholder="请选择字段类型"
                />
              </n-form-item>

              <n-form-item label="数据库字段" path="columnName">
                <n-input-group>
                  <n-input 
                    v-model:value="currentField.columnName" 
                    placeholder="请输入数据库字段名，如：name"
                    style="width: calc(100% - 80px)"
                  />
                  <n-button type="primary" style="width: 80px" @click="generateColumnName">自动生成</n-button>
                </n-input-group>
              </n-form-item>

              <n-form-item label="字段描述" path="fieldDesc">
                <n-input v-model:value="currentField.fieldDesc" placeholder="请输入字段描述，如：名称" />
              </n-form-item>

              <n-grid :cols="2" :x-gap="12">
                <n-grid-item>
                  <n-form-item label="是否必填">
                    <n-switch v-model:value="currentField.required" />
                  </n-form-item>
                </n-grid-item>

                <n-grid-item>
                  <n-form-item label="是否主键">
                    <n-switch v-model:value="currentField.isPrimaryKey" />
                  </n-form-item>
                </n-grid-item>

                <n-grid-item>
                  <n-form-item label="可搜索">
                    <n-switch v-model:value="currentField.isSearchable" />
                  </n-form-item>
                </n-grid-item>

                <n-grid-item>
                  <n-form-item label="可过滤">
                    <n-switch v-model:value="currentField.isFilterable" />
                  </n-form-item>
                </n-grid-item>

                <n-grid-item>
                  <n-form-item label="可排序">
                    <n-switch v-model:value="currentField.isSortable" />
                  </n-form-item>
                </n-grid-item>
              </n-grid>
            </template>

            <!-- 关系字段 -->
            <template v-else>
              <n-form-item label="关系类型" path="relationType">
                <n-select
                  v-model:value="currentField.relationType"
                  :options="relationTypeOptions"
                  placeholder="请选择关系类型"
                />
              </n-form-item>

              <n-form-item label="关联模型" path="relatedModel">
                <n-input v-model:value="currentField.relatedModel" placeholder="请输入关联模型名，如：User" />
              </n-form-item>

              <n-form-item label="字段描述" path="fieldDesc">
                <n-input v-model:value="currentField.fieldDesc" placeholder="请输入字段描述，如：用户" />
              </n-form-item>
              
              <n-form-item>
                <n-button type="primary" @click="generateRelationFields">自动生成关系字段</n-button>
              </n-form-item>

              <n-form-item label="外键字段" path="foreignKey">
                <n-input 
                  v-model:value="currentField.foreignKey" 
                  placeholder="请输入外键字段名，如：UserID" 
                />
              </n-form-item>

              <n-form-item label="引用字段" path="references">
                <n-input 
                  v-model:value="currentField.references" 
                  placeholder="请输入引用字段名，如：ID" 
                />
              </n-form-item>

              <template v-if="currentField.relationType === 'many_to_many'">
                <n-form-item label="关联表名" path="joinTable">
                  <n-input 
                    v-model:value="currentField.joinTable" 
                    placeholder="请输入关联表名，如：user_roles" 
                  />
                </n-form-item>
              </template>

              <n-form-item label="是否预加载">
                <n-switch v-model:value="currentField.preload" />
              </n-form-item>
              
              <n-form-item label="支持JOIN查询">
                <n-switch v-model:value="currentField.joinable" />
              </n-form-item>
              
              <template v-if="currentField.joinable">
                <n-form-item label="JOIN条件字段" path="joinCondition">
                  <n-input 
                    v-model:value="currentField.joinCondition" 
                    placeholder="请输入JOIN条件字段名，如：name" 
                  />
                </n-form-item>
                
                <n-form-item label="过滤条件字段" path="filterCondition">
                  <n-input 
                    v-model:value="currentField.filterCondition" 
                    placeholder="请输入过滤条件字段名，如：name" 
                  />
                </n-form-item>
              </template>

              <n-grid :cols="2" :x-gap="12">
                <n-grid-item>
                  <n-form-item label="可搜索">
                    <n-switch v-model:value="currentField.isSearchable" />
                  </n-form-item>
                </n-grid-item>

                <n-grid-item>
                  <n-form-item label="可过滤">
                    <n-switch v-model:value="currentField.isFilterable" />
                  </n-form-item>
                </n-grid-item>
              </n-grid>
            </template>
          </n-form>
        </n-modal>
      </n-spin>
    </n-drawer-content>
  </n-drawer>
</template>

<script lang="ts" setup>
import { computed, h, nextTick, onMounted, reactive, ref, watch } from 'vue';
import { useMessage } from 'naive-ui';
import type { FormInst, FormRules, DataTableColumns, DataTableBaseColumn } from 'naive-ui';
import { generateCode as apiGenerateCode, getAppList, getColumns, getTables, type ColumnInfo, type FieldInfo } from '@/service/api/codegen';

// Props 定义
const props = defineProps({
  isVisible: {
    type: Boolean,
    default: false
  }
});

// 事件定义
const emit = defineEmits(['cancel', 'success']);

// 本地状态
const localVisible = ref(props.isVisible);
const loading = ref(false);
const title = ref('代码生成器');
const currentStep = ref(0);
const currentStatus = ref<'process' | 'error' | 'finish'>('process');

// 应用列表
const appOptions = ref<{ label: string; value: string }[]>([]);
const appsLoading = ref(false);

// 数据表相关
const showImportTableModal = ref(false);
const tablesLoading = ref(false);
const tableOptions = ref<{ label: string; value: string }[]>([]);
const selectedTable = ref('');
const columnsLoading = ref(false);
const columnsData = ref<ColumnInfo[]>([]);

// 表单实例
const baseFormRef = ref<FormInst | null>(null);
const optionsFormRef = ref<FormInst | null>(null);
const fieldFormRef = ref<FormInst | null>(null);

// 字段相关
const showFieldModal = ref(false);
const fieldEditIndex = ref(-1);
const showNewAppField = ref(false);
const currentField = reactive<FieldInfo & { _id?: string }>({
  fieldName: '',
  fieldType: 'string',
  columnName: '',
  fieldDesc: '',
  required: false,
  isPrimaryKey: false,
  isSearchable: false,
  isFilterable: false,
  isSortable: false,
  joinable: false,
  joinCondition: '',
  filterCondition: '',
  _id: ''
});

// 字段类型选项
const fieldTypeOptions = [
  { label: 'string', value: 'string' },
  { label: 'int', value: 'int' },
  { label: 'uint', value: 'uint' },
  { label: 'int64', value: 'int64' },
  { label: 'uint64', value: 'uint64' },
  { label: 'float64', value: 'float64' },
  { label: 'bool', value: 'bool' },
  { label: 'time.Time', value: 'time.Time' }
];

// 关系类型选项
const relationTypeOptions = [
  { label: '从属于(BelongsTo)', value: 'belongs_to' },
  { label: '拥有一个(HasOne)', value: 'has_one' },
  { label: '拥有多个(HasMany)', value: 'has_many' },
  { label: '多对多(ManyToMany)', value: 'many_to_many' }
];

// 表单数据
const formData = reactive({
  appName: '',
  newAppName: '',
  packageName: 'admin',
  structName: '',
  tableName: '',
  description: '',
  apiPrefix: '',
  hasList: true,
  hasCreate: true,
  hasUpdate: true,
  hasDelete: true,
  hasDetail: true,
  hasPagination: true,
  fields: [] as (FieldInfo & { _id: string })[]
});

// 表单验证规则
const baseRules: FormRules = {
  appName: {
    required: true,
    trigger: ['blur', 'input'],
    validator(rule, value) {
      if (!value && !formData.newAppName) {
        return new Error('请选择应用或创建新应用');
      }
      return true;
    }
  },
  newAppName: {
    trigger: ['blur', 'input'],
    validator(rule, value) {
      if (showNewAppField.value && !value) {
        return new Error('请输入新应用名称');
      }
      if (showNewAppField.value && !/^[a-zA-Z][a-zA-Z0-9_]*$/.test(value)) {
        return new Error('应用名称只能包含字母、数字和下划线，且以字母开头');
      }
      return true;
    }
  },
  packageName: {
    required: true,
    trigger: ['blur', 'input'],
    message: '请输入包名'
  },
  structName: {
    required: true,
    trigger: ['blur', 'input'],
    validator(rule, value) {
      if (!value) {
        return new Error('请输入结构体名称');
      }
      if (!/^[A-Z][a-zA-Z0-9]*$/.test(value)) {
        return new Error('结构体名称应以大写字母开头，只能包含字母和数字');
      }
      return true;
    }
  },
  tableName: {
    required: true,
    trigger: ['blur', 'input'],
    message: '请输入表名'
  },
  description: {
    required: true,
    trigger: ['blur', 'input'],
    message: '请输入描述'
  }
};

const fieldRules: FormRules = {
  fieldName: {
    required: true,
    trigger: ['blur', 'input'],
    validator(rule, value) {
      if (!value) {
        return new Error('请输入字段名称');
      }
      if (!/^[A-Z][a-zA-Z0-9]*$/.test(value)) {
        return new Error('字段名称应以大写字母开头，只能包含字母和数字');
      }
      return true;
    }
  },
  fieldType: {
    required: true,
    trigger: ['blur', 'change'],
    message: '请选择字段类型'
  },
  columnName: {
    required: true,
    trigger: ['blur', 'input'],
    validator(rule, value) {
      if (!value && !currentField.isRelation) {
        return new Error('请输入数据库字段名');
      }
      if (value && !/^[a-z][a-z0-9_]*$/.test(value)) {
        return new Error('数据库字段名应以小写字母开头，只能包含小写字母、数字和下划线');
      }
      return true;
    }
  },
  fieldDesc: {
    required: true,
    trigger: ['blur', 'input'],
    message: '请输入字段描述'
  },
  // 关系字段验证
  relationType: {
    required: true,
    trigger: ['blur', 'change'],
    validator(rule, value) {
      if (currentField.isRelation && !value) {
        return new Error('请选择关系类型');
      }
      return true;
    }
  },
  relatedModel: {
    required: true,
    trigger: ['blur', 'input'],
    validator(rule, value) {
      if (currentField.isRelation && !value) {
        return new Error('请输入关联模型');
      }
      if (value && !/^[A-Z][a-zA-Z0-9]*$/.test(value)) {
        return new Error('关联模型名称应以大写字母开头，只能包含字母和数字');
      }
      return true;
    }
  },
  foreignKey: {
    trigger: ['blur', 'input'],
    validator(rule, value) {
      if (currentField.isRelation && 
          (currentField.relationType === 'belongs_to' || 
           currentField.relationType === 'has_one' || 
           currentField.relationType === 'has_many')) {
        if (value && !/^[A-Z][a-zA-Z0-9]*$/.test(value)) {
          return new Error('外键字段名应以大写字母开头，只能包含字母和数字');
        }
      }
      return true;
    }
  },
  joinTable: {
    trigger: ['blur', 'input'],
    validator(rule, value) {
      if (currentField.isRelation && currentField.relationType === 'many_to_many' && !value) {
        return new Error('多对多关系必须指定关联表名');
      }
      if (value && !/^[a-z][a-z0-9_]*$/.test(value)) {
        return new Error('关联表名应以小写字母开头，只能包含小写字母、数字和下划线');
      }
      return true;
    }
  }
};

// 消息组件
const message = useMessage();

// 工具函数：大驼峰转小驼峰
const toLowerCamel = (str: string) => {
  if (!str) return '';
  return str.charAt(0).toLowerCase() + str.slice(1);
};

// 工具函数：大驼峰转下划线
const toSnakeCase = (str: string) => {
  if (!str) return '';
  return str
    .replace(/([A-Z])/g, '_$1')
    .toLowerCase()
    .replace(/^_/, '');
};

// 监听结构体名称变化，自动填充表名和API前缀
watch(
  () => formData.structName,
  (newValue) => {
    // 移除自动填充，让用户自由决定是否使用自动生成的值
  }
);

// 移除字段名称自动转换监听器
// 改为添加辅助函数供模板中的按钮使用

// 自动生成表名函数（供按钮使用）
function generateTableName() {
  if (formData.structName) {
    formData.tableName = toLowerCamel(formData.structName) + 's';
  }
}

// 自动生成API前缀函数
function generateApiPrefix() {
  if (formData.structName) {
    formData.apiPrefix = toLowerCamel(formData.structName);
  }
}

// 自动生成描述函数
function generateDescription() {
  if (formData.structName) {
    formData.description = formData.structName.replace(/([A-Z])/g, ' $1').trim();
  }
}

// 自动生成数据库字段名函数
function generateColumnName() {
  if (currentField.fieldName && !currentField.isRelation) {
    currentField.columnName = toSnakeCase(currentField.fieldName);
  }
}

// 监听关系字段类型和关联模型变化，自动生成外键和引用字段
watch(
  [() => currentField.isRelation, () => currentField.relationType, () => currentField.relatedModel],
  ([isRelation, relationType, relatedModel]) => {
    // 移除自动生成，让用户手动触发
  }
);

// 自动生成关系字段的外键和引用字段
function generateRelationFields() {
  if (!currentField.isRelation || !currentField.relatedModel) return;

  const relationType = currentField.relationType;
  const relatedModel = currentField.relatedModel;

  if (relationType === 'belongs_to') {
    // 从属关系：自动生成外键 (RelatedModelID)
    currentField.foreignKey = relatedModel + 'ID';
    currentField.references = 'ID';
  } else if (relationType === 'has_one' || relationType === 'has_many') {
    // 拥有关系：自动生成当前结构体ID作为外键
    if (formData.structName) {
      currentField.foreignKey = formData.structName + 'ID';
      currentField.references = 'ID';
    }
  } else if (relationType === 'many_to_many') {
    // 多对多关系：自动生成关联表名
    if (formData.structName) {
      const names = [formData.structName, relatedModel].sort();
      currentField.joinTable = toSnakeCase(names.join('')).toLowerCase();
    }
  }
}

// 字段表格列定义
const fieldColumns = computed(() => {
  return [
    { 
      title: '字段名称',
      key: 'fieldName',
      width: 120
    },
    {
      title: '字段类型',
      key: 'fieldType',
      width: 100,
      render(row) {
        if (row.isRelation) {
          switch (row.relationType) {
            case 'belongs_to':
              return row.relatedModel;
            case 'has_one':
              return row.relatedModel;
            case 'has_many':
              return `[]${row.relatedModel}`;
            case 'many_to_many':
              return `[]${row.relatedModel}`;
            default:
              return row.fieldType;
          }
        } else {
          return row.fieldType;
        }
      }
    },
    {
      title: '关系类型',
      key: 'relationType',
      width: 120,
      render(row) {
        if (!row.isRelation) return '';
        
        switch (row.relationType) {
          case 'belongs_to':
            return '从属于';
          case 'has_one':
            return '拥有一个';
          case 'has_many':
            return '拥有多个';
          case 'many_to_many':
            return '多对多';
          default:
            return '';
        }
      }
    },
    {
      title: '数据库字段',
      key: 'columnName',
      width: 120
    },
    {
      title: '字段描述',
      key: 'fieldDesc',
      width: 120
    },
    {
      title: '必填',
      key: 'required',
      width: 80,
      render(row) {
        return row.required ? '是' : '否';
      }
    },
    {
      title: '主键',
      key: 'isPrimaryKey',
      width: 80,
      render(row) {
        return row.isPrimaryKey ? '是' : '否';
      }
    },
    {
      title: '可搜索',
      key: 'isSearchable',
      width: 80,
      render(row) {
        return row.isSearchable ? '是' : '否';
      }
    },
    {
      title: '可过滤',
      key: 'isFilterable',
      width: 80,
      render(row) {
        return row.isFilterable ? '是' : '否';
      }
    },
    {
      title: '操作',
      key: 'actions',
      fixed: 'right',
      width: 120,
      render(row, index) {
        return h('div', { class: 'flex gap-2' }, [
          h(
            'button',
            {
              class: 'text-blue-500 hover:text-blue-700',
              onClick: () => editField(index)
            },
            '编辑'
          ),
          h(
            'button',
            {
              class: 'text-red-500 hover:text-red-700',
              onClick: () => deleteField(index)
            },
            '删除'
          )
        ]);
      }
    }
  ] as DataTableColumns;
});

// 数据表列定义
const columnColumns: DataTableColumns = [
  { 
    title: '字段名',
    key: 'columnName',
    width: 120
  },
  {
    title: '数据类型',
    key: 'dataType',
    width: 100
  },
  {
    title: '字段注释',
    key: 'columnComment',
    width: 150
  },
  {
    title: '是否可为空',
    key: 'isNullable',
    width: 100
  },
  {
    title: '主键',
    key: 'columnKey',
    width: 80,
    render(row) {
      return row.columnKey === 'PRI' ? '是' : '否';
    }
  }
];

// 监听props变化
watch(
  () => props.isVisible,
  (val) => {
    localVisible.value = val;
    if (val) {
      currentStep.value = 0;
      currentStatus.value = 'process';
      resetForm();
      loadAppList();
    }
  }
);

// 加载应用列表
async function loadAppList() {
  try {
    appsLoading.value = true;
    const res = await getAppList();
    if (res.data) {
      appOptions.value = res.data.map(app => ({
        label: app,
        value: app
      }));
    }
  } catch (error) {
    message.error('获取应用列表失败');
  } finally {
    appsLoading.value = false;
  }
}

// 加载表列表
async function loadTables() {
  try {
    tablesLoading.value = true;
    const res = await getTables();
    if (res.data) {
      console.log('获取到表格数据:', JSON.stringify(res.data));
      
      tableOptions.value = res.data.map(table => {
        // 确保表名是字符串
        const tableName = String(table.tableName || '');
        const tableComment = String(table.tableComment || '');
        
        return {
          label: tableComment ? `${tableName} (${tableComment})` : tableName,
          value: tableName
        };
      });
      
      console.log('处理后的表格选项:', tableOptions.value);
    }
  } catch (error) {
    console.error('获取表列表失败:', error);
    message.error('获取表列表失败');
  } finally {
    tablesLoading.value = false;
  }
}

// 加载表字段
async function loadColumns(tableName: string) {
  try {
    columnsLoading.value = true;
    const res = await getColumns(tableName);
    if (res.data) {
      columnsData.value = res.data;
    }
  } catch (error) {
    message.error('获取表字段失败');
  } finally {
    columnsLoading.value = false;
  }
}

// 应用选择变化
function onAppChange(value: string) {
  if (value) {
    showNewAppField.value = false;
  }
  // 自动设置包名
  if (value && !formData.packageName) {
    formData.packageName = value;
  }
}

// 表选择变化
function onTableChange(value: string) {
  if (value) {
    loadColumns(value);
    // 自动设置结构体名和表名
    if (!formData.tableName) {
      formData.tableName = value;
    }
  }
}

// 从数据表导入字段
function importFromTable() {
  showImportTableModal.value = true;
  loadTables();
}

// 确认导入表字段
function confirmImportTable() {
  if (!selectedTable.value || columnsData.value.length === 0) {
    message.warning('请选择表并加载字段');
    return;
  }

  // 使用表信息填充表单
  if (!formData.tableName) {
    formData.tableName = selectedTable.value;
  }

  // 根据表名生成结构体名
  if (!formData.structName) {
    // 转换表名为大驼峰结构体名
    let structName = selectedTable.value
      .split('_')
      .map(part => part.charAt(0).toUpperCase() + part.slice(1))
      .join('');
    // 如果是复数形式，转为单数
    if (structName.endsWith('s')) {
      structName = structName.slice(0, -1);
    }
    formData.structName = structName;
    
    // 自动设置API前缀（小驼峰）
    if (!formData.apiPrefix) {
      formData.apiPrefix = toLowerCamel(structName);
    }
  }

  // 字段映射
  const fields: (FieldInfo & { _id: string })[] = columnsData.value.map(column => {
    // 将数据库类型转换为Go类型
    let fieldType = 'string';
    if (column.dataType.includes('int')) {
      fieldType = 'int';
    } else if (column.dataType.includes('float') || column.dataType.includes('double')) {
      fieldType = 'float64';
    } else if (column.dataType.includes('bool')) {
      fieldType = 'bool';
    } else if (column.dataType.includes('datetime') || column.dataType.includes('timestamp')) {
      fieldType = 'time.Time';
    }

    // 字段名转为大驼峰
    const fieldName = column.columnName
      .split('_')
      .map(part => part.charAt(0).toUpperCase() + part.slice(1))
      .join('');

    // 智能处理字段描述
    let fieldDesc = column.columnComment || '';
    if (!fieldDesc) {
      // 如果没有注释，根据字段名生成描述
      fieldDesc = column.columnName
        .split('_')
        .map(part => part.charAt(0).toUpperCase() + part.slice(1))
        .join(' ');
    }

    // 判断是否为外键关系字段
    const isForeignKey = column.columnName.endsWith('_id') && column.columnName !== 'id';
    const isPotentialRelation = isForeignKey && !column.columnKey;

    return {
      fieldName,
      fieldType,
      columnName: column.columnName,
      fieldDesc: fieldDesc,
      required: column.isNullable === 'NO',
      isPrimaryKey: column.columnKey === 'PRI',
      isSearchable: column.columnKey === 'PRI' || column.columnName.includes('name') || column.columnName.includes('title'),
      isFilterable: column.columnKey === 'PRI' || column.dataType === 'enum' || column.columnName.includes('status'),
      isSortable: column.columnKey === 'PRI' || column.columnName.includes('time') || column.columnName.includes('date'),
      _id: `field_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`
    };
  });

  formData.fields = fields;
  showImportTableModal.value = false;
  message.success('导入字段成功');
}

// 取消导入表字段
function cancelImportTable() {
  showImportTableModal.value = false;
}

// 添加字段
function addField() {
  showFieldModal.value = true;
  fieldEditIndex.value = -1;
  
  // 重置当前字段
  Object.assign(currentField, {
    fieldName: '',
    fieldType: 'string',
    columnName: '',
    fieldDesc: '',
    required: false,
    isPrimaryKey: false,
    isSearchable: false,
    isFilterable: false,
    isSortable: false,
    
    // 关系字段
    isRelation: false,
    relationType: 'belongs_to',
    relatedModel: '',
    foreignKey: '',
    references: '',
    preload: false,
    joinTable: '',
    joinable: false,
    joinCondition: '',
    filterCondition: '',

    _id: `field_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`
  });
  
  nextTick(() => {
    fieldFormRef.value?.restoreValidation();
  });
}

// 编辑字段
function editField(index: number) {
  const field = formData.fields[index];
  if (!field) return;
  
  showFieldModal.value = true;
  fieldEditIndex.value = index;
  
  // 复制字段数据
  Object.assign(currentField, field);
  
  nextTick(() => {
    fieldFormRef.value?.restoreValidation();
  });
}

// 删除字段
function deleteField(index: number) {
  formData.fields.splice(index, 1);
}

// 确认字段编辑
function confirmField() {
  fieldFormRef.value?.validate((errors) => {
    if (errors) return;
    
    // 创建基本字段对象
    const field = {
      fieldName: currentField.fieldName,
      fieldType: currentField.fieldType,
      columnName: currentField.columnName,
      fieldDesc: currentField.fieldDesc,
      required: currentField.required,
      isPrimaryKey: currentField.isPrimaryKey,
      isSearchable: currentField.isSearchable,
      isFilterable: currentField.isFilterable,
      isSortable: currentField.isSortable,
      _id: currentField._id || `field_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`
    } as (FieldInfo & { _id: string });
    
    // 添加关系字段属性
    if (currentField.isRelation) {
      field.isRelation = true;
      field.relationType = currentField.relationType;
      field.relatedModel = currentField.relatedModel;
      field.foreignKey = currentField.foreignKey;
      field.references = currentField.references;
      field.preload = currentField.preload;
      
      // 添加JOIN查询相关属性
      field.joinable = currentField.joinable;
      field.joinCondition = currentField.joinCondition;
      field.filterCondition = currentField.filterCondition;
      
      // 对于多对多关系，添加关联表属性
      if (currentField.relationType === 'many_to_many') {
        field.joinTable = currentField.joinTable;
      }
    }
    
    if (fieldEditIndex.value >= 0) {
      // 更新字段
      formData.fields[fieldEditIndex.value] = field;
    } else {
      // 添加新字段
      formData.fields.push(field);
    }
    
    showFieldModal.value = false;
  });
}

// 取消字段编辑
function cancelField() {
  showFieldModal.value = false;
}

// 上一步
function prevStep() {
  currentStep.value -= 1;
}

// 下一步
function nextStep() {
  if (currentStep.value === 0) {
    // 验证基本设置表单
    baseFormRef.value?.validate((errors) => {
      if (errors) {
        currentStatus.value = 'error';
        return;
      }
      
      // 如果没有填写API前缀，使用表名
      if (!formData.apiPrefix) {
        formData.apiPrefix = formData.tableName;
      }
      
      // 检查是否有字段
      if (formData.fields.length === 0) {
        // 添加一个ID字段作为主键
        formData.fields.push({
          fieldName: 'ID',
          fieldType: 'uint',
          columnName: 'id',
          fieldDesc: '主键ID',
          required: true,
          isPrimaryKey: true,
          isSearchable: true,
          isFilterable: true,
          isSortable: true,
          _id: `field_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`
        });
      }
      
      currentStep.value += 1;
      currentStatus.value = 'process';
    });
  } else if (currentStep.value === 1) {
    // 检查是否有字段
    if (formData.fields.length === 0) {
      message.error('请至少添加一个字段');
      currentStatus.value = 'error';
      return;
    }
    
    // 检查是否有主键
    const hasPrimaryKey = formData.fields.some(field => field.isPrimaryKey);
    if (!hasPrimaryKey) {
      message.error('请设置一个主键字段');
      currentStatus.value = 'error';
      return;
    }
    
    currentStep.value += 1;
    currentStatus.value = 'process';
  }
}

// 生成代码
async function generateCode() {
  try {
    loading.value = true;
    
    // 准备请求数据
    const requestData = {
      structName: formData.structName,
      tableName: formData.tableName,
      packageName: formData.packageName,
      description: formData.description,
      apiPrefix: formData.apiPrefix,
      appName: formData.newAppName || formData.appName,
      hasList: formData.hasList,
      hasCreate: formData.hasCreate,
      hasUpdate: formData.hasUpdate,
      hasDelete: formData.hasDelete,
      hasDetail: formData.hasDetail,
      hasPagination: formData.hasPagination,
      fields: formData.fields.map(field => ({
        fieldName: field.fieldName,
        fieldType: field.fieldType,
        columnName: field.columnName,
        fieldDesc: field.fieldDesc,
        required: field.required,
        isPrimaryKey: field.isPrimaryKey,
        isSearchable: field.isSearchable,
        isFilterable: field.isFilterable,
        isSortable: field.isSortable,
        // 添加关系字段属性
        isRelation: field.isRelation,
        relationType: field.relationType,
        relatedModel: field.relatedModel,
        foreignKey: field.foreignKey,
        references: field.references,
        preload: field.preload,
        joinTable: field.joinTable,
        joinable: field.joinable,
        joinCondition: field.joinCondition,
        filterCondition: field.filterCondition
      }))
    };
    
    await apiGenerateCode(requestData);
    message.success('代码生成成功');
    emit('success');
  } catch (error) {
    message.error('代码生成失败');
  } finally {
    loading.value = false;
  }
}

// 关闭
function handleClose() {
  emit('cancel');
}

// 重置表单
function resetForm() {
  Object.assign(formData, {
    appName: '',
    newAppName: '',
    packageName: 'admin',
    structName: '',
    tableName: '',
    description: '',
    apiPrefix: '',
    hasList: true,
    hasCreate: true,
    hasUpdate: true,
    hasDelete: true,
    hasDetail: true,
    hasPagination: true,
    fields: []
  });
  
  showNewAppField.value = false;
  
  nextTick(() => {
    baseFormRef.value?.restoreValidation();
    optionsFormRef.value?.restoreValidation();
  });
}

// 组件挂载时初始化
onMounted(() => {
  if (props.isVisible) {
    loadAppList();
  }
});
</script> 