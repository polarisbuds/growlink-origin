package repository_interface

//go:generate mockgen -typed -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock

type MeetingCommentRepository interface {
}
