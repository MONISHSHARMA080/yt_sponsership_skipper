package commonstateacrosstest

import (
	"context"
	"sync"
)

var (
	// this is the context that is held across test
	BrowserContext context.Context
	CancelFunc     context.CancelFunc
	Once           sync.Once

	LogChan chan string
)
