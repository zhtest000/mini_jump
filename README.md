# MiniJump HTTP 跳转服务

MiniJump 是一个轻量级的 HTTP 跳转服务，支持基于域名和路径的智能跳转配置，具备高性能和易管理特性。

## 功能特性

### 1. 跳转规则配置
- **域名跳转**：根据请求域名跳转到指定目标
- **路径跳转**：支持域名+子路径的精确匹配跳转
- **跳转方式**：支持 301、302、307、JavaScript 跳转
- **有效期控制**：支持设置跳转规则的有效期

### 2. 数据管理
- **内存存储**：使用 sync.Map 存储跳转规则，保证高性能和线程安全
- **实时更新**：规则修改即时生效
- **持久化**：支持配置持久化到文件

### 3. 日志系统
- **访问日志**：记录 IP、User-Agent、跳转详情等信息
- **缓冲机制**：每 1000 条或 3 分钟批量落盘
- **性能优化**：减少磁盘 I/O 操作

## 快速开始

### 安装

```bash
go mod download
go build -o minijump main.go
```

### 交叉编译

项目提供了多个构建脚本用于不同平台：

#### Linux x86_64（64位）
```bash
chmod +x build-linux-x86.sh
./build-linux-x86.sh
```

#### Linux x86（32位）
```bash
chmod +x build-linux-x86-32.sh
./build-linux-x86-32.sh
```

#### 构建所有平台
```bash
chmod +x build-all.sh
./build-all.sh
```

构建后的文件会输出到 `build/` 目录。

### 运行

```bash
./minijump -port 8080 -config rules.json -log access.log
```

### 命令行参数

- `-port`: 服务端口（默认：8080）
- `-config`: 配置文件路径（默认：rules.json）
- `-log`: 日志文件路径（默认：access.log）
- `-log-buffer`: 日志缓冲大小（默认：1000）
- `-log-flush`: 日志刷新间隔秒数（默认：180）

### 系统服务安装

MiniJump 支持安装为系统服务，实现开机自启动和后台运行。

#### Windows 服务安装

1. **以管理员身份运行命令提示符或 PowerShell**

2. **安装服务**
```powershell
# 使用默认参数
.\minijump.exe install

# 指定端口和配置文件
.\minijump.exe install -port 8080 -config "C:\path\to\rules.json" -log "C:\path\to\access.log"
```

3. **启动服务**
```powershell
sc start MiniJump
```

4. **停止服务**
```powershell
sc stop MiniJump
```

5. **卸载服务**
```powershell
.\minijump.exe uninstall
```

#### Linux 服务安装（systemd）

1. **使用 root 权限运行**

2. **安装服务**
```bash
# 使用默认参数
sudo ./minijump install

# 指定参数
sudo ./minijump install -port 8080 -config /etc/minijump/rules.json -log /var/log/minijump/access.log
```

3. **启动服务**
```bash
sudo systemctl start MiniJump
```

4. **查看服务状态**
```bash
sudo systemctl status MiniJump
```

5. **设置开机自启**
```bash
# 服务安装时会自动启用开机自启
sudo systemctl enable MiniJump
```

6. **卸载服务**
```bash
sudo ./minijump uninstall
```

#### 服务安装参数

安装服务时可以使用的参数：
- `-name`: 服务名称（默认：MiniJump）
- `-port`: 服务端口（默认：8080）
- `-config`: 配置文件路径（默认：rules.json）
- `-log`: 日志文件路径（默认：access.log）
- `-log-buffer`: 日志缓冲大小（默认：1000）
- `-log-flush`: 日志刷新间隔秒数（默认：180）

**注意**：
- Windows 需要管理员权限
- Linux 需要 root 权限（使用 sudo）
- 服务安装后会自动使用指定的参数启动

## 管理页面

访问 `http://localhost:8080/manager558630` 打开 Web 管理界面。

管理页面功能：
- 📋 **规则列表**：查看所有跳转规则，包括域名、路径、目标URL、跳转类型等
- ➕ **添加规则**：通过表单创建新的跳转规则
- ✏️ **编辑规则**：修改现有规则的配置
- 🗑️ **删除规则**：删除不需要的规则
- 🔄 **重新加载**：从配置文件重新加载规则
- 💾 **保存配置**：将当前规则保存到配置文件
- ⚠️ **冲突检测**：自动检测并提示规则冲突

### 规则冲突检测

系统会自动检测以下冲突情况：
- **完全重复**：相同的域名和路径组合
- **域名覆盖**：域名级别规则会覆盖该域名的所有路径规则
- **路径被覆盖**：如果存在域名级别规则，路径规则将无法匹配

当检测到冲突时，系统会：
- 显示冲突提示信息
- 列出所有冲突的规则详情
- 阻止保存冲突的规则

## API 接口

### 1. 列出所有规则

```bash
GET /api/rules
```

### 2. 创建规则

```bash
POST /api/rules
Content-Type: application/json

{
  "domain": "example.com",
  "path": "/old",
  "target": "https://example.com/new",
  "type": 302,
  "expires_at": "2024-12-31T23:59:59Z",
  "description": "临时跳转"
}
```

### 3. 获取规则

```bash
GET /api/rules/{id}
```

### 4. 更新规则

```bash
PUT /api/rules/{id}
Content-Type: application/json

{
  "domain": "example.com",
  "path": "/old",
  "target": "https://example.com/new",
  "type": 301,
  "description": "永久跳转"
}
```

### 5. 删除规则

```bash
DELETE /api/rules/{id}
```

### 6. 重新加载配置

```bash
POST /api/reload
```

### 7. 保存配置

```bash
POST /api/save
```

## 跳转类型

- `301`: HTTP 301 永久重定向
- `302`: HTTP 302 临时重定向
- `307`: HTTP 307 临时重定向（保持请求方法）
- `4`: JavaScript 跳转

## 配置文件格式

配置文件为 JSON 格式，存储在 `rules.json`（可通过 `-config` 参数指定）：

```json
[
  {
    "id": "example_com",
    "domain": "example.com",
    "path": "",
    "target": "https://www.example.com",
    "type": 301,
    "expires_at": null,
    "created_at": "2024-01-01T00:00:00Z",
    "description": "域名跳转"
  },
  {
    "id": "example_com_old",
    "domain": "example.com",
    "path": "/old",
    "target": "https://example.com/new",
    "type": 302,
    "expires_at": "2024-12-31T23:59:59Z",
    "created_at": "2024-01-01T00:00:00Z",
    "description": "路径跳转"
  }
]
```

## 访问日志格式

访问日志为 JSON Lines 格式，每条记录一行：

```json
{"timestamp":"2024-01-01T12:00:00Z","ip":"127.0.0.1","user_agent":"Mozilla/5.0...","method":"GET","domain":"example.com","path":"/old","target":"https://example.com/new","redirect_type":302,"status_code":302}
```

## 规则匹配优先级

1. 精确匹配：域名 + 路径
2. 域名匹配：仅域名

## 项目结构

```
mini_jump/
├── main.go          # 主程序入口
├── config/          # 配置管理模块
│   └── config.go
├── handler/         # HTTP 请求处理
│   └── handler.go
├── logger/          # 日志管理
│   └── logger.go
├── api/             # RESTful API
│   └── api.go
├── manager/         # 管理页面
│   └── manager.go
├── service/         # 系统服务管理
│   └── service.go
├── go.mod           # Go 模块定义
└── README.md        # 项目文档
```

## 开发

```bash
# 运行测试
go run main.go

# 构建
go build -o minijump main.go

# 运行测试服务器
./minijump -port 8080
```

## 许可证

MIT License
