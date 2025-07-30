package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

// testWebDAVAuthenticationFix 测试WebDAV认证修复
func testWebDAVAuthenticationFix() error {
	fmt.Println("🧪 测试WebDAV认证修复...")

	// 测试无认证访问（应该返回401）
	resp, err := http.Get("http://localhost:8080/webdav/")
	if err != nil {
		return fmt.Errorf("无法连接到WebDAV服务: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 401 {
		return fmt.Errorf("期望状态码401，实际收到: %d", resp.StatusCode)
	}

	// 检查WWW-Authenticate头是否包含charset
	authHeader := resp.Header.Get("WWW-Authenticate")
	if !strings.Contains(authHeader, "charset=\"UTF-8\"") {
		return fmt.Errorf("WWW-Authenticate头缺少charset参数: %s", authHeader)
	}

	// 检查Content-Type头
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html; charset=utf-8") {
		return fmt.Errorf("Content-Type头不正确: %s", contentType)
	}

	// 测试正确的认证
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:8080/webdav/", nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	// 添加Basic认证头
	auth := base64.StdEncoding.EncodeToString([]byte("webdav:webdav"))
	req.Header.Set("Authorization", "Basic "+auth)

	resp, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("认证请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 207 {
		return fmt.Errorf("认证后期望状态码200或207，实际收到: %d", resp.StatusCode)
	}

	fmt.Println("✅ WebDAV认证修复测试通过")
	return nil
}

// testHTTPProxyAuthentication 测试HTTP代理认证
func testHTTPProxyAuthentication() error {
	fmt.Println("🧪 测试HTTP代理认证...")

	// 测试无认证的代理请求（应该返回407）
	conn, err := net.Dial("tcp", "localhost:10808")
	if err != nil {
		return fmt.Errorf("无法连接到HTTP代理: %v", err)
	}
	defer conn.Close()

	// 发送CONNECT请求（无认证）
	request := "CONNECT httpbin.org:80 HTTP/1.1\r\nHost: httpbin.org:80\r\n\r\n"
	_, err = conn.Write([]byte(request))
	if err != nil {
		return fmt.Errorf("发送代理请求失败: %v", err)
	}

	// 读取响应
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return fmt.Errorf("读取代理响应失败: %v", err)
	}

	response := string(buffer[:n])
	if !strings.Contains(response, "407 Proxy Authentication Required") {
		return fmt.Errorf("期望407响应，实际收到: %s", response)
	}

	// 检查Proxy-Authenticate头
	if !strings.Contains(response, "Proxy-Authenticate: Basic") {
		return fmt.Errorf("缺少Proxy-Authenticate头: %s", response)
	}

	fmt.Println("✅ HTTP代理认证测试通过")
	return nil
}

// testSOCKS5Authentication 测试SOCKS5认证
func testSOCKS5Authentication() error {
	fmt.Println("🧪 测试SOCKS5认证...")

	conn, err := net.Dial("tcp", "localhost:1080")
	if err != nil {
		return fmt.Errorf("无法连接到SOCKS5代理: %v", err)
	}
	defer conn.Close()

	// 发送认证方法协商（只支持无认证）
	authReq := []byte{0x05, 0x01, 0x00} // VER=5, NMETHODS=1, METHOD=0(无认证)
	_, err = conn.Write(authReq)
	if err != nil {
		return fmt.Errorf("发送SOCKS5认证请求失败: %v", err)
	}

	// 读取认证响应
	authResp := make([]byte, 2)
	_, err = conn.Read(authResp)
	if err != nil {
		return fmt.Errorf("读取SOCKS5认证响应失败: %v", err)
	}

	// 应该返回0xFF（不支持的认证方法）
	if authResp[0] != 0x05 || authResp[1] != 0xFF {
		return fmt.Errorf("期望认证失败响应(05 FF)，实际收到: %02x %02x", authResp[0], authResp[1])
	}

	fmt.Println("✅ SOCKS5认证测试通过")
	return nil
}

