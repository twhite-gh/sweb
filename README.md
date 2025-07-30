# 🌐 简单Web文件服务器 v0.0.5

一个用Go语言编写的轻量级Web文件服务器，支持静态文件服务、文件上传、文件浏览、WebDAV协议、HTTPS加密连接、SOCKS5代理、HTTP代理和HTTPS代理服务。支持中英文双语界面，根据浏览器语言自动切换或手动选择。

## ✨ 项目特色

- 🚀 **零配置启动** - 开箱即用，自动创建必要的目录结构
- 🔒 **安全优先** - 高级功能默认禁用，通过命令行参数明确启用
- 🔐 **认证支持** - WebDAV、HTTP代理、HTTPS代理和SOCKS5代理支持Basic认证，保护服务安全
- 📂 **文件浏览** - 专门的文件浏览页面，支持拖拽多文件上传和进度显示
- 🌐 **WebDAV支持** - 完整的WebDAV协议实现，支持远程文件管理和认证
- 🔐 **HTTPS加密** - 支持SSL/TLS加密连接，确保数据传输安全
- 🌐 **代理服务** - 支持SOCKS5、HTTP和HTTPS代理，支持Basic认证保护
- 📱 **实时状态** - 动态显示功能状态，支持热更新
- 🎨 **现代界面** - 响应式设计，支持中英文双语界面
- 🌍 **多语言支持** - 根据浏览器语言自动切换，支持手动语言选择
- ⚙️ **灵活配置** - 支持HTTP服务开关，可仅启用HTTPS或代理服务

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
- **独立端口运行**：默认端口8081，不与HTTP服务共享端口
- 支持读写和只读两种模式
- 可配置挂载目录和自定义端口
- 支持Basic认证保护（默认用户名/密码：webdav/webdav）
- 兼容各种WebDAV客户端，特别优化Windows客户端兼容性
- Windows WebDAV客户端认证修复，支持网络驱动器挂载
- 独立的WebDAV信息页面，显示服务状态和挂载说明

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

   # 使用构建脚本编译所有平台
   # Windows
   build.bat

   # Linux/macOS
   make all

   # 或手动编译（包含所有源文件）
   go build -o sweb.exe ./src/main.go ./src/upload.go ./src/files.go ./src/webdav.go ./src/server.go ./src/utils.go ./src/socks5.go ./src/proxy.go ./src/https_proxy.go
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

# 启用WebDAV服务并开启认证（默认端口8081）
./sweb.exe -webdav -webdav-auth

# 启用WebDAV服务并指定自定义端口
./sweb.exe -webdav -webdav-port 8082

# 启用HTTPS服务
./sweb.exe -https

# 启用SOCKS5代理服务
./sweb.exe -socks5

# 启用HTTP代理服务
./sweb.exe -proxy

# 启用SOCKS5代理服务并开启认证
./sweb.exe -socks5 -socks5-auth

# 启用HTTP代理服务并开启认证
./sweb.exe -proxy -proxy-auth

# 启用HTTPS代理服务
./sweb.exe -https-proxy

# 启用HTTPS代理服务并开启认证
./sweb.exe -https-proxy -https-proxy-auth

# 仅启用HTTPS服务（关闭HTTP）
./sweb.exe -http=false -https

# 启用所有功能
./sweb.exe -upload -files -webdav -webdav-auth -https -socks5 -proxy -proxy-auth -https-proxy -https-proxy-auth

# 查看中文帮助信息（默认）
./sweb.exe -help

