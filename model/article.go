package model

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"SolarNotesAPI/db"

	"gorm.io/gorm"
)

var ArticleNotFoundErr = errors.New("article not found")

// Article holds the article metadata and its catalogue node's owning book ID.
type Article struct {
	ID        int64     `gorm:"column:id"`
	BookID    int64     `gorm:"column:book_id"`
	Path      string    `gorm:"column:path"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

// ReadArticleSnapshot reads metadata and the book's tree in a single snapshot.
func ReadArticleSnapshot(ctx context.Context, catalogueID int64) (Article, []Catalogue, error) {
	if db.DB == nil {
		return Article{}, nil, errors.New("database connection is not initialized")
	}
	var row Article
	var nodes []Catalogue
	opts := &sql.TxOptions{Isolation: sql.LevelRepeatableRead, ReadOnly: true}
	err := db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Table("article").
			Select("article.id, article.path, article.created_at, catalogue.book_id").
			Joins("JOIN catalogue ON catalogue.id = article.catalogue_id").
			Joins("JOIN book ON book.id = catalogue.book_id").
			Where("catalogue.id = ? AND catalogue.type = ? AND catalogue.display_status = ? AND book.display_status = ?",
				catalogueID, "file", "visible", "visible").
			Take(&row).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ArticleNotFoundErr
		}
		if err != nil {
			return err
		}
		nodes, err = readCatalogueRows(tx, row.BookID)
		return err
	}, opts)
	if err != nil {
		return Article{}, nil, err
	}
	return row, nodes, nil
}
