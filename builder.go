package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/bagaking/ankibuild/apkg"
	"github.com/bagaking/goulp/wlog"
)

// BuildAPKGsFromConfig searches for .apkg.md files in the current directory and subdirectories
// and generates .apkg files accordingly.
func BuildAPKGsFromConfig(ctx context.Context) error {
	log := wlog.ByCtx(ctx, "BuildAPKGsFromConfig")

	return filepath.Walk(".", func(pth string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(info.Name(), ".apkg.toml") {
			return nil
		}
		config, err := parseConfigFromFile(pth)
		log.Infof("find path= %s", pth)

		outDir := filepath.Dir(pth)

		if err != nil {
			return err
		}

		fileName := strings.TrimSuffix(info.Name(), ".apkg.toml")
		if config.Title != "" {
			fileName = config.Title
		}
		/* 这里注释下一步的代码，等到我们的 apkg 的包实现创建 apkg 文件的方法之后再解除注释 */
		pkgInfo, err := apkg.CreatePkgInfo(ctx, outDir) // 你的输出文件夹路径
		if err != nil {
			log.Fatalf("create pkg info failed, outPth= %s, err: %v", outDir, err)
		}
		defer pkgInfo.Close()

		//创建每个卡片并添加到 apkg 包中
		for _, card := range config.QnAs {
			log.Infof("create card，q= %s，a= %s", card.Question, card.Answer)
			// 添加卡片进入 apkg 包
			n, c, err := pkgInfo.CardService().CreateCard(card.Question, card.Answer)
			if err != nil {
				return err
			}
			log.Infof("card created，q= %s，a= %s, n= %+v, c= %+v", card.Question, card.Answer, n, c)

		}
		return pkgInfo.ExportToAPKG(filepath.Join(outDir, fileName+".apkg"))
	})

}

// parseConfigFromFile takes a file path and parses the .apkg.md config file.
func parseConfigFromFile(filePath string) (KnowledgePage, error) {
	// 读取 TOML 文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	var knowledge KnowledgePage

	// 使用toml库解析文件内容
	if _, err = toml.Decode(string(content), &knowledge); err != nil {
		return knowledge, err
	}

	return knowledge, nil
}
