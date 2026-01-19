# pt note

Record a note manually and append it to the event log. Planned to soon also accept text from STDIN. Please note that surrounding quotes are not necessary, however options cannot be parsed once the first non-flag argument is provided.

  `pt note [options] <comment...>

## Examples

```sh
# take a note
$ pt note Hello world!

# add a tag
$ pt note -t recon Port 80 is open.

# containing quotes will be ignored, so you can use brackets or other symbols without errors
$ pt note "I can type (almost) anything!"
```
