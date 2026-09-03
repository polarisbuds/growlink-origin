package converter

import (
	inputmodel "growlink/internal/application/model/input_model"
	outputmodel "growlink/internal/application/model/output_model"
	"growlink/internal/domain/entity"
	ldapentity "growlink/internal/domain/entity/ldap"
)

// --- Input -> Entity ---

func MeetingInputToEntity(in inputmodel.MeetingInput) *entity.Meeting {
	return &entity.Meeting{
		ImplementedAt: in.ImplementedAt,
		Status:        in.Status,
		TemplateID:    in.TemplateID,
		PairingID:     in.PairingID,
	}
}

func TemplateInputToEntity(in inputmodel.TemplateInput) *entity.Template {
	topics := make([]entity.Topic, len(in.Topics))
	for i, topic := range in.Topics {
		topics[i] = *TopicInputToEntity(topic)
	}
	return &entity.Template{
		Title:  in.Title,
		Topics: topics,
	}
}

func TopicInputToEntity(in inputmodel.TopicInput) *entity.Topic {
	return &entity.Topic{
		DisplayNo: in.DisplayNo,
		Question:  in.Question,
		IsRequire: in.IsRequire,
	}
}

func PairingInputToEntity(in inputmodel.PairingInput) *entity.Pairing {
	return &entity.Pairing{
		MentorID: in.MentorID,
		MenteeID: in.MenteeID,
	}
}

func UserInputToEntity(in inputmodel.UserInput) *entity.User {
	return &entity.User{
		Name:         in.Name,
		Email:        in.Email,
		ImageURI:     in.ImageURI,
		RoleID:       in.RoleID,
		DepartmentID: in.DepartmentID,
		Authority:    in.Authority,
	}
}

func AnswerInputToEntity(in inputmodel.AnswerInput) *entity.Answer {
	return &entity.Answer{
		Content:   in.Content,
		MeetingID: in.MeetingID,
		TopicID:   in.TopicID,
		UserID:    in.UserID,
	}
}

func AnswerCommentInputToEntity(in inputmodel.AnswerCommentInput) *entity.AnswerComment {
	return &entity.AnswerComment{
		Content:  in.Content,
		AnswerID: in.AnswerID,
		UserID:   in.UserID,
	}
}

func MeetingCommentInputToEntity(in inputmodel.MeetingCommentInput) *entity.MeetingComment {
	return &entity.MeetingComment{
		Content:   in.Content,
		MeetingID: in.MeetingID,
		UserID:    in.UserID,
	}
}

func ReadCommentableInputToEntity(in inputmodel.ReadCommentableInput) *entity.ReadCommentable {
	return &entity.ReadCommentable{
		MeetingID: in.MeetingID,
		UserID:    in.UserID,
	}
}

func UserRelationshipInputToEntity(in inputmodel.UserRelationshipInput) *entity.UserRelationship {
	return &entity.UserRelationship{
		RelationType: in.RelationType,
		SourceUserID: in.SourceUserID,
		TargetUserID: in.TargetUserID,
	}
}

func TeamInputToEntity(in inputmodel.TeamInput) *entity.Team {
	return &entity.Team{
		Name:         in.Name,
		DepartmentID: in.DepartmentID,
	}
}

func TeamMemberInputToEntity(in inputmodel.TeamMemberInput) *entity.TeamMember {
	return &entity.TeamMember{
		IsReader: in.IsReader,
		TeamID:   in.TeamID,
		MemberID: in.MemberID,
	}
}

func RoleInputToEntity(in inputmodel.RoleInput) *entity.Role {
	return &entity.Role{Name: in.Name}
}

func DepartmentInputToEntity(in inputmodel.DepartmentInput) *entity.Department {
	return &entity.Department{Name: in.Name}
}

func HistoryInputToEntity(in inputmodel.HistoryInput) *entity.History {
	return &entity.History{
		UserID:    in.UserID,
		EventType: in.EventType,
	}
}

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

// --- Entity -> Output (leaf) ---

