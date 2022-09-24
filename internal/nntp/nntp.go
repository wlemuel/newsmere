package nntp

import (
	"io"
	"log"
	"net/textproto"
	"newsmere/internal/types"
)

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

// MessageID provides convenient access to the article's Message ID.
func (a *Article) MessageID() string {
	return a.Header.Get("Message-Id")
}

func (c *NNTPClient) GetType() types.BackendType {
	return types.NNTP
}

func (c *NNTPClient) GetName() string {
	return "NNTP"
}

func (c *NNTPClient) IsRunning() bool {
	return c.running
}

func (c *NNTPClient) Run() error {
	c.running = true
	return nil
}

func (c *NNTPClient) Stop() error {
	c.Close()
	log.Println("NNTPClient stopped")
	c.running = false
	return nil
}

func (c *NNTPClient) Restart() error {
	c.running = true
	return nil
}
