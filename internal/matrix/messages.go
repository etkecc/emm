package matrix

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/etkecc/emm/internal/utils"
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
	// URLSafeID is a url-safe version of the ID
	URLSafeID string
	// Replace is a matrix ID of old (replaced) event
	Replace id.EventID
	// ReplacedNote is a text note usable from template to mark replaced message as updated
	ReplacedNote string
	// Author is a matrix id of the sender
	Author id.UserID
	// Title is a message title (first line of the message)
	Title string
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

func (m *Message) Vars() map[string]string {
	return map[string]string{
		"ID":            string(m.ID),
		"URLSafeID":     m.URLSafeID,
		"Replace":       string(m.Replace),
		"ReplacedNote":  m.ReplacedNote,
		"Author":        string(m.Author),
		"Title":         m.Title,
		"Text":          m.Text,
		"HTML":          m.HTML,
		"CreatedAtDate": m.CreatedAtFull.Format(time.DateOnly),
		"CreatedAt":     m.CreatedAt,
		"CreatedAtFull": m.CreatedAtFull.String(),
	}
}

// Messages of the room
// Note on limit - the output slice may be less size than limit you sent in the following cases:
// * room contains less messages than limit
// * some room messages don't contain body/formatted body
func Messages(limit int, since time.Time) (map[id.EventID]*Message, error) {
	var err error
	var sinceMS int64
	if !since.IsZero() {
		sinceMS = since.UnixMilli()
	}
	msgmap = make(map[id.EventID]*Message, 0)
	if limit > Page {
		err = paginate(limit, sinceMS)
	} else {
		err = load(sinceMS)
	}
	if err != nil {
		return nil, err
	}

	log.Println("loaded", len(msgmap), "messages total")

	return msgmap, nil
}

func paginate(limit int, sinceMS int64) error {
	ctx := context.Background()
	var token string
	page := 1
	for i := Page; i < limit; {
		var chunks *mautrix.RespMessages
		log.Println("requesting messages from", room, "page =", page)
		err := retry(func() error {
			var messagesErr error
			chunks, messagesErr = client.Messages(ctx, room, token, "", 'b', filter, Page)

			return messagesErr
		})
		if err != nil {
			return err
		}
		if len(chunks.Chunk) == 0 {
			log.Println("no more messages")
			break
		}

		processEvents(chunks, sinceMS)
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

func load(sinceMS int64) error {
	ctx := context.Background()
	var chunks *mautrix.RespMessages
	log.Println("requesting messages from", room, "without pagination")
	err := retry(func() error {
		var messagesErr error
		chunks, messagesErr = client.Messages(ctx, room, "", "", 'b', filter, Page)

		return messagesErr
	})
	if err != nil {
		return err
	}
	processEvents(chunks, sinceMS)
	return nil
}

func processEvents(resp *mautrix.RespMessages, sinceMS int64) {
	log.Println("parsing messages chunk:", len(resp.Chunk), "events")
	for _, evt := range resp.Chunk {
		_, ignore := ignored[evt.Sender]
		if ignore {
			continue
		}

		if evt.Timestamp < sinceMS {
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
	title := createdAt.Format("2006-01-02 15:04 MST")
	if strings.Contains(text, "\n") {
		title = sanitizeTitle(strings.Split(text, "\n")[0])
	}

	return &Message{
		ID:            evt.ID,
		URLSafeID:     utils.MakeURLSafe(string(evt.ID)),
		Replace:       replace,
		Author:        evt.Sender,
		Title:         title,
		Text:          text,
		HTML:          html,
		CreatedAt:     createdAt.Format("2006-01-02 15:04 MST"),
		CreatedAtFull: createdAt,
	}
}

func sanitizeTitle(text string) string {
	prefixes := []string{"> ", "# ", "## ", "### ", "#### ", "```", "~~", "| ", "* ", "- ", "+ ", "1. ", "1) "}
	kws := []string{"**", "__", "`", "```", "~~", "||", "'"}
	for _, prefix := range prefixes {
		text = strings.TrimPrefix(text, prefix)
	}
	for _, kw := range kws {
		text = strings.ReplaceAll(text, kw, "")
	}

	return text
}
