# MiniJump 系统服务安装指南

本文档介绍如何将 MiniJump 安装为系统服务。

## Windows 服务安装

### 前置要求
- Windows 操作系统
- 管理员权限

### 安装步骤

1. **以管理员身份打开命令提示符或 PowerShell**
   - 右键点击命令提示符或 PowerShell
   - 选择"以管理员身份运行"

2. **切换到 MiniJump 程序目录**

3. **安装服务**
   ```powershell
   # 使用默认配置（端口 8080）
   .\minijump.exe install
   
   # 自定义配置
   .\minijump.exe install -port 8080 -config "C:\MiniJump\rules.json" -log "C:\MiniJump\access.log"
   ```

4. **启动服务**
   ```powershell
   sc start MiniJump
   ```

5. **查看服务状态**
   ```powershell
   sc query MiniJump
   ```

6. **停止服务**
   ```powershell
   sc stop MiniJump
   ```

7. **卸载服务**
   ```powershell
   .\minijump.exe uninstall
   ```

### 服务管理命令

| 操作 | 命令 |
|------|------|
| 启动服务 | `sc start MiniJump` |
| 停止服务 | `sc stop MiniJump` |
| 查询状态 | `sc query MiniJump` |
| 删除服务 | `sc delete MiniJump` 或 `.\minijump.exe uninstall` |

### Windows 服务管理器

你也可以使用 Windows 服务管理器（services.msc）来管理服务：
1. 按 `Win + R`，输入 `services.msc`
2. 找到 "MiniJump HTTP Redirect Service"
3. 右键可以启动、停止、重启服务

## Linux 服务安装（systemd）

### 前置要求
- Linux 操作系统（使用 systemd）
- root 权限（使用 sudo）

### 安装步骤

1. **使用 root 权限运行安装命令**
   ```bash
   # 使用默认配置
   sudo ./minijump install
   
   # 自定义配置
   sudo ./minijump install -port 8080 -config /etc/minijump/rules.json -log /var/log/minijump/access.log
   ```

2. **启动服务**
   ```bash
   sudo systemctl start MiniJump
   ```

3. **查看服务状态**
   ```bash
   sudo systemctl status MiniJump
   ```

4. **停止服务**
   ```bash
   sudo systemctl stop MiniJump
   ```

5. **查看服务日志**
   ```bash
   sudo journalctl -u MiniJump -f
   ```

6. **卸载服务**
   ```bash
   sudo ./minijump uninstall
   ```

### 服务管理命令

| 操作 | 命令 |
|------|------|
| 启动服务 | `sudo systemctl start MiniJump` |
| 停止服务 | `sudo systemctl stop MiniJump` |
| 重启服务 | `sudo systemctl restart MiniJump` |
| 查看状态 | `sudo systemctl status MiniJump` |
| 查看日志 | `sudo journalctl -u MiniJump -f` |
| 开机自启 | `sudo systemctl enable MiniJump` |
| 禁用自启 | `sudo systemctl disable MiniJump` |
| 重新加载配置 | `sudo systemctl daemon-reload` |

## 安装参数说明

安装服务时可以指定以下参数：

- `-name`: 服务名称（默认：MiniJump）
- `-port`: 服务端口（默认：8080）
- `-config`: 配置文件路径（默认：rules.json，相对于可执行文件目录）
- `-log`: 日志文件路径（默认：access.log，相对于可执行文件目录）
- `-log-buffer`: 日志缓冲大小（默认：1000）
- `-log-flush`: 日志刷新间隔秒数（默认：180）

### 示例

```bash
# Windows - 自定义端口和配置路径
.\minijump.exe install -port 9090 -config "D:\Config\rules.json" -log "D:\Logs\access.log"

# Linux - 自定义配置
sudo ./minijump install -port 9090 -config /etc/minijump/rules.json -log /var/log/minijump/access.log
```

## 常见问题

### Windows

**Q: 安装时提示"需要管理员权限"**
A: 请以管理员身份运行命令提示符或 PowerShell。

**Q: 服务安装成功但无法启动**
A: 检查：
1. 可执行文件路径是否正确
2. 配置文件路径是否存在
3. 端口是否被占用
4. 查看事件查看器（eventvwr.msc）中的错误日志

**Q: 如何修改服务配置**
A: 卸载服务后重新安装，使用新的参数。

### Linux

**Q: 提示"未找到 systemctl"**
A: 确保系统使用 systemd（大多数现代 Linux 发行版都使用 systemd）。

**Q: 服务安装成功但无法启动**
A: 检查：
1. 使用 `sudo systemctl status MiniJump` 查看错误信息
2. 使用 `sudo journalctl -u MiniJump` 查看日志
3. 检查可执行文件权限：`chmod +x minijump`
4. 检查配置文件路径和权限

**Q: 如何修改服务配置**
A: 编辑 `/etc/systemd/system/MiniJump.service` 文件，然后运行 `sudo systemctl daemon-reload` 和 `sudo systemctl restart MiniJump`。

## 注意事项

1. **路径问题**：配置文件路径建议使用绝对路径，避免服务运行时找不到文件。

2. **权限问题**：
   - Windows：必须以管理员身份运行安装命令
   - Linux：必须使用 sudo 或 root 权限

3. **服务名称**：默认服务名称为 "MiniJump"，可以在安装时使用 `-name` 参数自定义。

4. **日志文件**：建议将日志文件放在专门的目录中，如 Windows 的 `C:\Logs\` 或 Linux 的 `/var/log/minijump/`。

5. **自动重启**：
   - Windows：服务配置为自动启动（start= auto）
   - Linux：服务配置了 Restart=always，服务异常退出时会自动重启
