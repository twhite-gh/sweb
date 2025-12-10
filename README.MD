# 🌐 Simple Web File Server v0.0.5

A lightweight Web file server written in Go, supporting static file serving, file upload, file browsing, WebDAV protocol, HTTPS encrypted connection, SOCKS5 proxy, HTTP proxy, and HTTPS proxy services. Supports bilingual interface (Chinese/English), with automatic switching based on browser language or manual selection.

## ✨ Project Features

  * 🚀 **Zero-Configuration Startup** - Ready to use out of the box, automatically creates the necessary directory structure.
  * 🔒 **Security First** - Advanced features are disabled by default and must be explicitly enabled via command-line arguments.
  * 🔐 **Authentication Support** - WebDAV, HTTP Proxy, HTTPS Proxy, and SOCKS5 Proxy support Basic authentication to protect service security.
  * 📂 **File Browsing** - Dedicated file browsing page, supports drag-and-drop multiple file uploads and progress display.
  * 🌐 **WebDAV Support** - Complete WebDAV protocol implementation, supporting remote file management and authentication.
  * 🔐 **HTTPS Encryption** - Supports SSL/TLS encrypted connections, ensuring data transfer security.
  * 🌐 **Proxy Services** - Supports SOCKS5, HTTP, and HTTPS proxies, protected by Basic authentication.
  * 📱 **Real-time Status** - Dynamically displays feature status, supports hot updates.
  * 🎨 **Modern Interface** - Responsive design, supports bilingual interface (Chinese/English).
  * 🌍 **Multilingual Support** - Automatic switching based on browser language, supports manual language selection.
  * ⚙️ **Flexible Configuration** - Supports HTTP service switch, allowing only HTTPS or proxy services to be enabled.

## 📋 Key Features

### 📁 Static File Serving

Automatically serves all files in the `web` directory, supporting various file types like HTML, CSS, JavaScript, images, etc.

### 📤 File Upload

  * Supports drag-and-drop multiple file uploads to the Web interface.
  * Displays real-time upload progress bar.
  * Supports various file formats.
  * Security consideration: **Disabled by default**, requires explicit enabling.

### 📂 File Browsing

  * Dedicated file browsing page, displaying all files in the `web` directory.
  * Supports file type icon recognition and online preview.
  * Clicking a file directly downloads it; displays file size and modification time.
  * Security consideration: **Disabled by default**, requires explicit enabling.

### 🌐 WebDAV Service

  * Complete WebDAV protocol support (RFC 4918).
  * **Runs on an independent port**: Default port is **8081**, not shared with the HTTP service.
  * Supports both read/write and read-only modes.
  * Configurable mount directory and custom port.
  * Supports Basic authentication protection (Default username/password: `webdav`/`webdav`).
  * Compatible with various WebDAV clients, with special optimization for Windows client compatibility.
  * Fixes for Windows WebDAV client authentication, supporting network drive mapping.
  * Independent WebDAV information page, displaying service status and mounting instructions.

### 🔐 HTTPS Service

  * Supports SSL/TLS encrypted connections.
  * Automatically loads certificate files from the `cert` directory.
  * Provides both HTTP and HTTPS services simultaneously.
  * Real-time display of certificate status.

### 🔧 Automatic Configuration

Automatically creates the necessary directory structure, making it usable without manual configuration.

## 🚀 Quick Start

### Download and Installation

1.  Download the executable file suitable for your system from the [可疑链接已删除] page.
2.  Alternatively, compile from source:
    ```bash
    git clone <repository-url>
    cd sweb

    # Use the build script to compile for all platforms
    # Windows
    build.bat

    # Linux/macOS
    make all

    # Or manually compile (including all source files)
    go build -o sweb.exe ./src/main.go ./src/upload.go ./src/files.go ./src/webdav.go ./src/server.go ./src/utils.go ./src/socks5.go ./src/proxy.go ./src/https_proxy.go
    ```

### Basic Usage