func RoleEntityToOutput(e entity.Role) outputmodel.RoleOutput {
	return outputmodel.RoleOutput{ID: e.ID, Name: e.Name}
}

func RoleEntitiesToOutputs(es []entity.Role) []outputmodel.RoleOutput {
	outputs := make([]outputmodel.RoleOutput, len(es))
	for i, e := range es {
		outputs[i] = RoleEntityToOutput(e)
	}
	return outputs
}

func DepartmentEntityToOutput(e entity.Department) outputmodel.DepartmentOutput {
	return outputmodel.DepartmentOutput{ID: e.ID, Name: e.Name}
}

func DepartmentEntitiesToOutputs(es []entity.Department) []outputmodel.DepartmentOutput {
	outputs := make([]outputmodel.DepartmentOutput, len(es))
	for i, e := range es {
		outputs[i] = DepartmentEntityToOutput(e)
	}
	return outputs
}

func UserEntityToOutput(e entity.User) outputmodel.UserOutput {
	return outputmodel.UserOutput{
		ID:       e.ID,
		Name:     e.Name,
		Email:    e.Email,
		ImageURI: e.ImageURI,
	}
}

func UserEntitiesToOutputs(es []entity.User) []outputmodel.UserOutput {
	outputs := make([]outputmodel.UserOutput, len(es))
	for i, e := range es {
		outputs[i] = UserEntityToOutput(e)
	}
	return outputs
}

func TopicEntityToOutput(e entity.Topic) outputmodel.TopicOutput {
	return outputmodel.TopicOutput{
		ID:        e.ID,
		DisplayNo: e.DisplayNo,
		Question:  e.Question,
		IsRequire: e.IsRequire,
	}
}

