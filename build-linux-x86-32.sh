#!/bin/bash

# MiniJump Linux x86 (32位) 构建脚本
# 用于交叉编译到 Linux 386 架构

set -e

echo "开始构建 MiniJump for Linux x86 (32位)..."

# 设置交叉编译环境变量
export GOOS=linux
export GOARCH=386

# 构建参数
APP_NAME="minijump"
VERSION=$(date +"%Y%m%d-%H%M%S")
BUILD_DIR="build"
OUTPUT_FILE="${BUILD_DIR}/minijump-linux-x86"

# 创建构建目录
mkdir -p ${BUILD_DIR}

# 构建
echo "正在编译..."
go build -ldflags "-s -w" -o ${OUTPUT_FILE} main.go

# 检查构建结果
if [ -f ${OUTPUT_FILE} ]; then
    echo "✓ 构建成功！"
    echo "  输出文件: ${OUTPUT_FILE}"
    
    # 显示文件信息
    file ${OUTPUT_FILE}
    ls -lh ${OUTPUT_FILE}
    
    echo ""
    echo "使用方法："
    echo "  chmod +x ${OUTPUT_FILE}"
    echo "  ./${OUTPUT_FILE} -port 8080"
else
    echo "✗ 构建失败！"
    exit 1
fi
