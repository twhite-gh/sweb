# WebDAV功能使用说明

## 概述

本Web文件服务器现已支持WebDAV (Web Distributed Authoring and Versioning) 协议，允许用户通过标准的文件管理客户端远程管理服务器上的文件。

## 启用WebDAV服务

### 基本启用
```bash
# 启用WebDAV服务（默认目录为当前目录）
sweb.exe -webdav

# 启用WebDAV服务并指定目录
sweb.exe -webdav -webdav-dir /path/to/directory

# 启用只读WebDAV服务
sweb.exe -webdav -webdav-readonly

# 同时启用文件上传和WebDAV服务
sweb.exe -upload -webdav
```

### 命令行参数说明
- `-webdav` 或 `--enable-webdav`: 启用WebDAV服务
- `-webdav-dir <目录>`: 指定WebDAV服务的根目录（默认为当前目录）
- `-webdav-readonly`: 设置WebDAV为只读模式

## WebDAV地址

当WebDAV服务启用后，可以通过以下地址访问：
```
http://localhost:8080/webdav
```

## 客户端连接方法

### Windows系统
1. 打开文件资源管理器
2. 右键点击"此电脑"
3. 选择"映射网络驱动器"
4. 在文件夹路径中输入：`http://localhost:8080/webdav`
5. 点击"完成"

### macOS系统
1. 打开Finder
2. 按下 `Cmd + K` 或选择"前往" → "连接服务器"
3. 在服务器地址中输入：`http://localhost:8080/webdav`
4. 点击"连接"

### Linux系统
使用davfs2挂载：
```bash
# 安装davfs2
sudo apt-get install davfs2  # Ubuntu/Debian
sudo yum install davfs2      # CentOS/RHEL

# 挂载WebDAV
sudo mount -t davfs http://localhost:8080/webdav /mnt/webdav
```

### 移动设备
在iOS或Android设备上，可以使用支持WebDAV的文件管理应用，如：
- iOS: Documents by Readdle, FileBrowser
- Android: Solid Explorer, FX File Explorer

## 支持的操作

### 读写模式（默认）
- ✅ 浏览文件和目录
- ✅ 下载文件
- ✅ 上传文件
- ✅ 创建目录
- ✅ 删除文件和目录
- ✅ 重命名/移动文件
- ✅ 复制文件

### 只读模式
- ✅ 浏览文件和目录
- ✅ 下载文件
- ❌ 上传文件
- ❌ 创建目录
- ❌ 删除文件和目录
- ❌ 重命名/移动文件
- ❌ 复制文件

## 实时状态检查

服务器主页会实时显示WebDAV服务的状态：
- 🔒 已禁用：WebDAV服务未启用
- ✅ 读写模式：WebDAV服务已启用，支持完整操作
- 📖 只读模式：WebDAV服务已启用，仅支持读取操作

## 安全注意事项

1. **默认禁用**：WebDAV服务默认禁用，需要明确启用
2. **网络访问**：确保只在可信网络环境中使用
3. **权限控制**：使用只读模式限制写入权限
4. **目录隔离**：使用`-webdav-dir`参数限制访问范围

## 故障排除

### 连接失败
- 确认服务器正在运行且WebDAV已启用
- 检查防火墙设置
- 确认使用正确的URL格式

### 权限错误
- 检查服务器对指定目录的读写权限
- 确认不是在只读模式下尝试写入操作

### 中文文件名问题
- 服务器支持UTF-8编码，应该能正确处理中文文件名
- 如有问题，请检查客户端的编码设置

## 技术规范

- 协议版本：WebDAV RFC 4918
- 支持的HTTP方法：OPTIONS, GET, PUT, DELETE, PROPFIND, PROPPATCH, MKCOL, COPY, MOVE, LOCK, UNLOCK
- 编码支持：UTF-8
- 锁定机制：内存锁定系统