```bash
# Start the basic file server
./sweb.exe

# Enable file upload feature
./sweb.exe -upload

# Enable file browsing feature
./sweb.exe -files

# Enable WebDAV service
./sweb.exe -webdav

# Enable WebDAV service and turn on authentication (default port 8081)
./sweb.exe -webdav -webdav-auth

# Enable WebDAV service and specify a custom port
./sweb.exe -webdav -webdav-port 8082

# Enable HTTPS service
./sweb.exe -https

# Enable SOCKS5 proxy service
./sweb.exe -socks5

# Enable HTTP proxy service
./sweb.exe -proxy

# Enable SOCKS5 proxy service and turn on authentication
./sweb.exe -socks5 -socks5-auth

# Enable HTTP proxy service and turn on authentication
./sweb.exe -proxy -proxy-auth

# Enable HTTPS proxy service
./sweb.exe -https-proxy

# Enable HTTPS proxy service and turn on authentication
./sweb.exe -https-proxy -https-proxy-auth

# Only enable HTTPS service (turn off HTTP)
./sweb.exe -http=false -https

# Enable all features
./sweb.exe -upload -files -webdav -webdav-auth -https -socks5 -proxy -proxy-auth -https-proxy -https-proxy-auth

# View Chinese help information (default)
./sweb.exe -help

# View English help information
./sweb.exe -help-lang en -help

# Specify port
./sweb.exe -port 9000

# View help
./sweb.exe -help
```

## 📖 Command Line Arguments

| Parameter | Shorthand | Description | Default Value |
| :--- | :--- | :--- | :--- |
| `--enable-upload` | `-upload` | Enable file upload feature | Disabled |
| `--enable-files` | `-files` | Enable file browsing feature | Disabled |
| `--enable-webdav` | `-webdav` | Enable WebDAV service | Disabled |
| `--webdav-dir` | | Root directory for WebDAV service | Current directory |
| `--webdav-port` | | WebDAV server port | 8081 |
| `--webdav-readonly` | | WebDAV service read-only mode | Read/Write mode |
| `--webdav-auth` | | Enable WebDAV authentication | Disabled |
| `--webdav-username` | | WebDAV authentication username | webdav |
| `--webdav-password` | | WebDAV authentication password | webdav |
| `--enable-http` | `-http` | Enable HTTP service | Enabled |
| `--https` | | Enable HTTPS service | Disabled |
| `--enable-socks5` | `-socks5` | Enable SOCKS5 proxy service | Disabled |
| `--socks5-auth` | | Enable SOCKS5 proxy authentication | Disabled |
| `--socks5-username` | | SOCKS5 proxy authentication username | socks5 |
| `--socks5-password` | | SOCKS5 proxy authentication password | socks5 |
| `--enable-proxy` | `-proxy` | Enable HTTP proxy service | Disabled |
| `--proxy-auth` | | Enable HTTP proxy authentication | Disabled |
| `--proxy-username` | | HTTP proxy authentication username | http |
| `--proxy-password` | | HTTP proxy authentication password | http |
| `--enable-https-proxy` | `-https-proxy` | Enable HTTPS proxy service | Disabled |
| `--https-proxy-auth` | | Enable HTTPS proxy authentication | Disabled |
| `--https-proxy-username` | | HTTPS proxy authentication username | https |
| `--https-proxy-password` | | HTTPS proxy authentication password | https |
| `--port` | `-p` | HTTP server port | 8080 |
| `--https-port` | | HTTPS server port | 8443 |
| `--socks5-port` | | SOCKS5 proxy port | 1080 |
| `--proxy-port` | | HTTP proxy port | 10808 |
| `--https-proxy-port` | | HTTPS proxy port | 10443 |
| `--cert-dir` | | SSL certificate directory | ./cert |
| `--help-lang` | | Help message language (zh/en) | zh |
| `--help` | `-h` | Display help information | |

## 📂 File Browsing Feature

### Enabling File Browsing

```bash
# Enable file browsing feature
./sweb.exe -files

# Enable both file browsing and upload
./sweb.exe -files -upload
```

### Feature Characteristics

  * **File List**: Grid layout displaying all files, supports file type icons.
  * **File Information**: Displays file size, modification time, and file type.
  * **Online Download**: Clicking a file directly downloads it.
  * **Real-time Refresh**: Supports manual refreshing of the file list.

