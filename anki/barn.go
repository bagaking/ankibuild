package anki

import (
	"context"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/bagaking/goulp/wlog"
)

type (
	BarnSetting struct {
		Title string `toml:"title,omitempty" json:"title,omitempty"`

		// RuntimeEnabled if set true, record runtime to original file when created
		// todo: 回写的时候会改变格式，不是很理想，再想想办法，或者开发其他 parser
		// todo: 目前只有回写，还不支持读取时沿用原来的 ID
		RuntimeEnabled bool `toml:"runtime,omitempty" json:"runtime,omitempty"`

		// Meta - 卡片 Meta 的默认配置，被具体的卡片覆盖
		Meta
	}

	// Barn - 一组卡片的配置
	Barn struct {
		BarnSetting
		QnAs []QnACard `toml:"q_a,omitempty" json:"q_a,omitempty"`
	}
)

// ParseTomlFile takes a file path and parses the `.apkg.md` config file.
func ParseTomlFile(ctx context.Context, filePath string) (*Barn, error) {
	log := wlog.ByCtx(ctx, "ParseTomlFile").WithField("path", filePath)
	// 读取 TOML 文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return ParseTomlContent(ctx, content)
}

// ParseTomlContent takes a file path and parses the .apkg.md config file.
func ParseTomlContent(ctx context.Context, content []byte) (*Barn, error) {
	var knowledge Barn
	// 使用toml库解析文件内容
	if err := toml.Unmarshal(content, &knowledge); err != nil {
		return nil, err
	}
	return &knowledge, nil
}
