package screenshot

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/chromedp/chromedp"
)

func Capture(ctx context.Context, selector string) (string, error) {
	var buf []byte
	var nodeFound bool

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
		chromedp.Screenshot(selector, &buf, chromedp.NodeVisible, chromedp.ByQuery),
	)

	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf), nil
}
