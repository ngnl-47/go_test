package service_test

import (
	"context"
	"testing"
	"time"

	appsvc "go_test/ranksys/app/service"
	"go_test/ranksys/domain/facade"
	"go_test/ranksys/domain/model"
	domainsvc "go_test/ranksys/domain/service"
)

type captureRankRepo struct {
	lastAdd    model.AddScoreInput
	lastAnchor model.RankScoreAnchorInput
	lastRange  model.RankScoreRangeInput
	lastIDs    model.RanksByPlayerIDsInput
	lastMyPage model.MyPageQuery
}

func (c *captureRankRepo) UploadScore(_ context.Context, _ model.RankStorageKeyTriplet, _ model.UploadScoreInput) (model.UploadScoreResult, error) {
	return model.UploadScoreResult{}, nil
}

func (c *captureRankRepo) GetRanks(_ context.Context, _ model.RankStorageKeyTriplet, _ model.RankRangeQuery) ([]*model.RankRow, *model.RankRow, error) {
	return nil, nil, nil
}

func (c *captureRankRepo) RemovePlayer(_ context.Context, _ model.RankStorageKeyTriplet, _ string) error {
	return nil
}

func (c *captureRankRepo) Clear(_ context.Context, _ model.RankStorageKeyTriplet) error {
	return nil
}

func (c *captureRankRepo) AddScore(_ context.Context, _ model.RankStorageKeyTriplet, in model.AddScoreInput) (model.AddScoreResult, error) {
	c.lastAdd = in
	return model.AddScoreResult{Ok: true}, nil
}

func (c *captureRankRepo) GetRanksByScoreAnchor(_ context.Context, _ model.RankStorageKeyTriplet, in model.RankScoreAnchorInput) ([]*model.RankRow, error) {
	c.lastAnchor = in
	return nil, nil
}

func (c *captureRankRepo) GetRanksByScoreRange(_ context.Context, _ model.RankStorageKeyTriplet, in model.RankScoreRangeInput) ([]*model.RankRow, error) {
	c.lastRange = in
	return nil, nil
}

func (c *captureRankRepo) GetRanksByPlayerIDs(_ context.Context, _ model.RankStorageKeyTriplet, in model.RanksByPlayerIDsInput) ([]*model.RankRow, error) {
	c.lastIDs = in
	return nil, nil
}

func (c *captureRankRepo) GetMyPage(_ context.Context, _ model.RankStorageKeyTriplet, q model.MyPageQuery) ([]*model.RankRow, error) {
	c.lastMyPage = q
	return nil, nil
}

var _ facade.RankBoardRepository = (*captureRankRepo)(nil)

func TestRankBoardAppService_AddScore_DisplayDeltaToRedis(t *testing.T) {
	ctx := context.Background()
	kb, err := domainsvc.NewRankStorageKeyBuilder("cap")
	if err != nil {
		t.Fatal(err)
	}
	cap := &captureRankRepo{}
	app := appsvc.NewRankBoardAppService(kb, cap)
	board, _ := model.ParseBoardName("b1")
	now := time.Unix(1_700_000_000, 0)

	_, err = app.AddScore(ctx, board, appsvc.AddScoreRequest{
		PlayerKey: "p1",
		Sort:      model.SortDesc,
		ScoreAdd:  7,
	}, now)
	if err != nil {
		t.Fatal(err)
	}
	if cap.lastAdd.RedisScoreDelta != -7 {
		t.Fatalf("desc delta: got %v want -7", cap.lastAdd.RedisScoreDelta)
	}

	_, err = app.AddScore(ctx, board, appsvc.AddScoreRequest{
		PlayerKey: "p1",
		Sort:      model.SortAsc,
		ScoreAdd:  7,
	}, now)
	if err != nil {
		t.Fatal(err)
	}
	if cap.lastAdd.RedisScoreDelta != 7 {
		t.Fatalf("asc delta: got %v want 7", cap.lastAdd.RedisScoreDelta)
	}
}

func TestRankBoardAppService_GetRanksByScoreAnchor_RedisAnchor(t *testing.T) {
	ctx := context.Background()
	kb, err := domainsvc.NewRankStorageKeyBuilder("cap")
	if err != nil {
		t.Fatal(err)
	}
	cap := &captureRankRepo{}
	app := appsvc.NewRankBoardAppService(kb, cap)
	board, _ := model.ParseBoardName("b1")

	_, err = app.GetRanksByScoreAnchor(ctx, board, 100, 3, 4, true, model.SortDesc)
	if err != nil {
		t.Fatal(err)
	}
	if cap.lastAnchor.RedisAnchorScore != 100 || cap.lastAnchor.High != 3 || cap.lastAnchor.Low != 4 || !cap.lastAnchor.IncludeInfo {
		t.Fatalf("anchor: %+v", cap.lastAnchor)
	}
}

func TestRankBoardAppService_GetRanksByScoreRange_RedisRange(t *testing.T) {
	ctx := context.Background()
	kb, err := domainsvc.NewRankStorageKeyBuilder("cap")
	if err != nil {
		t.Fatal(err)
	}
	cap := &captureRankRepo{}
	app := appsvc.NewRankBoardAppService(kb, cap)
	board, _ := model.ParseBoardName("b1")

	_, err = app.GetRanksByScoreRange(ctx, board, 50, 100, 0, 0, true, model.SortDesc)
	if err != nil {
		t.Fatal(err)
	}
	if cap.lastRange.MinRedisScore != 50 || cap.lastRange.MaxRedisScore != 100 || !cap.lastRange.IncludeInfo {
		t.Fatalf("score range: %+v", cap.lastRange)
	}
}

func TestRankBoardAppService_GetRanksByPlayerIDs_AndGetMyPage(t *testing.T) {
	ctx := context.Background()
	kb, err := domainsvc.NewRankStorageKeyBuilder("cap")
	if err != nil {
		t.Fatal(err)
	}
	cap := &captureRankRepo{}
	app := appsvc.NewRankBoardAppService(kb, cap)
	board, _ := model.ParseBoardName("b1")

	_, err = app.GetRanksByPlayerIDs(ctx, board, model.RanksByPlayerIDsInput{
		PlayerKeys:  []string{"a", "b"},
		IncludeInfo: true,
		IgnoreNil:   false,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(cap.lastIDs.PlayerKeys) != 2 || cap.lastIDs.PlayerKeys[0] != "a" || !cap.lastIDs.IncludeInfo {
		t.Fatalf("player ids: %+v", cap.lastIDs)
	}

	_, err = app.GetMyPage(ctx, board, model.MyPageQuery{MyPlayerKey: "me", PageSize: 10, IncludeInfo: false})
	if err != nil {
		t.Fatal(err)
	}
	if cap.lastMyPage.MyPlayerKey != "me" || cap.lastMyPage.PageSize != 10 || cap.lastMyPage.IncludeInfo {
		t.Fatalf("my page: %+v", cap.lastMyPage)
	}
}
