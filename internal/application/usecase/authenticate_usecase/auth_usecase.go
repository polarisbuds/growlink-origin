package authenticate_usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"growlink/internal/application/converter"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
	ldap_interface "growlink/internal/domain/interface/external/ldap_interface"
	repository_interface "growlink/internal/domain/interface/repository"
	"growlink/internal/domain/valueobject"
	"growlink/internal/shared/auth"
	"growlink/internal/shared/config"
	"growlink/internal/shared/crypto"
	logutils "growlink/internal/shared/logUtils"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthUsecase struct {
	ldapClient ldap_interface.LDAPClient
	userRepo   repository_interface.UserRepository
}

const opAuth = "auth_usecase"

func NewAuthUsecase(ldapClient ldap_interface.LDAPClient, userRepo repository_interface.UserRepository) *AuthUsecase {
	return &AuthUsecase{ldapClient: ldapClient, userRepo: userRepo}
}

func (uc *AuthUsecase) LoginLDAP(ctx context.Context, email, password string) (string, outputmodel.UserOutput, error) {
	const op = opAuth + ".LoginLDAP"
	logutils.Debug(ctx, logutils.LayerUsecase, op, "処理開始", slog.String("email", email))
	if uc.ldapClient == nil {
		return "", outputmodel.UserOutput{}, fmt.Errorf("LDAP認証は設定されていません")
	}
	loginName := strings.TrimSuffix(email, "@cap-net.co.jp")

	ldapUser, err := uc.ldapClient.Authenticate(loginName, password)
	if err != nil {
		logutils.Error(ctx, logutils.LayerUsecase, op, "LDAP認証失敗", slog.String("email", email), slog.String("err", err.Error()))
		return "", outputmodel.UserOutput{}, fmt.Errorf("LDAP認証に失敗しました: %w", err)
	}

	user, err := uc.userRepo.GetByEmail(ldapUser.Mail)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", outputmodel.UserOutput{}, err
		}
		user = &entity.User{
			Name:         ldapUser.Name,
			Email:        ldapUser.Mail,
			RoleID:       defaultRoleID(),
			DepartmentID: defaultDepartmentID(),
			Authority:    valueobject.AuthorityGeneral,
		}
		if err := uc.userRepo.Create(user); err != nil {
			return "", outputmodel.UserOutput{}, err
		}
	} else {
		user.Name = ldapUser.Name
		if err := uc.userRepo.UpdateByID(user); err != nil {
			return "", outputmodel.UserOutput{}, err
		}
	}
	if user.Suspended {
		return "", outputmodel.UserOutput{}, fmt.Errorf("アカウントは停止されています")
	}

	token, err := issueToken(user)
	if err != nil {
		return "", outputmodel.UserOutput{}, err
	}
	logutils.Info(ctx, logutils.LayerUsecase, op, "ログイン成功", slog.String("method", "ldap"), slog.Uint64("user_id", uint64(user.ID)))
	return token, converter.UserEntityToOutput(*user), nil
}

func (uc *AuthUsecase) LoginBasic(ctx context.Context, email, password string) (string, outputmodel.UserOutput, error) {
	const op = opAuth + ".LoginBasic"
	logutils.Debug(ctx, logutils.LayerUsecase, op, "処理開始", slog.String("email", email))
	user, err := uc.userRepo.GetByEmail(email)
	if err != nil {
		return "", outputmodel.UserOutput{}, fmt.Errorf("アカウントが存在しません")
	}
	if user.Password == nil {
		return "", outputmodel.UserOutput{}, fmt.Errorf("パスワードが入力されていません")
	}
	ok, err := crypto.VerifyPassword(password, config.GetPepper(), *user.Password)
	if err != nil || !ok {
		return "", outputmodel.UserOutput{}, fmt.Errorf("メールアドレスかパスワードが誤っています")
	}
	if user.Suspended {
		return "", outputmodel.UserOutput{}, fmt.Errorf("アカウントは停止されています")
	}
	token, err := issueToken(user)
	if err != nil {
		return "", outputmodel.UserOutput{}, err
	}
	logutils.Info(ctx, logutils.LayerUsecase, op, "ログイン成功", slog.String("method", "basic"), slog.Uint64("user_id", uint64(user.ID)))
	return token, converter.UserEntityToOutput(*user), nil
}

func issueToken(user *entity.User) (string, error) {
	claims := auth.CustomClaims{
		UserID:     user.ID,
		UUID:       user.UUID.String(),
		Email:      user.Email,
		Name:       user.Name,
		Authority:  user.Authority,
		SystemRole: user.SystemRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(config.GetJWTKey()))
}

func defaultRoleID() uint {
	if v := os.Getenv("LDAP_DEFAULT_ROLE_ID"); v != "" {
		var id uint
		if _, err := fmt.Sscan(v, &id); err == nil && id > 0 {
			return id
		}
	}
	return 1
}

func defaultDepartmentID() uint {
	if v := os.Getenv("LDAP_DEFAULT_DEPARTMENT_ID"); v != "" {
		var id uint
		if _, err := fmt.Sscan(v, &id); err == nil && id > 0 {
			return id
		}
	}
	return 1
}
