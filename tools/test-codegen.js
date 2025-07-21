#!/usr/bin/env node

/**
 * 代码生成器测试工具
 * 
 * 此脚本用于手动测试代码生成器功能并验证各种场景下的行为
 */

const fs = require('fs');
const path = require('path');
const axios = require('axios');
const { exec } = require('child_process');

// 配置
const config = {
  // 后端API地址
  apiBaseUrl: 'http://localhost:8080/admin',
  
  // 测试用例 - 包含基本字段
  basicFieldsTest: {
    structName: 'TestBasic',
    tableName: 'test_basics',
    packageName: 'admin',
    description: '基础测试',
    apiPrefix: 'testBasic',
    appName: 'admin',
    hasList: true,
    hasCreate: true,
    hasUpdate: true,
    hasDelete: true,
    hasDetail: true,
    hasPagination: true,
    fields: [
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
      },
      {
        fieldName: 'Status',
        fieldType: 'int',
        columnName: 'status',
        fieldDesc: '状态',
        required: false,
        isPrimaryKey: false,
        isSearchable: false,
        isFilterable: true,
        isSortable: false,
      }
    ]
  },
  
  // 测试用例 - 包含关系字段
  relationFieldsTest: {
    structName: 'TestRelation',
    tableName: 'test_relations',
    packageName: 'admin',
    description: '关系测试',
    apiPrefix: 'testRelation',
    appName: 'admin',
    hasList: true,
    hasCreate: true,
    hasUpdate: true,
    hasDelete: true,
    hasDetail: true,
    hasPagination: true,
    fields: [
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
        fieldName: 'Title',
        fieldType: 'string',
        columnName: 'title',
        fieldDesc: '标题',
        required: true,
        isPrimaryKey: false,
        isSearchable: true,
        isFilterable: true,
        isSortable: false,
      },
      {
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
        foreignKey: 'CategoryID',
        references: 'ID',
        preload: true,
        joinable: true,
      },
      {
        fieldName: 'Tags',
        fieldType: 'string', // 会被覆盖
        columnName: '',
        fieldDesc: '标签',
        required: false,
        isPrimaryKey: false,
        isSearchable: false,
        isFilterable: false,
        isSortable: false,
        isRelation: true,
        relationType: 'many_to_many',
        relatedModel: 'Tag',
        joinTable: 'test_relation_tags',
        preload: true,
        joinable: false,
      }
    ]
  },
  
  // 生成的文件路径
  generatedPaths: {
    model: 'server/apps/admin/models/',
    dto: 'server/apps/admin/dto/',
    controller: 'server/apps/admin/controllers/',
    routes: 'server/apps/admin/routes/routes.go',
    frontend: 'front-end/src/views/'
  }
};

