package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bagaking/ankibuild/apkg"
	"github.com/bagaking/goulp/wlog"
	"github.com/pelletier/go-toml"
)

// BuildAPKGsFromToml searches for .apkg.md files in the current directory and subdirectories
// and generates .apkg files accordingly.
func BuildAPKGsFromToml(ctx context.Context) error {
	log, ctx := wlog.ByCtxAndCache(ctx, "BuildAPKGsFromToml")

	return filepath.Walk(".", func(pth string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(info.Name(), ".apkg.toml") {
			return nil
		}

		confK, err := parseConfigFromFile(pth)
		if err != nil {
			return err
		}

		log.Infof("find path= %s", pth)
		outDir := filepath.Dir(pth)

		fileName := strings.TrimSuffix(info.Name(), ".apkg.toml")
		if confK.BarnSetting.Title != "" {
			fileName = confK.BarnSetting.Title
		}

		/* 这里注释下一步的代码，等到我们的 apkg 的包实现创建 apkg 文件的方法之后再解除注释 */
		pkgInfo, err := apkg.CreateDeck(ctx, outDir) // 你的输出文件夹路径
		if err != nil {
			log.WithField("pth", pth).Fatalf("create pkg info failed, outPth= %s, err: %v", outDir, err)
		}
		defer pkgInfo.Close()

		if err = insertCards(ctx, confK, pkgInfo); err != nil {
			log.WithField("pth", pth).Errorf("insert cards failed, err: %v", err)
			return err
		}

		// If RuntimeEnabled is true, write back the runtime information to the TOML configuration file.
		if confK.RuntimeEnabled {
			if err = writeRuntimeBack(&confK, pth); err != nil {
				log.WithField("pth", pth).Warnf("write runtime back failed, err: %v", err)
				return err
			}
		}
		return pkgInfo.ExportToAPKG(filepath.Join(outDir, fileName+".apkg"))
	})

}

func insertCards(ctx context.Context, confK Barn, pkgInfo *apkg.Deck) (err error) {
	log := wlog.ByCtx(ctx, "insertCards")
	log.Infof("confK.BarnSetting= %+v, runtime= %v", confK.BarnSetting, confK.RuntimeEnabled)

	//创建每个卡片并添加到 apkg 包中
	for i := range confK.QnAs {
		cardConf := confK.QnAs[i]
		log.Infof("to create card，q= %+v", cardConf.Question)

		// todo: remove repeated items
		combinedTags := append(confK.BarnSetting.Tags, cardConf.Tags...)
		contentFmt := cardConf.ContentFormatter
		if contentFmt == "" {
			contentFmt = confK.BarnSetting.ContentFormatter
		}

		// 检查是否存在runtime信息，如果存在，则使用现有的NID和CID
		noteID, noteGUID, cardID := cardConf.GetNoteID(), cardConf.GetNoteGUID(), cardConf.GetCardID()

		var n *apkg.Note
		var c *apkg.Card

		// 实际应该判空，只是当前的零值能兼容
		noteOpt := []apkg.NoteOption{
			apkg.NoteWithTags(combinedTags...),
			apkg.NoteWithNID(noteID),
			apkg.NoteWithGUID(noteGUID),
			apkg.NoteWithContentFormatter(contentFmt),
		}
		log.Infof("=== tags= %v, contentFmt= '%v'", combinedTags, contentFmt)

		if n, err = pkgInfo.NoteService().CreateNote(ctx, cardConf.Question, cardConf.Answer, noteOpt...); err != nil {
			return err
		}

		if c, err = pkgInfo.CardService().CreateCard(cardID, n); err != nil {
			return err
		}
		log.Infof("card created，q= %s, n.flds= %s", cardConf.Question, n.FLDs)

		// Update and save runtime information if enabled
		if confK.RuntimeEnabled {
			newCard, err := updateCardRuntime(cardConf, n, c)
			if err != nil {
				return err
			}
			confK.QnAs[i] = *newCard
		}
	}
	return nil
}

// parseConfigFromFile takes a file path and parses the .apkg.md config file.
func parseConfigFromFile(filePath string) (Barn, error) {
	// 读取 TOML 文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	var knowledge Barn

	// 使用toml库解析文件内容
	if err = toml.Unmarshal(content, &knowledge); err != nil {
		return knowledge, err
	}

	return knowledge, nil
}
