package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// uploadHandler 处理文件上传请求
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