// testWindowsWebDAVCompatibility 测试Windows WebDAV兼容性
func testWindowsWebDAVCompatibility() error {
	fmt.Println("🧪 测试Windows WebDAV兼容性...")

	// 模拟Windows WebDAV客户端的请求
	client := &http.Client{}
	
	// 测试PROPFIND请求（Windows WebDAV客户端常用）
	req, err := http.NewRequest("PROPFIND", "http://localhost:8080/webdav/", strings.NewReader(`<?xml version="1.0" encoding="utf-8"?>
<D:propfind xmlns:D="DAV:">
  <D:prop>
    <D:displayname/>
    <D:getcontentlength/>
    <D:getcontenttype/>
    <D:getlastmodified/>
    <D:resourcetype/>
  </D:prop>
</D:propfind>`))
	if err != nil {
		return fmt.Errorf("创建PROPFIND请求失败: %v", err)
	}

	// 添加Windows WebDAV客户端常用的头
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("Depth", "1")
	req.Header.Set("User-Agent", "Microsoft-WebDAV-MiniRedir/10.0.19041")
	
	// 添加认证
	auth := base64.StdEncoding.EncodeToString([]byte("webdav:webdav"))
	req.Header.Set("Authorization", "Basic "+auth)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("PROPFIND请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 207 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("PROPFIND期望状态码207，实际收到: %d, 响应: %s", resp.StatusCode, string(body))
	}

	fmt.Println("✅ Windows WebDAV兼容性测试通过")
	return nil
}

// startTestServer 启动测试服务器
func startTestServer() (*exec.Cmd, error) {
	fmt.Println("🚀 启动测试服务器...")
	
	// 构建项目
	buildCmd := exec.Command("go", "build", "-o", "sweb_test.exe", ".")
	buildCmd.Dir = ".."
	if err := buildCmd.Run(); err != nil {
		return nil, fmt.Errorf("构建失败: %v", err)
	}

	// 启动服务器（启用所有认证功能）
	cmd := exec.Command("./sweb_test.exe", 
		"-webdav", "-webdav-auth",
		"-proxy", "-proxy-auth", 
		"-socks5", "-socks5-auth")
	cmd.Dir = ".."
	
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("启动服务器失败: %v", err)
	}

	// 等待服务器启动
	time.Sleep(3 * time.Second)
	
	return cmd, nil
}

// stopTestServer 停止测试服务器
func stopTestServer(cmd *exec.Cmd) {
	if cmd != nil && cmd.Process != nil {
		cmd.Process.Kill()
		cmd.Wait()
	}
	
	// 清理测试文件
	os.Remove("../sweb_test.exe")
}

func main() {
	fmt.Println("🔧 开始认证功能修复测试")
	fmt.Println("=" * 50)

	// 启动测试服务器
	cmd, err := startTestServer()
	if err != nil {
		fmt.Printf("❌ 启动测试服务器失败: %v\n", err)
		os.Exit(1)
	}
	defer stopTestServer(cmd)

	// 运行测试
	tests := []struct {
		name string
		fn   func() error
	}{
		{"WebDAV认证修复", testWebDAVAuthenticationFix},
		{"HTTP代理认证", testHTTPProxyAuthentication},
		{"SOCKS5认证", testSOCKS5Authentication},
		{"Windows WebDAV兼容性", testWindowsWebDAVCompatibility},
	}

	var failedTests []string
	for _, test := range tests {
		if err := test.fn(); err != nil {
			fmt.Printf("❌ %s 失败: %v\n", test.name, err)
			failedTests = append(failedTests, test.name)
		}
	}

	fmt.Println("=" * 50)
	if len(failedTests) == 0 {
		fmt.Println("🎉 所有认证功能测试通过！")
	} else {
		fmt.Printf("❌ %d 个测试失败: %v\n", len(failedTests), failedTests)
		os.Exit(1)
	}
}
