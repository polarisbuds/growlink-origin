package main

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"time"

	"growlink/internal/adapter/handler"
	authenticate_usecase "growlink/internal/application/usecase/authenticate_usecase"
	management_usecase "growlink/internal/application/usecase/management_usecase"
	meeting_mng_usecase "growlink/internal/application/usecase/management_usecase/meeting_management_usecase"
	pairing_mng_usecase "growlink/internal/application/usecase/management_usecase/pairing_management_usecase"
	org_admin_usecase "growlink/internal/application/usecase/management_usecase/org_admin_usecase"
	remind_settings_usecase "growlink/internal/application/usecase/management_usecase/remind_settings_usecase"
	template_mng_usecase "growlink/internal/application/usecase/management_usecase/template_management_usecase"
	user_admin_usecase "growlink/internal/application/usecase/management_usecase/user_admin_usecase"
	"growlink/internal/application/usecase/meeting_usecase"
	"growlink/internal/application/usecase/pairing_usecase"
	"growlink/internal/application/usecase/setting_usecase"
	"growlink/internal/infrastructure/external/ldap"
	applogger "growlink/internal/infrastructure/logger"
	infrastructure "growlink/internal/infrastructure/postgresql"
	"growlink/internal/infrastructure/repository"
	"growlink/internal/shared/config"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(".env"); err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	time.Local = loc
}

