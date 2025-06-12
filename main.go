package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/webdav"
)

// 全局变量存储功能状态
var (
	uploadEnabled  bool
	webdavEnabled  bool
	webdavDir      string
	webdavReadonly bool
)

func main() {
	// 解析命令行参数
	var port int
	var showHelp bool

	flag.BoolVar(&uploadEnabled, "upload", false, "启用文件上传功能")
	flag.BoolVar(&uploadEnabled, "enable-upload", false, "启用文件上传功能")
	flag.BoolVar(&webdavEnabled, "webdav", false, "启用WebDAV服务")
	flag.BoolVar(&webdavEnabled, "enable-webdav", false, "启用WebDAV服务")
	flag.StringVar(&webdavDir, "webdav-dir", ".", "WebDAV服务的根目录")
	flag.BoolVar(&webdavReadonly, "webdav-readonly", false, "WebDAV服务只读模式")
	flag.IntVar(&port, "port", 8080, "指定服务器端口")
	flag.IntVar(&port, "p", 8080, "指定服务器端口")
	flag.BoolVar(&showHelp, "help", false, "显示帮助信息")
	flag.BoolVar(&showHelp, "h", false, "显示帮助信息")

	flag.Parse()

	// 显示帮助信息
	if showHelp {
		showHelpInfo()
		return
	}

	// 创建web目录（如果不存在）
	webDir := "./web"
	if _, err := os.Stat(webDir); os.IsNotExist(err) {
		err := os.Mkdir(webDir, 0755)
		if err != nil {
			log.Fatalf("无法创建web目录: %v", err)
		}
	}

	// 检查并创建默认页面
	createDefaultPageIfNeeded(webDir, uploadEnabled)

	// 处理静态文件（HTML, JS等）
	fileServer := http.FileServer(http.Dir(webDir))
	http.Handle("/", fileServer)

	// 添加上传状态API端点
	http.HandleFunc("/api/upload-status", uploadStatusHandler)

	// 根据参数决定是否启用文件上传
	if uploadEnabled {
		http.HandleFunc("/upload", uploadHandler)
		fmt.Println("✅ 文件上传功能已启用")
	} else {
		http.HandleFunc("/upload", uploadDisabledHandler)
		fmt.Println("🔒 文件上传功能已禁用 (使用 -upload 参数启用)")
	}

	// 根据参数决定是否启用WebDAV服务
	if webdavEnabled {
		setupWebDAVHandler()
		if webdavReadonly {
			fmt.Printf("✅ WebDAV服务已启用 (只读模式) - 目录: %s\n", webdavDir)
		} else {
			fmt.Printf("✅ WebDAV服务已启用 (读写模式) - 目录: %s\n", webdavDir)
		}
	} else {
		http.HandleFunc("/webdav", webdavDisabledHandler)
		fmt.Println("🔒 WebDAV服务已禁用 (使用 -webdav 参数启用)")
	}

	// 启动服务器
	fmt.Printf("服务器启动在 http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// 显示上传表单
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(`
            <!DOCTYPE html>
            <html>
            <head>
                <title>文件上传</title>
            </head>
            <body>
                <h2>文件上传</h2>
                <form method="post" enctype="multipart/form-data">
                    <input type="file" name="file">
                    <input type="submit" value="上传">
                </form>
            </body>
            </html>
        `))
	} else if r.Method == "POST" {
		// 处理文件上传
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "无法获取上传文件: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		// 创建目标文件
		dst, err := os.Create(filepath.Join("web", header.Filename))
		if err != nil {
			http.Error(w, "无法创建目标文件: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// 复制文件内容
		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "无法保存文件: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// 返回成功信息
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(fmt.Sprintf(`
            <!DOCTYPE html>
            <html>
            <head>
                <title>上传成功</title>
            </head>
            <body>
                <h2>文件上传成功!</h2>
                <p>文件名: %s</p>
                <p><a href="/%s">查看文件</a></p>
                <p><a href="/upload">继续上传</a></p>
            </body>
            </html>
        `, header.Filename, header.Filename)))
	} else {
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
	}
}

