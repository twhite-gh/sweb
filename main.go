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
	httpsEnabled   bool
	httpPort       int
	httpsPort      int
	certDir        string
	filesEnabled   bool
)

func main() {
	// 解析命令行参数
	var showHelp bool

	flag.BoolVar(&uploadEnabled, "upload", false, "启用文件上传功能")
	flag.BoolVar(&uploadEnabled, "enable-upload", false, "启用文件上传功能")
	flag.BoolVar(&filesEnabled, "files", false, "启用文件浏览功能")
	flag.BoolVar(&filesEnabled, "enable-files", false, "启用文件浏览功能")
	flag.BoolVar(&webdavEnabled, "webdav", false, "启用WebDAV服务")
	flag.BoolVar(&webdavEnabled, "enable-webdav", false, "启用WebDAV服务")
	flag.StringVar(&webdavDir, "webdav-dir", ".", "WebDAV服务的根目录")
	flag.BoolVar(&webdavReadonly, "webdav-readonly", false, "WebDAV服务只读模式")
	flag.BoolVar(&httpsEnabled, "https", false, "启用HTTPS服务")
	flag.BoolVar(&httpsEnabled, "enable-https", false, "启用HTTPS服务")
	flag.IntVar(&httpPort, "port", 8080, "指定HTTP服务器端口")
	flag.IntVar(&httpPort, "p", 8080, "指定HTTP服务器端口")
	flag.IntVar(&httpsPort, "https-port", 8443, "指定HTTPS服务器端口")
	flag.StringVar(&certDir, "cert-dir", "./cert", "SSL证书目录")
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

	// 添加状态API端点
	http.HandleFunc("/api/status", statusHandler)

	// 添加文件浏览API端点
	http.HandleFunc("/api/files", filesListHandler)

	// 根据参数决定是否启用文件浏览
	if filesEnabled {
		http.HandleFunc("/files", filesBrowseHandler)
		fmt.Println("✅ 文件浏览功能已启用")
	} else {
		http.HandleFunc("/files", filesDisabledHandler)
		fmt.Println("🔒 文件浏览功能已禁用 (使用 -files 参数启用)")
	}

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
	startServers()
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// 显示上传表单
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(generateUploadPageHTML()))
	} else if r.Method == "POST" {
		// 处理多文件上传
		handleMultipleFileUpload(w, r)
	} else {
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
	}
}

// handleMultipleFileUpload 处理多文件上传
func handleMultipleFileUpload(w http.ResponseWriter, r *http.Request) {
	// 解析多部分表单，限制最大内存为32MB
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "解析表单失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		http.Error(w, "没有选择文件", http.StatusBadRequest)
		return
	}

	var uploadResults []map[string]interface{}
	successCount := 0

	for _, fileHeader := range files {
		result := map[string]interface{}{
			"filename": fileHeader.Filename,
			"size":     fileHeader.Size,
		}

		// 打开上传的文件
		file, err := fileHeader.Open()
		if err != nil {
			result["success"] = false
			result["error"] = "无法打开文件: " + err.Error()
			uploadResults = append(uploadResults, result)
			continue
		}

		// 创建目标文件
		dst, err := os.Create(filepath.Join("web", fileHeader.Filename))
		if err != nil {
			file.Close()
			result["success"] = false
			result["error"] = "无法创建目标文件: " + err.Error()
			uploadResults = append(uploadResults, result)
			continue
		}

		// 复制文件内容
		_, err = io.Copy(dst, file)
		file.Close()
		dst.Close()

		if err != nil {
			result["success"] = false
			result["error"] = "无法保存文件: " + err.Error()
		} else {
			result["success"] = true
			result["url"] = "/" + fileHeader.Filename
			successCount++
		}

		uploadResults = append(uploadResults, result)
	}

	// 返回JSON响应
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"success":      successCount > 0,
		"total":        len(files),
		"successCount": successCount,
		"results":      uploadResults,
	}

	json.NewEncoder(w).Encode(response)
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
	fmt.Println("  -files, --enable-files      启用文件浏览功能 (默认: 禁用)")
	fmt.Println("  -webdav, --enable-webdav    启用WebDAV服务 (默认: 禁用)")
	fmt.Println("  -webdav-dir <目录>          WebDAV服务的根目录 (默认: 当前目录)")
	fmt.Println("  -webdav-readonly            WebDAV服务只读模式 (默认: 读写)")
	fmt.Println("  -https, --enable-https      启用HTTPS服务 (默认: 禁用)")
	fmt.Println("  -port, -p <端口>           指定HTTP服务器端口 (默认: 8080)")
	fmt.Println("  -https-port <端口>         指定HTTPS服务器端口 (默认: 8443)")
	fmt.Println("  -cert-dir <目录>           SSL证书目录 (默认: ./cert)")
	fmt.Println("  -help, -h                  显示此帮助信息")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  sweb.exe                           # 启动HTTP服务器，仅提供静态文件服务")
	fmt.Println("  sweb.exe -upload                   # 启动服务器并启用文件上传功能")
	fmt.Println("  sweb.exe -files                    # 启动服务器并启用文件浏览功能")
	fmt.Println("  sweb.exe -https                    # 启动HTTP和HTTPS服务器")
	fmt.Println("  sweb.exe -webdav                   # 启动服务器并启用WebDAV服务")
	fmt.Println("  sweb.exe -webdav -webdav-readonly  # 启动只读WebDAV服务")
	fmt.Println("  sweb.exe -webdav -webdav-dir /data # 指定WebDAV目录")
	fmt.Println("  sweb.exe -upload -files -webdav -https # 启用所有功能")
	fmt.Println("  sweb.exe -https -https-port 9443   # 指定HTTPS端口")
	fmt.Println()
	fmt.Println("HTTPS证书:")
	fmt.Println("  证书文件: ./cert/server.crt")
	fmt.Println("  私钥文件: ./cert/server.key")
	fmt.Println("  可以使用openssl生成自签名证书用于测试")
	fmt.Println()
	fmt.Println("访问地址:")
	fmt.Println("  HTTP: http://localhost:8080")
	fmt.Println("  HTTPS: https://localhost:8443 (如果启用)")
	fmt.Println("  WebDAV: http://localhost:8080/webdav")
	fmt.Println()
	fmt.Println("安全说明:")
	fmt.Println("  文件上传、文件浏览、WebDAV和HTTPS功能默认禁用以确保服务器安全。")
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

