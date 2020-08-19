![Version Badge](https://img.shields.io/badge/Version-4.1.1-informational)
![Maintenance Badge](https://img.shields.io/badge/Maintained-yes-success)
[![GoDoc reference example](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/snhilde/statusbar4)
[![GoReportCard example](https://goreportcard.com/badge/github.com/snhilde/statusbar4)](https://goreportcard.com/report/github.com/snhilde/statusbar4)


# statusbar4
`statusbar4` formats and displays information for the [`dwm`](https://dwm.suckless.org/) statusbar.


## Table of Contents
1. [Overview](#overview)
1. [Documentation](#documentation)
1. [Included modules](#included-modules)
1. [Changelog](#changelog)
	1. [4.1.2](#412)
	1. [4.1.1](#411)
	1. [4.1](#41)
	1. [4.0](#40)
	1. [3.0](#30)
	1. [2.0](#20)
	1. [1.0](#10)


## Overview
<!-- TODO -->
This is the framework that controls the modular routines for calculating, formatting, and displaying information on the statusbar.
For modules currently integrated with this framework, see [sb4routines](https://godoc.org/github.com/snhilde/sb4routines).
<!-- TODO -->


## Documentation
Usage guidelines and documentation are hosted at [GoDoc](https://godoc.org/github.com/snhilde/statusbar4).


## Included modules
`statusbar4` is modular by design, and it's simple to build and integrate modules; you only have to implement [two methods](https://godoc.org/github.com/snhilde/statusbar4#RoutineHandler).

This repository includes these modules to get up and running quickly:
| Module       | Documentation                                                            | Major usage          |
| ------------ | ------------------------------------------------------------------------ | -------------------- |
| `sbbattery`  | [GoDoc docs](https://godoc.org/github.com/snhilde/statusbar4/sbbattery)  | Battery usage        |
| `sbcputemp`  | [GoDoc docs](https://godoc.org/github.com/snhilde/statusbar4/sbcputemp)  | CPU temperature      |
| `sbcpuusage` | [GoDoc docs](https://godoc.org/github.com/snhilde/statusbar4/sbcpuusage) | CPU usage            |
| `sbdisk`     | [GoDoc docs](https://godoc.org/github.com/snhilde/statusbar4/sbdisk)     | Filesystem usage     |
| `sbfan`      | [GoDoc docs](https://godoc.org/github.com/snhilde/statusbar4/sbfan)      | Fan speed            |
| `sbload`     | [GoDoc docs](https://godoc.org/github.com/snhilde/statusbar4/sbload)     | System load averages |
| `sbnetwork`  | [GoDoc docs](https://godoc.org/github.com/snhilde/statusbar4/sbnetwork)  | Network usage        |
| `sbnordvpn`  | [GoDoc docs](https://godoc.org/github.com/snhilde/statusbar4/sbnordvpn)  | NordVPN status       |
| `sbram`      | [GoDoc docs](https://godoc.org/github.com/snhilde/statusbar4/sbram)      | RAM usage            |
| `sbtime`     | [GoDoc docs](https://godoc.org/github.com/snhilde/statusbar4/sbtime)     | Current date/time    |
| `sbtodo`     | [GoDoc docs](https://godoc.org/github.com/snhilde/statusbar4/sbtodo)     | TODO list display    |
| `sbvolume`   | [GoDoc docs](https://godoc.org/github.com/snhilde/statusbar4/sbvolume)   | Volume percentage    |
| `sbweather`  | [GoDoc docs](https://godoc.org/github.com/snhilde/statusbar4/sbweather)  | Weather information  |

## Changelog
### 4.1.2
* Updated documentation.

### 4.1.1
* Changes made to formatting and style as per `gofmt` and `golint` recommendations.

### 4.1
* Moved routines and engine into one common repository.

### 4.0
* Complete rewrite in `go`.
* Simpler formatting and customization.

### 3.0
* Added support for concurrency.
* Made routines modular.

### 2.0
* Ported to Linux.

### 1.0
* Initial release (OpenBSD only).
