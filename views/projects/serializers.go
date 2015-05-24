package projects

import (
	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/validator"
	"github.com/phonkee/patrol/types"
	"github.com/phonkee/patrol/views/auth"
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
	Type   models.MemberType `json:"type"`
	UserID types.PrimaryKey  `json:"user_id"`
}
