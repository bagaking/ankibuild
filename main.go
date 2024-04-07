package main

import (
	"context"
	"flag"

	"github.com/bagaking/goulp/wlog"
)

var exportFormat string

func init() {
	flag.StringVar(&exportFormat, "format", "apkg", "The export format: 'apkg' or 'excel'")
	flag.Parse()
}

func main() {
	ctx := context.Background()
	var err error

	switch exportFormat {
	case "apkg":
		err = BuildAPKGsFromToml(ctx)
	case "excel":
		err = BuildExcelsFromToml(ctx)
	default:
		wlog.ByCtx(ctx).Fatalf("Unsupported export format: %s", exportFormat)
	}

	if err != nil {
		wlog.ByCtx(ctx).Fatal(err)
	}
}
