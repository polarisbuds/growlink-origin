package valueobject

import "fmt"

type LineWorksMessageKind int

const (
	LineWorksMessageKindUnsubmitted LineWorksMessageKind = iota + 1
	LineWorksMessageKindSubmit
	LineWorksMessageKindNotConducted
	LineWorksMessageKindCancelSubmitted
	LineWorksMessageKindCreateMeeting
	LineWorksMessageKindSetImplementDate
)

func (k LineWorksMessageKind) IsValid() bool {
	switch k {
	case LineWorksMessageKindUnsubmitted, LineWorksMessageKindSubmit, LineWorksMessageKindNotConducted, LineWorksMessageKindCancelSubmitted, LineWorksMessageKindCreateMeeting, LineWorksMessageKindSetImplementDate:
		return true
	default:
		return false
	}
}

func (k LineWorksMessageKind) Int() int {
	return int(k)
}

func LineWorksMessageKindFromTitle(title string) (LineWorksMessageKind, bool) {
	switch title {
	case "未提出通知":
		return LineWorksMessageKindUnsubmitted, true
	case "提出通知":
		return LineWorksMessageKindSubmit, true
	case "未実施通知":
		return LineWorksMessageKindNotConducted, true
	case "提出取消通知":
		return LineWorksMessageKindCancelSubmitted, true
	case "面談作成通知":
		return LineWorksMessageKindCreateMeeting, true
	case "面談日設定通知":
		return LineWorksMessageKindSetImplementDate, true
	default:
		return 0, false
	}
}

func (k LineWorksMessageKind) DefaultTitle() string {
	switch k {
	case LineWorksMessageKindUnsubmitted:
		return "未提出通知"
	case LineWorksMessageKindSubmit:
		return "提出通知"
	case LineWorksMessageKindNotConducted:
		return "未実施通知"
	case LineWorksMessageKindCancelSubmitted:
		return "提出取消通知"
	case LineWorksMessageKindCreateMeeting:
		return "面談作成通知"
	case LineWorksMessageKindSetImplementDate:
		return "面談日設定通知"
	default:
		return ""
	}
}

func (k LineWorksMessageKind) DefaultContent() (string, bool) {
	switch k {
	case LineWorksMessageKindUnsubmitted:
		return "[year]年[month]月のログが未提出です。", true
	case LineWorksMessageKindSubmit:
		return "[mentee]さんが[month]月のログを提出しました。", true
	case LineWorksMessageKindNotConducted:
		return "[year]年[month]月の面談が未実施です。", true
	case LineWorksMessageKindCancelSubmitted:
		return "[mentee]さんが[month]月のログの提出を取り消しました。", true
	case LineWorksMessageKindCreateMeeting:
		return "[month]月の面談が作成されました。", true
	case LineWorksMessageKindSetImplementDate:
		return "[mentor]さんが[month]月の面談実施日を設定しました。", true
	default:
		return "", false
	}
}

func ParseLineWorksMessageKind(n int) (LineWorksMessageKind, error) {
	k := LineWorksMessageKind(n)
	if !k.IsValid() {
		return 0, fmt.Errorf("invalid line works message kind: %d", n)
	}
	return k, nil
}
