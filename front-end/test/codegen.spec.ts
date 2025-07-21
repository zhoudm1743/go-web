import { describe, it, expect, beforeEach, vi } from 'vitest';
import { mount } from '@vue/test-utils';
import { createTestingPinia } from '@pinia/testing';
import { nextTick } from 'vue';

// 由于组件可能有依赖，这里需要模拟
vi.mock('@/service/api/codegen', () => {
  return {
    getAppList: vi.fn().mockResolvedValue({ data: ['admin', 'api'] }),
    getTables: vi.fn().mockResolvedValue({ 
      data: [
        { tableName: 'users', tableComment: '用户表' },
        { tableName: 'products', tableComment: '产品表' }
      ] 
    }),
    getColumns: vi.fn().mockResolvedValue({
      data: [
        { columnName: 'id', dataType: 'int', columnComment: '主键ID', isNullable: 'NO', columnKey: 'PRI' },
        { columnName: 'name', dataType: 'varchar', columnComment: '名称', isNullable: 'NO', columnKey: '' }
      ]
    }),
    generateCode: vi.fn().mockResolvedValue({ code: 200, message: '生成成功' }),
    getHistoryList: vi.fn().mockResolvedValue({
      data: {
        list: [
          { 
            id: 1, 
            createdAt: '2023-01-01', 
            updatedAt: '2023-01-01', 
            table: 'users', 
            structName: 'User', 
            packageName: 'admin', 
            moduleName: 'github.com/zhoudm1743/go-web', 
            description: '用户', 
            flag: 0 
          }
        ],
        total: 1,
        page: 1,
        pageSize: 10
      }
    }),
    rollbackCode: vi.fn().mockResolvedValue({ code: 200, message: '回滚成功' }),
    deleteHistory: vi.fn().mockResolvedValue({ code: 200, message: '删除成功' })
  };
});

// 模拟 naive-ui 组件
vi.mock('naive-ui', () => {
  return {
    useMessage: () => ({
      success: vi.fn(),
      error: vi.fn(),
      warning: vi.fn()
    }),
    // 其他需要的组件可以在这里添加
  };
});

// 模拟路由
vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: vi.fn(),
    replace: vi.fn()
  })
}));

// 导入要测试的组件
// 注意：实际测试时需要根据项目结构导入正确的路径
import CodegenIndex from '@/views/setting/codegen/index.vue';
import GeneratorForm from '@/views/setting/codegen/components/GeneratorForm.vue';
import HistoryList from '@/views/setting/codegen/components/HistoryList.vue';