// statusHandler 处理状态查询请求
func statusHandler(w http.ResponseWriter, r *http.Request) {
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
		"files": map[string]interface{}{
			"enabled": filesEnabled,
			"status": func() string {
				if filesEnabled {
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
		"https": map[string]interface{}{
			"enabled":   httpsEnabled,
			"httpPort":  httpPort,
			"httpsPort": httpsPort,
			"certDir":   certDir,
			"certStatus": func() string {
				if !httpsEnabled {
					return "disabled"
				}
				if err := checkCertificates(); err != nil {
					return "cert-error"
				}
				return "enabled"
			}(),
		},
	}

	json.NewEncoder(w).Encode(response)
}

// filesListHandler 处理文件列表查询请求
func filesListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 读取web目录下的文件
	webDir := "./web"
	files, err := os.ReadDir(webDir)
	if err != nil {
		http.Error(w, "无法读取文件目录", http.StatusInternalServerError)
		return
	}

	var fileList []map[string]interface{}
	for _, file := range files {
		if file.Name() == "index.html" || file.Name() == "index.htm" {
			continue // 跳过默认页面
		}

		fileInfo, err := file.Info()
		if err != nil {
			continue
		}

		fileData := map[string]interface{}{
			"name":    file.Name(),
			"size":    fileInfo.Size(),
			"modTime": fileInfo.ModTime().Format("2006-01-02 15:04:05"),
			"isDir":   file.IsDir(),
			"url":     "/" + file.Name(),
		}

		// 如果是文件，添加文件类型信息
		if !file.IsDir() {
			ext := strings.ToLower(filepath.Ext(file.Name()))
			fileData["extension"] = ext
			fileData["type"] = getFileType(ext)
		}

		fileList = append(fileList, fileData)
	}

	response := map[string]interface{}{
		"files": fileList,
		"total": len(fileList),
	}

	json.NewEncoder(w).Encode(response)
}

// getFileType 根据文件扩展名返回文件类型
func getFileType(ext string) string {
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp":
		return "image"
	case ".mp4", ".avi", ".mov", ".wmv", ".flv", ".mkv":
		return "video"
	case ".mp3", ".wav", ".flac", ".aac", ".ogg":
		return "audio"
	case ".pdf":
		return "pdf"
	case ".doc", ".docx":
		return "document"
	case ".xls", ".xlsx":
		return "spreadsheet"
	case ".ppt", ".pptx":
		return "presentation"
	case ".txt", ".md":
		return "text"
	case ".zip", ".rar", ".7z", ".tar", ".gz":
		return "archive"
	case ".html", ".htm":
		return "html"
	case ".css":
		return "css"
	case ".js":
		return "javascript"
	case ".json":
		return "json"
	case ".xml":
		return "xml"
	default:
		return "file"
	}
}