# 查看英文帮助信息
./sweb.exe -help-lang en -help

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
| `--webdav-port` | | WebDAV服务器端口 | 8081 |
| `--webdav-readonly` | | WebDAV服务只读模式 | 读写模式 |
| `--webdav-auth` | | 启用WebDAV认证 | 禁用 |
| `--webdav-username` | | WebDAV认证用户名 | webdav |
| `--webdav-password` | | WebDAV认证密码 | webdav |
| `--enable-http` | `-http` | 启用HTTP服务 | 启用 |
| `--https` | | 启用HTTPS服务 | 禁用 |
| `--enable-socks5` | `-socks5` | 启用SOCKS5代理服务 | 禁用 |
| `--socks5-auth` | | 启用SOCKS5代理认证 | 禁用 |
| `--socks5-username` | | SOCKS5代理认证用户名 | socks5 |
| `--socks5-password` | | SOCKS5代理认证密码 | socks5 |
| `--enable-proxy` | `-proxy` | 启用HTTP代理服务 | 禁用 |
| `--proxy-auth` | | 启用HTTP代理认证 | 禁用 |
| `--proxy-username` | | HTTP代理认证用户名 | http |
| `--proxy-password` | | HTTP代理认证密码 | http |
| `--enable-https-proxy` | `-https-proxy` | 启用HTTPS代理服务 | 禁用 |
| `--https-proxy-auth` | | 启用HTTPS代理认证 | 禁用 |
| `--https-proxy-username` | | HTTPS代理认证用户名 | https |
| `--https-proxy-password` | | HTTPS代理认证密码 | https |
| `--port` | `-p` | HTTP服务器端口 | 8080 |
| `--https-port` | | HTTPS服务器端口 | 8443 |
| `--socks5-port` | | SOCKS5代理端口 | 1080 |
| `--proxy-port` | | HTTP代理端口 | 10808 |
| `--https-proxy-port` | | HTTPS代理端口 | 10443 |
| `--cert-dir` | | SSL证书目录 | ./cert |
| `--help-lang` | | 帮助信息语言 (zh/en) | zh |
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
# 基本启用（默认端口8081）
./sweb.exe -webdav

# 指定目录和端口
./sweb.exe -webdav -webdav-dir /path/to/files -webdav-port 8082

# 只读模式
./sweb.exe -webdav -webdav-readonly

# 启用认证
./sweb.exe -webdav -webdav-auth

# 完整配置示例
./sweb.exe -webdav -webdav-auth -webdav-port 8082 -webdav-dir ./data
```

### 客户端连接

#### Windows系统
1. 打开文件资源管理器
2. 右键点击"此电脑" → "映射网络驱动器"
3. 输入地址：`http://localhost:8081/webdav`（默认端口）
4. 如果启用了认证，输入用户名：`webdav`，密码：`webdav`

#### macOS系统
1. 打开Finder
2. 按 `Cmd + K` 或选择"前往" → "连接服务器"
3. 输入地址：`http://localhost:8081/webdav`（默认端口）
4. 如果启用了认证，输入用户名：`webdav`，密码：`webdav`

#### Linux系统
```bash
# 安装davfs2
sudo apt-get install davfs2

# 挂载WebDAV（默认端口）
sudo mount -t davfs http://localhost:8081/webdav /mnt/webdav

# 如果启用了认证，会提示输入用户名和密码
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

## 🔌 SOCKS5代理使用指南

### 启用SOCKS5代理服务

```bash
# 启用SOCKS5代理服务
./sweb.exe -socks5

# 启用SOCKS5代理服务并开启认证
./sweb.exe -socks5 -socks5-auth

# 指定自定义端口
./sweb.exe -socks5 -socks5-port 1080

# 自定义认证用户名和密码
./sweb.exe -socks5 -socks5-auth -socks5-username myuser -socks5-password mypass

# 与其他功能一起使用
./sweb.exe -socks5 -upload -files
```

### 客户端配置

#### 浏览器配置
1. **Chrome/Edge**: 设置 → 高级 → 系统 → 打开代理设置
2. **Firefox**: 设置 → 网络设置 → 手动代理配置
3. 配置SOCKS5代理：`localhost:1080`

#### 系统代理配置
```bash
# Windows (PowerShell)
netsh winhttp set proxy proxy-server="socks=localhost:1080" bypass-list="localhost;127.*"

# Linux (使用proxychains)
echo "socks5 127.0.0.1 1080" >> /etc/proxychains.conf

# macOS (网络偏好设置)
# 系统偏好设置 → 网络 → 高级 → 代理 → SOCKS代理
```

### 功能特性
- **高性能**: 基于Go协程的并发处理
- **标准协议**: 完整支持SOCKS5协议（RFC 1928）
- **认证支持**: 支持无认证和用户名密码认证两种模式
- **IPv4/IPv6**: 支持双栈网络
- **域名解析**: 支持远程域名解析

## 🌐 HTTP代理使用指南

### 启用HTTP代理服务

```bash
# 启用HTTP代理服务
./sweb.exe -proxy

