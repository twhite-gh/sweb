package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

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
	// 检查是否至少启用了一个服务器
	if !httpEnabled && !httpsEnabled {
		log.Fatal("错误: 必须至少启用HTTP或HTTPS服务器中的一个")
	}

	// 创建一个通道来接收错误
	errChan := make(chan error, 2)
	serverCount := 0

	// 如果启用了HTTP，启动HTTP服务器
	if httpEnabled {
		serverCount++
		go func() {
			fmt.Printf("🌐 HTTP服务器启动在 http://localhost:%d\n", httpPort)
			err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil)
			errChan <- fmt.Errorf("HTTP服务器错误: %v", err)
		}()
	} else {
		fmt.Println("🔒 HTTP服务已禁用 (使用 -http 参数启用)")
	}

	// 如果启用了HTTPS，启动HTTPS服务器
	if httpsEnabled {
		serverCount++
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

// serveIndexOrAbout 处理根路径的请求，根据 index.html 的存在情况返回不同页面
func serveIndexOrAbout(w http.ResponseWriter, r *http.Request, webFilesDir string) {
	indexPath := filepath.Join(webFilesDir, "index.html")
	aboutPath := filepath.Join(webFilesDir, "about.html")

	// 检查 index.html 是否存在
	_, err := os.Stat(indexPath)
	if err == nil {
		// index.html 存在，返回 index.html
		http.ServeFile(w, r, indexPath)
		log.Printf("Serving %s for path /", indexPath)
	} else if os.IsNotExist(err) {
		// index.html 不存在，返回 about.html
		// 再次检查 about.html 是否存在，以防万一
		_, err := os.Stat(aboutPath)
		if err == nil {
			http.ServeFile(w, r, aboutPath)
			log.Printf("Serving %s for path / (index.html not found)", aboutPath)
		} else if os.IsNotExist(err) {
			// about.html 也不存在
			http.NotFound(w, r)
			log.Printf("Neither index.html nor about.html found in %s", webFilesDir)
		} else {
			// 其他文件系统错误
			http.Error(w, "Error checking about.html: "+err.Error(), http.StatusInternalServerError)
			log.Printf("Error checking about.html: %v", err)
		}
	} else {
		// 其他文件系统错误（例如权限问题）
		http.Error(w, "Error checking index.html: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Error checking index.html: %v", err)
	}
}
