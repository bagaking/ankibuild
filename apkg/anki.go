package apkg

import (
	"context"

	"github.com/bagaking/goulp/wlog"
	"gorm.io/gorm"
)

const ApkgDBName = "collection.anki2"

// Deck Anki 的 apkg 管理实体, 一一对应 deck
// Ref:
// - official: https://docs.ankiweb.net/intro.html
// - open source:
// -- https://gist.github.com/sartak/3921255
// -- https://github.com/SergioFacchini/anki-cards-web-browser/blob/master/documentation/Processing%20Anki's%20.apkg%20files.md
type Deck struct {
	Path string
	*Storage
}

// CreateDeck initializes and returns a new Deck with a connected database.
func CreateDeck(ctx context.Context, fPath string) (*Deck, error) {
	var deckCol *Col
	storage, err := CreateStorage(ctx, fPath, ApkgDBName, func(db *gorm.DB) (eInit error) {
		// AutoMigrateModels migrates the provided models and creates necessary indexes.
		if eInit = AutoMigrateModels(db); eInit != nil {
			return eInit
		}
		if eInit = CreateIndexes(db); eInit != nil {
			return eInit
		}
		if deckCol, eInit = findOrMockSimpleDeck(db); eInit != nil {
			return eInit
		}
		return
	})
	if err != nil {
		return nil, err
	}
	wlog.ByCtx(ctx, "CreateDeck").Infof("deck initialed, %+v", deckCol.ID)

	return &Deck{Path: fPath, Storage: storage}, nil
}

func (deck *Deck) Col() *Col {
	deckCol, err := findOrMockSimpleDeck(deck.DB)
	if err == nil {
		return nil
	}
	return deckCol
}

func (deck *Deck) CardService() *CardService {
	return &CardService{DB: deck.DB}
}

func (deck *Deck) NoteService() *NoteService {
	return &NoteService{DB: deck.DB}
}
