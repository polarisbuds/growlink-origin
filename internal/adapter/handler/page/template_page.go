package page

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

func (h *Handler) TemplateListPage(c *echo.Context) error {
	data := h.baseData(c, "planning-templates")
	list, err := h.deps.Management.Template.ListTemplate()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	data["Templates"] = list
	data["TemplatesJSON"] = mustJSON(list)
	data["PlanningTab"] = "templates"
	data["Breadcrumb"] = "テンプレート"
	return c.Render(http.StatusOK, "template_list.html", data)
}

func (h *Handler) TemplateDetailPage(c *echo.Context) error {
	data := h.baseData(c, "planning-templates-detail")
	templateUUID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.Redirect(http.StatusFound, "/planning/templates")
	}
	tpl, err := h.deps.Management.Template.GetTemplateByUUID(templateUUID)
	if err != nil {
		return c.Redirect(http.StatusFound, "/planning/templates")
	}
	data["Template"] = tpl
	data["TemplateJSON"] = mustJSON(tpl)
	data["PlanningTab"] = "templates"
	data["Breadcrumb"] = "テンプレート詳細"
	return c.Render(http.StatusOK, "template_detail.html", data)
}

func (h *Handler) TemplateNewPage(c *echo.Context) error {
	data := h.baseData(c, "planning-templates-new")
	data["PlanningTab"] = "templates"
	data["Breadcrumb"] = "テンプレート作成"
	return c.Render(http.StatusOK, "template_new.html", data)
}
