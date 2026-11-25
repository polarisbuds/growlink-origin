package main

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"sync"
	"time"

	"growlink/internal/adapter/handler"
	"growlink/internal/application/usecase/authenticate_usecase"
	"growlink/internal/application/usecase/lineworks_usecase"
	"growlink/internal/application/usecase/management_usecase"
	"growlink/internal/application/usecase/management_usecase/line_works_messages_usecase"
	meeting_mng_usecase "growlink/internal/application/usecase/management_usecase/meeting_management_usecase"
	"growlink/internal/application/usecase/management_usecase/org_admin_usecase"
	pairing_mng_usecase "growlink/internal/application/usecase/management_usecase/pairing_management_usecase"
	progress_mng_usecase "growlink/internal/application/usecase/management_usecase/progress_management_usecase"
	"growlink/internal/application/usecase/management_usecase/remind_settings_usecase"
	"growlink/internal/application/usecase/management_usecase/role_admin_usecase"
	template_mng_usecase "growlink/internal/application/usecase/management_usecase/template_management_usecase"
	"growlink/internal/application/usecase/management_usecase/user_admin_usecase"
	"growlink/internal/application/usecase/meeting_usecase"
	"growlink/internal/application/usecase/pairing_usecase"
	"growlink/internal/application/usecase/planning_usecase"
	"growlink/internal/application/usecase/setting_usecase"
	"growlink/internal/infrastructure/external/ldap"
	"growlink/internal/infrastructure/external/lineworks"
	applogger "growlink/internal/infrastructure/logger"
	infrastructure "growlink/internal/infrastructure/postgresql"
	"growlink/internal/infrastructure/repository"
	"growlink/internal/shared/config"
	"growlink/internal/shared/utils"

	"github.com/google/uuid"
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
	userRepo := repository.NewUserRepository(ctx, db, config.GetUploadDir())
	meetingRepo := repository.NewMeetingRepository(ctx, db)
	pairingRepo := repository.NewPairingRepository(ctx, db)
	answerRepo := repository.NewAnswerRepository(ctx, db)
	templateRepo := repository.NewTemplateRepository(ctx, db)
	topicRepo := repository.NewTopicRepository(ctx, db)
	teamMemberRepo := repository.NewTeamMemberRepository(ctx, db)
	answerCommentRepo := repository.NewAnswerCommentRepository(ctx, db)
	meetingCommentRepo := repository.NewMeetingCommentRepository(ctx, db)
	uploadRepo := repository.NewStorageRepository(ctx, config.GetUploadDir())
	deptRepo := repository.NewDepartementRepository(ctx, db)
	teamRepo := repository.NewTeamsRepository(ctx, db)
	deptHierarchyRepo := repository.NewDepartmentHierarchyRepository(ctx, db)
	deptMemberRepo := repository.NewDepartmentMemberRepository(ctx, db)
	roleRepo := repository.NewRoleRepository(ctx, db)
	systemSettingsRepo := repository.NewSystemSettingsRepository(ctx, db)
	lineWorksMessageSettingRepo := repository.NewLineWorksMessageSettingRepository(ctx, db)

	ldapSetting, err := ldap.NewLDAPSetting()
	var ldapClient *ldap.LDAPClient
	if err != nil {
		slog.Warn("LDAP not configured", "err", err)
	} else {
		ldapClient = ldap.NewClient(ldapSetting)
	}

	// setting, ok, err := lineworks.Load()
	// if err != nil {
	// 	slog.Error("Line WorksのAPIの設定の読込に失敗しました: %w", slog.String("err", err.Error()))
	// 	os.Exit(2)
	// }
	// if !ok {
	// 	slog.Error("Line WorksのAPIの設定取得に失敗しました")
	// 	os.Exit(2)
	// }
	// client, err := lineworks.NewClient(setting)
	// if err != nil {
	// 	slog.Error("Line Worksのクライアント作成に失敗しました: %w", slog.String("err", err.Error()))
	// 	os.Exit(2)
	// }
	client, _ := lineworks.NewClient(lineworks.Setting{})

	authUC := authenticate_usecase.NewAuthUsecase(ldapClient, userRepo, deptRepo, deptMemberRepo)
	answerCommenter := meeting_usecase.NewAnswerCommenterUsecase(answerCommentRepo)
	meetingBundle := meeting_usecase.Bundle{
		List:                   meeting_usecase.NewMeetingListUsecase(meetingRepo, pairingRepo),
		AnswerReader:           meeting_usecase.NewAnswerReaderUsecase(answerRepo),
		DraftSaver:             meeting_usecase.NewDraftSaverUsecase(meetingRepo, answerRepo),
		Submitter:              meeting_usecase.NewSubmitterUsecase(meetingRepo, answerRepo),
		SubmitCanceller:        meeting_usecase.NewSubmitCancellerUsecase(meetingRepo),
		AnswerCommenter:        answerCommenter,
		MentorMeetingCompleter: meeting_usecase.NewMentorMeetingCompleterUsecase(meetingRepo, answerRepo, answerCommenter),
		MeetingCommenter:       meeting_usecase.NewMeetingCommenterUsecase(meetingCommentRepo),
	}
	meetingBundle.API = meeting_usecase.NewMeetingAPIUsecase(
		meetingRepo,
		userRepo,
		answerRepo,
		meetingBundle.AnswerReader,
		meetingBundle.DraftSaver,
		meetingBundle.Submitter,
		meetingBundle.SubmitCanceller,
		meetingBundle.MentorMeetingCompleter,
		meetingBundle.MeetingCommenter,
		meetingBundle.AnswerCommenter,
	)

	managementBundle := management_usecase.Bundle{
		Meeting:         meeting_mng_usecase.NewMeetingViewer(meetingRepo, pairingRepo, userRepo, teamMemberRepo),
		MeetingCreator:  meeting_mng_usecase.NewMeetingCreator(meetingRepo, pairingRepo, userRepo, templateRepo),
		MeetingUpdater:  meeting_mng_usecase.NewMeetingUpdater(meetingRepo, userRepo),
		MeetingDeleter:  meeting_mng_usecase.NewMeetingDeleter(meetingRepo),
		Pairing:         pairing_mng_usecase.NewPairingViewer(userRepo, pairingRepo, teamMemberRepo, teamRepo, deptRepo),
		PairingCreator:  pairing_mng_usecase.NewPairingCreator(pairingRepo, userRepo),
		PairingUpdater:  pairing_mng_usecase.NewPairingUpdater(pairingRepo),
		PairingDeleter:  pairing_mng_usecase.NewPairingDeleter(pairingRepo),
		Template:        template_mng_usecase.NewTemplateViewer(templateRepo),
		TemplateCreator: template_mng_usecase.NewTemplateCreator(templateRepo, topicRepo),
		TemplateUpdater: template_mng_usecase.NewTemplateUpdater(templateRepo, topicRepo),
		TemplateDeleter: template_mng_usecase.NewTemplateDeleter(templateRepo, topicRepo),
		TemplateCopier:  template_mng_usecase.NewTemplateCopier(templateRepo, topicRepo),
		UserAdmin: user_admin_usecase.NewUserAdmin(
			userRepo,
			deptMemberRepo,
			teamMemberRepo,
			deptRepo,
			teamRepo,
			roleRepo,
		),
		OrgAdmin:          org_admin_usecase.NewOrgAdmin(deptRepo, teamRepo, deptHierarchyRepo, deptMemberRepo, teamMemberRepo, userRepo),
		RemindSettings:    remind_settings_usecase.NewRemindSettings(systemSettingsRepo),
		LineWorksMessages: line_works_messages_usecase.NewLineWorksMessages(lineWorksMessageSettingRepo),
		RoleAdmin:         role_admin_usecase.NewRoleAdmin(roleRepo),
	}

	planningAccess := planning_usecase.NewPlanningAccessUsecase(
		meetingRepo,
		userRepo,
		deptRepo,
		deptMemberRepo,
		deptHierarchyRepo,
		teamRepo,
		teamMemberRepo,
		managementBundle.Meeting,
		meetingBundle.List,
		managementBundle.Pairing,
		managementBundle.PairingUpdater,
		pairingRepo,
	)
	meetingBundle.API.SetPlanningAccess(planningAccess)

	managementBundle.Progress = progress_mng_usecase.NewProgressViewer(
		meetingRepo,
		planningAccess,
		pairingRepo,
		deptRepo,
		deptHierarchyRepo,
		deptMemberRepo,
		teamRepo,
		teamMemberRepo,
	)

	pairingListCreatorUC := pairing_usecase.NewPairingListCreatorUsecase(pairingRepo, userRepo)
	deps := handler.Deps{
		Auth:           authUC,
		Meeting:        meetingBundle,
		Management:     managementBundle,
		Pairing:        pairingListCreatorUC,
		PairingBrk:     pairing_usecase.NewPairingBreakerUsecase(pairingRepo, userRepo),
		Setting:        setting_usecase.NewIconRegisterUsecase(uploadRepo, userRepo),
		SettingRm:      setting_usecase.NewIconRemoverUsecase(uploadRepo, userRepo),
		Profile:        setting_usecase.NewProfileUsecase(userRepo),
		MeetingRepo:    meetingRepo,
		AnswerRepo:     answerRepo,
		PairingRepo:    pairingRepo,
		TemplateRepo:   templateRepo,
		UserRepo:       userRepo,
		TeamRepo:       teamRepo,
		TeamMemberRepo: teamMemberRepo,
		DepartmentRepo: deptRepo,
		RoleRepo:       roleRepo,
		Planning:       planningAccess,
		LineWorks:      lineworks_usecase.New(client, lineWorksMessageSettingRepo, meetingRepo, pairingRepo),
	}

	handler.RegisterRoutes(ec, deps, logs.Audit)

	if ldapClient != nil {
		go runLDAPDirectorySync(authUC)
	}
	go runPeriodicReminder(managementBundle.RemindSettings, deps.LineWorks)
	go runPeriodicMeetingCreator(managementBundle.MeetingCreator, managementBundle.Pairing, managementBundle.Template, deps.LineWorks)

	if err := ec.Start(":" + config.GetEchoPort()); err != nil {
		slog.Error("server stopped", "err", err)
		os.Exit(1)
	}
}

