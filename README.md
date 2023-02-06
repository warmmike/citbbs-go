# citbbs-go

Go package to access the citbbs API.


## Install

```bash
go get github.com/warmmike/citbbs-go/citbbs
```

## Usage

Here is an example usage of the citbbs Go client. Please make sure to
handle errors in your production application.

```go
c, _ := citbbs.NewClient(citbbs.WithAccessToken(string(creds)))
// Get user resource observability
user, _ := c.Users.Get(ctx, &citbbs.GetUserRequest{
    User: user,
})
```
