# deployit

![CI/CD](https://github.com/NobleMajo/deployit/actions/workflows/go-bin-release.yml/badge.svg)
![CI/CD](https://github.com/NobleMajo/deployit/actions/workflows/go-test-build.yml/badge.svg)  
![MIT](https://img.shields.io/badge/license-MIT-blue.svg)
![](https://img.shields.io/badge/dynamic/json?color=green&label=watchers&query=watchers&suffix=x&url=https%3A%2F%2Fapi.github.com%2Frepos%2FNobleMajo%2Fdeployit)
![](https://img.shields.io/badge/dynamic/json?color=yellow&label=stars&query=stargazers_count&suffix=x&url=https%3A%2F%2Fapi.github.com%2Frepos%2FNobleMajo%2Fdeployit)
![](https://img.shields.io/badge/dynamic/json?color=navy&label=forks&query=forks&suffix=x&url=https%3A%2F%2Fapi.github.com%2Frepos%2FNobleMajo%2Fdeployit)

Uses ssh + sftp to deploy configurations to Linux servers and can execute simple commands.

## Table of Contents

- [Table of Contents](#table-of-contents)
- [Configuration](#configuration)
- [Requirements](#requirements)
- [Getting Started](#getting-started)
- [Quick help](#quick-help)
- [Install via go](#install-via-go)
- [Install via wget](#install-via-wget)
- [Build requirements](#build-requirements)
- [Build Instructions](#build-instructions)
- [Install go](#install-go)

## Configuration

DeployIT is easily configured using environment variables or a `.env` file in the working directory.

Here is a WireGuard example:

```bash
DIT_NODE1=ssh://<user>@<host1>*<priv-key-path>
DIT_NODE1_TASK1=UPLOAD@./node1.wg0.conf@/etc/wireguard/wg0.conf
DIT_NODE1_TASK2=CMD@sudo wg-quick down wg0 || true && sudo wg-quick up wg0
DIT_NODE1_TASK3=DOWNLOAD@/etc/wireguard/wg0.conf@./test.node1.wg0.conf

DIT_NODE2=ssh://<user>@<host2>*<priv-key-path>
DIT_NODE2_TASK1=UPLOAD@./node2.wg0.conf@/etc/wireguard/wg0.conf
DIT_NODE2_TASK2=CMD@sudo wg-quick down wg0 || true && sudo wg-quick up wg0
DIT_NODE2_TASK3=DOWNLOAD@/etc/wireguard/wg0.conf@./test.node2.wg0.conf
```

This example deploys two local WireGuard configs from `./nodeX.wg0.conf` to the selected hosts.
It then runs `wg-quick` down and up on each interface to reload the config.
To verify deployment, it downloads each remote config to `./test.node1.wg0.conf` and `./test.node2.wg0.conf`.

For this example, make sure the user you are using has permissions to access `/etc/wireguard` on the server.
If a password is used, use `!your-password` instead of `*<priv-key-path>`.

# User Guide

## Requirements

Linux- or macos-like systems with `go` or `wget & tar` installed.

## Getting Started

Start the latest repo version directly without leaving stuff in the current working dir:

```sh
go run github.com/NobleMajo/deployit@latest
```

## Quick help

```sh
go run github.com/NobleMajo/deployit@latest -h
```

## Install via go

###### _For this section go is required, check out the [install go guide](#install-go)._

```sh
go install github.com/NobleMajo/deployit@latest
```

## Install via wget

```sh
export CUSTOM_BIN_DIR="/usr/local/bin" # <- change if needed
export CUSTOM_VERSION="" # <- set latest version here

rm -rf $CUSTOM_BIN_DIR/deployit
wget https://github.com/NobleMajo/deployit/releases/download/v$CUSTOM_VERSION/deployit-v$CUSTOM_VERSION-linux-amd64.tar.gz -O /tmp/deployit.tar.gz
tar -xzvf /tmp/deployit.tar.gz -C $CUSTOM_BIN_DIR/ deployit
rm /tmp/deployit.tar.gz
```

# Build

## Build requirements

To build, you need to install go.
The required go version is in the `go.mod` file.

## Build Instructions

###### _For this section go is required, check out the [install go guide](#install-go)._

Clone the repo:

```sh
git clone https://github.com/NobleMajo/deployit.git
cd deployit
```

Build the deployit binary from source code:

```sh
make build
./deployit
```

# Development

###### _For this section go is required, check out the [install go guide](#install-go)._

This part is work in progress, I want to use 'AIR' as auto-reload tool:

```sh
make dev #WIP
```

## Install go

The required go version for this project is in the `go.mod` file.

To install and update go, I can recommend the following repo:

```sh
git clone git@github.com:udhos/update-golang.git golang-updater
cd golang-updater
sudo ./update-golang.sh
```

# Contributing

Contributions to this project are welcome!  
Interested users can refer to the guidelines provided in the [CONTRIBUTING.md](CONTRIBUTING.md) file to contribute to the project and help improve its functionality and features.

# License

This project is licensed under the [MIT license](LICENSE), providing users with flexibility and freedom to use and modify the software according to their needs.

# Disclaimer

This project is provided without warranties.  
Users are advised to review the accompanying license for more information on the terms of use and limitations of liability.
