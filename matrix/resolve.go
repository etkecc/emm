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
//
//nolint:goerr113 // no need
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

// resolveAlias resolves room alias to a room ID
func resolveAlias(alias id.RoomAlias) (id.RoomID, error) {
	var resp *mautrix.RespAliasResolve
	err := retry(func() error {
		var resolveErr error
		resp, resolveErr = client.ResolveAlias(alias)

		return resolveErr
	})
	if err != nil {
		return "", err
	}

	return resp.RoomID, err
}

// resolveIgnored parses string with list of MXIDs to the `ignore` map
func resolveIgnored(list string) {
	slice := strings.Split(list, ",")
	size := len(slice)
	if size == 0 {
		return
	}
	v := struct{}{}
	ignored = make(map[id.UserID]struct{}, len(slice))
	for _, user := range slice {
		mxid := id.UserID(strings.TrimSpace(user))
		ignored[mxid] = v
	}
}
