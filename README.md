# Overseerr to Apprise

[![Release](https://github.com/USA-RedDragon/overseerr-to-apprise/actions/workflows/release.yaml/badge.svg)](https://github.com/USA-RedDragon/overseerr-to-apprise/actions/workflows/release.yaml) [![License](https://badgen.net/github/license/USA-RedDragon/overseerr-to-apprise)](https://github.com/USA-RedDragon/overseerr-to-apprise/blob/master/LICENSE) [![go.mod version](https://img.shields.io/github/go-mod/go-version/USA-RedDragon/overseerr-to-apprise.svg)](https://github.com/USA-RedDragon/overseerr-to-apprise) [![GoReportCard](https://goreportcard.com/badge/github.com/USA-RedDragon/overseerr-to-apprise)](https://goreportcard.com/report/github.com/USA-RedDragon/overseerr-to-apprise) [![codecov](https://codecov.io/gh/USA-RedDragon/overseerr-to-apprise/graph/badge.svg?token=6ASKMAKOZE)](https://codecov.io/gh/USA-RedDragon/overseerr-to-apprise)

## Setup

### Server Configuration

The service is configured via environment variables, a configuration YAML file, or command line flags. The [`config.example.yaml`](config.example.yaml) file shows the available configuration options. The command line flags match the schema of the YAML file, i.e. `--http.ipv4_host='0.0.0.0'` would equate to `http.ipv4_host: "0.0.0.0"`. Environment variables are in the same format, however they are uppercase and replace hyphens with underscores and dots with double underscores, i.e. `HTTP__IPV4_HOSTS="0.0.0.0"`.
