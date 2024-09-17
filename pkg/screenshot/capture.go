package screenshot

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
)

func Capture(ctx context.Context, selector string) (string, error) {
	var buf []byte
	var nodeFound bool

	// Create a new context with a longer timeout
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Evaluate(fmt.Sprintf(`document.querySelector('%s') !== null`, selector), &nodeFound),
	)

	if err != nil {
		return "", err
	}

	if !nodeFound {
		return "", nil
	}

	err = chromedp.Run(ctx,
		chromedp.WaitVisible(selector, chromedp.ByQuery),
		// Sleep for 2 seconds to allow the page to load
		chromedp.Sleep(2*time.Second),
		chromedp.Screenshot(selector, &buf, chromedp.NodeVisible, chromedp.ByQuery),
	)

	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf), nil
}