# 启用HTTP代理服务并开启认证
./sweb.exe -proxy -proxy-auth

# 指定自定义端口
./sweb.exe -proxy -proxy-port 10808

# 自定义认证用户名和密码
./sweb.exe -proxy -proxy-auth -proxy-username myuser -proxy-password mypass

# 与其他功能一起使用
./sweb.exe -proxy -proxy-auth -https -webdav -webdav-auth
```

### 客户端配置

#### 浏览器配置
1. 在浏览器代理设置中配置HTTP代理
2. 代理服务器：`localhost`
3. 端口：`10808`（或自定义端口）

#### 命令行工具配置
```bash
# curl使用HTTP代理（无认证）
curl --proxy http://localhost:10808 https://example.com

# curl使用HTTP代理（带认证）
curl --proxy http://http:http@localhost:10808 https://example.com

# wget使用HTTP代理（无认证）
wget --proxy=on --http-proxy=localhost:10808 https://example.com

# wget使用HTTP代理（带认证）
wget --proxy=on --http-proxy=http:http@localhost:10808 https://example.com

# 设置环境变量（无认证）
export http_proxy=http://localhost:10808
export https_proxy=http://localhost:10808

# 设置环境变量（带认证）
export http_proxy=http://http:http@localhost:10808
export https_proxy=http://http:http@localhost:10808
```

### 功能特性
- **HTTP/HTTPS支持**: 同时支持HTTP和HTTPS协议
- **CONNECT方法**: 支持HTTPS隧道代理
- **Basic认证**: 支持用户名密码认证保护代理服务
- **高并发**: 基于Go协程的高性能处理
- **透明代理**: 完整转发HTTP头部信息
- **错误处理**: 完善的错误处理和日志记录

## 🔒 HTTPS代理使用指南

### 启用HTTPS代理服务

```bash
# 启用HTTPS代理服务
./sweb.exe -https-proxy

# 启用HTTPS代理服务并开启认证
./sweb.exe -https-proxy -https-proxy-auth

# 指定自定义端口
./sweb.exe -https-proxy -https-proxy-port 10443

# 自定义认证用户名和密码
./sweb.exe -https-proxy -https-proxy-auth -https-proxy-username myuser -https-proxy-password mypass

# 与其他功能一起使用
./sweb.exe -https-proxy -upload -files
```

### 客户端配置

#### 浏览器配置
1. 在浏览器代理设置中配置HTTPS代理
2. 代理服务器：`localhost`
3. 端口：`10443`（或自定义端口）

#### 命令行工具配置
```bash
# curl使用HTTPS代理（无认证）
curl --proxy https://localhost:10443 https://example.com

# curl使用HTTPS代理（带认证）
curl --proxy https://https:https@localhost:10443 https://example.com

# 设置环境变量（无认证）
export https_proxy=https://localhost:10443

# 设置环境变量（带认证）
export https_proxy=https://https:https@localhost:10443
```

### 功能特性
- **TLS加密**: 代理连接本身使用TLS加密
- **CONNECT方法**: 支持HTTPS隧道代理
- **Basic认证**: 支持用户名密码认证保护代理服务
- **高并发**: 基于Go协程的高性能处理
- **证书支持**: 自动加载SSL证书文件
- **错误处理**: 完善的错误处理和日志记录

### 证书要求
HTTPS代理需要SSL证书文件：
- 证书文件：`./cert/server.crt`
- 私钥文件：`./cert/server.key`

可以使用以下命令生成自签名证书：
```bash
openssl req -x509 -newkey rsa:4096 -keyout cert/server.key -out cert/server.crt -days 365 -nodes
```

## 🛠️ 技术特性

- **语言**: Go语言
- **依赖**: 最小化外部依赖，主要使用标准库
- **协议**: HTTP/1.1, WebDAV RFC 4918, SOCKS5 RFC 1928
- **代理**: 高性能SOCKS5、HTTP和HTTPS代理服务器
- **编码**: UTF-8支持，完美处理中文
- **平台**: 跨平台兼容（Windows、Linux、macOS）
- **架构支持**:
  - Windows: 32位(386)、64位(amd64)
  - Linux: 32位(386)、64位(amd64)
  - macOS: Intel(amd64)、Apple Silicon(arm64)
- **部署**: 轻量级，单文件部署
- **并发**: 基于协程的高并发处理
- **构建**: 支持交叉编译，一次构建多平台

## 🌍 多语言支持

sweb支持中英文双语界面，提供完整的国际化体验：

### 🔧 语言设置

- **自动检测**: 根据浏览器语言自动选择界面语言
- **手动切换**: 页面右上角提供语言切换按钮
- **本地存储**: 用户选择的语言偏好会保存在浏览器中

### 📖 帮助信息语言

```bash
# 显示中文帮助信息（默认）
./sweb.exe -help

