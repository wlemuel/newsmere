package domain

import (
	"io"
	"net/textproto"
)

// An Article that may appear in one or more groups.
type NNTPArticle struct {
	// The article's headers
	Header textproto.MIMEHeader
	Body   io.Reader
	Bytes  int
	Lines  int
}

// nntpGroup represents a usenet newsgroup.
type NNTPGroup struct {
	Name        string
	Description string
	Count       int64
	High        int64
	Low         int64
}
