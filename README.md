# Nowcoder CLI

牛客网面经爬虫命令行工具 / A CLI tool for crawling interview experiences from Nowcoder.

## 功能特点 / Features

- 🔍 根据公司和岗位搜索面经 / Search interview experiences by company and position
- 📄 自动翻页抓取帖子 / Auto-paginate through posts
- 📝 获取帖子详细内容 / Fetch detailed post content
- 🔄 自动去重 / Automatic deduplication
- 💾 结果保存到文件 / Save results to file

## 使用方法 / Usage

```bash
# 搜索腾讯后端面经，抓取2页
./nowcoder interviewCmd -c 腾讯 -p 后端 -l 2

# Search ByteDance backend interviews, 3 pages
./nowcoder interviewCmd -c 字节跳动 -p 后端 -l 3
```

## 参数说明 / Parameters

| 参数 / Param | 缩写 / Short | 必填 / Required | 说明 / Description |
|-------------|--------------|-----------------|---------------------|
| --company   | -c           | 是 / Yes        | 目标公司 / Target company |
| --position  | -p           | 否 / No         | 目标岗位 / Target position |
| --limit     | -l           | 否 / No         | 抓取页数 / Number of pages (默认/Default: 5) |

## 前提条件 / Prerequisites

- Chrome/Chromium 浏览器 / Chrome/Chromium browser
- 牛客网 Cookie 保存到 `nowcoder_cookie.json`（与二进制文件同目录）/ Save Nowcoder cookies to `nowcoder_cookie.json` (same directory as the binary)

### 获取 Cookie 方法 / How to Get Cookies

1. 在浏览器中登录 [牛客网](https://www.nowcoder.com)
2. 按 F12 打开开发者工具 → Application（应用）→ Cookies → `https://www.nowcoder.com`
3. 右键表格 → Export（导出），得到 JSON 文件
4. 将 JSON 文件重命名为 `nowcoder_cookie.json` 放到二进制文件同目录

### Cookie 文件格式 / Cookie File Format

从 Chrome 导出的格式需要包含以下字段 / The exported cookie JSON needs to include these fields:

```json
[
  {
    "domain": ".nowcoder.com",
    "expirationDate": 1748000000,
    "hostOnly": false,
    "httpOnly": true,
    "name": "NOWCODERUID",
    "path": "/",
    "sameSite": "lax",
    "secure": true,
    "session": false,
    "storeId": "0",
    "value": "your_cookie_value_here"
  },
  {
    "domain": ".nowcoder.com",
    "name": "NOWCODERCOOKIE",
    "value": "another_cookie_value"
  }
]
```

**必需字段 / Required fields:** `domain`, `name`, `value`
**可选字段 / Optional fields:** `path`, `secure`, `httpOnly`, `sameSite`, `expirationDate`, `session`

## 构建 / Build

```bash
# macOS
go build -o nowcoder ./cmd/cli/main.go

# Linux (x86_64)
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o nowcoder ./cmd/cli/main.go
```

## 文件结构 / Project Structure

```
.
├── cmd/cli/main.go           # CLI 入口 / CLI entry point
├── service/nowcoder/          # 核心爬虫逻辑 / Core crawler logic
├── internal/commands/        # 命令行命令 / CLI commands
├── docs/usage.md             # 详细使用文档 / Detailed usage documentation
└── skills/                   # 技能文档 / Skills documentation
```
