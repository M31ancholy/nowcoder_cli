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

			err := cookieCmd.Do(ctx)
			if err != nil {
				log.Printf("⚠️  Cookie %q 设置失败: %v", c.Name, err)
			} else {
				log.Printf("✅ Cookie %q 设置成功", c.Name)
			}
		}
		return nil
	}
}

func main() {
	cookies, err := loadCookiesFromFile("./cmd/tmp/nowcoder_cookie.json")
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

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var title string

	err = chromedp.Run(ctx,
		chromedp.Navigate("https://www.nowcoder.com"),
		chromedp.WaitReady("body"),

		setCookies(cookies),

		chromedp.Reload(),
		chromedp.WaitReady("body"),

		chromedp.Title(&title),
	)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(5 * time.Second)
	fmt.Println("页面标题:", title)
}
