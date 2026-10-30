package api

import (
	"net/http"
	"strconv"
	"strings"

	user_admin_usecase "growlink/internal/application/usecase/management_usecase/user_admin_usecase"
	org_admin_usecase "growlink/internal/application/usecase/management_usecase/org_admin_usecase"
	remind_settings_usecase "growlink/internal/application/usecase/management_usecase/remind_settings_usecase"
	"growlink/internal/adapter/handler/authctx"
	"growlink/internal/domain/entity"

	"github.com/labstack/echo/v5"
)

func (h *Handler) CreateManagementUser(c *echo.Context) error {
	var req user_admin_usecase.CreateUserInput
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	user, err := h.deps.Management.UserAdmin.CreateUser(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, user)
}

func (h *Handler) ListManagementUsers(c *echo.Context) error {
	users, err := h.deps.Management.UserAdmin.ListUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, users)
}

func (h *Handler) UpdateManagementUser(c *echo.Context) error {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	var req user_admin_usecase.UpdateUserInput
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := h.deps.Management.UserAdmin.UpdateUser(uint(id64), req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) SuspendManagementUsers(c *echo.Context) error {
	var req struct {
		IDs       []uint `json:"ids"`
		Suspended *bool  `json:"suspended"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	suspended := true
	if req.Suspended != nil {
		suspended = *req.Suspended
	}
	if err := h.deps.Management.UserAdmin.SetUsersSuspended(req.IDs, suspended); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) DeleteManagementUserIcon(c *echo.Context) error {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	if err := h.deps.SettingRm.RemoveIcon(uint(id64)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) ListOrgTree(c *echo.Context) error {
	nodes, err := h.deps.Management.OrgAdmin.ListTree()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, nodes)
}

func (h *Handler) CreateDepartment(c *echo.Context) error {
	var req org_admin_usecase.CreateDepartmentInput
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	dept, err := h.deps.Management.OrgAdmin.CreateDepartment(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, dept)
}

func (h *Handler) GetDepartmentDetail(c *echo.Context) error {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	detail, err := h.deps.Management.OrgAdmin.GetDepartmentDetail(uint(id64))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, detail)
}

func (h *Handler) GetTeamDetail(c *echo.Context) error {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	detail, err := h.deps.Management.OrgAdmin.GetTeamDetail(uint(id64))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, detail)
}

func (h *Handler) UpdateDepartment(c *echo.Context) error {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	var req org_admin_usecase.UpdateDepartmentInput
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := h.deps.Management.OrgAdmin.UpdateDepartment(uint(id64), req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) UpdateTeam(c *echo.Context) error {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	var req org_admin_usecase.UpdateTeamInput
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := h.deps.Management.OrgAdmin.UpdateTeam(uint(id64), req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) DeleteDepartment(c *echo.Context) error {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	if err := h.deps.Management.OrgAdmin.DeleteDepartment(uint(id64)); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) CreateTeam(c *echo.Context) error {
	var req org_admin_usecase.CreateTeamInput
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	team, err := h.deps.Management.OrgAdmin.CreateTeam(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, team)
}

func (h *Handler) CreateRole(c *echo.Context) error {
	var req struct {
		Name string `json:"name"`
		Rank int    `json:"rank"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "役職名を入力してください"})
	}
	role := &entity.Role{Name: name, Rank: req.Rank}
	if err := h.deps.RoleRepo.Create(role); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, map[string]interface{}{"id": role.ID, "name": role.Name, "rank": role.Rank})
}

func (h *Handler) DeleteTeam(c *echo.Context) error {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	if err := h.deps.Management.OrgAdmin.DeleteTeam(uint(id64)); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) GetRemindSettings(c *echo.Context) error {
	settings, err := h.deps.Management.RemindSettings.Get()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, settings)
}

func (h *Handler) UpdateRemindSettings(c *echo.Context) error {
	var req remind_settings_usecase.UpdateRemindInput
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := h.deps.Management.RemindSettings.Update(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) ListManagementOptions(c *echo.Context) error {
	_ = authctx.Claims(c)
	depts, err := h.deps.Management.OrgAdmin.ListDepartmentOptions()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	teams, err := h.deps.Management.OrgAdmin.ListTeamOptions()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	roles, err := h.deps.RoleRepo.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	roleOpts := make([]map[string]interface{}, 0, len(roles))
	for _, r := range roles {
		roleOpts = append(roleOpts, map[string]interface{}{"id": r.ID, "name": r.Name, "rank": r.Rank})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"departments": depts,
		"teams":       teams,
		"roles":       roleOpts,
	})
}
