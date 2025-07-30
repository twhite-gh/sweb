package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func main() {
	fmt.Println("=== WebDAV Windows兼容性测试 ===")

	// 测试服务器地址
	baseURL := "http://localhost:8082/webdav"
	username := "webdav"
	password := "webdav"

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 测试1: OPTIONS请求（Windows WebDAV客户端的第一步）
	fmt.Println("\n1. 测试OPTIONS请求...")
	testOPTIONS(client, baseURL, username, password)

	// 测试2: PROPFIND请求（获取目录信息）
	fmt.Println("\n2. 测试PROPFIND请求...")
	testPROPFIND(client, baseURL, username, password)

	// 测试3: 无认证的OPTIONS请求
	fmt.Println("\n3. 测试无认证的OPTIONS请求...")
	testOPTIONSNoAuth(client, baseURL)

	// 测试4: 错误认证的请求
	fmt.Println("\n4. 测试错误认证...")
	testWrongAuth(client, baseURL)

	fmt.Println("\n=== 测试完成 ===")
}

func testOPTIONS(client *http.Client, baseURL, username, password string) {
	req, err := http.NewRequest("OPTIONS", baseURL, nil)
	if err != nil {
		fmt.Printf("❌ 创建OPTIONS请求失败: %v\n", err)
		return
	}

	// 添加Basic认证
	auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("User-Agent", "Microsoft-WebDAV-MiniRedir/10.0.19041")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ OPTIONS请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("✅ OPTIONS响应状态: %s\n", resp.Status)

	// 检查关键的WebDAV头部
	checkHeader(resp, "DAV")
	checkHeader(resp, "MS-Author-Via")
	checkHeader(resp, "Allow")
	checkHeader(resp, "Server")
}

func testPROPFIND(client *http.Client, baseURL, username, password string) {
	propfindBody := `<?xml version="1.0" encoding="utf-8"?>
<D:propfind xmlns:D="DAV:">
  <D:prop>
    <D:displayname/>
    <D:getcontentlength/>
    <D:getcontenttype/>
    <D:getlastmodified/>
    <D:resourcetype/>
  </D:prop>
</D:propfind>`

	req, err := http.NewRequest("PROPFIND", baseURL, strings.NewReader(propfindBody))
	if err != nil {
		fmt.Printf("❌ 创建PROPFIND请求失败: %v\n", err)
		return
	}

	// 添加Basic认证
	auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/xml; charset=utf-8")
	req.Header.Set("Depth", "1")
	req.Header.Set("User-Agent", "Microsoft-WebDAV-MiniRedir/10.0.19041")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ PROPFIND请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("✅ PROPFIND响应状态: %s\n", resp.Status)

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ 读取PROPFIND响应失败: %v\n", err)
		return
	}

	if len(body) > 0 {
		fmt.Printf("✅ PROPFIND响应长度: %d bytes\n", len(body))
		// 检查是否包含XML响应
		if strings.Contains(string(body), "<?xml") {
			fmt.Println("✅ 收到有效的XML响应")
		}
	}
}

func testOPTIONSNoAuth(client *http.Client, baseURL string) {
	req, err := http.NewRequest("OPTIONS", baseURL, nil)
	if err != nil {
		fmt.Printf("❌ 创建OPTIONS请求失败: %v\n", err)
		return
	}

	req.Header.Set("User-Agent", "Microsoft-WebDAV-MiniRedir/10.0.19041")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ OPTIONS请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("✅ 无认证OPTIONS响应状态: %s\n", resp.Status)

	if resp.StatusCode == 401 {
		fmt.Println("✅ 正确返回401认证要求")
		checkHeader(resp, "WWW-Authenticate")
	}
}

func testWrongAuth(client *http.Client, baseURL string) {
	req, err := http.NewRequest("OPTIONS", baseURL, nil)
	if err != nil {
		fmt.Printf("❌ 创建OPTIONS请求失败: %v\n", err)
		return
	}

	// 使用错误的认证信息
	auth := base64.StdEncoding.EncodeToString([]byte("wrong:wrong"))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("User-Agent", "Microsoft-WebDAV-MiniRedir/10.0.19041")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ OPTIONS请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("✅ 错误认证响应状态: %s\n", resp.Status)

	if resp.StatusCode == 401 {
		fmt.Println("✅ 正确拒绝错误认证")
	}
}

func checkHeader(resp *http.Response, headerName string) {
	value := resp.Header.Get(headerName)
	if value != "" {
		fmt.Printf("✅ %s: %s\n", headerName, value)
	} else {
		fmt.Printf("❌ 缺少头部: %s\n", headerName)
	}
}