// filesBrowseHandler 处理文件浏览页面请求
func filesBrowseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>文件浏览 - 简单Web文件服务器</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            max-width: 1000px;
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
            margin-bottom: 30px;
        }
        .toolbar {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 20px;
            padding: 15px;
            background: #f8f9fa;
            border-radius: 5px;
        }
        .file-count {
            color: #666;
            font-size: 14px;
        }
        .refresh-btn {
            background: #007acc;
            color: white;
            border: none;
            padding: 8px 16px;
            border-radius: 5px;
            cursor: pointer;
            transition: background 0.3s;
        }
        .refresh-btn:hover {
            background: #005a9e;
        }
        .files-grid {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
            gap: 15px;
            margin-bottom: 30px;
        }
        .file-card {
            border: 1px solid #ddd;
            border-radius: 8px;
            padding: 15px;
            background: white;
            transition: all 0.3s;
            cursor: pointer;
        }
        .file-card:hover {
            box-shadow: 0 4px 12px rgba(0,0,0,0.15);
            transform: translateY(-2px);
        }
        .file-icon {
            font-size: 48px;
            text-align: center;
            margin-bottom: 10px;
        }
        .file-name {
            font-weight: bold;
            margin-bottom: 5px;
            word-break: break-all;
        }
        .file-info {
            font-size: 12px;
            color: #666;
            display: flex;
            justify-content: space-between;
        }
        .file-size {
            font-weight: 500;
        }
        .file-date {
            font-style: italic;
        }
        .empty-state {
            text-align: center;
            padding: 60px 20px;
            color: #666;
        }
        .empty-icon {
            font-size: 64px;
            margin-bottom: 20px;
        }
        .back-link {
            display: inline-block;
            background: #6c757d;
            color: white;
            padding: 10px 20px;
            text-decoration: none;
            border-radius: 5px;
            margin-bottom: 20px;
            transition: background 0.3s;
        }
        .back-link:hover {
            background: #545b62;
        }
        .loading {
            text-align: center;
            padding: 40px;
            color: #666;
        }
        .error {
            background: #f8d7da;
            color: #721c24;
            padding: 15px;
            border-radius: 5px;
            margin-bottom: 20px;
            border: 1px solid #f5c6cb;
        }
    </style>
</head>
<body>
    <div class="container">
        <a href="/" class="back-link">← 返回首页</a>

        <h1>📂 文件浏览</h1>

        <div class="toolbar">
            <div class="file-count" id="file-count">正在加载...</div>
            <button class="refresh-btn" onclick="loadFiles()">🔄 刷新</button>
        </div>

        <div id="error-message" class="error" style="display: none;"></div>

        <div id="loading" class="loading">
            <div>🔄 正在加载文件列表...</div>
        </div>

        <div id="files-container" class="files-grid" style="display: none;"></div>

        <div id="empty-state" class="empty-state" style="display: none;">
            <div class="empty-icon">📁</div>
            <h3>暂无文件</h3>
            <p>web目录下还没有文件。</p>
            <p>您可以通过<a href="/upload">文件上传</a>功能添加文件。</p>
        </div>
    </div>

    <script>
        // 文件类型图标映射
        const fileIcons = {
            'image': '🖼️',
            'video': '🎬',
            'audio': '🎵',
            'pdf': '📄',
            'document': '📝',
            'spreadsheet': '📊',
            'presentation': '📋',
            'text': '📄',
            'archive': '📦',
            'html': '🌐',
            'css': '🎨',
            'javascript': '⚡',
            'json': '📋',
            'xml': '📄',
            'file': '📄'
        };

        // 格式化文件大小
        function formatFileSize(bytes) {
            if (bytes === 0) return '0 B';
            const k = 1024;
            const sizes = ['B', 'KB', 'MB', 'GB'];
            const i = Math.floor(Math.log(bytes) / Math.log(k));
            return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
        }

        // 加载文件列表
        function loadFiles() {
            const loading = document.getElementById('loading');
            const filesContainer = document.getElementById('files-container');
            const emptyState = document.getElementById('empty-state');
            const fileCount = document.getElementById('file-count');
            const errorMessage = document.getElementById('error-message');

            // 显示加载状态
            loading.style.display = 'block';
            filesContainer.style.display = 'none';
            emptyState.style.display = 'none';
            errorMessage.style.display = 'none';

            fetch('/api/files')
                .then(response => {
                    if (!response.ok) {
                        throw new Error('网络请求失败');
                    }
                    return response.json();
                })
                .then(data => {
                    loading.style.display = 'none';

                    if (data.files && data.files.length > 0) {
                        displayFiles(data.files);
                        fileCount.textContent = '共 ' + data.total + ' 个文件';
                        filesContainer.style.display = 'grid';
                    } else {
                        emptyState.style.display = 'block';
                        fileCount.textContent = '共 0 个文件';
                    }
                })
                .catch(error => {
                    loading.style.display = 'none';
                    errorMessage.textContent = '加载文件列表失败: ' + error.message;
                    errorMessage.style.display = 'block';
                    fileCount.textContent = '加载失败';
                });
        }

        // 显示文件列表
        function displayFiles(files) {
            const container = document.getElementById('files-container');
            container.innerHTML = '';

            files.forEach(file => {
                const fileCard = document.createElement('div');
                fileCard.className = 'file-card';
                fileCard.onclick = () => {
                    if (!file.isDir) {
                        window.open(file.url, '_blank');
                    }
                };

                const icon = fileIcons[file.type] || fileIcons['file'];

                fileCard.innerHTML = ` + "`" + `
                    <div class="file-icon">${icon}</div>
                    <div class="file-name">${file.name}</div>
                    <div class="file-info">
                        <span class="file-size">${file.isDir ? '文件夹' : formatFileSize(file.size)}</span>
                        <span class="file-date">${file.modTime}</span>
                    </div>
                ` + "`" + `;

                container.appendChild(fileCard);
            });
        }

        // 页面加载时获取文件列表
        document.addEventListener('DOMContentLoaded', function() {
            loadFiles();
        });
    </script>