func runLDAPDirectorySync(uc *authenticate_usecase.AuthUsecase) {
	var mu sync.Mutex
	run := func() {
		mu.Lock()
		defer mu.Unlock()
		if err := uc.SyncDirectory(context.Background()); err != nil {
			slog.Error("ldap directory sync failed", "err", err)
		}
	}
	run()
	for {
		time.Sleep(time.Until(authenticate_usecase.NextDirectorySyncAt(time.Now())))
		run()
	}
}

func runPeriodicReminder(remindSettingUC *remind_settings_usecase.RemindSettings, notifyUC *lineworks_usecase.LineWorksUsecase) {
	var mu sync.Mutex
	run := func() {
		mu.Lock()
		defer mu.Unlock()

		ctx := context.Background()
		remindSetting, err := remindSettingUC.Get()
		if err != nil {
			slog.Error("remind settings load failed", "err", err)
			return
		}
		now := time.Now()
		if !utils.MatchesRemindInterval(now, remindSetting.NotifyStartWeek, remindSetting.NotifyIntervalDays) {
			slog.Debug("remind skipped: cycle does not match today",
				"notifyStartWeek", remindSetting.NotifyStartWeek,
				"notifyIntervalDays", remindSetting.NotifyIntervalDays,
			)
			return
		}
		if err := notifyUC.SendOverdueMeetingReminders(ctx); err != nil {
			slog.Error("remind send failed", "err", err)
		}
	}
	for {
		notifyTime := "10:00"
		if setting, err := remindSettingUC.Get(); err == nil && setting.NotifyTime != "" {
			notifyTime = setting.NotifyTime
		} else if err != nil {
			slog.Error("remind settings load failed before schedule", "err", err)
		}
		time.Sleep(time.Until(utils.NextRemindNotifyAt(time.Now(), notifyTime)))
		run()
	}
}

