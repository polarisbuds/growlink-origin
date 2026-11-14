package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/require"
)

func newEchoGETPairings(target string) *echo.Context {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, target, nil)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec)
}

func TestParsePairingsListPage_Missing(t *testing.T) {
	c := newEchoGETPairings("/api/v1/pairings")
	_, hasPage, err := parsePairingsListPage(c)
	require.NoError(t, err)
	require.False(t, hasPage)
}

func TestParsePairingsListPage_Invalid(t *testing.T) {
	c := newEchoGETPairings("/api/v1/pairings?page=0")
	_, hasPage, err := parsePairingsListPage(c)
	require.Error(t, err)
	require.True(t, hasPage)
}

func TestParsePlanningPairingsListFilter_InvalidUUID(t *testing.T) {
	c := newEchoGETPairings("/api/v1/pairings?teamUuids=bad")
	_, err := parsePlanningPairingsListFilter(c)
	require.Error(t, err)
}

func TestPairingsListQueryHasFilters(t *testing.T) {
	c := newEchoGETPairings("/api/v1/pairings?departmentUuids=")
	require.False(t, pairingsListQueryHasFilters(c))
	c2 := newEchoGETPairings("/api/v1/pairings?mentorUuids=550e8400-e29b-41d4-a716-446655440000")
	require.True(t, pairingsListQueryHasFilters(c2))
}