// 工具函数
const utils = {
  log(message, type = 'info') {
    const types = {
      info: '\x1b[36m%s\x1b[0m', // 青色
      success: '\x1b[32m%s\x1b[0m', // 绿色
      warning: '\x1b[33m%s\x1b[0m', // 黄色
      error: '\x1b[31m%s\x1b[0m', // 红色
    };
    console.log(types[type], message);
  },
  
  async makeRequest(endpoint, data = null, method = 'get') {
    try {
      const url = `${config.apiBaseUrl}${endpoint}`;
      let response;
      
      if (method.toLowerCase() === 'get') {
        response = await axios.get(url, { params: data });
      } else if (method.toLowerCase() === 'post') {
        response = await axios.post(url, data);
      } else {
        throw new Error(`不支持的请求方法: ${method}`);
      }
      
      return response.data;
    } catch (error) {
      utils.log(`请求失败: ${error.message}`, 'error');
      if (error.response) {
        utils.log(`状态码: ${error.response.status}`, 'error');
        utils.log(`响应: ${JSON.stringify(error.response.data)}`, 'error');
      }
      throw error;
    }
  },
  
  checkFile(filepath) {
    try {
      if (fs.existsSync(filepath)) {
        const content = fs.readFileSync(filepath, 'utf8');
        return { exists: true, content };
      }
      return { exists: false, content: null };
    } catch (error) {
      utils.log(`检查文件失败: ${error.message}`, 'error');
      return { exists: false, error };
    }
  },
  
  validateGeneratedFiles(testCase, checkOnly = false) {
    const workDir = process.cwd();
    const results = [];
    
    // 检查模型文件
    const modelFileName = testCase.structName.charAt(0).toLowerCase() + testCase.structName.slice(1) + '.go';
    const modelPath = path.join(workDir, config.generatedPaths.model, modelFileName);
    const modelResult = utils.checkFile(modelPath);
    
    if (modelResult.exists) {
      utils.log(`模型文件已生成: ${modelPath}`, 'success');
      
      // 检查关系字段
      testCase.fields.forEach(field => {
        if (field.isRelation) {
          if (field.relationType === 'belongs_to') {
            // 检查关联模型字段
            const relationFieldPattern = `${field.fieldName} \\*${field.relatedModel}`;
            if (!modelResult.content.match(new RegExp(relationFieldPattern))) {
              utils.log(`模型文件中缺少关系字段: ${field.fieldName} *${field.relatedModel}`, 'error');
              results.push({
                type: 'error',
                file: 'model',
                message: `缺少关系字段: ${field.fieldName} *${field.relatedModel}`
              });
            }
            
            // 检查外键字段
            const foreignKeyPattern = `${field.foreignKey || field.relatedModel + 'ID'} uint`;
            if (!modelResult.content.match(new RegExp(foreignKeyPattern))) {
              utils.log(`模型文件中缺少外键字段: ${field.foreignKey || field.relatedModel + 'ID'} uint`, 'error');
              results.push({
                type: 'error',
                file: 'model',
                message: `缺少外键字段: ${field.foreignKey || field.relatedModel + 'ID'} uint`
              });
            }
          }
          
          // 检查预加载
          if (field.preload) {
            const preloadPattern = `Preload\\("${field.fieldName}"\\)`;
            if (!modelResult.content.match(new RegExp(preloadPattern))) {
              utils.log(`模型文件中缺少预加载代码: Preload("${field.fieldName}")`, 'error');
              results.push({
                type: 'error',
                file: 'model',
                message: `缺少预加载代码: Preload("${field.fieldName}")`
              });
            }
          }
        }
      });
    } else {
      utils.log(`模型文件未找到: ${modelPath}`, 'error');
      results.push({
        type: 'error',
        file: 'model',
        message: '文件不存在'
      });
    }
    
    // 检查DTO文件
    const dtoFileName = testCase.structName.charAt(0).toLowerCase() + testCase.structName.slice(1) + '.go';
    const dtoPath = path.join(workDir, config.generatedPaths.dto, dtoFileName);
    const dtoResult = utils.checkFile(dtoPath);
    
    if (dtoResult.exists) {
      utils.log(`DTO文件已生成: ${dtoPath}`, 'success');
      
      // 检查关联过滤字段
      testCase.fields.forEach(field => {
        if (field.isRelation && field.isFilterable) {
          const filterFieldPattern = `${field.fieldName}Filter string`;
          if (!dtoResult.content.match(new RegExp(filterFieldPattern))) {
            utils.log(`DTO文件中缺少关联过滤字段: ${field.fieldName}Filter string`, 'error');
            results.push({
              type: 'error',
              file: 'dto',
              message: `缺少关联过滤字段: ${field.fieldName}Filter string`
            });
          }
        }
      });
    } else {
      utils.log(`DTO文件未找到: ${dtoPath}`, 'error');
      results.push({
        type: 'error',
        file: 'dto',
        message: '文件不存在'
      });
    }
    
    // 检查控制器文件
    const controllerFileName = testCase.structName.charAt(0).toLowerCase() + testCase.structName.slice(1) + '_controller.go';
    const controllerPath = path.join(workDir, config.generatedPaths.controller, controllerFileName);
    const controllerResult = utils.checkFile(controllerPath);
    
    if (controllerResult.exists) {
      utils.log(`控制器文件已生成: ${controllerPath}`, 'success');
      
      // 检查预加载代码
      const hasPreloadFields = testCase.fields.some(f => f.isRelation && f.preload);
      if (hasPreloadFields) {
        const preloadPattern = `Preload\\("`;
        if (!controllerResult.content.match(new RegExp(preloadPattern))) {
          utils.log(`控制器文件中缺少预加载代码`, 'error');
          results.push({
            type: 'error',
            file: 'controller',
            message: '缺少预加载代码'
          });
        }
      }
      
      // 检查JOIN查询代码
      const hasJoinableFields = testCase.fields.some(f => f.isRelation && f.joinable);
      if (hasJoinableFields) {
        const joinPattern = `Joins\\(`;
        if (!controllerResult.content.match(new RegExp(joinPattern))) {
          utils.log(`控制器文件中缺少JOIN查询代码`, 'error');
          results.push({
            type: 'error',
            file: 'controller',
            message: '缺少JOIN查询代码'
          });
        }
      }
    } else {
      utils.log(`控制器文件未找到: ${controllerPath}`, 'error');
      results.push({
        type: 'error',
        file: 'controller',
        message: '文件不存在'
      });
    }
    
    // 检查路由文件
    const routesPath = path.join(workDir, config.generatedPaths.routes);
    const routesResult = utils.checkFile(routesPath);
    
    if (routesResult.exists) {
      utils.log(`路由文件已生成: ${routesPath}`, 'success');
      
      // 检查控制器变量
      const controllerVarPattern = `${testCase.structName.charAt(0).toLowerCase() + testCase.structName.slice(1)}Controller`;
      if (!routesResult.content.includes(controllerVarPattern)) {
        utils.log(`路由文件中缺少控制器变量: ${controllerVarPattern}`, 'error');
        results.push({
          type: 'error',
          file: 'routes',
          message: `缺少控制器变量: ${controllerVarPattern}`
        });
      }
      
      // 检查路由注册
      const routeCommentPattern = `// ${testCase.description}路由`;
      if (!routesResult.content.includes(routeCommentPattern)) {
        utils.log(`路由文件中缺少路由注册: ${routeCommentPattern}`, 'error');
        results.push({
          type: 'error',
          file: 'routes',
          message: `缺少路由注册: ${routeCommentPattern}`
        });
      }
      
      // 检查路由是否有重复
      const routePattern = `GET\\("/${testCase.apiPrefix.toLowerCase()}s"`;
      const routeMatches = routesResult.content.match(new RegExp(routePattern, 'g'));
      if (routeMatches && routeMatches.length > 1) {
        utils.log(`路由文件中存在重复路由: ${routePattern}`, 'error');
        results.push({
          type: 'error',
          file: 'routes',
          message: `存在重复路由: ${routePattern}`
        });
      }
    } else {
      utils.log(`路由文件未找到: ${routesPath}`, 'error');
      results.push({
        type: 'error',
        file: 'routes',
        message: '文件不存在'
      });
    }
    
    return results;
  }
};

