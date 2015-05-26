package projects

import (
	"net/http"

	"github.com/golang/glog"
	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/metadata"
	"github.com/phonkee/patrol/rest/response"
	"github.com/phonkee/patrol/rest/views"
	"github.com/phonkee/patrol/serializers"
	"github.com/phonkee/patrol/settings"
	"github.com/phonkee/patrol/views/mixins"
)

/*
ProjectListAPIView
List of projects rest api endpoint
*/
type ProjectListAPIView struct {
	views.APIView

	// mixins used
	mixins.AuthUserMixin

	context *context.Context

	user *models.User
}

func (p *ProjectListAPIView) Before(w http.ResponseWriter, r *http.Request) (err error) {
	p.context = p.GetContext(r)

	// get authenticated user
	p.user = models.NewUser()
	if err = p.GetAuthUser(p.user, w, r); err != nil {
		return
	}

	return
}

/*
Retrieve list of projects

@TODO: handle superuser and all other users
*/
func (p *ProjectListAPIView) GET(w http.ResponseWriter, r *http.Request) {
	manager := models.NewProjectManager(p.context)
	paginator := manager.NewPaginatorFromRequest(r)
	projects := manager.NewProjectList()

	if err := manager.FilterPaged(&projects, paginator); err != nil {
		glog.Error(err)
		response.New(http.StatusInternalServerError).Write(w, r)
		return
	}

	// update paginator
	response.New(http.StatusOK).Result(projects).Paginator(paginator).Write(w, r)
}

/*
Create new project

Check permissions for user identified by token
Create team and add project to database
*/
func (p *ProjectListAPIView) POST(w http.ResponseWriter, r *http.Request) {

	var err error

	// unmarshal posted data
	serializer := serializers.ProjectsProjectCreateSerializer{}

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
	if !p.user.IsSuperuser {
		if mt, err = tmm.MemberType(team, p.user); err != nil || mt != models.MEMBER_TYPE_ADMIN {
			response.New(http.StatusForbidden).Write(w, r)
			return
		}
	}

	var project *models.Project

	if project, err = serializer.Save(p.context, team, p.user); err != nil {
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
		create := md.ActionCreate().From(serializers.ProjectsProjectCreateSerializer{})
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
