package meeting_usecase

// Bundle は Meeting 関連ユースケースを HTTP 層へまとめて渡すための構造体。
type Bundle struct {
	List             *MeetingListUsecase
	AnswerReader     *AnswerReaderUsecase
	DraftSaver       *DraftSaverUsecase
	Submitter        *SubmitterUsecase
	SubmitCanceller  *SubmitCancellerUsecase
	AnswerCommenter  *AnswerCommenterUsecase
	MeetingCommenter *MeetingCommenterUsecase
}
