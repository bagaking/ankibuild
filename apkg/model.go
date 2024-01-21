package apkg

import (
	"strings"

	"gorm.io/gorm"
)

const SplitFieldOfNote = "\x1f"

// Col represents the 'col' table which contains information about the collection like settings and statistics.
type Col struct {
	ID int `gorm:"primaryKey;column:id" json:"id"` // Unique identifier for the collection

	Crt    int    `gorm:"column:crt" json:"crt"`       // Collection created time (timestamp)
	Mod    int64  `gorm:"column:mod" json:"mod"`       // Last modified in milliseconds
	Scm    int64  `gorm:"column:scm" json:"scm"`       // Schema modified time (used for syncing)
	Ver    int    `gorm:"column:ver" json:"ver"`       // Version of the collection
	Dty    int    `gorm:"column:dty" json:"dty"`       // Dirty (need to be synced)
	Usn    int    `gorm:"column:usn" json:"usn"`       // Update sequence number (for finding diffs when syncing)
	Ls     int    `gorm:"column:ls" json:"ls"`         // Last sync time
	Conf   string `gorm:"column:conf" json:"conf"`     // Configuration JSON object
	Models string `gorm:"column:models" json:"models"` // JSON object containing note models
	Decks  string `gorm:"column:decks" json:"decks"`   // JSON object containing deck information
	Dconf  string `gorm:"column:dconf" json:"dconf"`   // JSON object containing deck options
	Tags   string `gorm:"column:tags" json:"tags"`     // A cache of tags used in the collection
}

// Note represents the 'notes' table which stores all the information of notes including metadata and the content itself.
type Note struct {
	ID int `gorm:"primaryKey;column:id" json:"id"` // Unique identifier for the note

	Guid string `gorm:"column:guid" json:"guid"` // Globally unique id, used for syncing

	// Mid - Model ID,
	// links to card templates which define how to use the fields to generate "front" and "back" sides of cards.
	Mid int `gorm:"column:mid" json:"mid"`
	// Mod - Last modified in milliseconds
	Mod int64 `gorm:"column:mod" json:"mod"`

	// Usn - Update sequence number
	Usn  int    `gorm:"column:usn" json:"usn"`
	Tags string `gorm:"column:tags" json:"tags"` // Space-separated string of tags.

	// Flds - Fields of the note joined by 0x1f character.
	// They are used by card templates to generate "front" and "back" sides of cards.
	//
	// - "Basic" notes have "Front" and "Back" fields for basic question-answer cards.
	// - "Basic (and reversed card)" notes add a reversed card from the answer to the question alongside the basic card.
	// - "Basic (optional reversed card)" has "Front", "Back", and "Add Reverse" fields, creating reversed cards when "Add Reverse" is filled.
	// - "Cloze" notes are used to create fill-in-the-blank cards where text is omitted.
	// - "Image Occlusion" notes utilize images with sections blocked out to test recognition of image parts.
	Flds string `gorm:"column:flds" json:"flds"`

	// Sfld - Sort field: the value of the field by which notes are sorted in the browser.
	Sfld int `gorm:"column:sfld" json:"sfld"`

	Csum  int64  `gorm:"column:csum" json:"csum"`   // Checksum used for duplicate check.
	Flags int    `gorm:"column:flags" json:"flags"` // Flags
	Data  string `gorm:"column:data" json:"data"`   // Unused, currently just an empty string.
}

func (n *Note) Front() string {
	return strings.Split(n.Flds, SplitFieldOfNote)[0]
}

// todo: 猜测
const (
	CardTypeNew = iota
	CardTypeLearning
	CardTypeReview
	CardTypeRelearning
	CardTypeDue
)

// todo: 猜测
const (
	CardQueueTypeNew = iota
	CardQueueTypeLearning
	CardQueueTypeReview
	CardQueueTypeRelearning
	CardQueueTypeDue
)

