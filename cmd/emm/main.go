package main

import (
	"log"

	"maunium.net/go/mautrix/id"

	"github.com/etkecc/emm/internal/export"
	"github.com/etkecc/emm/internal/flags"
	"github.com/etkecc/emm/internal/matrix"
)

var cfg *flags.Config

func main() {
	var err error
	log.Println("parsing command line arguments..")
	cfg, err = flags.Parse()
	if err != nil {
		panic(err)
	}

	log.Println("initializing client...")
	resolveHS()
	err = matrix.Init(*cfg.HS, *cfg.Login, *cfg.Password, *cfg.Ignore)
	if err != nil {
		panic(err)
	}
	defer matrix.Exit()

	resolveRooms()

	messages := map[id.EventID]*matrix.Message{}
	for _, roomID := range cfg.RoomIDs {
		log.Println("loading messages from " + roomID + "...")
		roomMessages, err := matrix.Messages(roomID, *cfg.Limit, cfg.StartAt)
		if err != nil {
			panic(err)
		}
		for id, msg := range roomMessages {
			messages[id] = msg
		}
	}

	err = export.Run(*cfg.Template, *cfg.Output, messages, *cfg.Append)
	if err != nil {
		panic(err)
	}
}

func resolveHS() {
	if cfg.NoDelegation == nil || !*cfg.NoDelegation {
		log.Println("resolving homeserver...")
		hs, err := matrix.ResolveServer(*cfg.HS)
		if err != nil {
			panic(err)
		}
		cfg.HS = &hs
	}
}

func resolveRooms() {
	for _, room := range cfg.Rooms.Slice() {
		alias, err := matrix.IsRoom(room)
		if err != nil {
			panic(err)
		}
		if !alias {
			cfg.RoomIDs = append(cfg.RoomIDs, id.RoomID(room))
			continue
		}
		log.Println("resolving " + room + " alias...")
		roomID, err := matrix.ResolveAlias(id.RoomAlias(room))
		if err != nil {
			panic(err)
		}
		cfg.RoomIDs = append(cfg.RoomIDs, roomID)
	}
}
