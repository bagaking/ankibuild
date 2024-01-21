package main

type (
	// CardRuntime - to record already created card and note
	CardRuntime struct {
		CardID int `toml:"cid,omitempty"`
		NoteID int `toml:"nid,omitempty"`
	}

	// QnACard - 问答格式的卡片
	QnACard struct {
		Question string   `toml:"question,omitempty"`
		Answer   string   `toml:"answer,omitempty"`
		Tags     []string `toml:"tags,omitempty"`

		Runtime CardRuntime `toml:"runtime,omitempty"`
	}

	// KnowledgePage - 一组卡片的配置
	KnowledgePage struct {
		Title string    `toml:"title,omitempty" json:"title,omitempty"`
		Tags  []string  `toml:"tags,omitempty" json:"tags,omitempty"`
		QnAs  []QnACard `toml:"q_a,omitempty" json:"q_a,omitempty"`

		// RuntimeEnabled if set true, record runtime to original file when created
		// todo: 回写的时候会改变格式，不是很理想，再想想办法，或者开发其他 parser
		// todo: 目前只有回写，还不支持读取时沿用原来的 ID
		RuntimeEnabled bool `toml:"runtime,omitempty" json:"runtime,omitempty"`
	}
)

func (c *QnACard) GetNoteID() int {
	if c.Runtime.NoteID != 0 {
		return c.Runtime.NoteID
	}
	return 0
}

func (c *QnACard) GetCardID() int {
	if c.Runtime.CardID != 0 {
		return c.Runtime.CardID
	}
	return 0
}
