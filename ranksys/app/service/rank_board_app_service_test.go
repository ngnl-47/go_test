package service

import (
	"strings"
	"testing"
	"time"

	"go_test/ranksys/domain/model"
)

func TestBuildUploadInput_Desc(t *testing.T) {
	svc := &RankBoardAppService{}
	now := time.Unix(1_700_000_000, 123_000_000)
	in, err := svc.buildUploadInput(UploadScoreRequest{
		PlayerID: 1001,
		Sort:     model.SortDesc,
		Score:    88.5,
		Info:     map[string]int{"x": 1},
	}, now)
	if err != nil {
		t.Fatal(err)
	}
	if in.PlayerKey != "1001" {
		t.Fatal(in.PlayerKey)
	}
	if in.RedisSortScore != -88.5 {
		t.Fatalf("redis score %v", in.RedisSortScore)
	}
	if len(in.PrefixedZMember) < 18 || in.PrefixedZMember[:1] != "-" {
		t.Fatalf("prefixed member %q", in.PrefixedZMember)
	}
	if !strings.HasSuffix(in.PrefixedZMember, ":1001") {
		t.Fatalf("suffix %q", in.PrefixedZMember)
	}
}

func TestBuildAddScoreInput_Asc(t *testing.T) {
	svc := &RankBoardAppService{}
	now := time.Unix(1_700_000_000, 0)
	in, err := svc.buildAddScoreInput(AddScoreRequest{
		PlayerKey: "p1",
		Sort:      model.SortAsc,
		ScoreAdd:  5,
	}, now)
	if err != nil {
		t.Fatal(err)
	}
	if in.RedisScoreDelta != 5 {
		t.Fatalf("delta %v", in.RedisScoreDelta)
	}
	if in.PrefixedZMember[0] != '+' {
		t.Fatalf("want + prefix got %q", in.PrefixedZMember)
	}
}
