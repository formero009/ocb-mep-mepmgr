package common

//page struct
type Page struct {
	PageNo     int64       `json:"page_no"`
	PageSize   int64       `json:"page_size"`
	TotalPage  int64       `json:"total_page"`
	TotalCount int64       `json:"total_count"`
	List       interface{} `json:"list"`
}

const (
	DefaultPageNumber = 1
	DefaultPageSize   = 10
	MaxPageSize       = 100
)

const (
	CurrentPage = "currentPage"
	PageSize    = "pageSize"
)