# 显示英文帮助信息
./sweb.exe -help-lang en -help
```

### 🌐 Web界面语言

- **中文界面**: 浏览器语言为中文时自动显示
- **英文界面**: 其他语言浏览器自动显示英文界面
- **语言切换**: 点击页面右上角的"中文"/"English"按钮切换
- **状态同步**: 所有功能状态信息都支持双语显示

### 📱 支持的语言

| 语言 | 代码 | 支持范围 |
|------|------|----------|
| 中文 | zh | 完整支持（帮助信息、Web界面） |
| 英文 | en | 完整支持（帮助信息、Web界面） |

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
  },
  "socks5": {
    "enabled": true,
    "status": "running",
    "port": 1080,
    "totalConnections": 25,
    "activeConnections": 3,
    "successfulConnections": 22,
    "failedConnections": 0,
    "bytesTransferred": 1048576
  },
  "proxy": {
    "enabled": true,
    "status": "running",
    "port": 10808,
    "totalConnections": 15,
    "activeConnections": 2,
    "successfulConnections": 13,
    "failedConnections": 0,
    "httpRequests": 8,
    "httpsRequests": 5,
    "bytesTransferred": 2097152
  }
}
```

## 🔒 安全考虑

### 默认安全策略
- 文件上传功能默认**禁用**
- 文件浏览功能默认**禁用**
- WebDAV服务默认**禁用**
- WebDAV认证默认**禁用**
- HTTPS服务默认**禁用**
- SOCKS5代理服务默认**禁用**
- SOCKS5代理认证默认**禁用**
- HTTP代理服务默认**禁用**
- HTTP代理认证默认**禁用**
- HTTPS代理服务默认**禁用**
- HTTPS代理认证默认**禁用**
- 需要通过命令行参数明确启用高级功能

### 权限控制
- WebDAV支持只读模式和Basic认证，已优化Windows客户端兼容性
- 可限制WebDAV访问目录范围
- HTTP代理支持Basic认证保护
- HTTPS代理支持Basic认证保护和TLS加密
- SOCKS5代理支持Basic认证保护
- 建议在生产环境中启用认证功能

### 最佳实践
```bash
# 生产环境推荐配置（只读WebDAV + HTTPS）
./sweb.exe -webdav -webdav-readonly -webdav-dir /safe/directory -https

# 开发环境配置（启用所有功能）
./sweb.exe -upload -files -webdav -https -socks5 -proxy

# 文件共享配置（只启用文件浏览）
./sweb.exe -files

# 代理服务器配置（仅代理功能）
./sweb.exe -socks5 -proxy

# 内网代理配置（指定端口）
./sweb.exe -socks5 -socks5-port 1080 -proxy -proxy-port 8080
```

## 📁 目录结构