// createDefaultPageIfNeeded 检查并创建默认页面
func createDefaultPageIfNeeded(webDir string, uploadEnabled bool) {
	// 检查是否存在默认页面
	indexFiles := []string{"index.html", "index.htm"}
	hasDefaultPage := false

	for _, indexFile := range indexFiles {
		indexPath := filepath.Join(webDir, indexFile)
		if _, err := os.Stat(indexPath); err == nil {
			hasDefaultPage = true
			break
		}
	}

	// 如果没有默认页面，创建一个
	if !hasDefaultPage {
		indexPath := filepath.Join(webDir, "index.html")
		indexContent := generateEnhancedDefaultPageContent()

		err := os.WriteFile(indexPath, []byte(indexContent), 0644)
		if err != nil {
			log.Printf("警告：无法创建默认页面: %v", err)
		} else {
			fmt.Println("已创建默认页面: index.html (支持上传和WebDAV状态检查)")
		}
	}
}

// showHelpInfo 显示帮助信息
func showHelpInfo() {
	fmt.Println("简单Web文件服务器 - 基于Go语言开发")
	fmt.Println()
	fmt.Println("用法:")
	fmt.Println("  sweb.exe [选项]")
	fmt.Println()
	fmt.Println("选项:")
	fmt.Println("  -upload, --enable-upload    启用文件上传功能 (默认: 禁用)")
	fmt.Println("  -webdav, --enable-webdav    启用WebDAV服务 (默认: 禁用)")
	fmt.Println("  -webdav-dir <目录>          WebDAV服务的根目录 (默认: 当前目录)")
	fmt.Println("  -webdav-readonly            WebDAV服务只读模式 (默认: 读写)")
	fmt.Println("  -port, -p <端口>           指定服务器端口 (默认: 8080)")
	fmt.Println("  -help, -h                  显示此帮助信息")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  sweb.exe                           # 启动服务器，仅提供静态文件服务")
	fmt.Println("  sweb.exe -upload                   # 启动服务器并启用文件上传功能")
	fmt.Println("  sweb.exe -webdav                   # 启动服务器并启用WebDAV服务")
	fmt.Println("  sweb.exe -webdav -webdav-readonly  # 启动只读WebDAV服务")
	fmt.Println("  sweb.exe -webdav -webdav-dir /data # 指定WebDAV目录")
	fmt.Println("  sweb.exe -upload -webdav -p 9000   # 启用所有功能并指定端口")
	fmt.Println()
	fmt.Println("WebDAV访问:")
	fmt.Println("  WebDAV地址: http://localhost:8080/webdav")
	fmt.Println("  可以使用支持WebDAV的客户端连接，如Windows资源管理器、")
	fmt.Println("  macOS Finder、或专用的WebDAV客户端软件。")
	fmt.Println()
	fmt.Println("安全说明:")
	fmt.Println("  文件上传和WebDAV功能默认禁用以确保服务器安全。")
	fmt.Println("  只有在明确需要时才使用相应参数启用。")
}

// uploadDisabledHandler 处理上传功能被禁用时的请求
func uploadDisabledHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte(`
        <!DOCTYPE html>
        <html lang="zh-CN">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>上传功能已禁用</title>
            <style>
                body {
                    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
                    max-width: 600px;
                    margin: 50px auto;
                    padding: 20px;
                    text-align: center;
                    background-color: #f8f9fa;
                }
                .container {
                    background: white;
                    padding: 40px;
                    border-radius: 10px;
                    box-shadow: 0 2px 10px rgba(0,0,0,0.1);
                    border-left: 5px solid #dc3545;
                }
                h1 {
                    color: #dc3545;
                    margin-bottom: 20px;
                }
                .icon {
                    font-size: 64px;
                    margin-bottom: 20px;
                }
                .command {
                    background: #f8f9fa;
                    padding: 10px;
                    border-radius: 5px;
                    font-family: 'Courier New', monospace;
                    margin: 10px 0;
                    border: 1px solid #dee2e6;
                }
                .back-link {
                    display: inline-block;
                    background: #007acc;
                    color: white;
                    padding: 10px 20px;
                    text-decoration: none;
                    border-radius: 5px;
                    margin-top: 20px;
                }
                .back-link:hover {
                    background: #005a9e;
                }
            </style>
        </head>
        <body>
            <div class="container">
                <div class="icon">🔒</div>
                <h1>文件上传功能已禁用</h1>
                <p>出于安全考虑，文件上传功能默认处于禁用状态。</p>
                <p>如需启用文件上传功能，请使用以下命令重新启动服务器：</p>

                <div class="command">sweb.exe -upload</div>
                <p>或</p>
                <div class="command">sweb.exe --enable-upload</div>

                <p>您也可以使用 <code>sweb.exe -help</code> 查看所有可用选项。</p>

                <a href="/" class="back-link">← 返回首页</a>
            </div>
        </body>
        </html>
    `))
}

