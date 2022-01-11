package main

import (
	"maunium.net/go/mautrix/id"

	"gitlab.com/etke.cc/emm/export"
	"gitlab.com/etke.cc/emm/flags"
	"gitlab.com/etke.cc/emm/matrix"
)

var cfg *flags.Config

func main() {
	var err error
	cfg, err = flags.Parse()
	if err != nil {
		panic(err)
	}
	resolve()

	err = matrix.Init(*cfg.HS, *cfg.Login, *cfg.Password, cfg.RoomID, cfg.RoomAlias)
	if err != nil {
		panic(err)
	}

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
	// Resolve homeserver
	hs, err := matrix.ResolveServer(*cfg.HS)
	if err != nil {
		panic(err)
	}
	cfg.HS = &hs

	// Resolve room type
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
