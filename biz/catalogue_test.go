package biz

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"

	"SolarNotesAPI/model"
	"SolarNotesAPI/sys_error"
)

func catalogueTestRow(id, parentID int64, kind, name, status string, sort *int32) model.Catalogue {
	return model.Catalogue{
		ID: id, BookID: 7, ParentID: parentID, Type: kind, Name: name,
		DisplayStatus: status, Sort: sort,
	}
}

func TestBuildCatalogueTreePublicBranches(t *testing.T) {
	const firstID int64 = 9007199254740996
	const secondID int64 = firstID - 1
	zero, one, two, three, four := int32(0), int32(1), int32(2), int32(3), int32(4)
	root := catalogueTestRow(10, 0, "folder", "Original root", "visible", &zero)
	root.CreatedAt = time.Date(2026, 1, 2, 12, 30, 0, 0, time.UTC)
	root.Intro = "Root introduction"
	rows := []model.Catalogue{
		root,
		catalogueTestRow(11, 10, "folder", "Hidden", "hidden", &one),
		catalogueTestRow(12, 11, "folder", "Below hidden folder", "visible", &one),
		catalogueTestRow(15, 12, "file", "Hidden by grandparent", "visible", &one),
		catalogueTestRow(14, 13, "file", "Below img folder", "visible", &one),
		catalogueTestRow(13, 10, "folder", "img", "visible", &two),
		catalogueTestRow(firstID, 10, "folder", "IMG", "visible", &three),
		catalogueTestRow(secondID, 10, "file", "img", "visible", &four),
	}
	got, err := buildCatalogueTree(7, rows)
	if err != nil {
		t.Fatal(err)
	}
	want := &CatalogueNode{
		ID: root.ID, Type: "folder", Name: root.Name, Intro: root.Intro, CreatedAt: root.CreatedAt,
		Children: []*CatalogueNode{
			{ID: firstID, Type: "folder", Name: "IMG", Children: []*CatalogueNode{}},
			{ID: secondID, Type: "file", Name: "img"},
		},
	}
	if !reflect.DeepEqual(got, want) {
		tree, _ := json.Marshal(got)
		t.Errorf("unexpected public catalogue: %s", tree)
	}
}

func TestBuildCatalogueTreeWithoutPublicRoot(t *testing.T) {
	zero := int32(0)
	for name, rows := range map[string][]model.Catalogue{
		"no rows":     nil,
		"hidden root": {catalogueTestRow(1, 0, "folder", "Root", "hidden", &zero)},
		"img root":    {catalogueTestRow(1, 0, "folder", "img", "visible", &zero)},
	} {
		t.Run(name, func(t *testing.T) {
			got, err := buildCatalogueTree(7, rows)
			if got != nil || !errors.Is(err, sys_error.CatalogueNotFoundErr) {
				t.Errorf("expected missing public catalogue, got root %v, error %v", got, err)
			}
		})
	}
}

func TestBuildCatalogueTreeRejectsCorruptRows(t *testing.T) {
	zero, one := int32(0), int32(1)
	root := catalogueTestRow(1, 0, "folder", "Root", "visible", &zero)
	for name, rows := range map[string][]model.Catalogue{
		"disconnected hidden cycle": {
			root,
			catalogueTestRow(2, 3, "folder", "A", "hidden", &one),
			catalogueTestRow(3, 2, "folder", "B", "hidden", &one),
		},
		"visible node without sort": {
			catalogueTestRow(2, 1, "file", "Unsorted", "visible", nil),
			root,
		},
		"duplicate sibling sort": {
			root,
			catalogueTestRow(2, 1, "file", "A", "visible", &one),
			catalogueTestRow(3, 1, "file", "B", "hidden", &one),
		},
	} {
		t.Run(name, func(t *testing.T) {
			got, err := buildCatalogueTree(7, rows)
			if got != nil || err == nil || errors.Is(err, sys_error.CatalogueNotFoundErr) {
				t.Errorf("expected a tree integrity error, got root %v, error %v", got, err)
			}
		})
	}
}
