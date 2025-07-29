package main

import (
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// testWebDAVAuth 测试WebDAV认证功能
func testWebDAVAuth() {
	fmt.Println("🔐 开始测试WebDAV认证功能...")

	// 测试无认证访问（应该失败）
	fmt.Println("1. 测试无认证访问WebDAV...")
	resp, err := http.Get("http://localhost:8080/webdav")
	if err != nil {
		fmt.Printf("❌ 请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		fmt.Println("✅ 无认证访问被正确拒绝 (401 Unauthorized)")
	} else {
		fmt.Printf("❌ 期望401状态码，但得到: %d\n", resp.StatusCode)
		return
	}

	// 测试错误认证（应该失败）
	fmt.Println("2. 测试错误认证...")
	req, _ := http.NewRequest("GET", "http://localhost:8080/webdav", nil)
	auth := base64.StdEncoding.EncodeToString([]byte("wrong:wrong"))
	req.Header.Set("Authorization", "Basic "+auth)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("❌ 请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		fmt.Println("✅ 错误认证被正确拒绝 (401 Unauthorized)")
	} else {
		fmt.Printf("❌ 期望401状态码，但得到: %d\n", resp.StatusCode)
		return
	}

	// 测试正确认证（应该成功）
	fmt.Println("3. 测试正确认证...")
	req, _ = http.NewRequest("GET", "http://localhost:8080/webdav", nil)
	auth = base64.StdEncoding.EncodeToString([]byte("webdav:webdav"))
	req.Header.Set("Authorization", "Basic "+auth)

	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("❌ 请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 || resp.StatusCode == 207 {
		fmt.Println("✅ 正确认证访问成功")
	} else {
		fmt.Printf("❌ 期望200/207状态码，但得到: %d\n", resp.StatusCode)
		return
	}

	fmt.Println("✅ WebDAV认证功能测试通过")
}

// testProxyAuth 测试HTTP代理认证功能
func testProxyAuth() {
	fmt.Println("🔐 开始测试HTTP代理认证功能...")

	// 测试无认证访问（应该失败）
	fmt.Println("1. 测试无认证代理访问...")
	proxyURL, _ := url.Parse("http://localhost:10808")
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get("http://httpbin.org/ip")
	if err != nil {
		if strings.Contains(err.Error(), "407") {
			fmt.Println("✅ 无认证代理访问被正确拒绝 (407 Proxy Authentication Required)")
		} else {
			fmt.Printf("❌ 期望407错误，但得到: %v\n", err)
			return
		}
	} else {
		defer resp.Body.Close()
		fmt.Printf("❌ 期望认证失败，但请求成功，状态码: %d\n", resp.StatusCode)
		return
	}

	// 测试正确认证（应该成功）
	fmt.Println("2. 测试正确代理认证...")
	proxyURLWithAuth, _ := url.Parse("http://http:http@localhost:10808")

	client = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURLWithAuth),
		},
		Timeout: 10 * time.Second,
	}

	resp, err = client.Get("http://www.baidu.com")
	if err != nil {
		fmt.Printf("❌ 代理请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Println("✅ 正确认证代理访问成功")
	} else {
		fmt.Printf("❌ 期望200状态码，但得到: %d\n", resp.StatusCode)
		return
	}

	fmt.Println("✅ HTTP代理认证功能测试通过")
}

// testSOCKS5Auth 测试SOCKS5代理认证功能
func testSOCKS5Auth() {
	fmt.Println("🔐 开始测试SOCKS5代理认证功能...")
	fmt.Println("注意: SOCKS5认证测试需要专门的客户端库，这里只做基本连接测试")

	// 测试连接到SOCKS5代理
	fmt.Println("1. 测试连接到SOCKS5代理...")
	conn, err := net.DialTimeout("tcp", "localhost:1080", 5*time.Second)
	if err != nil {
		fmt.Printf("❌ 无法连接到SOCKS5代理: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Println("✅ 成功连接到SOCKS5代理服务器")

	// 发送SOCKS5认证协商请求（支持用户名密码认证）
	fmt.Println("2. 测试SOCKS5认证协商...")
	authRequest := []byte{0x05, 0x02, 0x00, 0x02} // 版本5，2种方法：无认证和用户名密码
	_, err = conn.Write(authRequest)
	if err != nil {
		fmt.Printf("❌ 发送认证请求失败: %v\n", err)
		return
	}

	// 读取服务器响应
	response := make([]byte, 2)
	_, err = conn.Read(response)
	if err != nil {
		fmt.Printf("❌ 读取认证响应失败: %v\n", err)
		return
	}

	if response[0] != 0x05 {
		fmt.Printf("❌ 无效的SOCKS5版本响应: %d\n", response[0])
		return
	}

	if response[1] == 0x02 {
		fmt.Println("✅ 服务器要求用户名密码认证")

		// 发送用户名密码
		username := "socks5"
		password := "socks5"
		authData := []byte{0x01} // 认证版本
		authData = append(authData, byte(len(username)))
		authData = append(authData, []byte(username)...)
		authData = append(authData, byte(len(password)))
		authData = append(authData, []byte(password)...)

		_, err = conn.Write(authData)
		if err != nil {
			fmt.Printf("❌ 发送认证数据失败: %v\n", err)
			return
		}

		// 读取认证结果
		authResult := make([]byte, 2)
		_, err = conn.Read(authResult)
		if err != nil {
			fmt.Printf("❌ 读取认证结果失败: %v\n", err)
			return
		}

		if authResult[1] == 0x00 {
			fmt.Println("✅ SOCKS5用户名密码认证成功")
		} else {
			fmt.Printf("❌ SOCKS5用户名密码认证失败: %d\n", authResult[1])
			return
		}
	} else if response[1] == 0x00 {
		fmt.Println("✅ 服务器接受无认证连接")
	} else {
		fmt.Printf("❌ 服务器拒绝认证方法: %d\n", response[1])
		return
	}

	fmt.Println("✅ SOCKS5代理认证功能测试通过")
}

func main() {
	fmt.Println("🧪 开始认证功能测试")
	fmt.Println("请确保服务器已启动:")
	fmt.Println("  WebDAV认证: ./sweb.exe -webdav -webdav-auth")
	fmt.Println("  HTTP代理认证: ./sweb.exe -proxy -proxy-auth")
	fmt.Println("  SOCKS5代理认证: ./sweb.exe -socks5 -socks5-auth")
	fmt.Println()

	// 等待用户确认
	fmt.Print("按回车键开始测试...")
	fmt.Scanln()

	// 测试WebDAV认证
	testWebDAVAuth()
	fmt.Println()

	// 测试HTTP代理认证
	testProxyAuth()
	fmt.Println()

	// 测试SOCKS5代理认证
	testSOCKS5Auth()
	fmt.Println()

	fmt.Println("🎉 所有认证功能测试完成")
}
