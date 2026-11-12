package api

import (
	"net/http"
	"strings"

	"growlink/internal/adapter/handler/authctx"
	org_admin_usecase "growlink/internal/application/usecase/management_usecase/org_admin_usecase"
	remind_settings_usecase "growlink/internal/application/usecase/management_usecase/remind_settings_usecase"
	user_admin_usecase "growlink/internal/application/usecase/management_usecase/user_admin_usecase"
	"growlink/internal/domain/entity"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func (h *Handler) resolveUserID(param string) (uint, error) {
	userUUID, err := uuid.Parse(param)
	if err != nil {
		return 0, err
	}
	user, err := h.deps.UserRepo.GetByUUID(userUUID)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (h *Handler) resolveDepartmentID(param string) (uint, error) {
	deptUUID, err := uuid.Parse(param)
	if err != nil {
		return 0, err
	}
	dept, err := h.deps.DepartmentRepo.GetByUUID(deptUUID)
	if err != nil {
		return 0, err
	}
	return dept.ID, nil
}

func (h *Handler) resolveTeamID(param string) (uint, error) {
	teamUUID, err := uuid.Parse(param)
	if err != nil {
		return 0, err
	}
	team, err := h.deps.TeamRepo.GetByUUID(teamUUID)
	if err != nil {
		return 0, err
	}
	return team.ID, nil
}

func (h *Handler) resolveRoleID(param string) (uint, error) {
	roleUUID, err := uuid.Parse(param)
	if err != nil {
		return 0, err
	}
	role, err := h.deps.RoleRepo.GetByUUID(roleUUID)
	if err != nil {
		return 0, err
	}
	if role == nil {
		return 0, gorm.ErrRecordNotFound
	}
	return role.ID, nil
}

func (h *Handler) roleNameInUse(name string, excludeID uint) (bool, error) {
	existing, err := h.deps.RoleRepo.GetByName(name)
	if err != nil {
		return false, err
	}
	if existing == nil {
		return false, nil
	}
	return existing.ID != excludeID, nil
}

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
	users, err := h.deps.Management.UserAdmin.ListUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, users)
}

func (h *Handler) UpdateManagementUser(c *echo.Context) error {
	userID, err := h.resolveUserID(c.Param("uuid"))
	if respErr := managementUUIDResolveError(c, err); respErr != nil {
		return respErr
	}
	var req user_admin_usecase.UpdateUserInput
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := h.deps.Management.UserAdmin.UpdateUser(userID, req); err != nil {
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
	userID, err := h.resolveUserID(c.Param("uuid"))
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
	deptID, err := h.resolveDepartmentID(c.Param("uuid"))
	if respErr := managementUUIDResolveError(c, err); respErr != nil {
		return respErr
	}
	detail, err := h.deps.Management.OrgAdmin.GetDepartmentDetail(deptID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, detail)
}

func (h *Handler) GetTeamDetail(c *echo.Context) error {
	teamID, err := h.resolveTeamID(c.Param("uuid"))
	if respErr := managementUUIDResolveError(c, err); respErr != nil {
		return respErr
	}
	detail, err := h.deps.Management.OrgAdmin.GetTeamDetail(teamID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, detail)
}

func (h *Handler) UpdateDepartment(c *echo.Context) error {
	deptID, err := h.resolveDepartmentID(c.Param("uuid"))
	if respErr := managementUUIDResolveError(c, err); respErr != nil {
		return respErr
	}
	var req org_admin_usecase.UpdateDepartmentInput
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := h.deps.Management.OrgAdmin.UpdateDepartment(deptID, req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) UpdateTeam(c *echo.Context) error {
	teamID, err := h.resolveTeamID(c.Param("uuid"))
	if respErr := managementUUIDResolveError(c, err); respErr != nil {
		return respErr
	}
	var req org_admin_usecase.UpdateTeamInput
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := h.deps.Management.OrgAdmin.UpdateTeam(teamID, req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) DeleteDepartment(c *echo.Context) error {
	deptID, err := h.resolveDepartmentID(c.Param("uuid"))
	if respErr := managementUUIDResolveError(c, err); respErr != nil {
		return respErr
	}
	if err := h.deps.Management.OrgAdmin.DeleteDepartment(deptID); err != nil {
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
	roles, err := h.deps.RoleRepo.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	out := make([]map[string]interface{}, 0, len(roles))
	for _, r := range roles {
		out = append(out, map[string]interface{}{"uuid": r.UUID.String(), "name": r.Name})
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
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "役職名を入力してください"})
	}
	inUse, err := h.roleNameInUse(name, 0)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if inUse {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "同じ名前の役職が既に存在します"})
	}
	role := &entity.Role{Name: name}
	if err := h.deps.RoleRepo.Create(role); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, map[string]interface{}{"uuid": role.UUID.String(), "name": role.Name})
}

func (h *Handler) UpdateRole(c *echo.Context) error {
	roleUUID, parseErr := uuid.Parse(c.Param("uuid"))
	if parseErr != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid uuid"})
	}
	roleID, err := h.resolveRoleID(c.Param("uuid"))
	if respErr := managementUUIDResolveError(c, err); respErr != nil {
		return respErr
	}
	var req struct {
		Name string `json:"name"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "役職名を入力してください"})
	}
	inUse, err := h.roleNameInUse(name, roleID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if inUse {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "同じ名前の役職が既に存在します"})
	}
	if err := h.deps.RoleRepo.Update(&entity.Role{Model: gorm.Model{ID: roleID}, Name: name}); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	role, err := h.deps.RoleRepo.GetByUUID(roleUUID)
	if err != nil || role == nil {
		return c.JSON(http.StatusOK, map[string]interface{}{"uuid": roleUUID.String(), "name": name})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"uuid": role.UUID.String(), "name": role.Name})
}

func (h *Handler) DeleteRole(c *echo.Context) error {
	roleID, err := h.resolveRoleID(c.Param("uuid"))
	if respErr := managementUUIDResolveError(c, err); respErr != nil {
		return respErr
	}
	count, err := h.deps.RoleRepo.CountMembersByRoleID(roleID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if count > 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "この役職はユーザに割り当てられているため削除できません"})
	}
	if err := h.deps.RoleRepo.Delete(roleID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handler) DeleteTeam(c *echo.Context) error {
	teamID, err := h.resolveTeamID(c.Param("uuid"))
	if respErr := managementUUIDResolveError(c, err); respErr != nil {
		return respErr
	}
	if err := h.deps.Management.OrgAdmin.DeleteTeam(teamID); err != nil {
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
		roleOpts = append(roleOpts, map[string]interface{}{"uuid": r.UUID.String(), "name": r.Name})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"departments": depts,
		"teams":       teams,
		"roles":       roleOpts,
	})
}
