# HTTPS代理配置指南

## 概述

SWeb的HTTP代理服务器同时支持HTTP和HTTPS协议。HTTPS代理通过HTTP CONNECT方法建立SSL/TLS隧道，实现安全的HTTPS流量转发。

## 启动HTTPS代理服务

```bash
# 启动HTTP代理服务（同时支持HTTPS）
./sweb.exe -proxy

# 指定自定义端口
./sweb.exe -proxy -proxy-port 10808

# 与其他功能一起使用
./sweb.exe -proxy -socks5 -upload -files
```

## 浏览器配置

### Chrome/Edge配置

1. **打开代理设置**
   - 设置 → 高级 → 系统 → 打开代理设置
   - 或直接访问：`chrome://settings/system`

2. **手动代理配置**
   - 选择"手动代理配置"
   - HTTP代理：`127.0.0.1:10808`
   - HTTPS代理：`127.0.0.1:10808` (与HTTP相同)
   - 勾选"对所有协议使用此代理服务器"

3. **验证配置**
   - 访问：`https://www.baidu.com`
   - 访问：`https://httpbin.org/ip`

### Firefox配置

1. **网络设置**
   - 设置 → 网络设置 → 手动代理配置

2. **代理配置**
   - HTTP代理：`127.0.0.1` 端口：`10808`
   - SSL代理：`127.0.0.1` 端口：`10808`
   - 勾选"为所有协议使用此代理"

3. **DNS设置**
   - 勾选"通过代理服务器进行DNS查询"

## 命令行工具测试

### 使用curl测试HTTPS代理

```bash
# 通过HTTP代理访问HTTPS网站
curl --proxy http://127.0.0.1:10808 https://www.baidu.com

# 显示详细信息
curl -v --proxy http://127.0.0.1:10808 https://httpbin.org/ip

# 忽略SSL证书验证（测试用）
curl -k --proxy http://127.0.0.1:10808 https://www.google.com
```

### 使用wget测试

```bash
# 设置代理环境变量
export http_proxy=http://127.0.0.1:10808
export https_proxy=http://127.0.0.1:10808

# 测试HTTPS访问
wget https://www.baidu.com
```

## 工作原理

### HTTP CONNECT方法

1. **建立隧道**
   ```
   客户端 → 代理服务器: CONNECT www.example.com:443 HTTP/1.1
   代理服务器 → 客户端: HTTP/1.1 200 Connection established
   ```

2. **SSL/TLS握手**
   ```
   客户端 ←→ 目标服务器 (通过代理隧道进行SSL握手)
   ```

3. **数据传输**
   ```
   客户端 ←→ 代理服务器 ←→ 目标服务器 (加密数据透明转发)
   ```

### 支持的功能

- ✅ HTTP GET/POST请求
- ✅ HTTPS CONNECT隧道
- ✅ SSL/TLS透明代理
- ✅ 并发连接处理
- ✅ 错误处理和重试
- ✅ 连接统计和监控

## 测试验证

### 自动化测试

```bash
# 运行HTTPS代理测试
test_https_proxy.bat

# 运行完整代理测试套件
test_proxy_all.bat
```

### 手动验证步骤

1. **启动服务器**
   ```bash
   ./sweb.exe -proxy
   ```

2. **检查端口监听**
   ```bash
   netstat -an | grep :10808
   ```

3. **测试HTTP代理**
   ```bash
   curl --proxy http://127.0.0.1:10808 http://httpbin.org/ip
   ```

4. **测试HTTPS代理**
   ```bash
   curl --proxy http://127.0.0.1:10808 https://httpbin.org/ip
   ```

5. **浏览器测试**
   - 配置浏览器代理
   - 访问HTTPS网站
   - 检查服务器日志

## 性能特性

### 并发处理
- 基于Go协程的高并发处理
- 每个连接独立的goroutine
- 高效的内存使用

### 连接管理
- 自动连接超时处理
- 优雅的连接关闭
- 错误恢复机制

### 监控统计
- 实时连接数统计
- HTTP/HTTPS请求分类统计
- 数据传输量统计
- 成功/失败连接统计

## 安全考虑

### 默认安全策略
- 代理功能默认禁用
- 无认证机制（适用于可信网络）
- 仅监听本地地址

### 使用建议
- 仅在可信网络环境中使用
- 不建议暴露到公网
- 定期检查连接日志
- 监控异常流量

### 限制说明
- 不支持用户认证
- 不支持访问控制列表
- 不支持流量加密（HTTPS内容本身已加密）

## 故障排除

### 常见问题

1. **HTTPS网站无法访问**
   - 检查代理配置是否正确
   - 确认目标网站可直接访问
   - 查看服务器错误日志

2. **SSL证书错误**
   - 某些网站的SSL证书可能有问题
   - 可以使用curl的-k参数跳过验证测试

3. **连接超时**
   - 检查网络连接
   - 增加超时时间设置
   - 尝试其他测试网站

4. **代理拒绝连接**
   - 确认代理服务器正在运行
   - 检查端口配置
   - 验证防火墙设置
