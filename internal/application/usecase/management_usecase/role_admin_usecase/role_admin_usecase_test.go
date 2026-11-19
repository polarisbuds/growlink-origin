package role_admin_usecase_test

import (
	"testing"

	role_admin_usecase "growlink/internal/application/usecase/management_usecase/role_admin_usecase"
	"growlink/internal/domain/entity"
	repomock "growlink/internal/domain/interface/repository/mock"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestRoleAdmin_Create_and_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := repomock.NewMockRoleRepository(ctrl)
	uc := role_admin_usecase.NewRoleAdmin(repo)

	repo.EXPECT().GetByName("Lead").Return(nil, nil)
	repo.EXPECT().Create(gomock.Any()).DoAndReturn(func(r *entity.Role) error {
		r.UUID = uuid.New()
		return nil
	})

	out, err := uc.Create("Lead")
	require.NoError(t, err)
	assert.Equal(t, "Lead", out.Name)

	repo.EXPECT().GetAll().Return([]entity.Role{{Name: "Lead", Model: gorm.Model{ID: 1}}}, nil)
	list, err := uc.List()
	require.NoError(t, err)
	require.Len(t, list, 1)
}