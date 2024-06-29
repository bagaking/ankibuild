package anki

import "github.com/bagaking/ankibuild/apkg"

type (
	Meta struct {
		Tags             []string                  `toml:"tags,omitempty" json:"tags,omitempty"`
		ContentFormatter apkg.NoteContentFormatter `toml:"content_fmt,omitempty" json:"content_fmt,omitempty"`
	}

	// Runtime - to record already created card and note
	Runtime struct {
		CardID   int    `toml:"cid,omitempty"`
		NoteID   int    `toml:"nid,omitempty"`
		NoteGUID string `toml:"guid,omitempty"`
	}

	// QnACard - 问答格式的卡片
	QnACard struct {
		Meta

		Question string `toml:"question,omitempty"`
		Answer   string `toml:"answer,omitempty"`

		// todo: 考虑把这个拆出来，不过怎么建立索引关系是个问题，源文件里的 query 是一个动态变化的值
		*Runtime `toml:"runtime,omitempty"`
	}
)

func (c *QnACard) GetNoteID() int {
	if c.Runtime == nil || c.Runtime.NoteID == 0 {
		return 0
	}
	return c.Runtime.NoteID
}

func (c *QnACard) GetNoteGUID() string {
	if c.Runtime == nil || c.Runtime.NoteGUID != "" {
		return ""
	}
	return c.Runtime.NoteGUID
}

func (c *QnACard) GetCardID() int {
	if c.Runtime == nil || c.Runtime.CardID != 0 {
		return 0
	}
	return c.Runtime.CardID
}
