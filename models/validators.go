package models

import (
	"errors"

	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/rest/validator"
	"github.com/phonkee/patrol/types"
)

var (
	ErrInvalidTeamID     = errors.New("invalid_team")
	ErrInvalidUserID     = errors.New("invalid_user")
	ErrInvalidMemberType = errors.New("invalid_member_type")
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

/*
Validate team name
*/
func ValidateTeamName() validator.ValidatorFunc {
	return validator.Any(
		validator.ValidateStringMinLength(5),
		validator.ValidateStringMaxLength(200),
	)
}

/*
Validate user id
*/
func ValidateUserID(context *context.Context) validator.ValidatorFunc {
	um := NewUserManager(context)

	return func(value interface{}) (err error) {
		user := um.NewUser()
		id := value.(types.ForeignKey)
		if err = um.GetByID(user, id); err != nil {
			if err == ErrObjectDoesNotExists {
				return ErrInvalidUserID
			}
		}
		return
	}
}

func ValidateMemberType() validator.ValidatorFunc {
	return func(value interface{}) (err error) {
		mt := value.(MemberType)
		if !mt.IsValid() {
			return ErrInvalidMemberType
		}
		return
	}
}
