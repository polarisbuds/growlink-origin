package authenticate_usecase

import (
	"fmt"
	"growlink/internal/application/converter"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/interface/external/ldap_interface"
)

type LoginUsecase struct {
	ldapClient ldap_interface.LDAPClient
}

func NewLoginUsecase(ldapClient ldap_interface.LDAPClient) *LoginUsecase {
	return &LoginUsecase{
		ldapClient: ldapClient,
	}
}

func (uc *LoginUsecase) Authenticate(id, password string) (*outputmodel.LDAPUserOutput, error) {
	ldapUser, err := uc.ldapClient.Authenticate(id, password)
	if err != nil {
		return nil, fmt.Errorf("LDAP認証に失敗しました: %w", err)
	}
	user := converter.LDAPUserEntityToOutput(*ldapUser)
	return &user, nil
}
