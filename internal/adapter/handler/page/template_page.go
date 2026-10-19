package page

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
)

func (h *Handler) TemplateListPage(c *echo.Context) error {
	data := h.baseData(c, "template")
	list, err := h.deps.Management.Template.ListTemplate()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	data["Templates"] = list
	data["TemplatesJSON"] = mustJSON(list)
	return c.Render(http.StatusOK, "template_list.html", data)
}

func (h *Handler) TemplateDetailPage(c *echo.Context) error {
	data := h.baseData(c, "template-detail")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.Redirect(http.StatusFound, "/template")
	}
	tpl, err := h.deps.Management.Template.GetTemplate(uint(id))
	if err != nil {
		return c.Redirect(http.StatusFound, "/template")
	}
	data["Template"] = tpl
	data["TemplateJSON"] = mustJSON(tpl)
	return c.Render(http.StatusOK, "template_detail.html", data)
}

func (h *Handler) TemplateNewPage(c *echo.Context) error {
	data := h.baseData(c, "template-new")
	return c.Render(http.StatusOK, "template_new.html", data)
}
