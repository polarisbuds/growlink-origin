package api

import (
	"net/http"

	"growlink/internal/adapter/handler/authctx"
	outputmodel "growlink/internal/application/model/output_model"
	line_works_messages_usecase "growlink/internal/application/usecase/management_usecase/line_works_messages_usecase"
	org_admin_usecase "growlink/internal/application/usecase/management_usecase/org_admin_usecase"
	remind_settings_usecase "growlink/internal/application/usecase/management_usecase/remind_settings_usecase"
	role_admin_usecase "growlink/internal/application/usecase/management_usecase/role_admin_usecase"
	user_admin_usecase "growlink/internal/application/usecase/management_usecase/user_admin_usecase"

	"errors"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func managementUUIDResolveError(c *echo.Context, err error) error {
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid uuid"})
	}
	return nil
}

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
	q, err := parseManagementUsersListQuery(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if q.All {
		users, err := h.deps.Management.UserAdmin.ListUsers()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		ps := len(users)
		if ps == 0 || q.Size == 0 {
			ps = user_admin_usecase.ManagementUsersPageSize
		}
		return c.JSON(http.StatusOK, outputmodel.ManagementUsersPageOutput{
			Items:    users,
			Total:    len(users),
			Page:     1,
			PageSize: ps,
		})
	}

	pageSize := q.Size
	if pageSize == 0 {
		pageSize = user_admin_usecase.ManagementUsersPageSize
	}
	out, err := h.deps.Management.UserAdmin.ListUsersPaged(user_admin_usecase.ManagementUsersListFilter{
		Name:            q.Name,
		Email:           q.Email,
		DepartmentUUIDs: q.DepartmentUUIDs,
	}, q.Page, pageSize)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, out)
}

func (h *Handler) UpdateManagementUser(c *echo.Context) error {
	var req user_admin_usecase.UpdateUserInput
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := h.deps.Management.UserAdmin.UpdateUserByUUID(c.Param("uuid"), req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) SuspendManagementUsers(c *echo.Context) error {
	var req struct {
		UUIDs     []string `json:"uuids"`
		Suspended *bool    `json:"suspended"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	suspended := true
	if req.Suspended != nil {
		suspended = *req.Suspended
	}
	if err := h.deps.Management.UserAdmin.SetUsersSuspended(req.UUIDs, suspended); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) DeleteManagementUserIcon(c *echo.Context) error {
	userID, err := h.deps.Management.OrgAdmin.ResolveUserID(c.Param("uuid"))
	if respErr := managementUUIDResolveError(c, err); respErr != nil {
		return respErr
	}
	if err := h.deps.SettingRm.RemoveIcon(userID); err != nil {
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
	detail, err := h.deps.Management.OrgAdmin.GetDepartmentDetailByUUID(c.Param("uuid"))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, detail)
}

func (h *Handler) GetTeamDetail(c *echo.Context) error {
	detail, err := h.deps.Management.OrgAdmin.GetTeamDetailByUUID(c.Param("uuid"))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, detail)
}

func (h *Handler) UpdateDepartment(c *echo.Context) error {
	var req org_admin_usecase.UpdateDepartmentInput
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := h.deps.Management.OrgAdmin.UpdateDepartmentByUUID(c.Param("uuid"), req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) UpdateTeam(c *echo.Context) error {
	var req org_admin_usecase.UpdateTeamInput
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := h.deps.Management.OrgAdmin.UpdateTeamByUUID(c.Param("uuid"), req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) DeleteDepartment(c *echo.Context) error {
	if err := h.deps.Management.OrgAdmin.DeleteDepartmentByUUID(c.Param("uuid")); err != nil {
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

func (h *Handler) ListRoles(c *echo.Context) error {
	roles, err := h.deps.Management.RoleAdmin.List()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	out := make([]map[string]interface{}, len(roles))
	for idx, r := range roles {
		out[idx] = map[string]interface{}{"uuid": r.UUID, "name": r.Name}
	}
	return c.JSON(http.StatusOK, out)
}

func (h *Handler) CreateRole(c *echo.Context) error {
	var req struct {
		Name string `json:"name"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	role, err := h.deps.Management.RoleAdmin.Create(req.Name)
	if err != nil {
		if errors.Is(err, role_admin_usecase.ErrRoleNameEmpty) || errors.Is(err, role_admin_usecase.ErrRoleNameInUse) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, map[string]interface{}{"uuid": role.UUID, "name": role.Name})
}

func (h *Handler) UpdateRole(c *echo.Context) error {
	roleUUID, parseErr := uuid.Parse(c.Param("uuid"))
	if parseErr != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid uuid"})
	}
	var req struct {
		Name string `json:"name"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	role, err := h.deps.Management.RoleAdmin.Update(roleUUID, req.Name)
	if respErr := roleAdminError(c, err); respErr != nil {
		return respErr
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"uuid": role.UUID, "name": role.Name})
}

func (h *Handler) DeleteRole(c *echo.Context) error {
	roleUUID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid uuid"})
	}
	if err := h.deps.Management.RoleAdmin.DeleteByUUID(roleUUID); err != nil {
		return roleAdminError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func roleAdminError(c *echo.Context, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	if errors.Is(err, role_admin_usecase.ErrRoleNameEmpty) ||
		errors.Is(err, role_admin_usecase.ErrRoleNameInUse) ||
		errors.Is(err, role_admin_usecase.ErrRoleHasMembers) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
}

func (h *Handler) DeleteTeam(c *echo.Context) error {
	if err := h.deps.Management.OrgAdmin.DeleteTeamByUUID(c.Param("uuid")); err != nil {
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

func (h *Handler) GetLineWorksMessages(c *echo.Context) error {
	messages, err := h.deps.Management.LineWorksMessages.List()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]any{"messages": messages})
}

func (h *Handler) UpdateLineWorksMessages(c *echo.Context) error {
	claims := authctx.Claims(c)
	var req line_works_messages_usecase.UpdateLineWorksMessagesInput
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := h.deps.Management.LineWorksMessages.Update(claims.UserID, req); err != nil {
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
	roles, err := h.deps.Management.RoleAdmin.List()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	roleOpts := make([]map[string]interface{}, len(roles))
	for idx, r := range roles {
		roleOpts[idx] = map[string]interface{}{"uuid": r.UUID, "name": r.Name}
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"departments": depts,
		"teams":       teams,
		"roles":       roleOpts,
	})
}
