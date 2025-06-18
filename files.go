package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

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
