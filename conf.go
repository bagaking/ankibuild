package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/bagaking/ankibuild/anki"
	"github.com/bagaking/goulp/wlog"
)

type TomlProcessor func(ctx context.Context, confK anki.Barn, pth, outDir, fileName string) error

func WalkTomlFiles(ctx context.Context, processor TomlProcessor) error {
	return filepath.Walk(".", func(pth string, info os.FileInfo, err error) error {
		if info == nil || !strings.HasSuffix(info.Name(), ".apkg.toml") {
			return nil
		}

		if err != nil {
			return err
		}

		log, ctxIn := wlog.ByCtxAndCache(ctx, "WalkTomlFiles")
		log.WithField("pth", pth)

		confK, err := anki.ParseTomlFile(ctx, pth)
		if err != nil {
			return err
		}

		log.Infof("find conf at path %s", pth)
		outDir := filepath.Dir(pth)

		fileName := strings.TrimSuffix(info.Name(), ".apkg.toml")
		if confK.BarnSetting.Title != "" {
			fileName = confK.BarnSetting.Title
		}

		return processor(ctxIn, *confK, pth, outDir, fileName)
	})
}
