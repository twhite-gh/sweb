# Windows WebDAV 挂载指南

## 问题描述
在Windows系统中使用WebDAV认证时，可能会遇到"输入的文件夹似乎无效"的错误。这通常是由于Windows WebDAV客户端对服务器响应头部的严格要求导致的。

## 解决方案

### 1. 服务器端修复（已实现）
我们已经对WebDAV服务器进行了以下优化：

- ✅ 添加了`charset="UTF-8"`到WWW-Authenticate头部
- ✅ 设置了`Server: Microsoft-IIS/10.0`头部以提高Windows兼容性
- ✅ 为OPTIONS请求添加了完整的WebDAV头部支持
- ✅ 优化了认证流程以符合Windows WebDAV客户端期望

### 2. Windows客户端配置

#### 方法1: 使用文件资源管理器
1. 打开文件资源管理器
2. 右键点击"此电脑"
3. 选择"映射网络驱动器"
4. 点击"连接到可用于存储文档和图片的网站"
5. 点击"下一步"
6. 选择"选择自定义网络位置"
7. 输入WebDAV地址：`http://localhost:8081/webdav`
8. 输入用户名：`webdav`
9. 输入密码：`webdav`

#### 方法2: 使用命令行
```cmd
# 映射网络驱动器
net use Z: http://localhost:8081/webdav /user:webdav webdav

# 取消映射
net use Z: /delete
```

#### 方法3: 使用PowerShell
```powershell
# 创建网络驱动器
New-PSDrive -Name "WebDAV" -PSProvider FileSystem -Root "http://localhost:8081/webdav" -Credential (Get-Credential)

# 删除网络驱动器
Remove-PSDrive -Name "WebDAV"
```

### 3. Windows注册表优化（可选）

如果仍然遇到问题，可以修改以下注册表设置：

```reg
Windows Registry Editor Version 5.00

[HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Services\WebClient\Parameters]
"BasicAuthLevel"=dword:00000002
"FileSizeLimitInBytes"=dword:ffffffff
"FileAttributesLimitInBytes"=dword:00100000
```

### 4. 故障排除

#### 常见错误及解决方案

**错误1: "输入的文件夹似乎无效"**
- 确保服务器正在运行：`.\sweb.exe -webdav -webdav-auth -port 8082`
- 检查URL格式：`http://localhost:8082/webdav`
- 确认用户名密码：`webdav/webdav`

**错误2: "错误代码：0x80070043 找不到网络名"**
- **启用WebClient服务**：运行 `fix_webdav_service.bat`
- **手动启动服务**：以管理员身份运行 `sc start WebClient`
- **设置自动启动**：`sc config WebClient start= auto`
- **导入注册表设置**：双击 `webdav_registry_fix.reg`
- **重启服务**：`sc stop WebClient && sc start WebClient`

**错误3: "网络路径未找到"**
- 检查防火墙设置
- 确认端口8082未被占用
- 尝试使用IP地址：`http://127.0.0.1:8082/webdav`
- 检查WebClient服务是否运行：`sc query WebClient`

**错误4: 认证失败**
- 确认WebDAV认证已启用：`-webdav-auth`
- 检查用户名密码是否正确
- 查看服务器日志输出

**错误5: WebClient服务无法启动**
- 以管理员身份运行命令提示符
- 检查依赖服务：`sc query MRxDAV`
- 重启计算机后重试
- 检查Windows功能是否启用WebDAV

#### 自动修复工具

我们提供了几个自动修复工具来解决常见问题：

**1. WebClient服务修复工具**
```cmd
fix_webdav_service.bat
```
此工具会：
- 检查并启动WebClient服务
- 设置服务为自动启动
- 添加必要的注册表项
- 重启服务以应用更改

**2. 注册表修复文件**
```cmd
webdav_registry_fix.reg
```
双击导入此文件以优化WebDAV客户端设置

**3. 增强测试工具**
```cmd
test_webdav_mount_enhanced.bat
```
此工具会尝试多种挂载方法并提供详细的错误诊断

#### 测试连接
使用我们提供的测试工具验证连接：
```bash
cd test
go run test_webdav_windows.go
```

### 5. 高级配置

#### 自定义端口
```bash
# 使用自定义端口
go run . -webdav -webdav-auth -port 9090
# WebDAV地址变为：http://localhost:9090/webdav
```

#### 自定义认证信息
```bash
# 使用自定义用户名密码
go run . -webdav -webdav-auth -webdav-username myuser -webdav-password mypass
```

#### 只读模式
```bash
# 启用只读WebDAV
go run . -webdav -webdav-auth -webdav-readonly
```

### 6. 安全建议

1. **生产环境使用HTTPS**
   ```bash
   go run . -webdav -webdav-auth -https -https-port 8443
   # 使用：https://localhost:8443/webdav
   ```

2. **修改默认密码**
   ```bash
   go run . -webdav -webdav-auth -webdav-username admin -webdav-password your_secure_password
   ```

3. **限制访问目录**
   ```bash
   # 指定特定目录（如果支持）
   go run . -webdav -webdav-auth -webdav-dir /path/to/safe/directory
   ```

## 验证成功挂载

挂载成功后，您应该能够：
- 在文件资源管理器中看到WebDAV驱动器
- 浏览、创建、编辑和删除文件
- 拖拽文件进行上传和下载

## 技术说明

我们的WebDAV实现包含以下Windows兼容性优化：
- RFC 4918 WebDAV协议完整支持
- Windows WebDAV客户端特定的HTTP头部
- 优化的OPTIONS和PROPFIND响应
- Basic认证的正确实现
- 错误处理和日志记录

如果您仍然遇到问题，请检查服务器日志输出或联系技术支持。
