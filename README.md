# DeployIT
![CI/CD](https://github.com/CoreUnit-NET/deployit/actions/workflows/go-bin-release.yml/badge.svg)
![CI/CD](https://github.com/CoreUnit-NET/deployit/actions/workflows/go-test-build.yml/badge.svg)  
![MIT](https://img.shields.io/badge/license-MIT-blue.svg)
![](https://img.shields.io/badge/dynamic/json?color=green&label=watchers&query=watchers&suffix=x&url=https%3A%2F%2Fapi.github.com%2Frepos%2Fnoblemajo%2Fdeployit)
![](https://img.shields.io/badge/dynamic/json?color=yellow&label=stars&query=stargazers_count&suffix=x&url=https%3A%2F%2Fapi.github.com%2Frepos%2Fnoblemajo%2Fdeployit)
![](https://img.shields.io/badge/dynamic/json?color=navy&label=forks&query=forks&suffix=x&url=https%3A%2F%2Fapi.github.com%2Frepos%2Fnoblemajo%2Fdeployit)

Uses ssh + sftp to deploy configurations to Linux servers and can execute simple commands.

# Config
DeployIT is easily configured using environment variables or an .env file.
Here is a wireguard example:
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

This example deploys 2 different local Wireguard configs from `./nodeX.wg0.conf` to the selected host.
It then runs a wg-quick down and up on that interface to reload the config.
To test if the config and deployment were successful, it downloads the config to `./test.node2.wg0.conf`.

For this example, make sure the user you are using has permissions to access `/etc/wireguard` on the server.
If a password is used, use `!your-password` instead of `*<priv-key-path>`.

# Table of Contents
- [DeployIT](#deployit)
- [Config](#config)
- [Table of Contents](#table-of-contents)
- [Getting Started](#getting-started)
  - [Requirements](#requirements)
  - [Install via go](#install-via-go)
  - [Install via wget](#install-via-wget)
  - [Build](#build)
  - [Install go](#install-go)
- [Contributing](#contributing)
- [License](#license)
- [Disclaimer](#disclaimer)

# Getting Started

## Requirements
None windows system with `go` or `wget & tar` installed.

## Install via go
###### *For this section go is required, check out the [install go guide](#install-go).*

```sh
go install https://github.com/CoreUnit-NET/deployit
```

## Install via wget
```sh
BIN_DIR="/usr/local/bin"
DIT_VERSION="1.3.3"

rm -rf $BIN_DIR/deployit
wget https://github.com/CoreUnit-NET/deployit/releases/download/v$DIT_VERSION/deployit-v$DIT_VERSION-linux-amd64.tar.gz -O /tmp/deployit.tar.gz
tar -xzvf /tmp/deployit.tar.gz -C $BIN_DIR/ deployit
rm /tmp/deployit.tar.gz
```

## Build
###### *For this section go is required, check out the [install go guide](#install-go).*

Clone the repo:
```sh
git clone https://github.com/CoreUnit-NET/deployit.git
cd deployit
```

Build the deployit binary from source code:
```sh
make build
./deployit
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
