package apkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

const (
	VirtualDeckID   = 170214000000
	VirtualDeckName = "BK.EasyExport"
	SimpleTplID     = 1000
)

type Template struct {
	Ord   int    `json:"ord"`
	Bsize int    `json:"bsize"`
	Did   *int   `json:"did"`
	Bafmt string `json:"bafmt"`
	Qfmt  string `json:"qfmt"`
	Afmt  string `json:"afmt"`
	Bfont string `json:"bfont"`
	Bqfmt string `json:"bqfmt"`
	Name  string `json:"name"`
}

type Field struct {
	Ord    int      `json:"ord"`
	Sticky bool     `json:"sticky"`
	Rtl    bool     `json:"rtl"`
	Media  []string `json:"media"`
	Size   int      `json:"size"`
	Font   string   `json:"font"`
	Name   string   `json:"name"`
}

type TplModel struct {
	ID        int           `json:"id"`
	Tmpls     []Template    `json:"tmpls"`
	LatexPre  string        `json:"latexPre"`
	Req       []any         `json:"req"`
	Flds      []Field       `json:"flds"`
	Tags      []string      `json:"tags"`
	Type      int           `json:"type"`
	Mod       int64         `json:"mod"`
	LatexSVG  int           `json:"latexsvg"`
	Sortf     int           `json:"sortf"`
	Usn       int           `json:"usn"`
	Did       int           `json:"did"`
	Vers      []interface{} `json:"vers"`
	LatexPost string        `json:"latexPost"`
	Name      string        `json:"name"`
	Css       string        `json:"css"`
	Gf        bool          `json:"gf"`
}

func (cs *PkgInfo) CreateSimpleDeck(id int) (*Col, error) {
	// 创建一个简单模板实例
	simpleTpl := TplModel{
		ID:   id,
		Name: VirtualDeckName + ".TPL",
		Tmpls: []Template{
			{
				Qfmt: "{{正面}}",              // 前面的模板
				Afmt: "{{FrontSide}}{{背面}}", // 后面的模板
				Name: "问答卡片",
			},
		},
		Req: []any{0, "all", []any{0}},
		Flds: []Field{
			{
				Ord:  0,
				Size: 20,
				Name: "正面",
			},
			{
				Ord:  1,
				Size: 20,
				Name: "背面",
			},
		},
		Mod: time.Now().Unix(),
	}

	// 序列化模板为JSON字符串
	templatesJSON, err := json.Marshal(map[string]TplModel{
		fmt.Sprintf("%d", id): simpleTpl,
	})
	if err != nil {
		return nil, err
	}

	// 创建deck模型，包括刚刚生成的模板
	col := &Col{
		ID:     genID(),
		Mod:    time.Now().Unix(),
		Models: string(templatesJSON), // 将序列化后的JSON字符串作为模型
		Decks: fmt.Sprintf(
			`{"%d":{"newToday":[314,0],"revToday":[314,0],"lrnToday":[314,0],"timeToday":[314,0],"conf":1702143523045,"usn":9,"desc":"","dyn":0,"collapsed":false,"extendNew":10,"extendRev":50,"name":"%s","id":%d,"mod":1702208073}}`,
			VirtualDeckID,
			VirtualDeckName,
			VirtualDeckID),
	}

	// 在数据库中创建牌组
	if err := cs.DB.Create(col).Error; err != nil {
		return nil, err
	}

	return col, nil
}

func (cs *PkgInfo) FindOrCreateSimpleDeck() (*Col, error) {
	var col Col
	// 在数据库中尝试找到现有的Col
	err := cs.DB.Where("id > ?", 0).First(&col).Error

	// 如果找到了，直接返回找到的Col
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 如果不存在，则创建一个新的Col
		return cs.CreateSimpleDeck(SimpleTplID)
	} else if err != nil {
		// 如果发生了其他错误，则返回错误
		return nil, err
	}

	return &col, nil
}
