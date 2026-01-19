# Format flag

To change the format of any command's output, specify the `-f` or `--format` flag. Currently, only JSON ("j", "js", "jsn", "json") and terminal (default; "t", "term", "terminal") are supported, but Markdown support will come soon. Suggestions for formats are more than welcome.

  `pt <command> [--format <format>]`

## Examples

```bash
# render in JSON
$ pt status --format json

$ pt note -f j Hello, World!

# render for terminal if your default format is JSON or other
$ pt run --format terminal echo Hello!

$ nmap 127.0.0.1 | pt capture -f t
```
