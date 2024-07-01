# Table of Contents

- [Table of Contents](#table-of-contents)
- [:eyes: What Is This?](#eyes-what-is-this)
- [:sparkles: Features](#sparkles-features)
- [:computer: Installation](#computer-installation)
  - [Downloading \& Installing Go](#downloading--installing-go)
  - [Go CLI](#go-cli)
  - [Standalone Binary](#standalone-binary)
  - [:whale: Docker](#whale-docker)
- [Usage](#usage)
  - [:penguin: Linux](#penguin-linux)
  - [:whale: Docker](#whale-docker-1)
- [:handshake: Contributing](#handshake-contributing)
- [:exclamation: Disclaimer](#exclamation-disclaimer)
- [:rocket: Contributors](#rocket-contributors)


# :eyes: What Is This?

Gomap is a network port scanning tool similar to that of [nmap](https://nmap.org). There's also the very popular [RustScan](https://github.com/rustscan/rustscan). I wanted to experiment with Go and see what it was all about. But before I knew it, I wanted to do more and see what else I could do with the language.

Gomap does not aim to replace `nmap` or `RustScan`, but rather be a coexisting partner in the same space. I'm not trying to be the best network scanner or the fastest (though speed is a goal). I started this as a way to learn the Go programming language and having literally never published anything Open Source of my own design ever, wanted to see what others thought.

# :sparkles: Features

- Faster alternative to `nmap`
- [Goroutines](https://go.dev/doc/effective_go#goroutines) for faster runtimes
- Service fingerprinting (not yet implemented)
- OS fingerprinting (wishlist)

# :computer: Installation

## Downloading & Installing Go

To download Go, follow the instructions at Golang's own website, [here](https://go.dev/dl). Then follow the [install instructions](https://go.dev/doc/install) for your OS.

> [!NOTE]
> Linux users can also install via their package manager  (e.g.`apt install golang-go`). This does not appear to be present on Go's website, but it's how I installed it and have had no issues with Kali. YMMV depending on distro.

## Go CLI

To install Gomap, you can use the Go CLI:

```bash
go install github.com/0niSec/gomap@latest
```

## Standalone Binary

Standalone binaries can be downloaded from the [releases](https://github.com/0niSec/gomap/releases) page for Windows, Linux and MacOS. For other OS-specific methods, please see the corresponding section for your OS.

## :whale: Docker

> [!WARNING]
> The Dockerfile has been created and added and tested, but I have not yet published it to Docker Hub so these instructions will be valid, but **do not work yet.**

Docker is used a lot nowadays and I wanted to include it as an option because:

- It works on all systems
- The Docker image will use the latest build from Go. You'll always be using the latest version.
- No need to install Go

To install Docker, follow [their guide](https://docs.docker.com/engine/install/).

```bash
docker pull 0niSec/gomap:latest
```

# Usage

## :penguin: Linux

Running on Linux is simple! With Gomap installed using `go install`, and the binary added to your PATH:

```text
USAGE:
   gomap [global options] command [command options] 

VERSION:
   0.1.0

AUTHOR:
   0niSec

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --ports value, -p value    Port ranges to scan
   --quiet, -q                Don't print the banner and other noise (default: false)
   --target value, -t value   The target to scan
   --timeout value, -T value  Timeout for the connection (default: 10s)
   --output value, -o value   Output file
   --help, -h                 show help
   --version, -v              print the version
```

## :whale: Docker

Running the binary using Docker is as easy as

```bash
docker run 0niSec/gomap -p <PORTS> -t <TARGET>
```

# :handshake: Contributing

This started as a solo project and I'd love to accept any help people are willing to provide. If you're interested in helping, take a look at the [issues](https://github.com/0niSec/gomap/issues) for anything you'd like to tackle. Please also read the [Code of Conduct](CODE_OF_CONDUCT.md) and [Contributing](CONTRIBUTING.md) for more information.

**By actively participating in contributing to this project, you agree to all of the rules and guidelines set therein.**

# :exclamation: Disclaimer

This tool is meant to be used ethically in Capture the Flag programs such as [MetaCTF](https://metactf.com), [Hack the Box](https://app.hackthebox.com), or [TryHackMe](https://tryhackme.com) (to name a few) or on sanctioned penetration tests that have a formal contract and drawn out engagement. Please do not use this tool on infrastructure that you do not have permission to.

# :rocket: Contributors

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->