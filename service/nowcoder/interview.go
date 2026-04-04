package nowcoder

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/chromedp"
)

type Interview struct {
	Title   string
	Content string
}

type Service struct {
	allocCtx context.Context
	cancel   context.CancelFunc
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Start() error {
	ctx := context.Background()
	opts := chromedp.DefaultExecAllocatorOptions[:]
	opts = append(opts, chromedp.NoFirstRun, chromedp.NoDefaultBrowserCheck)

	s.allocCtx, s.cancel = chromedp.NewExecAllocator(ctx, opts...)
	s.allocCtx, s.cancel = context.WithTimeout(s.allocCtx, 30*time.Second)
	return nil
}

func (s *Service) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
}

func (s *Service) Hunt(company, position, outputDir string) (*Interview, error) {
	if err := s.Start(); err != nil {
		return nil, fmt.Errorf("start browser failed: %w", err)
	}
	defer s.Stop()

	ctx, cancel := chromedp.NewContext(s.allocCtx)
	defer cancel()

	timestamp := time.Now().Format("20060102_150405")
	searchQuery := fmt.Sprintf("%s %s", company, position)

	if err := s.openNowcoder(ctx); err != nil {
		return nil, err
	}

	if err := s.search(ctx, searchQuery); err != nil {
		return nil, err
	}

	if err := s.filterInterview(ctx); err != nil {
		return nil, err
	}

	screenshotPath := filepath.Join(outputDir, fmt.Sprintf("search_%s.png", timestamp))
	if err := s.screenshot(ctx, screenshotPath); err != nil {
		return nil, err
	}

	if err := s.clickFirstInterview(ctx); err != nil {
		return nil, err
	}

	interview, err := s.extractInterviewDetail(ctx, screenshotPath)
	if err != nil {
		return nil, err
	}

	rawFile := filepath.Join(outputDir, fmt.Sprintf("%s_%s_raw_%s.txt", company, position, timestamp))
	if err := os.WriteFile(rawFile, []byte(fmt.Sprintf("Title: %s\n\n%s", interview.Title, interview.Content)), 0644); err != nil {
		return nil, fmt.Errorf("save raw content failed: %w", err)
	}

	return interview, nil
}

func (s *Service) openNowcoder(ctx context.Context) error {
	if err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.nowcoder.com"),
		chromedp.Sleep(2*time.Second),
	); err != nil {
		return fmt.Errorf("open nowcoder failed: %w", err)
	}
	return nil
}

func (s *Service) search(ctx context.Context, query string) error {
	if err := chromedp.Run(ctx,
		chromedp.Sleep(1*time.Second),
	); err != nil {
		return err
	}

	if err := chromedp.Run(ctx,
		chromedp.Click(`input[placeholder="搜索面经/职位/试题/公司"]`, chromedp.ByQuery),
		chromedp.Sleep(1*time.Second),
	); err != nil {
		return fmt.Errorf("click search box failed: %w", err)
	}

	if err := chromedp.Run(ctx,
		chromedp.SendKeys(`input[placeholder="搜索面经/职位/试题/公司"]`, query+"\n", chromedp.ByQuery),
		chromedp.Sleep(3*time.Second),
	); err != nil {
		return fmt.Errorf("type search query failed: %w", err)
	}

	return nil
}

func (s *Service) filterInterview(ctx context.Context) error {
	if err := chromedp.Run(ctx,
		chromedp.Sleep(2*time.Second),
	); err != nil {
		return err
	}

	if err := chromedp.Run(ctx,
		chromedp.Click(`@e52`, chromedp.ByQuery),
	); err != nil {
		return fmt.Errorf("click interview filter failed: %w", err)
	}

	return nil
}

func (s *Service) screenshot(ctx context.Context, path string) error {
	var buf []byte
	if err := chromedp.Run(ctx,
		chromedp.CaptureScreenshot(&buf),
	); err != nil {
		return fmt.Errorf("screenshot failed: %w", err)
	}

	if err := os.WriteFile(path, buf, 0644); err != nil {
		return fmt.Errorf("write screenshot failed: %w", err)
	}
	return nil
}

func (s *Service) clickFirstInterview(ctx context.Context) error {
	if err := chromedp.Run(ctx,
		chromedp.Sleep(3*time.Second),
	); err != nil {
		return err
	}

	if err := chromedp.Run(ctx,
		chromedp.Click(`@e110`, chromedp.ByQuery),
		chromedp.Sleep(3*time.Second),
	); err != nil {
		return fmt.Errorf("click first interview failed: %w", err)
	}

	return nil
}

func (s *Service) extractInterviewDetail(ctx context.Context, screenshotPath string) (*Interview, error) {
	var title string
	if err := chromedp.Run(ctx,
		chromedp.Title(&title),
		chromedp.Sleep(2*time.Second),
	); err != nil {
		return nil, fmt.Errorf("get title failed: %w", err)
	}

	if err := chromedp.Run(ctx,
		chromedp.EvaluateAsDevTools("window.scrollTo(0, 1500)", nil),
		chromedp.Sleep(1*time.Second),
	); err != nil {
		return nil, fmt.Errorf("scroll failed: %w", err)
	}

	var content string
	if err := chromedp.Run(ctx,
		chromedp.OuterHTML("html", &content, chromedp.ByQuery),
	); err != nil {
		return nil, fmt.Errorf("get content failed: %w", err)
	}

	screenshotPath = filepath.Join(filepath.Dir(screenshotPath), "detail_"+filepath.Base(screenshotPath))
	if err := s.screenshot(ctx, screenshotPath); err != nil {
		return nil, err
	}

	return &Interview{
		Title:   title,
		Content: content,
	}, nil
}
