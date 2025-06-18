@echo off
SETLOCAL

:: Define project name and source files
SET PROJECT_NAME=sweb
SET SOURCE_FILES=./main.go ./upload.go ./files.go ./webdav.go ./server.go ./utils.go
:: NOTE: All Go files in the main package need to be included for compilation

SET BUILD_DIR=.\bin

echo Building for Windows...
IF NOT EXIST %BUILD_DIR% MD %BUILD_DIR%
IF NOT EXIST %BUILD_DIR%\windows MD %BUILD_DIR%\windows
SET GOOS=windows
SET GOARCH=amd64
go build -ldflags "-s -w" -o %BUILD_DIR%\windows\%PROJECT_NAME%-%GOOS%-%GOARCH%.exe %SOURCE_FILES%
IF %ERRORLEVEL% NEQ 0 (
    echo Windows build failed!
    GOTO :EOF
)
echo Windows build complete: %BUILD_DIR%\windows\%PROJECT_NAME%.exe

echo.
echo Building for Linux...
IF NOT EXIST %BUILD_DIR%\linux MD %BUILD_DIR%\linux
SET GOOS=linux
SET GOARCH=amd64
go build -ldflags "-s -w" -o %BUILD_DIR%\linux\%PROJECT_NAME%-%GOOS%-%GOARCH% %SOURCE_FILES%
IF %ERRORLEVEL% NEQ 0 (
    echo Linux build failed!
    GOTO :EOF
)
echo Linux build complete: %BUILD_DIR%\linux\%PROJECT_NAME%

echo.
echo All builds complete!

ENDLOCAL