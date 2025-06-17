# 🌐 简单Web文件服务器

一个用Go语言编写的轻量级Web文件服务器，支持静态文件服务、文件上传、文件浏览、WebDAV协议和HTTPS加密连接。

## ✨ 项目特色

- 🚀 **零配置启动** - 开箱即用，自动创建必要的目录结构
- 🔒 **安全优先** - 高级功能默认禁用，通过命令行参数明确启用
- 📂 **文件浏览** - 专门的文件浏览页面，支持拖拽多文件上传和进度显示
- 🌐 **WebDAV支持** - 完整的WebDAV协议实现，支持远程文件管理
- 🔐 **HTTPS加密** - 支持SSL/TLS加密连接，确保数据传输安全
- 📱 **实时状态** - 动态显示功能状态，支持热更新
- 🎨 **现代界面** - 响应式设计，支持中文界面

## 📋 主要功能

### 📁 静态文件服务
自动服务web目录下的所有文件，支持HTML、CSS、JavaScript、图片等各种文件类型。

### 📤 文件上传
- 支持拖拽多文件上传到Web界面
- 实时显示上传进度条
- 支持各种文件格式
- 安全考虑：默认禁用，需要明确启用

### 📂 文件浏览
- 专门的文件浏览页面，展示web目录下的所有文件
- 支持文件类型图标识别和在线预览
- 点击文件直接下载，支持文件大小和修改时间显示
- 安全考虑：默认禁用，需要明确启用

### 🌐 WebDAV服务
- 完整的WebDAV协议支持（RFC 4918）
- 支持读写和只读两种模式
- 可配置挂载目录
- 兼容各种WebDAV客户端

### 🔐 HTTPS服务
- 支持SSL/TLS加密连接
- 自动从cert目录加载证书文件
- 同时提供HTTP和HTTPS服务
- 实时显示证书状态

### 🔧 自动配置
自动创建必要的目录结构，无需手动配置即可使用。

## 🚀 快速开始

### 下载和安装

1. 从[Releases](../../releases)页面下载适合您系统的可执行文件
2. 或者从源码编译：
   ```bash
   git clone <repository-url>
   cd sweb
   go build -o sweb.exe main.go
   ```

### 基本使用

```bash
# 启动基本文件服务器
./sweb.exe

# 启用文件上传功能
./sweb.exe -upload

# 启用文件浏览功能
./sweb.exe -files

# 启用WebDAV服务
./sweb.exe -webdav

# 启用HTTPS服务
./sweb.exe -https

# 启用所有功能
./sweb.exe -upload -files -webdav -https

# 指定端口
./sweb.exe -port 9000

# 查看帮助
./sweb.exe -help
```

## 📖 命令行参数

| 参数 | 简写 | 说明 | 默认值 |
|------|------|------|--------|
| `--enable-upload` | `-upload` | 启用文件上传功能 | 禁用 |
| `--enable-files` | `-files` | 启用文件浏览功能 | 禁用 |
| `--enable-webdav` | `-webdav` | 启用WebDAV服务 | 禁用 |
| `--webdav-dir` | | WebDAV服务的根目录 | 当前目录 |
| `--webdav-readonly` | | WebDAV服务只读模式 | 读写模式 |
| `--https` | | 启用HTTPS服务 | 禁用 |
| `--http-port` | | HTTP服务器端口 | 8080 |
| `--https-port` | | HTTPS服务器端口 | 8443 |
| `--cert-dir` | | SSL证书目录 | ./cert |
| `--help` | `-h` | 显示帮助信息 | |

## 📂 文件浏览功能

### 启用文件浏览

```bash
# 启用文件浏览功能
./sweb.exe -files

# 同时启用文件浏览和上传
./sweb.exe -files -upload
```

### 功能特性
- **文件列表**: 网格布局展示所有文件，支持文件类型图标
- **文件信息**: 显示文件大小、修改时间和文件类型
- **在线下载**: 点击文件直接下载
- **实时刷新**: 支持手动刷新文件列表

## 🌐 WebDAV使用指南

### 启用WebDAV服务

```bash
# 基本启用
./sweb.exe -webdav

# 指定目录
./sweb.exe -webdav -webdav-dir /path/to/files

# 只读模式
./sweb.exe -webdav -webdav-readonly
```

