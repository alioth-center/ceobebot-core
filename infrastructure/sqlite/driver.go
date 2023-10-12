package sqlite

import (
	"fmt"
	"github.com/ceobebot/qqchannel/infrastructure/log"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"os"
	"path/filepath"
)

type Database interface {
	init(dataSource string) error
	migrate(models ...any) error
	Has(table string, query string, args ...any) (exist bool, err error)
	Count(table string, query string, args ...any) (count int64, err error)
	GetOne(receiver any, query string, args ...any) error
	GetAll(receiver any, query string, args ...any) error
	PickOne(receiver any, query string, args ...any) error
	PickAll(receiver any, length int, query string, args ...any) error
	InsertOne(data any) error
	InsertAll(data any) error
	UpdateOne(data any, query string, args ...any) error
	UpdateAll(data any, query string, args ...any) error
	DeleteOne(query string, args ...any) error
	DeleteAll(query string, args ...any) error
	ExecRaw(sql string, args ...any) error
	QueryRaw(receiver any, sql string, args ...any) error
}

type sqliteDb struct {
	db *gorm.DB
}

func (s *sqliteDb) init(dataSource string) error {
	if db, openErr := gorm.Open(sqlite.Open(filepath.Join("data", dataSource))); openErr != nil {
		if !os.IsNotExist(openErr) {
			err := fmt.Errorf("open sqliteDb database error: %w", openErr)
			log.Error(log.NewFieldsWithError(err).With("data_source", dataSource))
			return err
		} else if _, createErr := os.Create(filepath.Join("data", dataSource)); createErr != nil {
			err := fmt.Errorf("create sqliteDb database error: %w", createErr)
			log.Error(log.NewFieldsWithError(err).With("data_source", dataSource))
			return err
		} else {
			return s.init(dataSource)
		}
	} else {
		s.db = db
		log.Info(log.NewFieldsWithMessage("successfully open sqliteDb database").With("data_source", dataSource))
		return nil
	}
}

func (s *sqliteDb) migrate(models ...any) error {
	return s.db.AutoMigrate(models...)
}

func (s *sqliteDb) exec(command func(tx *gorm.DB) *gorm.DB) error {
	sql := s.db.ToSQL(command)
	log.Debug(log.NewFieldsWithMessage("sql executed").With("sql", sql))
	err := command(s.db).Error
	if err != nil {
		log.Error(log.NewFieldsWithError(err).With("sql", sql).With("message", "sql execution failed"))
	}
	return err
}

func (s *sqliteDb) Has(table, query string, args ...any) (exist bool, err error) {
	var count int64
	err = s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.Table(table).Where(query, args...).Limit(100).Count(&count)
	})
	return count > 0, err
}

func (s *sqliteDb) Count(table, query string, args ...any) (count int64, err error) {
	err = s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.Table(table).Where(query, args...).Count(&count)
	})
	return count, err
}

func (s *sqliteDb) GetOne(receiver any, query string, args ...any) error {
	return s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.First(receiver, append([]any{query}, args...)...)
	})
}

func (s *sqliteDb) GetAll(receiver any, query string, args ...any) error {
	return s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.Find(receiver, append([]any{query}, args...)...)
	})
}

func (s *sqliteDb) PickOne(receiver any, query string, args ...any) error {
	return s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.Order(gorm.Expr("random()")).Take(receiver, append([]any{query}, args...)...)
	})
}

func (s *sqliteDb) PickAll(receiver any, length int, query string, args ...any) error {
	return s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.Order(gorm.Expr("random()")).Limit(length).Find(receiver, append([]any{query}, args...)...)
	})
}

func (s *sqliteDb) InsertOne(data any) error {
	return s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.Create(data)
	})
}

func (s *sqliteDb) InsertAll(data any) error {
	return s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.CreateInBatches(data, 100)
	})
}

func (s *sqliteDb) UpdateOne(data any, query string, args ...any) error {
	return s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(data).Where(query, args...).Limit(1).Updates(data)
	})
}

func (s *sqliteDb) UpdateAll(data any, query string, args ...any) error {
	return s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(data).Where(query, args...).Updates(data)
	})
}

func (s *sqliteDb) DeleteOne(query string, args ...any) error {
	return s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.Delete(query, args...).Limit(1)
	})
}

func (s *sqliteDb) DeleteAll(query string, args ...any) error {
	return s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.Delete(query, args...)
	})
}

func (s *sqliteDb) ExecRaw(sql string, args ...any) error {
	return s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.Exec(sql, args...)
	})
}

func (s *sqliteDb) QueryRaw(receiver any, sql string, args ...any) error {
	return s.exec(func(tx *gorm.DB) *gorm.DB {
		return tx.Raw(sql, args...).Scan(receiver)
	})
}

func NewSqliteDb(path string, models ...any) (db Database, err error) {
	sqliteDb := &sqliteDb{}
	if initErr := sqliteDb.init(path); initErr != nil {
		return nil, fmt.Errorf("init sqliteDb database error: %w", initErr)
	} else if migrateErr := sqliteDb.migrate(models...); migrateErr != nil {
		return nil, fmt.Errorf("migrate sqliteDb database error: %w", migrateErr)
	} else {
		return sqliteDb, nil
	}
}
