# Tag(s) flag

All events have an optional `[]string` property called "tags", which allow you to further sort and filter events. Any command that appends or interacts with events can use the tag flag by simply adding `--tag <custom-tag-name>` or `-t <tag>`. Any string can be specified of any length to be used as a tag, and multiple tags can be added to a single event.

  `pt <command> [--tag <tag-name>]...`

## Examples

```bash
# add a tag to a note
$ pt note --tag recon Port 80 is open

# add multiple tags
$ pt note -t recon -t ftp --tag critical FTP allows anonymous login
```