// uploadStatusHandler 处理上传状态查询请求
func uploadStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	response := map[string]interface{}{
		"upload": map[string]interface{}{
			"enabled": uploadEnabled,
			"status": func() string {
				if uploadEnabled {
					return "enabled"
				}
				return "disabled"
			}(),
		},
		"webdav": map[string]interface{}{
			"enabled":   webdavEnabled,
			"readonly":  webdavReadonly,
			"directory": webdavDir,
			"status": func() string {
				if webdavEnabled {
					if webdavReadonly {
						return "enabled-readonly"
					}
					return "enabled-readwrite"
				}
				return "disabled"
			}(),
		},
	}

	json.NewEncoder(w).Encode(response)
}

// setupWebDAVHandler 设置WebDAV处理器
func setupWebDAVHandler() {
	// 确保WebDAV目录存在
	if _, err := os.Stat(webdavDir); os.IsNotExist(err) {
		err := os.MkdirAll(webdavDir, 0755)
		if err != nil {
			log.Fatalf("无法创建WebDAV目录: %v", err)
		}
	}

	// 创建WebDAV处理器
	handler := &webdav.Handler{
		Prefix:     "/webdav",
		FileSystem: webdav.Dir(webdavDir),
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			if err != nil {
				// 过滤掉一些常见的非关键错误
				errStr := err.Error()
				// 忽略文件不存在的PROPFIND错误（这在文件创建过程中是正常的）
				if r.Method == "PROPFIND" && (strings.Contains(errStr, "cannot find the file specified") ||
					strings.Contains(errStr, "no such file or directory") ||
					strings.Contains(errStr, "file does not exist")) {
					// 这些是正常的操作流程，不记录错误
					return
				}
				// 记录其他重要错误
				log.Printf("WebDAV操作: %s %s - %v", r.Method, r.URL.Path, err)
			}
		},
	}

	// 如果是只读模式，包装处理器
	if webdavReadonly {
		http.HandleFunc("/webdav/", func(w http.ResponseWriter, r *http.Request) {
			// 只允许GET、HEAD、OPTIONS、PROPFIND方法
			switch r.Method {
			case "GET", "HEAD", "OPTIONS", "PROPFIND":
				handler.ServeHTTP(w, r)
			default:
				http.Error(w, "WebDAV服务处于只读模式", http.StatusMethodNotAllowed)
			}
		})
		// 处理根路径
		http.HandleFunc("/webdav", func(w http.ResponseWriter, r *http.Request) {
			// 只允许GET、HEAD、OPTIONS、PROPFIND方法
			switch r.Method {
			case "GET", "HEAD", "OPTIONS", "PROPFIND":
				handler.ServeHTTP(w, r)
			default:
				http.Error(w, "WebDAV服务处于只读模式", http.StatusMethodNotAllowed)
			}
		})
	} else {
		http.Handle("/webdav/", handler)
		http.Handle("/webdav", handler)
	}
}

// webdavDisabledHandler 处理WebDAV功能被禁用时的请求
func webdavDisabledHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte(`
        <!DOCTYPE html>
        <html lang="zh-CN">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>WebDAV服务已禁用</title>
            <style>
                body {
                    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
                    max-width: 600px;
                    margin: 50px auto;
                    padding: 20px;
                    text-align: center;
                    background-color: #f8f9fa;
                }
                .container {
                    background: white;
                    padding: 40px;
                    border-radius: 10px;
                    box-shadow: 0 2px 10px rgba(0,0,0,0.1);
                    border-left: 5px solid #dc3545;
                }
                h1 {
                    color: #dc3545;
                    margin-bottom: 20px;
                }
                .icon {
                    font-size: 64px;
                    margin-bottom: 20px;
                }
                .command {
                    background: #f8f9fa;
                    padding: 10px;
                    border-radius: 5px;
                    font-family: 'Courier New', monospace;
                    margin: 10px 0;
                    border: 1px solid #dee2e6;
                }
                .back-link {
                    display: inline-block;
                    background: #007acc;
                    color: white;
                    padding: 10px 20px;
                    text-decoration: none;
                    border-radius: 5px;
                    margin-top: 20px;
                }
                .back-link:hover {
                    background: #005a9e;
                }
            </style>
        </head>
        <body>
            <div class="container">
                <div class="icon">🔒</div>
                <h1>WebDAV服务已禁用</h1>
                <p>出于安全考虑，WebDAV服务默认处于禁用状态。</p>
                <p>如需启用WebDAV服务，请使用以下命令重新启动服务器：</p>

                <div class="command">sweb.exe -webdav</div>
                <p>或</p>
                <div class="command">sweb.exe --enable-webdav</div>

                <p><strong>可选参数：</strong></p>
                <div class="command">sweb.exe -webdav -webdav-dir /path/to/directory</div>
                <div class="command">sweb.exe -webdav -webdav-readonly</div>

                <p>您也可以使用 <code>sweb.exe -help</code> 查看所有可用选项。</p>

                <a href="/" class="back-link">← 返回首页</a>
            </div>
        </body>
        </html>
    `))
}

