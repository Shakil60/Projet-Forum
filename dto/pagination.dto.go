package dto

// Gere le calcul de la pagination (pages, limite, decalage).

type Pagination struct {
	Page       int
	Size       int
	Total      int
	TotalPages int
	ShowAll    bool
}

// Construit une pagination en bornant la page et en calculant le nombre total de pages.
func NewPagination(page int, size int, total int) Pagination {
	showAll := size <= 0

	if page < 1 {
		page = 1
	}

	totalPages := 1
	if !showAll && size > 0 {
		totalPages = (total + size - 1) / size
		if totalPages < 1 {
			totalPages = 1
		}
		if page > totalPages {
			page = totalPages
		}
	}

	return Pagination{
		Page:       page,
		Size:       size,
		Total:      total,
		TotalPages: totalPages,
		ShowAll:    showAll,
	}
}

func (p Pagination) Offset() int {
	if p.ShowAll {
		return 0
	}
	return (p.Page - 1) * p.Size
}

func (p Pagination) Limit() int {
	if p.ShowAll {
		return 0
	}
	return p.Size
}

func (p Pagination) HasPrev() bool {
	return !p.ShowAll && p.Page > 1
}

func (p Pagination) HasNext() bool {
	return !p.ShowAll && p.Page < p.TotalPages
}

func (p Pagination) PrevPage() int {
	if p.Page > 1 {
		return p.Page - 1
	}
	return 1
}

func (p Pagination) NextPage() int {
	if p.Page < p.TotalPages {
		return p.Page + 1
	}
	return p.TotalPages
}

func (p Pagination) Pages() []int {
	pages := make([]int, 0, p.TotalPages)
	for i := 1; i <= p.TotalPages; i++ {
		pages = append(pages, i)
	}
	return pages
}
