package biz

import (
	"context"
	"errors"
	"fmt"
	"time"

	"SolarNotesAPI/model"
	"SolarNotesAPI/sys_error"
)

type CatalogueNode struct {
	ID        int64
	Type      string
	Name      string
	Intro     string
	CreatedAt time.Time
	Children  []*CatalogueNode
}

// ReadCatalogue reads the public catalogue rooted at the given book.
func ReadCatalogue(ctx context.Context, bookID int64) (*CatalogueNode, error) {
	if ctx == nil || bookID <= 0 {
		return nil, errors.New("biz.ReadCatalogue: invalid context or book ID")
	}
	var dao model.Catalogue
	rows, err := dao.ReadBookTree(ctx, bookID)
	if errors.Is(err, model.BookNotFoundErr) {
		return nil, sys_error.BookNotFoundErr
	}
	if err != nil {
		return nil, &sys_error.DatabaseError{Cause: err}
	}
	return buildCatalogueTree(bookID, rows)
}

// buildCatalogueTree expects the DAO's ascending sort order. Filtering by
// parent preserves that order within each set of siblings.
func buildCatalogueTree(bookID int64, rows []model.Catalogue) (*CatalogueNode, error) {
	if len(rows) == 0 {
		return nil, sys_error.CatalogueNotFoundErr
	}

	byID := make(map[int64]*model.Catalogue, len(rows))
	children := make(map[int64][]*model.Catalogue)
	var root *model.Catalogue
	for i := range rows {
		row := &rows[i]
		if row.ID <= 0 || row.BookID != bookID || row.ParentID < 0 {
			return nil, fmt.Errorf("book %d: invalid catalogue identity at node %d", bookID, row.ID)
		}
		if byID[row.ID] != nil {
			return nil, fmt.Errorf("book %d: duplicate node %d", bookID, row.ID)
		}
		if row.Type != "folder" && row.Type != "file" {
			return nil, fmt.Errorf("node %d: unknown type %q", row.ID, row.Type)
		}
		if row.DisplayStatus != "visible" && row.DisplayStatus != "hidden" {
			return nil, fmt.Errorf("node %d: unknown display status %q", row.ID, row.DisplayStatus)
		}
		byID[row.ID] = row
		if row.ParentID == 0 {
			if root != nil || row.Type != "folder" || row.Sort == nil || *row.Sort != 0 {
				return nil, fmt.Errorf("book %d: invalid or duplicate root %d", bookID, row.ID)
			}
			root = row
			continue
		}
		if (row.Sort == nil && row.DisplayStatus != "hidden") || (row.Sort != nil && *row.Sort <= 0) {
			return nil, fmt.Errorf("node %d: invalid non-root sort", row.ID)
		}
		siblings := children[row.ParentID]
		if row.Sort != nil && len(siblings) > 0 {
			previous := siblings[len(siblings)-1].Sort
			if previous != nil && *previous == *row.Sort {
				return nil, fmt.Errorf("parent %d: duplicate sort %d", row.ParentID, *row.Sort)
			}
		}
		children[row.ParentID] = append(siblings, row)
	}
	if root == nil {
		return nil, fmt.Errorf("book %d: catalogue root missing", bookID)
	}
	for parentID := range children {
		if parent := byID[parentID]; parent == nil || parent.Type != "folder" {
			return nil, fmt.Errorf("book %d: missing or non-folder parent %d", bookID, parentID)
		}
	}

	// Visit even unpublished subtrees so that reachability covers every row.
	type visit struct {
		row    *model.Catalogue
		parent *CatalogueNode
	}
	queue := []visit{{row: root}}
	var result *CatalogueNode
	for i := 0; i < len(queue); i++ {
		item := queue[i]
		row := item.row
		var node *CatalogueNode
		publish := row.DisplayStatus == "visible" && !(row.Type == "folder" && row.Name == "img")
		if publish && (row == root || item.parent != nil) {
			node = &CatalogueNode{ID: row.ID, Type: row.Type, Name: row.Name}
			if row.Type == "folder" {
				node.Intro = row.Intro
				node.CreatedAt = row.CreatedAt
				node.Children = make([]*CatalogueNode, 0)
			}
			if row == root {
				result = node
			} else {
				item.parent.Children = append(item.parent.Children, node)
			}
		}
		for _, child := range children[row.ID] {
			queue = append(queue, visit{row: child, parent: node})
		}
	}
	if len(queue) != len(rows) {
		return nil, fmt.Errorf("book %d: unreachable catalogue cycle", bookID)
	}
	if result == nil {
		return nil, sys_error.CatalogueNotFoundErr
	}
	return result, nil
}
