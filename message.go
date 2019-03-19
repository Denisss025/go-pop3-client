package pop3

import (
	"fmt"
	"net/mail"
)

// Message is a struct that helps working with POP3 messages.
type Message struct {
	c *Client

	i int
	n int64
}

func parseListResponse(list string) (num int, size int64, err error) {
	_, err = fmt.Sscanf(list, "%d %d", &num, &size)

	if err != nil {
		err = fmt.Errorf("parse %q: %v", list, err)
	}
	return
}

// Index returns an index number of the message.
func (m Message) Index() (i int) { return m.i }

// GetSize returns a size of the message.
func (m *Message) GetSize() (n int64, err error) {
	if m.n > 0 {
		n = m.n
		return
	}

	line := ""
	if line, err = m.c.command(CommandList, m.i); err == nil {
		_, n, err = parseListResponse(line)
	}

	m.n = n
	if err != nil {
		m.n = 0
		err = newError("list", err)
	}
	return
}

// Retrieve retrieves the message body from the POP3 server.
func (m *Message) Retrieve() (line string, msg *mail.Message, err error) {
	if line, err = m.c.command(CommandRetrieve, m.i); err != nil {
		return
	}

	if msg, err = mail.ReadMessage(m.c.Reader.DotReader()); err != nil {
		err = newError("retrieve", err)
	}

	return
}

// Delete removes the message from the POP3 server.
func (m *Message) Delete() (msg string, err error) {
	return m.c.command(CommandDelete, m.i)
}
