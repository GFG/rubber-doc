# RubberDoc!
[![Build Status](https://travis-ci.org/rocket-internet-berlin/RocketLabsRubberDoc.svg?branch=master)](https://travis-ci.org/rocket-internet-berlin/RocketLabsRubberDoc)

Table of Contents
=================

* [What is RubberDoc](#what-is-rubberdoc) 
* [Installation](#installation)
  * [Using the release](#using-the-release)
  * [Using the source](#using-the-source)
* [Configuration](#configuration)
* [Usage](#usage)
* [Help](#help)
* [Examples](#examples)
* [Contribute](#contribute)
* [CREDITS](#credits)


## What is RubberDoc
RubberDoc was designed to support api's documentation generation based on RAML and Blueprint.

## Installation

The latest executables for supported platforms are available from the [release page](https://github.com/rocket-internet-berlin/RocketLabsRubberDoc/releases).

### Using the release

Just extract and start using it:

```
$ wget https://github.com/rocket-internet-berlin/RocketLabsRubberDoc/releases/download/${version}/rubberdoc-${version}.${os}-${arch}.tar.gz
$ tar -zxvf rubberdoc-${version}.${os}-${arch}.tar.gz
$ ./rubberdoc -h
```
### Using the source

```sh
$ git clone https://github.com/rocket-internet-berlin/RocketLabsRubberDoc.git
$ cd RocketLabsRubberDoc
$ make install
```

> Note: Ensure you have installed [Go](https://golang.org/doc/install#tarball) and configured your `GOPATH` and `PATH`.

## Configuration
The configuration is used to provide to the generator the location os the templates and how the will be generate and also the output's destination.

##### Global
| Property  | Description |
|:----------|:----------|
| combined | The option is a flag (true or false) and allows you to decide if the output will be combined in one file or will be located in different files defined on the Templates configuration.
| srcDir | Defines the location of the templates. If a relative path is given then the absolute path will be resolved using the config's file absolute's path.
| dstDir | Defines the output directory. If a relative path is given then the absolute path will be resolved using the config's file absolute's path.
| output | Destination of the combined output. In case the combined property is true, this property should be set.
| templates | Configuration for each template. See section Templates below.

##### Templates
| Property  | Description |
|:----------|:----------|
| src | Template's location.
| dst | Templates's output destination. This property will be omitted when the combined property is set to true.

This example shows an configuration for a single template -> output:
```yaml
combined: false
srcDir: "__TEMPLATES_SOURCE_DIRECTORY__"
dstDir: "__OUTPUT_DESTINATION_DIRECTORY__"
templates:
  -
    src: "simple.tmpl"
    dst: "simple.html"
```

This example shows an configuration for a multiple templates -> output.
```yaml
combined: true
srcDir: "__TEMPLATES_SOURCE_DIRECTORY__"
dstDir: "__OUTPUT_DESTINATION_DIRECTORY__"
output: "__OUTPUT_FILENAME__"
templates:
  -
    src: "base.tmpl"
  -
    src: "title.tmpl"
  -
    src: "version.tmpl"
  -
    src: "baseUri.tmpl"
  -
    src: "protocols.tmpl"
  -
    src: "mediaTypes.tmpl"
```
To see how the configuration looks like, you can see it for the `try-it-out` located in [try-it-out/templates/config.yaml](try-it-out/templates/config.yaml).

## Usage

HTML from a RAML's specification:

```
$ rubberdoc generate --spec=API.raml --config=config.yml
```

HTML from a Blueprint's specification:

```
$ rubberdoc generate --spec=API.apib --config=config.yml
```

> Note: Check [Configuration](#configuration) section to how to build your config.yml file.

## Help

As usual, you can also see all supported flags by passing `-h`:

```
NAME:
   RubberDoc - A documentation generator for RAML and Blueprint

USAGE:
   rubberdoc [global options] command [command options] [arguments...]

VERSION:
   v0.1-alpha-2

DESCRIPTION:
   Documentation's generator for RAML and Blueprint.

COMMANDS:
     generate  This command receives a configuration file and a specification file written in RAML or Blueprint.
     help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug, -d    Enable debug logging
   --help, -h     show help
   --version, -v  print the version
```

## Examples

Using RAML's specification:

```
$ rubberdoc generate --spec=examples/spec/raml/simple.raml --config=try-it-out/templates/config.yaml
```

Using Blueprint's specification:

```
$ rubberdoc generate --spec=examples/spec/blueprint/simple.apib --config=try-it-out/templates/config.yaml
```

## Contribute

When you want to contribute to the development, follow the [contribution guidelines](contributing.md).

## CREDITS

- [Jumpscale/go-raml](https://github.com/Jumpscale/go-raml) for the raml parser.
- [apiaryio/drafter](https://github.com/apiaryio/drafter) for the Drafter library used by blueprint parser.
- [Sirupsen/logrus](https://github.com/Sirupsen/logrus) for the logging library.
- [urfave/cli](https://github.com/urfave/cli) for the cli library.
- [stretchr/testify](https://github.com/stretchr/testify) for the Assert package.
