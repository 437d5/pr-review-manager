package routers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/437d5/pr-review-manager/internal/application/http/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitRouter_RoutePrefixes(t *testing.T) {
	router := InitRouter(&handlers.TeamHandler{}, &handlers.UserHandler{}, &handlers.PRHandler{})

	tests := []struct {
		name          string
		method        string
		path          string
		expectedRoute string
	}{
		{name: "team add", method: http.MethodPost, path: "/team/add", expectedRoute: "/team/"},
		{name: "team get", method: http.MethodGet, path: "/team/get?team_name=backend", expectedRoute: "/team/"},
		{name: "user set active", method: http.MethodPost, path: "/users/setIsActive", expectedRoute: "/users/"},
		{name: "user get review", method: http.MethodGet, path: "/users/getReview?user_id=u1", expectedRoute: "/users/"},
		{name: "pr create", method: http.MethodPost, path: "/pullRequest/create", expectedRoute: "/pullRequest/"},
		{name: "pr merge", method: http.MethodPost, path: "/pullRequest/merge", expectedRoute: "/pullRequest/"},
		{name: "pr reassign", method: http.MethodPost, path: "/pullRequest/reassign", expectedRoute: "/pullRequest/"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			handler, pattern := router.Handler(req)
			require.NotNil(t, handler)
			assert.Equal(t, tt.expectedRoute, pattern)
		})
	}
}

