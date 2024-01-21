package apkg

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

// ErrNoteNotFound is returned when a note cannot be found.
var ErrNoteNotFound = errors.New("note not found")
var GlobalTags = []string{}

// CardService provides methods to work with notes and cards.
type CardService struct {
	DB *gorm.DB
}

func (p *PkgInfo) CardService() *CardService {
	return &CardService{DB: p.DB}
}

func (cs *CardService) CreateNote(flds string, tags ...string) (*Note, error) {
	noteID := genID()
	guid, err := genGUID()
	if err != nil {
		return nil, err
	}
	tags = append(GlobalTags, tags...)
	// 创建Note
	note := &Note{
		ID:   noteID,
		Guid: guid,
		Mid:  SimpleTplID, //  collection.ModelsID, // 模板ID
		Mod:  time.Now().Unix(),
		//Usn:  cs.DB.GetNextUsn(), // Update Sequence Number
		Tags: strings.Join(tags, " "),
		Flds: flds,
		Sfld: generateSortField(flds), // 生成排序字段
		Csum: calculateChecksum(flds), // 内容校验和
	}

	if err := cs.DB.Create(note).Error; err != nil {
		return nil, err
	}

	return note, nil
}

// CreateCard creates a new card based on the given front and back information.
func (cs *CardService) CreateCard(front, back string, tags ...string) (*Note, *Card, error) {
	// 创建Note
	note, err := cs.CreateNote(front + "\x1f" + back)
	if err != nil {
		return nil, nil, err
	}

	card := &Card{
		ID:  genID(),
		Nid: note.ID,
		Did: VirtualDeckID, // Deck ID
		Mod: time.Now().Unix(),
		//Usn:   cs.DB.GetNextUsn(), // Update Sequence Number
		Type:  CardTypeNew,
		Queue: CardQueueTypeNew,
		//Due:   cs.DB.GetNextDue(deck),
		//Ivl:   DefaultInitialInterval,
		// Set other necessary Card fields based on your business logic
	}

	if err := cs.DB.Create(card).Error; err != nil {
		return nil, nil, err
	}

	return note, card, nil
}

// GetAllFronts returns a slice of all fronts from notes.
func (cs *CardService) GetAllFronts() ([]string, error) {
	var notes []Note
	if err := cs.DB.Find(&notes).Error; err != nil {
		return nil, err
	}

	var fronts []string
	for _, note := range notes {
		fields := strings.Split(note.Flds, "\x1f")
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
