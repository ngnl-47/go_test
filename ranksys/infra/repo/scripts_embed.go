package repo

import _ "embed"

//go:embed lua/upload.lua
var scriptUpload string

//go:embed lua/get_rank.lua
var scriptGetRank string

//go:embed lua/remove.lua
var scriptRemove string

//go:embed lua/clear.lua
var scriptClear string

//go:embed lua/add.lua
var scriptAdd string

//go:embed lua/get_ranks_by_score_anchor.lua
var scriptGetRanksByScoreAnchor string

//go:embed lua/get_ranks_by_score_range.lua
var scriptGetRanksByScoreRange string

//go:embed lua/get_ranks_by_ids.lua
var scriptGetRanksByIDs string

//go:embed lua/get_my_page.lua
var scriptGetMyPage string