// generateDynamicDefaultPageContent 生成支持动态状态检查的默认页面HTML内容
func generateDynamicDefaultPageContent() string {
	return `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>简单Web文件服务器</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            line-height: 1.6;
            background-color: #f5f5f5;
        }
        .container {
            background: white;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        h1 {
            color: #333;
            text-align: center;
            border-bottom: 3px solid #007acc;
            padding-bottom: 10px;
        }
        h2 {
            color: #555;
            margin-top: 30px;
        }
        .feature {
            background: #f8f9fa;
            padding: 15px;
            margin: 10px 0;
            border-left: 4px solid #007acc;
            border-radius: 5px;
        }
        .button {
            display: inline-block;
            background: #007acc;
            color: white;
            padding: 10px 20px;
            text-decoration: none;
            border-radius: 5px;
            margin: 10px 5px;
            transition: background 0.3s;
        }
        .button:hover {
            background: #005a9e;
        }
        .button.disabled {
            background: #6c757d;
            cursor: not-allowed;
        }
        .footer {
            text-align: center;
            margin-top: 30px;
            padding-top: 20px;
            border-top: 1px solid #eee;
            color: #666;
        }
        code {
            background: #f4f4f4;
            padding: 2px 5px;
            border-radius: 3px;
            font-family: 'Courier New', monospace;
        }
        .status-indicator {
            font-weight: bold;
            padding: 2px 6px;
            border-radius: 3px;
        }
        .status-enabled {
            color: #28a745;
        }
        .status-disabled {
            color: #dc3545;
        }
        .loading {
            color: #6c757d;
        }
        .hidden {
            display: none;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🌐 简单Web文件服务器</h1>

        <p>欢迎使用这个简单而实用的Web文件服务器！这是一个用Go语言编写的轻量级文件管理工具。</p>

        <h2>📋 项目功能</h2>

        <div class="feature">
            <strong>📁 静态文件服务</strong><br>
            自动服务web目录下的所有文件，支持HTML、CSS、JavaScript、图片等各种文件类型。
        </div>

        <div class="feature">
            <strong>📤 文件上传</strong><br>
            <span id="upload-feature-description">文件上传功能可通过命令行参数启用，确保服务器安全。</span>
            <span id="upload-status" class="status-indicator loading">🔄 检查中...</span>
        </div>

        <div class="feature">
            <strong>📂 目录浏览</strong><br>
            当没有默认页面时，自动显示目录内容，方便浏览和下载文件。
        </div>

        <div class="feature">
            <strong>🔧 自动配置</strong><br>
            自动创建必要的目录结构，无需手动配置即可使用。
        </div>

        <h2>🚀 快速开始</h2>

        <div id="upload-section">
            <p><strong>文件上传：</strong></p>
            <div id="upload-enabled-content" class="hidden">
                <a href="/upload" class="button" id="upload-button">📤 上传文件</a>
            </div>
            <div id="upload-disabled-content" class="hidden">
                <p>要启用文件上传功能，请使用以下命令启动服务器：</p>
                <code>sweb.exe -upload</code> 或 <code>sweb.exe --enable-upload</code>
                <br><br>
                <a href="/upload" class="button disabled" id="upload-button-disabled">📤 上传功能已禁用</a>
            </div>
        </div>

        <p><strong>服务器信息：</strong></p>
        <ul>
            <li>服务端口: <code>8080</code></li>
            <li>文件目录: <code>./web</code></li>
            <li>上传地址: <code>/upload</code></li>
        </ul>

        <h2>💡 使用说明</h2>
        <div id="usage-instructions">
            <div id="usage-enabled" class="hidden">
                <ol>
                    <li><strong>上传文件</strong>：点击上方的"上传文件"按钮，选择要上传的文件</li>
                    <li><strong>访问文件</strong>：上传成功后，文件将保存在web目录下，可以直接通过URL访问</li>
                    <li><strong>管理文件</strong>：所有上传的文件都会显示在主页的文件列表中</li>
                </ol>
            </div>
            <div id="usage-disabled" class="hidden">
                <ol>
                    <li><strong>启用上传</strong>：使用 <code>-upload</code> 参数启动服务器以启用文件上传功能</li>
                    <li><strong>浏览文件</strong>：当前可以浏览和下载web目录中的现有文件</li>
                    <li><strong>安全考虑</strong>：文件上传功能默认禁用，确保服务器安全</li>
                </ol>
            </div>
        </div>

        <h2>🛠️ 技术特性</h2>

        <ul>
            <li>使用Go语言标准库开发，无外部依赖</li>
            <li>支持多部分表单数据上传</li>
            <li>自动MIME类型检测</li>
            <li>UTF-8编码支持，完美处理中文</li>
            <li>跨平台兼容（Windows、Linux、macOS）</li>
        </ul>

        <div class="footer">
            <p>🔗 <strong>简单Web文件服务器</strong> | 基于Go语言开发</p>
            <div id="footer-content">
                <div id="footer-enabled" class="hidden">
                    <p>开始使用：<a href="/upload" class="button">上传第一个文件</a></p>
                </div>
                <div id="footer-disabled" class="hidden">
                    <p>安全模式：文件上传功能已禁用</p>
                    <p>使用 <code>sweb.exe -help</code> 查看所有可用选项</p>
                </div>
            </div>
        </div>
    </div>

    <script>
        // 检查上传功能状态
        function checkUploadStatus() {
            fetch('/api/upload-status')
                .then(response => response.json())
                .then(data => {
                    updateUploadStatus(data.enabled);
                })
                .catch(error => {
                    console.error('检查上传状态失败:', error);
                    // 如果API调用失败，显示默认的禁用状态
                    updateUploadStatus(false);
                });
        }

        // 更新页面上的上传状态显示
        function updateUploadStatus(enabled) {
            const statusElement = document.getElementById('upload-status');
            const descriptionElement = document.getElementById('upload-feature-description');
            const uploadEnabledContent = document.getElementById('upload-enabled-content');
            const uploadDisabledContent = document.getElementById('upload-disabled-content');
            const usageEnabled = document.getElementById('usage-enabled');
            const usageDisabled = document.getElementById('usage-disabled');
            const footerEnabled = document.getElementById('footer-enabled');
            const footerDisabled = document.getElementById('footer-disabled');

            if (enabled) {
                // 上传功能已启用
                statusElement.textContent = '✅ 已启用';
                statusElement.className = 'status-indicator status-enabled';
                descriptionElement.textContent = '通过简单的Web界面上传文件到服务器，支持各种文件格式。';

                uploadEnabledContent.classList.remove('hidden');
                uploadDisabledContent.classList.add('hidden');
                usageEnabled.classList.remove('hidden');
                usageDisabled.classList.add('hidden');
                footerEnabled.classList.remove('hidden');
                footerDisabled.classList.add('hidden');
            } else {
                // 上传功能已禁用
                statusElement.textContent = '🔒 已禁用';
                statusElement.className = 'status-indicator status-disabled';
                descriptionElement.textContent = '文件上传功能可通过命令行参数启用，确保服务器安全。';

                uploadEnabledContent.classList.add('hidden');
                uploadDisabledContent.classList.remove('hidden');
                usageEnabled.classList.add('hidden');
                usageDisabled.classList.remove('hidden');
                footerEnabled.classList.add('hidden');
                footerDisabled.classList.remove('hidden');
            }
        }

        // 页面加载时检查状态
        document.addEventListener('DOMContentLoaded', function() {
            checkUploadStatus();

            // 每30秒检查一次状态（可选）
            setInterval(checkUploadStatus, 30000);
        });
    </script>
</body>
</html>`
}

