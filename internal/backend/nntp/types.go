package nntp

import (
	"io"
	"net"
	"net/textproto"
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

// NNTPClient is an NNTP client.
type NNTPClient struct {
	conn         *textproto.Conn
	netconn      net.Conn
	tls          bool
	Banner       string
	capabilities []string
}

type Backend struct {
	Name   string `json:"name"`
	User   string `json:"user,omitempty"`
	Pass   string `json:"pass,omitempty"`
	Server string `json:"server"`
	Port   int    `json:"port,omitempty"`

	client *NNTPClient
}
