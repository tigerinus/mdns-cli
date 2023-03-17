# mdns-cli

[![Go Report Card](https://goreportcard.com/badge/github.com/tigerinus/mdns-cli)](https://goreportcard.com/report/github.com/tigerinus/mdns-cli) [![goreleaser](https://github.com/tigerinus/mdns-cli/actions/workflows/release.yml/badge.svg)](https://github.com/tigerinus/mdns-cli/actions/workflows/release.yml)

A cross-platform CLI for mdns

## Motivation

There is no `avahi-browse` under Windows, so I had to create one myself.

## Usage

```text
Usage:
  mdns-cli [command]

Commands
  browse      Browse for services and instances for each service

Additional Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Show version

Flags:
  -h, --help   help for mdns-cli

Use "mdns-cli [command] --help" for more information about a command.
```
