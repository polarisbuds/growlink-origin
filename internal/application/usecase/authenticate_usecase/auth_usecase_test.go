package authenticate_usecase_test

import (
	"context"
	"os"
	"testing"

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

func TestMain(m *testing.M) {
	_ = os.Setenv("JWK_KEY", "test-jwt-key-for-unit-tests")
	os.Exit(m.Run())
}

func TestLoginLDAP_NewUser_CreatesDepartmentAndAssignsMembership(t *testing.T) {
	ctrl := gomock.NewController(t)
	ldap := ldapmock.NewMockLDAPClient(ctrl)
	userRepo := repomock.NewMockUserRepository(ctrl)
	deptRepo := repomock.NewMockDepartmentRepository(ctrl)
	deptMember := repomock.NewMockDepartmentMemberRepository(ctrl)

	const deptName = "新規部署"
	ldapUser := &ldap_entity.LDAPUser{
		Mail:       "new@cap-net.co.jp",
		Name:       "新規ユーザー",
		Department: deptName,
	}
	ldap.EXPECT().Authenticate("new", "secret").Return(ldapUser, nil)

	dept := &entity.Department{Name: deptName}
	gomock.InOrder(
		deptRepo.EXPECT().GetByName(deptName).Return(nil, gorm.ErrRecordNotFound),
		deptRepo.EXPECT().Create(gomock.AssignableToTypeOf(dept)).DoAndReturn(func(d *entity.Department) error {
			d.ID = 99
			return nil
		}),
		deptRepo.EXPECT().GetByName(deptName).Return(&entity.Department{Model: gorm.Model{ID: 99}, Name: deptName}, nil),
	)

	userRepo.EXPECT().GetByEmail(ldapUser.Mail).Return(nil, gorm.ErrRecordNotFound)
	userRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(u *entity.User) error {
		u.ID = 7
		u.UUID = uuid.New()
		return nil
	})

	deptMember.EXPECT().GetByDepartmentID(uint(99)).Return(nil, nil)
	deptMember.EXPECT().ReplaceByDepartment(uint(99), gomock.Any()).DoAndReturn(func(_ uint, members []entity.DepartmentMember) error {
		require.Len(t, members, 1)
		assert.Equal(t, uint(99), members[0].DepartmentID)
		assert.Equal(t, uint(7), members[0].UserID)
		assert.Equal(t, uint(1), members[0].RoleID)
		return nil
	})

	uc := authenticate_usecase.NewAuthUsecase(ldap, userRepo, deptRepo, deptMember)
	token, _, err := uc.LoginLDAP(context.Background(), "new@cap-net.co.jp", "secret")
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestLoginLDAP_NewUser_ExistingDepartmentName(t *testing.T) {
	ctrl := gomock.NewController(t)
	ldap := ldapmock.NewMockLDAPClient(ctrl)
	userRepo := repomock.NewMockUserRepository(ctrl)
	deptRepo := repomock.NewMockDepartmentRepository(ctrl)
	deptMember := repomock.NewMockDepartmentMemberRepository(ctrl)

	const deptName = "開発部"
	existingDept := &entity.Department{Model: gorm.Model{ID: 2}, Name: deptName}
	ldapUser := &ldap_entity.LDAPUser{
		Mail:       "dev@cap-net.co.jp",
		Name:       "開発メンバー",
		Department: deptName,
	}
	ldap.EXPECT().Authenticate("dev", "secret").Return(ldapUser, nil)

	deptRepo.EXPECT().GetByName(deptName).Return(existingDept, nil).Times(2)
	deptRepo.EXPECT().Create(gomock.Any()).Times(0)

	userRepo.EXPECT().GetByEmail(ldapUser.Mail).Return(nil, gorm.ErrRecordNotFound)
	userRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(u *entity.User) error {
		u.ID = 8
		u.UUID = uuid.New()
		return nil
	})

	deptMember.EXPECT().GetByDepartmentID(uint(2)).Return(nil, nil)
	deptMember.EXPECT().ReplaceByDepartment(uint(2), gomock.Any()).DoAndReturn(func(_ uint, members []entity.DepartmentMember) error {
		require.Len(t, members, 1)
		assert.Equal(t, uint(2), members[0].DepartmentID)
		return nil
	})

	uc := authenticate_usecase.NewAuthUsecase(ldap, userRepo, deptRepo, deptMember)
	_, _, err := uc.LoginLDAP(context.Background(), "dev@cap-net.co.jp", "secret")
	require.NoError(t, err)
}

