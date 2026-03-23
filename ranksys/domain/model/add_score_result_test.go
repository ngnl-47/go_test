package model

import (
	"encoding/json"
	"testing"
)

func TestAddScoreResultJSON(t *testing.T) {
	var r AddScoreResult
	if err := json.Unmarshal([]byte(`{"ok":true,"rank":2,"score":12.5}`), &r); err != nil {
		t.Fatal(err)
	}
	if !r.Ok || r.Rank != 2 || r.Score != 12.5 {
		t.Fatalf("got %+v", r)
	}
}