// 测试函数
const tests = {
  async runBasicFieldsTest() {
    utils.log('开始测试基本字段...', 'info');
    try {
      // 清理测试环境
      // 这里可能需要删除之前生成的文件
      
      // 发送代码生成请求
      const response = await utils.makeRequest('/codegen/generate', config.basicFieldsTest, 'post');
      utils.log('代码生成请求成功', 'success');
      utils.log(`响应: ${JSON.stringify(response)}`, 'info');
      
      // 验证生成的文件
      const results = utils.validateGeneratedFiles(config.basicFieldsTest);
      
      return { success: results.filter(r => r.type === 'error').length === 0, results };
    } catch (error) {
      utils.log(`测试失败: ${error.message}`, 'error');
      return { success: false, error };
    }
  },
  
  async runRelationFieldsTest() {
    utils.log('开始测试关系字段...', 'info');
    try {
      // 清理测试环境
      // 这里可能需要删除之前生成的文件
      
      // 发送代码生成请求
      const response = await utils.makeRequest('/codegen/generate', config.relationFieldsTest, 'post');
      utils.log('代码生成请求成功', 'success');
      utils.log(`响应: ${JSON.stringify(response)}`, 'info');
      
      // 验证生成的文件
      const results = utils.validateGeneratedFiles(config.relationFieldsTest);
      
      return { success: results.filter(r => r.type === 'error').length === 0, results };
    } catch (error) {
      utils.log(`测试失败: ${error.message}`, 'error');
      return { success: false, error };
    }
  },
  
  analyzeServerSupport() {
    utils.log('分析服务器API支持...', 'info');
    
    // 检查关联查询参数和JOIN支持
    const hasJoinFieldsInApi = config.relationFieldsTest.fields.some(f => f.joinable);
    utils.log(`API配置中${hasJoinFieldsInApi ? '已包含' : '未包含'}JOIN查询参数`, hasJoinFieldsInApi ? 'success' : 'warning');
    
    // 检查DTO格式
    const hasForeignKeys = config.relationFieldsTest.fields.some(f => f.isRelation && f.foreignKey);
    utils.log(`API配置中${hasForeignKeys ? '已包含' : '未包含'}外键定义`, hasForeignKeys ? 'success' : 'warning');
    
    // 检查预加载参数
    const hasPreloadFlags = config.relationFieldsTest.fields.some(f => f.isRelation && f.preload !== undefined);
    utils.log(`API配置中${hasPreloadFlags ? '已包含' : '未包含'}预加载标志`, hasPreloadFlags ? 'success' : 'warning');
    
    return {
      joinSupport: hasJoinFieldsInApi,
      foreignKeySupport: hasForeignKeys,
      preloadSupport: hasPreloadFlags
    };
  }
};

