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
	// Replace is a matrix ID of old (replaced) event
	Replace id.EventID
	// Author is a matrix id of the sender
	Author id.UserID
	// Text is the message body in plaintext/markdown format
	Text string
	// HTML is the message body in html format
	HTML string
	// CreatedAt is a timestamp, format: 2006-01-02 15:04 UTC
	CreatedAt string
	// CreatedAtFull is a time.Time object
	CreatedAtFull time.Time
}

var msgmap map[id.EventID]*Message

// Messages of the room
// Note on limit - the output slice may be less size than limit you sent in the following cases:
// * room contains less messages than limit
// * some room messages don't contain body/formatted body
func Messages(limit int) ([]*Message, error) {
	msgmap = make(map[id.EventID]*Message)
	filter := &mautrix.FilterPart{
		Types:    []event.Type{event.EventMessage},
		NotTypes: []event.Type{event.EventReaction, event.StateMember},
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
		msgmap[message.ID] = message
	}

	removeReplaced()

	var messages []*Message
	for _, message := range msgmap {
		messages = append(messages, message)
	}

	return messages, nil
}

func parseMessage(evt *event.Event) *Message {
	err := evt.Content.ParseRaw(event.EventMessage)
	if err != nil {
		return nil
	}

	var replace id.EventID
	content := evt.Content.AsMessage()
	text := content.Body
	html := content.FormattedBody
	if content.NewContent != nil {
		text = content.NewContent.Body
		html = content.NewContent.FormattedBody
	}
	if content.RelatesTo != nil {
		replace = content.RelatesTo.GetReplaceID()
	}

	if text == "" && html == "" {
		return nil
	}

	createdAt := time.UnixMilli(evt.Timestamp).UTC()
	return &Message{
		ID:            evt.ID,
		Replace:       replace,
		Author:        evt.Sender,
		Text:          text,
		HTML:          html,
		CreatedAt:     createdAt.Format("2006-01-02 15:04 MST"),
		CreatedAtFull: createdAt,
	}
}

func removeReplaced() {
	var list []id.EventID
	for _, message := range msgmap {
		if message.Replace != "" {
			list = append(list, message.Replace)
		}
	}
	for _, eventID := range list {
		delete(msgmap, eventID)
	}
}
