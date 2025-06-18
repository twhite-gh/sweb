package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// 创建证书目录
	certDir := "./cert"
	if err := os.MkdirAll(certDir, 0755); err != nil {
		fmt.Printf("创建证书目录失败: %v\n", err)
		return
	}

	// 生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Printf("生成私钥失败: %v\n", err)
		return
	}

	// 创建证书模板
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization:  []string{"Test Company"},
			Country:       []string{"CN"},
			Province:      []string{"Beijing"},
			Locality:      []string{"Beijing"},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour), // 1年有效期
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		DNSNames:     []string{"localhost"},
	}

	// 生成证书
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		fmt.Printf("生成证书失败: %v\n", err)
		return
	}

	// 保存证书文件
	certOut, err := os.Create(filepath.Join(certDir, "server.crt"))
	if err != nil {
		fmt.Printf("创建证书文件失败: %v\n", err)
		return
	}
	defer certOut.Close()

	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certDER}); err != nil {
		fmt.Printf("写入证书文件失败: %v\n", err)
		return
	}

	// 保存私钥文件
	keyOut, err := os.Create(filepath.Join(certDir, "server.key"))
	if err != nil {
		fmt.Printf("创建私钥文件失败: %v\n", err)
		return
	}
	defer keyOut.Close()

	privateKeyDER, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		fmt.Printf("序列化私钥失败: %v\n", err)
		return
	}

	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privateKeyDER}); err != nil {
		fmt.Printf("写入私钥文件失败: %v\n", err)
		return
	}

	fmt.Println("✅ SSL证书生成成功!")
	fmt.Printf("📁 证书目录: %s\n", certDir)
	fmt.Printf("📄 证书文件: %s\n", filepath.Join(certDir, "server.crt"))
	fmt.Printf("🔑 私钥文件: %s\n", filepath.Join(certDir, "server.key"))
	fmt.Println("⚠️  这是自签名证书，仅用于测试目的")
	fmt.Println("🌐 现在可以使用 -https 参数启动HTTPS服务")
}
