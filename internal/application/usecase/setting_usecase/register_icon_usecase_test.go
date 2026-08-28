package setting_usecase_test

import (
	"errors"
	"io"
	"strings"
	"testing"

	"growlink/internal/application/usecase/setting_usecase"
	"growlink/internal/domain/entity"
	repository_interface "growlink/internal/domain/interface/repository"
	"growlink/internal/shared/testutils/stab"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func newTestUser(id int64, imageURI *string) *entity.User {
	return &entity.User{
		Model:    gorm.Model{ID: uint(id)},
		Name:     "icon-user",
		Email:    "icon@example.com",
		ImageURI: imageURI,
	}
}

func TestRegisterIconUsecase_UploadIcon(t *testing.T) {
	const validUserID int64 = 10
	notFoundErr := errors.New("user not found")
	newFileErr := errors.New("new file failed")
	updateErr := errors.New("update failed")
	existsErr := errors.New("exists check failed")
	deleteErr := errors.New("delete failed")

	tests := []struct {
		name        string
		fileRepo    repository_interface.UploadFileRepository
		userRepo    repository_interface.UserRepository
		userID      int64
		filename    string
		content     io.Reader
		wantErr     bool
		errContains string
		assertFn    func(t *testing.T, fileRepo *stab.UploadFileRepository, userRepo *stab.UserRepository)
	}{
		{
			name: "通常系_拡張子ありのアイコン登録",
			fileRepo: &stab.UploadFileRepository{
				ExistsFunc:     func(string) (bool, error) { return true, nil },
				NewFileFunc:    func(string, io.Reader) error { return nil },
				DeleteFileFunc: func(string) error { return nil },
			},
			userRepo: &stab.UserRepository{
				GetByIDFunc: func(id int64) (*entity.User, error) {
					if id != validUserID {
						return nil, notFoundErr
					}
					return newTestUser(id, nil), nil
				},
			},
			userID:   validUserID,
			filename: "avatar.png",
			content:  strings.NewReader("png-bytes"),
			wantErr:  false,
			assertFn: func(t *testing.T, fileRepo *stab.UploadFileRepository, userRepo *stab.UserRepository) {
				t.Helper()
				require.Len(t, fileRepo.NewFileNames, 1)
				gotName := fileRepo.NewFileNames[0]
				assertUUIDFilename(t, gotName, ".png")
				require.NotNil(t, userRepo.UpdatedUser)
				require.NotNil(t, userRepo.UpdatedUser.ImageURI)
				assert.Equal(t, gotName, *userRepo.UpdatedUser.ImageURI)
				require.Len(t, fileRepo.ExistsNames, 1)
				assert.Equal(t, gotName, fileRepo.ExistsNames[0])
				require.Len(t, fileRepo.DeleteFileNames, 1)
				assert.Equal(t, gotName, fileRepo.DeleteFileNames[0])
			},
		},
		{
			name: "通常系_既存アイコンあり",
			fileRepo: &stab.UploadFileRepository{
				ExistsFunc:     func(string) (bool, error) { return true, nil },
				NewFileFunc:    func(string, io.Reader) error { return nil },
				DeleteFileFunc: func(string) error { return nil },
			},
			userRepo: func() *stab.UserRepository {
				old := "old-icon.jpg"
				return &stab.UserRepository{
					GetByIDFunc: func(id int64) (*entity.User, error) {
						return newTestUser(id, &old), nil
					},
				}
			}(),
			userID:   validUserID,
			filename: "next.webp",
			content:  strings.NewReader("webp-bytes"),
			wantErr:  false,
			assertFn: func(t *testing.T, fileRepo *stab.UploadFileRepository, userRepo *stab.UserRepository) {
				t.Helper()
				require.Len(t, fileRepo.NewFileNames, 1)
				assertUUIDFilename(t, fileRepo.NewFileNames[0], ".webp")
				require.NotNil(t, userRepo.UpdatedUser)
				require.NotNil(t, userRepo.UpdatedUser.ImageURI)
				assert.Equal(t, fileRepo.NewFileNames[0], *userRepo.UpdatedUser.ImageURI)
			},
		},
		{
			name: "通常系_存在チェックがfalseなら削除しない",
			fileRepo: &stab.UploadFileRepository{
				ExistsFunc:  func(string) (bool, error) { return false, nil },
				NewFileFunc: func(string, io.Reader) error { return nil },
			},
			userRepo: &stab.UserRepository{
				GetByIDFunc: func(id int64) (*entity.User, error) {
					return newTestUser(id, nil), nil
				},
			},
			userID:   validUserID,
			filename: "avatar.jpeg",
			content:  strings.NewReader("jpeg-bytes"),
			wantErr:  false,
			assertFn: func(t *testing.T, fileRepo *stab.UploadFileRepository, _ *stab.UserRepository) {
				t.Helper()
				require.Len(t, fileRepo.ExistsNames, 1)
				assert.Empty(t, fileRepo.DeleteFileNames)
			},
		},
		{
			name: "異常系_ユーザ取得失敗",
			fileRepo: &stab.UploadFileRepository{
				NewFileFunc: func(string, io.Reader) error {
					t.Fatal("NewFile must not be called")
					return nil
				},
			},
			userRepo: &stab.UserRepository{
				GetByIDFunc: func(int64) (*entity.User, error) {
					return nil, notFoundErr
				},
			},
			userID:      validUserID,
			filename:    "avatar.png",
			content:     strings.NewReader("png-bytes"),
			wantErr:     true,
			errContains: "ユーザアイコンのアップロードで存在しないユーザIDが指定されました",
			assertFn: func(t *testing.T, fileRepo *stab.UploadFileRepository, userRepo *stab.UserRepository) {
				t.Helper()
				assert.Empty(t, fileRepo.NewFileNames)
				assert.Nil(t, userRepo.UpdatedUser)
			},
		},
		{
			name: "異常系_ファイル新規作成失敗",
			fileRepo: &stab.UploadFileRepository{
				NewFileFunc: func(string, io.Reader) error { return newFileErr },
			},
			userRepo: &stab.UserRepository{
				GetByIDFunc: func(id int64) (*entity.User, error) {
					return newTestUser(id, nil), nil
				},
				UpdateByIDFunc: func(*entity.User) error {
					t.Fatal("UpdateByID must not be called")
					return nil
				},
			},
			userID:      validUserID,
			filename:    "avatar.png",
			content:     strings.NewReader("png-bytes"),
			wantErr:     true,
			errContains: "ユーザアイコンの新規作成処理に失敗しました",
			assertFn: func(t *testing.T, fileRepo *stab.UploadFileRepository, userRepo *stab.UserRepository) {
				t.Helper()
				require.Len(t, fileRepo.NewFileNames, 1)
				assert.Nil(t, userRepo.UpdatedUser)
			},
		},
		{
			name: "異常系_ユーザ更新失敗",
			fileRepo: &stab.UploadFileRepository{
				NewFileFunc: func(string, io.Reader) error { return nil },
				ExistsFunc: func(string) (bool, error) {
					t.Fatal("Exists must not be called")
					return false, nil
				},
			},
			userRepo: &stab.UserRepository{
				GetByIDFunc: func(id int64) (*entity.User, error) {
					return newTestUser(id, nil), nil
				},
				UpdateByIDFunc: func(*entity.User) error { return updateErr },
			},
			userID:      validUserID,
			filename:    "avatar.png",
			content:     strings.NewReader("png-bytes"),
			wantErr:     true,
			errContains: "ユーザデータの更新（ユーザアイコン）に失敗しました",
		},
		{
			name: "異常系_既存ファイルの存在チェック失敗",
			fileRepo: &stab.UploadFileRepository{
				NewFileFunc: func(string, io.Reader) error { return nil },
				ExistsFunc:  func(string) (bool, error) { return false, existsErr },
			},
			userRepo: &stab.UserRepository{
				GetByIDFunc: func(id int64) (*entity.User, error) {
					return newTestUser(id, nil), nil
				},
			},
			userID:      validUserID,
			filename:    "avatar.png",
			content:     strings.NewReader("png-bytes"),
			wantErr:     true,
			errContains: "アップロード済みのイメージファイルの存在チェックに失敗しました",
			assertFn: func(t *testing.T, fileRepo *stab.UploadFileRepository, _ *stab.UserRepository) {
				t.Helper()
				assert.Empty(t, fileRepo.DeleteFileNames)
			},
		},
		{
			name: "異常系_既存ファイルの削除失敗",
			fileRepo: &stab.UploadFileRepository{
				NewFileFunc:    func(string, io.Reader) error { return nil },
				ExistsFunc:     func(string) (bool, error) { return true, nil },
				DeleteFileFunc: func(string) error { return deleteErr },
			},
			userRepo: &stab.UserRepository{
				GetByIDFunc: func(id int64) (*entity.User, error) {
					return newTestUser(id, nil), nil
				},
			},
			userID:      validUserID,
			filename:    "avatar.png",
			content:     strings.NewReader("png-bytes"),
			wantErr:     true,
			errContains: "アップロード済みのイメージファイルの削除に失敗しました",
			assertFn: func(t *testing.T, fileRepo *stab.UploadFileRepository, _ *stab.UserRepository) {
				t.Helper()
				require.Len(t, fileRepo.DeleteFileNames, 1)
			},
		},
		{
			name: "不正値_ユーザIDが0",
			fileRepo: &stab.UploadFileRepository{
				NewFileFunc: func(string, io.Reader) error {
					t.Fatal("NewFile must not be called")
					return nil
				},
			},
			userRepo: &stab.UserRepository{
				GetByIDFunc: func(id int64) (*entity.User, error) {
					if id <= 0 {
						return nil, errors.New("invalid user id")
					}
					return newTestUser(id, nil), nil
				},
			},
			userID:      0,
			filename:    "avatar.png",
			content:     strings.NewReader("png-bytes"),
			wantErr:     true,
			errContains: "ユーザアイコンのアップロードで存在しないユーザIDが指定されました",
		},
		{
			name: "不正値_ユーザIDが負数",
			fileRepo: &stab.UploadFileRepository{
				NewFileFunc: func(string, io.Reader) error {
					t.Fatal("NewFile must not be called")
					return nil
				},
			},
			userRepo: &stab.UserRepository{
				GetByIDFunc: func(id int64) (*entity.User, error) {
					if id <= 0 {
						return nil, errors.New("invalid user id")
					}
					return newTestUser(id, nil), nil
				},
			},
			userID:      -1,
			filename:    "avatar.png",
			content:     strings.NewReader("png-bytes"),
			wantErr:     true,
			errContains: "ユーザアイコンのアップロードで存在しないユーザIDが指定されました",
		},
		{
			name: "不正値_ファイル名が空",
			fileRepo: &stab.UploadFileRepository{
				ExistsFunc:     func(string) (bool, error) { return false, nil },
				NewFileFunc:    func(string, io.Reader) error { return nil },
				DeleteFileFunc: func(string) error { return nil },
			},
			userRepo: &stab.UserRepository{
				GetByIDFunc: func(id int64) (*entity.User, error) {
					return newTestUser(id, nil), nil
				},
			},
			userID:   validUserID,
			filename: "",
			content:  strings.NewReader("bytes"),
			wantErr:  false,
			assertFn: func(t *testing.T, fileRepo *stab.UploadFileRepository, _ *stab.UserRepository) {
				t.Helper()
				require.Len(t, fileRepo.NewFileNames, 1)
				gotName := fileRepo.NewFileNames[0]
				require.True(t, strings.HasSuffix(gotName, "."))
				_, err := uuid.Parse(strings.TrimSuffix(gotName, "."))
				require.NoError(t, err)
			},
		},
		{
			name: "不正値_拡張子なしのファイル名",
			fileRepo: &stab.UploadFileRepository{
				ExistsFunc:  func(string) (bool, error) { return false, nil },
				NewFileFunc: func(string, io.Reader) error { return nil },
			},
			userRepo: &stab.UserRepository{
				GetByIDFunc: func(id int64) (*entity.User, error) {
					return newTestUser(id, nil), nil
				},
			},
			userID:   validUserID,
			filename: "avatar",
			content:  strings.NewReader("bytes"),
			wantErr:  false,
			assertFn: func(t *testing.T, fileRepo *stab.UploadFileRepository, _ *stab.UserRepository) {
				t.Helper()
				require.Len(t, fileRepo.NewFileNames, 1)
				gotName := fileRepo.NewFileNames[0]
				require.True(t, strings.HasSuffix(gotName, "."))
				_, err := uuid.Parse(strings.TrimSuffix(gotName, "."))
				require.NoError(t, err)
			},
		},
		{
			name: "不正値_contentがnil",
			fileRepo: &stab.UploadFileRepository{
				NewFileFunc: func(_ string, content io.Reader) error {
					if content == nil {
						return errors.New("content is nil")
					}
					return nil
				},
			},
			userRepo: &stab.UserRepository{
				GetByIDFunc: func(id int64) (*entity.User, error) {
					return newTestUser(id, nil), nil
				},
			},
			userID:      validUserID,
			filename:    "avatar.png",
			content:     nil,
			wantErr:     true,
			errContains: "ユーザアイコンの新規作成処理に失敗しました",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := setting_usecase.NewRegisterIconUsecase(tt.fileRepo, tt.userRepo)
			gotErr := uc.UploadIcon(tt.userID, tt.filename, tt.content)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("UploadIcon() failed: %v", gotErr)
					return
				}
				if tt.errContains != "" {
					assert.ErrorContains(t, gotErr, tt.errContains)
				}
			} else if tt.wantErr {
				t.Fatal("UploadIcon() succeeded unexpectedly")
			}

			if tt.assertFn != nil {
				fileRepo, okFile := tt.fileRepo.(*stab.UploadFileRepository)
				userRepo, okUser := tt.userRepo.(*stab.UserRepository)
				require.True(t, okFile && okUser)
				tt.assertFn(t, fileRepo, userRepo)
			}
		})
	}
}

func assertUUIDFilename(t *testing.T, filename, ext string) {
	t.Helper()
	require.True(t, strings.HasSuffix(filename, ext), "filename %q should end with %q", filename, ext)
	id, err := uuid.Parse(strings.TrimSuffix(filename, ext))
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, id)
}
