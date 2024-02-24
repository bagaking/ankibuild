package apkg

import (
	"context"
	"strings"
	"time"

	"github.com/russross/blackfriday/v2"
	"gorm.io/gorm"
)

type (
	// NoteService provides methods to work with notes and cards.
	NoteService struct {
		DB *gorm.DB
	}

	NoteContentFormatter string

	// NoteOptions includes options for creating a new note.
	NoteOptions struct {
		NID              int                  `json:"nid,omitempty"`
		GUID             string               `json:"guid,omitempty"`
		Tags             []string             `json:"tags,omitempty"`
		ContentFormatter NoteContentFormatter `json:"note_content_formatter,omitempty"`
	}

	// NoteOption configures NoteOptions.
	NoteOption func(*NoteOptions)
)

const (
	NoteCFmtPlainText NoteContentFormatter = ""
	NoteCFmtMarkdown                       = "markdown"
)

func (no *NoteOptions) Use(opts ...NoteOption) *NoteOptions {
	// Apply all options
	for _, opt := range opts {
		opt(no)
	}
	return no
}

func NoteWithNID(nid int) NoteOption {
	return func(opts *NoteOptions) {
		opts.NID = nid
	}
}

func NoteWithTags(tags ...string) NoteOption {
	return func(opts *NoteOptions) {
		opts.Tags = append(opts.Tags, tags...)
	}
}

func NoteWithGUID(guid string) NoteOption {
	return func(opts *NoteOptions) {
		opts.GUID = guid
	}
}

func NoteWithContentFormatter(contentFmt NoteContentFormatter) NoteOption {
	return func(opts *NoteOptions) {
		opts.ContentFormatter = contentFmt
	}
}

func (cs *NoteService) CreateNote(ctx context.Context, front, back string, opts ...NoteOption) (*Note, error) {
	options := (&NoteOptions{}).Use(opts...)

	// Ensure we have a NID
	if options.NID == 0 {
		options.NID = genID()
	}

	if options.GUID == "" {
		guid, err := genGUID()
		if err != nil {
			return nil, err
		}
		options.GUID = guid
	}

	switch options.ContentFormatter {
	case NoteCFmtMarkdown:
		back = string(blackfriday.Run([]byte(back)))
	case NoteCFmtPlainText:
	default:
	}

	tags := append(GlobalTags, options.Tags...)

	flds := makeFlds(front, back)
	// 创建Note
	note := &Note{
		ID:   options.NID,
		Guid: options.GUID,
		Mid:  SimpleTplID, //  collection.ModelsID, // 模板ID
		Mod:  time.Now().Unix(),
		//Usn:  cs.DB.GetNextUsn(), // Update Sequence Number
		Tags: strings.Join(tags, " "),
		FLDs: flds,
		SFLD: generateSortField(flds), // 生成排序字段
		CSum: calculateChecksum(flds), // 内容校验和
	}

	if err := cs.DB.Create(note).Error; err != nil {
		return nil, err
	}

	return note, nil
}
