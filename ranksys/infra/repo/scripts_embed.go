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
