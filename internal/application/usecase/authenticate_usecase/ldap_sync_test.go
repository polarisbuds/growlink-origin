package authenticate_usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"growlink/internal/application/usecase/authenticate_usecase"
	"growlink/internal/domain/entity"
	ldap_entity "growlink/internal/domain/entity/ldap"
	ldapmock "growlink/internal/domain/interface/external/ldap_interface/mock"
	repomock "growlink/internal/domain/interface/repository/mock"
	"growlink/internal/domain/valueobject"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestSyncDirectory_CreatesMissingUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	ldap := ldapmock.NewMockLDAPClient(ctrl)
	userRepo := repomock.NewMockUserRepository(ctrl)
	deptRepo := repomock.NewMockDepartmentRepository(ctrl)
	deptMember := repomock.NewMockDepartmentMemberRepository(ctrl)

	ldap.EXPECT().ListProperActive().Return([]ldap_entity.LDAPUser{{
		Mail:       "new@cap-net.co.jp",
		Name:       "新規ユーザー",
		Department: "新規部署",
	}}, nil)
	userRepo.EXPECT().GetAll().Return([]entity.User{}, nil)

	dept := &entity.Department{Name: "新規部署"}
	gomock.InOrder(
		deptRepo.EXPECT().GetByName("新規部署").Return(nil, gorm.ErrRecordNotFound),
		deptRepo.EXPECT().Create(gomock.AssignableToTypeOf(dept)).DoAndReturn(func(d *entity.Department) error {
			d.ID = 99
			return nil
		}),
		deptRepo.EXPECT().GetByName("新規部署").Return(&entity.Department{Model: gorm.Model{ID: 99}, Name: "新規部署"}, nil),
	)
	userRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(u *entity.User) error {
		assert.Equal(t, "新規ユーザー", u.Name)
		assert.Equal(t, "new@cap-net.co.jp", u.Email)
		assert.Equal(t, valueobject.AuthorityGeneral, u.Authority)
		assert.Nil(t, u.Password)
		u.ID = 7
		u.UUID = uuid.New()
		return nil
	})
	deptMember.EXPECT().ReplaceByUser(uint(7), gomock.Any()).DoAndReturn(func(_ uint, member entity.DepartmentMember) error {
		assert.Equal(t, uint(99), member.DepartmentID)
		assert.Equal(t, uint(7), member.UserID)
		assert.Equal(t, uint(1), member.RoleID)
		assert.False(t, member.IsHead)
		assert.False(t, member.IsManager)
		return nil
	})

	uc := authenticate_usecase.NewAuthUsecase(ldap, userRepo, deptRepo, deptMember)
	require.NoError(t, uc.SyncDirectory(context.Background()))
}

func TestSyncDirectory_MovesDepartmentAndKeepsRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	ldap := ldapmock.NewMockLDAPClient(ctrl)
	userRepo := repomock.NewMockUserRepository(ctrl)
	deptRepo := repomock.NewMockDepartmentRepository(ctrl)
	deptMember := repomock.NewMockDepartmentMemberRepository(ctrl)

	ldap.EXPECT().ListProperActive().Return([]ldap_entity.LDAPUser{{
		Mail:       "User@cap-net.co.jp",
		Name:       "更新後",
		Department: "営業部",
	}}, nil)
	userRepo.EXPECT().GetAll().Return([]entity.User{{
		Model: gorm.Model{ID: 4},
		Name:  "更新前",
		Email: "user@cap-net.co.jp",
		DepartmentMembers: []entity.DepartmentMember{{
			Model:        gorm.Model{ID: 10},
			DepartmentID: 2,
			UserID:       4,
			RoleID:       8,
			IsHead:       true,
		}},
	}}, nil)
	userRepo.EXPECT().UpdateColumns(uint(4), map[string]any{"name": "更新後"}).Return(nil)
	deptRepo.EXPECT().GetByName("営業部").Return(&entity.Department{Model: gorm.Model{ID: 3}, Name: "営業部"}, nil)
	deptMember.EXPECT().ReplaceByUser(uint(4), gomock.Any()).DoAndReturn(func(_ uint, member entity.DepartmentMember) error {
		assert.Equal(t, uint(3), member.DepartmentID)
		assert.Equal(t, uint(8), member.RoleID)
		assert.False(t, member.IsHead)
		assert.False(t, member.IsManager)
		return nil
	})

	uc := authenticate_usecase.NewAuthUsecase(ldap, userRepo, deptRepo, deptMember)
	require.NoError(t, uc.SyncDirectory(context.Background()))
}

