// Package flags manages command line arguments and flags
package flags

import (
	"errors"
	"flag"

	"maunium.net/go/mautrix/id"
)

// Config from the command line args and flags
type Config struct {
	// HS is matrix homeserver (supports delegation)
	HS *string
	// Login is matrix user login
	Login *string
	// Password is matrix user password
	Password *string
	// Room id or alias (raw)
	Room *string
	// Room ID
	RoomID id.RoomID
	// RoomAlias
	RoomAlias id.RoomAlias
	// Ignore messages by following MXIDs
	Ignore *string
	// Limit of messages
	Limit *int
	// Template file
	Template *string
	// Output file name. If filename contains %s it will be replaced with event ID (in case of one file per message)
	Output *string
}

//nolint // as-is for now
func (cfg *Config) validate() error {
	if cfg.HS == nil || *cfg.HS == "" {
		return errors.New("-hs is not set. You must specify homeserver URL")
	}
	if cfg.Login == nil || *cfg.Login == "" {
		return errors.New("-u is not set. You must specify username/login of the matrix user")
	}
	if cfg.Password == nil || *cfg.Password == "" {
		return errors.New("-p is not set. You must specify password of the matrix user")
	}
	if cfg.Room == nil || *cfg.Room == "" {
		return errors.New("-r is not set. You must specify room id or alias")
	}
	if cfg.Output == nil || *cfg.Output == "" {
		return errors.New("-o is not set. You must specify output filename")
	}
	if cfg.Limit == nil || *cfg.Limit < 0 {
		limit := 0
		cfg.Limit = &limit
	}
	if cfg.Template == nil {
		empty := ""
		cfg.Template = &empty
	}

	return nil
}

// Parse command line arguments and flags
func Parse() (*Config, error) {
	cfg := &Config{
		HS:       flag.String("hs", "", "Homeserver URL (supports delegation)"),
		Login:    flag.String("u", "", "Username/Login of the matrix user"),
		Password: flag.String("p", "", "Password of the matrix user"),
		Room:     flag.String("r", "", "Room ID or alias"),
		Limit:    flag.Int("l", 0, "Messages limit"),
		Ignore:   flag.String("i", "", "Ignore messages by following MXIDs, separated by comma"),
		Template: flag.String("t", "", "Template file. Default is JSON message struct"),
		Output:   flag.String("o", "", "Output filename. If it contains %s, it will be replaced with event ID (one message per file)"),
	}
	flag.Parse()
	err := cfg.validate()

	return cfg, err
}
