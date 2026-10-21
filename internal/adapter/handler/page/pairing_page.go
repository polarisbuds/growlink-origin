package page

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

func (h *Handler) PairingListPage(c *echo.Context) error {
	data := h.baseData(c, "planning-pairs")
	list, err := h.deps.Management.Pairing.ListPairingsByAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	data["Pairings"] = list
	data["PairingsJSON"] = mustJSON(list)
	data["PlanningTab"] = "pairs"
	data["Breadcrumb"] = "ペア一覧"
	return c.Render(http.StatusOK, "pairing_list.html", data)
}

func (h *Handler) PairingNewPage(c *echo.Context) error {
	data := h.baseData(c, "planning-pairs-new")
	data["Mode"] = c.QueryParam("mode")
	if data["Mode"] == "" {
		data["Mode"] = "single"
	}
	data["PlanningTab"] = "pairs"
	data["Breadcrumb"] = "ペア作成"
	return c.Render(http.StatusOK, "pairing_new.html", data)
}

func (h *Handler) PairingEditPage(c *echo.Context) error {
	data := h.baseData(c, "planning-pairs-edit")
	pairUUID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.Redirect(http.StatusFound, "/planning/pairs")
	}
	pair, err := h.deps.Management.Pairing.GetPairingByUUID(pairUUID)
	if err != nil {
		return c.Redirect(http.StatusFound, "/planning/pairs")
	}
	data["Pairing"] = pair
	data["PairingJSON"] = mustJSON(pair)
	data["PlanningTab"] = "pairs"
	data["Breadcrumb"] = "ペア編集"
	return c.Render(http.StatusOK, "pairing_edit.html", data)
}