### 客户端连接

#### Windows系统
1. 打开文件资源管理器
2. 右键点击"此电脑" → "映射网络驱动器"
3. 输入地址：`http://localhost:8080/webdav`

#### macOS系统
1. 打开Finder
2. 按 `Cmd + K` 或选择"前往" → "连接服务器"
3. 输入地址：`http://localhost:8080/webdav`

#### Linux系统
```bash
# 安装davfs2
sudo apt-get install davfs2

# 挂载WebDAV
sudo mount -t davfs http://localhost:8080/webdav /mnt/webdav
```

#### 移动设备
使用支持WebDAV的文件管理应用：
- **iOS**: Documents by Readdle, FileBrowser
- **Android**: Solid Explorer, FX File Explorer

## 🔐 HTTPS使用指南

### 启用HTTPS服务

```bash
# 启用HTTPS服务
./sweb.exe -https

# 指定证书目录
./sweb.exe -https -cert-dir /path/to/certs

# 指定端口
./sweb.exe -https -http-port 8080 -https-port 8443
```

### 证书配置
1. 在当前目录创建 `cert` 文件夹
2. 将SSL证书文件放入cert目录：
   - `server.crt` - 证书文件
   - `server.key` - 私钥文件
3. 启动服务器时会自动加载证书

### 生成自签名证书
```bash
# 使用OpenSSL生成自签名证书
openssl req -x509 -newkey rsa:4096 -keyout cert/server.key -out cert/server.crt -days 365 -nodes
```

## 🛠️ 技术特性

- **语言**: Go语言
- **依赖**: 最小化外部依赖，主要使用标准库
- **协议**: HTTP/1.1, WebDAV RFC 4918
- **编码**: UTF-8支持，完美处理中文
- **平台**: 跨平台兼容（Windows、Linux、macOS）
- **架构**: 轻量级，单文件部署

## 📊 实时状态监控

服务器提供实时状态API和Web界面：

- **状态API**: `GET /api/status`
- **Web界面**: 主页自动显示当前功能状态
- **自动更新**: 每30秒检查一次状态变化

### API响应示例

```json
{
  "upload": {
    "enabled": true,
    "status": "enabled"
  },
  "files": {
    "enabled": true,
    "status": "enabled"
  },
  "webdav": {
    "enabled": true,
    "readonly": false,
    "directory": "./files",
    "status": "enabled-readwrite"
  },
  "https": {
    "enabled": true,
    "httpPort": 8080,
    "httpsPort": 8443,
    "certDir": "./cert",
    "certStatus": "enabled"
  }
}
```

## 🔒 安全考虑

### 默认安全策略
- 文件上传功能默认**禁用**
- 文件浏览功能默认**禁用**
- WebDAV服务默认**禁用**
- HTTPS服务默认**禁用**
- 需要通过命令行参数明确启用高级功能

### 权限控制
- WebDAV支持只读模式
- 可限制WebDAV访问目录范围
- 建议在可信网络环境中使用

### 最佳实践
```bash
# 生产环境推荐配置（只读WebDAV + HTTPS）
./sweb.exe -webdav -webdav-readonly -webdav-dir /safe/directory -https

# 开发环境配置（启用所有功能）
./sweb.exe -upload -files -webdav -https

# 文件共享配置（只启用文件浏览）
./sweb.exe -files
```

## 📁 目录结构

```
sweb/
├── main.go                 # 主程序文件
├── go.mod                  # Go模块文件
├── go.sum                  # 依赖校验文件
├── README.md               # 项目说明
├── WebDAV使用说明.md       # WebDAV详细说明
├── web/                    # Web文件目录（自动创建）
│   └── index.html          # 默认主页（自动生成）
├── cert/                   # SSL证书目录（HTTPS功能需要）
│   ├── server.crt          # SSL证书文件
│   └── server.key          # SSL私钥文件
└── test-webdav/           # WebDAV测试目录
    ├── readme.txt
    └── sample.json
```

## 🎯 使用场景

