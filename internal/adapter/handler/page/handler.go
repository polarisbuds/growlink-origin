package page

import (
	"encoding/json"
	"html/template"
	"log/slog"
	"strings"

	"growlink/internal/adapter/handler/authctx"
	authenticate_usecase "growlink/internal/application/usecase/authenticate_usecase"
	management_usecase "growlink/internal/application/usecase/management_usecase"
	"growlink/internal/application/usecase/meeting_usecase"
	repository_interface "growlink/internal/domain/interface/repository"

	"github.com/labstack/echo/v5"
)

type Deps struct {
	Meeting     meeting_usecase.Bundle
	Management  management_usecase.Bundle
	Auth        *authenticate_usecase.AuthUsecase
	UserRepo    repository_interface.UserRepository
	MeetingRepo repository_interface.MeetingRepository
}

type Handler struct {
	deps Deps
}

type NavItem struct {
	Href   string
	Label  string
	Icon   string
	Active bool
}

func New(deps Deps) *Handler {
	return &Handler{deps: deps}
}

func (h *Handler) navigation(active string, isSystemAdmin, isAnalyst, canPairing bool) []NavItem {
	planningActive := strings.HasPrefix(active, "planning")
	nav := []NavItem{
		{Href: "/", Label: "ワークスペース", Icon: "home", Active: active == "home" || active == "workspace"},
		{Href: "/planning/meetings", Label: "プランニング", Icon: "calendar-days", Active: planningActive},
		{Href: "/setting", Label: "設定", Icon: "settings", Active: active == "setting"},
	}
	if isSystemAdmin {
		nav = append(nav, NavItem{Href: "/management", Label: "管理", Icon: "shield", Active: active == "management"})
	}
	if isSystemAdmin || isAnalyst {
		nav = append(nav, NavItem{Href: "/progress", Label: "進捗", Icon: "bar-chart-2", Active: active == "progress"})
	}
	_ = canPairing
	return nav
}

func (h *Handler) baseData(c *echo.Context, active string) map[string]interface{} {
	claims := authctx.Claims(c)
	data := map[string]interface{}{
		"Active":         active,
		"UserName":       "",
		"UserEmail":      "",
		"IsSystemAdmin":  false,
		"IsAnalyst":      false,
		"CanPairing":     false,
		"Navigation":     []NavItem{},
	}
	if claims != nil {
		isSystemAdmin := authctx.IsSystemAdmin(c)
		isAnalyst := authctx.IsAnalyst(c)
		canPairing := authctx.CanPairing(c)
		data["UserName"] = claims.Name
		data["UserEmail"] = claims.Email
		data["IsSystemAdmin"] = isSystemAdmin
		data["IsAnalyst"] = isAnalyst
		data["CanPairing"] = canPairing
		data["Navigation"] = h.navigation(active, isSystemAdmin, isAnalyst, canPairing)
	}
	return data
}

func mustJSON(v interface{}) template.JS {
	b, err := json.Marshal(v)
	if err != nil {
		slog.Error("json marshal failed", slog.String("err", err.Error()))
		return template.JS("[]")
	}
	return template.JS(b)
}
