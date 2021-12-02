package web

import (
	"net/http"

	"github.com/seftomsk/abf/internal/access"
	"github.com/seftomsk/abf/internal/limiter"
)

func getAuthHandler(
	a *access.IPAccess,
	l *limiter.MultiLimiter) http.HandlerFunc {
	return CheckRequest(BlackAndWhite(a, Limiter(l)))
}
