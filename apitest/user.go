package apitest

import (
	"math/rand"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/utils"
)

func CreateUser(context *context.Context, funcs ...func(user *models.User)) (*models.User, error) {
	var err error
	manager := models.NewUserManager(context)

	rand.Seed(time.Now().UnixNano())

	/*
		@TODO: big todo, create user should use rest api. som this will be doing two things
		first to create superuser and login under this account
	*/

	funcs = append(funcs, func(user *models.User) {
		user.Email = utils.RandomString(20) + randomdata.Email()
		user.IsActive = true
		user.Username = utils.RandomString(20)
		user.Name = randomdata.FullName(randomdata.Male)
		user.SetPassword("")
	})

	user := manager.NewUser(funcs...)
	err = user.Insert(context)
	return user, err
}

func CreateUserWithToken(context *context.Context, funcs ...func(user *models.User)) (user *models.User, token string, err error) {
	if user, err = CreateUser(context, funcs...); err != nil {
		return
	}
	manager := models.NewUserManager(context)
	token, err = manager.Login(user)
	return
}
