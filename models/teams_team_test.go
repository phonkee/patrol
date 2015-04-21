package models

import (
	"testing"

	"github.com/Pallinder/go-randomdata"
	"github.com/phonkee/patrol/utils"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTeam(t *testing.T) {
	context, err := utils.NewTestPatrolContext()
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	Convey("Test insert/update valid team", t, func() {

		var (
			err      error
			password = "password"
		)
		manager := NewUserManager(context)
		user, err := manager.GetTestUser(func(user *User) {
			user.SetPassword(password)
		})
		So(err, ShouldBeNil)

		teamm := NewTeamManager(context)
		team := teamm.NewTeam(func(te *Team) {
			te.Name = randomdata.SillyName()
			te.OwnerID = user.ID
		})
		err = team.Insert(teamm)
		So(err, ShouldBeNil)

		return

		err = team.Insert(teamm)
		So(err, ShouldBeNil)
		So(team.ID, ShouldNotEqual, 0)

		var teams []*Team
		err = teamm.Filter(&teams)
		// @TODO: implement this when teammembeer will be implemented
		// , teamm.QueryFilterUser(user)
		So(err, ShouldBeNil)
		So(len(teams), ShouldBeGreaterThan, 0)

	})

}
