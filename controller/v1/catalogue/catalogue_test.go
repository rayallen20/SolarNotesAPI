package catalogue

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"SolarNotesAPI/biz"
	reqArticle "SolarNotesAPI/param/req/v1/article"

	"github.com/gin-gonic/gin/binding"
)

func TestCatalogueJSONAndArticleID(t *testing.T) {
	created := time.Date(2026, 1, 2, 0, 0, 0, 0, time.FixedZone("UTC+8", 8*60*60))
	root := &biz.CatalogueNode{
		ID: 9007199254740996, Type: "folder", Name: "Root", CreatedAt: created,
		Children: []*biz.CatalogueNode{
			{ID: 9007199254740997, Type: "file", Name: "Page"},
			{ID: 9007199254740998, Type: "folder", Name: "Empty", CreatedAt: created},
		},
	}
	raw, err := json.Marshal(toResponse(root))
	if err != nil {
		t.Fatal(err)
	}
	var got, want map[string]any
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	expected := `{
		"id":"9007199254740996","type":"folder","name":"Root","intro":"","createdAt":"2026-01-02",
		"children":[
			{"id":"9007199254740997","type":"file","name":"Page"},
			{"id":"9007199254740998","type":"folder","name":"Empty","intro":"","createdAt":"2026-01-02","children":[]}
		]
	}`
	if err := json.Unmarshal([]byte(expected), &want); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected public JSON: %s", raw)
	}

	leafID := got["children"].([]any)[0].(map[string]any)["id"]
	body, err := json.Marshal(map[string]any{"catalogue": map[string]any{"id": leafID}})
	if err != nil {
		t.Fatal(err)
	}
	var req reqArticle.ShowReq
	if err := binding.JSON.BindBody(body, &req); err != nil || req.Catalogue == nil || req.Catalogue.ID != root.Children[0].ID {
		t.Fatalf("article request did not preserve the leaf ID: %s, error: %v", body, err)
	}
	var invalid reqArticle.ShowReq
	if err := binding.JSON.BindBody([]byte(`{"catalogue":{"id":9007199254740997}}`), &invalid); err == nil {
		t.Fatal("article request accepted a numeric ID")
	}
}
