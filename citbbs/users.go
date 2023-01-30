package citbbs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// CreateUserRequest encapsulates the request for creating a new user.
type CreateUserRequest struct {
	Name string `json:"name"`
}

// UserRequest encapsulates the request for getting a single user.
type GetUserRequest struct {
	User string
}

// ListUsersRequest encapsulates the request for listing all users in an
// organization.
type ListUsersRequest struct {
	Organization string
}

// DeleteUserRequest encapsulates the request for deleting a user from
// an organization.
type DeleteUserRequest struct {
	Organization string
	User         string
}

// UserService is an interface for communicating with the PlanetScale
// Users API endpoint.
type UsersService interface {
	Create(context.Context, *CreateUserRequest) (*User, error)
	Get(context.Context, *GetUserRequest) (*User, error)
	List(context.Context, *ListUsersRequest, ...ListOption) ([]*User, error)
	Delete(context.Context, *DeleteUserRequest) (*UserDeletionRequest, error)
}

// UserDeletionRequest encapsulates the request for deleting a user from
// an organization.
type UserDeletionRequest struct {
	User string `json:"name"`
}

// UserState represents the state of a user
type UserState string

const (
	UserPending         UserState = "pending"
	UserImporting       UserState = "importing"
	UserAwakening       UserState = "awakening"
	UserSleepInProgress UserState = "sleep_in_progress"
	UserSleeping        UserState = "sleeping"
	UserReady           UserState = "ready"
)

// User represents a PlanetScale user
//
//	type User struct {
//		Name      string    `json:"name"`
//		Notes     string    `json:"notes"`
//		State     UserState `json:"state"`
//		HtmlURL   string    `json:"html_url"`
//		CreatedAt time.Time `json:"created_at"`
//		UpdatedAt time.Time `json:"updated_at"`
//	}
type User struct {
	Name string `json:"name"`
}

// User represents a list of PlanetScale users
type usersResponse struct {
	Users []*User `json:"data"`
}

type usersService struct {
	client *Client
}

var _ UsersService = &usersService{}

func NewUsersService(client *Client) *usersService {
	return &usersService{
		client: client,
	}
}

func (ds *usersService) List(ctx context.Context, listReq *ListUsersRequest, opts ...ListOption) ([]*User, error) {
	path := usersAPIPath(listReq.Organization)

	defaultOpts := defaultListOptions(WithPerPage(100))
	for _, opt := range opts {
		opt(defaultOpts)
	}

	if vals := defaultOpts.URLValues.Encode(); vals != "" {
		path += "?" + vals
	}

	req, err := ds.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating http request")
	}

	dbResponse := usersResponse{}
	err = ds.client.do(ctx, req, &dbResponse)
	if err != nil {
		return nil, err
	}

	return dbResponse.Users, nil
}

func (ds *usersService) Create(ctx context.Context, createReq *CreateUserRequest) (*User, error) {
	path := fmt.Sprintf("users/%s", createReq.Name)
	req, err := ds.client.newRequest(http.MethodPost, path, createReq)
	if err != nil {
		return nil, errors.Wrap(err, "error creating request for create user")
	}

	user := &User{}
	err = ds.client.do(ctx, req, &user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ds *usersService) Get(ctx context.Context, getReq *GetUserRequest) (*User, error) {
	path := fmt.Sprintf("users/%s", getReq.User)
	//fmt.Println(path)
	req, err := ds.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating request for get user")
	}
	//fmt.Println(req, err)
	user := &User{}
	err = ds.client.do(ctx, req, &user)
	//fmt.Println(err)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ds *usersService) Delete(ctx context.Context, deleteReq *DeleteUserRequest) (*UserDeletionRequest, error) {
	path := fmt.Sprintf("users/%s", deleteReq.User)
	req, err := ds.client.newRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating request for delete user")
	}

	var udr *UserDeletionRequest
	err = ds.client.do(ctx, req, &udr)
	if err != nil {
		return nil, err
	}

	return udr, nil
}

func usersAPIPath(org string) string {
	return fmt.Sprintf("v1/organizations/%s/users", org)
}
