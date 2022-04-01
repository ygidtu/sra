package client

import (
	"context"
	"github.com/chromedp/chromedp"
	"go.uber.org/zap"
)

func SetChromeClient(open bool, proxy string, sugar *zap.SugaredLogger) (context.Context, context.CancelFunc) {
	// create context
	sugar.Debug("Start chrome under headless ? ", !open)
	opts := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", !open),
		chromedp.DisableGPU,
		chromedp.NoSandbox,
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
	}

	if proxy != "" {
		opts = append(
			opts,
			chromedp.ProxyServer(proxy),
		)
	}

	contextOpts := []chromedp.ContextOption{
		chromedp.WithLogf(sugar.Infof),
		chromedp.WithErrorf(sugar.Infof),
	}
	allocContext, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(allocContext, contextOpts...)
	return ctx, cancel
}
