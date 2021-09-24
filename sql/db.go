package sql

import (
	"errors"
	"time"

	"github.com/doug-martin/goqu/v9"
	"gorm.io/gorm"
)

// DefaultDB ..
var DefaultDB *DataBase

// Gorm get gorm db instance
func Gorm() *gorm.DB {
	return DefaultDB.Gorm()
}

// Goqu get gorm db instance
func Goqu() *goqu.Database {
	return DefaultDB.Goqu()
}

// IsNotFound ..
func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// Config ..
type Config struct {
	Dialect         Dialect
	URL             string
	TransTimeout    time.Duration
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	Debug           bool
}

// Open get opened db instance
func Open(cfg *Config, opt ...gorm.Option) (*DataBase, error) {
	db, err := NewGorm(cfg.Dialect, cfg.URL, opt...)
	if err != nil {
		return nil, err
	}

	if cfg.Debug {
		db = db.Debug()
	}

	conn, err := db.DB()
	if err != nil {
		return nil, err
	}

	if cfg.MaxOpenConns != 0 {
		conn.SetMaxOpenConns(cfg.MaxOpenConns)
	}

	if cfg.MaxIdleConns != 0 {
		conn.SetMaxIdleConns(cfg.MaxIdleConns)
	}

	if cfg.ConnMaxLifetime != 0 {
		conn.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}

	goquDB, err := NewGoqu(cfg.Dialect, conn)
	if err != nil {
		return nil, err
	}
	return &DataBase{
		DB:   db,
		cfg:  cfg,
		goqu: goquDB,
	}, nil
}

// DataBase ...
type DataBase struct {
	*gorm.DB
	cfg  *Config
	goqu *goqu.Database
}

// Gorm ...
func (db *DataBase) Gorm() *gorm.DB {
	return db.DB
}

// Goqu ...
func (db *DataBase) Goqu() *goqu.Database {
	return db.goqu
}

// Close ...
func (db *DataBase) Close() {
	conn, err := db.Gorm().DB()
	if err == nil {
		conn.Close()
	}
}