- **开发测试**: 快速搭建本地文件服务器
- **文件共享**: 在局域网内安全共享文件
- **远程管理**: 通过WebDAV远程管理文件
- **静态网站**: 托管简单的静态网站
- **文件备份**: 作为简单的文件上传服务
- **团队协作**: 通过文件浏览功能快速查看和下载文件
- **安全传输**: 使用HTTPS确保文件传输安全

## 🤝 贡献指南

欢迎提交Issue和Pull Request！

1. Fork本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启Pull Request

## 📄 许可证

本项目采用MIT许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

- [golang.org/x/net/webdav](https://pkg.go.dev/golang.org/x/net/webdav) - WebDAV协议实现
- Go语言标准库 - 提供了强大的HTTP服务器功能

## 📸 功能截图

### 主页界面
- 实时显示所有功能状态（上传、文件浏览、WebDAV、HTTPS）
- 响应式设计，支持移动设备
- 中文界面，操作简单
- 动态按钮状态，根据功能启用情况显示

### 文件浏览页面
- 网格布局展示文件，支持文件类型图标
- 显示文件大小、修改时间等详细信息
- 点击文件直接下载，支持实时刷新

### 文件上传页面
- 支持拖拽多文件上传
- 实时显示上传进度条
- 文件列表管理，支持添加/删除文件

### 功能状态显示
- 🔒 已禁用：功能未启用
- ✅ 已启用：功能正常工作
- 📖 只读模式：WebDAV只读模式
- ⚠️ 证书错误：HTTPS证书有问题

## 🔧 开发信息

### 构建要求
- Go 1.19 或更高版本
- 网络连接（用于下载依赖）

### 依赖包
```go
require (
    golang.org/x/net v0.x.x // WebDAV协议支持
)
```

### 编译命令
```bash
# 当前平台
go build -o sweb main.go

# 交叉编译
# Windows
GOOS=windows GOARCH=amd64 go build -o sweb.exe main.go

# Linux
GOOS=linux GOARCH=amd64 go build -o sweb main.go

# macOS
GOOS=darwin GOARCH=amd64 go build -o sweb main.go
```

## 🐛 故障排除

### 常见问题

#### 1. 端口被占用
```bash
# 检查端口占用
netstat -an | grep :8080

# 使用其他端口
./sweb.exe -port 9000
```

#### 2. WebDAV连接失败
- 确认服务器正在运行且WebDAV已启用
- 检查防火墙设置
- 确认使用正确的URL格式：`http://localhost:8080/webdav`

#### 3. 文件上传失败
- 确认上传功能已启用（使用`-upload`参数）
- 检查目标目录的写入权限
- 确认磁盘空间充足

#### 4. 文件浏览页面无法访问
- 确认文件浏览功能已启用（使用`-files`参数）
- 检查web目录是否存在
- 确认有读取权限

#### 5. HTTPS连接失败
- 确认HTTPS功能已启用（使用`-https`参数）
- 检查cert目录下是否有server.crt和server.key文件
- 确认证书文件格式正确

#### 6. 中文文件名显示异常
- 服务器支持UTF-8编码
- 检查客户端的编码设置
- 确保文件系统支持Unicode

### 日志调试
服务器会输出详细的操作日志，包括：
- 功能启用/禁用状态
- WebDAV操作记录
- 文件上传状态
- HTTPS证书状态
- 错误信息

## 🔄 版本历史

### v1.0.0 (当前版本)
- ✅ 基础静态文件服务
- ✅ 拖拽多文件上传功能（带进度条）
- ✅ 专门的文件浏览页面
- ✅ 完整WebDAV协议支持
- ✅ HTTPS加密连接支持
- ✅ 实时状态监控
- ✅ 响应式Web界面
- ✅ 命令行参数配置
- ✅ 安全默认策略

### 计划功能
- [ ] 用户认证和权限管理
- [ ] 文件预览功能
- [ ] 批量文件操作
- [ ] 配置文件支持
- [ ] 文件搜索功能
- [ ] 文件版本管理

## 📞 联系方式

如果您有任何问题或建议，请通过以下方式联系：

- 提交 [Issue](../../issues)
- 发起 [Discussion](../../discussions)

## ⭐ 支持项目

如果这个项目对您有帮助，请考虑给它一个星标 ⭐

---

**简单Web文件服务器** - 让文件管理变得简单高效！ 🚀
