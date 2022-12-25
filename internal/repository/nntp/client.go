package nntp

import (
	"errors"
	"io"
	"net"
	"net/textproto"
	"newsmere/internal/domain"
	"strconv"
	"strings"
)

// newClient for create a nntp client to connect to server
func newClient(network, addr string) (*nntpClient, error) {
	netconn, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	return connect(netconn)
}

// newConnClient wraps an existing connection, for example one opened with tls.Dial
// func newConnClient(netconn net.Conn) (*nntpClient, error) {
// 	client, err := connect(netconn)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if _, ok := netconn.(*tls.Conn); ok {
// 		client.tls = true
// 	}
// 	return client, nil
// }

func connect(netconn net.Conn) (*nntpClient, error) {
	conn := textproto.NewConn(netconn)
	_, msg, err := conn.ReadCodeLine(200)
	if err != nil {
		return nil, err
	}

	return &nntpClient{
		conn:    conn,
		netconn: netconn,
		Banner:  msg,
	}, nil
}

// CLose the nntp client
func (c *nntpClient) Close() error {
	return c.conn.Close()
}

// Authenticate against an NNTP server using authinfo user/pass.
func (c *nntpClient) Authenticate(user, pass string) (msg string, err error) {
	err = c.conn.PrintfLine("authinfo user %s", user)
	if err != nil {
		return
	}

	_, _, err = c.conn.ReadCodeLine(381)
	if err != nil {
		return
	}

	err = c.conn.PrintfLine("authinfo pass %s", pass)
	if err != nil {
		return
	}

	_, msg, err = c.conn.ReadCodeLine(281)
	return
}

// Command sends a low-level command and get a response.
func (c *nntpClient) Command(cmd string, expectCode int) (int, string, error) {
	err := c.conn.PrintfLine(cmd)
	if err != nil {
		return 0, "", err
	}
	return c.conn.ReadCodeLine(expectCode)
}

// asLines issues a command and returns the response's data block as lines.
func (c *nntpClient) asLines(cmd string, expectCode int) ([]string, error) {
	_, _, err := c.Command(cmd, expectCode)
	if err != nil {
		return nil, err
	}
	return c.conn.ReadDotLines()
}

// capabilities retrieves a list of supported capabilities.
func (c *nntpClient) Capabilities() ([]string, error) {
	caps, err := c.asLines("CAPABILITIES", 101)
	if err != nil {
		return nil, err
	}
	for i, line := range caps {
		caps[i] = strings.ToUpper(line)
	}
	c.capabilities = caps
	return caps, nil
}

// GetCapability returns a complete capability line.
func (c *nntpClient) GetCapability(capability string) string {
	capability = strings.ToUpper(capability)
	for _, capa := range c.capabilities {
		i := strings.IndexAny(capa, "\t ")
		if i != -1 && capa[:i] == capability {
			return capa
		}
		if capa == capability {
			return capa
		}
	}
	return ""
}

// HasCapabilityArgument indicates whether a capability arg is supported
func (c *nntpClient) HasCapabilityArgument(capability, argument string) (bool, error) {
	if c.capabilities == nil {
		return false, errors.New("Capabilities unpopulated")
	}
	capLine := c.GetCapability(capability)
	if capLine == "" {
		return false, errors.New("no such capability")
	}
	argument = strings.ToUpper(argument)
	for _, capArg := range strings.Fields(capLine)[1:] {
		if capArg == argument {
			return true, nil
		}
	}
	return false, nil
}

// List groups
func (c *nntpClient) List(sub string) (rv []domain.NNTPGroup, err error) {
	cmd := "LIST"
	if sub != "" {
		cmd = "LIST " + sub
	}
	_, _, err = c.Command(cmd, 215)
	if err != nil {
		return
	}

	var groupLines []string
	groupLines, err = c.conn.ReadDotLines()
	if err != nil {
		return
	}

	rv = make([]domain.NNTPGroup, 0, len(groupLines))
	for _, l := range groupLines {
		parts := strings.Split(l, " ")
		high, errh := strconv.ParseInt(parts[1], 10, 64)
		low, errl := strconv.ParseInt(parts[2], 10, 64)
		if errh == nil && errl == nil {
			rv = append(rv, domain.NNTPGroup{
				Name: parts[0],
				High: high,
				Low:  low,
			})
		}
	}
	return
}

// ListOverviewFmt performs a LIST OVERVIEW.FMT query.
func (c *nntpClient) ListOverviewFmt() ([]string, error) {
	fields, err := c.asLines("LIST OVERVIEW.FMT", 215)
	if err != nil {
		return nil, err
	}
	return fields, nil
}

// Over returns a list of raw overview lines with tab-separated fields.
func (c *nntpClient) Over(specifier string) ([]string, error) {
	lines, err := c.asLines("OVER "+specifier, 224)
	if err != nil {
		return nil, err
	}
	return lines, nil
}

// Group select a group.
func (c *nntpClient) Group(name string) (rv domain.NNTPGroup, err error) {
	var msg string
	_, msg, err = c.Command("GROUP "+name, 211)
	if err != nil {
		return
	}

	parts := strings.Split(msg, " ")
	if len(parts) != 4 {
		err = errors.New("Don't know how to parse result: " + msg)
		return
	}

	rv.Count, err = strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return
	}

	rv.Low, err = strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return
	}

	rv.High, err = strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return
	}

	rv.Name = parts[3]

	return
}

func (c *nntpClient) articleish(expected int) (int64, string, io.Reader, error) {
	_, msg, err := c.conn.ReadCodeLine(expected)
	if err != nil {
		return 0, "", nil, err
	}
	parts := strings.SplitN(msg, " ", 2)
	n, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, "", nil, err
	}
	return n, parts[1], c.conn.DotReader(), nil
}

// Article grabs an article
func (c *nntpClient) Article(specifier string) (int64, string, io.Reader, error) {
	err := c.conn.PrintfLine("ARTICLE %s", specifier)
	if err != nil {
		return 0, "", nil, err
	}
	return c.articleish(220)
}

// Head gets the headers of an article
func (c *nntpClient) Head(specifier string) (int64, string, io.Reader, error) {
	err := c.conn.PrintfLine("HEAD %s", specifier)
	if err != nil {
		return 0, "", nil, err
	}
	return c.articleish(221)
}

// Body gets the body of an article
func (c *nntpClient) Body(specifier string) (int64, string, io.Reader, error) {
	err := c.conn.PrintfLine("BODY %s", specifier)
	if err != nil {
		return 0, "", nil, err
	}
	return c.articleish(222)
}

// HasTLS checks whether tls supported.
func (c *nntpClient) HasTLS() bool {
	return c.tls
}
