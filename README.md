# DeployIT
![MIT](https://img.shields.io/badge/license-MIT-blue.svg)
![](https://img.shields.io/badge/dynamic/json?color=green&label=watchers&query=watchers&suffix=x&url=https%3A%2F%2Fapi.github.com%2Frepos%2Fnoblemajo%2Fdeployit)
![](https://img.shields.io/badge/dynamic/json?color=yellow&label=stars&query=stargazers_count&suffix=x&url=https%3A%2F%2Fapi.github.com%2Frepos%2Fnoblemajo%2Fdeployit)
![](https://img.shields.io/badge/dynamic/json?color=navy&label=forks&query=forks&suffix=x&url=https%3A%2F%2Fapi.github.com%2Frepos%2Fnoblemajo%2Fdeployit)

![](https://github.com/noblemajo/deployit/actions/workflows/go-test-build.yml/badge.svg)  
![](https://github.com/noblemajo/deployit/actions/workflows/go-tag-release.yml/badge.svg)  
![](https://github.com/noblemajo/deployit/actions/workflows/go-bin-release.yml/badge.svg)  

Uses ssh + sftp to deploy configurations to Linux servers and can execute simple commands.

## Table of Contents
- [DeployIT](#deployit)
  - [Table of Contents](#table-of-contents)
  - [Getting Started](#getting-started)
    - [Install via Go](#install-via-go)
    - [Download executables](#download-executables)
    - [Build via Go](#build-via-go)
    - [Build via docker](#build-via-docker)
  - [Config](#config)
  - [Contributing](#contributing)
    - [Makefile](#makefile)
  - [License](#license)
  - [Disclaimer](#disclaimer)

## Getting Started

### Install via Go

```sh
go install github.com/noblemajo/deployit@latest
```

### Download executables

Go to the [releases](./releases) page and copy the download URL for your target system.

### Build via Go

Clone the repo:
```sh
git clone https://github.com/noblemajo/deployit.git
cd deployit
```

Build the **`bin`** binary from source:
```sh
make build

./bin
```

### Build via docker

*You’re not required to install Go on the host if you have **Make** and **Docker**.*

Use `make docker` to start an interactive shell inside the Go dev container.

Use Make targets such as `make run` in the container, then stop with `exit`.

## Config
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

## Contributing
Contributions to this project are welcome!  
Interested users can refer to the guidelines provided in the [CONTRIBUTING.md](CONTRIBUTING.md) file to contribute to the project and help improve its functionality and features.

### Makefile

The `Makefile` contains many useful Make targets.
It also includes Make targets to trigger releases and support development. Check it out as a starting point.

## License
This project is licensed under the [MIT license](LICENSE), providing users with flexibility and freedom to use and modify the software according to their needs.

## Disclaimer
This project is provided without warranties.  
Users are advised to review the accompanying license for more information on the terms of use and limitations of liability.
