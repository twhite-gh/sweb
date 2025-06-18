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
