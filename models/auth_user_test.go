package models

import (
	"testing"

	"github.com/Pallinder/go-randomdata"
	"github.com/phonkee/patrol/utils"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUser(t *testing.T) {
	context, err := utils.NewTestPatrolContext()
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	Convey("Test insert/update valid user", t, func() {
		var err error
		var changed bool
		password := utils.RandomString(32)
		manager := NewUserManager(context)

		user, err := manager.GetTestUnsavedUser(func(user *User) {
			user.SetPassword(password)
		})
		So(err, ShouldBeNil)

		_, err = user.Update(manager)
		So(err, ShouldNotBeNil)

		err = user.Insert(manager)

		So(err, ShouldBeNil)
		So(user.ID, ShouldNotEqual, 0)

		err = user.Insert(manager)
		So(err, ShouldNotBeNil)

		user.Name = randomdata.FullName(randomdata.Male)
		changed, err = user.Update(manager, "name")
		So(err, ShouldBeNil)
		So(changed, ShouldBeTrue)

		changed, err = user.Update(manager, "name", "last_login")
		So(err, ShouldBeNil)
		_, err = user.Update(context, "nonexistingfield")
		So(err, ShouldNotBeNil)

	})

	Convey("login user", t, func() {
		var err error
		var changed bool
		password := utils.RandomString(32)
		manager := NewUserManager(context)
		user := manager.NewUser(func(u *User) {
			u.Email = utils.RandomString(20) + randomdata.Email()
			u.IsActive = false
			u.Username = utils.RandomString(20)
			u.Name = randomdata.FullName(randomdata.Male)
			u.SetPassword(password)
		})

		err = user.Insert(manager)
		So(err, ShouldBeNil)

		_, err = manager.Login(user)
		So(err, ShouldNotBeNil)
		So(err, ShouldEqual, ErrCannotLoginUser)

		user.IsActive = true
		changed, err = user.Update(manager, "is_active")
		So(changed, ShouldNotEqual, 0)
		So(err, ShouldBeNil)

		_, err = manager.Login(user)
		So(err, ShouldBeNil)

		blank := manager.NewUser()
		_, err = manager.Login(blank)
		So(err, ShouldNotBeNil)

	})

	Convey("test verify password", t, func() {
		var err error
		password := "dfaspekur3irkjf;lsadkjf"
		manager := NewUserManager(context)
		user, err := manager.GetTestUser(func(user *User) {
			user.SetPassword(password)
		})

		So(err, ShouldBeNil)
		var verified bool
		verified, err = user.VerifyPassword(password)
		So(verified, ShouldBeTrue)

	})

	Convey("test get", t, func() {
		manager := NewUserManager(context)
		user, err := manager.GetTestUser()

		fetched := manager.NewUser()

		err = manager.GetByID(fetched, user.ID)
		So(err, ShouldBeNil)
		So(user.Email, ShouldEqual, fetched.Email)

		blank := manager.NewUser()
		err = manager.Get(blank, manager.QueryFilterEmail(user.Email))
		So(err, ShouldBeNil)

		blank = manager.NewUser()
		err = manager.Get(blank, manager.QueryFilterUsername(user.Username))
		So(err, ShouldBeNil)

	})

	Convey("test new user", t, func() {
		email := "testing@example.com"
		name := "Poweruser"
		manager := NewUserManager(context)
		user := manager.NewUser(func(user *User) {
			user.Email = email
		}, func(user *User) {
			user.Name = name
		})

		So(user.Email, ShouldEqual, email)
		So(user.Name, ShouldEqual, name)
	})

}