## 🌐 WebDAV Usage Guide

### Enabling WebDAV Service

```bash
# Basic enablement (default port 8081)
./sweb.exe -webdav

# Specify directory and port
./sweb.exe -webdav -webdav-dir /path/to/files -webdav-port 8082

# Read-only mode
./sweb.exe -webdav -webdav-readonly

# Enable authentication
./sweb.exe -webdav -webdav-auth

# Full configuration example
./sweb.exe -webdav -webdav-auth -webdav-port 8082 -webdav-dir ./data
```

### Client Connection

#### Windows System

1.  Open File Explorer.
2.  Right-click "This PC" → "Map network drive".
3.  Enter the address: `http://localhost:8081/webdav` (default port).
4.  If authentication is enabled, enter username: `webdav`, password: `webdav`.

#### macOS System

1.  Open Finder.
2.  Press `Cmd + K` or select "Go" → "Connect to Server".
3.  Enter the address: `http://localhost:8081/webdav` (default port).
4.  If authentication is enabled, enter username: `webdav`, password: `webdav`.

#### Linux System

```bash
# Install davfs2
sudo apt-get install davfs2

# Mount WebDAV (default port)
sudo mount -t davfs http://localhost:8081/webdav /mnt/webdav

# If authentication is enabled, you will be prompted for username and password
```

#### Mobile Devices

Use file management applications that support WebDAV:

  * **iOS**: Documents by Readdle, FileBrowser
  * **Android**: Solid Explorer, FX File Explorer

## 🔐 HTTPS Usage Guide

### Enabling HTTPS Service

```bash
# Enable HTTPS service
./sweb.exe -https

# Specify certificate directory
./sweb.exe -https -cert-dir /path/to/certs

# Specify ports
./sweb.exe -https -http-port 8080 -https-port 8443
```

### Certificate Configuration

1.  Create a `cert` folder in the current directory.
2.  Place SSL certificate files into the `cert` directory:
      * `server.crt` - Certificate file
      * `server.key` - Private key file
3.  The server will automatically load the certificates upon startup.

### Generating Self-Signed Certificates

```bash
# Use OpenSSL to generate a self-signed certificate
openssl req -x509 -newkey rsa:4096 -keyout cert/server.key -out cert/server.crt -days 365 -nodes
```

## 🔌 SOCKS5 Proxy Usage Guide

### Enabling SOCKS5 Proxy Service

```bash
# Enable SOCKS5 proxy service
./sweb.exe -socks5

# Enable SOCKS5 proxy service and turn on authentication
./sweb.exe -socks5 -socks5-auth

# Specify custom port
./sweb.exe -socks5 -socks5-port 1080

# Custom authentication username and password
./sweb.exe -socks5 -socks5-auth -socks5-username myuser -socks5-password mypass

# Use with other features
./sweb.exe -socks5 -upload -files
```

### Client Configuration

#### Browser Configuration

1.  **Chrome/Edge**: Settings → Advanced → System → Open proxy settings
2.  **Firefox**: Settings → Network Settings → Manual proxy configuration
3.  Configure SOCKS5 proxy: `localhost:1080`

#### System Proxy Configuration

```bash
# Windows (PowerShell)
netsh winhttp set proxy proxy-server="socks=localhost:1080" bypass-list="localhost;127.*"

# Linux (using proxychains)
echo "socks5 127.0.0.1 1080" >> /etc/proxychains.conf

# macOS (Network Preferences)
# System Preferences → Network → Advanced → Proxies → SOCKS Proxy
```

### Feature Characteristics

  * **High Performance**: Concurrent processing based on Go routines.
  * **Standard Protocol**: Complete support for the SOCKS5 protocol (RFC 1928).
  * **Authentication Support**: Supports both no authentication and username/password authentication modes.
  * **IPv4/IPv6**: Supports dual-stack networks.
  * **Domain Resolution**: Supports remote domain resolution.

## 🌐 HTTP Proxy Usage Guide