func TestLoginLDAP_NewUser_EmptyDepartmentUsesDefault(t *testing.T) {
	ctrl := gomock.NewController(t)
	ldap := ldapmock.NewMockLDAPClient(ctrl)
	userRepo := repomock.NewMockUserRepository(ctrl)
	deptRepo := repomock.NewMockDepartmentRepository(ctrl)
	deptMember := repomock.NewMockDepartmentMemberRepository(ctrl)

	ldapUser := &ldap_entity.LDAPUser{
		Mail:       "blank@cap-net.co.jp",
		Name:       "部署なし",
		Department: "   ",
	}
	ldap.EXPECT().Authenticate("blank", "secret").Return(ldapUser, nil)

	deptRepo.EXPECT().GetByName(gomock.Any()).Times(0)
	deptRepo.EXPECT().Create(gomock.Any()).Times(0)

	userRepo.EXPECT().GetByEmail(ldapUser.Mail).Return(nil, gorm.ErrRecordNotFound)
	userRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(u *entity.User) error {
		u.ID = 9
		u.UUID = uuid.New()
		return nil
	})

	deptMember.EXPECT().GetByDepartmentID(uint(1)).Return(nil, nil)
	deptMember.EXPECT().ReplaceByDepartment(uint(1), gomock.Any()).Return(nil)

	uc := authenticate_usecase.NewAuthUsecase(ldap, userRepo, deptRepo, deptMember)
	_, _, err := uc.LoginLDAP(context.Background(), "blank@cap-net.co.jp", "secret")
	require.NoError(t, err)
}

func TestLoginLDAP_ExistingUser_CreatesDepartmentWithoutChangingMembership(t *testing.T) {
	ctrl := gomock.NewController(t)
	ldap := ldapmock.NewMockLDAPClient(ctrl)
	userRepo := repomock.NewMockUserRepository(ctrl)
	deptRepo := repomock.NewMockDepartmentRepository(ctrl)
	deptMember := repomock.NewMockDepartmentMemberRepository(ctrl)

	const deptName = "新規部署"
	ldapUser := &ldap_entity.LDAPUser{
		Mail:       "exist@cap-net.co.jp",
		Name:       "更新後",
		Department: deptName,
	}
	existing := &entity.User{
		Model:     gorm.Model{ID: 5},
		UUID:      uuid.New(),
		Name:      "更新前",
		Email:     ldapUser.Mail,
		Authority: valueobject.AuthorityGeneral,
	}
	ldap.EXPECT().Authenticate("exist", "secret").Return(ldapUser, nil)

	deptRepo.EXPECT().GetByName(deptName).Return(nil, gorm.ErrRecordNotFound)
	deptRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(d *entity.Department) error {
		d.ID = 50
		return nil
	})

	userRepo.EXPECT().GetByEmail(ldapUser.Mail).Return(existing, nil)
	userRepo.EXPECT().UpdateByID(gomock.Any()).Return(nil)

	deptMember.EXPECT().GetByDepartmentID(gomock.Any()).Times(0)
	deptMember.EXPECT().ReplaceByDepartment(gomock.Any(), gomock.Any()).Times(0)

	uc := authenticate_usecase.NewAuthUsecase(ldap, userRepo, deptRepo, deptMember)
	_, _, err := uc.LoginLDAP(context.Background(), "exist@cap-net.co.jp", "secret")
	require.NoError(t, err)
}