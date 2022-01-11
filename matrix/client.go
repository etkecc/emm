package matrix

import (
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/id"
)

var (
	client *mautrix.Client
	room   id.RoomID
)

// Init matrix client
func Init(hs string, login string, password string, roomID id.RoomID, alias id.RoomAlias) error {
	var err error
	client, err = mautrix.NewClient(hs, "", "")
	if err != nil {
		return err
	}
	_, err = client.Login(&mautrix.ReqLogin{
		Type: "m.login.password",
		Identifier: mautrix.UserIdentifier{
			Type: mautrix.IdentifierTypeUser,
			User: login,
		},
		Password:         password,
		StoreCredentials: true,
	})
	if err != nil {
		return err
	}

	if roomID == "" {
		roomID, err = resolveAlias(alias)
		if err != nil {
			return err
		}
	}
	room = roomID

	return nil
}
