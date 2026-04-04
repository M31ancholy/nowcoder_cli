package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

type BrowserCookie struct {
	Domain         string  `json:"domain"`
	ExpirationDate float64 `json:"expirationDate"`
	HostOnly       bool    `json:"hostOnly"`
	HTTPOnly       bool    `json:"httpOnly"`
	Name           string  `json:"name"`
	Path           string  `json:"path"`
	SameSite       string  `json:"sameSite"`
	Secure         bool    `json:"secure"`
	Session        bool    `json:"session"`
	StoreId        string  `json:"storeId"`
	Value          string  `json:"value"`
}

// PostLink 表示一条面经帖子
type PostLink struct {
	Title string `json:"title"`
	Href  string `json:"href"`
	Desc  string `json:"desc"`
}

func loadCookiesFromFile(filename string) ([]BrowserCookie, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}
	var cookies []BrowserCookie
	if err := json.Unmarshal(data, &cookies); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %w", err)
	}
	return cookies, nil
}

func convertSameSite(sameSite string) network.CookieSameSite {
	switch sameSite {
	case "no_restriction":
		return network.CookieSameSiteNone
	case "lax":
		return network.CookieSameSiteLax
	case "strict":
		return network.CookieSameSiteStrict
	default:
		return ""
	}
}

func setCookies(cookies []BrowserCookie) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		for _, c := range cookies {
			cookieCmd := network.SetCookie(c.Name, c.Value).
				WithDomain(c.Domain).
				WithPath(c.Path).
				WithHTTPOnly(c.HTTPOnly).
				WithSecure(c.Secure)
			if ss := convertSameSite(c.SameSite); ss != "" {
				cookieCmd = cookieCmd.WithSameSite(ss)
			}
			if !c.Session && c.ExpirationDate > 0 {
				expr := cdp.TimeSinceEpoch(time.Unix(int64(c.ExpirationDate), 0))
				cookieCmd = cookieCmd.WithExpires(&expr)
			}
			_ = cookieCmd.Do(ctx)
		}
		return nil
	}
}

// extractPostLinks 从当前页面提取所有帖子链接（去重，只取标题链接）
func extractPostLinks() chromedp.ActionFunc {
	return func(ctx context.Context) error {
		return nil // 占位，实际提取在 Evaluate 里做
	}
}

// collectCurrentPage 用 JS 提取当前页面的帖子列表
func collectCurrentPage(result *string) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		return chromedp.Evaluate(`
			(() => {
				// 每张卡片的容器: div.tw-bg-white.tw-mt-3.tw-rounded-xl
				const cards = document.querySelectorAll('div.tw-bg-white.tw-mt-3.tw-rounded-xl');
				const posts = [];
				
				cards.forEach(card => {
					// 标题链接: 在 tw-font-bold tw-text-lg 的父 div 下
					const titleParent = card.querySelector('.tw-font-bold.tw-text-lg');
					const titleLink = titleParent ? titleParent.querySelector('a') : null;
					
					// 内容预览: tw-mb-2 下的 a 标签
					const descEl = card.querySelector('.tw-mb-2 a') || card.querySelector('.feed-text a');
					
					if (titleLink) {
						posts.push({
							title: titleLink.innerText.trim(),
							href: titleLink.getAttribute('href') || '',
							desc: descEl ? descEl.innerText.trim().substring(0, 200) : '',
						});
					}
				});
				
				return JSON.stringify(posts);
			})()
		`, result).Do(ctx)
	}
}

