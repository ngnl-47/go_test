package service

import (
	"testing"

	"go_test/ranksys/domain/model"
)

func TestDisplayScoreFromRedis(t *testing.T) {
	if g, w := DisplayScoreFromRedis(model.SortDesc, -42), float64(-42); g != w {
		t.Fatalf("desc: got %v want %v", g, w)
	}
	if g, w := DisplayScoreFromRedis(model.SortAsc, 42), float64(-42); g != w {
		t.Fatalf("asc: got %v want %v", g, w)
	}
}

func TestRedisAnchorScore(t *testing.T) {
	if g, w := RedisAnchorScore(model.SortDesc, 100), float64(100); g != w {
		t.Fatalf("desc anchor: got %v want %v", g, w)
	}
	if g, w := RedisAnchorScore(model.SortAsc, 100), float64(-100); g != w {
		t.Fatalf("asc anchor: got %v want %v", g, w)
	}
}

func TestRedisScoreRange(t *testing.T) {
	min, max, ok := RedisScoreRange(model.SortDesc, 10, 100)
	if !ok || min != 10 || max != 100 {
		t.Fatalf("desc range: %v %v %v", min, max, ok)
	}
	min, max, ok = RedisScoreRange(model.SortAsc, 10, 100)
	if !ok || min != -100 || max != -10 {
		t.Fatalf("asc range: %v %v %v", min, max, ok)
	}
	_, _, ok = RedisScoreRange(model.SortAsc, 100, 10)
	if ok {
		t.Fatal("invalid range should be !ok")
	}
}
