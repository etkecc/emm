# emm: Export Matrix Messages

A CLI tool that joins the room and exports last N messages to the file you specified.

## Features

* Get messages from any matrix room with pagination (if limit greather than page, to prevent timeout errors) or without it (if limit less or equals page)
* Export messages to one file for all messages
* Export each message in separate file
* Custom templates supported (`contrib` contains an example of hugo post template, [etke.cc/news](https://etke.cc/news) can be used as reference)
* Delegation and aliases supported
* `Anyone`/`world_readable` access supported without invite

## Usage

### Full example

That's how [etke.cc/news](https://etke.cc/news) generated

```bash
# using event ID for filename
emm -hs hs.url -u user -p pass -r "#room:hs.url" -t contrib/hugo-post-template.md -o /tmp/%s.md
# using a custom slug for filename
emm -hs hs.url -u user -p pass -r "#room:hs.url" -t contrib/hugo-post-template.md -o /tmp/{{ .CreatedAtDate }}.md
```

> **NOTE**: we at [etke.cc](https://etke.cc) do not spam customers with messages, attempting to be as less annoying as possible,
> so we send about 1-2 messages **per week**.
> If you are sending more than 1 message per day, the `{{ .CreatedAtDate }}` will not work as expected,
> because it will store only a single message per file. Instead, you can either use other template variables (see below),
> or use `-a` flag to append messages to the same file.

### Documentation

```bash
Usage of emm:
  -a	Append to the output file. Useful when you are using multi output mode, but want to have multiple messages in a single file
  -hs string
    	Homeserver URL (supports delegation)
  -l int
    	Messages limit
  -o string
    	Output filename. If it contains %s, it will be replaced with event ID (one message per file, old way), or you can use Go template syntax to use all fields (new way)
  -p string
    	Password of the matrix user
  -r string
    	Room ID or alias
  -s string
    	Load messages since this timestamp. Format (RFC3339): YYYY-MM-DDTHH:MM:SSZ, e.g. 2023-10-01T00:00:00Z
  -t string
    	Template file. Default is JSON message struct
  -u string
    	Username/Login of the matrix user
```

**Template syntax**

You can modify the output file name using Go template syntax. The following fields are available:

* `{{ .ID }}` - event ID
* `{{ .URLSafeID }}` - event ID with URL-safe encoding
* `{{ .Replace }}` - old event ID that was replaced by this message
* `{{ .ReplacedNote }}` - a simple ` (updated)` note if this message replaces another one
* `{{ .Author }}` - author's MXID
* `{{ .Title }}` - first line of the message with partially-stripped markdown
* `{{ .Text }}` - full message text (markdown plain-text)
* `{{ .HTML }}` - full message text (HTML, if available)
* `{{ .CreatedAtDate }}` - created at date in 2006-01-02 format (date only)
* `{{ .CreatedAt }}` - created at date in 2006-01-02 15:04:05 UTC format (date and time)
* `{{ .CreatedAtFull }}` - created at date in default Go's representation format

Using those vars you could have output files likes:

```
-o /tmp/{{ .CreatedAtDate }}-{{ .Title }}.md
```

All fields used in the output filename will be converted to URL-safe format, and the length of the filename will be truncated to 100 characters (to prevent filesystem errors).


## How to get

* Arch Linux [AUR](https://aur.archlinux.org/packages/export-matrix-messages-git/)
* [Releases](https://github.com/etkecc/emm/releases)
* or `go install github.com/etkecc/emm/cmd/emm@latest`
* or `just install` from source code