### Enabling HTTP Proxy Service

```bash
# Enable HTTP proxy service
./sweb.exe -proxy

# Enable HTTP proxy service and turn on authentication
./sweb.exe -proxy -proxy-auth

# Specify custom port
./sweb.exe -proxy -proxy-port 10808

# Custom authentication username and password
./sweb.exe -proxy -proxy-auth -proxy-username myuser -proxy-password mypass

# Use with other features
./sweb.exe -proxy -proxy-auth -https -webdav -webdav-auth
```

### Client Configuration

#### Browser Configuration

1.  Configure HTTP proxy in browser proxy settings.
2.  Proxy server: `localhost`
3.  Port: `10808` (or custom port)

#### Command Line Tool Configuration

```bash
# curl using HTTP proxy (no authentication)
curl --proxy http://localhost:10808 https://example.com

# curl using HTTP proxy (with authentication)
curl --proxy http://http:http@localhost:10808 https://example.com

# wget using HTTP proxy (no authentication)
wget --proxy=on --http-proxy=localhost:10808 https://example.com

# wget using HTTP proxy (with authentication)
wget --proxy=on --http-proxy=http:http@localhost:10808 https://example.com

# Set environment variables (no authentication)
export http_proxy=http://localhost:10808
export https_proxy=http://localhost:10808

# Set environment variables (with authentication)
export http_proxy=http://http:http@localhost:10808
export https_proxy=http://http:http@localhost:10808
```

### Feature Characteristics

  * **HTTP/HTTPS Support**: Supports both HTTP and HTTPS protocols.
  * **CONNECT Method**: Supports HTTPS tunnel proxy.
  * **Basic Authentication**: Supports username/password authentication to protect the proxy service.
  * **High Concurrency**: High-performance processing based on Go routines.
  * **Transparent Proxy**: Completely forwards HTTP header information.
  * **Error Handling**: Comprehensive error handling and logging.

## 🔒 HTTPS Proxy Usage Guide

### Enabling HTTPS Proxy Service

```bash
# Enable HTTPS proxy service
./sweb.exe -https-proxy

# Enable HTTPS proxy service and turn on authentication
./sweb.exe -https-proxy -https-proxy-auth

# Specify custom port
./sweb.exe -https-proxy -https-proxy-port 10443

# Custom authentication username and password
./sweb.exe -https-proxy -https-proxy-auth -https-proxy-username myuser -https-proxy-password mypass

# Use with other features
./sweb.exe -https-proxy -upload -files
```

### Client Configuration

#### Browser Configuration

1.  Configure HTTPS proxy in browser proxy settings.
2.  Proxy server: `localhost`
3.  Port: `10443` (or custom port)

#### Command Line Tool Configuration

```bash
# curl using HTTPS proxy (no authentication)
curl --proxy https://localhost:10443 https://example.com

# curl using HTTPS proxy (with authentication)
curl --proxy https://https:https@localhost:10443 https://example.com

# Set environment variables (no authentication)
export https_proxy=https://localhost:10443

# Set environment variables (with authentication)
export https_proxy=https://https:https@localhost:10443
```

### Feature Characteristics

  * **TLS Encryption**: The proxy connection itself uses TLS encryption.
  * **CONNECT Method**: Supports HTTPS tunnel proxy.
  * **Basic Authentication**: Supports username/password authentication to protect the proxy service.
  * **High Concurrency**: High-performance processing based on Go routines.
  * **Certificate Support**: Automatically loads SSL certificate files.
  * **Error Handling**: Comprehensive error handling and logging.

### Certificate Requirements

The HTTPS proxy requires SSL certificate files:

  * Certificate file: `./cert/server.crt`
  * Private key file: `./cert/server.key`

You can use the following command to generate a self-signed certificate:

```bash
openssl req -x509 -newkey rsa:4096 -keyout cert/server.key -out cert/server.crt -days 365 -nodes
```