func main() {
	// ========== 配置 ==========
	query := "面经" // 搜索关键词
	limit := 3    // 遍历几页
	cookieFile := "./cmd/tmp/nowcoder_cookie.json"
	outputFile := "./cmd/tmp/post_links.json"
	// ===========================

	cookies, err := loadCookiesFromFile(cookieFile)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("共加载 %d 个 Cookie\n", len(cookies))

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// 总超时 = 每页约15秒
	timeout := time.Duration(limit*15+30) * time.Second
	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	searchURL := fmt.Sprintf(
		"https://www.nowcoder.com/search/all?query=%s&type=all&searchType=%%E9%%A1%%B6%%E9%%83%%A8%%E5%%AF%%BC%%E8%%88%%AA%%E6%%A0%%8F",
		query,
	)

	// 1. 打开页面 + 注入 Cookie + 刷新
	err = chromedp.Run(ctx,
		chromedp.Navigate(searchURL),
		chromedp.WaitReady("body"),
		setCookies(cookies),
		chromedp.Reload(),
		chromedp.WaitReady("body"),
		chromedp.Sleep(3*time.Second),
	)
	if err != nil {
		log.Fatalf("初始化页面失败: %v", err)
	}

	var allPosts []PostLink
	seen := make(map[string]bool) // 用 href 去重

	for page := 1; page <= limit; page++ {
		log.Printf("========== 正在抓取第 %d / %d 页 ==========", page, limit)

		// 等待卡片渲染
		err = chromedp.Run(ctx,
			chromedp.Sleep(2*time.Second),
		)
		if err != nil {
			log.Printf("第 %d 页等待失败: %v", page, err)
			break
		}

		// 提取当前页帖子
		var rawJSON string
		err = chromedp.Run(ctx,
			collectCurrentPage(&rawJSON),
		)
		if err != nil {
			log.Printf("第 %d 页提取失败: %v", page, err)
			break
		}

		var pagePosts []PostLink
		if err := json.Unmarshal([]byte(rawJSON), &pagePosts); err != nil {
			log.Printf("第 %d 页 JSON 解析失败: %v, raw: %s", page, err, rawJSON)
			break
		}

		newCount := 0
		for _, p := range pagePosts {
			if p.Href == "" || seen[p.Href] {
				continue
			}
			seen[p.Href] = true
			// 补全链接
			if len(p.Href) > 0 && p.Href[0] == '/' {
				p.Href = "https://www.nowcoder.com" + p.Href
			}
			allPosts = append(allPosts, p)
			newCount++
		}

		log.Printf("第 %d 页: 提取到 %d 条，新增 %d 条（累计 %d 条）",
			page, len(pagePosts), newCount, len(allPosts))

		// 打印本页帖子
		for i, p := range pagePosts {
			if seen[p.Href] || p.Href == "" {
				// 已经统计过的跳过打印
			}
			fmt.Printf("  [%d-%02d] %s\n", page, i+1, p.Title)
			fmt.Printf("          %s\n", p.Href)
			if p.Desc != "" {
				desc := p.Desc
				if len(desc) > 80 {
					desc = desc[:80] + "..."
				}
				fmt.Printf("          %s\n", desc)
			}
			fmt.Println()
		}

		// 如果不是最后一页，点击下一页
		if page < limit {
			log.Printf("点击下一页...")

			// 点击分页组件的 "下一页" 按钮
			err = chromedp.Run(ctx,
				// 牛客网用的 el-pagination，下一页按钮是 btn-next
				chromedp.Evaluate(`
					(() => {
						// 方法1: el-pagination 的下一页按钮
						const nextBtn = document.querySelector('.search-agination .btn-next') ||
							document.querySelector('.el-pagination .btn-next') ||
							document.querySelector('button.btn-next');
						if (nextBtn && !nextBtn.disabled) {
							nextBtn.click();
							return 'clicked_btn_next';
						}

						// 方法2: 找到当前页码，点击下一个数字
						const activeNum = document.querySelector('.el-pagination .number.active, .el-pager .number.active');
						if (activeNum && activeNum.nextElementSibling) {
							activeNum.nextElementSibling.click();
							return 'clicked_next_number';
						}

						// 方法3: 找带 > 或 下一页 文字的按钮
						const allBtns = document.querySelectorAll('button, a, span');
						for (const btn of allBtns) {
							const txt = btn.innerText.trim();
							if (txt === '>' || txt === '下一页' || txt === 'Next') {
								btn.click();
								return 'clicked_text_' + txt;
							}
						}
						
						return 'not_found';
					})()
				`, nil),

				// 等待页面刷新 / 数据更新
				chromedp.Sleep(3*time.Second),
			)
			if err != nil {
				log.Printf("翻页失败: %v", err)
				break
			}
		}
	}

	// ========== 输出结果 ==========
	fmt.Println("\n============================================")
	fmt.Printf("✅ 共抓取 %d 页, 累计 %d 条不重复帖子\n", limit, len(allPosts))
	fmt.Println("============================================")

	for i, p := range allPosts {
		fmt.Printf("[%3d] %s\n      %s\n", i+1, p.Title, p.Href)
	}

	// 保存到文件
	output, _ := json.MarshalIndent(allPosts, "", "  ")
	if err := os.WriteFile(outputFile, output, 0644); err != nil {
		log.Printf("保存文件失败: %v", err)
	} else {
		fmt.Printf("\n✅ 结果已保存到 %s\n", outputFile)
	}

	time.Sleep(5 * time.Second)
}