// 运行所有测试
async function runAllTests() {
  utils.log('开始运行所有测试...', 'info');
  
  // 分析API支持
  const apiSupport = tests.analyzeServerSupport();
  utils.log(`API支持分析: ${JSON.stringify(apiSupport)}`, 'info');
  
  // 运行基本字段测试
  const basicResult = await tests.runBasicFieldsTest();
  utils.log(`基本字段测试结果: ${basicResult.success ? '通过' : '失败'}`, basicResult.success ? 'success' : 'error');
  
  // 运行关系字段测试
  const relationResult = await tests.runRelationFieldsTest();
  utils.log(`关系字段测试结果: ${relationResult.success ? '通过' : '失败'}`, relationResult.success ? 'success' : 'error');
  
  // 输出测试摘要
  utils.log('\n===== 测试摘要 =====', 'info');
  utils.log(`基本字段测试: ${basicResult.success ? '通过' : '失败'}`, basicResult.success ? 'success' : 'error');
  utils.log(`关系字段测试: ${relationResult.success ? '通过' : '失败'}`, relationResult.success ? 'success' : 'error');
  
  // 输出建议
  utils.log('\n===== 改进建议 =====', 'info');
  
  if (!basicResult.success || !relationResult.success) {
    utils.log('1. 检查生成器的路径设置，确保生成文件到正确位置', 'warning');
    utils.log('2. 检查DTO和模型生成器中是否正确处理关系字段', 'warning');
    utils.log('3. 验证预加载(Preload)选项是否正确传递到后端', 'warning');
    utils.log('4. 验证JOIN查询配置是否正确传递到后端', 'warning');
    utils.log('5. 检查路由生成器是否存在重复注册问题', 'warning');
    
    // 输出详细错误
    utils.log('\n===== 详细错误 =====', 'error');
    
    if (!basicResult.success && basicResult.results) {
      utils.log('基本字段测试错误:', 'error');
      basicResult.results.forEach((result, index) => {
        if (result.type === 'error') {
          utils.log(`${index + 1}. ${result.file}: ${result.message}`, 'error');
        }
      });
    }
    
    if (!relationResult.success && relationResult.results) {
      utils.log('关系字段测试错误:', 'error');
      relationResult.results.forEach((result, index) => {
        if (result.type === 'error') {
          utils.log(`${index + 1}. ${result.file}: ${result.message}`, 'error');
        }
      });
    }
  } else {
    utils.log('所有测试通过!', 'success');
  }
}

// 入口函数
function main() {
  const args = process.argv.slice(2);
  
  if (args.includes('--help') || args.includes('-h')) {
    utils.log('代码生成器测试工具', 'info');
    utils.log('用法:', 'info');
    utils.log('  node test-codegen.js [选项]', 'info');
    utils.log('选项:', 'info');
    utils.log('  --basic      只运行基本字段测试', 'info');
    utils.log('  --relation   只运行关系字段测试', 'info');
    utils.log('  --analyze    只分析API支持', 'info');
    utils.log('  --help, -h   显示此帮助信息', 'info');
    return;
  }
  
  if (args.includes('--basic')) {
    tests.runBasicFieldsTest();
    return;
  }
  
  if (args.includes('--relation')) {
    tests.runRelationFieldsTest();
    return;
  }
  
  if (args.includes('--analyze')) {
    tests.analyzeServerSupport();
    return;
  }
  
  // 默认运行所有测试
  runAllTests();
}

// 执行入口函数
main(); 