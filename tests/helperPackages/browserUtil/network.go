package browserutil

import (
	"context"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// disables the network in  the tab
func DisableNetwork(ctx context.Context) error {
	return chromedp.Run(ctx,
		// 1) enable network domain
		network.Enable(),
		// 2) emulate offline (no network)
		network.EmulateNetworkConditions(
			true, // offline
			0,    // latency (ms)
			0,    // downloadThroughput
			0,    // uploadThroughput
		),
	)
}