func TestSyncDirectory_SkipsPasswordUsersAndSuspendsMissing(t *testing.T) {
	ctrl := gomock.NewController(t)
	ldap := ldapmock.NewMockLDAPClient(ctrl)
	userRepo := repomock.NewMockUserRepository(ctrl)
	deptRepo := repomock.NewMockDepartmentRepository(ctrl)
	deptMember := repomock.NewMockDepartmentMemberRepository(ctrl)

	password := "hashed"
	ldap.EXPECT().ListProperActive().Return([]ldap_entity.LDAPUser{{
		Mail:       "kept@cap-net.co.jp",
		Name:       "在籍",
		Department: "開発部",
	}}, nil)
	userRepo.EXPECT().GetAll().Return([]entity.User{
		{Model: gorm.Model{ID: 1}, Name: "ローカル", Email: "local@example.com", Password: &password},
		{Model: gorm.Model{ID: 2}, Name: "在籍", Email: "kept@cap-net.co.jp", DepartmentMembers: []entity.DepartmentMember{{
			Model: gorm.Model{ID: 1}, DepartmentID: 5, UserID: 2, RoleID: 1,
		}}},
		{Model: gorm.Model{ID: 3}, Name: "退職", Email: "gone@cap-net.co.jp"},
		{Model: gorm.Model{ID: 6}, Name: "停止済み", Email: "stopped@cap-net.co.jp", Suspended: true},
	}, nil)
	deptRepo.EXPECT().GetByName("開発部").Return(&entity.Department{Model: gorm.Model{ID: 5}, Name: "開発部"}, nil)
	userRepo.EXPECT().SetSuspendedByIDs([]uint{3}, true).Return(nil)

	uc := authenticate_usecase.NewAuthUsecase(ldap, userRepo, deptRepo, deptMember)
	require.NoError(t, uc.SyncDirectory(context.Background()))
}

func TestSyncDirectory_DoesNotUnsuspend(t *testing.T) {
	ctrl := gomock.NewController(t)
	ldap := ldapmock.NewMockLDAPClient(ctrl)
	userRepo := repomock.NewMockUserRepository(ctrl)
	deptRepo := repomock.NewMockDepartmentRepository(ctrl)
	deptMember := repomock.NewMockDepartmentMemberRepository(ctrl)

	ldap.EXPECT().ListProperActive().Return([]ldap_entity.LDAPUser{{
		Mail:       "back@cap-net.co.jp",
		Name:       "復帰候補",
		Department: "開発部",
	}}, nil)
	userRepo.EXPECT().GetAll().Return([]entity.User{{
		Model:     gorm.Model{ID: 9},
		Name:      "復帰候補",
		Email:     "back@cap-net.co.jp",
		Suspended: true,
		DepartmentMembers: []entity.DepartmentMember{{
			Model: gorm.Model{ID: 1}, DepartmentID: 5, UserID: 9, RoleID: 1,
		}},
	}}, nil)
	deptRepo.EXPECT().GetByName("開発部").Return(&entity.Department{Model: gorm.Model{ID: 5}, Name: "開発部"}, nil)

	uc := authenticate_usecase.NewAuthUsecase(ldap, userRepo, deptRepo, deptMember)
	require.NoError(t, uc.SyncDirectory(context.Background()))
}

func TestSyncDirectory_EmptyOrFailedListChangesNothing(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepo := repomock.NewMockUserRepository(ctrl)
	deptRepo := repomock.NewMockDepartmentRepository(ctrl)
	deptMember := repomock.NewMockDepartmentMemberRepository(ctrl)

	ldapEmpty := ldapmock.NewMockLDAPClient(ctrl)
	ldapEmpty.EXPECT().ListProperActive().Return([]ldap_entity.LDAPUser{}, nil)
	uc := authenticate_usecase.NewAuthUsecase(ldapEmpty, userRepo, deptRepo, deptMember)
	require.NoError(t, uc.SyncDirectory(context.Background()))

	ldapErr := ldapmock.NewMockLDAPClient(ctrl)
	ldapErr.EXPECT().ListProperActive().Return(nil, errors.New("dial failed"))
	uc = authenticate_usecase.NewAuthUsecase(ldapErr, userRepo, deptRepo, deptMember)
	require.NoError(t, uc.SyncDirectory(context.Background()))
}

func TestNextDirectorySyncAt(t *testing.T) {
	loc := time.FixedZone("JST", 9*60*60)
	before := time.Date(2026, 11, 17, 8, 59, 0, 0, loc)
	atNine := time.Date(2026, 11, 17, 9, 0, 0, 0, loc)
	after := time.Date(2026, 11, 17, 9, 0, 1, 0, loc)

	assert.Equal(t, time.Date(2026, 11, 17, 9, 0, 0, 0, loc), authenticate_usecase.NextDirectorySyncAt(before))
	assert.Equal(t, time.Date(2026, 11, 18, 9, 0, 0, 0, loc), authenticate_usecase.NextDirectorySyncAt(atNine))
	assert.Equal(t, time.Date(2026, 11, 18, 9, 0, 0, 0, loc), authenticate_usecase.NextDirectorySyncAt(after))
}