```
sweb/
├── src/                        # 源代码目录
│   ├── main.go                 # 主程序入口文件
│   ├── upload.go               # 文件上传功能模块
│   ├── files.go                # 文件浏览功能模块
│   ├── webdav.go               # WebDAV服务模块
│   ├── server.go               # HTTP/HTTPS服务器模块
│   ├── utils.go                # 工具函数和页面生成
│   ├── socks5.go               # SOCKS5代理服务器模块
│   ├── proxy.go                # HTTP代理服务器模块
│   └── https_proxy.go          # HTTPS代理服务器模块
├── test/                       # 测试代码目录
│   ├── test_all.bat            # Windows完整测试套件
│   ├── test_all.sh             # Linux/macOS完整测试套件
│   ├── test_all.go             # 综合功能测试
│   ├── test_socks5.go          # SOCKS5代理测试
│   ├── test_proxy.go           # HTTP代理测试
│   ├── debug_socks5.go         # SOCKS5调试工具
│   └── README.md               # 测试说明文档
├── doc/                        # 文档目录
│   ├── 项目结构说明.md         # 项目结构详细说明
│   ├── HTTPS代理配置指南.md    # HTTPS代理配置指南
│   ├── HTTPS功能说明.md        # HTTPS功能说明
│   ├── WebDAV使用说明.md       # WebDAV详细说明
│   ├── 测试说明.md             # 测试说明文档
│   ├── 浏览器配置说明.md       # 浏览器配置说明
│   └── 目录重构完成说明.md     # 目录重构说明
├── tools/                      # 工具目录
│   └── generate-cert.go        # SSL证书生成工具
├── bin/                        # 编译输出目录
│   ├── windows/                # Windows可执行文件
│   ├── linux/                  # Linux可执行文件
│   └── macos/                  # macOS可执行文件
├── cert/                       # SSL证书目录（HTTPS功能需要）
│   ├── server.crt              # SSL证书文件
│   └── server.key              # SSL私钥文件
├── web/                        # Web文件目录（自动创建）
│   └── index.html              # 默认主页（自动生成）
├── go.mod                      # Go模块文件
├── go.sum                      # 依赖校验文件
├── Makefile                    # 构建脚本（Linux/macOS）
├── build.bat                   # 构建脚本（Windows）
└── README.md                   # 项目说明
```

## 🎯 使用场景

- **开发测试**: 快速搭建本地文件服务器
- **文件共享**: 在局域网内安全共享文件
- **远程管理**: 通过WebDAV远程管理文件
- **静态网站**: 托管简单的静态网站
- **文件备份**: 作为简单的文件上传服务
- **团队协作**: 通过文件浏览功能快速查看和下载文件
- **安全传输**: 使用HTTPS确保文件传输安全
- **网络代理**: 提供SOCKS5和HTTP代理服务
- **内网穿透**: 在受限网络环境中提供代理访问
- **开发调试**: 通过代理服务器调试网络请求
- **流量转发**: 高性能的TCP和HTTP流量转发

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
- 实时显示所有功能状态（上传、文件浏览、WebDAV、HTTPS、SOCKS5、HTTP代理）
- 响应式设计，支持移动设备
- 中文界面，操作简单
- 动态按钮状态，根据功能启用情况显示
- 代理服务状态和统计信息展示

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
# 使用Makefile（推荐）
make all                    # 编译所有平台（Windows、Linux、macOS）
make windows               # 仅编译Windows版本（64位和32位）
make linux                 # 仅编译Linux版本（64位和32位）
make macos                 # 仅编译macOS版本（Intel和Apple Silicon）

# 单独编译特定架构
make windows-64            # Windows 64位
make windows-32            # Windows 32位
make linux-64              # Linux 64位
make linux-32              # Linux 32位
make macos-intel           # macOS Intel
make macos-arm64           # macOS Apple Silicon

# 使用build.bat（Windows）
build.bat                  # 编译所有平台版本

# 手动编译（包含所有源文件）
go build -o sweb ./src/main.go ./src/upload.go ./src/files.go ./src/webdav.go ./src/server.go ./src/utils.go ./src/socks5.go ./src/proxy.go ./src/https_proxy.go

