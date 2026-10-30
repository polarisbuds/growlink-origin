package management_usecase

import (
	meeting_mng_usecase "growlink/internal/application/usecase/management_usecase/meeting_management_usecase"
	org_admin_usecase "growlink/internal/application/usecase/management_usecase/org_admin_usecase"
	pairing_mng_usecase "growlink/internal/application/usecase/management_usecase/pairing_management_usecase"
	progress_mng_usecase "growlink/internal/application/usecase/management_usecase/progress_management_usecase"
	remind_settings_usecase "growlink/internal/application/usecase/management_usecase/remind_settings_usecase"
	template_mng_usecase "growlink/internal/application/usecase/management_usecase/template_management_usecase"
	user_admin_usecase "growlink/internal/application/usecase/management_usecase/user_admin_usecase"
)

// Bundle は管理系ユースケースを HTTP 層へまとめて渡すための構造体。
type Bundle struct {
	Meeting         *meeting_mng_usecase.MeetingViewer
	MeetingCreator  *meeting_mng_usecase.MeetingCreator
	MeetingUpdater  *meeting_mng_usecase.MeetingUpdater
	MeetingDeleter  *meeting_mng_usecase.MeetingDeleter
	Pairing         *pairing_mng_usecase.PairingViewer
	PairingCreator  *pairing_mng_usecase.PairingCreator
	PairingUpdater  *pairing_mng_usecase.PairingUpdater
	PairingDeleter  *pairing_mng_usecase.PairingDeleter
	Template        *template_mng_usecase.TemplateViewer
	TemplateCreator *template_mng_usecase.TemplateCreator
	TemplateUpdater *template_mng_usecase.TemplateUpdater
	TemplateDeleter *template_mng_usecase.TemplateDeleter
	TemplateCopier  *template_mng_usecase.TemplateCopier
	UserAdmin       *user_admin_usecase.UserAdmin
	OrgAdmin        *org_admin_usecase.OrgAdmin
	RemindSettings  *remind_settings_usecase.RemindSettings
	Progress        *progress_mng_usecase.ProgressViewer
}