// Card represents the 'cards' table which stores the information about the review cards generated from notes.
// Anki 卡片是通过笔记（Notes）和卡片类型（Card Types）相结合来实现的。每个笔记可以生成一个或多个卡
type Card struct {
	ID int `gorm:"primaryKey;column:id" json:"id"` // Unique identifier for the card

	// NID - Note ID
	NID int `gorm:"column:nid" json:"nid"`
	// DID - Deck ID
	DID int `gorm:"column:did" json:"did"`

	// Ord - Ordinal, identifies card's template
	Ord int `gorm:"column:ord" json:"ord"`

	Mod int64 `gorm:"column:mod" json:"mod"` // Last modified in milliseconds
	Usn int   `gorm:"column:usn" json:"usn"` // Update sequence number

	// Type of the card: 0: New, 1: Learning, 2: Review, 3: Relearning
	// 表示卡片的类型（如新的、学习中、待复习等）
	Type int `gorm:"column:type" json:"type"`

	// Queue the card is in (new, learning, due, etc.)
	// 表示卡片当前所处的队列（如新的、学习中、待复习等）
	Queue int `gorm:"column:queue" json:"queue"`

	// // Due date for review. For new cards and learning cards this is the order in which they will be shown.
	// 表示卡片的到期日期
	Due int `gorm:"column:due" json:"due"`

	// Interval (used in SRS algorithm). Determines if the review card is young (Ivl < 21) or mature (Ivl >= 21).
	// 用于SRS算法的间隔时间，用于决定卡片是"年轻"还是"成熟"
	Ivl int `gorm:"column:ivl" json:"ivl"`

	Factor int `gorm:"column:factor" json:"factor"` // Ease factor (used in SRS algorithm)
	Reps   int `gorm:"column:reps" json:"reps"`     // Number of reviews
	Lapses int `gorm:"column:lapses" json:"lapses"` // Number of lapses
	Left   int `gorm:"column:left" json:"left"`     // Steps left to graduation (in learning)
	Odue   int `gorm:"column:odue" json:"odue"`     // Original due date (for cards in relearning)
	Odid   int `gorm:"column:odid" json:"odid"`     // Original deck ID (for cards in filtered decks)

	Flags int    `gorm:"column:flags" json:"flags"` // Flags
	Data  string `gorm:"column:data" json:"data"`   // Unused, currently just an empty string.
}

// Revlog represents the 'revlog' table which stores the review logs of cards.
type Revlog struct {
	ID int64 `gorm:"primaryKey;column:id" json:"id"` // Timestamp (based on 13 digits), used as ID

	Cid int `gorm:"column:cid" json:"cid"` // Card ID

	// Usn - Update sequence number
	Usn     int `gorm:"column:usn" json:"usn"`
	Ease    int `gorm:"column:ease" json:"ease"`       // Ease button pressed (again, hard, good, easy)
	Ivl     int `gorm:"column:ivl" json:"ivl"`         // New interval
	LastIvl int `gorm:"column:lastIvl" json:"lastIvl"` // Last interval
	Factor  int `gorm:"column:factor" json:"factor"`   // New ease factor
	Time    int `gorm:"column:time" json:"time"`       // Review time in milliseconds
	Type    int `gorm:"column:type" json:"type"`       // Review type
}

// Grave represents the 'graves' table which logs deleted cards, notes, and decks.
type Grave struct {
	ID   int `gorm:"primaryKey;column:id" json:"id"` // Unused, currently just set to zero
	Usn  int `gorm:"column:usn" json:"usn"`          // Update sequence number
	Oid  int `gorm:"column:oid" json:"oid"`          // Original ID (card, note, or deck id)
	Type int `gorm:"column:type" json:"type"`        // Type of object (card, note, or deck)
}

func (Col) TableName() string {
	return "col"
}

func (Note) TableName() string {
	return "notes"
}

func (Card) TableName() string {
	return "cards"
}

func (Revlog) TableName() string {
	return "revlog"
}

func (Grave) TableName() string {
	return "grave"
}

// CreateIndexes creates database indexes to improve query performance.
func CreateIndexes(db *gorm.DB) error {
	indexStatements := []string{
		"CREATE INDEX IF NOT EXISTS idx_notes_usn ON notes (usn);",
		"CREATE INDEX IF NOT EXISTS idx_cards_usn ON cards (usn);",
		"CREATE INDEX IF NOT EXISTS idx_revlog_usn ON revlog (usn);",
		"CREATE INDEX IF NOT EXISTS idx_cards_nid ON cards (nid);",
		"CREATE INDEX IF NOT EXISTS idx_cards_sched ON cards (did, queue, due);",
		"CREATE INDEX IF NOT EXISTS idx_revlog_cid ON revlog (cid);",
		"CREATE INDEX IF NOT EXISTS idx_notes_csum ON notes (csum);",
	}

	for _, stmt := range indexStatements {
		err := db.Exec(stmt).Error
		if err != nil {
			return err
		}
	}

	return nil
}
