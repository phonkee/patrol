package projects

import (
	"net/http"

	"github.com/golang/glog"
	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/core"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/metadata"
	"github.com/phonkee/patrol/rest/response"
	"github.com/phonkee/patrol/rest/validator"
	"github.com/phonkee/patrol/settings"
	"github.com/phonkee/patrol/types"
)

type ProjectCreateSerializer struct {
	Name     string           `json:"name" validator:"name"`
	Platform string           `json:"platform"`
	TeamID   types.ForeignKey `json:"team_id" validator:"team_id"`
}

func (p ProjectCreateSerializer) Validate(context *context.Context) *validator.Result {
	validator := validator.New()
	validator["name"] = models.ValidateProjectName()
	validator["team_id"] = models.ValidateTeamID(context)
	return validator.Validate(p)
}

/*
ProjectListAPIView
List of projects rest api endpoint
*/
type ProjectListAPIView struct {
	core.JSONView

	context *context.Context
}

func (p *ProjectListAPIView) Before(w http.ResponseWriter, r *http.Request) (err error) {
	p.context = p.Context(r)
	return
}

/*
Retrieve list of projects
*/
func (p *ProjectListAPIView) GET(w http.ResponseWriter, r *http.Request) {
	manager := models.NewProjectManager(p.context)
	paging := manager.NewPagingFromRequest(r)
	projects := manager.NewProjectList()

	if err := manager.FilterPaged(&projects, paging); err != nil {
		glog.Error(err)
		response.New(http.StatusInternalServerError).Write(w, r)
		return
	}

	// update paging
	response.New(http.StatusOK).Result(projects).Paging(paging).Write(w, r)
}

/*
Create new project

Check permissions for user identified by token
Create team and add project to database
*/
func (p *ProjectListAPIView) POST(w http.ResponseWriter, r *http.Request) {

	var err error

	// get auth user to check permissions
	usermanager := models.NewUserManager(p.context)
	user := usermanager.NewUser()
	if err = usermanager.GetAuthUser(user, r); err != nil {
		// this should not happened
		response.New(http.StatusUnauthorized).Write(w, r)
		return
	}

	// unmarshal posted data
	serializer := ProjectCreateSerializer{}

	if err = p.context.Bind(&serializer); err != nil {
		response.New(http.StatusBadRequest).Write(w, r)
	}

	// validate struct
	if vr := serializer.Validate(p.context); !vr.IsValid() {
		response.New(http.StatusBadRequest).Error(vr).Write(w, r)
		return
	}

	// get team from serializer.TeamID
	team := models.NewTeam()
	if err = team.Manager(p.context).GetByID(team, serializer.TeamID); err != nil {
		response.New(http.StatusNotFound).Write(w, r)
		return
	}

	// check permissions (IsSuperuser, member type)
	tmm := models.NewTeamMemberManager(p.context)
	var mt models.MemberType
	// not superuser so we have to check member type for permissions
	if !user.IsSuperuser {
		if mt, err = tmm.MemberType(team, user); err != nil || mt != models.MEMBER_TYPE_ADMIN {
			response.New(http.StatusForbidden).Write(w, r)
			return
		}
	}

	// // create project
	project := models.NewProject(func(proj *models.Project) {
		proj.Name = serializer.Name
		proj.Platform = serializer.Platform
		proj.TeamID = types.ForeignKey(team.ID)
	})

	if vr, _ := project.Validate(p.context); !vr.IsValid() {
		response.New(http.StatusBadRequest).Error(vr).Write(w, r)
		return
	}

	if err = project.Insert(p.context); err != nil {
		response.New(http.StatusInternalServerError).Error(err).Write(w, r)
		return
	}

	// create new project key
	pk := models.NewProjectKey(func(projectKey *models.ProjectKey) {
		projectKey.UserID = types.ForeignKey(user.ID)
		projectKey.UserAddedID = types.ForeignKey(user.ID)
		projectKey.ProjectID = project.ID.ToForeignKey()
	})

	if err = pk.Insert(p.context); err != nil {
		response.New(http.StatusInternalServerError).Error(err).Write(w, r)
		return
	}

	// everything went ok
	response.New(http.StatusCreated).Result(project).Write(w, r)
}

/*
OPTIONS

Metadata about possible methods.
*/
func (p *ProjectListAPIView) OPTIONS(w http.ResponseWriter, r *http.Request) {
	var err error

	// create metadata object
	md := metadata.New("List projects endpoint")

	// get auth user to check permissions
	usermanager := models.NewUserManager(p.context)
	user := usermanager.NewUser()
	if err = usermanager.GetAuthUser(user, r); err != nil {
		// this should not happened
		response.New(http.StatusUnauthorized).Write(w, r)
		return
	}

	// if user has permission to create new project
	if user.IsSuperuser || user.Permissions.Has(settings.PERMISSION_PROJECTS_PROJECT_ADD) {
		create := md.ActionCreate().From(ProjectCreateSerializer{})
		create.Field("name").Update(models.UpdateProjectNameMetadata)
		create.Field("platform").Update(models.UpdateProjectPlatformMetadata)

		tm := models.NewTeamManager(p.context)
		tl := tm.NewTeamList()
		if err = tm.Filter(&tl); err != nil {
			response.New(http.StatusInternalServerError).Write(w, r)
			return
		}

		// add teams choices
		for _, team := range tl {
			create.Field("team_id").Choices.Add(team.ID, team.Name)
		}

	}

	// write metadata response
	response.New(http.StatusOK).Raw(md).Write(w, r)
}
