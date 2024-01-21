package main

// QnACard - 问答格式的卡片
type QnACard struct {
	Question string `toml:"question,omitempty"`
	Answer   string `toml:"answer,omitempty"`
}

// KnowledgePage - 一组卡片的配置
type KnowledgePage struct {
	Title string    `toml:"title"`
	QnAs  []QnACard `toml:"q_a"`
}
