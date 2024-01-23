package main

import (
	"context"

	"github.com/bagaking/goulp/wlog"
)

func main() {
	ctx := context.Background()

	err := BuildAPKGsFromToml(ctx)
	if err != nil {
		wlog.ByCtx(ctx).Fatal(err)
	}
}
