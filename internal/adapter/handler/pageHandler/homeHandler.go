package page

import (
	usecase_interface "growlink/internal/adapter/interface/usecase"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
)

type Home struct {
	BaseHandler
	meetingList usecase_interface.MeetingList
	meetingUC   usecase_interface.MeetingUsecase
}

func (handler Home) Handle() echo.HandlerFunc {
	return func(c *echo.Context) error {
		var myID uint = 10
		myMeetingList, err := handler.meetingList.ListMyMeeting(myID)
		if err != nil {
			slog.Error("failed to marshal history data",
				slog.String("error", err.Error()),
			)

			return c.JSON(http.StatusInternalServerError, map[string]any{
				"error":   "failed_to_serialize_history_data",
				"message": "履歴情報の変換に失敗しました。",
			})
		}

		return nil
	}
}
