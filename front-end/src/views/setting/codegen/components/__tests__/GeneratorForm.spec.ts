import { describe, it, expect, beforeEach, vi } from 'vitest';
import { mount } from '@vue/test-utils';
import { NButton, NSwitch } from 'naive-ui';
import { nextTick } from 'vue';
import GeneratorForm from '../GeneratorForm.vue';

// Mock API调用
vi.mock('@/service/api/codegen', () => ({
  getAppList: vi.fn().mockResolvedValue({ code: 0, data: ['admin', 'api'] }),
  getTables: vi.fn().mockResolvedValue({ 
    code: 0, 
    data: [
      { tableName: 'users', tableComment: '用户表' },
      { tableName: 'products', tableComment: '产品表' }
    ] 
  }),
  getColumns: vi.fn().mockResolvedValue({
    code: 0, 
    data: [
      { columnName: 'id', dataType: 'int', columnComment: '主键ID', isNullable: 'NO', columnKey: 'PRI' },
      { columnName: 'name', dataType: 'varchar', columnComment: '名称', isNullable: 'NO', columnKey: '' }
    ]
  }),
  generateCode: vi.fn().mockResolvedValue({ code: 0, message: '生成成功' })
}));

describe('GeneratorForm组件', () => {
  // 测试辅助函数
  function createWrapper(props = {}) {
    return mount(GeneratorForm, {
      props: {
        isVisible: true,
        ...props
      },
      global: {
        stubs: {
          // 存根Naive UI组件
          'n-drawer': true,
          'n-drawer-content': true,
          'n-form': true,
          'n-form-item': true,
          'n-input': true,
          'n-button': NButton,
          'n-select': true,
          'n-steps': true,
          'n-step': true,
          'n-spin': true,
          'n-switch': NSwitch,
          'n-grid': true,
          'n-grid-item': true,
          'n-modal': true,
          'n-data-table': true,
          'n-alert': true,
          'n-icon': true,
          'n-space': true,
          'icon-ant-design:info-circle-outlined': true
        }
      }
    });
  }

  it('应该正确初始化', async () => {
    const wrapper = createWrapper();
    
    // 验证组件已经渲染
    expect(wrapper.exists()).toBe(true);
    
    // 验证初始状态
    expect(wrapper.vm.currentStep).toBe(0);
    expect(wrapper.vm.formData.fields.length).toBe(0);
  });

  it('应该能正确设置结构体名称和表名', async () => {
    const wrapper = createWrapper();
    
    // 设置结构体名称
    await wrapper.vm.formData.structName = 'Product';
    
    // 测试自动生成表名
    await wrapper.vm.generateTableName();
    expect(wrapper.vm.formData.tableName).toBe('products');
  });

  it('应该正确处理普通字段', async () => {
    const wrapper = createWrapper();
    
    // 添加普通字段
    const field = {
      fieldName: 'Name',
      fieldType: 'string',
      columnName: 'name',
      fieldDesc: '名称',
      required: true,
      isPrimaryKey: false,
      isSearchable: true,
      isFilterable: true,
      isSortable: false
    };
    
    wrapper.vm.currentField = { ...field };
    wrapper.vm.confirmField();
    
    // 验证字段是否被添加
    const addedField = wrapper.vm.formData.fields.find(f => f.fieldName === 'Name');
    expect(addedField).toBeDefined();
    expect(addedField?.fieldType).toBe('string');
    expect(addedField?.isSearchable).toBe(true);
  });

  it('应该正确处理关系字段', async () => {
    const wrapper = createWrapper();
    
    // 添加关系字段
    const relationField = {
      fieldName: 'Category',
      fieldType: 'string',
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
      foreignKey: 'CategoryID',
      references: 'ID',
      preload: true,
      joinable: true
    };
    
    wrapper.vm.currentField = { ...relationField };
    wrapper.vm.confirmField();
    
    // 验证关系字段
    const addedField = wrapper.vm.formData.fields.find(f => f.fieldName === 'Category');
    expect(addedField).toBeDefined();
    expect(addedField?.isRelation).toBe(true);
    expect(addedField?.preload).toBe(true);
    expect(addedField?.joinable).toBe(true);
  });

  it('应该能生成相应的前端代码', async () => {
    const wrapper = createWrapper();
    
    // 填写表单
    wrapper.vm.formData.structName = 'Product';
    wrapper.vm.formData.tableName = 'products';
    wrapper.vm.formData.packageName = 'admin';
    wrapper.vm.formData.description = '产品';
    wrapper.vm.formData.apiPrefix = 'product';
    
    // 添加字段
    const fields = [
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
      }
    ];
    
    fields.forEach(field => {
      wrapper.vm.currentField = { ...field };
      wrapper.vm.confirmField();
    });
    
    // 进行到最后一步
    wrapper.vm.currentStep = 2;
    
    // 触发代码生成
    await wrapper.vm.generateCode();
    
    // 验证API调用
    const { generateCode } = await import('@/service/api/codegen');
    expect(generateCode).toHaveBeenCalled();
    
    // 验证传递给API的数据
    const callData = vi.mocked(generateCode).mock.calls[0][0];
    expect(callData.structName).toBe('Product');
    expect(callData.tableName).toBe('products');
    expect(callData.fields.length).toBe(2);
    expect(callData.fields[0].isPrimaryKey).toBe(true);
    expect(callData.fields[1].fieldName).toBe('Name');
  });
  
  it('应该检查并验证必填字段', async () => {
    const wrapper = createWrapper();
    
    // 不填写任何数据
    wrapper.vm.validateStep();
    
    // 应该有验证错误
    expect(wrapper.vm.formErrors.length).toBeGreaterThan(0);
    expect(wrapper.vm.formErrors.some(e => e.includes('结构体名称'))).toBe(true);
    
    // 填写部分数据
    wrapper.vm.formData.structName = 'Product';
    wrapper.vm.validateStep();
    
    // 应该仍然有错误，但少了一个
    expect(wrapper.vm.formErrors.some(e => e.includes('结构体名称'))).toBe(false);
  });
  
  it('自动生成功能应该正常工作', async () => {
    const wrapper = createWrapper();
    
    // 设置结构体名称
    wrapper.vm.formData.structName = 'ProductCategory';
    
    // 测试自动生成表名
    await wrapper.vm.generateTableName();
    expect(wrapper.vm.formData.tableName).toBe('product_categories');
    
    // 测试自动生成API前缀
    await wrapper.vm.generateApiPrefix();
    expect(wrapper.vm.formData.apiPrefix).toBe('productCategory');
    
    // 测试自动生成描述
    await wrapper.vm.generateDescription();
    expect(wrapper.vm.formData.description).toBe('ProductCategory');
  });
  
  it('工具函数应该正确处理命名转换', () => {
    // 测试大驼峰转小驼峰
    expect(GeneratorForm.toSnakeCase('ProductCategory')).toBe('product_category');
    expect(GeneratorForm.toLowerCamel('ProductCategory')).toBe('productCategory');
    expect(GeneratorForm.toPlural('category')).toBe('categories');
    expect(GeneratorForm.toPlural('product')).toBe('products');
  });
}); 