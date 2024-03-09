package scrapers

type LiveCheckResult struct {
	ViewCount int32
	Nickname  string
	Err       error
	URL       string
}

type ByViewCount []LiveCheckResult

func (a ByViewCount) Len() int           { return len(a) }
func (a ByViewCount) Less(i, j int) bool { return a[i].ViewCount > a[j].ViewCount }
func (a ByViewCount) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
