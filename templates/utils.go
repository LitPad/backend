package templates

import "html/template"

// paginationRange generates a slice of page numbers based on the current page and total pages.
func paginationRange(currentPage, totalPages int) []int {
	var pages []int
	start := currentPage - 2
	if start < 1 {
		start = 1
	}
	end := currentPage + 2
	if end > totalPages {
		end = totalPages
	}
	for i := start; i <= end; i++ {
		pages = append(pages, i)
	}
	return pages
}

// Add helper functions for pagination
var TemplateFuncMap = template.FuncMap{
	"paginationRange": paginationRange,
	"add": func(a, b int) int {
		return a + b
	},
	"sub": func(a, b int) int {
		return a - b
	},
	"sequence": func(start, end int) []int {
		var seq []int
		for i := start; i <= end; i++ {
			seq = append(seq, i)
		}
		return seq
	},
}
