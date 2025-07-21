#!/usr/bin/env node

/**
 * 代码生成器端到端测试工具
 * 
 * 此脚本用于测试整个代码生成流程，包括前端到后端的完整测试
 * 使用方法: node e2e-codegen-test.js
 */

const fs = require('fs');
const path = require('path');
const axios = require('axios');
const { exec, spawn } = require('child_process');
const readline = require('readline');

// 配置
const config = {
  // 前端服务配置
  frontend: {
    port: 3000,
    buildCommand: 'cd front-end && npm run build',
    startCommand: 'cd front-end && npm run dev',
    url: 'http://localhost:3000'
  },
  
  // 后端服务配置
  backend: {
    port: 8080,
    buildCommand: 'cd server && go build -o go-web',
    startCommand: 'cd server && go run main.go',
    url: 'http://localhost:8080'
  },
  
  // 测试用例
  testCase: {
    basic: {
      structName: 'E2ETest',
      tableName: 'e2e_tests',
      packageName: 'admin',
      description: 'E2E测试',
      apiPrefix: 'e2eTest',
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
          joinCondition: 'categories.name LIKE ?',
          filterCondition: 'categories.id = ?',
        }
      ]
    }
  },
  
  // 检查点
  checkpoints: [
    // 前端生成的文件
    {
      path: 'front-end/src/views/e2etest/index.vue',
      check: (content) => content.includes('<template>') && content.includes('列表')
    },
    {
      path: 'front-end/src/views/e2etest/components/TableModal.vue',
      check: (content) => content.includes('<template>') && content.includes('Modal')
    },
    {
      path: 'front-end/src/service/api/e2etest.ts',
      check: (content) => content.includes('export') && content.includes('e2eTest')
    },
    
    // 后端生成的文件
    {
      path: 'server/apps/admin/models/e2eTest.go',
      check: (content) => content.includes('type E2ETest struct') && 
                          content.includes('Category *Category') && 
                          content.includes('CategoryID uint')
    },
    {
      path: 'server/apps/admin/dto/e2etest.go',
      check: (content) => content.includes('type E2ETestResponse struct') && 
                          content.includes('CategoryFilter string')
    },
    {
      path: 'server/apps/admin/controllers/e2etest_controller.go',
      check: (content) => content.includes('func NewE2ETestController') && 
                          content.includes('Preload("Category")') && 
                          content.includes('Joins(')
    },
    {
      path: 'server/apps/admin/routes/routes.go',
      check: (content) => content.includes('e2etestGroup') && 
                          content.includes('e2etestController.GetE2ETests')
    }
  ]
};

