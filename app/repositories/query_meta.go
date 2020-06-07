package repositories

type queryMeta struct {
	CurrentPage int `json:"current_page"`
	Total       int `json:"total"`
	From        int `json:"from"`
	To          int `json:"to"`
	PerPage     int `json:"per_page"`
	LastPage    int `json:"last_page"`
}

func buildQueryMeta(limit, offset, total int) *queryMeta {
	currentPage := (offset / limit) + 1

	if offset >= total {
		currentPage = -1
	}

	return &queryMeta{
		CurrentPage: currentPage,
		From:        offset,
		PerPage:     limit,
		To:          offset + limit,
		LastPage:    total / limit,
		Total:       total,
	}
}