describe('代码生成可视化测试', () => {
  // 测试主页面
  describe('Codegen Index', () => {
    it('应该渲染代码生成器和历史记录选项卡', async () => {
      const wrapper = mount(CodegenIndex, {
        global: {
          plugins: [createTestingPinia({ createSpy: vi.fn })],
          stubs: {
            'n-card': true,
            'n-tabs': true,
            'n-tab-pane': true,
            'n-button': true,
            'n-result': true,
            'icon-mdi:plus-circle': true,
            'GeneratorForm': true,
            'HistoryList': true
          }
        }
      });

      // 确保主要组件被渲染
      expect(wrapper.findComponent({ name: 'n-tabs' }).exists()).toBe(true);
      expect(wrapper.findComponent({ name: 'n-tab-pane', props: { name: 'generator' } }).exists()).toBe(true);
      expect(wrapper.findComponent({ name: 'n-tab-pane', props: { name: 'history' } }).exists()).toBe(true);
    });

    it('点击创建按钮应该显示表单', async () => {
      const wrapper = mount(CodegenIndex, {
        global: {
          plugins: [createTestingPinia({ createSpy: vi.fn })],
          stubs: {
            'n-card': true,
            'n-tabs': true,
            'n-tab-pane': true,
            'n-button': { template: '<button @click="$emit(\'click\')"><slot /></button>' },
            'n-result': true,
            'icon-mdi:plus-circle': true,
            'GeneratorForm': true,
            'HistoryList': true
          }
        }
      });

      // 初始不应该显示表单
      expect(wrapper.vm.showGeneratorForm).toBe(false);
      
      // 点击创建按钮
      await wrapper.findComponent('button').trigger('click');
      
      // 表单应该显示
      expect(wrapper.vm.showGeneratorForm).toBe(true);
    });
  });

  // 测试生成器表单组件
  describe('Generator Form', () => {
    let wrapper;
    
    beforeEach(() => {
      wrapper = mount(GeneratorForm, {
        props: {
          isVisible: true
        },
        global: {
          plugins: [createTestingPinia({ createSpy: vi.fn })],
          stubs: {
            // 这里列出所有GeneratorForm中使用的naive-ui组件
            'n-drawer': true,
            'n-drawer-content': true,
            'n-form': true,
            'n-form-item': true,
            'n-input': true,
            'n-button': true,
            'n-select': true,
            'n-steps': true,
            'n-step': true,
            'n-spin': true,
            'n-switch': true,
            'n-grid': true,
            'n-grid-item': true,
            'n-modal': true,
            'n-data-table': true,
            'n-alert': true,
            'n-icon': true,
            'icon-ant-design:info-circle-outlined': true
          }
        }
      });
    });

    it('应该正确加载应用列表', async () => {
      // 测试应用列表加载
      expect(wrapper.vm.appOptions).toHaveLength(2); // 应该有两个应用选项
      expect(wrapper.vm.appOptions[0].value).toBe('admin');
      expect(wrapper.vm.appOptions[1].value).toBe('api');
    });

    it('应该能正确填写表单并提交', async () => {
      // 填写表单
      wrapper.vm.formData.structName = 'TestModel';
      wrapper.vm.formData.tableName = 'test_models';
      wrapper.vm.formData.description = '测试模型';
      wrapper.vm.formData.apiPrefix = 'testModel';
      wrapper.vm.formData.appName = 'admin';
      wrapper.vm.formData.packageName = 'admin';
      
      // 添加字段
      wrapper.vm.formData.fields = [
        {
          fieldName: 'ID',
          fieldType: 'uint',
          columnName: 'id',
          fieldDesc: '主键ID',
          required: true,
          isPrimaryKey: true,
          isSearchable: true,
          isFilterable: true,
          isSortable: true,
          _id: 'test_id_1'
        },
        {
          fieldName: 'Name',
          fieldType: 'string',
          columnName: 'name',
          fieldDesc: '名称',
          required: true,
          isPrimaryKey: false,
          isSearchable: true,
          isFilterable: true,
          isSortable: false,
          _id: 'test_id_2'
        }
      ];
      
      // 进入下一步
      wrapper.vm.currentStep = 1; // 字段设置
      await nextTick();
      wrapper.vm.nextStep(); // 进入最后一步
      await nextTick();
      
      // 生成代码
      wrapper.vm.generateCode();
      await nextTick();
      
      // 验证generateCode方法被调用
      const { generateCode } = await import('@/service/api/codegen');
      expect(generateCode).toHaveBeenCalledWith(expect.objectContaining({
        structName: 'TestModel',
        tableName: 'test_models',
        description: '测试模型',
        apiPrefix: 'testModel',
        fields: expect.arrayContaining([
          expect.objectContaining({
            fieldName: 'ID',
            fieldType: 'uint'
          }),
          expect.objectContaining({
            fieldName: 'Name',
            fieldType: 'string'
          })
        ])
      }));
    });

    it('应该能正确处理关系字段', async () => {
      // 添加一个关系字段
      const relationField = {
        fieldName: 'Category',
        fieldType: 'string', // 会被覆盖
        columnName: '',
        fieldDesc: '分类',
        required: false,
        isPrimaryKey: false,
        isSearchable: false,
        isFilterable: true,
        isSortable: false,
        isRelation: true,
        relationType: 'belongs_to',
        relatedModel: 'Category',
        preload: true,
        joinable: true,
        _id: 'test_relation_id'
      };
      
      // 使用编辑字段功能测试
      wrapper.vm.currentField = { ...relationField };
      wrapper.vm.confirmField();
      await nextTick();
      
      // 验证字段是否添加到表单数据中
      const fieldExists = wrapper.vm.formData.fields.some(f => 
        f.fieldName === 'Category' && 
        f.isRelation === true && 
        f.relationType === 'belongs_to' &&
        f.preload === true &&
        f.joinable === true
      );
      
      expect(fieldExists).toBe(true);
    });

    it('应该能使用自动生成按钮', async () => {
      // 设置结构体名称
      wrapper.vm.formData.structName = 'Product';
      await nextTick();
      
      // 测试自动生成表名
      wrapper.vm.generateTableName();
      expect(wrapper.vm.formData.tableName).toBe('products');
      
      // 测试自动生成API前缀
      wrapper.vm.generateApiPrefix();
      expect(wrapper.vm.formData.apiPrefix).toBe('product');
      
      // 测试自动生成描述
      wrapper.vm.generateDescription();
      expect(wrapper.vm.formData.description).toBe('Product');
    });
  });

  // 测试历史记录组件
  describe('History List', () => {
    let wrapper;
    
    beforeEach(() => {
      wrapper = mount(HistoryList, {
        props: {
          onRefresh: vi.fn()
        },
        global: {
          plugins: [createTestingPinia({ createSpy: vi.fn })],
          stubs: {
            'n-card': true,
            'n-data-table': true,
            'n-modal': true,
            'n-checkbox': true,
            'n-tabs': true,
            'n-tab-pane': true,
            'n-code': true,
            'n-button': true
          }
        }
      });
    });

    it('应该正确加载历史记录', async () => {
      // 等待历史记录加载
      await wrapper.vm.loadHistoryData();
      await nextTick();
      
      // 验证getHistoryList方法被调用
      const { getHistoryList } = await import('@/service/api/codegen');
      expect(getHistoryList).toHaveBeenCalled();
      
      // 验证历史记录数据
      expect(wrapper.vm.historyData).toHaveLength(1);
      expect(wrapper.vm.historyData[0].structName).toBe('User');
    });

    it('应该能处理回滚操作', async () => {
      // 设置回滚选项
      wrapper.vm.rollbackOptions = {
        deleteFiles: true,
        deleteApi: false,
        deleteMenu: false,
        deleteTable: false
      };
      
      // 设置要回滚的历史记录
      wrapper.vm.currentHistoryId = 1;
      
      // 调用回滚方法
      await wrapper.vm.confirmRollback();
      await nextTick();
      
      // 验证rollbackCode方法被调用
      const { rollbackCode } = await import('@/service/api/codegen');
      expect(rollbackCode).toHaveBeenCalledWith({
        id: 1,
        deleteFiles: true,
        deleteApi: false,
        deleteMenu: false,
        deleteTable: false
      });
    });
  });
});

// 测试工具函数
describe('工具函数测试', () => {
  it('应该正确转换大驼峰为小驼峰', () => {
    const { toLowerCamel } = await import('@/views/setting/codegen/components/GeneratorForm.vue');
    expect(toLowerCamel('UserInfo')).toBe('userInfo');
    expect(toLowerCamel('Product')).toBe('product');
  });
  
  it('应该正确转换大驼峰为下划线形式', () => {
    const { toSnakeCase } = await import('@/views/setting/codegen/components/GeneratorForm.vue');
    expect(toSnakeCase('UserInfo')).toBe('user_info');
    expect(toSnakeCase('ProductCategory')).toBe('product_category');
  });
}); 