# 交叉编译示例
# Windows 64位
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o bin/windows/sweb-windows-amd64.exe ./src/*.go

# Windows 32位
GOOS=windows GOARCH=386 go build -ldflags "-s -w" -o bin/windows/sweb-windows-386.exe ./src/*.go

# Linux 64位
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o bin/linux/sweb-linux-amd64 ./src/*.go

# Linux 32位
GOOS=linux GOARCH=386 go build -ldflags "-s -w" -o bin/linux/sweb-linux-386 ./src/*.go

# macOS Intel
GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o bin/macos/sweb-darwin-amd64 ./src/*.go

# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o bin/macos/sweb-darwin-arm64 ./src/*.go
```

### 运行测试
```bash
# 进入测试目录
cd test

# Windows - 运行完整测试套件
test_all.bat

# Linux/macOS - 运行完整测试套件
chmod +x test_all.sh
./test_all.sh

# 单独运行测试
go run test_all.go          # 综合功能测试
go run test_socks5.go       # SOCKS5代理测试
go run test_proxy.go        # HTTP代理测试
go run debug_socks5.go      # SOCKS5调试工具
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
- 确认使用正确的URL格式：`http://localhost:8081/webdav`（默认端口）

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

#### 7. SOCKS5代理连接失败
- 确认SOCKS5代理功能已启用（使用`-socks5`参数）
- 检查端口是否被占用：`netstat -an | grep :1080`
- 确认客户端SOCKS5配置正确
- 检查防火墙是否阻止连接

#### 8. HTTP代理无法访问网站
- 确认HTTP代理功能已启用（使用`-proxy`参数）
- 检查端口是否被占用：`netstat -an | grep :10808`
- 确认浏览器代理设置正确
- 检查目标网站是否可访问

### 日志调试
服务器会输出详细的操作日志，包括：
- 功能启用/禁用状态
- WebDAV操作记录
- 文件上传状态
- HTTPS证书状态
- SOCKS5代理连接和数据转发
- HTTP代理请求处理
- 错误信息和性能统计

## 🔄 版本历史

### v0.0.5 (当前版本)
- ✅ WebDAV服务Basic认证支持，优化Windows客户端兼容性
- ✅ HTTP代理服务Basic认证支持
- ✅ HTTPS代理服务独立实现，支持TLS加密和Basic认证
- ✅ SOCKS5代理服务Basic认证支持
- ✅ HTTP服务开关功能（可仅启用HTTPS或代理服务）
- ✅ WebDAV Windows挂载认证修复
- ✅ 认证功能测试套件
- ✅ 多语言支持（中英文双语界面）
- ✅ 更新文档和帮助信息

### v0.0.4
- ✅ 基础静态文件服务
- ✅ 拖拽多文件上传功能（带进度条）
- ✅ 专门的文件浏览页面
- ✅ 完整WebDAV协议支持
- ✅ HTTPS加密连接支持
- ✅ 高性能SOCKS5代理服务器
- ✅ HTTP/HTTPS代理服务器
- ✅ 实时状态监控和统计
- ✅ 响应式Web界面
- ✅ 命令行参数配置
- ✅ 安全默认策略
- ✅ 模块化代码结构
- ✅ 重构目录结构（src/、test/、doc/、tools/）
- ✅ 多平台编译支持（Windows/Linux/macOS，32位和64位）
- ✅ 完整的测试套件和文档

### v0.0.3
- ✅ 基础文件服务功能
- ✅ WebDAV和HTTPS支持
- ✅ Web界面和状态监控

### v0.0.2
- ✅ 文件上传和浏览功能
- ✅ 代理服务器支持

### v0.0.1
- ✅ 初始版本
- ✅ 基础HTTP服务器

### 计划功能
- [ ] 用户认证和权限管理系统
- [x] ~~代理服务认证机制~~ (v0.0.5已实现)
- [ ] 文件预览功能
- [ ] 批量文件操作
- [ ] 配置文件支持
- [ ] 文件搜索功能
- [ ] 文件版本管理
- [ ] 代理服务访问控制和白名单

## 🔌 默认端口说明

| 服务类型 | 默认端口 | 说明 |
|---------|---------|------|
| HTTP服务 | 8080 | 主要的Web服务端口 |
| **WebDAV服务** | **8081** | **独立的WebDAV服务端口** |
| HTTPS服务 | 8443 | 加密的Web服务端口 |
| SOCKS5代理 | 1080 | SOCKS5代理服务端口 |
| HTTP代理 | 10808 | HTTP代理服务端口 |
| HTTPS代理 | 10443 | HTTPS代理服务端口 |

### 重要说明
- **WebDAV服务现在运行在独立的8081端口**，不再与HTTP服务共享8080端口
- 所有端口都可以通过相应的命令行参数自定义
- WebDAV独立端口设计提高了服务的稳定性和安全性
- 访问WebDAV服务请使用：`http://localhost:8081/webdav`

## 📞 联系方式

如果您有任何问题或建议，请通过以下方式联系：

- 提交 [Issue](../../issues)
- 发起 [Discussion](../../discussions)

## ⭐ 支持项目

如果这个项目对您有帮助，请考虑给它一个星标 ⭐

---

**简单Web文件服务器** - 让文件管理变得简单高效！ 🚀
