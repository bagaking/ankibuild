package apkg

import (
	"context"
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/bagaking/goulp/wlog"
)

const (
	ApkgDBName = "collection.anki2"
)

// PkgInfo Anki 的 apkg 管理实体
// Ref:
// - official: https://docs.ankiweb.net/intro.html
// - open source:
// -- https://gist.github.com/sartak/3921255
// -- https://github.com/SergioFacchini/anki-cards-web-browser/blob/master/documentation/Processing%20Anki's%20.apkg%20files.md
type PkgInfo struct {
	Path string
	DB   *gorm.DB
}

// CreatePkgInfo initializes and returns a new PkgInfo with a connected database.
func CreatePkgInfo(ctx context.Context, path string) (*PkgInfo, error) {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create directory: %v", err)
	}

	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("%s/%s", path, ApkgDBName)), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %v", err)
	}

	// Migrate the schema
	err = AutoMigrateModels(db)
	if err != nil {
		return nil, fmt.Errorf("failed to automigrate: %v", err)
	}

	pkgInfo := &PkgInfo{Path: path, DB: db}
	col, err := pkgInfo.FindOrCreateSimpleDeck()
	if err != nil {
		defer pkgInfo.Close()
		return nil, err
	}
	wlog.ByCtx(ctx, "CreatePkgInfo").Infof("deck initialed, %+v", col)

	return pkgInfo, nil
}

// AutoMigrateModels migrates the provided models and creates necessary indexes.
func AutoMigrateModels(db *gorm.DB) error {
	err := db.AutoMigrate(&Col{}, &Note{}, &Card{}, &Revlog{}, &Grave{})
	if err != nil {
		return err
	}

	return CreateIndexes(db)
}

func (p *PkgInfo) Close() error {
	sqlDB, err := p.DB.DB()

	if err != nil {
		return fmt.Errorf("failed to get generic database object: %v", err)
	}

	return sqlDB.Close()
}
