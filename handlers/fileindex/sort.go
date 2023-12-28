package fileindex

type sortStruct struct {
	Sort bool
	IsAsc bool
}

type SortBy string
const (
	SortByName SortBy = "name"
	SortBySize SortBy = "size"
	SortByDate SortBy = "date"
)

type sortParams struct {
	IsDesc bool
	Field SortBy
}