// 工具类
const utils = {
  // 日志输出
  log(message, type = 'info') {
    const types = {
      info: '\x1b[36m%s\x1b[0m', // 青色
      success: '\x1b[32m%s\x1b[0m', // 绿色
      warning: '\x1b[33m%s\x1b[0m', // 黄色
      error: '\x1b[31m%s\x1b[0m', // 红色
      highlight: '\x1b[35m%s\x1b[0m', // 紫色
    };
    console.log(types[type] || types.info, message);
  },
  
  // 运行shell命令
  async runCommand(command) {
    return new Promise((resolve, reject) => {
      utils.log(`执行命令: ${command}`, 'info');
      
      exec(command, (error, stdout, stderr) => {
        if (error) {
          utils.log(`命令执行失败: ${error.message}`, 'error');
          utils.log(stderr, 'error');
          reject(error);
          return;
        }
        
        if (stderr) {
          utils.log(stderr, 'warning');
        }
        
        utils.log(stdout);
        resolve(stdout);
      });
    });
  },
  
  // 启动服务
  startService(command, name) {
    utils.log(`启动${name}服务: ${command}`, 'info');
    
    const childProcess = spawn(command, {
      shell: true,
      stdio: 'pipe',
    });
    
    childProcess.stdout.on('data', (data) => {
      utils.log(`[${name}] ${data}`, 'info');
    });
    
    childProcess.stderr.on('data', (data) => {
      utils.log(`[${name}] ${data}`, 'warning');
    });
    
    childProcess.on('error', (error) => {
      utils.log(`${name}服务启动失败: ${error.message}`, 'error');
    });
    
    childProcess.on('close', (code) => {
      utils.log(`${name}服务退出，退出码: ${code}`, 'info');
    });
    
    return childProcess;
  },
  
  // 等待服务就绪
  async waitForService(url, retries = 30, delay = 1000) {
    for (let i = 0; i < retries; i++) {
      try {
        utils.log(`等待服务就绪: ${url} (${i + 1}/${retries})`, 'info');
        const response = await axios.get(url);
        if (response.status === 200) {
          utils.log(`服务已就绪: ${url}`, 'success');
          return true;
        }
      } catch (error) {
        await new Promise(resolve => setTimeout(resolve, delay));
      }
    }
    
    throw new Error(`服务未就绪: ${url}`);
  },
  
  // 发送API请求
  async makeRequest(url, data = null, method = 'get') {
    try {
      utils.log(`发送请求: ${method.toUpperCase()} ${url}`, 'info');
      
      let response;
      if (method.toLowerCase() === 'get') {
        response = await axios.get(url, { params: data });
      } else if (method.toLowerCase() === 'post') {
        response = await axios.post(url, data);
      } else {
        throw new Error(`不支持的请求方法: ${method}`);
      }
      
      utils.log(`请求成功: ${response.status}`, 'success');
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
  
  // 检查生成的文件
  async checkGeneratedFiles() {
    const results = [];
    
    for (const checkpoint of config.checkpoints) {
      try {
        utils.log(`检查文件: ${checkpoint.path}`, 'info');
        
        const filePath = path.resolve(checkpoint.path);
        if (!fs.existsSync(filePath)) {
          results.push({
            path: checkpoint.path,
            success: false,
            message: '文件不存在',
          });
          utils.log(`文件不存在: ${filePath}`, 'error');
          continue;
        }
        
        const content = fs.readFileSync(filePath, 'utf8');
        const checkResult = checkpoint.check(content);
        
        results.push({
          path: checkpoint.path,
          success: checkResult,
          message: checkResult ? '检查通过' : '内容不符合预期',
        });
        
        if (checkResult) {
          utils.log(`检查通过: ${checkpoint.path}`, 'success');
        } else {
          utils.log(`内容不符合预期: ${checkpoint.path}`, 'error');
        }
      } catch (error) {
        results.push({
          path: checkpoint.path,
          success: false,
          message: `检查失败: ${error.message}`,
        });
        utils.log(`检查失败: ${checkpoint.path} - ${error.message}`, 'error');
      }
    }
    
    return results;
  },
  
  // 清理生成的文件
  async cleanupGeneratedFiles() {
    for (const checkpoint of config.checkpoints) {
      try {
        const filePath = path.resolve(checkpoint.path);
        if (fs.existsSync(filePath)) {
          fs.unlinkSync(filePath);
          utils.log(`删除文件: ${filePath}`, 'success');
        }
      } catch (error) {
        utils.log(`删除文件失败: ${checkpoint.path} - ${error.message}`, 'error');
      }
    }
  }
};

// 测试类
const tester = {
  // 前端服务进程
  frontendProcess: null,
  
  // 后端服务进程
  backendProcess: null,
  
  // 启动服务
  async startServices() {
    utils.log('启动服务...', 'highlight');
    
    // 构建前端
    try {
      await utils.runCommand(config.frontend.buildCommand);
    } catch (error) {
      utils.log('前端构建失败', 'error');
      return false;
    }
    
    // 构建后端
    try {
      await utils.runCommand(config.backend.buildCommand);
    } catch (error) {
      utils.log('后端构建失败', 'error');
      return false;
    }
    
    // 启动服务
    this.backendProcess = utils.startService(config.backend.startCommand, '后端');
    this.frontendProcess = utils.startService(config.frontend.startCommand, '前端');
    
    // 等待服务就绪
    try {
      await utils.waitForService(`${config.backend.url}/api/health`);
      await utils.waitForService(config.frontend.url);
    } catch (error) {
      utils.log('服务启动失败', 'error');
      this.stopServices();
      return false;
    }
    
    utils.log('服务启动成功', 'success');
    return true;
  },
  
  // 停止服务
  stopServices() {
    utils.log('停止服务...', 'highlight');
    
    if (this.frontendProcess) {
      this.frontendProcess.kill();
      this.frontendProcess = null;
    }
    
    if (this.backendProcess) {
      this.backendProcess.kill();
      this.backendProcess = null;
    }
    
    utils.log('服务已停止', 'success');
  },
  
  // 登录
  async login() {
    utils.log('登录系统...', 'highlight');
    
    try {
      const response = await utils.makeRequest(
        `${config.backend.url}/admin/login`,
        { username: 'admin', password: '123456' },
        'post'
      );
      
      if (response.code !== 0 || !response.data.token) {
        utils.log('登录失败', 'error');
        return null;
      }
      
      utils.log('登录成功', 'success');
      return response.data.token;
    } catch (error) {
      utils.log('登录失败', 'error');
      return null;
    }
  },
  
  // 生成代码
  async generateCode(token) {
    utils.log('生成代码...', 'highlight');
    
    try {
      const response = await utils.makeRequest(
        `${config.backend.url}/admin/codegen/generate`,
        config.testCase.basic,
        'post',
        {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        }
      );
      
      if (response.code !== 0) {
        utils.log('代码生成失败', 'error');
        return false;
      }
      
      utils.log('代码生成成功', 'success');
      return true;
    } catch (error) {
      utils.log('代码生成失败', 'error');
      return false;
    }
  },
  
  // 运行测试
  async runTest() {
    utils.log('开始端到端测试...', 'highlight');
    
    // 清理之前的文件
    await utils.cleanupGeneratedFiles();
    
    // 启动服务
    if (!await this.startServices()) {
      return false;
    }
    
    // 登录
    const token = await this.login();
    if (!token) {
      this.stopServices();
      return false;
    }
    
    // 生成代码
    if (!await this.generateCode(token)) {
      this.stopServices();
      return false;
    }
    
    // 检查生成的文件
    const checkResults = await utils.checkGeneratedFiles();
    const success = checkResults.every(result => result.success);
    
    // 输出结果
    utils.log('\n===== 测试结果 =====', 'highlight');
    
    if (success) {
      utils.log('测试通过!', 'success');
    } else {
      utils.log('测试失败!', 'error');
      
      utils.log('\n===== 失败项 =====', 'error');
      checkResults
        .filter(result => !result.success)
        .forEach(result => {
          utils.log(`${result.path}: ${result.message}`, 'error');
        });
    }
    
    // 停止服务
    this.stopServices();
    
    return success;
  }
};

// 命令行交互
const cli = {
  // 创建命令行接口
  createInterface() {
    return readline.createInterface({
      input: process.stdin,
      output: process.stdout
    });
  },
  
  // 显示菜单
  async showMenu() {
    const rl = this.createInterface();
    
    utils.log('代码生成器端到端测试工具', 'highlight');
    utils.log('1. 运行完整端到端测试', 'info');
    utils.log('2. 仅检查生成的文件', 'info');
    utils.log('3. 仅清理生成的文件', 'info');
    utils.log('0. 退出', 'info');
    
    const answer = await new Promise(resolve => {
      rl.question('请选择操作: ', resolve);
    });
    
    rl.close();
    
    switch (answer) {
      case '1':
        await tester.runTest();
        break;
      case '2':
        await utils.checkGeneratedFiles();
        break;
      case '3':
        await utils.cleanupGeneratedFiles();
        break;
      case '0':
        utils.log('再见!', 'highlight');
        return;
      default:
        utils.log('无效的选择', 'error');
        await this.showMenu();
        return;
    }
    
    // 返回菜单
    await new Promise(resolve => {
      const rl = this.createInterface();
      rl.question('\n按回车键继续...', () => {
        rl.close();
        resolve();
      });
    });
    
    await this.showMenu();
  }
};

// 主函数
async function main() {
  const args = process.argv.slice(2);
  
  if (args.includes('--help') || args.includes('-h')) {
    utils.log('代码生成器端到端测试工具', 'highlight');
    utils.log('用法: node e2e-codegen-test.js [选项]', 'info');
    utils.log('选项:', 'info');
    utils.log('  --run         直接运行测试', 'info');
    utils.log('  --check       仅检查文件', 'info');
    utils.log('  --clean       仅清理文件', 'info');
    utils.log('  --help, -h    显示此帮助', 'info');
    return;
  }
  
  if (args.includes('--run')) {
    await tester.runTest();
    return;
  }
  
  if (args.includes('--check')) {
    await utils.checkGeneratedFiles();
    return;
  }
  
  if (args.includes('--clean')) {
    await utils.cleanupGeneratedFiles();
    return;
  }
  
  // 显示交互菜单
  await cli.showMenu();
}

// 启动
main().catch(error => {
  utils.log(`致命错误: ${error.message}`, 'error');
  console.error(error);
  process.exit(1);
}); 