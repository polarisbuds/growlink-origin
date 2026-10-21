package page

import (
	"encoding/json"
	"html/template"
	"log/slog"

	"growlink/internal/adapter/handler/authctx"
	authenticate_usecase "growlink/internal/application/usecase/authenticate_usecase"
	management_usecase "growlink/internal/application/usecase/management_usecase"
	"growlink/internal/application/usecase/meeting_usecase"
	repository_interface "growlink/internal/domain/interface/repository"
	"growlink/internal/domain/valueobject"

	"github.com/labstack/echo/v5"
)

type Deps struct {
	Meeting    meeting_usecase.Bundle
	Management management_usecase.Bundle
	Auth       *authenticate_usecase.AuthUsecase
	UserRepo   repository_interface.UserRepository
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

func (h *Handler) navigation(active string, isBoss, isOwner bool) []NavItem {
	nav := []NavItem{
		{Href: "/", Label: "あなたの1on1", Icon: "home", Active: active == "home"},
		{Href: "/meeting", Label: "面談一覧", Icon: "calendar-days", Active: active == "meeting" || active == "meeting-detail" || active == "meeting-answer" || active == "meeting-setting"},
	}
	if isBoss {
		nav = append(nav, NavItem{Href: "/meeting/all", Label: "全ての面談", Icon: "layers", Active: active == "meeting-all"})
		nav = append(nav, NavItem{Href: "/pairing", Label: "ペア一覧", Icon: "users", Active: active == "pairing" || active == "pairing-new" || active == "pairing-edit"})
	}
	if isOwner {
		nav = append(nav, NavItem{Href: "/template", Label: "テンプレート", Icon: "file-text", Active: active == "template" || active == "template-new" || active == "template-detail"})
		nav = append(nav, NavItem{Href: "/management", Label: "管理", Icon: "shield", Active: active == "management"})
	}
	nav = append(nav, NavItem{Href: "/setting", Label: "設定", Icon: "settings", Active: active == "setting"})
	return nav
}

func (h *Handler) baseData(c *echo.Context, active string) map[string]interface{} {
	claims := authctx.Claims(c)
	data := map[string]interface{}{
		"Active":    active,
		"UserName":  "",
		"UserEmail": "",
		"IsOwner":   false,
		"IsBoss":    false,
		"Navigation": []NavItem{},
	}
	if claims != nil {
		isOwner := claims.Authority == valueobject.AuthorityOwner
		isBoss := claims.Authority >= valueobject.AuthorityBoss
		data["UserName"] = claims.Name
		data["UserEmail"] = claims.Email
		data["IsOwner"] = isOwner
		data["IsBoss"] = isBoss
		data["Navigation"] = h.navigation(active, isBoss, isOwner)
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