func main() {
	level := slog.LevelInfo
	if config.IsDevelop() {
		level = slog.LevelDebug
	}
	logs, err := applogger.New(applogger.Options{
		Dir:          config.GetLogDir(),
		Level:        level,
		AccessStdout: true,
	})
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := logs.Close(); err != nil {
			slog.Error("failed to close log files", "err", err)
		}
	}()
	slog.SetDefault(logs.App)

	maxProc := runtime.GOMAXPROCS(runtime.NumCPU() / 2)
	slog.Info("Use processor limit", "maxProc", maxProc)

	if err := os.MkdirAll(config.GetUploadDir(), 0755); err != nil {
		slog.Error("failed to create uploads dir", "err", err)
		os.Exit(1)
	}

	ec := handler.NewEcho(logs.Access, logs.Panic)

	db, err := infrastructure.NewPostgresqlDB(config.GetPostgresqlDSN())
	if err != nil {
		slog.Error("failed to connect to postgresql", "err", err)
		os.Exit(1)
	}

	ctx := context.Background()
	userRepo := repository.NewUserRepository(ctx, db)
	meetingRepo := repository.NewMeetingRepository(ctx, db)
	pairingRepo := repository.NewPairingRepository(ctx, db)
	answerRepo := repository.NewAnswerRepository(ctx, db)
	templateRepo := repository.NewTemplateRepository(ctx, db)
	topicRepo := repository.NewTopicRepository(ctx, db)
	relationshipRepo := repository.NewUserRelationshipRepository(ctx, db)
	teamMemberRepo := repository.NewTeamMemberRepository(ctx, db)
	answerCommentRepo := repository.NewAnswerCommentRepository(ctx, db)
	meetingCommentRepo := repository.NewMeetingCommentRepository(ctx, db)
	uploadRepo := repository.NewStorageRepository(ctx, config.GetUploadDir())
	deptRepo := repository.NewDepartementRepository(ctx, db)
	teamRepo := repository.NewTeamsRepository(ctx, db)
	deptHierarchyRepo := repository.NewDepartmentHierarchyRepository(ctx, db)
	deptRoleRepo := repository.NewDepartmentRoleRepository(ctx, db)
	deptRefRepo := repository.NewUserReferenceDepartmentRepository(ctx, db)
	teamRefRepo := repository.NewUserReferenceTeamRepository(ctx, db)
	roleRepo := repository.NewRoleRepository(ctx, db)
	systemSettingsRepo := repository.NewSystemSettingsRepository(ctx, db)

	ldapSetting, err := ldap.NewLDAPSetting()
	var ldapClient *ldap.LDAPClient
	if err != nil {
		slog.Warn("LDAP not configured", "err", err)
	} else {
		ldapClient = ldap.NewClient(ldapSetting)
	}

	authUC := authenticate_usecase.NewAuthUsecase(ldapClient, userRepo)

	meetingBundle := meeting_usecase.Bundle{
		List:             meeting_usecase.NewMeetingListUsecase(meetingRepo, pairingRepo),
		AnswerReader:     meeting_usecase.NewAnswerReaderUsecase(answerRepo),
		DraftSaver:       meeting_usecase.NewDraftSaverUsecase(meetingRepo, answerRepo),
		Submitter:        meeting_usecase.NewSubmitterUsecase(meetingRepo, answerRepo),
		SubmitCanceller:  meeting_usecase.NewSubmitCancellerUsecase(meetingRepo),
		AnswerCommenter:  meeting_usecase.NewAnswerCommenterUsecase(answerCommentRepo),
		MeetingCommenter: meeting_usecase.NewMeetingCommenterUsecase(meetingCommentRepo),
	}

	managementBundle := management_usecase.Bundle{
		Meeting: meeting_mng_usecase.NewMeetingViewer(meetingRepo, pairingRepo, relationshipRepo, userRepo, teamMemberRepo),
		MeetingCreator:  meeting_mng_usecase.NewMeetingCreator(meetingRepo, pairingRepo),
		MeetingUpdater:  meeting_mng_usecase.NewMeetingUpdater(meetingRepo),
		MeetingDeleter:  meeting_mng_usecase.NewMeetingDeleter(meetingRepo),
		Pairing:         pairing_mng_usecase.NewPairingViewer(userRepo, relationshipRepo, pairingRepo, teamMemberRepo),
		PairingCreator:  pairing_mng_usecase.NewPairingCreator(pairingRepo),
		PairingUpdater:  pairing_mng_usecase.NewPairingUpdater(pairingRepo),
		PairingDeleter:  pairing_mng_usecase.NewPairingDeleter(pairingRepo),
		Template:        template_mng_usecase.NewTemplateViewer(templateRepo),
		TemplateCreator: template_mng_usecase.NewTemplateCreator(templateRepo, topicRepo),
		TemplateUpdater: template_mng_usecase.NewTemplateUpdater(templateRepo, topicRepo),
		TemplateDeleter: template_mng_usecase.NewTemplateDeleter(templateRepo, topicRepo),
		TemplateCopier:  template_mng_usecase.NewTemplateCopier(templateRepo, topicRepo),
		UserAdmin: user_admin_usecase.NewUserAdmin(userRepo, deptRefRepo, teamRefRepo, relationshipRepo),
		OrgAdmin: org_admin_usecase.NewOrgAdmin(deptRepo, teamRepo, deptHierarchyRepo, deptRoleRepo, teamMemberRepo, userRepo),
		RemindSettings: remind_settings_usecase.NewRemindSettings(systemSettingsRepo),
	}

	deps := handler.Deps{
		Auth:        authUC,
		Meeting:     meetingBundle,
		Management:  managementBundle,
		Pairing:     pairing_usecase.NewPairingListCreatorUsecase(pairingRepo, userRepo),
		PairingBrk:  pairing_usecase.NewPairingBreakerUsecase(pairingRepo, userRepo),
		Setting:     setting_usecase.NewIconRegisterUsecase(uploadRepo, userRepo),
		SettingRm:   setting_usecase.NewIconRemoverUsecase(uploadRepo, userRepo),
		MeetingRepo: meetingRepo,
		UserRepo:    userRepo,
		RoleRepo:    roleRepo,
	}

	handler.RegisterRoutes(ec, deps, logs.Audit)

	if err := ec.Start(":" + config.GetEchoPort()); err != nil {
		slog.Error("server stopped", "err", err)
		os.Exit(1)
	}
}
