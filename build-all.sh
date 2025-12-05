#!/bin/bash

# MiniJump 多平台构建脚本
# 构建 Windows、Linux (x86_64 和 32位) 可执行文件

set -e

echo "=========================================="
echo "  MiniJump 多平台构建脚本"
echo "=========================================="
echo ""

BUILD_DIR="build"
mkdir -p ${BUILD_DIR}

# Linux x86_64
echo "[1/3] 构建 Linux x86_64..."
export GOOS=linux
export GOARCH=amd64
go build -ldflags "-s -w" -o ${BUILD_DIR}/minijump-linux-x86_64 main.go
echo "✓ Linux x86_64 构建完成"
echo ""

# Linux x86 (32位)
echo "[2/3] 构建 Linux x86 (32位)..."
export GOOS=linux
export GOARCH=386
go build -ldflags "-s -w" -o ${BUILD_DIR}/minijump-linux-x86 main.go
echo "✓ Linux x86 (32位) 构建完成"
echo ""

# Windows
echo "[3/3] 构建 Windows x86_64..."
export GOOS=windows
export GOARCH=amd64
go build -ldflags "-s -w" -o ${BUILD_DIR}/minijump-windows-x86_64.exe main.go
echo "✓ Windows x86_64 构建完成"
echo ""

echo "=========================================="
echo "所有平台构建完成！"
echo "=========================================="
echo ""
echo "输出文件："
ls -lh ${BUILD_DIR}/
echo ""
echo "使用方法："
echo "  Linux:   chmod +x ${BUILD_DIR}/minijump-linux-x86_64 && ./${BUILD_DIR}/minijump-linux-x86_64"
echo "  Windows: ${BUILD_DIR}/minijump-windows-x86_64.exe"
