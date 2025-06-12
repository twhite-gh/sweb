@echo off
SETLOCAL

:: Define project name and main file
SET PROJECT_NAME=sweb
SET MAIN_FILE=./main.go
:: NOTE: If your main.go is in a subdirectory like 'cmd/sweb/main.go',
::       you'd change the above to: SET MAIN_FILE=./cmd/sweb/main.go

SET BUILD_DIR=.\bin

echo Building for Windows...
IF NOT EXIST %BUILD_DIR% MD %BUILD_DIR%
IF NOT EXIST %BUILD_DIR%\windows MD %BUILD_DIR%\windows
SET GOOS=windows
SET GOARCH=amd64
go build -o %BUILD_DIR%\windows\%PROJECT_NAME%.exe %MAIN_FILE%
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
go build -o %BUILD_DIR%\linux\%PROJECT_NAME% %MAIN_FILE%
IF %ERRORLEVEL% NEQ 0 (
    echo Linux build failed!
    GOTO :EOF
)
echo Linux build complete: %BUILD_DIR%\linux\%PROJECT_NAME%

echo.
echo All builds complete!

ENDLOCAL