## 🛠️ Technical Specifications

  * **Language**: Go language
  * **Dependencies**: Minimized external dependencies, primarily using the standard library.
  * **Protocols**: HTTP/1.1, WebDAV RFC 4918, SOCKS5 RFC 1928
  * **Proxy**: High-performance SOCKS5, HTTP, and HTTPS proxy server.
  * **Encoding**: UTF-8 support, perfect handling of Chinese characters.
  * **Platform**: Cross-platform compatibility (Windows, Linux, macOS).
  * **Architecture Support**:
      * Windows: 32-bit (386), 64-bit (amd64)
      * Linux: 32-bit (386), 64-bit (amd64)
      * macOS: Intel (amd64), Apple Silicon (arm64)
  * **Deployment**: Lightweight, single-file deployment.
  * **Concurrency**: High-concurrency processing based on Go routines.
  * **Build**: Supports cross-compilation, one build for multiple platforms.

## 🌍 Multilingual Support

sweb supports a bilingual interface (Chinese/English), providing a complete international experience:

### 🔧 Language Settings

  * **Automatic Detection**: Automatically selects the interface language based on the browser's language.
  * **Manual Switching**: A language switch button is provided in the upper right corner of the page.
  * **Local Storage**: The user's language preference is saved in the browser.

### 📖 Help Message Language

```bash
# Display Chinese help information (default)
./sweb.exe -help

# Display English help information
./sweb.exe -help-lang en -help
```

### 🌐 Web Interface Language

  * **Chinese Interface**: Automatically displayed when the browser language is Chinese.
  * **English Interface**: Automatically displayed for other browser languages.
  * **Language Switching**: Click the "中文"/"English" button in the upper right corner of the page to switch.
  * **Status Synchronization**: All feature status information supports bilingual display.

### 📱 Supported Languages

| Language | Code | Support Scope |
| :--- | :--- | :--- |
| Chinese | zh | Full support (Help information, Web interface) |
| English | en | Full support (Help information, Web interface) |

## 📊 Real-time Status Monitoring

The server provides a real-time status API and a Web interface:

  * **Status API**: `GET /api/status`
  * **Web Interface**: Home page automatically displays the current feature status.
  * **Automatic Update**: Checks for status changes every 30 seconds.

### API Response Example

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
    "port": 8081,
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

## 🔒 Security Considerations

### Default Security Policy

  * File upload feature is **disabled** by default.
  * File browsing feature is **disabled** by default.
  * WebDAV service is **disabled** by default.
  * WebDAV authentication is **disabled** by default.
  * HTTPS service is **disabled** by default.
  * SOCKS5 proxy service is **disabled** by default.
  * SOCKS5 proxy authentication is **disabled** by default.
  * HTTP proxy service is **disabled** by default.
  * HTTP proxy authentication is **disabled** by default.
  * HTTPS proxy service is **disabled** by default.
  * HTTPS proxy authentication is **disabled** by default.
  * Advanced features require explicit enabling via command-line arguments.

### Permission Control

  * WebDAV supports read-only mode and Basic authentication, with optimized Windows client compatibility.
  * The WebDAV access directory range can be limited.
  * HTTP proxy supports Basic authentication protection.
  * HTTPS proxy supports Basic authentication protection and TLS encryption.
  * SOCKS5 proxy supports Basic authentication protection.
  * It is recommended to enable authentication in production environments.

### Best Practices

```bash
# Recommended production configuration (Read-only WebDAV + HTTPS)
./sweb.exe -webdav -webdav-readonly -webdav-dir /safe/directory -https

# Development environment configuration (Enable all features)
./sweb.exe -upload -files -webdav -https -socks5 -proxy

# File sharing configuration (Only file browsing enabled)
./sweb.exe -files

# Proxy server configuration (Proxy features only)
./sweb.exe -socks5 -proxy

# Internal network proxy configuration (Specify ports)
./sweb.exe -socks5 -socks5-port 1080 -proxy -proxy-port 8080
```

## 📁 Directory Structure

