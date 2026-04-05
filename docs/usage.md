# Nowcoder CLI 使用文档

## 安装

```bash
go build -o nowcoder ./cmd/cli/main.go
```

## 使用方法

### 基本命令

```bash
./nowcoder interviewCmd -c <公司> -p <岗位> -l <页数>
```

### 参数说明

| 参数 | 缩写 | 必填 | 说明 | 默认值 |
|------|------|------|------|--------|
| --company | -c | 是 | 目标公司名称 | - |
| --position | -p | 否 | 目标岗位名称 | - |
| --limit | -l | 否 | 抓取页数 | 5 |

### 示例

```bash
# 搜索腾讯后端面经，抓取2页
./nowcoder interviewCmd -c 腾讯 -p 后端 -l 2

# 搜索字节跳动后端面经
./nowcoder interviewCmd -c 字节跳动 -p 后端 -l 3
```

## 功能说明

1. **搜索面经帖子** - 根据公司名称和岗位关键词搜索牛客网面经
2. **分页抓取** - 自动翻页抓取指定页数的帖子列表
3. **帖子详情** - 访问每个帖子链接，获取完整面经内容
4. **去重处理** - 自动去除重复帖子
5. **结果保存** - 帖子列表保存到 `cmd/tmp/post_links.json`

## 注意事项

- 需要提前获取牛客网 Cookie 并保存到 `cmd/tmp/nowcoder_cookie.json`
- 建议使用 headless=false 模式运行以便调试
- 抓取过程中会打开浏览器窗口
