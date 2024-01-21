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

// BuildAPKGsFromToml searches for .apkg.md files in the current directory and subdirectories
// and generates .apkg files accordingly.
func BuildAPKGsFromToml(ctx context.Context) error {
	log := wlog.ByCtx(ctx, "BuildAPKGsFromToml")

	return filepath.Walk(".", func(pth string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(info.Name(), ".apkg.toml") {
			return nil
		}

		conf, err := parseConfigFromFile(pth)
		log.Infof("find path= %s", pth)

		outDir := filepath.Dir(pth)

		if err != nil {
			return err
		}

		fileName := strings.TrimSuffix(info.Name(), ".apkg.toml")
		if conf.Title != "" {
			fileName = conf.Title
		}

		/* 这里注释下一步的代码，等到我们的 apkg 的包实现创建 apkg 文件的方法之后再解除注释 */
		pkgInfo, err := apkg.CreatePkgInfo(ctx, outDir) // 你的输出文件夹路径
		if err != nil {
			log.Fatalf("create pkg info failed, outPth= %s, err: %v", outDir, err)
		}
		defer pkgInfo.Close()

		//创建每个卡片并添加到 apkg 包中
		for i := range conf.QnAs {
			card := conf.QnAs[i]
			log.Infof("create card，q= %s，a= %s", card.Question, card.Answer)

			combinedTags := append(conf.Tags, card.Tags...)

			// 检查是否存在runtime信息，如果存在，则使用现有的NID和CID
			noteID := card.GetNoteID()
			cardID := card.GetCardID()

			var n *apkg.Note
			var c *apkg.Card

			if n, err = pkgInfo.CardService().CreateNote(noteID, card.Question, card.Answer, combinedTags...); err != nil {
				return err
			}

			if c, err = pkgInfo.CardService().CreateCard(cardID, n); err != nil {
				return err
			}
			log.Infof("card created，q= %s，a= %s, n= %+v, c= %+v", card.Question, card.Answer, n, c)

			// Update and save runtime information if enabled
			if conf.RuntimeEnabled {
				newCard, err := updateCardRuntime(n, c, card)
				if err != nil {
					return err
				}
				conf.QnAs[i] = *newCard
			}
		}

		// If RuntimeEnabled is true, write back the runtime information to the TOML configuration file.
		if conf.RuntimeEnabled {
			if err = writeRuntimeBack(&conf, pth); err != nil {
				return err
			}
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
