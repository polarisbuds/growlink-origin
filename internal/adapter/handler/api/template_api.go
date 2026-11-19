package api

import (
	"errors"
	"net/http"
	"growlink/internal/adapter/handler/authctx"
	inputmodel "growlink/internal/application/model/input_model"
	template_mng_usecase "growlink/internal/application/usecase/management_usecase/template_management_usecase"

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
	tpl, err := h.deps.Management.Template.GetTemplateByParam(c.Param("uuid"))
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
	var req inputmodel.TemplateInput
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := h.deps.Management.TemplateUpdater.UpdateTemplateByParam(c.Param("uuid"), req); err != nil {
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
	var req struct {
		Title string `json:"title"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := h.deps.Management.TemplateCopier.CopyTemplateByParam(c.Param("uuid"), req.Title, claims.UserID); err != nil {
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
	actor := template_mng_usecase.TemplateActor{
		UserID:     claims.UserID,
		SystemRole: claims.SystemRole,
	}
	if err := h.deps.Management.TemplateDeleter.DeleteTemplateByParam(c.Param("uuid"), actor); err != nil {
		if errors.Is(err, template_mng_usecase.ErrTemplatePermissionDenied) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
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
	var req struct {
		IsLock bool `json:"isLock"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	actor := template_mng_usecase.TemplateActor{UserID: claims.UserID, SystemRole: claims.SystemRole}
	if err := h.deps.Management.TemplateUpdater.UpdateTemplateLockByParam(c.Param("uuid"), actor, req.IsLock); err != nil {
		if errors.Is(err, template_mng_usecase.ErrTemplatePermissionDenied) {
			return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) UpdateTemplateDefault(c *echo.Context) error {
	var req struct {
		IsDefault bool `json:"isDefault"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := h.deps.Management.TemplateUpdater.UpdateTemplateDefaultByParam(c.Param("uuid"), req.IsDefault); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}
