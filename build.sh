#!/bin/bash

# 检查是否安装了 poetry
if ! command -v poetry &> /dev/null; then
    echo "正在安装 Poetry..."
    curl -sSL https://install.python-poetry.org | python3 -
fi

# 确保使用系统 Python
poetry config virtualenvs.create true
poetry config virtualenvs.in-project true

# 安装依赖
echo "安装项目依赖..."
poetry install --no-root

# 使用 Poetry 环境运行 PyInstaller
echo "开始构建二进制文件..."
poetry run pyinstaller \
    --onefile \
    --name lottery7 \
    --hidden-import mysql.connector \
    --hidden-import mysql.connector.locales \
    --hidden-import mysql.connector.locales.eng.client_error \
    lottery_scraper.py

echo "构建完成！可执行文件位于 dist/lottery7" 