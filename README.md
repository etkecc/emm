# emm: Export Matrix Messages [![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/etkecc) [![coverage report](https://gitlab.com/etke.cc/emm/badges/main/coverage.svg)](https://gitlab.com/etke.cc/emm/-/commits/main) [![Go Report Card](https://goreportcard.com/badge/gitlab.com/etke.cc/emm)](https://goreportcard.com/report/gitlab.com/etke.cc/emm) [![Go Reference](https://pkg.go.dev/badge/gitlab.com/etke.cc/emm.svg)](https://pkg.go.dev/gitlab.com/etke.cc/emm)

A CLI tool that joins the room and exports last N messages to the file you specified.

## Features

* Get messages from any matrix room
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

* [Releases](https://gitlab.com/etke.cc/emm/-/releases) for freebsd, linux and MacOS
* or `go install gitlab.com/etke.cc/emm@latest` / `make install`
* or from source code
