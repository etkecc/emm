# emm: Export Matrix Messages

A CLI tool that joins the room and exports last N messages to the file you specified.

## Features

* Get messages from any matrix room with pagination (if limit greather than page, to prevent timeout errors) or without it (if limit less or equals page)
* Export messages to one file for all messages
* Export each message in separate file
* Custom templates supported (`contrib` contains an example of hugo post template, [etke.cc/webite](https://gitlab.com/etke.cc/website) can be used as reference)
* Delegation and aliases supported
* `Anyone`/`world_readable` access supported without invite

## Usage

### Full example

That's how [etke.cc/website](https://gitlab.com/etke.cc/website) news generated

```bash
emm -hs hs.url -u user -p pass -r "#room:hs.url" -t contrib/hugo-post-template.md -o /tmp/%s.md
```

### Documentation

```bash
Usage of emm:
  -hs string
    	Homeserver URL (supports delegation)
  -l int
    	Messages limit
  -o string
    	Output filename. If it contains %s, it will be replaced with event ID (one message per file)
  -p string
    	Password of the matrix user
  -r string
    	Room ID or alias
  -t string
    	Template file. Default is JSON message struct
  -u string
    	Username/Login of the matrix user
```

## How to get

* Arch Linux [AUR](https://aur.archlinux.org/packages/export-matrix-messages-git/)
* [Releases](https://github.com/etkecc/emm/releases)
* or `go install github.com/etkecc/emm/cmd/emm@latest`
* or from source code
