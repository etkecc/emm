package main

import (
	"log"

	"maunium.net/go/mautrix/id"

	"gitlab.com/etke.cc/emm/export"
	"gitlab.com/etke.cc/emm/flags"
	"gitlab.com/etke.cc/emm/matrix"
)

var cfg *flags.Config

func main() {
	var err error
	log.Println("parsing command line arguments..")
	cfg, err = flags.Parse()
	if err != nil {
		panic(err)
	}
	resolve()

	log.Println("initializing client...")
	err = matrix.Init(*cfg.HS, *cfg.Login, *cfg.Password, cfg.RoomID, cfg.RoomAlias, *cfg.Ignore)
	if err != nil {
		panic(err)
	}
	defer matrix.Exit()

	log.Println("loading messages...")
	messages, err := matrix.Messages(*cfg.Limit)
	if err != nil {
		panic(err)
	}
	err = export.Run(*cfg.Template, *cfg.Output, messages)
	if err != nil {
		panic(err)
	}
}

func resolve() {
	log.Println("resolving homeserver...")
	hs, err := matrix.ResolveServer(*cfg.HS)
	if err != nil {
		panic(err)
	}
	cfg.HS = &hs

	log.Println("resolving room type...")
	alias, err := matrix.IsRoom(*cfg.Room)
	if err != nil {
		panic(err)
	}
	if alias {
		cfg.RoomAlias = id.RoomAlias(*cfg.Room)
	} else {
		cfg.RoomID = id.RoomID(*cfg.Room)
	}
}
