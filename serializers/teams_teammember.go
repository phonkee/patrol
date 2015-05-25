package serializers

import (
	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/validator"
	"github.com/phonkee/patrol/types"
)

/*
	Teammember list serializer
*/
type TeamsTeamMemberDetailSerializer struct {
	models.TeamMember
	User *AuthUserDetailSerializer `json:"user,omitempty"`
}

/*
Loads User from database, store it to serializer and returns it
*/
func (t *TeamsTeamMemberDetailSerializer) LoadUser(context *context.Context) (user *AuthUserDetailSerializer, err error) {
	manager := models.NewUserManager(context)
	t.User = &AuthUserDetailSerializer{}
	err = manager.GetByID(t.User, t.UserID)
	return t.User, err
}

/*
TeamsTeamMemberCreateSerializer
	create team member
*/
type TeamsTeamMemberCreateSerializer struct {
	Type   models.MemberType `json:"type"    validator:"type"`
	UserID types.ForeignKey  `json:"user_id" validator:"user_id"`
}

/*
Validate
	validates data
*/
func (t *TeamsTeamMemberCreateSerializer) Validate(context *context.Context, team *models.Team) *validator.Result {
	validator := validator.New()
	validator["user_id"] = models.ValidateUserID(context)
	validator["type"] = models.ValidateMemberType()

	result := validator.Validate(t)
	if !result.IsValid() {
		return result
	}

	tmm := models.NewTeamMemberManager(context)
	tmlist := tmm.NewTeamMemberList()
	if err := tmm.Filter(&tmlist, tmm.QueryFilterWhere("team_id = ? AND user_id = ?", team.ID, t.UserID)); err != nil {
		result.AddUnboundError(err)
	}

	if len(tmlist) > 0 {
		result.AddFieldError("user_id", ErrUserAlreadyMember)
	}

	return result
}

/*
	Saves team member to database
*/
func (t *TeamsTeamMemberCreateSerializer) Save(context *context.Context, team *models.Team) (result *models.TeamMember, err error) {

	tm := models.NewTeamMember(func(tm *models.TeamMember) {
		tm.TeamID = team.ID.ToForeignKey()
		tm.UserID = t.UserID
		tm.Type = t.Type
	})

	if err = tm.Insert(context); err != nil {
		return
	}

	return
}

/*
	TeamsTeamMemberUpdateSerializer
	Update serializer, updates only type.
*/
type TeamsTeamMemberUpdateSerializer struct {
	Type models.MemberType `json:"type"    validator:"type"`
}

/*
Validate method
*/
func (t *TeamsTeamMemberUpdateSerializer) Validate(context *context.Context) *validator.Result {
	v := validator.New()
	v["type"] = models.ValidateMemberType()
	return v.Validate(t)
}

/*
Updates team member
*/
func (t *TeamsTeamMemberUpdateSerializer) Update(context *context.Context, teammember *models.TeamMember) (err error) {
	teammember.Type = t.Type
	_, err = teammember.Update(context, "type")
	return
}
