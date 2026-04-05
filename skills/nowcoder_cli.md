# Nowcoder CLI 使用指南

## 使用方法

```bash
./nowcoder interviewCmd -c <公司> -p <岗位> -l <页数>
```

## 参数说明

| 参数 | 缩写 | 必填 | 说明 | 默认值 |
|------|------|----|------|--------|
| --company | -c | 是  | 目标公司名称 | - |
| --position | -p | 是  | 目标岗位名称 | - |
| --limit | -l | 是  | 抓取页数 | 5 |

## 示例

```bash
# 搜索腾讯后端面经，抓取2页
./nowcoder interviewCmd -c 腾讯 -p 后端 -l 2

# 搜索字节跳动后端面经，抓取3页
./nowcoder interviewCmd -c 字节跳动 -p 后端 -l 3
```

## 前提条件

- 已获取牛客网 Cookie 并保存到 `cmd/tmp/nowcoder_cookie.json`
