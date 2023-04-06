package matrix

import (
	"log"
	"time"

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

// Page is a amount of messages per page
const Page = 100

// Message struct
type Message struct {
	// ID is a matrix event id of the message
	ID id.EventID
	// Replace is a matrix ID of old (replaced) event
	Replace id.EventID
	// ReplacedNote is a text note usable from template to mark replaced message as updated
	ReplacedNote string
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

var (
	msgmap map[id.EventID]*Message
	filter = &mautrix.FilterPart{
		Types:    []event.Type{event.EventMessage},
		NotTypes: []event.Type{event.EventReaction, event.StateMember},
	}
)

// Messages of the room
// Note on limit - the output slice may be less size than limit you sent in the following cases:
// * room contains less messages than limit
// * some room messages don't contain body/formatted body
func Messages(limit int) ([]*Message, error) {
	var err error
	msgmap = make(map[id.EventID]*Message, 0)
	if limit > Page {
		err = paginate(limit)
	} else {
		err = load()
	}
	if err != nil {
		return nil, err
	}

	var messages []*Message
	for _, message := range msgmap {
		messages = append(messages, message)
	}
	log.Println("loaded", len(messages), "messages total")

	return messages, nil
}

func paginate(limit int) error {
	var token string
	page := 1
	for i := Page; i < limit; {
		var chunks *mautrix.RespMessages
		log.Println("requesting messages from", room, "page =", page)
		err := retry(func() error {
			var messagesErr error
			chunks, messagesErr = client.Messages(room, token, "", 'b', filter, Page)

			return messagesErr
		})
		if err != nil {
			return err
		}
		if len(chunks.Chunk) == 0 {
			log.Println("no more messages")
			break
		}

		processEvents(chunks)
		token = chunks.End
		if len(chunks.Chunk) < Page {
			log.Println("it was the last page")
			break
		}

		i += Page
		page++
	}

	return nil
}

func load() error {
	var chunks *mautrix.RespMessages
	log.Println("requesting messages from", room, "without pagination")
	err := retry(func() error {
		var messagesErr error
		chunks, messagesErr = client.Messages(room, "", "", 'b', filter, Page)

		return messagesErr
	})
	if err != nil {
		return err
	}
	processEvents(chunks)
	return nil
}

func processEvents(resp *mautrix.RespMessages) {
	log.Println("parsing messages chunk:", len(resp.Chunk), "events")
	for _, evt := range resp.Chunk {
		_, ignore := ignored[evt.Sender]
		if ignore {
			continue
		}

		message := parseMessage(evt)
		if message == nil {
			continue
		}
		addMessage(message)
	}
}

func addMessage(message *Message) {
	if message.Replace != "" {
		message.ID = message.Replace
		message.ReplacedNote = " (updated)"
	}

	msg, ok := msgmap[message.ID]
	if !ok {
		msgmap[message.ID] = message
		return
	}

	if msg.CreatedAtFull.Before(message.CreatedAtFull) {
		msgmap[message.ID] = message
	}
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
