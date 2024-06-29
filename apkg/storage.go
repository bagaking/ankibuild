package apkg

import (
	"context"
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/bagaking/goulp/wlog"
	"github.com/khicago/irr"
)

type Storage struct {
	DB *gorm.DB
}

// CreateStorage initializes and returns a new Storage with a connected database.
func CreateStorage(ctx context.Context, fPath, dbName string, initial func(db *gorm.DB) error) (*Storage, error) {
	log := wlog.ByCtx(ctx, "CreateStorage")
	err := os.MkdirAll(fPath, 0o755)
	if err != nil {
		return nil, fmt.Errorf("failed to create directory: %v", err)
	}

	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("%s/%s", fPath, dbName)), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %v", err)
	}
	storage := &Storage{DB: db}

	// Migrate the schema
	if initial != nil {
		err = initial(db)
		if err != nil {
			if e := storage.Close(); e != nil {
				log.WithError(e).Warnf("after migration failed, failed to close storage")
			}
			return nil, irr.Wrap(err, "failed to auto migrate")
		}
	}

	return storage, nil
}

func CreateInMemStorage(ctx context.Context, initial func(db *gorm.DB) error) (*Storage, error) {
	log := wlog.ByCtx(ctx, "CreateInMemStorage")
	db, err := gorm.Open(sqlite.Open("file:memdb1?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to inmemory database: %v", err)
	}
	storage := &Storage{DB: db}

	// Migrate the schema
	if initial != nil {
		err = initial(db)
		if err != nil {
			if e := storage.Close(); e != nil {
				log.WithError(e).Warnf("after migration failed, failed to close storage")
			}
			return nil, irr.Wrap(err, "failed to auto migrate")
		}
	}

	return storage, nil
}

// 暂时没有特别好的解法
//func (p *Storage) Dump(dest io.Writer) error {
//	sqlDB, err := p.DB.DB()
//	if err != nil {
//		return irr.Wrap(err, "unable to get database from GORM")
//	}
//
//	return &out, nil
//}

func (p *Storage) Close() error {
	sqlDB, err := p.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get generic database object: %v", err)
	}

	return sqlDB.Close()
}
