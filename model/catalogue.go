package model

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"SolarNotesAPI/db"

	"gorm.io/gorm"
)

var BookNotFoundErr = errors.New("visible book not found")

type Catalogue struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement:false"`
	BookID        int64     `gorm:"column:book_id"`
	ParentID      int64     `gorm:"column:parent_id"`
	Type          string    `gorm:"column:type"`
	Name          string    `gorm:"column:name"`
	Intro         string    `gorm:"column:intro"`
	Sort          *int32    `gorm:"column:sort"`
	DisplayStatus string    `gorm:"column:display_status"`
	CreatedAt     time.Time `gorm:"column:created_at"`
}

// ReadBookTree reads the visible book and all its nodes from one snapshot.
func (Catalogue) ReadBookTree(ctx context.Context, bookID int64) ([]Catalogue, error) {
	if db.DB == nil {
		return nil, errors.New("database connection is not initialized")
	}

	var rows []Catalogue
	opts := &sql.TxOptions{Isolation: sql.LevelRepeatableRead, ReadOnly: true}
	err := db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var book Book
		err := tx.Select("id").Where("id = ? AND display_status = ?", bookID, "visible").Take(&book).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return BookNotFoundErr
		}
		if err != nil {
			return err
		}

		return tx.Select("id", "book_id", "parent_id", "type", "name", "intro", "sort", "display_status", "created_at").
			Where("book_id = ?", bookID).
			Order("sort ASC").
			Find(&rows).Error
	}, opts)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
