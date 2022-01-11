package matrix

import (
	"errors"
	"strings"

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/id"
)

// ResolveServer resolves actual homeserver URL from possible delegated host
func ResolveServer(homeserver string) (string, error) {
	discover, err := mautrix.DiscoverClientAPI(homeserver)
	if err != nil {
		return "", err
	}
	return discover.Homeserver.BaseURL, nil
}

// IsRoom checks if room is valid room alias or room ID. True = alias, False = room ID
func IsRoom(room string) (bool, error) {
	if room == "" {
		return false, errors.New("room is not set")
	}

	if strings.LastIndex(room, ":") == -1 {
		return false, errors.New("not a valid room id or alias")
	}

	if strings.HasPrefix(room, "#") {
		return true, nil
	}

	if strings.HasPrefix(room, "!") {
		return false, nil
	}

	return false, errors.New("not a valid room id or alias")
}

func resolveAlias(alias id.RoomAlias) (id.RoomID, error) {
	resp, err := client.ResolveAlias(alias)
	if err != nil {
		return "", err
	}

	return resp.RoomID, err
}
