package progress_mng_usecase

import (
	"fmt"
	"time"

	outputmodel "growlink/internal/application/model/output_model"
)

type monthWindow struct {
	year  int
	month int
	label string
	from  time.Time
	to    time.Time
}

func enumerateMonthWindows(rangeFrom, rangeTo time.Time) []monthWindow {
	loc := rangeFrom.Location()
	y, m, _ := rangeFrom.Date()
	endY, endM, _ := rangeTo.Date()

	var out []monthWindow
	for {
		monthStart := time.Date(y, m, 1, 0, 0, 0, 0, loc)
		monthEnd := monthStart.AddDate(0, 1, 0).Add(-time.Nanosecond)

		clipFrom := monthStart
		if rangeFrom.After(monthStart) {
			clipFrom = rangeFrom
		}
		clipTo := monthEnd
		if rangeTo.Before(monthEnd) {
			clipTo = rangeTo
		}
		if !clipFrom.After(clipTo) {
			out = append(out, monthWindow{
				year:  y,
				month: int(m),
				label: fmt.Sprintf("%d年%d月", y, m),
				from:  clipFrom,
				to:    clipTo,
			})
		}

		if y == endY && m == time.Month(endM) {
			break
		}
		m++
		if m > 12 {
			m = 1
			y++
		}
	}
	return out
}

func (uc *ProgressViewer) DoneTrendByMonth(filter ProgressFilter) ([]outputmodel.ProgressDoneMonthOutput, error) {
	if filter.From == nil || filter.To == nil {
		return []outputmodel.ProgressDoneMonthOutput{}, nil
	}
	if uc.deptMemberRepo == nil || uc.teamMemberRepo == nil || uc.pairingRepo == nil {
		return nil, fmt.Errorf("完了推移の集計に必要なリポジトリが設定されていません")
	}

	visible, err := uc.visibleUserIDs(filter)
	if err != nil {
		return nil, err
	}

	deptMembers, err := uc.deptMemberRepo.GetAll()
	if err != nil {
		return nil, err
	}
	teamMembers, err := uc.teamMemberRepo.GetAll()
	if err != nil {
		return nil, err
	}
	pairings, err := uc.pairingRepo.GetAll()
	if err != nil {
		return nil, err
	}
	meetings, err := uc.meetingRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("面談一覧の取得に失敗しました: %w", err)
	}

	affiliated := allAffiliatedUserIDs(deptMembers, teamMembers)
	totalUsers := intersectUsers(affiliated, visible)
	if len(totalUsers) == 0 && visible != nil {
		windows := enumerateMonthWindows(*filter.From, *filter.To)
		out := make([]outputmodel.ProgressDoneMonthOutput, 0, len(windows))
		for _, w := range windows {
			out = append(out, outputmodel.ProgressDoneMonthOutput{
				Year:  w.year,
				Month: w.month,
				Label: w.label,
				Done:  0,
			})
		}
		return out, nil
	}

	windows := enumerateMonthWindows(*filter.From, *filter.To)
	out := make([]outputmodel.ProgressDoneMonthOutput, 0, len(windows))
	for _, w := range windows {
		sub := filter
		sub.From = &w.from
		sub.To = &w.to
		stats := statsForUsers(totalUsers, pairings, meetings, sub)
		out = append(out, outputmodel.ProgressDoneMonthOutput{
			Year:  w.year,
			Month: w.month,
			Label: w.label,
			Done:  stats.doneCount,
		})
	}
	return out, nil
}