</body>
</html>`

	w.Write([]byte(html))
}

// filesDisabledHandler 处理文件浏览功能被禁用时的请求
func filesDisabledHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte(`
        <!DOCTYPE html>
        <html lang="zh-CN">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>文件浏览功能已禁用</title>
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
                <h1>文件浏览功能已禁用</h1>
                <p>出于安全考虑，文件浏览功能默认处于禁用状态。</p>
                <p>如需启用文件浏览功能，请使用以下命令重新启动服务器：</p>

                <div class="command">sweb.exe -files</div>
                <p>或</p>
                <div class="command">sweb.exe --enable-files</div>

                <p>您也可以使用 <code>sweb.exe -help</code> 查看所有可用选项。</p>

                <a href="/" class="back-link">← 返回首页</a>
            </div>
        </body>
        </html>
    `))
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
            <strong>🔒 HTTPS服务</strong><br>
            <span id="https-feature-description">HTTPS服务可通过命令行参数启用，提供加密的安全连接。</span>
            <span id="https-status" class="status-indicator loading">🔄 检查中...</span>
            <div id="https-info" class="webdav-info hidden">
                <strong>HTTP地址:</strong> <code id="http-url">http://localhost:8080</code><br>
                <strong>HTTPS地址:</strong> <code id="https-url">https://localhost:8443</code><br>
                <strong>证书目录:</strong> <code id="cert-directory">./cert</code><br>
                <strong>证书状态:</strong> <span id="cert-status">检查中</span>
            </div>
        </div>

        <div class="feature">
            <strong>📊 功能状态展示</strong><br>
            默认页面实时显示服务器各项功能的启用状态，包括文件上传、WebDAV和HTTPS服务。
        </div>

        <div class="feature">
            <strong>📂 文件浏览下载</strong><br>
            <span id="files-feature-description">专门的文件浏览页面，展示web目录下的所有文件，支持在线浏览和下载。</span>
            <span id="files-status" class="status-indicator loading">🔄 检查中...</span>
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

        <div id="files-section">
            <p><strong>文件浏览：</strong></p>
            <div id="files-enabled-content" class="hidden">
                <a href="/files" class="button" id="files-button">📂 浏览文件</a>
                <p><small>查看和下载web目录中的所有文件</small></p>
            </div>
            <div id="files-disabled-content" class="hidden">
                <p>要启用文件浏览功能，请使用以下命令启动服务器：</p>
                <code>sweb.exe -files</code> 或 <code>sweb.exe --enable-files</code>
                <br><br>
                <a href="/files" class="button disabled" id="files-button-disabled">📂 文件浏览已禁用</a>
            </div>
        </div>

        <p><strong>服务器信息：</strong></p>
        <ul>
            <li>服务端口: <code>8080</code></li>
            <li>文件目录: <code>./web</code></li>
            <li>上传地址: <code>/upload</code></li>
            <li>文件浏览: <code>/files</code></li>
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
                        <span id="footer-files-link" class="hidden"><a href="/files" class="button">浏览文件</a></span>
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
            fetch('/api/status')
                .then(response => response.json())
                .then(data => {
                    updateUploadStatus(data.upload);
                    updateFilesStatus(data.files);
                    updateWebDAVStatus(data.webdav);
                    updateHTTPSStatus(data.https);
                    updateUsageInstructions(data.upload, data.files, data.webdav, data.https);
                })
                .catch(error => {
                    console.error('检查服务状态失败:', error);
                    // 如果API调用失败，显示默认的禁用状态
                    updateUploadStatus({enabled: false, status: 'disabled'});
                    updateFilesStatus({enabled: false, status: 'disabled'});
                    updateWebDAVStatus({enabled: false, status: 'disabled'});
                    updateHTTPSStatus({enabled: false, status: 'disabled'});
                    updateUsageInstructions({enabled: false}, {enabled: false}, {enabled: false}, {enabled: false});
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

        // 更新页面上的文件浏览状态显示
        function updateFilesStatus(filesData) {
            const statusElement = document.getElementById('files-status');
            const descriptionElement = document.getElementById('files-feature-description');
            const filesEnabledContent = document.getElementById('files-enabled-content');
            const filesDisabledContent = document.getElementById('files-disabled-content');

            if (filesData.enabled) {
                statusElement.textContent = '✅ 已启用';
                statusElement.className = 'status-indicator status-enabled';
                descriptionElement.textContent = '专门的文件浏览页面，展示web目录下的所有文件，支持在线浏览和下载。';

                filesEnabledContent.classList.remove('hidden');
                filesDisabledContent.classList.add('hidden');
            } else {
                statusElement.textContent = '🔒 已禁用';
                statusElement.className = 'status-indicator status-disabled';
                descriptionElement.textContent = '文件浏览功能可通过命令行参数启用，确保服务器安全。';

                filesEnabledContent.classList.add('hidden');
                filesDisabledContent.classList.remove('hidden');
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

        // 更新页面上的HTTPS状态显示
        function updateHTTPSStatus(httpsData) {
            const statusElement = document.getElementById('https-status');
            const descriptionElement = document.getElementById('https-feature-description');
            const httpsInfo = document.getElementById('https-info');
            const httpUrl = document.getElementById('http-url');
            const httpsUrl = document.getElementById('https-url');
            const certDirectory = document.getElementById('cert-directory');
            const certStatus = document.getElementById('cert-status');

            if (httpsData.enabled) {
                if (httpsData.certStatus === 'enabled') {
                    statusElement.textContent = '✅ 已启用';
                    statusElement.className = 'status-indicator status-enabled';
                    descriptionElement.textContent = 'HTTPS服务已启用，提供加密的安全连接。';
                    certStatus.textContent = '证书正常';
                    certStatus.className = 'status-enabled';
                } else if (httpsData.certStatus === 'cert-error') {
                    statusElement.textContent = '⚠️ 证书错误';
                    statusElement.className = 'status-indicator status-readonly';
                    descriptionElement.textContent = 'HTTPS服务已启用，但证书文件有问题。';
                    certStatus.textContent = '证书错误';
                    certStatus.className = 'status-disabled';
                }

                httpUrl.textContent = 'http://localhost:' + httpsData.httpPort;
                httpsUrl.textContent = 'https://localhost:' + httpsData.httpsPort;
                certDirectory.textContent = httpsData.certDir || './cert';
                httpsInfo.classList.remove('hidden');
            } else {
                statusElement.textContent = '🔒 已禁用';
                statusElement.className = 'status-indicator status-disabled';
                descriptionElement.textContent = 'HTTPS服务可通过命令行参数启用，提供加密的安全连接。';

                httpsInfo.classList.add('hidden');
            }
        }

        // 更新使用说明和页脚
        function updateUsageInstructions(uploadData, filesData, webdavData, httpsData) {
            const usageUploadEnabled = document.getElementById('usage-upload-enabled');
            const usageWebdavEnabled = document.getElementById('usage-webdav-enabled');
            const usageDisabled = document.getElementById('usage-disabled');
            const footerEnabled = document.getElementById('footer-enabled');
            const footerDisabled = document.getElementById('footer-disabled');
            const footerFilesLink = document.getElementById('footer-files-link');
            const footerUploadLink = document.getElementById('footer-upload-link');
            const footerWebdavLink = document.getElementById('footer-webdav-link');

            const anyEnabled = uploadData.enabled || filesData.enabled || webdavData.enabled || httpsData.enabled;

            if (anyEnabled) {
                footerEnabled.classList.remove('hidden');
                footerDisabled.classList.add('hidden');
                usageDisabled.classList.add('hidden');

                if (filesData.enabled) {
                    footerFilesLink.classList.remove('hidden');
                } else {
                    footerFilesLink.classList.add('hidden');
                }

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

// checkCertificates 检查SSL证书文件
func checkCertificates() error {
	certFile := filepath.Join(certDir, "server.crt")
	keyFile := filepath.Join(certDir, "server.key")

	// 检查证书目录是否存在
	if _, err := os.Stat(certDir); os.IsNotExist(err) {
		return fmt.Errorf("证书目录不存在: %s", certDir)
	}

	// 检查证书文件是否存在
	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		return fmt.Errorf("证书文件不存在: %s", certFile)
	}

	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		return fmt.Errorf("私钥文件不存在: %s", keyFile)
	}

	return nil
}

// startServers 启动HTTP和HTTPS服务器
func startServers() {
	// 创建一个通道来接收错误
	errChan := make(chan error, 2)

	// 启动HTTP服务器
	go func() {
		fmt.Printf("🌐 HTTP服务器启动在 http://localhost:%d\n", httpPort)
		err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil)
		errChan <- fmt.Errorf("HTTP服务器错误: %v", err)
	}()

	// 如果启用了HTTPS，启动HTTPS服务器
	if httpsEnabled {
		go func() {
			// 检查证书文件
			if err := checkCertificates(); err != nil {
				errChan <- fmt.Errorf("HTTPS证书检查失败: %v", err)
				return
			}

			certFile := filepath.Join(certDir, "server.crt")
			keyFile := filepath.Join(certDir, "server.key")

			fmt.Printf("🔒 HTTPS服务器启动在 https://localhost:%d\n", httpsPort)
			fmt.Printf("📁 证书目录: %s\n", certDir)
			err := http.ListenAndServeTLS(fmt.Sprintf(":%d", httpsPort), certFile, keyFile, nil)
			errChan <- fmt.Errorf("HTTPS服务器错误: %v", err)
		}()
	} else {
		fmt.Println("🔒 HTTPS服务已禁用 (使用 -https 参数启用)")
	}

	// 等待任一服务器出错
	log.Fatal(<-errChan)
}

// generateUploadPageHTML 生成支持拖拽多文件上传的页面HTML
func generateUploadPageHTML() string {
	return `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>文件上传 - 简单Web文件服务器</title>
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
            margin-bottom: 30px;
        }
        .upload-area {
            border: 3px dashed #ccc;
            border-radius: 10px;
            padding: 60px 20px;
            text-align: center;
            margin-bottom: 30px;
            transition: all 0.3s;
            cursor: pointer;
            background: #fafafa;
        }
        .upload-area.dragover {
            border-color: #007acc;
            background: #e7f3ff;
        }
        .upload-area:hover {
            border-color: #007acc;
            background: #f0f8ff;
        }
        .upload-icon {
            font-size: 48px;
            margin-bottom: 20px;
            color: #666;
        }
        .upload-text {
            font-size: 18px;
            color: #666;
            margin-bottom: 10px;
        }
        .upload-hint {
            font-size: 14px;
            color: #999;
        }
        .file-input {
            display: none;
        }
        .upload-btn {
            background: #007acc;
            color: white;
            border: none;
            padding: 12px 24px;
            border-radius: 5px;
            cursor: pointer;
            font-size: 16px;
            margin: 10px;
            transition: background 0.3s;
        }
        .upload-btn:hover {
            background: #005a9e;
        }
        .upload-btn:disabled {
            background: #ccc;
            cursor: not-allowed;
        }
        .file-list {
            margin-top: 20px;
        }
        .file-item {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 10px;
            border: 1px solid #ddd;
            border-radius: 5px;
            margin-bottom: 10px;
            background: white;
        }
        .file-info {
            flex: 1;
        }
        .file-name {
            font-weight: bold;
            margin-bottom: 5px;
        }
        .file-size {
            font-size: 12px;
            color: #666;
        }
        .file-status {
            padding: 5px 10px;
            border-radius: 3px;
            font-size: 12px;
            font-weight: bold;
        }
        .status-waiting {
            background: #fff3cd;
            color: #856404;
        }
        .status-uploading {
            background: #cce5ff;
            color: #004085;
        }
        .status-success {
            background: #d4edda;
            color: #155724;
        }
        .status-error {
            background: #f8d7da;
            color: #721c24;
        }
        .progress-bar {
            width: 100%;
            height: 6px;
            background: #f0f0f0;
            border-radius: 3px;
            margin: 5px 0;
            overflow: hidden;
        }
        .progress-fill {
            height: 100%;
            background: #007acc;
            width: 0%;
            transition: width 0.3s;
        }
        .remove-btn {
            background: #dc3545;
            color: white;
            border: none;
            padding: 5px 10px;
            border-radius: 3px;
            cursor: pointer;
            font-size: 12px;
        }
        .remove-btn:hover {
            background: #c82333;
        }
        .back-link {
            display: inline-block;
            background: #6c757d;
            color: white;
            padding: 10px 20px;
            text-decoration: none;
            border-radius: 5px;
            margin-bottom: 20px;
            transition: background 0.3s;
        }
        .back-link:hover {
            background: #545b62;
        }
        .upload-summary {
            background: #f8f9fa;
            padding: 15px;
            border-radius: 5px;
            margin-top: 20px;
            display: none;
        }
        .links-section {
            margin-top: 30px;
            padding-top: 20px;
            border-top: 1px solid #eee;
        }
        .link-btn {
            display: inline-block;
            background: #28a745;
            color: white;
            padding: 8px 16px;
            text-decoration: none;
            border-radius: 5px;
            margin: 5px;
            font-size: 14px;
            transition: background 0.3s;
        }
        .link-btn:hover {
            background: #218838;
        }
    </style>
