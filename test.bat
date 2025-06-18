@echo off
echo ========================================
echo SWeb 代理功能测试脚本
echo ========================================
echo.

echo 请确保服务器已启动并启用了代理功能：
echo ./sweb.exe -socks5 -proxy
echo.

echo 1. 测试综合功能状态...
go run test_all.go
echo.

echo 2. 测试SOCKS5代理功能...
go run test_socks5.go
echo.

echo 3. 测试HTTP代理功能...
go run test_proxy.go
echo.

echo ========================================
echo 测试完成！
echo ========================================
pause
