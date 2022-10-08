package nntp

import (
	"io"
	"net/textproto"
	"newsmere/internal/storage"
)

const Type = "nntp"

// Group represents a usenet newsgroup.
type Group struct {
	Name        string
	Description string
	Count       int64
	High        int64
	Low         int64
}

// An Article that may appear in one or more groups.
type Article struct {
	// The article's headers
	Header textproto.MIMEHeader
	Body   io.Reader
	Bytes  int
	Lines  int
}

// NNTPError is coded NNTP error message.
type NNTPError struct {
	Code int
	Msg  string
}

// ErrNoSuchGroup is returned for a request for a group that can't be found.
var ErrNoSuchGroup = &NNTPError{411, "No such newsgroup"}

// ErrNoSuchGroup is returned for a request that requires a current
// group when none has been selected.
var ErrNoGroupSelected = &NNTPError{412, "No newsgroup selected"}

// ErrInvalidMessageID is returned when a message is requested that can't be found.
var ErrInvalidMessageID = &NNTPError{430, "No article with that message-id"}

// ErrInvalidArticleNumber is returned when an article is requested that can't be found.
var ErrInvalidArticleNumber = &NNTPError{423, "No article with that number"}

// ErrNoCurrentArticle is returned when a command is executed that
// requires a current article when one has not been selected.
var ErrNoCurrentArticle = &NNTPError{420, "Current article number is invalid"}

// ErrUnknownCommand is returned for unknown comands.
var ErrUnknownCommand = &NNTPError{500, "Unknown command"}

// ErrSyntax is returned when a command can't be parsed.
var ErrSyntax = &NNTPError{501, "not supported, or syntax error"}

// ErrPostingNotPermitted is returned as the response to an attempt to
// post an article where posting is not permitted.
var ErrPostingNotPermitted = &NNTPError{440, "Posting not permitted"}

// ErrPostingFailed is returned when an attempt to post an article fails.
var ErrPostingFailed = &NNTPError{441, "posting failed"}

// ErrNotWanted is returned when an attempt to post an article is
// rejected due the server not wanting the article.
var ErrNotWanted = &NNTPError{435, "Article not wanted"}

// ErrAuthRequired is returned to indicate authentication is required
// to proceed.
var ErrAuthRequired = &NNTPError{450, "authorization required"}

// ErrAuthRejected is returned for invalid authentication.
var ErrAuthRejected = &NNTPError{452, "authorization rejected"}

// ErrNotAuthenticated is returned when a command is issued that requires
// authentication, but authentication was not provided.
var ErrNotAuthenticated = &NNTPError{480, "authentication required"}

// Handler is a low-level protocol handler
type Handler func(args []string, s *session, c *textproto.Conn) error

// A NumberedArticle provides local sequence nubers to articles When
// listing articles in a group.
type NumberedArticle struct {
	Num     int64
	Article *Article
}

// The Operator that provides the things and does the stuff.
type Operator interface {
	ListGroups(max int) ([]*storage.Group, error)
	GetGroup(name string) (*storage.Group, error)
	GetArticle(group *Group, id string) (*Article, error)
	GetArticles(group *Group, from, to int64) ([]NumberedArticle, error)
	Authorized() bool
	Authenticate(user, pass string) (Operator, error)
}

type session struct {
	server   *NNTPServer
	operator Operator
	group    *Group
}

// The Server handle.
type NNTPServer struct {
	Handlers map[string]Handler
	Operator Operator
	group    *Group
}

type Service struct {
	Host string `json:"Host"`
	Port int    `json:"Port"`

	operator Operator
	server   *NNTPServer
}
