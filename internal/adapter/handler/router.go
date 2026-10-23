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
	repository_interface "growlink/internal/domain/interface/repository"

	"github.com/labstack/echo/v5"
)

type Deps struct {
	Auth           *authenticate_usecase.AuthUsecase
	Meeting        meeting_usecase.Bundle
	Management     management_usecase.Bundle
	Pairing        *pairing_usecase.PairingListCreatorUsecase
	PairingBrk     *pairing_usecase.PairingBreakerUsecase
	Setting        *setting_usecase.IconRegisterUsecase
	SettingRm      *setting_usecase.IconRemoverUsecase
	MeetingRepo    repository_interface.MeetingRepository
	PairingRepo    repository_interface.PairingRepository
	TemplateRepo   repository_interface.TemplateRepository
	UserRepo       repository_interface.UserRepository
	TeamRepo       repository_interface.TeamRepository
	TeamMemberRepo repository_interface.TeamMemberRepository
	DepartmentRepo repository_interface.DepartmentRepository
	RoleRepo       repository_interface.RoleRepository
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
		Auth:           deps.Auth,
		Meeting:        deps.Meeting,
		Management:     deps.Management,
		Pairing:        deps.Pairing,
		PairingBrk:     deps.PairingBrk,
		Setting:        deps.Setting,
		SettingRm:      deps.SettingRm,
		MeetingRepo:    deps.MeetingRepo,
		PairingRepo:    deps.PairingRepo,
		TemplateRepo:   deps.TemplateRepo,
		UserRepo:       deps.UserRepo,
		TeamRepo:       deps.TeamRepo,
		TeamMemberRepo: deps.TeamMemberRepo,
		DepartmentRepo: deps.DepartmentRepo,
		RoleRepo:       deps.RoleRepo,
	})

	e.GET("/login", ph.LoginPage)
	e.POST("/api/v1/login", ah.Login, auditLogging(auditLogger, true))

	spa := page.ServeSPA
	auth := e.Group("", JWTConfig(auditLogger), ActiveUserCheck(deps.UserRepo, auditLogger))
	auth.GET("/", spa, auditLogging(auditLogger, false))
	auth.GET("/meeting", spa, auditLogging(auditLogger, false))
	auth.GET("/meeting/new", spa, PairingPlanningPage(), auditLogging(auditLogger, false))
	auth.GET("/meeting/:uuid", spa, auditLogging(auditLogger, false))
	auth.GET("/meeting/:uuid/answer", spa, auditLogging(auditLogger, false))
	auth.GET("/meeting/:uuid/setting", spa, auditLogging(auditLogger, false))

	auth.GET("/planning/meetings", spa, PairingPlanningPage(), auditLogging(auditLogger, false))
	auth.GET("/planning/meetings/new", spa, PairingPlanningPage(), auditLogging(auditLogger, false))
	auth.GET("/planning/pairs", spa, PairingPlanningPage(), auditLogging(auditLogger, false))
	auth.GET("/planning/pairs/new", spa, PairingPlanningPage(), auditLogging(auditLogger, false))
	auth.GET("/planning/pairs/:uuid/edit", spa, PairingPlanningPage(), auditLogging(auditLogger, false))
	auth.GET("/planning/templates", spa, auditLogging(auditLogger, false))
	auth.GET("/planning/templates/new", spa, auditLogging(auditLogger, false))
	auth.GET("/planning/templates/:uuid", spa, auditLogging(auditLogger, false))

	auth.GET("/meeting/all", redirectPage("/planning/meetings"), auditLogging(auditLogger, false))
	auth.GET("/pairing", redirectPage("/planning/pairs"), auditLogging(auditLogger, false))
	auth.GET("/pairing/new", redirectPage("/planning/pairs/new"), auditLogging(auditLogger, false))
	auth.GET("/pairing/:id/edit", redirectPage("/planning/pairs"), auditLogging(auditLogger, false))
	auth.GET("/template", redirectPage("/planning/templates"), auditLogging(auditLogger, false))
	auth.GET("/template/new", redirectPage("/planning/templates/new"), auditLogging(auditLogger, false))
	auth.GET("/template/:id", redirectPage("/planning/templates"), auditLogging(auditLogger, false))

	auth.GET("/setting", spa, auditLogging(auditLogger, false))
	auth.GET("/management", spa, SystemAdminPage(), auditLogging(auditLogger, false))
	auth.GET("/progress", spa, ProgressPage(), auditLogging(auditLogger, false))
	auth.GET("/logout", ph.LogoutPage, auditLogging(auditLogger, true))

	apiAuth := e.Group("/api/v1", APIJWTConfig(auditLogger), ActiveUserCheckAPI(deps.UserRepo, auditLogger))
	apiAuth.POST("/logout", ah.Logout, auditLogging(auditLogger, true))
	apiAuth.GET("/avatar/me", ah.AvatarMe, auditLogging(auditLogger, false))
	apiAuth.GET("/meetings/involved", ah.ListInvolvedMeetings, auditLogging(auditLogger, false))
	apiAuth.GET("/meetings/shared", ah.ListSharedMeetings, auditLogging(auditLogger, false))
	apiAuth.GET("/meetings/:uuid", ah.GetMeeting, auditLogging(auditLogger, false))
	apiAuth.POST("/meetings/:uuid/draft", ah.SaveDraft, auditLogging(auditLogger, false))
	apiAuth.POST("/meetings/:uuid/submit", ah.SubmitMeeting, auditLogging(auditLogger, false))
	apiAuth.POST("/meetings/:uuid/draft/cancel", ah.CancelSubmit, auditLogging(auditLogger, false))
	apiAuth.POST("/meetings/:uuid/meeting-comment", ah.PostMeetingComment, auditLogging(auditLogger, false))
	apiAuth.POST("/meetings/answer-comment", ah.PostAnswerComment, auditLogging(auditLogger, false))
	apiAuth.GET("/meetings/all", ah.ListMeetingsAll, PairingPlanningAPI(), auditLogging(auditLogger, false))
	apiAuth.POST("/meetings", ah.CreateMeeting, PairingPlanningAPI(), auditLogging(auditLogger, false))
	apiAuth.POST("/meetings/all", ah.CreateMeetingByAll, PairingPlanningAPI(), auditLogging(auditLogger, false))
	apiAuth.POST("/meetings/monthly", ah.CreateMeetingMonthly, PairingPlanningAPI(), auditLogging(auditLogger, false))

	apiAuth.GET("/pairings", ah.ListPairings, PairingPlanningAPI(), auditLogging(auditLogger, false))
	apiAuth.GET("/pairings/:uuid", ah.GetPairing, PairingPlanningAPI(), auditLogging(auditLogger, false))
	apiAuth.POST("/pairings", ah.CreatePairing, PairingPlanningAPI(), auditLogging(auditLogger, false))
	apiAuth.PUT("/pairings/:uuid", ah.UpdatePairing, PairingPlanningAPI(), auditLogging(auditLogger, false))
	apiAuth.DELETE("/pairings/:uuid", ah.DeletePairing, PairingPlanningAPI(), auditLogging(auditLogger, false))

	apiAuth.GET("/templates", ah.ListTemplates, auditLogging(auditLogger, false))
	apiAuth.GET("/templates/:uuid", ah.GetTemplate, auditLogging(auditLogger, false))
	apiAuth.POST("/templates", ah.CreateTemplate, auditLogging(auditLogger, false))
	apiAuth.PUT("/templates/:uuid", ah.UpdateTemplate, auditLogging(auditLogger, false))
	apiAuth.DELETE("/templates/:uuid", ah.DeleteTemplate, auditLogging(auditLogger, false))
	apiAuth.POST("/templates/:uuid/copy", ah.CopyTemplate, auditLogging(auditLogger, false))
	apiAuth.PUT("/templates/reorder", ah.ReorderTemplates, auditLogging(auditLogger, false))

	apiAuth.GET("/planning/directory", ah.ListPlanningDirectory, PairingPlanningAPI(), auditLogging(auditLogger, false))
	apiAuth.GET("/progress/summary", ah.ProgressSummary, ProgressAPI(), auditLogging(auditLogger, false))

	apiAuth.GET("/setting/profile", ah.GetProfile, auditLogging(auditLogger, false))
	apiAuth.PUT("/setting", ah.UpdateSetting, auditLogging(auditLogger, false))
	apiAuth.POST("/setting/icon", ah.UploadIcon, auditLogging(auditLogger, false))
	apiAuth.DELETE("/setting/icon", ah.DeleteIcon, auditLogging(auditLogger, false))

	apiAuth.GET("/management/users", ah.ListManagementUsers, SystemAdminAPI(), auditLogging(auditLogger, false))
	apiAuth.POST("/management/users", ah.CreateManagementUser, SystemAdminAPI(), auditLogging(auditLogger, false))
	apiAuth.PUT("/management/users/:id", ah.UpdateManagementUser, SystemAdminAPI(), auditLogging(auditLogger, false))
	apiAuth.POST("/management/users/suspend", ah.SuspendManagementUsers, SystemAdminAPI(), auditLogging(auditLogger, false))
	apiAuth.DELETE("/management/users/:id/icon", ah.DeleteManagementUserIcon, SystemAdminAPI(), auditLogging(auditLogger, false))
	apiAuth.GET("/management/org/tree", ah.ListOrgTree, SystemAdminAPI(), auditLogging(auditLogger, false))
	apiAuth.POST("/management/departments", ah.CreateDepartment, SystemAdminAPI(), auditLogging(auditLogger, false))
	apiAuth.GET("/management/departments/:id", ah.GetDepartmentDetail, SystemAdminAPI(), auditLogging(auditLogger, false))
	apiAuth.PUT("/management/departments/:id", ah.UpdateDepartment, SystemAdminAPI(), auditLogging(auditLogger, false))
	apiAuth.DELETE("/management/departments/:id", ah.DeleteDepartment, SystemAdminAPI(), auditLogging(auditLogger, false))
	apiAuth.GET("/management/teams/:id", ah.GetTeamDetail, SystemAdminAPI(), auditLogging(auditLogger, false))
	apiAuth.PUT("/management/teams/:id", ah.UpdateTeam, SystemAdminAPI(), auditLogging(auditLogger, false))
	apiAuth.DELETE("/management/teams/:id", ah.DeleteTeam, SystemAdminAPI(), auditLogging(auditLogger, false))
	apiAuth.GET("/management/settings/remind", ah.GetRemindSettings, SystemAdminAPI(), auditLogging(auditLogger, false))
	apiAuth.PUT("/management/settings/remind", ah.UpdateRemindSettings, SystemAdminAPI(), auditLogging(auditLogger, false))
	apiAuth.GET("/management/options", ah.ListManagementOptions, SystemAdminAPI(), auditLogging(auditLogger, false))
}
