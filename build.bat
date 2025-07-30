@echo off
SETLOCAL

:: Define project name and source files
SET PROJECT_NAME=sweb
SET SOURCE_FILES=./src/main.go ./src/upload.go ./src/files.go ./src/webdav.go ./src/server.go ./src/utils.go ./src/socks5.go ./src/proxy.go ./src/https_proxy.go
:: NOTE: All Go files in the main package need to be included for compilation

SET BUILD_DIR=.\bin

echo Building for Windows (64-bit)...
IF NOT EXIST %BUILD_DIR% MD %BUILD_DIR%
IF NOT EXIST %BUILD_DIR%\windows MD %BUILD_DIR%\windows
SET GOOS=windows
SET GOARCH=amd64
go build -ldflags "-s -w" -o %BUILD_DIR%\windows\%PROJECT_NAME%-%GOOS%-%GOARCH%.exe %SOURCE_FILES%
IF %ERRORLEVEL% NEQ 0 (
    echo Windows 64-bit build failed!
    GOTO :EOF
)
echo Windows 64-bit build complete: %BUILD_DIR%\windows\%PROJECT_NAME%-%GOOS%-%GOARCH%.exe

echo.
echo Building for Windows (32-bit)...
SET GOOS=windows
SET GOARCH=386
go build -ldflags "-s -w" -o %BUILD_DIR%\windows\%PROJECT_NAME%-%GOOS%-%GOARCH%.exe %SOURCE_FILES%
IF %ERRORLEVEL% NEQ 0 (
    echo Windows 32-bit build failed!
    GOTO :EOF
)
echo Windows 32-bit build complete: %BUILD_DIR%\windows\%PROJECT_NAME%-%GOOS%-%GOARCH%.exe

echo.
echo Building for Linux (64-bit)...
IF NOT EXIST %BUILD_DIR%\linux MD %BUILD_DIR%\linux
SET GOOS=linux
SET GOARCH=amd64
go build -ldflags "-s -w" -o %BUILD_DIR%\linux\%PROJECT_NAME%-%GOOS%-%GOARCH% %SOURCE_FILES%
IF %ERRORLEVEL% NEQ 0 (
    echo Linux 64-bit build failed!
    GOTO :EOF
)
echo Linux 64-bit build complete: %BUILD_DIR%\linux\%PROJECT_NAME%-%GOOS%-%GOARCH%

echo.
echo Building for Linux (32-bit)...
SET GOOS=linux
SET GOARCH=386
go build -ldflags "-s -w" -o %BUILD_DIR%\linux\%PROJECT_NAME%-%GOOS%-%GOARCH% %SOURCE_FILES%
IF %ERRORLEVEL% NEQ 0 (
    echo Linux 32-bit build failed!
    GOTO :EOF
)
echo Linux 32-bit build complete: %BUILD_DIR%\linux\%PROJECT_NAME%-%GOOS%-%GOARCH%

echo.
echo Building for macOS (Intel)...
IF NOT EXIST %BUILD_DIR%\macos MD %BUILD_DIR%\macos
SET GOOS=darwin
SET GOARCH=amd64
go build -ldflags "-s -w" -o %BUILD_DIR%\macos\%PROJECT_NAME%-%GOOS%-%GOARCH% %SOURCE_FILES%
IF %ERRORLEVEL% NEQ 0 (
    echo macOS Intel build failed!
    GOTO :EOF
)
echo macOS Intel build complete: %BUILD_DIR%\macos\%PROJECT_NAME%-%GOOS%-%GOARCH%

echo.
echo Building for macOS (Apple Silicon)...
SET GOOS=darwin
SET GOARCH=arm64
go build -ldflags "-s -w" -o %BUILD_DIR%\macos\%PROJECT_NAME%-%GOOS%-%GOARCH% %SOURCE_FILES%
IF %ERRORLEVEL% NEQ 0 (
    echo macOS Apple Silicon build failed!
    GOTO :EOF
)
echo macOS Apple Silicon build complete: %BUILD_DIR%\macos\%PROJECT_NAME%-%GOOS%-%GOARCH%

echo.
echo All builds complete!

ENDLOCAL