```
sweb/
├── src/                        # Source code directory
│   ├── main.go                 # Main program entry file
│   ├── upload.go               # File upload feature module
│   ├── files.go                # File browsing feature module
│   ├── webdav.go               # WebDAV service module
│   ├── server.go               # HTTP/HTTPS server module
│   ├── utils.go                # Utility functions and page generation
│   ├── socks5.go               # SOCKS5 proxy server module
│   ├── proxy.go                # HTTP proxy server module
│   └── https_proxy.go          # HTTPS proxy server module
├── test/                       # Test code directory
│   ├── test_all.bat            # Windows complete test suite
│   ├── test_all.sh             # Linux/macOS complete test suite
│   ├── test_all.go             # Comprehensive feature testing
│   ├── test_socks5.go          # SOCKS5 proxy testing
│   ├── test_proxy.go           # HTTP proxy testing
│   ├── debug_socks5.go         # SOCKS5 debugging tool
│   └── README.md               # Testing documentation
├── doc/                        # Documentation directory
│   ├── 项目结构说明.md         # Detailed project structure
│   ├── HTTPS代理配置指南.md    # HTTPS proxy configuration guide
│   ├── HTTPS功能说明.md        # HTTPS feature documentation
│   ├── WebDAV使用说明.md       # Detailed WebDAV documentation
│   ├── 测试说明.md             # Testing documentation
│   ├── 浏览器配置说明.md       # Browser configuration instructions
│   └── 目录重构完成说明.md     # Directory refactoring notes
├── tools/                      # Tools directory
│   └── generate-cert.go        # SSL certificate generation tool
├── bin/                        # Compilation output directory
│   ├── windows/                # Windows executable files
│   ├── linux/                  # Linux executable files
│   └── macos/                  # macOS executable files
├── cert/                       # SSL Certificate directory (Required for HTTPS feature)
│   ├── server.crt              # SSL Certificate file
│   └── server.key              # SSL Private key file
├── web/                        # Web file directory (Automatically created)
│   └── index.html              # Default home page (Automatically generated)
├── go.mod                      # Go module file
├── go.sum                      # Dependency checksum file
├── Makefile                    # Build script (Linux/macOS)
├── build.bat                   # Build script (Windows)
└── README.md                   # Project description
```

## 🎯 Use Cases

  * **Development Testing**: Quickly set up a local file server.
  * **File Sharing**: Securely share files within a local network.
  * **Remote Management**: Remotely manage files via WebDAV.
  * **Static Websites**: Host simple static websites.
  * **File Backup**: Serve as a simple file upload service.
  * **Team Collaboration**: Quickly view and download files via the file browsing feature.
  * **Secure Transfer**: Use HTTPS to ensure secure file transfer.
  * **Network Proxy**: Provide SOCKS5 and HTTP proxy services.
  * **Intranet Penetration**: Provide proxy access in restricted network environments.
  * **Development Debugging**: Debug network requests through a proxy server.
  * **Traffic Forwarding**: High-performance TCP and HTTP traffic forwarding.

## 🤝 Contribution Guide

Welcome to submit Issues and Pull Requests\!

