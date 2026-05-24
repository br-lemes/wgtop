# wgtop

A lightweight, real-time terminal monitor for WireGuard interfaces, written in
Go. It provides a quick, `top`-like overview of connected devices, their
status, and traffic statistics.

This tool is tailored for hub-and-spoke (central server) network topologies,
specifically monitoring peers configured with `/32` (IPv4).

## Features

* **Real-time Updates:** Automatically refreshes device statistics every few
  seconds.
* **Filtering:** Filter peers by status (`online`, `offline`, or `all`).
* **Sorting:** Sort devices by IP, hostname, connection time, bytes sent, or
  bytes received.
* **Clean CLI Interface:** Designed to be minimal, scannable, and easy to read.

## Installation

Build and install directly:

```sh
go install github.com/br-lemes/wgtop@latest
```

Or clone the source code first:

```sh
git clone https://github.com/br-lemes/wgtop
cd wgtop
go build
```

> **Note:** Since it interacts with WireGuard interfaces via `wgctrl`, you might
need administrative privileges (`sudo`) to run the tool depending on your system
configuration.

## Usage

Run the tool without arguments to start with default settings (online devices,
sorted by IP):

```bash
sudo ./wgtop
```

### Command Line Flags

You can customize the output using the following flags:

```text
-sort string
      Sort devices by: ip, name, time, sent, recv (default "ip")
-status string
      Filter devices by status: online, offline, all (default "online")
-version
      Show version information
```

### Examples

Sort devices by downloaded data (`recv`) and show all peers (including offline
ones):

```bash
sudo ./wgtop -sort recv -status all
```

Monitor only currently active peers sorted by connection time:

```bash
sudo ./wgtop -sort time -status online
```

## License

This project is licensed under the BSD Zero Clause License (0BSD).
