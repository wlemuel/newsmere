package nntp

import (
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"net/textproto"
	"strconv"
	"strings"
)

func NewServer(operator Operator) *NNTPServer {
	rv := NNTPServer{
		Handlers: make(map[string]Handler),
		Operator: operator,
	}
	rv.Handlers[""] = handleDefault
	rv.Handlers["quit"] = handleQuit
	rv.Handlers["group"] = handleGroup
	rv.Handlers["list"] = handleList
	rv.Handlers["head"] = handleHead
	rv.Handlers["body"] = handleBody
	rv.Handlers["article"] = handleArticle
	rv.Handlers["capabilities"] = handleCap
	rv.Handlers["mode"] = handleMode
	rv.Handlers["authinfo"] = handleAuthInfo
	rv.Handlers["newgroups"] = handleNewGroups
	rv.Handlers["over"] = handleOver
	rv.Handlers["xover"] = handleOver

	return &rv
}

func (e *NNTPError) Error() string {
	return fmt.Sprintf("%d %s", e.Code, e.Msg)
}

func (s *session) dispatchCommand(cmd string, args []string, c *textproto.Conn) (err error) {
	handler, found := s.server.Handlers[strings.ToLower(cmd)]
	if !found {
		handler, found = s.server.Handlers[""]
		if !found {
			panic("No default handler.")
		}
	}
	return handler(args, s, c)
}

// Process an NNTP session.
func (s *NNTPServer) Process(nc net.Conn) {
	defer nc.Close()
	c := textproto.NewConn(nc)

	sess := &session{
		server:   s,
		operator: s.Operator,
		group:    nil,
	}

	c.PrintfLine("200 Hello!")
	for {
		l, err := c.ReadLine()
		if err != nil {
			log.Printf("Error reading from client, dropping conn: %v", err)
			return
		}
		cmd := strings.Split(l, " ")
		log.Printf("Got cmd: %+v", cmd)
		args := []string{}
		if len(cmd) > 1 {
			args = cmd[1:]
		}
		err = sess.dispatchCommand(cmd[0], args, c)
		if err != nil {
			_, isNNTPError := err.(*NNTPError)
			switch {
			case err == io.EOF:
				// Drop this connection silently. They hung up
				return
			case isNNTPError:
				c.PrintfLine(err.Error())
			default:
				log.Printf("Error dispatching command, dropping conn: %v", err)
				return
			}
		}
	}
}

func parseRange(spec string) (low, high int64) {
	if spec == "" {
		return 0, math.MaxInt64
	}
	parts := strings.Split(spec, "-")
	if len(parts) == 1 {
		h, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			h = math.MaxInt64
		}
		return 0, h
	}
	l, _ := strconv.ParseInt(parts[0], 10, 64)
	h, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		h = math.MaxInt64
	}
	return l, h
}

func handleOver(args []string, s *session, c *textproto.Conn) error {
	if s.group == nil {
		return ErrNoGroupSelected
	}
	from, to := parseRange(args[0])
	articles, err := s.operator.GetArticles(s.group, from, to)
	if err != nil {
		return err
	}

	c.PrintfLine("224 here it comes")
	dw := c.DotWriter()
	defer dw.Close()
	for _, a := range articles {
		fmt.Fprintf(dw, "%d\t%s\t%s\t%s\t%s\t%s\t%d\t%d\n", a.Num,
			a.Article.Header.Get("Subject"),
			a.Article.Header.Get("From"),
			a.Article.Header.Get("Date"),
			a.Article.Header.Get("Message-Id"),
			a.Article.Header.Get("References"),
			a.Article.Bytes, a.Article.Lines)
	}
	return nil
}

func handleListOverviewFmt(c *textproto.Conn) error {
	err := c.PrintfLine("215 Order of fields in overview database.")
	if err != nil {
		return err
	}
	dw := c.DotWriter()
	defer dw.Close()
	_, err = fmt.Fprintln(dw, `Subject:
From:
Date:
Message-ID:
References:
:bytes
:lines`)
	return err
}

