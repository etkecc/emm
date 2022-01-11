package matrix

import (
	"time"

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

// Message struct
type Message struct {
	// ID is a matrix event id of the message
	ID id.EventID
	// Author is a matrix id of the sender
	Author id.UserID
	// Text is the message body in plaintext/markdown format
	Text string
	// HTML is the message body in html format
	HTML string
	// CreatedAt is a timestamp
	CreatedAt time.Time
}

// Messages of the room
// Note on limit - the output slice may be less size than limit you sent in the following cases:
// * room contains less messages than limit
// * some room messages don't contain body/formatted body
func Messages(limit int) ([]*Message, error) {
	var messages []*Message
	filter := &mautrix.FilterPart{
		Types:    []event.Type{event.EventMessage},
		NotTypes: []event.Type{event.EventRedaction, event.EventReaction, event.StateMember},
	}

	chunks, err := client.Messages(room, "", "", 'b', filter, limit)
	if err != nil {
		return nil, err
	}

	for _, chunk := range chunks.Chunk {
		message := parseMessage(chunk)
		if message == nil {
			continue
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func parseMessage(evt *event.Event) *Message {
	err := evt.Content.ParseRaw(event.EventMessage)
	if err != nil {
		return nil
	}

	content := evt.Content.AsMessage()
	if content.Body == "" && content.FormattedBody == "" {
		return nil
	}

	return &Message{
		ID:        evt.ID,
		Author:    evt.Sender,
		Text:      content.Body,
		HTML:      content.FormattedBody,
		CreatedAt: time.UnixMilli(evt.Timestamp),
	}
}
