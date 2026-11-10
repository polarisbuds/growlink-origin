package api

import (
	"errors"
	"net/http"
	"strconv"

	"growlink/internal/adapter/handler/authctx"
	inputmodel "growlink/internal/application/model/input_model"
	template_mng_usecase "growlink/internal/application/usecase/management_usecase/template_management_usecase"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

func (h *Handler) resolveTemplateID(param string) (uint, error) {
	if id64, err := strconv.ParseUint(param, 10, 64); err == nil {
		return uint(id64), nil
	}
	templateUUID, err := uuid.Parse(param)
	if err != nil {
		return 0, err
	}
	tpl, err := h.deps.TemplateRepo.GetByUUID(templateUUID)
	if err != nil {
		return 0, err
	}
	return tpl.ID, nil
}

func (h *Handler) ListTemplates(c *echo.Context) error {
	list, err := h.deps.Management.Template.ListTemplate()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, list)
}

func (h *Handler) GetTemplate(c *echo.Context) error {
	templateID, err := h.resolveTemplateID(c.Param("uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	tpl, err := h.deps.Management.Template.GetTemplate(templateID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	return c.JSON(http.StatusOK, tpl)
}

func (h *Handler) CreateTemplate(c *echo.Context) error {
	claims := authctx.Claims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	var req inputmodel.TemplateInput
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := h.deps.Management.TemplateCreator.CreateTemplate(req, claims.UserID); err != nil {
		if errors.Is(err, template_mng_usecase.ErrInvalidTemplateInput) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) UpdateTemplate(c *echo.Context) error {
	claims := authctx.Claims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	templateID, err := h.resolveTemplateID(c.Param("uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	if _, err := h.deps.TemplateRepo.Get(templateID); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	var req inputmodel.TemplateInput
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := h.deps.Management.TemplateUpdater.UpdateTemplate(templateID, req); err != nil {
		if errors.Is(err, template_mng_usecase.ErrInvalidTemplateInput) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		if errors.Is(err, template_mng_usecase.ErrTemplateLocked) {
			return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) CopyTemplate(c *echo.Context) error {
	claims := authctx.Claims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	templateID, err := h.resolveTemplateID(c.Param("uuid"))
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
		tpl, err := h.deps.Management.Template.GetTemplate(templateID)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
		}
		title = tpl.Title + " のコピー"
	}
	if err := h.deps.Management.TemplateCopier.CopyTemplate(templateID, title, claims.UserID); err != nil {
		if errors.Is(err, template_mng_usecase.ErrInvalidTemplateInput) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	_ = claims
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) DeleteTemplate(c *echo.Context) error {
	claims := authctx.Claims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	templateID, err := h.resolveTemplateID(c.Param("uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	t, err := h.deps.TemplateRepo.Get(templateID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	if t.UserID != claims.UserID && !claims.SystemRole.IsAdmin() {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "permission denied"})
	}

	if err := h.deps.Management.TemplateDeleter.DeleteTemplate(templateID); err != nil {
		if errors.Is(err, template_mng_usecase.ErrTemplateLocked) {
			return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) UpdateTemplateLock(c *echo.Context) error {
	claims := authctx.Claims(c)
	if claims == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	templateID, err := h.resolveTemplateID(c.Param("uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	tpl, err := h.deps.TemplateRepo.Get(templateID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	if !authctx.CanModify(claims.UserID, claims.SystemRole, tpl.UserID) {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "permission denied"})
	}
	var req struct {
		IsLock bool `json:"isLock"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := h.deps.Management.TemplateUpdater.UpdateTemplateLock(templateID, req.IsLock); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) UpdateTemplateDefault(c *echo.Context) error {
	templateID, err := h.resolveTemplateID(c.Param("uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	if _, err := h.deps.TemplateRepo.Get(templateID); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	var req struct {
		IsDefault bool `json:"isDefault"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := h.deps.Management.TemplateUpdater.UpdateTemplateDefault(templateID, req.IsDefault); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}
