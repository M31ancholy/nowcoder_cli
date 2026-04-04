package nowcoder

import (
	"context"
)

type Interview struct {
	Title   string
	Content string
}

type Service struct {
	allocCtx context.Context
	cancel   context.CancelFunc
	cookies  []BrowserCookie
}

func NewService() *Service {
	return &Service{}
}
