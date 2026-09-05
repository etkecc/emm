// Package flags manages command line arguments and flags
package flags

import (
	"errors"
	"flag"
	"strings"
	"time"

	"maunium.net/go/mautrix/id"
)

// Config from the command line args and flags
type Config struct {
	// HS is matrix homeserver (supports delegation)
	HS *string
	// NoDelegation disables /.well-known delegation support
	NoDelegation *bool
	// Login is matrix user login
	Login *string
	// Password is matrix user password
	Password *string
	// Rooms ids or aliases (raw)
	Rooms *StringSliceFlag
	// Room ID
	RoomIDs []id.RoomID
	// Since load only messages since this timestamp. Format (RFC3339): YYYY-MM-DDTHH:MM:SSZ, e.g. 2023-10-01T00:00:00Z
	Since *string
	// StartAt is a converted time.Time from Since
	StartAt time.Time
	// Ignore messages by following MXIDs
	Ignore *string
	// Append enables appending to the output file, for multi-output mode with several messages in one file
	Append *bool
	// Limit of messages
	Limit *int
	// Template file
	Template *string
	// Output file name. If filename contains %s it will be replaced with event ID (in case of one file per message)
	Output *string
}

func (cfg *Config) validate() error {
	if err := cfg.validateCredentials(); err != nil {
		return err
	}
	if cfg.Rooms == nil || cfg.Rooms.String() == "" {
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

	if cfg.Append == nil {
		appendVal := false
		cfg.Append = &appendVal
	}

	if cfg.Since != nil && *cfg.Since != "" {
		t, err := time.Parse(time.RFC3339, *cfg.Since)
		if err != nil {
			return errors.New("-s is not valid. Must be in RFC3339 format: YYYY-MM-DDTHH:MM:SSZ")
		}
		cfg.StartAt = t
	}

	return nil
}

func (cfg *Config) validateCredentials() error {
	if cfg.HS == nil || *cfg.HS == "" {
		return errors.New("-hs is not set. You must specify homeserver URL")
	}
	if cfg.Login == nil || *cfg.Login == "" {
		return errors.New("-u is not set. You must specify username/login of the matrix user")
	}
	if cfg.Password == nil || *cfg.Password == "" {
		return errors.New("-p is not set. You must specify password of the matrix user")
	}
	return nil
}

// StringSliceFlag used for a flag with multiple values, eg: -c one -c two -c three
type StringSliceFlag []string

// String from slice
func (f *StringSliceFlag) String() string {
	return strings.Join(*f, "; ")
}

// Set (append) to a string slice
func (f *StringSliceFlag) Set(value string) error {
	*f = append(*f, value)

	return nil
}

// Slice returns the underlying slice of strings
func (f *StringSliceFlag) Slice() []string {
	if f == nil {
		return nil
	}
	return *f
}

// Parse command line arguments and flags
func Parse() (*Config, error) {
	rooms := &StringSliceFlag{}
	flag.Var(rooms, "r", "Room ID or alias, you can add multiple rooms with this flag, e.g. -r #room1:example.com -r !pHSBxS_8cAoZXw63TskCZ4W5Ea1_LYn9kOlL9zPt0Vk")
	cfg := &Config{
		HS:           flag.String("hs", "", "Homeserver URL (supports delegation)"),
		NoDelegation: flag.Bool("no-delegation", false, "Disable /.well-known delegation support"),
		Login:        flag.String("u", "", "Username/Login of the matrix user"),
		Password:     flag.String("p", "", "Password of the matrix user"),
		Rooms:        rooms,
		Limit:        flag.Int("l", 0, "Messages limit"),
		Append:       flag.Bool("a", false, "Append to the output file. Useful when you are using multi output mode, but want to have multiple messages in a single file"),
		Ignore:       flag.String("i", "", "Ignore messages by following MXIDs, separated by comma"),
		Since:        flag.String("s", "", "Load messages since this timestamp. Format (RFC3339): YYYY-MM-DDTHH:MM:SSZ, e.g. 2023-10-01T00:00:00Z"),
		Template:     flag.String("t", "", "Template file. Default is JSON message struct"),
		Output:       flag.String("o", "", "Output filename. If it contains %s, it will be replaced with event ID (one message per file)"),
	}
	flag.Parse()
	err := cfg.validate()

	return cfg, err
}
