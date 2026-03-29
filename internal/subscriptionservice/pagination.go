package subscriptionservice

import (
	"net/url"
	"strconv"
	"strings"
)

const (
	defaultListPage    = 1
	defaultListPerPage = 20
	maxListPerPage     = 100
)

// parseListPagination разбирает page и per_page из query. При ошибке ok == false.
func parseListPagination(q url.Values) (page, perPage int, ok bool) {
	page = defaultListPage
	perPage = defaultListPerPage
	if p := strings.TrimSpace(q.Get("page")); p != "" {
		v, err := strconv.Atoi(p)
		if err != nil || v < 1 {
			return 0, 0, false
		}
		page = v
	}
	if pp := strings.TrimSpace(q.Get("per_page")); pp != "" {
		v, err := strconv.Atoi(pp)
		if err != nil || v < 1 || v > maxListPerPage {
			return 0, 0, false
		}
		perPage = v
	}
	return page, perPage, true
}

func totalPages(total, perPage int) int {
	if total <= 0 || perPage <= 0 {
		return 0
	}
	return (total + perPage - 1) / perPage
}
