package api

import (
	"net/http"
	"strconv"

	inputmodel "growlink/internal/application/model/input_model"

	"github.com/labstack/echo/v5"
)

func (h *Handler) ListTemplates(c *echo.Context) error {
	list, err := h.deps.Management.Template.ListTemplate()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, list)
}

func (h *Handler) GetTemplate(c *echo.Context) error {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	tpl, err := h.deps.Management.Template.GetTemplate(uint(id64))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	return c.JSON(http.StatusOK, tpl)
}

func (h *Handler) CreateTemplate(c *echo.Context) error {
	var req inputmodel.TemplateInput
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := h.deps.Management.TemplateCreator.CreateTemplate(req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) UpdateTemplate(c *echo.Context) error {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	var req inputmodel.TemplateInput
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := h.deps.Management.TemplateUpdater.UpdateTemplate(uint(id64), req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) CopyTemplate(c *echo.Context) error {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	var req struct {
		Title string `json:"title"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	title := req.Title
	if title == "" {
		tpl, err := h.deps.Management.Template.GetTemplate(uint(id64))
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
		}
		title = tpl.Title + " のコピー"
	}
	if err := h.deps.Management.TemplateCopier.CopyTemplate(uint(id64), title); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) ReorderTemplates(c *echo.Context) error {
	var req struct {
		Items []struct {
			ID        uint `json:"id"`
			DisplayNo int  `json:"displayNo"`
		} `json:"items"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	for _, item := range req.Items {
		if err := h.deps.Management.TemplateUpdater.UpdateDisplayNo(item.ID, item.DisplayNo); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) DeleteTemplate(c *echo.Context) error {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	if err := h.deps.Management.TemplateDeleter.DeleteTemplate(uint(id64)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}
