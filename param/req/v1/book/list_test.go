package book

import (
	"math"
	"testing"

	"github.com/gin-gonic/gin/binding"
)

func TestListReqPlanetID(t *testing.T) {
	for _, tc := range []struct {
		body string
		want int64
	}{
		{`{"planet":{"id":"9007199254740993"}}`, 9007199254740993},
		{`{"planet":{"id":"9223372036854775807"}}`, math.MaxInt64},
	} {
		var req ListReq
		if err := binding.JSON.BindBody([]byte(tc.body), &req); err != nil {
			t.Fatalf("bind %s: %v", tc.body, err)
		}
		if req.Planet.ID != tc.want {
			t.Errorf("got planet ID %d, want %d", req.Planet.ID, tc.want)
		}
	}
}

func TestListReqRejectsInvalidPlanetID(t *testing.T) {
	for _, body := range []string{
		`{"planet":{"id":1}}`,
		`{"planet":{"id":1.5}}`,
		`{"planet":{"id":9007199254740993}}`,
		`{}`,
		`{"planet":null}`,
		`{"planet":{}}`,
		`{"planet":{"id":null}}`,
		`{"planet":{"id":""}}`,
		`{"planet":{"id":"0"}}`,
		`{"planet":{"id":"-1"}}`,
		`{"planet":{"id":"0x10"}}`,
		`{"planet":{"id":"1e3"}}`,
		`{"planet":{"id":"9223372036854775808"}}`,
	} {
		var req ListReq
		if err := binding.JSON.BindBody([]byte(body), &req); err == nil {
			t.Errorf("unexpectedly accepted %s", body)
		}
	}
}
