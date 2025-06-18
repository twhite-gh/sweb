# SWeb 测试套件

## 概述

这个目录包含了SWeb项目的所有测试文件和脚本。测试套件已经重新组织，提供了统一的测试入口和交互式菜单。

## 文件说明

### 测试脚本
- `test_all.bat` - Windows完整测试套件（交互式菜单）
- `test_all.sh` - Linux/macOS完整测试套件（交互式菜单）

### 测试程序
- `test_all.go` - 综合功能测试，检查服务器状态和基础功能
- `test_socks5.go` - SOCKS5代理专项测试
- `test_proxy.go` - HTTP/HTTPS代理专项测试
- `debug_socks5.go` - SOCKS5连接调试工具

## 快速开始

### Windows用户

1. 打开命令提示符或PowerShell
2. 进入测试目录：
   ```cmd
   cd test
   ```
3. 运行测试套件：
   ```cmd
   test_all.bat
   ```
4. 根据菜单选择要运行的测试

### Linux/macOS用户

1. 打开终端
2. 进入测试目录：
   ```bash
   cd test
   ```
3. 给脚本执行权限：
   ```bash
   chmod +x test_all.sh
   ```
4. 运行测试套件：
   ```bash
   ./test_all.sh
   ```
5. 根据菜单选择要运行的测试

## 测试菜单选项

1. **Run all tests** - 运行所有测试（推荐用于完整验证）
2. **Basic functionality test** - 基础功能测试（检查服务器状态）
3. **SOCKS5 proxy test** - SOCKS5代理专项测试
4. **HTTP/HTTPS proxy test** - HTTP/HTTPS代理专项测试
5. **Debug connection test** - 调试连接测试
6. **Quick SOCKS5 test** - SOCKS5快速测试（包含编译和运行）
0. **Exit** - 退出测试套件

## 前置条件

### 启动服务器

在运行测试之前，请确保SWeb服务器已启动并启用了相应功能：

```bash
# 启用所有功能（推荐）
./sweb -upload -files -webdav -https -socks5 -proxy

# 或者根据需要启用特定功能
./sweb -socks5                    # 仅SOCKS5代理
./sweb -proxy                     # 仅HTTP代理
./sweb -upload -files             # 仅文件功能
```

### 网络要求

- 确保端口8080（HTTP）、8443（HTTPS）、1080（SOCKS5）、10808（HTTP代理）未被占用
- 测试需要访问外部网站（如baidu.com、httpbin.org）
- 确保防火墙允许相关端口的连接

## 单独运行测试

如果您不想使用交互式菜单，也可以直接运行单个测试：

```bash
# 综合功能测试
go run test_all.go

# SOCKS5代理测试
go run test_socks5.go

# HTTP代理测试
go run test_proxy.go

# 调试工具
go run debug_socks5.go
```

## 测试结果说明

### 成功标志
- ✅ 或 √ 表示测试通过
- 绿色文字或正常输出表示功能正常

### 失败标志
- ❌ 或 X 表示测试失败
- 红色文字或错误信息表示需要检查

### 常见错误

1. **连接被拒绝**
   - 确保服务器已启动
   - 检查端口是否被占用
   - 确认防火墙设置

2. **编译失败**
   - 检查Go环境是否正确安装
   - 确保在正确的目录中运行
   - 检查源文件是否完整

3. **网络超时**
   - 检查网络连接
   - 尝试访问其他测试网站
   - 某些网站可能被防火墙阻止

## 浏览器配置验证

测试通过后，您可以配置浏览器验证代理功能：

### SOCKS5代理配置
- 类型：SOCKS5
- 地址：127.0.0.1
- 端口：1080

### HTTP代理配置
- HTTP代理：127.0.0.1:10808
- HTTPS代理：127.0.0.1:10808（与HTTP相同）

## 故障排除

如果遇到问题，请：

1. 检查服务器是否正常启动
2. 查看服务器日志输出
3. 使用debug_socks5.go进行详细诊断
4. 检查网络连接和防火墙设置
5. 参考主目录下的相关文档

## 贡献

如果您发现测试中的问题或想要添加新的测试，请：

1. 在test目录下创建新的测试文件
2. 更新test_all.bat和test_all.sh脚本
3. 更新本README文件
4. 提交Pull Request

## 更多信息

- 项目主页：../README.md
- 测试详细说明：../测试说明.md
- HTTPS代理配置：../HTTPS代理配置指南.md
- 浏览器配置：../浏览器配置说明.md
