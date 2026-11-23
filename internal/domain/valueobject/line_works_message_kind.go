package valueobject

import "fmt"

type LineWorksMessageKind int

const (
	LineWorksMessageKindDraft LineWorksMessageKind = iota + 1
	LineWorksMessageKindSubmit
	LineWorksMessageKindNotConducted
)

func (k LineWorksMessageKind) IsValid() bool {
	switch k {
	case LineWorksMessageKindDraft, LineWorksMessageKindSubmit, LineWorksMessageKindNotConducted:
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
		return LineWorksMessageKindDraft, true
	case "提出通知":
		return LineWorksMessageKindSubmit, true
	case "未実施通知":
		return LineWorksMessageKindNotConducted, true
	default:
		return 0, false
	}
}

func (k LineWorksMessageKind) DefaultTitle() string {
	switch k {
	case LineWorksMessageKindDraft:
		return "未提出通知"
	case LineWorksMessageKindSubmit:
		return "提出通知"
	case LineWorksMessageKindNotConducted:
		return "未実施通知"
	default:
		return ""
	}
}

func (k LineWorksMessageKind) DefaultContent() (string, bool) {
	switch k {
	case LineWorksMessageKindDraft:
		return "[year]年[month]月のログが未提出です。", true
	case LineWorksMessageKindSubmit:
		return "[mentee]さんが[month]月のログを提出しました。", true
	case LineWorksMessageKindNotConducted:
		return "[year]年[month]月の面談が未実施です。", true
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