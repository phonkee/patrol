package models

import (
	"errors"

	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/rest/validator"
	"github.com/phonkee/patrol/types"
)

var (
	ErrInvalidTeamID = errors.New("invalid_team")
)

/*
Validate project name
*/
func ValidateProjectName() validator.ValidatorFunc {
	return validator.Any(
		validator.ValidateStringMinLength(5),
		validator.ValidateStringMaxLength(255),
	)
}

/*
Validate team id
*/
func ValidateTeamID(context *context.Context) validator.ValidatorFunc {
	tm := NewTeamManager(context)

	return func(value interface{}) (err error) {
		team := tm.NewTeam()
		teamid := value.(types.ForeignKey)
		if err = tm.GetByID(team, teamid); err != nil {
			if err == ErrObjectDoesNotExists {
				return ErrInvalidTeamID
			}
		}
		return
	}
}
