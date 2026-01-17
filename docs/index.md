# Praetor CLI Manual

Praetor, or  `pt`, is a command-line interface for managing penetration testing engagements and automating note-taking. 

## Installation

You can find installation instructions on the [README](https://github.com/lachlanharrisdev/praetor) or for the corresponding release on the repo [releases page](https://github.com/lachlanharrisdev/praetor/releases).

## Configuration

Configuration for individual engagements lives inside the `{engagement}/.praetor/` directory. This directory is automatically created and contains the event log and other files necessary for Praetor's operation.

Global configuration files are located in the user's home directory under `~/.config/praetor/`. The primary configuration file is stored in `config.json`, and the engagement folder template is located in the `template/` directory. The contents of this folder will be cloned into every new engagement created with the `pt start` command.

## Support

* Ask for support in [discussions](https://github.com/lachlanharrisdev/praetor/discussions)
* Report a bug [here](https://github.com/lachlanharrisdev/praetor/issues/new?template=bug_report.md)
* Request a feature [here](https://github.com/lachlanharrisdev/praetor/issues/new?template=feature_request.md)
* Report a security vulnerability [here](https://github.com/lachlanharrisdev/praetor/security/advisories)
