package setting_usecase_test

import (
	"errors"
	"io"
	"strings"
	"testing"

	"growlink/internal/application/usecase/setting_usecase"
	"growlink/internal/domain/entity"
	repomock "growlink/internal/domain/interface/repository/mock"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type iconRegisterCaptures struct {
	newFileNames    []string
	existsNames     []string
	deleteFileNames []string
	updatedUser     *entity.User
}

type iconRegisterSetupFn func(t *testing.T, user *repomock.MockUserRepository, file *repomock.MockUploadFileRepository, cap *iconRegisterCaptures)
type iconRegisterAssertFn func(t *testing.T, cap *iconRegisterCaptures)

func newIconRegisterUsecase(t *testing.T) (
	*setting_usecase.IconRegisterUsecase,
	*repomock.MockUserRepository,
	*repomock.MockUploadFileRepository,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	userRepo := repomock.NewMockUserRepository(ctrl)
	fileRepo := repomock.NewMockUploadFileRepository(ctrl)
	return setting_usecase.NewIconRegisterUsecase(fileRepo, userRepo), userRepo, fileRepo
}

func TestIconRegisterUsecase_RegistIcon(t *testing.T) {
	tests := []struct {
		name        string
		userID      uint
		filename    string
		content     io.Reader
		wantErr     bool
		errContains string
		setup       iconRegisterSetupFn
		assertFn    iconRegisterAssertFn
	}{
		{
			name:     "通常系_拡張子ありのアイコン登録",
			userID:   validUserID,
			filename: "avatar.png",
			content:  strings.NewReader("png-bytes"),
			wantErr:  false,
			setup: func(_ *testing.T, user *repomock.MockUserRepository, file *repomock.MockUploadFileRepository, cap *iconRegisterCaptures) {
				user.EXPECT().GetByID(validUserID).Return(newTestUser(validUserID, nil), nil)
				file.EXPECT().NewFile(gomock.Any(), gomock.Any()).DoAndReturn(func(name string, _ io.Reader) error {
					cap.newFileNames = append(cap.newFileNames, name)
					return nil
				})
				user.EXPECT().UpdateByID(gomock.Any()).DoAndReturn(func(u *entity.User) error {
					cap.updatedUser = u
					return nil
				})
			},
			assertFn: func(t *testing.T, cap *iconRegisterCaptures) {
				t.Helper()
				require.Len(t, cap.newFileNames, 1)
				gotName := cap.newFileNames[0]
				assertUUIDFilename(t, gotName, ".png")
				require.NotNil(t, cap.updatedUser)
				require.NotNil(t, cap.updatedUser.ImageURI)
				assert.Equal(t, gotName, *cap.updatedUser.ImageURI)
				assert.Empty(t, cap.deleteFileNames)
			},
		},
		{
			name:     "通常系_既存アイコンあり",
			userID:   validUserID,
			filename: "next.webp",
			content:  strings.NewReader("webp-bytes"),
			wantErr:  false,
			setup: func(_ *testing.T, user *repomock.MockUserRepository, file *repomock.MockUploadFileRepository, cap *iconRegisterCaptures) {
				old := "old-icon.jpg"
				user.EXPECT().GetByID(validUserID).Return(newTestUser(validUserID, &old), nil)
				file.EXPECT().NewFile(gomock.Any(), gomock.Any()).DoAndReturn(func(name string, _ io.Reader) error {
					cap.newFileNames = append(cap.newFileNames, name)
					return nil
				})
				user.EXPECT().UpdateByID(gomock.Any()).DoAndReturn(func(u *entity.User) error {
					cap.updatedUser = u
					return nil
				})
				file.EXPECT().Exists("old-icon.jpg").Return(true, nil)
				file.EXPECT().DeleteFile("old-icon.jpg").DoAndReturn(func(name string) error {
					cap.deleteFileNames = append(cap.deleteFileNames, name)
					return nil
				})
			},
			assertFn: func(t *testing.T, cap *iconRegisterCaptures) {
				t.Helper()
				require.Len(t, cap.newFileNames, 1)
				assertUUIDFilename(t, cap.newFileNames[0], ".webp")
				require.NotNil(t, cap.updatedUser)
				require.NotNil(t, cap.updatedUser.ImageURI)
				assert.Equal(t, cap.newFileNames[0], *cap.updatedUser.ImageURI)
				require.Equal(t, []string{"old-icon.jpg"}, cap.deleteFileNames)
			},
		},
		{
			name:     "通常系_古いアイコンが無ければ削除処理しない",
			userID:   validUserID,
			filename: "avatar.jpeg",
			content:  strings.NewReader("jpeg-bytes"),
			wantErr:  false,
			setup: func(_ *testing.T, user *repomock.MockUserRepository, file *repomock.MockUploadFileRepository, _ *iconRegisterCaptures) {
				user.EXPECT().GetByID(validUserID).Return(newTestUser(validUserID, nil), nil)
				file.EXPECT().NewFile(gomock.Any(), gomock.Any()).Return(nil)
				user.EXPECT().UpdateByID(gomock.Any()).Return(nil)
			},
			assertFn: func(t *testing.T, cap *iconRegisterCaptures) {
				t.Helper()
				assert.Empty(t, cap.deleteFileNames)
			},
		},
		{
			name:        "異常系_ユーザ取得失敗",
			userID:      validUserID,
			filename:    "avatar.png",
			content:     strings.NewReader("png-bytes"),
			wantErr:     true,
			errContains: "ユーザアイコンのアップロードで存在しないユーザIDが指定されました",
			setup: func(_ *testing.T, user *repomock.MockUserRepository, _ *repomock.MockUploadFileRepository, _ *iconRegisterCaptures) {
				user.EXPECT().GetByID(validUserID).Return(nil, errUserNotFound)
			},
			assertFn: func(t *testing.T, cap *iconRegisterCaptures) {
				t.Helper()
				assert.Empty(t, cap.newFileNames)
				assert.Nil(t, cap.updatedUser)
			},
		},
		{
			name:        "異常系_ファイル新規作成失敗",
			userID:      validUserID,
			filename:    "avatar.png",
			content:     strings.NewReader("png-bytes"),
			wantErr:     true,
			errContains: "ユーザアイコンの新規作成処理に失敗しました",
			setup: func(_ *testing.T, user *repomock.MockUserRepository, file *repomock.MockUploadFileRepository, cap *iconRegisterCaptures) {
				user.EXPECT().GetByID(validUserID).Return(newTestUser(validUserID, nil), nil)
				file.EXPECT().NewFile(gomock.Any(), gomock.Any()).DoAndReturn(func(name string, _ io.Reader) error {
					cap.newFileNames = append(cap.newFileNames, name)
					return errNewFileFailed
				})
			},
			assertFn: func(t *testing.T, cap *iconRegisterCaptures) {
				t.Helper()
				require.Len(t, cap.newFileNames, 1)
				assert.Nil(t, cap.updatedUser)
			},
		},
		{
			name:        "異常系_ユーザ更新失敗",
			userID:      validUserID,
			filename:    "avatar.png",
			content:     strings.NewReader("png-bytes"),
			wantErr:     true,
			errContains: "ユーザデータの更新（ユーザアイコン）に失敗しました",
			setup: func(_ *testing.T, user *repomock.MockUserRepository, file *repomock.MockUploadFileRepository, _ *iconRegisterCaptures) {
				user.EXPECT().GetByID(validUserID).Return(newTestUser(validUserID, nil), nil)
				file.EXPECT().NewFile(gomock.Any(), gomock.Any()).Return(nil)
				user.EXPECT().UpdateByID(gomock.Any()).Return(errUpdateFailed)
			},
		},
		{
			name:        "異常系_既存ファイルの存在チェック失敗",
			userID:      validUserID,
			filename:    "avatar.png",
			content:     strings.NewReader("png-bytes"),
			wantErr:     true,
			errContains: "アップロード済みのイメージファイルの存在チェックに失敗しました",
			setup: func(_ *testing.T, user *repomock.MockUserRepository, file *repomock.MockUploadFileRepository, _ *iconRegisterCaptures) {
				old := "old.png"
				user.EXPECT().GetByID(validUserID).Return(newTestUser(validUserID, &old), nil)
				file.EXPECT().NewFile(gomock.Any(), gomock.Any()).Return(nil)
				user.EXPECT().UpdateByID(gomock.Any()).Return(nil)
				file.EXPECT().Exists(gomock.Any()).Return(false, errExistsCheckFailed)
			},
			assertFn: func(t *testing.T, cap *iconRegisterCaptures) {
				t.Helper()
				assert.Empty(t, cap.deleteFileNames)
			},
		},
		{
			name:        "異常系_既存ファイルの削除失敗",
			userID:      validUserID,
			filename:    "avatar.png",
			content:     strings.NewReader("png-bytes"),
			wantErr:     true,
			errContains: "アップロード済みのイメージファイルの削除に失敗しました",
			setup: func(_ *testing.T, user *repomock.MockUserRepository, file *repomock.MockUploadFileRepository, cap *iconRegisterCaptures) {
				old := "old.png"
				user.EXPECT().GetByID(validUserID).Return(newTestUser(validUserID, &old), nil)
				file.EXPECT().NewFile(gomock.Any(), gomock.Any()).Return(nil)
				user.EXPECT().UpdateByID(gomock.Any()).Return(nil)
				file.EXPECT().Exists(gomock.Any()).Return(true, nil)
				file.EXPECT().DeleteFile(gomock.Any()).DoAndReturn(func(name string) error {
					cap.deleteFileNames = append(cap.deleteFileNames, name)
					return errDeleteFailed
				})
			},
			assertFn: func(t *testing.T, cap *iconRegisterCaptures) {
				t.Helper()
				require.Len(t, cap.deleteFileNames, 1)
			},
		},
		{
			name:        "不正値_ユーザIDが0",
			userID:      0,
			filename:    "avatar.png",
			content:     strings.NewReader("png-bytes"),
			wantErr:     true,
			errContains: "ユーザアイコンのアップロードで存在しないユーザIDが指定されました",
			setup: func(_ *testing.T, user *repomock.MockUserRepository, _ *repomock.MockUploadFileRepository, _ *iconRegisterCaptures) {
				user.EXPECT().GetByID(uint(0)).Return(nil, errors.New("invalid user id"))
			},
		},
		{
			name:     "不正値_ファイル名が空",
			userID:   validUserID,
			filename: "",
			content:  strings.NewReader("bytes"),
			wantErr:  false,
			setup: func(_ *testing.T, user *repomock.MockUserRepository, file *repomock.MockUploadFileRepository, cap *iconRegisterCaptures) {
				user.EXPECT().GetByID(validUserID).Return(newTestUser(validUserID, nil), nil)
				file.EXPECT().NewFile(gomock.Any(), gomock.Any()).DoAndReturn(func(name string, _ io.Reader) error {
					cap.newFileNames = append(cap.newFileNames, name)
					return nil
				})
				user.EXPECT().UpdateByID(gomock.Any()).Return(nil)
			},
			assertFn: func(t *testing.T, cap *iconRegisterCaptures) {
				t.Helper()
				require.Len(t, cap.newFileNames, 1)
				gotName := cap.newFileNames[0]
				require.True(t, strings.HasSuffix(gotName, "."))
				_, err := uuid.Parse(strings.TrimSuffix(gotName, "."))
				require.NoError(t, err)
			},
		},
		{
			name:     "不正値_拡張子なしのファイル名",
			userID:   validUserID,
			filename: "avatar",
			content:  strings.NewReader("bytes"),
			wantErr:  false,
			setup: func(_ *testing.T, user *repomock.MockUserRepository, file *repomock.MockUploadFileRepository, cap *iconRegisterCaptures) {
				user.EXPECT().GetByID(validUserID).Return(newTestUser(validUserID, nil), nil)
				file.EXPECT().NewFile(gomock.Any(), gomock.Any()).DoAndReturn(func(name string, _ io.Reader) error {
					cap.newFileNames = append(cap.newFileNames, name)
					return nil
				})
				user.EXPECT().UpdateByID(gomock.Any()).Return(nil)
			},
			assertFn: func(t *testing.T, cap *iconRegisterCaptures) {
				t.Helper()
				require.Len(t, cap.newFileNames, 1)
				gotName := cap.newFileNames[0]
				require.True(t, strings.HasSuffix(gotName, "."))
				_, err := uuid.Parse(strings.TrimSuffix(gotName, "."))
				require.NoError(t, err)
			},
		},
		{
			name:        "不正値_contentがnil",
			userID:      validUserID,
			filename:    "avatar.png",
			content:     nil,
			wantErr:     true,
			errContains: "ユーザアイコンの新規作成処理に失敗しました",
			setup: func(_ *testing.T, user *repomock.MockUserRepository, file *repomock.MockUploadFileRepository, _ *iconRegisterCaptures) {
				user.EXPECT().GetByID(validUserID).Return(newTestUser(validUserID, nil), nil)
				file.EXPECT().NewFile(gomock.Any(), gomock.Any()).DoAndReturn(func(_ string, content io.Reader) error {
					if content == nil {
						return errors.New("content is nil")
					}
					return nil
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, userRepo, fileRepo := newIconRegisterUsecase(t)
			cap := &iconRegisterCaptures{}
			if tt.setup != nil {
				tt.setup(t, userRepo, fileRepo, cap)
			}

			gotErr := uc.RegistIcon(tt.userID, tt.filename, tt.content)
			assertUsecaseErr(t, gotErr, tt.wantErr, tt.errContains, "UploadIcon")
			if t.Failed() {
				return
			}
			if tt.assertFn != nil {
				tt.assertFn(t, cap)
			}
		})
	}
}
