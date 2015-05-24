package projects

import (
	"errors"

	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/validator"
	"github.com/phonkee/patrol/types"
	"github.com/phonkee/patrol/views/auth"
)

var (
	ErrUserAlreadyMember = errors.New("user_already_member")
)

type ProjectCreateSerializer struct {
	Name     string           `json:"name" validator:"name"`
	Platform string           `json:"platform"`
	TeamID   types.ForeignKey `json:"team_id" validator:"team_id"`
}

func (p ProjectCreateSerializer) Validate(context *context.Context) *validator.Result {
	validator := validator.New()
	validator["name"] = models.ValidateProjectName()
	validator["team_id"] = models.ValidateTeamID(context)
	return validator.Validate(p)
}

type ProjectMemberListItem struct {
	ID     types.PrimaryKey           `json:"id"`
	Type   models.MemberType          `json:"type"`
	UserID types.ForeignKey           `json:"user_id"`
	User   *auth.UserDetailSerializer `json:"user"`
}

type ProjectMemberCreate struct {
	Type   models.MemberType `json:"type" validator:"type"`
	UserID types.ForeignKey  `json:"user_id" validator:"user_id"`
}

func (p *ProjectMemberCreate) Validate(context *context.Context, teamid types.ForeignKey) *validator.Result {
	validator := validator.New()
	validator["user_id"] = models.ValidateUserID(context)
	validator["type"] = models.ValidateMemberType()

	result := validator.Validate(p)
	if !result.IsValid() {
		return result
	}

	tmm := models.NewTeamMemberManager(context)
	tmlist := tmm.NewTeamMemberList()
	if err := tmm.Filter(&tmlist, tmm.QueryFilterWhere("team_id = ? AND user_id = ?", teamid, p.UserID)); err != nil {
		result.AddUnboundError(err)
	}

	if len(tmlist) > 0 {
		result.AddFieldError("user_id", ErrUserAlreadyMember)
	}

	return result
}