// generateEnhancedDefaultPageContent 生成包含WebDAV功能的增强版默认页面
func generateEnhancedDefaultPageContent() string {
	return `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>简单Web文件服务器</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            line-height: 1.6;
            background-color: #f5f5f5;
        }
        .container {
            background: white;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        h1 {
            color: #333;
            text-align: center;
            border-bottom: 3px solid #007acc;
            padding-bottom: 10px;
        }
        h2 {
            color: #555;
            margin-top: 30px;
        }
        .feature {
            background: #f8f9fa;
            padding: 15px;
            margin: 10px 0;
            border-left: 4px solid #007acc;
            border-radius: 5px;
        }
        .button {
            display: inline-block;
            background: #007acc;
            color: white;
            padding: 10px 20px;
            text-decoration: none;
            border-radius: 5px;
            margin: 10px 5px;
            transition: background 0.3s;
        }
        .button:hover {
            background: #005a9e;
        }
        .button.disabled {
            background: #6c757d;
            cursor: not-allowed;
        }
        .footer {
            text-align: center;
            margin-top: 30px;
            padding-top: 20px;
            border-top: 1px solid #eee;
            color: #666;
        }
        code {
            background: #f4f4f4;
            padding: 2px 5px;
            border-radius: 3px;
            font-family: 'Courier New', monospace;
        }
        .status-indicator {
            font-weight: bold;
            padding: 2px 6px;
            border-radius: 3px;
        }
        .status-enabled {
            color: #28a745;
        }
        .status-disabled {
            color: #dc3545;
        }
        .status-readonly {
            color: #ffc107;
        }
        .loading {
            color: #6c757d;
        }
        .hidden {
            display: none;
        }
        .webdav-info {
            background: #e7f3ff;
            border-left: 4px solid #007acc;
            padding: 10px;
            margin: 10px 0;
            border-radius: 5px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🌐 简单Web文件服务器</h1>

        <p>欢迎使用这个简单而实用的Web文件服务器！这是一个用Go语言编写的轻量级文件管理工具。</p>

        <h2>📋 项目功能</h2>

        <div class="feature">
            <strong>📁 静态文件服务</strong><br>
            自动服务web目录下的所有文件，支持HTML、CSS、JavaScript、图片等各种文件类型。
        </div>

        <div class="feature">
            <strong>📤 文件上传</strong><br>
            <span id="upload-feature-description">文件上传功能可通过命令行参数启用，确保服务器安全。</span>
            <span id="upload-status" class="status-indicator loading">🔄 检查中...</span>
        </div>

        <div class="feature">
            <strong>🌐 WebDAV服务</strong><br>
            <span id="webdav-feature-description">WebDAV服务可通过命令行参数启用，支持文件管理客户端连接。</span>
            <span id="webdav-status" class="status-indicator loading">🔄 检查中...</span>
            <div id="webdav-info" class="webdav-info hidden">
                <strong>WebDAV地址:</strong> <code id="webdav-url">http://localhost:8080/webdav</code><br>
                <strong>挂载目录:</strong> <code id="webdav-directory">.</code><br>
                <strong>访问模式:</strong> <span id="webdav-mode">读写</span>
            </div>
        </div>

        <div class="feature">
            <strong>📂 目录浏览</strong><br>
            当没有默认页面时，自动显示目录内容，方便浏览和下载文件。
        </div>

        <div class="feature">
            <strong>🔧 自动配置</strong><br>
            自动创建必要的目录结构，无需手动配置即可使用。
        </div>

        <h2>🚀 快速开始</h2>

        <div id="upload-section">
            <p><strong>文件上传：</strong></p>
            <div id="upload-enabled-content" class="hidden">
                <a href="/upload" class="button" id="upload-button">📤 上传文件</a>
            </div>
            <div id="upload-disabled-content" class="hidden">
                <p>要启用文件上传功能，请使用以下命令启动服务器：</p>
                <code>sweb.exe -upload</code> 或 <code>sweb.exe --enable-upload</code>
                <br><br>
                <a href="/upload" class="button disabled" id="upload-button-disabled">📤 上传功能已禁用</a>
            </div>
        </div>

        <div id="webdav-section">
            <p><strong>WebDAV服务：</strong></p>
            <div id="webdav-enabled-content" class="hidden">
                <a href="/webdav" class="button" id="webdav-button">🌐 访问WebDAV</a>
                <p><small>可以在文件管理器中添加网络位置：<code id="webdav-mount-url">http://localhost:8080/webdav</code></small></p>
            </div>
            <div id="webdav-disabled-content" class="hidden">
                <p>要启用WebDAV服务，请使用以下命令启动服务器：</p>
                <code>sweb.exe -webdav</code> 或 <code>sweb.exe --enable-webdav</code>
                <br><br>
                <a href="/webdav" class="button disabled" id="webdav-button-disabled">🌐 WebDAV服务已禁用</a>
            </div>
        </div>

        <p><strong>服务器信息：</strong></p>
        <ul>
            <li>服务端口: <code>8080</code></li>
            <li>文件目录: <code>./web</code></li>
            <li>上传地址: <code>/upload</code></li>
            <li>WebDAV地址: <code>/webdav</code></li>
        </ul>

        <h2>💡 使用说明</h2>
        <div id="usage-instructions">
            <div id="usage-upload-enabled" class="hidden">
                <h3>文件上传</h3>
                <ol>
                    <li><strong>上传文件</strong>：点击上方的"上传文件"按钮，选择要上传的文件</li>
                    <li><strong>访问文件</strong>：上传成功后，文件将保存在web目录下，可以直接通过URL访问</li>
                    <li><strong>管理文件</strong>：所有上传的文件都会显示在主页的文件列表中</li>
                </ol>
            </div>
            <div id="usage-webdav-enabled" class="hidden">
                <h3>WebDAV服务</h3>
                <ol>
                    <li><strong>Windows</strong>：在文件资源管理器中，右键"此电脑" → "映射网络驱动器" → 输入WebDAV地址</li>
                    <li><strong>macOS</strong>：在Finder中，按Cmd+K → 输入WebDAV地址</li>
                    <li><strong>Linux</strong>：使用davfs2或其他WebDAV客户端挂载</li>
                    <li><strong>移动设备</strong>：使用支持WebDAV的文件管理应用</li>
                </ol>
            </div>
            <div id="usage-disabled" class="hidden">
                <ol>
                    <li><strong>启用功能</strong>：使用相应的命令行参数启动服务器</li>
                    <li><strong>浏览文件</strong>：当前可以浏览和下载web目录中的现有文件</li>
                    <li><strong>安全考虑</strong>：高级功能默认禁用，确保服务器安全</li>
                </ol>
            </div>
        </div>

        <h2>🛠️ 技术特性</h2>

        <ul>
            <li>使用Go语言标准库开发，轻量级无外部依赖</li>
            <li>支持多部分表单数据上传</li>
            <li>完整的WebDAV协议支持（RFC 4918）</li>
            <li>可配置的读写权限控制</li>
            <li>自动MIME类型检测</li>
            <li>UTF-8编码支持，完美处理中文</li>
            <li>跨平台兼容（Windows、Linux、macOS）</li>
        </ul>

        <div class="footer">
            <p>🔗 <strong>简单Web文件服务器</strong> | 基于Go语言开发</p>
            <div id="footer-content">
                <div id="footer-enabled" class="hidden">
                    <p>开始使用：
                        <span id="footer-upload-link" class="hidden"><a href="/upload" class="button">上传文件</a></span>
                        <span id="footer-webdav-link" class="hidden"><a href="/webdav" class="button">访问WebDAV</a></span>
                    </p>
                </div>
                <div id="footer-disabled" class="hidden">
                    <p>安全模式：高级功能已禁用</p>
                    <p>使用 <code>sweb.exe -help</code> 查看所有可用选项</p>
                </div>
            </div>
        </div>
    </div>

    <script>
        // 检查服务状态
        function checkServiceStatus() {
            fetch('/api/upload-status')
                .then(response => response.json())
                .then(data => {
                    updateUploadStatus(data.upload);
                    updateWebDAVStatus(data.webdav);
                    updateUsageInstructions(data.upload, data.webdav);
                })
                .catch(error => {
                    console.error('检查服务状态失败:', error);
                    // 如果API调用失败，显示默认的禁用状态
                    updateUploadStatus({enabled: false, status: 'disabled'});
                    updateWebDAVStatus({enabled: false, status: 'disabled'});
                    updateUsageInstructions({enabled: false}, {enabled: false});
                });
        }

        // 更新页面上的上传状态显示
        function updateUploadStatus(uploadData) {
            const statusElement = document.getElementById('upload-status');
            const descriptionElement = document.getElementById('upload-feature-description');
            const uploadEnabledContent = document.getElementById('upload-enabled-content');
            const uploadDisabledContent = document.getElementById('upload-disabled-content');

            if (uploadData.enabled) {
                statusElement.textContent = '✅ 已启用';
                statusElement.className = 'status-indicator status-enabled';
                descriptionElement.textContent = '通过简单的Web界面上传文件到服务器，支持各种文件格式。';

                uploadEnabledContent.classList.remove('hidden');
                uploadDisabledContent.classList.add('hidden');
            } else {
                statusElement.textContent = '🔒 已禁用';
                statusElement.className = 'status-indicator status-disabled';
                descriptionElement.textContent = '文件上传功能可通过命令行参数启用，确保服务器安全。';

                uploadEnabledContent.classList.add('hidden');
                uploadDisabledContent.classList.remove('hidden');
            }
        }

        // 更新页面上的WebDAV状态显示
        function updateWebDAVStatus(webdavData) {
            const statusElement = document.getElementById('webdav-status');
            const descriptionElement = document.getElementById('webdav-feature-description');
            const webdavEnabledContent = document.getElementById('webdav-enabled-content');
            const webdavDisabledContent = document.getElementById('webdav-disabled-content');
            const webdavInfo = document.getElementById('webdav-info');
            const webdavDirectory = document.getElementById('webdav-directory');
            const webdavMode = document.getElementById('webdav-mode');

            if (webdavData.enabled) {
                if (webdavData.readonly) {
                    statusElement.textContent = '📖 只读模式';
                    statusElement.className = 'status-indicator status-readonly';
                    descriptionElement.textContent = 'WebDAV服务已启用（只读模式），支持文件浏览和下载。';
                    webdavMode.textContent = '只读';
                } else {
                    statusElement.textContent = '✅ 读写模式';
                    statusElement.className = 'status-indicator status-enabled';
                    descriptionElement.textContent = 'WebDAV服务已启用（读写模式），支持完整的文件管理操作。';
                    webdavMode.textContent = '读写';
                }

                webdavDirectory.textContent = webdavData.directory || '.';
                webdavInfo.classList.remove('hidden');
                webdavEnabledContent.classList.remove('hidden');
                webdavDisabledContent.classList.add('hidden');
            } else {
                statusElement.textContent = '🔒 已禁用';
                statusElement.className = 'status-indicator status-disabled';
                descriptionElement.textContent = 'WebDAV服务可通过命令行参数启用，支持文件管理客户端连接。';

                webdavInfo.classList.add('hidden');
                webdavEnabledContent.classList.add('hidden');
                webdavDisabledContent.classList.remove('hidden');
            }
        }

        // 更新使用说明和页脚
        function updateUsageInstructions(uploadData, webdavData) {
            const usageUploadEnabled = document.getElementById('usage-upload-enabled');
            const usageWebdavEnabled = document.getElementById('usage-webdav-enabled');
            const usageDisabled = document.getElementById('usage-disabled');
            const footerEnabled = document.getElementById('footer-enabled');
            const footerDisabled = document.getElementById('footer-disabled');
            const footerUploadLink = document.getElementById('footer-upload-link');
            const footerWebdavLink = document.getElementById('footer-webdav-link');

            const anyEnabled = uploadData.enabled || webdavData.enabled;

            if (anyEnabled) {
                footerEnabled.classList.remove('hidden');
                footerDisabled.classList.add('hidden');
                usageDisabled.classList.add('hidden');

                if (uploadData.enabled) {
                    usageUploadEnabled.classList.remove('hidden');
                    footerUploadLink.classList.remove('hidden');
                } else {
                    usageUploadEnabled.classList.add('hidden');
                    footerUploadLink.classList.add('hidden');
                }

                if (webdavData.enabled) {
                    usageWebdavEnabled.classList.remove('hidden');
                    footerWebdavLink.classList.remove('hidden');
                } else {
                    usageWebdavEnabled.classList.add('hidden');
                    footerWebdavLink.classList.add('hidden');
                }
            } else {
                footerEnabled.classList.add('hidden');
                footerDisabled.classList.remove('hidden');
                usageDisabled.classList.remove('hidden');
                usageUploadEnabled.classList.add('hidden');
                usageWebdavEnabled.classList.add('hidden');
            }
        }

        // 页面加载时检查状态
        document.addEventListener('DOMContentLoaded', function() {
            checkServiceStatus();

            // 每30秒检查一次状态
            setInterval(checkServiceStatus, 30000);
        });
    </script>
</body>
</html>`
}
