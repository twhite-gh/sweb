@echo off
chcp 65001 >nul 2>&1
cls
echo ========================================
echo SWeb Complete Test Suite
echo ========================================
echo.

echo Please ensure the server is started with all features enabled:
echo ./sweb.exe -upload -files -webdav -https -socks5 -proxy
echo.

echo Starting comprehensive tests...
echo.

:menu
echo ========================================
echo Test Menu
echo ========================================
echo 1. Run all tests
echo 2. Basic functionality test
echo 3. SOCKS5 proxy test
echo 4. HTTP/HTTPS proxy test
echo 5. Debug connection test
echo 6. Quick SOCKS5 test
echo 0. Exit
echo ========================================
set /p choice="Please select an option (0-6): "

if "%choice%"=="0" goto :end
if "%choice%"=="1" goto :all_tests
if "%choice%"=="2" goto :basic_test
if "%choice%"=="3" goto :socks5_test
if "%choice%"=="4" goto :proxy_test
if "%choice%"=="5" goto :debug_test
if "%choice%"=="6" goto :quick_test

echo Invalid option
goto :menu

:all_tests
echo ========================================
echo Running All Tests
echo ========================================
call :basic_test
call :socks5_test
call :proxy_test
call :debug_test
goto :menu

:basic_test
echo.
echo ========================================
echo 1. Basic Functionality Test
echo ========================================
echo.
echo Testing server status and basic features...
go run test_all.go
echo.
pause
goto :menu

:socks5_test
echo.
echo ========================================
echo 2. SOCKS5 Proxy Test
echo ========================================
echo.
echo Compiling SOCKS5 test...
go build -o test_socks5_temp.exe test_socks5.go
if %errorlevel% neq 0 (
    echo X SOCKS5 test compilation failed
    pause
    goto :menu
)

echo Running SOCKS5 proxy test...
test_socks5_temp.exe
if exist test_socks5_temp.exe del test_socks5_temp.exe
echo.
pause
goto :menu

:proxy_test
echo.
echo ========================================
echo 3. HTTP/HTTPS Proxy Test
echo ========================================
echo.
echo Compiling HTTP proxy test...
go build -o test_proxy_temp.exe test_proxy.go
if %errorlevel% neq 0 (
    echo X HTTP proxy test compilation failed
    pause
    goto :menu
)

echo Running HTTP/HTTPS proxy test...
test_proxy_temp.exe
if exist test_proxy_temp.exe del test_proxy_temp.exe
echo.
pause
goto :menu

:debug_test
echo.
echo ========================================
echo 4. Debug Connection Test
echo ========================================
echo.
echo Compiling debug test...
go build -o debug_socks5_temp.exe debug_socks5.go
if %errorlevel% neq 0 (
    echo X Debug test compilation failed
    pause
    goto :menu
)

echo Running debug connection test...
debug_socks5_temp.exe
if exist debug_socks5_temp.exe del debug_socks5_temp.exe
echo.
pause
goto :menu

:quick_test
echo.
echo ========================================
echo 5. Quick SOCKS5 Test
echo ========================================
echo.
echo Compiling main program...
cd ..
go build -o sweb_temp.exe main.go upload.go files.go webdav.go server.go utils.go socks5.go proxy.go
if %errorlevel% neq 0 (
    echo X Main program compilation failed
    cd test
    pause
    goto :menu
)
cd test

echo Compilation successful
echo.
echo Please run in another command window: ../sweb_temp.exe -socks5
echo Then press any key to continue testing...
pause

echo.
echo Testing SOCKS5 connection...
go run debug_socks5.go

if exist ..\sweb_temp.exe del ..\sweb_temp.exe
echo.
pause
goto :menu

:end
echo.
echo ========================================
echo Test Suite Completed!
echo ========================================
echo.
echo Summary:
echo - Server status: http://localhost:8080/api/status
echo - SOCKS5 proxy: localhost:1080
echo - HTTP proxy: localhost:10808
echo.
echo Browser configuration tips:
echo 1. For SOCKS5: Set SOCKS5 proxy to 127.0.0.1:1080
echo 2. For HTTP/HTTPS: Set HTTP proxy to 127.0.0.1:10808
echo.
echo Thank you for using SWeb!
echo.
pause
