@echo off
REM 构建脚本 - Windows 入口

setlocal enabledelayedexpansion

set SCRIPT_DIR=%~dp0
set PROJECT_ROOT=%SCRIPT_DIR%..

cd /d "%PROJECT_ROOT%"

echo === GitHub Issue Pack 构建 ===

REM 读取版本（简化处理）
for /f "tokens=2 delims=:" %%a in ('findstr /c:"version" version.json') do (
    set VERSION=%%a
    set VERSION=!VERSION:"=!
    set VERSION=!VERSION: =!
    set VERSION=!VERSION:,=!
    goto :found_version
)
:found_version

echo 版本: %VERSION%

REM 创建输出目录
if not exist bin mkdir bin

echo 构建 windows/amd64...
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-s -w" -o bin\github-issue-windows-amd64.exe .\cmd\github-issue

echo 构建 darwin/amd64...
set GOOS=darwin
set GOARCH=amd64
go build -ldflags="-s -w" -o bin\github-issue-darwin-amd64 .\cmd\github-issue

echo 构建 darwin/arm64...
set GOOS=darwin
set GOARCH=arm64
go build -ldflags="-s -w" -o bin\github-issue-darwin-arm64 .\cmd\github-issue

echo 构建 linux/amd64...
set GOOS=linux
set GOARCH=amd64
go build -ldflags="-s -w" -o bin\github-issue-linux-amd64 .\cmd\github-issue

echo 构建 linux/arm64...
set GOOS=linux
set GOARCH=arm64
go build -ldflags="-s -w" -o bin\github-issue-linux-arm64 .\cmd\github-issue

REM 复制为默认名称
copy /y bin\github-issue-windows-amd64.exe bin\github-issue.exe >nul

echo.
echo ✅ 构建完成!
dir bin\

endlocal
