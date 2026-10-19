package page

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
)

func (h *Handler) PairingListPage(c *echo.Context) error {
	data := h.baseData(c, "pairing")
	list, err := h.deps.Management.Pairing.ListPairingsByAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	data["Pairings"] = list
	data["PairingsJSON"] = mustJSON(list)
	return c.Render(http.StatusOK, "pairing_list.html", data)
}

func (h *Handler) PairingNewPage(c *echo.Context) error {
	data := h.baseData(c, "pairing-new")
	data["Mode"] = c.QueryParam("mode")
	if data["Mode"] == "" {
		data["Mode"] = "single"
	}
	return c.Render(http.StatusOK, "pairing_new.html", data)
}

func (h *Handler) PairingEditPage(c *echo.Context) error {
	data := h.baseData(c, "pairing-edit")
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.Redirect(http.StatusFound, "/pairing")
	}
	list, err := h.deps.Management.Pairing.ListPairingsByAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	for _, p := range list {
		if p.ID == uint(id64) {
			data["Pairing"] = p
			data["PairingJSON"] = mustJSON(p)
			data["Breadcrumb"] = "ペア編集"
			return c.Render(http.StatusOK, "pairing_edit.html", data)
		}
	}
	return c.Redirect(http.StatusFound, "/pairing")
}
