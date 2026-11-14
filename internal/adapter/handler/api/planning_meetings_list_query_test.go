package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/require"
)

func newEchoGET(target string) *echo.Context {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, target, nil)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec)
}

func TestParsePlanningMeetingsListQuery_Defaults(t *testing.T) {
	c := newEchoGET("/api/v1/meetings/all")
	filter, page, dirDesc, err := parsePlanningMeetingsListQuery(c)
	require.NoError(t, err)
	require.Equal(t, 1, page)
	require.True(t, dirDesc)
	require.Nil(t, filter.From)
}

func TestParsePlanningMeetingsListQuery_InvalidDir(t *testing.T) {
	c := newEchoGET("/api/v1/meetings/all?dir=up")
	_, _, _, err := parsePlanningMeetingsListQuery(c)
	require.Error(t, err)
}

func TestParsePlanningMeetingsListQuery_InvalidPage(t *testing.T) {
	c := newEchoGET("/api/v1/meetings/all?page=0")
	_, _, _, err := parsePlanningMeetingsListQuery(c)
	require.Error(t, err)
}

func TestParsePlanningMeetingsListQuery_InvalidStatus(t *testing.T) {
	c := newEchoGET("/api/v1/meetings/all?status=4")
	_, _, _, err := parsePlanningMeetingsListQuery(c)
	require.Error(t, err)
}

func TestParsePlanningMeetingsListQuery_InvalidUUID(t *testing.T) {
	c := newEchoGET("/api/v1/meetings/all?mentorUuids=not-a-uuid")
	_, _, _, err := parsePlanningMeetingsListQuery(c)
	require.Error(t, err)
}

func TestParsePlanningMeetingsListQuery_DateRangeInverted(t *testing.T) {
	c := newEchoGET("/api/v1/meetings/all?from=2025-02-01&to=2025-01-01")
	filter, _, _, err := parsePlanningMeetingsListQuery(c)
	require.NoError(t, err)
	require.True(t, filter.DateRangeEmpty)
}