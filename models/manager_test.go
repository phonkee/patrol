package models

import (
	"testing"

	"github.com/phonkee/patrol/utils"
	. "github.com/smartystreets/goconvey/convey"
)

func TestManager(t *testing.T) {
	context, err := utils.NewTestPatrolContext()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	Convey("Test ChangedModelFields", t, func() {

		um := NewUserManager(context)
		tu, err := um.GetTestUser()
		So(err, ShouldBeNil)

		tu2 := um.NewUser()
		// clone object
		*tu2 = *tu

		tu2.Username = "test username"
		tu.Email = "testemail@example.com"

		fields, errf := ChangedModelFields(tu, tu2)
		So(errf, ShouldBeNil)
		So(fields, ShouldContain, "username")
		So(fields, ShouldContain, "email")

		pm := NewProjectManager(context)
		project := pm.NewProject()

		// different models
		_, errf = ChangedModelFields(tu, project)
		So(errf, ShouldNotBeNil)

	})

}
