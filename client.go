// Package pop3 provides simple POP3 client.
package pop3

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/textproto"
)

const (
	// CommandDelete is a command to delete message from POP3 server.
	CommandDelete = "DELE"
	// CommandList is a command to get list of messages from POP3 server.
	CommandList = "LIST"
	// CommandNoop is a ping-like command that tells POP3 to do nothing.
	// (i.e. send something line pong-response).
	CommandNoop = "NOOP"
	// CommandPassword is a command to send user password to POP3 server.
	CommandPassword = "PASS"
	// CommandQuit is a command to tell POP3 server that you are quitting.
	CommandQuit = "QUIT"
	// CommandRetrieve is a command to retrieve POP3 message from server.
	CommandRetrieve = "RETR"
	// CommandUser is a command to send user login to POP3 server.
	CommandUser = "USER"
)

// LineReader is an interface that has a ReadLine method.
type LineReader interface {
	// ReadLine reads a single line from reader.
	ReadLine() (line string, err error)
}

// Client is a POP3 client.
type Client struct {
	conn io.ReadWriteCloser
	err  error

	// Reader is a pointer to a textproto.Reader struct.
	Reader *textproto.Reader
	// Writer is a pointer to a textproto.Writer struct.
	Writer *textproto.Writer
}

// NewClient creates a new POP3 client.
func NewClient(conn io.ReadWriteCloser) (c *Client, err error) {
	c = &Client{
		conn:   conn,
		Reader: textproto.NewReader(bufio.NewReader(conn)),
		Writer: textproto.NewWriter(bufio.NewWriter(conn)),
	}

	if _, err = parseResponseLine(c.Reader); err != nil {
		c, err = nil, newError("new", err)
	}
	return
}

// Dial opens new connection and creates a new POP3 client.
func Dial(addr string) (c *Client, err error) {
	var conn net.Conn
	if conn, err = net.Dial("tcp", addr); err != nil {
		return nil, newError("dial", err)
	}

	return NewClient(conn)
}

// DialTLS opens new TLS connection and creates a new POP3 client.
func DialTLS(addr string) (c *Client, err error) {
	var conn *tls.Conn
	if conn, err = tls.Dial("tcp", addr, nil); err != nil {
		return nil, newError("dial-tls", err)
	}
	return NewClient(conn)
}

// Quit sends "QUIT" command to the POP3 server.
func (c *Client) Quit() (err error) {
	if c.err != ErrAlreadyQuit {
		_, err = c.command(CommandQuit)
	}
	c.err = ErrAlreadyQuit

	return
}

func parseResponseLine(r LineReader) (message string, err error) {
	message, err = r.ReadLine()
	if err != nil {
		return "", fmt.Errorf("read line: %v", err)
	}

	if len(message) < 1 {
		return "", fmt.Errorf("empty message")
	}

	if len(message) < 3 {
		return "", fmt.Errorf("line too short: %q", message)
	}

	if message[:3] == "+OK" {
		if len(message) == 3 {
			message = ""
		} else {
			message = message[4:]
		}

		return
	}

	if len(message) > 3 && message[:4] == "-ERR" {
		if len(message) == 4 {
			return "", fmt.Errorf("error")
		}
		return "", fmt.Errorf("%s", message[5:])
	}

	return "", fmt.Errorf("unexpected response: %s", message)
}

func (c *Client) command(cmd string, args ...interface{}) (resp string,
	err error) {
	if err = c.err; err != nil {
		return
	}

	s := cmd
	if len(args) > 0 {
		s += " " + fmt.Sprintln(args...)
		s = s[:len(s)-1]
	}

	if err = c.Writer.PrintfLine("%s", s); err != nil {
		return "", fmt.Errorf("%s: %v", cmd, err)
	}

	resp, err = parseResponseLine(c.Reader)
	if err != nil {
		err = fmt.Errorf("%s: %v", cmd, err)
	}
	return
}

// Close sends "QUIT" command to POP3 server and closes connection.
func (c *Client) Close() (err error) {
	if c.err == ErrWriteAfterClose {
		return
	}

	err = c.Quit()
	if err != nil {
		_ = c.conn.Close()
		return
	}

	return c.conn.Close()
}

// Login logs into POP3 server with login and password.
func (c *Client) Login(user, pass string) (err error) {
	if _, err = c.command(CommandUser, user); err != nil {
		return
	}

	_, err = c.command(CommandPassword, pass)
	return
}

// GetMessages requests messages lists from POP3 server.
func (c *Client) GetMessages() (list []Message, err error) {
	if _, err = c.command(CommandList); err != nil {
		return
	}

	s := ""
	i := 0
	var n int64

	list = make([]Message, 0, 100)
	for {
		if s, err = c.Reader.ReadLine(); err != nil {
			return nil, newError("messages", err)
		}

		if s == "." {
			return
		}

		if i, n, err = parseListResponse(s); err != nil {
			return nil, newError("messages", err)
		}

		list = append(list, Message{c: c, i: i, n: n})
	}
}

// NoOperation sends ping-like request to the POP3-server.
func (c *Client) NoOperation() (msg string, err error) {
	return c.command(CommandNoop)
}
