package citbbs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// CreateMessageRequest encapsulates the request for creating a new message.
type CreateMessageRequest struct {
	Name string `json:"name"`
}

// MessageRequest encapsulates the request for getting a single message.
type GetMessageRequest struct {
	User string
}

// ListMessagesRequest encapsulates the request for listing all messages in an
// organization.
type ListMessagesRequest struct {
	Organization string
}

// DeleteMessageRequest encapsulates the request for deleting a message from
// an organization.
type DeleteMessageRequest struct {
	Organization string
	User         string
}

// MessageService is an interface for communicating with the PlanetScale
// Messages API endpoint.
type MessagesService interface {
	Create(context.Context, *CreateMessageRequest) (*User, error)
	Get(context.Context, *GetMessageRequest) (*User, error)
	List(context.Context, *ListMessagesRequest, ...ListOption) ([]*User, error)
	Delete(context.Context, *DeleteMessageRequest) (*MessageDeletionRequest, error)
}

// MessageDeletionRequest encapsulates the request for deleting a message from
// an organization.
type MessageDeletionRequest struct {
	User string `json:"name"`
}

// MessageState represents the state of a message
type MessageState string

const (
	MessagePending         MessageState = "pending"
	MessageImporting       MessageState = "importing"
	MessageAwakening       MessageState = "awakening"
	MessageSleepInProgress MessageState = "sleep_in_progress"
	MessageSleeping        MessageState = "sleeping"
	MessageReady           MessageState = "ready"
)

// Message represents a citbbs message
//
//	type Message struct {
//		Name      string    `json:"name"`
//		Notes     string    `json:"notes"`
//		State     MessageState `json:"state"`
//		HtmlURL   string    `json:"html_url"`
//		CreatedAt time.Time `json:"created_at"`
//		UpdatedAt time.Time `json:"updated_at"`
//	}
//type Message struct {
//	Name string `json:"name"`
//}

// Message represents a list of citbbs messages
type messagesResponse struct {
	Users []*User `json:"data"`
}

type messagesService struct {
	client *Client
}

var _ MessagesService = &messagesService{}

func NewMessagesService(client *Client) *messagesService {
	return &messagesService{
		client: client,
	}
}

func (ds *messagesService) List(ctx context.Context, listReq *ListMessagesRequest, opts ...ListOption) ([]*User, error) {
	path := messagesAPIPath(listReq.Organization)

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

	dbResponse := messagesResponse{}
	err = ds.client.do(ctx, req, &dbResponse)
	if err != nil {
		return nil, err
	}

	return dbResponse.Users, nil
}

func (ds *messagesService) Create(ctx context.Context, createReq *CreateMessageRequest) (*User, error) {
	path := fmt.Sprintf("messages/%s", createReq.Name)
	req, err := ds.client.newRequest(http.MethodPost, path, createReq)
	if err != nil {
		return nil, errors.Wrap(err, "error creating request for create message")
	}

	user := &User{}
	err = ds.client.do(ctx, req, &user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ds *messagesService) Get(ctx context.Context, getReq *GetMessageRequest) (*User, error) {
	path := fmt.Sprintf("messages/%s", getReq.User)
	//fmt.Println(path)
	req, err := ds.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating request for get message")
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

func (ds *messagesService) Delete(ctx context.Context, deleteReq *DeleteMessageRequest) (*MessageDeletionRequest, error) {
	path := fmt.Sprintf("messages/%s", deleteReq.User)
	req, err := ds.client.newRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating request for delete user")
	}

	var udr *MessageDeletionRequest
	err = ds.client.do(ctx, req, &udr)
	if err != nil {
		return nil, err
	}

	return udr, nil
}

func messagesAPIPath(org string) string {
	return fmt.Sprintf("v1/organizations/%s/messages", org)
}
