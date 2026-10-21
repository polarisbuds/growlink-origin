package handler

import (
	"log/slog"

	"growlink/internal/adapter/handler/api"
	"growlink/internal/adapter/handler/page"
	authenticate_usecase "growlink/internal/application/usecase/authenticate_usecase"
	management_usecase "growlink/internal/application/usecase/management_usecase"
	"growlink/internal/application/usecase/meeting_usecase"
	"growlink/internal/application/usecase/pairing_usecase"
	"growlink/internal/application/usecase/setting_usecase"
	"growlink/internal/domain/valueobject"
	repository_interface "growlink/internal/domain/interface/repository"

	"github.com/labstack/echo/v5"
)

type Deps struct {
	Auth       *authenticate_usecase.AuthUsecase
	Meeting    meeting_usecase.Bundle
	Management management_usecase.Bundle
	Pairing    *pairing_usecase.PairingListCreatorUsecase
	PairingBrk *pairing_usecase.PairingBreakerUsecase
	Setting    *setting_usecase.IconRegisterUsecase
	SettingRm  *setting_usecase.IconRemoverUsecase
	MeetingRepo repository_interface.MeetingRepository
	UserRepo    repository_interface.UserRepository
}

func RegisterRoutes(e *echo.Echo, deps Deps, auditLogger *slog.Logger) {
	e.Static("/static", "static")
	e.Static("/uploads", "uploads")

	ph := page.New(page.Deps{
		Meeting:     deps.Meeting,
		Management:  deps.Management,
		Auth:        deps.Auth,
		UserRepo:    deps.UserRepo,
		MeetingRepo: deps.MeetingRepo,
	})
	ah := api.New(api.Deps{
		Auth:        deps.Auth,
		Meeting:     deps.Meeting,
		Management:  deps.Management,
		Pairing:     deps.Pairing,
		PairingBrk:  deps.PairingBrk,
		Setting:     deps.Setting,
		SettingRm:   deps.SettingRm,
		MeetingRepo: deps.MeetingRepo,
		UserRepo:    deps.UserRepo,
	})

	e.GET("/login", ph.LoginPage)
	e.POST("/api/v1/login", ah.Login, auditLogging(auditLogger, true))

	auth := e.Group("", JWTConfig(auditLogger))
	auth.GET("/", ph.HomePage, auditLogging(auditLogger, false))
	auth.GET("/meeting", ph.MeetingListPage, auditLogging(auditLogger, false))
	auth.GET("/meeting/all", ph.MeetingAllPage, AuthorityAtLeast(valueobject.AuthorityBoss), auditLogging(auditLogger, false))
	auth.GET("/meeting/new", ph.MeetingNewPage, AuthorityAtLeast(valueobject.AuthorityBoss), auditLogging(auditLogger, false))
	auth.GET("/meeting/:uuid", ph.MeetingDetailPage, auditLogging(auditLogger, false))
	auth.GET("/meeting/:uuid/answer", ph.MeetingAnswerPage, auditLogging(auditLogger, false))
	auth.GET("/meeting/:uuid/setting", ph.MeetingSettingPage, auditLogging(auditLogger, false))
	auth.GET("/pairing", ph.PairingListPage, AuthorityAtLeast(valueobject.AuthorityBoss), auditLogging(auditLogger, false))
	auth.GET("/pairing/new", ph.PairingNewPage, AuthorityAtLeast(valueobject.AuthorityBoss), auditLogging(auditLogger, false))
	auth.GET("/pairing/:id/edit", ph.PairingEditPage, AuthorityAtLeast(valueobject.AuthorityBoss), auditLogging(auditLogger, false))
	auth.GET("/template", ph.TemplateListPage, AuthorityAtLeast(valueobject.AuthorityOwner), auditLogging(auditLogger, false))
	auth.GET("/template/new", ph.TemplateNewPage, AuthorityAtLeast(valueobject.AuthorityOwner), auditLogging(auditLogger, false))
	auth.GET("/template/:id", ph.TemplateDetailPage, AuthorityAtLeast(valueobject.AuthorityOwner), auditLogging(auditLogger, false))
	auth.GET("/setting", ph.SettingPage, auditLogging(auditLogger, false))
	auth.GET("/logout", ph.LogoutPage, auditLogging(auditLogger, true))

	apiAuth := e.Group("/api/v1", APIJWTConfig(auditLogger))
	apiAuth.POST("/logout", ah.Logout, auditLogging(auditLogger, true))
	apiAuth.GET("/meetings/:uuid", ah.GetMeeting, auditLogging(auditLogger, false))
	apiAuth.POST("/meetings/:uuid/draft", ah.SaveDraft, auditLogging(auditLogger, false))
	apiAuth.POST("/meetings/:uuid/submit", ah.SubmitMeeting, auditLogging(auditLogger, false))
	apiAuth.POST("/meetings/:uuid/draft/cancel", ah.CancelSubmit, auditLogging(auditLogger, false))
	apiAuth.POST("/meetings/:uuid/meeting-comment", ah.PostMeetingComment, auditLogging(auditLogger, false))
	apiAuth.POST("/meetings/answer-comment", ah.PostAnswerComment, auditLogging(auditLogger, false))
	apiAuth.GET("/meetings/all", ah.ListMeetingsAll, AuthorityAtLeastAPI(valueobject.AuthorityBoss), auditLogging(auditLogger, false))
	apiAuth.POST("/meetings", ah.CreateMeeting, AuthorityAtLeastAPI(valueobject.AuthorityBoss), auditLogging(auditLogger, false))
	apiAuth.POST("/meetings/all", ah.CreateMeetingByAll, AuthorityAtLeastAPI(valueobject.AuthorityBoss), auditLogging(auditLogger, false))
	apiAuth.POST("/meetings/monthly", ah.CreateMeetingMonthly, AuthorityAtLeastAPI(valueobject.AuthorityBoss), auditLogging(auditLogger, false))

	apiAuth.GET("/pairings", ah.ListPairings, AuthorityAtLeastAPI(valueobject.AuthorityBoss), auditLogging(auditLogger, false))
	apiAuth.POST("/pairings", ah.CreatePairing, AuthorityAtLeastAPI(valueobject.AuthorityBoss), auditLogging(auditLogger, false))
	apiAuth.PUT("/pairings/:id", ah.UpdatePairing, AuthorityAtLeastAPI(valueobject.AuthorityBoss), auditLogging(auditLogger, false))
	apiAuth.DELETE("/pairings/:id", ah.DeletePairing, AuthorityAtLeastAPI(valueobject.AuthorityBoss), auditLogging(auditLogger, false))

	apiAuth.GET("/templates", ah.ListTemplates, AuthorityAtLeastAPI(valueobject.AuthorityOwner), auditLogging(auditLogger, false))
	apiAuth.GET("/templates/:id", ah.GetTemplate, AuthorityAtLeastAPI(valueobject.AuthorityOwner), auditLogging(auditLogger, false))
	apiAuth.POST("/templates", ah.CreateTemplate, AuthorityAtLeastAPI(valueobject.AuthorityOwner), auditLogging(auditLogger, false))
	apiAuth.PUT("/templates/:id", ah.UpdateTemplate, AuthorityAtLeastAPI(valueobject.AuthorityOwner), auditLogging(auditLogger, false))
	apiAuth.DELETE("/templates/:id", ah.DeleteTemplate, AuthorityAtLeastAPI(valueobject.AuthorityOwner), auditLogging(auditLogger, false))
	apiAuth.POST("/templates/:id/copy", ah.CopyTemplate, AuthorityAtLeastAPI(valueobject.AuthorityOwner), auditLogging(auditLogger, false))
	apiAuth.PUT("/templates/reorder", ah.ReorderTemplates, AuthorityAtLeastAPI(valueobject.AuthorityOwner), auditLogging(auditLogger, false))

	apiAuth.GET("/setting/profile", ah.GetProfile, auditLogging(auditLogger, false))
	apiAuth.PUT("/setting", ah.UpdateSetting, auditLogging(auditLogger, false))
	apiAuth.POST("/setting/icon", ah.UploadIcon, auditLogging(auditLogger, false))
	apiAuth.DELETE("/setting/icon", ah.DeleteIcon, auditLogging(auditLogger, false))
}
