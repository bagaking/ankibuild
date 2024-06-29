package apkg

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

// ErrNoteNotFound is returned when a note cannot be found.
var (
	ErrNoteNotFound = errors.New("note not found")
	GlobalTags      = []string{}
)

// CardService provides methods to work with notes and cards.
type CardService struct {
	DB *gorm.DB
}

// CreateCard creates a new card based on the given front and back information.
func (cs *CardService) CreateCard(cid int, note *Note) (*Card, error) {
	// 创建Note
	if cid == 0 {
		cid = genID()
	}

	card := &Card{
		ID:  cid,
		NID: note.ID,
		DID: VirtualDeckID, // Deck ID
		Mod: time.Now().Unix(),
		// Usn:   cs.DB.GetNextUsn(), // Update Sequence Number
		Type:  CardTypeNew,
		Queue: CardQueueTypeNew,
		// Due:   cs.DB.GetNextDue(deck),
		// Ivl:   DefaultInitialInterval,
		// Set other necessary Card fields based on your business logic
	}

	if err := cs.DB.Create(card).Error; err != nil {
		return nil, err
	}

	return card, nil
}

// GetAllFronts returns a slice of all fronts from notes.
func (cs *CardService) GetAllFronts() ([]string, error) {
	var notes []Note
	if err := cs.DB.Find(&notes).Error; err != nil {
		return nil, err
	}

	var fronts []string
	for _, note := range notes {
		fields := strings.Split(note.FLDs, "\x1f")
		if len(fields) > 0 {
			fronts = append(fronts, fields[0])
		}
	}

	return fronts, nil
}

// FindCardByFront finds notes and cards with the given front.
func (cs *CardService) FindCardByFront(front string) ([]Note, []Card, error) {
	var notes []Note
	if err := cs.DB.Where("flds LIKE ?", front+"%").Find(&notes).Error; err != nil {
		return nil, nil, err
	}

	if len(notes) == 0 {
		return nil, nil, ErrNoteNotFound
	}

	var cards []Card
	for _, note := range notes {
		card := Card{}
		if err := cs.DB.Where("nid = ?", note.ID).First(&card).Error; err != nil {
			return nil, nil, err
		}
		cards = append(cards, card)
	}

	return notes, cards, nil
}

func makeFlds(front, back string) string {
	return strings.Join([]string{front, back}, SplitFieldOfNote)
}
