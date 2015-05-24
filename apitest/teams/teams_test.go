package teams

// import (
// 	"testing"

// 	"github.com/phonkee/patrol"
// 	"github.com/phonkee/patrol/models"
// 	. "github.com/smartystreets/goconvey/convey"
// )

// func trunc() {
// 	patrol.Context.DB.Queryx("TRUNCATE TABLE " + models.TEAMS_TEAM_DB_TABLE + " CASCADE")
// 	patrol.Context.DB.Queryx("TRUNCATE TABLE " + models.TEAMS_TEAMMEMBER_DB_TABLE + " CASCADE")
// }

// func TestTeams(t *testing.T) {
// 	// if err := patrol.Setup(); err != nil {
// 	// 	if err != patrol.ErrPatrolAlreadySetup {
// 	// 		fmt.Printf("patrol: setup failed with error: %s", err)
// 	// 	}
// 	// }

// 	// patrol.Run([]string{"migrate"})

// 	// Convey("Test options/metadata", t, func() {
// 	// 	sudoSession := NewRequestsSession(patrol.Context).WithNewUser(func(user *models.User) {
// 	// 		user.Password = "password"
// 	// 		user.IsSuperuser = true
// 	// 	})
// 	// 	req := sudoSession.Request("OPTIONS", settings.ROUTE_TEAMS_TEAM_LIST).Do()
// 	// 	So(req.Response().Code, ShouldEqual, http.StatusOK)
// 	// })

// 	// Convey("Test list teams (team members)", t, func() {
// 	// 	trunc()

// 	// 	sudoSession := NewRequestsSession(patrol.Context).WithNewUser(func(user *models.User) {
// 	// 		user.Password = "password"
// 	// 		user.IsSuperuser = true
// 	// 	})
// 	// 	req := sudoSession.Request("GET", settings.ROUTE_TEAMS_TEAM_LIST).Do()
// 	// 	So(req.Response().Code, ShouldEqual, http.StatusOK)

// 	// 	reqNew := sudoSession.Request("POST", settings.ROUTE_TEAMS_TEAM_LIST)

// 	// 	s := teams.TeamCreateSerializer{Name: "new team"}

// 	// 	reqNew.JSONBody(s).Do()

// 	// 	x := struct {
// 	// 		Result *models.Team `json:"result"`
// 	// 	}{}

// 	// 	reqNew.Scan(&x)

// 	// 	So(reqNew.Response().Code, ShouldEqual, http.StatusCreated)
// 	// 	So(x.Result.Name, ShouldEqual, s.Name)
// 	// 	So(x.Result.ID, ShouldNotEqual, 0)

// 	// 	normalSession := NewRequestsSession(patrol.Context).WithNewUser(func(user *models.User) {
// 	// 		user.IsSuperuser = false
// 	// 	})

// 	// 	y := struct {
// 	// 		Result []*models.Team `json:"result"`
// 	// 	}{}

// 	// 	reqNoTeams := normalSession.Request("POST", settings.ROUTE_TEAMS_TEAM_LIST).Do().Scan(&y)
// 	// 	So(reqNoTeams.Error(), ShouldBeNil)
// 	// 	So(len(y.Result), ShouldEqual, 0)

// 	// 	unauthorized := NewRequestsSession(patrol.Context)
// 	// 	So(
// 	// 		unauthorized.Request("POST", settings.ROUTE_TEAMS_TEAM_LIST).Do().Response().Code,
// 	// 		ShouldEqual, http.StatusUnauthorized,
// 	// 	)

// 	// 	So(
// 	// 		unauthorized.Request("GET", settings.ROUTE_TEAMS_TEAM_LIST).Do().Response().Code,
// 	// 		ShouldEqual, http.StatusUnauthorized,
// 	// 	)

// 	// })

// 	// Convey("test create new team", t, func() {

// 	// 	Convey("test unauthorized", func() {
// 	// 		normal := NewRequestsSession(patrol.Context)
// 	// 		code := normal.Request("POST", settings.ROUTE_TEAMS_TEAM_LIST).
// 	// 			JSONBody(teams.TeamCreateSerializer{Name: "new team"}).Do().Response().Code
// 	// 		So(code, ShouldEqual, http.StatusUnauthorized)
// 	// 	})

// 	// 	Convey("test forbidden access", func() {
// 	// 		normal := NewRequestsSession(patrol.Context).WithNewUser()
// 	// 		code := normal.Request("POST", settings.ROUTE_TEAMS_TEAM_LIST).
// 	// 			JSONBody(teams.TeamCreateSerializer{Name: "new team"}).Do().Response().Code
// 	// 		So(code, ShouldEqual, http.StatusForbidden)
// 	// 	})

// 	// 	Convey("test sudo access", func() {
// 	// 		sudo := NewRequestsSession(patrol.Context).WithNewUser(func(user *models.User) {
// 	// 			user.IsSuperuser = true
// 	// 		})

// 	// 		req := sudo.Request("POST", settings.ROUTE_TEAMS_TEAM_LIST).JSONBody(teams.TeamCreateSerializer{Name: "new team"}).Do()
// 	// 		So(req.Response().Code, ShouldEqual, http.StatusCreated)
// 	// 		x := struct {
// 	// 			Result *models.Team `json:"result"`
// 	// 		}{}

// 	// 		req.Scan(&x)
// 	// 		So(x.Result.PrimaryKey(), ShouldBeGreaterThan, 0)

// 	// 	})

// 	// })

// 	// Convey("test get team", t, func() {

// 	// 	sudo := NewRequestsSession(patrol.Context).WithNewUser(func(user *models.User) { user.IsSuperuser = true })
// 	// 	normal := NewRequestsSession(patrol.Context).WithNewUser()
// 	// 	unauthorized := NewRequestsSession(patrol.Context)

// 	// 	_, _, _ = sudo, normal, unauthorized

// 	// 	req := sudo.Request("POST", settings.ROUTE_TEAMS_TEAM_LIST).JSONBody(teams.TeamCreateSerializer{Name: "new team"}).Do()
// 	// 	So(req.Response().Code, ShouldEqual, http.StatusCreated)
// 	// 	x := struct {
// 	// 		Result *models.Team `json:"result"`
// 	// 	}{}

// 	// 	So(req.Error(), ShouldBeNil)
// 	// 	req.Scan(&x)

// 	// 	Convey("test unauthorized", func() {
// 	// 		code := unauthorized.Request("POST", settings.ROUTE_TEAMS_TEAM_DETAIL, "team_id", x.Result.ID.String()).Do().Response().Code
// 	// 		So(code, ShouldEqual, http.StatusUnauthorized)
// 	// 	})

// 	// 	Convey("test sudo", func() {
// 	// 		code := sudo.Request("GET", settings.ROUTE_TEAMS_TEAM_DETAIL, "team_id", x.Result.ID.String()).Do().Response().Code
// 	// 		So(code, ShouldEqual, http.StatusOK)

// 	// 		y := struct {
// 	// 			Result *models.Team `json:"result"`
// 	// 		}{}

// 	// 		So(req.Error(), ShouldBeNil)
// 	// 		req.Scan(&y)
// 	// 		So(y.Result.ID, ShouldEqual, x.Result.ID)
// 	// 	})

// 	// })

// }
