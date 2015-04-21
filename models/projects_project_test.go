package models

import (
	"testing"

	"github.com/Pallinder/go-randomdata"
	"github.com/phonkee/patrol/utils"
	. "github.com/smartystreets/goconvey/convey"
)

func TestInsertProject(t *testing.T) {
	context, err := utils.NewTestPatrolContext()
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	Convey("Test insert valid project", t, func() {
		var err error
		manager := NewProjectManager(context)
		project := manager.NewProject()
		project.Name = "test project"
		project.Platform = "go"

		validator, errValidator := project.Validate(context)
		So(errValidator, ShouldBeNil)

		valid := validator.IsValid()
		So(valid, ShouldBeFalse)

		So(project.ID, ShouldEqual, 0)
		err = project.Insert(manager)
		// team id not given
		So(err, ShouldNotBeNil)

		validator.AddPostgresError(err)

		return

		// So(project.ID, ShouldNotEqual, 0)

		// project.Platform = "python"

		// var changed bool

		// changed, err = project.Update(manager)
		// So(err, ShouldBeNil)
		// So(changed, ShouldBeTrue)

		// projectNew := manager.NewProject()
		// err = manager.Get(projectNew, manager.QueryFilterID(project.ID))
		// So(err, ShouldBeNil)
		// So(project.ID, ShouldEqual, projectNew.ID)
	})

	Convey("List projects + delete project", t, func() {
		password := "password"

		projectmanager := NewProjectManager(context)
		teammanager := NewTeamManager(context)
		teammembermanager := NewTeamMemberManager(context)
		usermanager := NewUserManager(context)

		owner, errOwner := usermanager.GetTestUser(func(user *User) {
			user.SetPassword(password)
			user.IsActive = true
			user.IsSuperuser = true
		})
		So(errOwner, ShouldBeNil)

		user, err := usermanager.GetTestUser(func(user *User) {
			user.SetPassword(password)
			user.IsActive = true
			user.IsSuperuser = false
		})
		So(err, ShouldBeNil)

		team := teammanager.NewTeam(func(team *Team) {
			team.Name = "some team"
			team.OwnerID = owner.ID
		})
		errTeamInsert := team.Insert(teammanager)
		So(errTeamInsert, ShouldBeNil)

		err = teammembermanager.SetTeamMemberType(team, user, MEMBER_TYPE_MEMBER)
		So(err, ShouldBeNil)

		project := projectmanager.NewProject(func(p *Project) {
			p.Platform = "go"
			p.Name = randomdata.SillyName()
			p.TeamID = team.ID
		})

		errProjectInsert := project.Insert(teammanager)
		So(errProjectInsert, ShouldBeNil)

		projects := projectmanager.NewProjectList()
		errFilter := projectmanager.Filter(&projects, projectmanager.QueryFilterUser(user))
		So(errFilter, ShouldBeNil)
		So(len(projects), ShouldEqual, 1)

		another, errAnother := usermanager.GetTestUser(func(user *User) {
			user.SetPassword(password)
			user.IsActive = true
			user.IsSuperuser = false
		})
		So(errAnother, ShouldBeNil)

		projects = projectmanager.NewProjectList()
		errFilter = projectmanager.Filter(&projects, projectmanager.QueryFilterUser(another))
		So(errFilter, ShouldBeNil)
		So(len(projects), ShouldEqual, 0)

		superuser, errSuperuser := usermanager.GetTestUser(func(user *User) {
			user.SetPassword(password)
			user.IsActive = true
			user.IsSuperuser = true
		})
		So(errSuperuser, ShouldBeNil)

		projects = projectmanager.NewProjectList()
		errFilter = projectmanager.Filter(&projects, projectmanager.QueryFilterUser(superuser))
		So(errFilter, ShouldBeNil)

		suprojectsLen := len(projects)

		projects = projectmanager.NewProjectList()
		errFilter = projectmanager.Filter(&projects)
		So(errFilter, ShouldBeNil)

		So(len(projects), ShouldEqual, suprojectsLen)

		allProjects := len(projects)

		errDelete := project.Delete(context)
		So(errDelete, ShouldBeNil)
		So(project.ID, ShouldEqual, 0)

		projects = projectmanager.NewProjectList()
		errFilter = projectmanager.Filter(&projects)
		So(errFilter, ShouldBeNil)

		So(allProjects, ShouldEqual, len(projects)+1)
	})

	Convey("Test update only defined fields", t, func() {

	})
}
