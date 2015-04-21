package commands

import (
	"fmt"

	"code.google.com/p/gopass"
	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/core"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/validator"
)

/*
Commands
*/
func NewAuthCreateSuperuserCommand(context *context.Context) *AuthCreateSuperuserCommand {
	return &AuthCreateSuperuserCommand{
		context: context,
	}
}

type AuthCreateSuperuserCommand struct {
	core.Command
	context *context.Context
}

func (m *AuthCreateSuperuserCommand) ID() string          { return "create_superuser" }
func (m *AuthCreateSuperuserCommand) Description() string { return "Creates superuser" }
func (m *AuthCreateSuperuserCommand) Run() (err error) {

	var (
		username       string
		password       string
		passwordretype string
		email          string
	)

	fmt.Printf("Enter username: ")
	fmt.Scanf("%s", &username)

	fmt.Printf("Enter email: ")
	fmt.Scanf("%s", &email)

	if password, err = gopass.GetPass("Enter password: "); err != nil {
		return
	}
	if passwordretype, err = gopass.GetPass("Retype password: "); err != nil {
		return
	}

	if password != passwordretype {
		return fmt.Errorf("Passwords don't match!")
	}

	manager := models.NewUserManager(m.context)
	user := manager.NewUser(func(u *models.User) {
		u.Username = username
		u.Email = email
		u.IsActive = true
		u.IsSuperuser = true
		u.SetPassword(password)
	})

	var result *validator.Result
	result, err = user.Validate(m.context)

	if err != nil || !result.IsValid() {
		fmt.Printf("validation errors:\n")
		for field, errs := range result.Fields {
			fmt.Printf("    field: %s, errors: %v\n", field, errs)
		}

		if len(result.Unbound) > 0 {
			fmt.Printf("    unbound: %v\n", result.Unbound)
		}

		return
	}

	if err = user.Insert(m.context); err != nil {
		return
	}
	fmt.Printf("Superuser %s successfuly created.\n", user.Username)

	return
}
