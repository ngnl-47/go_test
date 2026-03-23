package service

import (
	"testing"

	"go_test/ranksys/domain/model"
)

func TestRankStorageKeyBuilder_Build_ClusterTagShared(t *testing.T) {
	b := MustNewRankStorageKeyBuilder("xyjh")
	name, _ := model.ParseBoardName("rank:server:level:{1}")
	k := b.Build(name)

	r, m, i := k.RankZSetKey(), k.IDToZMemberMapKey(), k.InfoHashKey()
	if r == "" || m == "" || i == "" {
		t.Fatal("empty key")
	}
	want := "xyjh:leaderboard:{rank:server:level:{1}}"
	if r != want+":rank" || m != want+":map" || i != want+":info" {
		t.Fatalf("unexpected keys:\n%s\n%s\n%s", r, m, i)
	}
	if k.BaseKey() != want {
		t.Fatalf("BaseKey: got %q want %q", k.BaseKey(), want)
	}
}
