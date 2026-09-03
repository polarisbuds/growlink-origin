package converter

import (
	inputmodel "growlink/internal/application/model/input_model"
	outputmodel "growlink/internal/application/model/output_model"
	ldapentity "growlink/internal/domain/entity/ldap"
)

func LDAPUserInputToEntity(in inputmodel.LDAPUserInput) *ldapentity.LDAPUser {
	return &ldapentity.LDAPUser{
		UserName:   in.UserName,
		Name:       in.Name,
		Department: in.Department,
		EmployeeID: in.EmployeeID,
		Mail:       in.Mail,
		DN:         in.DN,
	}
}

func LDAPUserEntityToOutput(e ldapentity.LDAPUser) outputmodel.LDAPUserOutput {
	return outputmodel.LDAPUserOutput{
		UserName:   e.UserName,
		Name:       e.Name,
		Department: e.Department,
		EmployeeID: e.EmployeeID,
		Mail:       e.Mail,
		DN:         e.DN,
	}
}

func LDAPUserEntitiesToOutputs(es []ldapentity.LDAPUser) []outputmodel.LDAPUserOutput {
	outputs := make([]outputmodel.LDAPUserOutput, len(es))
	for i, e := range es {
		outputs[i] = LDAPUserEntityToOutput(e)
	}
	return outputs
}
