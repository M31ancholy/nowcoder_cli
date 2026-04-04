package nowcoder

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

type Interview struct{}

func GetInterviews() ([]Interview, error) {
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
		setCookies(cookies),
		// 导航到面经搜索之后的界面
		chromedp.Navigate("https://www.nowcoder.com/search/all?query=面经&type=all&searchType=顶部导航栏"),
		chromedp.WaitReady("body"),
		chromedp.Reload(),
		chromedp.WaitReady("body"),

		chromedp.Title(&title),
	)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(5 * time.Second)
	fmt.Println("页面标题:", title)
	return nil, err
}