func runPeriodicMeetingCreator(meetingUC *meeting_mng_usecase.MeetingCreator, pairingUC *pairing_mng_usecase.PairingViewer, templateUC *template_mng_usecase.TemplateViewer, notifyUC *lineworks_usecase.LineWorksUsecase) {
	var mux sync.Mutex
	run := func() {
		mux.Lock()
		defer mux.Unlock()

		template, err := templateUC.GetDefaultTemplate()
		if err != nil {
			slog.Error("get default template failed", slog.String("err", err.Error()))
			return
		}
		pairings, err := pairingUC.ListPairingsByAllExTemporary()
		if err != nil {
			slog.Error("list pairings failed", "err", err)
		}
		pairingIDs := make([]uint, len(pairings))
		for idx, pairing := range pairings {
			pairingIDs[idx] = pairing.ID
		}

		now := time.Now()
		implementeDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
		meetings, err := meetingUC.CreateMeetingByAll(template.ID, implementeDate, 1, nil, pairingIDs)
		if err != nil {
			slog.Error("failed periodically meeting create", "err", err)
		}
		meetingUUIDs := make([]uuid.UUID, len(meetings))
		for idx, meeting := range meetings {
			id, err := uuid.Parse(meeting.UUID)
			if err != nil {
				slog.Error("cannot parse uuid", "err", err)
				return
			}
			meetingUUIDs[idx] = id
		}

		if err := notifyUC.SendCreateMeetingMsgsBySystem(context.Background(), meetingUUIDs); err != nil {
			slog.Error("failed create meeting", "err", err)
			return
		}
	}

	for {
		time.Sleep(time.Until(utils.NextMonthFirstAt(time.Now(), 10, 0)))
		run()
	}
}