1.  Fork this project.
2.  Create a feature branch (`git checkout -b feature/AmazingFeature`).
3.  Commit your changes (`git commit -m 'Add some AmazingFeature'`).
4.  Push to the branch (`git push origin feature/AmazingFeature`).
5.  Open a Pull Request.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](https://www.google.com/search?q=LICENSE) file for details.

## 🙏 Acknowledgements

  * [golang.org/x/net/webdav](https://pkg.go.dev/golang.org/x/net/webdav) - WebDAV protocol implementation
  * Go Standard Library - Provides powerful HTTP server functionality

## 📸 Feature Screenshots

### Home Page Interface

  * Real-time display of all feature statuses (Upload, File Browsing, WebDAV, HTTPS, SOCKS5, HTTP Proxy).
  * Responsive design, supports mobile devices.
  * Chinese interface, simple operation.
  * Dynamic button status, displayed based on feature enablement.
  * Proxy service status and statistics display.

### File Browsing Page

  * Grid layout for files, supports file type icons.
  * Displays file size, modification time, and other details.
  * Clicking a file directly downloads it, supports real-time refresh.

### File Upload Page

  * Supports drag-and-drop multiple file uploads.
  * Displays real-time upload progress bar.
  * File list management, supports adding/removing files.

### Feature Status Display

  * 🔒 Disabled: Feature is not enabled.
  * ✅ Enabled: Feature is working normally.
  * 📖 Read-Only Mode: WebDAV read-only mode.
  * ⚠️ Certificate Error: HTTPS certificate issue.

## 🔧 Development Information

### Build Requirements

  * Go 1.19 or higher
  * Network connection (for downloading dependencies)

### Dependencies

```go
require (
    golang.org/x/net v0.x.x // WebDAV protocol support
)
```

### Compilation Commands

```bash
# Using Makefile (Recommended)
make all                    # Compile for all platforms (Windows, Linux, macOS)
make windows               # Only compile Windows versions (64-bit and 32-bit)
make linux                 # Only compile Linux versions (64-bit and 32-bit)
make macos                 # Only compile macOS versions (Intel and Apple Silicon)

# Separately compile specific architectures
make windows-64            # Windows 64-bit
make windows-32            # Windows 32-bit
make linux-64              # Linux 64-bit
make linux-32              # Linux 32-bit
make macos-intel           # macOS Intel
make macos-arm64           # macOS Apple Silicon

# Using build.bat (Windows)
build.bat                  # Compile all platform versions

# Manual compilation (including all source files)
go build -o sweb ./src/main.go ./src/upload.go ./src/files.go ./src/webdav.go ./src/server.go ./src/utils.go ./src/socks5.go ./src/proxy.go ./src/https_proxy.go

# Cross-compilation examples
# Windows 64-bit
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o bin/windows/sweb-windows-amd64.exe ./src/*.go

# Windows 32-bit
GOOS=windows GOARCH=386 go build -ldflags "-s -w" -o bin/windows/sweb-windows-386.exe ./src/*.go

# Linux 64-bit
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o bin/linux/sweb-linux-amd64 ./src/*.go

# Linux 32-bit
GOOS=linux GOARCH=386 go build -ldflags "-s -w" -o bin/linux/sweb-linux-386 ./src/*.go

# macOS Intel
GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o bin/macos/sweb-darwin-amd64 ./src/*.go

# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o bin/macos/sweb-darwin-arm64 ./src/*.go
```

### Running Tests

```bash
# Navigate to the test directory
cd test

# Windows - Run the complete test suite
test_all.bat

# Linux/macOS - Run the complete test suite
chmod +x test_all.sh
./test_all.sh

# Run tests individually
go run test_all.go          # Comprehensive feature testing
go run test_socks5.go       # SOCKS5 proxy testing
go run test_proxy.go        # HTTP proxy testing
go run debug_socks5.go      # SOCKS5 debugging tool
```

## 🐛 Troubleshooting

### Common Issues

#### 1\. Port is already in use

```bash
# Check port usage
netstat -an | grep :8080

# Use a different port
./sweb.exe -port 9000
```

#### 2\. WebDAV connection failed

  * Confirm the server is running and WebDAV is enabled.
  * Check firewall settings.
  * Confirm the correct URL format is used: `http://localhost:8081/webdav` (default port).

#### 3\. File upload failed

  * Confirm upload feature is enabled (using the `-upload` parameter).
  * Check write permissions for the target directory.
  * Confirm sufficient disk space.

#### 4\. File browsing page is inaccessible

  * Confirm file browsing feature is enabled (using the `-files` parameter).
  * Check if the `web` directory exists.
  * Confirm read permissions.

#### 5\. HTTPS connection failed

  * Confirm HTTPS feature is enabled (using the `-https` parameter).
  * Check for `server.crt` and `server.key` files in the `cert` directory.
  * Confirm the certificate file format is correct.

#### 6\. Chinese file names display abnormally

  * The server supports UTF-8 encoding.
  * Check the client's encoding settings.
  * Ensure the file system supports Unicode.

#### 7\. SOCKS5 proxy connection failed

  * Confirm SOCKS5 proxy feature is enabled (using the `-socks5` parameter).
  * Check if the port is in use: `netstat -an | grep :1080`
  * Confirm client SOCKS5 configuration is correct.
  * Check if the firewall is blocking connections.

#### 8\. HTTP proxy cannot access websites

  * Confirm HTTP proxy feature is enabled (using the `-proxy` parameter).
  * Check if the port is in use: `netstat -an | grep :10808`
  * Confirm browser proxy settings are correct.
  * Check if the target website is accessible.

### Log Debugging

The server outputs detailed operation logs, including:

  * Feature enablement/disablement status
  * WebDAV operation records
  * File upload status
  * HTTPS certificate status
  * SOCKS5 proxy connection and data forwarding
  * HTTP proxy request processing
  * Error messages and performance statistics

## 🔄 Version History

### v0.0.5 (Current Version)

  * ✅ WebDAV service Basic authentication support, optimized Windows client compatibility.
  * ✅ HTTP proxy service Basic authentication support.
  * ✅ Independent HTTPS proxy service implementation, supporting TLS encryption and Basic authentication.
  * ✅ SOCKS5 proxy service Basic authentication support.
  * ✅ HTTP service switch feature (allows only HTTPS or proxy services to be enabled).
  * ✅ WebDAV Windows mount authentication fix.
  * ✅ Authentication feature test suite.
  * ✅ Multilingual support (Chinese/English bilingual interface).
  * ✅ Updated documentation and help information.

### v0.0.4

  * ✅ Basic static file serving.
  * ✅ Drag-and-drop multiple file upload feature (with progress bar).
  * ✅ Dedicated file browsing page.
  * ✅ Complete WebDAV protocol support.
  * ✅ HTTPS encrypted connection support.
  * ✅ High-performance SOCKS5 proxy server.
  * ✅ HTTP/HTTPS proxy server.
  * ✅ Real-time status monitoring and statistics.
  * ✅ Responsive Web interface.
  * ✅ Command-line argument configuration.
  * ✅ Secure default policy.
  * ✅ Modular code structure.
  * ✅ Refactored directory structure (src/, test/, doc/, tools/).
  * ✅ Multi-platform compilation support (Windows/Linux/macOS, 32-bit and 64-bit).
  * ✅ Complete test suite and documentation.

### v0.0.3

  * ✅ Basic file service features.
  * ✅ WebDAV and HTTPS support.
  * ✅ Web interface and status monitoring.

### v0.0.2

  * ✅ File upload and browsing features.
  * ✅ Proxy server support.

### v0.0.1

  * ✅ Initial version.
  * ✅ Basic HTTP server.

### Planned Features

  * [ ] User authentication and permission management system
  * [x] \~\~Proxy service authentication mechanism\~\~ (Implemented in v0.0.5)
  * [ ] File preview feature
  * [ ] Batch file operations
  * [ ] Configuration file support
  * [ ] File search feature
  * [ ] File version management
  * [ ] Proxy service access control and whitelisting

## 🔌 Default Port Information

| Service Type | Default Port | Description |
| :--- | :--- | :--- |
| HTTP Service | 8080 | Primary Web service port |
| **WebDAV Service** | **8081** | **Independent WebDAV service port** |
| HTTPS Service | 8443 | Encrypted Web service port |
| SOCKS5 Proxy | 1080 | SOCKS5 proxy service port |
| HTTP Proxy | 10808 | HTTP proxy service port |
| HTTPS Proxy | 10443 | HTTPS proxy service port |

### Important Notes

  * **WebDAV service now runs on the independent port 8081**, no longer sharing port 8080 with the HTTP service.
  * All ports can be customized using the corresponding command-line arguments.
  * The independent WebDAV port design improves service stability and security.
  * To access the WebDAV service, please use: `http://localhost:8081/webdav`

## 📞 Contact Information

If you have any questions or suggestions, please contact us through:

  * Submitting an [可疑链接已删除]
  * Starting a [可疑链接已删除]

## ⭐ Support the Project

If this project has been helpful to you, please consider giving it a star ⭐

-----

**Simple Web File Server** - Making file management simple and efficient\! 🚀

-----

Would you like me to translate any other sections or provide a summary of the key features?