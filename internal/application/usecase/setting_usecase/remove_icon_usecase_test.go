package setting_usecase_test

import (
	"errors"
	"testing"

	"growlink/internal/application/usecase/setting_usecase"
	"growlink/internal/domain/entity"
	repomock "growlink/internal/domain/interface/repository/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type removeIconCaptures struct {
	existsNames     []string
	deleteFileNames []string
	updatedUser     *entity.User
}

type removeIconSetupFn func(t *testing.T, user *repomock.MockUserRepository, file *repomock.MockUploadFileRepository, cap *removeIconCaptures)
type removeIconAssertFn func(t *testing.T, cap *removeIconCaptures)

func newRemoveIconUsecase(t *testing.T) (
	*setting_usecase.RemoveIconUsecase,
	*repomock.MockUserRepository,
	*repomock.MockUploadFileRepository,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	userRepo := repomock.NewMockUserRepository(ctrl)
	fileRepo := repomock.NewMockUploadFileRepository(ctrl)
	return setting_usecase.NewRemoveIconUsecase(fileRepo, userRepo), userRepo, fileRepo
}

func TestRemoveIconUsecase_RemoveIcon(t *testing.T) {
	existingIcon := "current-icon.png"

	tests := []struct {
		name        string
		userID      int64
		wantErr     bool
		errContains string
		setup       removeIconSetupFn
		assertFn    removeIconAssertFn
	}{
		{
			name:    "通常系_既存アイコンあり_ファイルも存在する",
			userID:  validUserID,
			wantErr: false,
			setup: func(_ *testing.T, user *repomock.MockUserRepository, file *repomock.MockUploadFileRepository, cap *removeIconCaptures) {
				user.EXPECT().GetByID(validUserID).Return(newTestUser(validUserID, &existingIcon), nil)
				file.EXPECT().Exists(existingIcon).DoAndReturn(func(name string) (bool, error) {
					cap.existsNames = append(cap.existsNames, name)
					return true, nil
				})
				file.EXPECT().DeleteFile(existingIcon).DoAndReturn(func(name string) error {
					cap.deleteFileNames = append(cap.deleteFileNames, name)
					return nil
				})
				user.EXPECT().UpdateByID(gomock.Any()).DoAndReturn(func(u *entity.User) error {
					cap.updatedUser = u
					return nil
				})
			},
			assertFn: func(t *testing.T, cap *removeIconCaptures) {
				t.Helper()
				require.Equal(t, []string{existingIcon}, cap.existsNames)
				require.Equal(t, []string{existingIcon}, cap.deleteFileNames)
				require.NotNil(t, cap.updatedUser)
				assert.Nil(t, cap.updatedUser.ImageURI)
			},
		},
		{
			name:    "通常系_既存アイコンあり_ファイルが存在しない",
			userID:  validUserID,
			wantErr: false,
			setup: func(_ *testing.T, user *repomock.MockUserRepository, file *repomock.MockUploadFileRepository, cap *removeIconCaptures) {
				user.EXPECT().GetByID(validUserID).Return(newTestUser(validUserID, &existingIcon), nil)
				file.EXPECT().Exists(existingIcon).DoAndReturn(func(name string) (bool, error) {
					cap.existsNames = append(cap.existsNames, name)
					return false, nil
				})
				user.EXPECT().UpdateByID(gomock.Any()).DoAndReturn(func(u *entity.User) error {
					cap.updatedUser = u
					return nil
				})
			},
			assertFn: func(t *testing.T, cap *removeIconCaptures) {
				t.Helper()
				require.Equal(t, []string{existingIcon}, cap.existsNames)
				assert.Empty(t, cap.deleteFileNames)
				require.NotNil(t, cap.updatedUser)
				assert.Nil(t, cap.updatedUser.ImageURI)
			},
		},
		{
			name:    "通常系_アイコン未設定",
			userID:  validUserID,
			wantErr: false,
			setup: func(_ *testing.T, user *repomock.MockUserRepository, _ *repomock.MockUploadFileRepository, cap *removeIconCaptures) {
				user.EXPECT().GetByID(validUserID).Return(newTestUser(validUserID, nil), nil)
				user.EXPECT().UpdateByID(gomock.Any()).DoAndReturn(func(u *entity.User) error {
					cap.updatedUser = u
					return nil
				})
			},
			assertFn: func(t *testing.T, cap *removeIconCaptures) {
				t.Helper()
				assert.Empty(t, cap.existsNames)
				assert.Empty(t, cap.deleteFileNames)
				require.NotNil(t, cap.updatedUser)
				assert.Nil(t, cap.updatedUser.ImageURI)
			},
		},
		{
			name:        "異常系_ユーザ取得失敗",
			userID:      validUserID,
			wantErr:     true,
			errContains: "ユーザアイコンの削除処理で存在しないユーザIDが指定されました",
			setup: func(_ *testing.T, user *repomock.MockUserRepository, _ *repomock.MockUploadFileRepository, _ *removeIconCaptures) {
				user.EXPECT().GetByID(validUserID).Return(nil, errUserNotFound)
			},
			assertFn: func(t *testing.T, cap *removeIconCaptures) {
				t.Helper()
				assert.Empty(t, cap.existsNames)
				assert.Empty(t, cap.deleteFileNames)
				assert.Nil(t, cap.updatedUser)
			},
		},
		{
			name:        "異常系_存在チェック失敗",
			userID:      validUserID,
			wantErr:     true,
			errContains: "削除対象のイメージファイルの存在チェックに失敗しました",
			setup: func(_ *testing.T, user *repomock.MockUserRepository, file *repomock.MockUploadFileRepository, cap *removeIconCaptures) {
				user.EXPECT().GetByID(validUserID).Return(newTestUser(validUserID, &existingIcon), nil)
				file.EXPECT().Exists(existingIcon).DoAndReturn(func(name string) (bool, error) {
					cap.existsNames = append(cap.existsNames, name)
					return false, errExistsCheckFailed
				})
			},
			assertFn: func(t *testing.T, cap *removeIconCaptures) {
				t.Helper()
				require.Equal(t, []string{existingIcon}, cap.existsNames)
				assert.Empty(t, cap.deleteFileNames)
				assert.Nil(t, cap.updatedUser)
			},
		},
		{
			name:        "異常系_ファイル削除失敗",
			userID:      validUserID,
			wantErr:     true,
			errContains: "削除対象のイメージファイルの削除に失敗しました",
			setup: func(_ *testing.T, user *repomock.MockUserRepository, file *repomock.MockUploadFileRepository, cap *removeIconCaptures) {
				user.EXPECT().GetByID(validUserID).Return(newTestUser(validUserID, &existingIcon), nil)
				file.EXPECT().Exists(existingIcon).Return(true, nil)
				file.EXPECT().DeleteFile(existingIcon).DoAndReturn(func(name string) error {
					cap.deleteFileNames = append(cap.deleteFileNames, name)
					return errDeleteFailed
				})
			},
			assertFn: func(t *testing.T, cap *removeIconCaptures) {
				t.Helper()
				require.Equal(t, []string{existingIcon}, cap.deleteFileNames)
				assert.Nil(t, cap.updatedUser)
			},
		},
		{
			name:        "異常系_ユーザ更新失敗",
			userID:      validUserID,
			wantErr:     true,
			errContains: "ユーザデータの更新（ユーザアイコン）に失敗しました",
			setup: func(_ *testing.T, user *repomock.MockUserRepository, file *repomock.MockUploadFileRepository, cap *removeIconCaptures) {
				user.EXPECT().GetByID(validUserID).Return(newTestUser(validUserID, &existingIcon), nil)
				file.EXPECT().Exists(existingIcon).Return(true, nil)
				file.EXPECT().DeleteFile(existingIcon).DoAndReturn(func(name string) error {
					cap.deleteFileNames = append(cap.deleteFileNames, name)
					return nil
				})
				user.EXPECT().UpdateByID(gomock.Any()).Return(errUpdateFailed)
			},
			assertFn: func(t *testing.T, cap *removeIconCaptures) {
				t.Helper()
				require.Equal(t, []string{existingIcon}, cap.deleteFileNames)
			},
		},
		{
			name:        "不正値_ユーザIDが0",
			userID:      0,
			wantErr:     true,
			errContains: "ユーザアイコンの削除処理で存在しないユーザIDが指定されました",
			setup: func(_ *testing.T, user *repomock.MockUserRepository, _ *repomock.MockUploadFileRepository, _ *removeIconCaptures) {
				user.EXPECT().GetByID(int64(0)).Return(nil, errors.New("invalid user id"))
			},
			assertFn: func(t *testing.T, cap *removeIconCaptures) {
				t.Helper()
				assert.Empty(t, cap.existsNames)
				assert.Nil(t, cap.updatedUser)
			},
		},
		{
			name:        "不正値_ユーザIDが負数",
			userID:      -1,
			wantErr:     true,
			errContains: "ユーザアイコンの削除処理で存在しないユーザIDが指定されました",
			setup: func(_ *testing.T, user *repomock.MockUserRepository, _ *repomock.MockUploadFileRepository, _ *removeIconCaptures) {
				user.EXPECT().GetByID(int64(-1)).Return(nil, errors.New("invalid user id"))
			},
			assertFn: func(t *testing.T, cap *removeIconCaptures) {
				t.Helper()
				assert.Empty(t, cap.existsNames)
				assert.Nil(t, cap.updatedUser)
			},
		},
		{
			name:    "不正値_ImageURIが空文字",
			userID:  validUserID,
			wantErr: false,
			setup: func(_ *testing.T, user *repomock.MockUserRepository, file *repomock.MockUploadFileRepository, cap *removeIconCaptures) {
				empty := ""
				user.EXPECT().GetByID(validUserID).Return(newTestUser(validUserID, &empty), nil)
				file.EXPECT().Exists("").DoAndReturn(func(name string) (bool, error) {
					cap.existsNames = append(cap.existsNames, name)
					return false, nil
				})
				user.EXPECT().UpdateByID(gomock.Any()).DoAndReturn(func(u *entity.User) error {
					cap.updatedUser = u
					return nil
				})
			},
			assertFn: func(t *testing.T, cap *removeIconCaptures) {
				t.Helper()
				require.Equal(t, []string{""}, cap.existsNames)
				assert.Empty(t, cap.deleteFileNames)
				require.NotNil(t, cap.updatedUser)
				assert.Nil(t, cap.updatedUser.ImageURI)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, userRepo, fileRepo := newRemoveIconUsecase(t)
			cap := &removeIconCaptures{}
			if tt.setup != nil {
				tt.setup(t, userRepo, fileRepo, cap)
			}

			gotErr := uc.RemoveIcon(tt.userID)
			assertUsecaseErr(t, gotErr, tt.wantErr, tt.errContains, "RemoveIcon")
			if t.Failed() {
				return
			}
			if tt.assertFn != nil {
				tt.assertFn(t, cap)
			}
		})
	}
}
