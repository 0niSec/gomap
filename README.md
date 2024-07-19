# Table of Contents

- [Table of Contents](#table-of-contents)
- [:eyes: What Is This?](#eyes-what-is-this)
- [:sparkles: Features](#sparkles-features)
- [:computer: Installation](#computer-installation)
  - [:exclamation: Important Note](#exclamation-important-note)
  - [Standalone Binary](#standalone-binary)
    - [Alternative - Go CLI](#alternative---go-cli)
      - [Downloading \& Installing Go](#downloading--installing-go)
      - [Installing Gomap](#installing-gomap)
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

- [Goroutines](https://go.dev/doc/effective_go#goroutines) for fast runtimes
- Utilizes TCP SYN scanning methods for fast results. See more information on the topic and how `nmap` does it [here]("https://nmap.org/book/synscan.html")
- Service fingerprinting (not yet implemented)
- OS fingerprinting (not yet implemented)
- Banner grabbing (not yet implemented)

# :computer: Installation

## :exclamation: Important Note

<font color="red">
Gomap was built with Linux in mind since as a penetration testing tool, the typical demographic OS is some form of Linux. With that being said, Gomap <strong>will not work on Windows or Mac</strong>.
</font>

<br />
Gomap was developed on 64-bit Kali Linux. It should work on other Linux distros, but I haven't tested it.

## Standalone Binary

Standalone binaries can be downloaded from the [releases](https://github.com/0niSec/gomap/releases) page. You may also use the Go CLI following the below steps.

### Alternative - Go CLI

While I do provide the binary through Github, the binary can also be installed through the Go CLI. For this method, you will need to have Go installed.

#### Downloading & Installing Go

To download Go, follow the instructions at Golang's own website, [here](https://go.dev/dl). Then follow the [install instructions](https://go.dev/doc/install) for your OS.

> [!NOTE]
> Linux users can also install via their package manager  (e.g.`apt install golang-go`). This does not appear to be present on Go's website, but it's how I installed it and have had no issues with Kali. YMMV depending on distro.

#### Installing Gomap

To install Gomap, you can use the Go CLI:

```bash
go install github.com/0niSec/gomap@latest
```

## :whale: Docker

Docker is used a lot nowadays and I wanted to include it as an option because:

- It works on all systems
- The Docker image will use the latest build from Go. You'll always be using the latest version.
- No need to install Go

To install Docker, follow [their guide](https://docs.docker.com/engine/install/).

```bash
docker pull 0niSec/gomap:latest
```

# Usage

> [!WARNING]
> Gomap requires elevated privileges to run. This is due to the fact that it uses raw sockets to send and receive packets. Any command will need to be run with `sudo`. If you want to be able to not enter a password, add your user to the sudoers file and specify NOPASSWD for gomap.

## :penguin: Linux

Running on Linux is simple! With Gomap installed using `go install` or downloaded directly, and the binary added to your PATH (`go install` installs directly to `$GOPATH/bin`):

```text
NAME:
   gomap - The Go port scanner

USAGE:
   gomap [global options] command [command options] 

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --ports value, -p value    Port ranges to scan (e.g. 80,443,8000-8100)
   --quiet, -q                Don't print the banner and other noise (default: false)
   --target value, -t value   The target to scan
   --timeout value, -T value  Timeout for the connection (default: 10s)
   --output value, -o value   Output file
   --help, -h                 show help
```

> [!NOTE]
> More flags and features will be added here as they are developed.

## :whale: Docker

Running using Docker is as easy as

```bash
docker run 0niSec/gomap -p <PORTS> -t <TARGET_IP>
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
