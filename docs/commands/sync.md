# pt sync family

Sync has multiple different utilities for syncing files with the event log. This command and its subcommands are what allow Praetor to integrate with almost any existing notetaking workflow.

## pt sync

The root `pt sync` command will sync all tracked files. It will check for changes and append events for changed files.

  `pt sync [options]`

### Examples

```bash
# sync tracked files
$ pt sync

# show output in JSON format
$ pt sync -f json
```

## pt sync add

Add a file from a filepath to sync. Running this command will also trigger an initial sync event and append it to the event log. Note that while it's recommended to have tracked files inside your engagement directory, they don't have to be.

  `pt sync add [options] <filepath>`

### Examples

```bash
# add a basic file
$ pt sync add report.md

# add a file with a tag
$ pt sync add -t critical vulnerabilities/ftp.md
```

## pt sync remove

Stop a file from syncing. Running this command will also send an event to the event log that shows the file is no longer tracked. You can only remove files from being synced that have been added with the `pt sync add` command.

  `pt sync remove [options] <filepath>`
  
### Examples

```bash
# stop a file from being tracked
$ pt sync remove report.md
```

## pt sync list

Shows a list of all files currently being tracked in this engagement. 

  `pt sync list [options]`

### Examples

```bash
# show all tracked files
$ pt sync list
```
