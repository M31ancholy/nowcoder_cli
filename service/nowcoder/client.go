package nowcoder

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

type Service struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Start(ctx context.Context) error {
	opts := chromedp.DefaultExecAllocatorOptions[:]
	opts = append(opts, chromedp.NoFirstRun, chromedp.NoDefaultBrowserCheck)

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	s.ctx, s.cancel = chromedp.NewContext(allocCtx)
	defer cancel()

	s.ctx, cancel = context.WithTimeout(s.ctx, 30*time.Second)
	return nil
}

func (s *Service) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
}

func (s *Service) Login(username, password string) error {
	if err := chromedp.Run(s.ctx,
		chromedp.Navigate("https://www.nowcoder.com/login"),
		chromedp.Sleep(2*time.Second),
	); err != nil {
		return fmt.Errorf("navigate failed: %w", err)
	}

	log.Printf("Login called with username: %s", username)
	return nil
}

func (s *Service) GetProfile() error {
	if err := chromedp.Run(s.ctx,
		chromedp.Navigate("https://www.nowcoder.com/profile"),
		chromedp.Sleep(2*time.Second),
	); err != nil {
		return fmt.Errorf("navigate to profile failed: %w", err)
	}

	var username string
	if err := chromedp.Run(s.ctx,
		chromedp.Text(".username", &username, chromedp.ByQuery),
	); err != nil {
		return fmt.Errorf("get username failed: %w", err)
	}

	log.Printf("Username: %s", username)
	return nil
}
