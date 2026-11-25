package model

import (
	"context"
	"database/sql"
	"errors"

	"SolarNotesAPI/db"

	"gorm.io/gorm"
)

var PlanetNotFoundErr = errors.New("planet not found")

// Book 映射书单查询需要返回的两列.
type Book struct {
	ID    int64  `gorm:"column:id;primaryKey;autoIncrement:false"`
	Title string `gorm:"column:title"`
}

// ListVisibleByPlanet 在同一只读快照内检查天体并查询书单.
func (Book) ListVisibleByPlanet(ctx context.Context, planetID int64) ([]Book, error) {
	if db.DB == nil {
		return nil, errors.New("database connection is not initialized")
	}

	var books []Book
	opts := &sql.TxOptions{Isolation: sql.LevelRepeatableRead, ReadOnly: true}
	err := db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var planet Planet
		err := tx.Select("id").Where("id = ?", planetID).Take(&planet).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return PlanetNotFoundErr
		}
		if err != nil {
			return err
		}

		return tx.Select("id", "title").
			Where("planet_id = ? AND display_status = ?", planetID, "visible").
			Order("sort ASC").
			Find(&books).Error
	}, opts)
	if err != nil {
		return nil, err
	}
	return books, nil
}