</head>
<body>
    <div class="container">
        <a href="/" class="back-link">← 返回首页</a>

        <h1>📤 文件上传</h1>

        <div class="upload-area" id="uploadArea">
            <div class="upload-icon">📁</div>
            <div class="upload-text">拖拽文件到此处或点击选择文件</div>
            <div class="upload-hint">支持多文件同时上传</div>
            <input type="file" id="fileInput" class="file-input" multiple>
            <button type="button" class="upload-btn" onclick="document.getElementById('fileInput').click()">
                选择文件
            </button>
        </div>

        <div class="file-list" id="fileList"></div>

        <div style="text-align: center; margin-top: 20px;">
            <button type="button" class="upload-btn" id="uploadBtn" onclick="uploadFiles()" disabled>
                开始上传
            </button>
            <button type="button" class="upload-btn" onclick="clearFiles()">
                清空列表
            </button>
        </div>

        <div class="upload-summary" id="uploadSummary"></div>

        <div class="links-section">
            <strong>相关链接：</strong><br>
            <a href="/files" class="link-btn">📂 浏览文件</a>
            <a href="/" class="link-btn">🏠 返回首页</a>
        </div>
    </div>

    <script>
        let selectedFiles = [];
        let uploadInProgress = false;

        // 格式化文件大小
        function formatFileSize(bytes) {
            if (bytes === 0) return '0 B';
            const k = 1024;
            const sizes = ['B', 'KB', 'MB', 'GB'];
            const i = Math.floor(Math.log(bytes) / Math.log(k));
            return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
        }

        // 初始化拖拽功能
        function initDragAndDrop() {
            const uploadArea = document.getElementById('uploadArea');

            ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
                uploadArea.addEventListener(eventName, preventDefaults, false);
            });

            function preventDefaults(e) {
                e.preventDefault();
                e.stopPropagation();
            }

            ['dragenter', 'dragover'].forEach(eventName => {
                uploadArea.addEventListener(eventName, highlight, false);
            });

            ['dragleave', 'drop'].forEach(eventName => {
                uploadArea.addEventListener(eventName, unhighlight, false);
            });

            function highlight(e) {
                uploadArea.classList.add('dragover');
            }

            function unhighlight(e) {
                uploadArea.classList.remove('dragover');
            }

            uploadArea.addEventListener('drop', handleDrop, false);

            function handleDrop(e) {
                const dt = e.dataTransfer;
                const files = dt.files;
                handleFiles(files);
            }
        }

        // 处理文件选择
        function handleFiles(files) {
            for (let i = 0; i < files.length; i++) {
                const file = files[i];
                if (!selectedFiles.find(f => f.name === file.name && f.size === file.size)) {
                    selectedFiles.push(file);
                }
            }
            updateFileList();
            updateUploadButton();
        }

        // 更新文件列表显示
        function updateFileList() {
            const fileList = document.getElementById('fileList');
            fileList.innerHTML = '';

            selectedFiles.forEach((file, index) => {
                const fileItem = document.createElement('div');
                fileItem.className = 'file-item';
                fileItem.innerHTML = ` + "`" + `
                    <div class="file-info">
                        <div class="file-name">${file.name}</div>
                        <div class="file-size">${formatFileSize(file.size)}</div>
                        <div class="progress-bar" id="progress-${index}" style="display: none;">
                            <div class="progress-fill" id="progress-fill-${index}"></div>
                        </div>
                    </div>
                    <div class="file-status status-waiting" id="status-${index}">等待上传</div>
                    <button class="remove-btn" onclick="removeFile(${index})" ${uploadInProgress ? 'disabled' : ''}>
                        删除
                    </button>
                ` + "`" + `;
                fileList.appendChild(fileItem);
            });
        }

        // 移除文件
        function removeFile(index) {
            if (uploadInProgress) return;
            selectedFiles.splice(index, 1);
            updateFileList();
            updateUploadButton();
        }

        // 清空文件列表
        function clearFiles() {
            if (uploadInProgress) return;
            selectedFiles = [];
            updateFileList();
            updateUploadButton();
            document.getElementById('uploadSummary').style.display = 'none';
        }

        // 更新上传按钮状态
        function updateUploadButton() {
            const uploadBtn = document.getElementById('uploadBtn');
            uploadBtn.disabled = selectedFiles.length === 0 || uploadInProgress;
        }

        // 上传文件
        async function uploadFiles() {
            if (selectedFiles.length === 0 || uploadInProgress) return;

            uploadInProgress = true;
            updateUploadButton();

            const formData = new FormData();
            selectedFiles.forEach(file => {
                formData.append('files', file);
            });

            // 更新所有文件状态为上传中
            selectedFiles.forEach((file, index) => {
                const statusElement = document.getElementById(` + "`" + `status-${index}` + "`" + `);
                const progressElement = document.getElementById(` + "`" + `progress-${index}` + "`" + `);
                statusElement.textContent = '上传中...';
                statusElement.className = 'file-status status-uploading';
                progressElement.style.display = 'block';
            });

            try {
                const response = await fetch('/upload', {
                    method: 'POST',
                    body: formData
                });

                const result = await response.json();
                handleUploadResult(result);
            } catch (error) {
                console.error('上传失败:', error);
                // 更新所有文件状态为错误
                selectedFiles.forEach((file, index) => {
                    const statusElement = document.getElementById(` + "`" + `status-${index}` + "`" + `);
                    statusElement.textContent = '上传失败';
                    statusElement.className = 'file-status status-error';
                });
                showUploadSummary(0, selectedFiles.length, []);
            }

            uploadInProgress = false;
            updateUploadButton();
        }

        // 处理上传结果
        function handleUploadResult(result) {
            result.results.forEach((fileResult, index) => {
                const statusElement = document.getElementById(` + "`" + `status-${index}` + "`" + `);
                const progressFill = document.getElementById(` + "`" + `progress-fill-${index}` + "`" + `);

                progressFill.style.width = '100%';

                if (fileResult.success) {
                    statusElement.textContent = '上传成功';
                    statusElement.className = 'file-status status-success';
                } else {
                    statusElement.textContent = '上传失败';
                    statusElement.className = 'file-status status-error';
                }
            });

            showUploadSummary(result.successCount, result.total, result.results.filter(r => r.success));
        }

        // 显示上传摘要
        function showUploadSummary(successCount, totalCount, successFiles) {
            const summary = document.getElementById('uploadSummary');
            let summaryHTML = ` + "`" + `
                <h3>上传完成</h3>
                <p>成功上传 ${successCount} 个文件，共 ${totalCount} 个文件</p>
            ` + "`" + `;

            if (successFiles.length > 0) {
                summaryHTML += '<p><strong>成功上传的文件：</strong></p><ul>';
                successFiles.forEach(file => {
                    summaryHTML += ` + "`" + `<li><a href="${file.url}" target="_blank">${file.filename}</a></li>` + "`" + `;
                });
                summaryHTML += '</ul>';
            }

            summary.innerHTML = summaryHTML;
            summary.style.display = 'block';
        }

        // 文件输入变化事件
        document.getElementById('fileInput').addEventListener('change', function(e) {
            handleFiles(e.target.files);
        });

        // 初始化
        document.addEventListener('DOMContentLoaded', function() {
            initDragAndDrop();
        });
    </script>
</body>
</html>`
}
