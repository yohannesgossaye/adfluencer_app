package common

type ListMeta struct {
	TotalDocs     int  `json:"total_docs"`
	Limit         int  `json:"limit"`
	TotalPages    int  `json:"total_pages"`
	Page          int  `json:"page"`
	PagingCounter int  `json:"paging_counter"`
	HasPrevPage   bool `json:"has_prev_page"`
	HasNextPage   bool `json:"has_next_page"`
	NextPage      *int `json:"next_page,omitempty"`
}

func NewListMeta(page, limit, total int) ListMeta {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + limit - 1) / limit
	}

	pagingCounter := 0
	if total > 0 {
		pagingCounter = (page-1)*limit + 1
	}

	hasNext := totalPages > 0 && page < totalPages
	var nextPage *int
	if hasNext {
		n := page + 1
		nextPage = &n
	}

	return ListMeta{
		TotalDocs:     total,
		Limit:         limit,
		TotalPages:    totalPages,
		Page:          page,
		PagingCounter: pagingCounter,
		HasPrevPage:   page > 1,
		HasNextPage:   hasNext,
		NextPage:      nextPage,
	}
}
