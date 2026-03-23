package repo

import (
	"testing"

	"go_test/ranksys/domain/model"
)

func TestRowsFromLuaRankObjectJSON_ScoreAdjust(t *testing.T) {
	raw := `{"1":{"key":"1","score":-99,"rank":1,"info":"{\"n\":\"a\"}"},"2":{"key":"2","score":-50,"rank":2}}`
	rows, err := rowsFromLuaRankObjectJSON(raw, func(s float64) float64 {
		return -s // asc-style flip
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 || rows[0].Rank != 1 || rows[0].PlayerKey != "1" || rows[0].Score != 99 {
		t.Fatalf("row0: %+v", rows[0])
	}
	if rows[1].Score != 50 {
		t.Fatalf("row1 score: %v", rows[1].Score)
	}
}

func TestRowsFromScoreQueryJSON_Desc(t *testing.T) {
	raw := `{"u1":{"key":"u1","score":-10,"rank":1}}`
	rows, err := rowsFromScoreQueryJSON(raw, model.SortDesc)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Score != -10 {
		t.Fatalf("got %+v", rows[0])
	}
}

func TestParseMyPageJSON_Empty(t *testing.T) {
	rows, err := parseMyPageJSON("[]")
	if err != nil || rows != nil {
		t.Fatalf("got %v %v", rows, err)
	}
}

func TestParseMyPageJSON_Object(t *testing.T) {
	raw := `{"9":{"key":"9","score":100,"rank":2,"info":false}}`
	rows, err := parseMyPageJSON(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].PlayerKey != "9" || rows[0].Rank != 2 {
		t.Fatalf("got %+v", rows[0])
	}
}

func TestParseRanksByPlayerIDsJSON_PreserveNil(t *testing.T) {
	raw := `[null,{"roleIdStr":"2","score":10,"rank":1}]`
	rows, err := parseRanksByPlayerIDsJSON(raw, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 || rows[0] != nil || rows[1].PlayerKey != "2" {
		t.Fatalf("got %#v", rows)
	}
}

func TestParseRanksByPlayerIDsJSON_IgnoreNil(t *testing.T) {
	raw := `[null,{"roleIdStr":"2","score":10,"rank":1}]`
	rows, err := parseRanksByPlayerIDsJSON(raw, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].PlayerKey != "2" {
		t.Fatalf("got %#v", rows)
	}
}
