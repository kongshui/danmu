package sse

import "context"

var (
	first_ctx = context.Background()
	ChanPool  = NewChanPool(10)
)
