<!--<br />-->
<div align="center">
  <a href="https://github.com/lachlanharrisdev/praetor">
    <img src=".github/praetor-white-transparent.png" alt="Logo" width="80" height="80"/>
  </a>

  <h1 align="center">Praetor</h1>

  <p align="center">
    Sophisticated engagement management & low-friction notetaking for penetration testing
    <br />
    <a href="https://github.com/lachlanharrisdev/praetor"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/lachlanharrisdev/praetor/blob/main/CONTRIBUTING.md">Contribute</a>
    &middot;
    <a href="https://github.com/lachlanharrisdev/praetor/issues/new?template=bug_report.md">Report Bug</a>
    &middot;
    <a href="https://github.com/lachlanharrisdev/praetor/issues/new?template=feature_request.md">Request Feature</a>
  </p>
</div>

<br/>
<h2>Introduction</h2>

Praetor is a CLI tool built solely to reduce administrative friction in penetration testing. `pt` offers a clean set of integrated utilities to:
* Manage engagement contexts, directories & archiving
* Provide intuitive commands to record, update, report on, format and export commands, notes & outputs
* Manage a forward-immutable, centralised & lightweight event log to contain context for each engagement
* Give endless possibilities to customise to suit Praetor to your team's needs and reduce the cognitive load of administrative work.

<br/>
<h2 align="right">Get Started</h2>

<h3>Installation</h3>

1. Go to the [releases](https://github.com/lachlanharrisdev/praetor/releases/) page and download the desired version & `checksums.txt` file. It should look like `praetor_{version}_{os}_{arch}.targ.gz`
2. Run the following commands in your shell to extract and move to your `bin`
```sh
tar xzf praetor_{version}_{os}_{arch}.tar.gz
sudo mv pt /usr/local/bin/
```
3. (Optional) Verify the checksums before moving to `bin`
```sh
sha256sum -c checksums.txt
```
4. Verify the installation succeeded
```sh
pt version
```

<br/>
<h3>Usage</h3>

`pt` has countless methods of use. There's no one correct way to use it and it all depends on your environment, existing methods of administration and the needs of your team. It's best to keep up with the documentation and refer to each individual commands use.

Some basic usage could look as follows:

1. Create a new engagement directory
```bash
$ pt init test-eng
/home/{user}/engagements/test-eng/

$ cd test-eng
```
2. Take your first note
```bash
$ pt note Engagement begun. Provided IP: 123.45.67.89
```
3. Record a tool output
```bash
$ nmap -sC 123.45.67.89 | pt capture

# or:
$ nmap -sC -o nmap_result.txt 123.45.67.89
$ pt capture nmap_result.txt

# or:
$ pt run nmap -sC 123.45.67.89
```
4. View the last few events
```bash
$ pt list 3
```

<br/>

<h2 align="right">Contributing</h2>

<br/>

Praetor follows most standard conventions for contributing, and accepts any contributions from documentation improvements, bug triage / fixes, small features or any updates for [issues in the backlog](https://github.com/lachlanharrisdev/praetor/issues?q=is%3Aissue%20state%3Aopen%20label%3A%22status%3A%20backlog%22). For more information on contributing please see [CONTRIBUTING.md](https://github.com/lachlanharrisdev/praetor/blob/main/.github/CONTRIBUTING.md).

<br/>
<h3>Codespaces</h3>

Praetor has full support for Github Codespaces. These are recommended for small changes or devices with no access to a Linux environment. You can use the buttons below to open the repository in a web-based editor and get started.

[![Open in GitHub Codespaces](https://github.com/codespaces/badge.svg)](https://codespaces.new/lachlanharrisdev/praetor?quickstart=1)

<h3>Dev Containers</h3>

We also have full support for Dev Containers. These provide a reproducible development environment that automatically isolates the project and installs the officially supported toolchain. 

Clicking the below button will open up VS Code on your local machine, clone this repository and open it automatically inside a development container.

[![Open in Dev Containers](https://img.shields.io/badge/Open%20In%20Dev%20Container-0078D4?style=for-the-badge&logo=visual%20studio%20code&logoColor=white)](https://vscode.dev/redirect?url=vscode://ms-vscode-remote.remote-containers/cloneInVolume?url=https://github.com/lachlanharrisdev/praetor)

<h3>Local Development</h3>

For local development, please refer to [CONTRIBUTING.md](https://github.com/lachlanharrisdev/praetor/blob/main/.github/CONTRIBUTING.md). Again, we follow most conventions so local development involves the standard flow of `fork-PR-merge`.

<br/>

---

<br/>

> This project is in a WIP state. It is currently unstable and not recommended for use within automated systems or under strict compliance policies. All code is open source and aims to have a minimal, secure footprint, but in it's pre-release stages no guarantees can be made.
>
> This project is licensed under the GPL-3.0 License. Please see [LICENSE](https://github.com/lachlanharrisdev/praetor?tab=GPL-3.0-1-ov-file) for more info.
>
> Copyright © Lachlan Harris 2026. All Rights Reserved.
