package biz

import (
	"context"
	"errors"
	"time"

	"SolarNotesAPI/articlecontent"
	"SolarNotesAPI/model"
	"SolarNotesAPI/sys_error"
)

type Article struct {
	ID        int64
	Content   string
	CreatedAt time.Time
}

// ReadArticle accesses storage only after the requested file is found in the public tree.
func ReadArticle(ctx context.Context, catalogueID int64, readMarkdown func(string) (string, error)) (*Article, error) {
	if ctx == nil || catalogueID <= 0 || readMarkdown == nil {
		return nil, &sys_error.InternalError{Cause: errors.New("invalid article read arguments")}
	}
	row, nodes, err := model.ReadArticleSnapshot(ctx, catalogueID)
	if errors.Is(err, model.ArticleNotFoundErr) {
		return nil, sys_error.ArticleUnavailableErr
	}
	if err != nil {
		return nil, &sys_error.DatabaseError{Cause: err}
	}
	root, err := buildCatalogueTree(row.BookID, nodes)
	if errors.Is(err, sys_error.CatalogueNotFoundErr) {
		return nil, sys_error.ArticleUnavailableErr
	}
	if err != nil {
		return nil, &sys_error.InternalError{Cause: err}
	}
	if !containsPublicFile(root, catalogueID) || row.Path == "" {
		return nil, sys_error.ArticleUnavailableErr
	}
	if row.ID <= 0 || row.CreatedAt.IsZero() {
		return nil, &sys_error.InternalError{Cause: errors.New("invalid article metadata")}
	}
	content, err := readMarkdown(row.Path)
	if errors.Is(err, articlecontent.ErrUnavailable) {
		return nil, sys_error.ArticleUnavailableErr
	}
	if err != nil {
		return nil, &sys_error.InternalError{Cause: err}
	}
	return &Article{ID: row.ID, Content: content, CreatedAt: row.CreatedAt}, nil
}

func containsPublicFile(root *CatalogueNode, id int64) bool {
	stack := []*CatalogueNode{root}
	for len(stack) > 0 {
		last := len(stack) - 1
		node := stack[last]
		stack = stack[:last]
		if node.ID == id {
			return node.Type == "file"
		}
		stack = append(stack, node.Children...)
	}
	return false
}
