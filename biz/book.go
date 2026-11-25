package biz

import (
	"context"
	"errors"

	"SolarNotesAPI/model"
	"SolarNotesAPI/sys_error"
)

type Book struct {
	ID    int64
	Title string
}

// List 保留 DAO 提供的展示顺序, 将数据访问错误转换为接口可识别的类别.
func (Book) List(ctx context.Context, planetID int64) ([]Book, error) {
	if ctx == nil || planetID <= 0 {
		return nil, &sys_error.InternalError{Cause: errors.New("biz.Book.List: invalid context or planet ID")}
	}

	var dao model.Book
	rows, err := dao.ListVisibleByPlanet(ctx, planetID)
	if errors.Is(err, model.PlanetNotFoundErr) {
		return nil, sys_error.PlanetNotFoundErr
	}
	if err != nil {
		return nil, &sys_error.DatabaseError{Cause: err}
	}

	books := make([]Book, 0, len(rows))
	for _, row := range rows {
		books = append(books, Book{ID: row.ID, Title: row.Title})
	}
	return books, nil
}
