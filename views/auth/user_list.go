package auth

import (
	"net/http"

	"github.com/phonkee/patrol/context"
	"github.com/phonkee/patrol/core"
	"github.com/phonkee/patrol/models"
	"github.com/phonkee/patrol/rest/metadata"
	"github.com/phonkee/patrol/rest/response"
	"github.com/phonkee/patrol/serializers"
	"github.com/phonkee/patrol/views"
	"github.com/phonkee/patrol/views/mixins"
)

/*
Constructor that returns New UserListAPIView
*/
func NewUserListAPIView() *UserListAPIView {
	return &UserListAPIView{}
}

/*
User list view
	GET - list of users
	POST - create new user
*/
type UserListAPIView struct {
	core.JSONView

	mixins.AuthUserMixin

	// store context
	context *context.Context

	user *models.User
}

// Before every method
func (u *UserListAPIView) Before(w http.ResponseWriter, r *http.Request) (err error) {
	u.context = u.GetContext(r)
	u.user = models.NewUser()

	if err = u.GetAuthUser(u.user, w, r); err != nil {
		return
	}

	if !u.user.IsSuperuser {
		response.New(http.StatusForbidden).Write(w, r)
		return views.ErrForbidden
	}

	return
}

// returns available metadata
func (u *UserListAPIView) OPTIONS(w http.ResponseWriter, r *http.Request) {

	// @TODO: check permissions

	md := metadata.New("User list").SetDescription("User list endpoint")
	md.ActionRetrieve().From(serializers.AuthUserDetailSerializer{})

	response.New(http.StatusOK).Metadata(md).Write(w, r)
}

/*
Retrieve list of registered users
*/
func (u *UserListAPIView) GET(w http.ResponseWriter, r *http.Request) {
	var err error
	manager := models.NewUserManager(u.context)
	paginator := manager.NewPaginator()

	// @TODO: check permissions

	list := []*serializers.AuthUserDetailSerializer{}
	if err = manager.FilterPaged(&list, paginator); err != nil {
		response.New(http.StatusInternalServerError).Write(w, r)
		return
	}

	response.New(http.StatusOK).Result(list).Paging(paginator).Write(w, r)
}

/*
Create new user
*/
func (u *UserListAPIView) POST(w http.ResponseWriter, r *http.Request) {
	var err error

	serializer := &serializers.AuthUserCreateSerializer{}
	if err = u.context.Bind(serializer); err != nil {
		response.New(http.StatusBadRequest).Write(w, r)
		return
	}

	// invalid
	if vr := serializer.Validate(u.context); !vr.IsValid() {
		response.New(http.StatusBadRequest).Error(vr).Write(w, r)
		return
	}

	var result *serializers.AuthUserDetailSerializer

	// save user to database
	if result, err = serializer.Save(u.context); err != nil {
		response.New(http.StatusInternalServerError).Error(err).Write(w, r)
		return
	}

	// clear password - no need to send it through wire
	response.New(http.StatusCreated).Result(result).Write(w, r)
}