func handleList(args []string, s *session, c *textproto.Conn) error {
	ltype := "active"
	if len(args) > 0 {
		ltype = strings.ToLower(args[0])
	}

	if ltype == "overview.fmt" {
		return handleListOverviewFmt(c)
	}

	groups, err := s.operator.ListGroups(-1)
	if err != nil {
		return err
	}
	c.PrintfLine("215 list of newsgroups follows")
	dw := c.DotWriter()
	defer dw.Close()
	for _, g := range groups {
		switch ltype {
		case "active":
			fmt.Fprintf(dw, "%s.%s %d %d %v\r\n",
				g.Source, g.Name, g.High, g.Low, g.High-g.Low)
		case "newsgroups":
			fmt.Fprintf(dw, "%s %s\r\n", g.Name, g.Description)
		}
	}
	return nil
}

func handleNewGroups(args []string, s *session, c *textproto.Conn) error {
	c.PrintfLine("231 list of newsgroups follows")
	c.PrintfLine(".")
	return nil
}

func handleDefault(args []string, s *session, c *textproto.Conn) error {
	return ErrUnknownCommand
}

func handleQuit(args []string, s *session, c *textproto.Conn) error {
	c.PrintfLine("205 bye")
	return io.EOF
}

func handleGroup(args []string, s *session, c *textproto.Conn) error {
	if len(args) < 1 {
		return ErrNoSuchGroup
	}

	group, err := s.operator.GetGroup(args[0])
	if err != nil {
		return err
	}

	s.group = &Group{
		Name:        group.Name,
		Description: group.Description,
		High:        int64(group.High),
		Low:         int64(group.Low),
		Count:       int64(group.High) - int64(group.Low),
	}

	c.PrintfLine("211 %d %d %d %s", s.group.Count, group.Low, group.High, group.Name)
	return nil
}

func (s *session) getArticle(args []string) (*Article, error) {
	if s.group == nil {
		return nil, ErrNoGroupSelected
	}
	return s.operator.GetArticle(s.group, args[0])
}

func handleHead(args []string, s *session, c *textproto.Conn) error {
	article, err := s.getArticle(args)
	if err != nil {
		return err
	}
	c.PrintfLine("221 1 %s", article.MessageID())
	dw := c.DotWriter()
	defer dw.Close()
	for k, v := range article.Header {
		fmt.Fprintf(dw, "%s: %s\r\n", k, v[0])
	}
	return nil
}

func handleBody(args []string, s *session, c *textproto.Conn) error {
	article, err := s.getArticle(args)
	if err != nil {
		return err
	}
	c.PrintfLine("222 1 %s", article.MessageID())
	dw := c.DotWriter()
	defer dw.Close()
	_, err = io.Copy(dw, article.Body)
	return err
}

func handleArticle(args []string, s *session, c *textproto.Conn) error {
	article, err := s.getArticle(args)
	if err != nil {
		return err
	}
	c.PrintfLine("220 1 %s", article.MessageID())
	dw := c.DotWriter()
	defer dw.Close()

	for k, v := range article.Header {
		fmt.Fprintf(dw, "%s: %s\r\n", k, v[0])
	}

	fmt.Fprintln(dw, "")

	_, err = io.Copy(dw, article.Body)
	return err
}

func handleCap(args []string, s *session, c *textproto.Conn) error {
	c.PrintfLine("101 Capability list:")
	dw := c.DotWriter()
	defer dw.Close()

	fmt.Fprintf(dw, "VERSION 2\n")
	fmt.Fprintf(dw, "READER\n")
	fmt.Fprintf(dw, "OVER\n")
	fmt.Fprintf(dw, "XOVER\n")
	fmt.Fprintf(dw, "LIST ACTIVE NEWSGROUPS OVERVIEW.FMT\n")
	return nil
}

func handleMode(args []string, s *session, c *textproto.Conn) error {
	c.PrintfLine("201 Posting prohibited")
	return nil
}

func handleAuthInfo(args []string, s *session, c *textproto.Conn) error {
	if len(args) < 2 {
		return ErrSyntax
	}
	if strings.ToLower(args[0]) != "user" {
		return ErrSyntax
	}
	if s.operator.Authorized() {
		return c.PrintfLine("250 authenticated")
	}

	c.PrintfLine("350 Continue")
	a, err := c.ReadLine()
	if err != nil {
		return err
	}

	parts := strings.SplitN(a, " ", 3)
	if strings.ToLower(parts[0]) != "authinfo" || strings.ToLower(parts[1]) != "pass" {
		return ErrSyntax
	}

	b, err := s.operator.Authenticate(args[1], parts[2])
	if err != nil {
		c.PrintfLine("250 authenticated")
		if b != nil {
			s.operator = b
		}
	}
	return err
}
