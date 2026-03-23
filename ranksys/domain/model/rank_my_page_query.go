package model

// MyPageQuery get_my_page.lua：按 PageSize 分页，返回「我」所在页内的条目（与 mandarin GetMyPageArgs.PageSize 一致）。
type MyPageQuery struct {
	MyPlayerKey string
	PageSize    int
	IncludeInfo bool
}