func TopicEntitiesToOutputs(es []entity.Topic) []outputmodel.TopicOutput {
	outputs := make([]outputmodel.TopicOutput, len(es))
	for i, e := range es {
		outputs[i] = TopicEntityToOutput(e)
	}
	return outputs
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

// --- Entity -> Output (intermediate) ---

func TemplateEntityToOutput(e entity.Template) outputmodel.TemplateOutput {
	return outputmodel.TemplateOutput{
		ID:     e.ID,
		Title:  e.Title,
		Topics: TopicEntitiesToOutputs(e.Topics),
	}
}

func TemplateEntitiesToOutputs(es []entity.Template) []outputmodel.TemplateOutput {
	outputs := make([]outputmodel.TemplateOutput, len(es))
	for i, e := range es {
		outputs[i] = TemplateEntityToOutput(e)
	}
	return outputs
}

func PairingEntityToOutput(e entity.Pairing) outputmodel.PairingOutput {
	out := outputmodel.PairingOutput{ID: e.ID}
	if e.Mentor != nil {
		out.Mentor = UserEntityToOutput(*e.Mentor)
	}
	if e.Mentee != nil {
		out.Mentee = UserEntityToOutput(*e.Mentee)
	}
	return out
}

func PairingEntitiesToOutputs(es []entity.Pairing) []outputmodel.PairingOutput {
	outputs := make([]outputmodel.PairingOutput, len(es))
	for i, e := range es {
		outputs[i] = PairingEntityToOutput(e)
	}
	return outputs
}

func teamEntityToOutputShallow(e entity.Team) outputmodel.TeamOutput {
	out := outputmodel.TeamOutput{
		ID:           e.ID,
		Name:         e.Name,
		DepartmentID: e.DepartmentID,
	}
	if e.Department != nil {
		out.Department = DepartmentEntityToOutput(*e.Department)
	}
	return out
}

func TeamMemberEntityToOutput(e entity.TeamMember) outputmodel.TeamMemberOutput {
	out := outputmodel.TeamMemberOutput{
		ID:       e.ID,
		IsReader: e.IsReader,
		TeamID:   e.TeamID,
		MemberID: e.MemberID,
	}
	if e.Team != nil {
		out.Team = teamEntityToOutputShallow(*e.Team)
	} else if e.TeamID != 0 {
		out.Team = outputmodel.TeamOutput{ID: e.TeamID}
	}
	if e.Member != nil {
		out.Member = UserEntityToOutput(*e.Member)
	}
	return out
}

func TeamMemberEntitiesToOutputs(es []entity.TeamMember) []outputmodel.TeamMemberOutput {
	outputs := make([]outputmodel.TeamMemberOutput, len(es))
	for i, e := range es {
		outputs[i] = TeamMemberEntityToOutput(e)
	}
	return outputs
}

func TeamEntityToOutput(e entity.Team) outputmodel.TeamOutput {
	out := teamEntityToOutputShallow(e)
	if len(e.TeamMembers) > 0 {
		out.TeamMembers = TeamMemberEntitiesToOutputs(e.TeamMembers)
	}
	return out
}

func TeamEntitiesToOutputs(es []entity.Team) []outputmodel.TeamOutput {
	outputs := make([]outputmodel.TeamOutput, len(es))
	for i, e := range es {
		outputs[i] = TeamEntityToOutput(e)
	}
	return outputs
}

func UserRelationshipEntityToOutput(e entity.UserRelationship) outputmodel.UserRelationshipOutput {
	out := outputmodel.UserRelationshipOutput{
		ID:           e.ID,
		RelationType: e.RelationType,
		SourceUserID: e.SourceUserID,
		TargetUserID: e.TargetUserID,
	}
	if e.SourceUser != nil {
		out.SourceUser = UserEntityToOutput(*e.SourceUser)
	}
	if e.TargetUser != nil {
		out.TargetUser = UserEntityToOutput(*e.TargetUser)
	}
	return out
}

func UserRelationshipEntitiesToOutputs(es []entity.UserRelationship) []outputmodel.UserRelationshipOutput {
	outputs := make([]outputmodel.UserRelationshipOutput, len(es))
	for i, e := range es {
		outputs[i] = UserRelationshipEntityToOutput(e)
	}
	return outputs
}

func HistoryEntityToOutput(e entity.History) outputmodel.HistoryOutput {
	out := outputmodel.HistoryOutput{
		ID:        e.ID,
		UserID:    e.UserID,
		EventType: e.EventType,
	}
	if e.User != nil {
		out.User = UserEntityToOutput(*e.User)
	}
	return out
}

func HistoryEntitiesToOutputs(es []entity.History) []outputmodel.HistoryOutput {
	outputs := make([]outputmodel.HistoryOutput, len(es))
	for i, e := range es {
		outputs[i] = HistoryEntityToOutput(e)
	}
	return outputs
}

// --- Entity -> Output (complex) ---

func MeetingEntityToOutput(e entity.Meeting) outputmodel.MeetingOutput {
	out := outputmodel.MeetingOutput{
		ID:            e.ID,
		ImplementedAt: e.ImplementedAt,
		Status:        e.Status,
	}
	if e.Template != nil {
		out.Template = TemplateEntityToOutput(*e.Template)
	} else if e.TemplateID != 0 {
		out.Template = outputmodel.TemplateOutput{ID: e.TemplateID}
	}
	if e.Pairing != nil {
		out.Pairing = PairingEntityToOutput(*e.Pairing)
	} else if e.PairingID != 0 {
		out.Pairing = outputmodel.PairingOutput{ID: e.PairingID}
	}
	return out
}

func MeetingEntityToOutputPtr(e *entity.Meeting) *outputmodel.MeetingOutput {
	if e == nil {
		return nil
	}
	out := MeetingEntityToOutput(*e)
	return &out
}

func MeetingEntitiesToOutputs(es []entity.Meeting) []outputmodel.MeetingOutput {
	outputs := make([]outputmodel.MeetingOutput, len(es))
	for i, e := range es {
		outputs[i] = MeetingEntityToOutput(e)
	}
	return outputs
}

func answerEntityToOutputShallow(e entity.Answer) outputmodel.AnswerOutput {
	out := outputmodel.AnswerOutput{
		ID:        e.ID,
		Content:   e.Content,
		MeetingID: e.MeetingID,
		TopicID:   e.TopicID,
		UserID:    e.UserID,
	}
	if e.Meeting != nil {
		out.Meeting = MeetingEntityToOutput(*e.Meeting)
	} else if e.MeetingID != 0 {
		out.Meeting = outputmodel.MeetingOutput{ID: e.MeetingID}
	}
	if e.Topic != nil {
		out.Topic = TopicEntityToOutput(*e.Topic)
	} else if e.TopicID != 0 {
		out.Topic = outputmodel.TopicOutput{ID: e.TopicID}
	}
	if e.User != nil {
		out.User = UserEntityToOutput(*e.User)
	}
	return out
}

func AnswerEntityToOutput(e entity.Answer) outputmodel.AnswerOutput {
	out := answerEntityToOutputShallow(e)
	if len(e.AnswerComments) > 0 {
		out.AnswerComments = AnswerCommentEntitiesToOutputs(e.AnswerComments)
	}
	return out
}

func AnswerEntitiesToOutputs(es []entity.Answer) []outputmodel.AnswerOutput {
	outputs := make([]outputmodel.AnswerOutput, len(es))
	for i, e := range es {
		outputs[i] = AnswerEntityToOutput(e)
	}
	return outputs
}

func AnswerCommentEntityToOutput(e entity.AnswerComment) outputmodel.AnswerCommentOutput {
	out := outputmodel.AnswerCommentOutput{
		ID:       e.ID,
		Content:  e.Content,
		AnswerID: e.AnswerID,
		UserID:   e.UserID,
	}
	if e.Answer != nil {
		out.Answer = answerEntityToOutputShallow(*e.Answer)
	} else if e.AnswerID != 0 {
		out.Answer = outputmodel.AnswerOutput{ID: e.AnswerID}
	}
	if e.User != nil {
		out.User = UserEntityToOutput(*e.User)
	}
	return out
}

func AnswerCommentEntitiesToOutputs(es []entity.AnswerComment) []outputmodel.AnswerCommentOutput {
	outputs := make([]outputmodel.AnswerCommentOutput, len(es))
	for i, e := range es {
		outputs[i] = AnswerCommentEntityToOutput(e)
	}
	return outputs
}

func MeetingCommentEntityToOutput(e entity.MeetingComment) outputmodel.MeetingCommentOutput {
	out := outputmodel.MeetingCommentOutput{
		ID:        e.ID,
		Content:   e.Content,
		MeetingID: e.MeetingID,
		UserID:    e.UserID,
	}
	if e.Meeting != nil {
		out.Meeting = MeetingEntityToOutput(*e.Meeting)
	} else if e.MeetingID != 0 {
		out.Meeting = outputmodel.MeetingOutput{ID: e.MeetingID}
	}
	if e.User != nil {
		out.User = UserEntityToOutput(*e.User)
	}
	return out
}

func MeetingCommentEntitiesToOutputs(es []entity.MeetingComment) []outputmodel.MeetingCommentOutput {
	outputs := make([]outputmodel.MeetingCommentOutput, len(es))
	for i, e := range es {
		outputs[i] = MeetingCommentEntityToOutput(e)
	}
	return outputs
}

func ReadCommentableEntityToOutput(e entity.ReadCommentable) outputmodel.ReadCommentableOutput {
	out := outputmodel.ReadCommentableOutput{
		ID:        e.ID,
		MeetingID: e.MeetingID,
		UserID:    e.UserID,
	}
	if e.Meeting != nil {
		out.Meeting = MeetingEntityToOutput(*e.Meeting)
	} else if e.MeetingID != 0 {
		out.Meeting = outputmodel.MeetingOutput{ID: e.MeetingID}
	}
	if e.User != nil {
		out.User = UserEntityToOutput(*e.User)
	}
	return out
}

func ReadCommentableEntitiesToOutputs(es []entity.ReadCommentable) []outputmodel.ReadCommentableOutput {
	outputs := make([]outputmodel.ReadCommentableOutput, len(es))
	for i, e := range es {
		outputs[i] = ReadCommentableEntityToOutput(e)
	}
	return outputs
}
