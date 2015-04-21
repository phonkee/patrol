package models

import "github.com/phonkee/patrol/rest/validator"

const (
	UserUsernameMinLength = 3
	UserUsernameMaxLength = 32

	UserNameMinLength = 4
	UserNameMaxLength = 32

	UserPasswordMinLength = 5
)

func ValidateUserName() validator.ValidatorFunc {
	return validator.Any(
		validator.ValidateStringMinLength(UserNameMinLength),
		validator.ValidateStringMaxLength(UserNameMaxLength),
	)
}

func ValidateUserUsername() validator.ValidatorFunc {
	return validator.Any(
		validator.ValidateStringMinLength(UserUsernameMinLength),
		validator.ValidateStringMaxLength(UserUsernameMaxLength),
	)
}

func ValidatePassword() validator.ValidatorFunc {
	return validator.ValidateStringMinLength(UserPasswordMinLength)
}
