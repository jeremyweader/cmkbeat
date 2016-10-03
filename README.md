# Cmkbeat

Welcome to Cmkbeat.

CMKBeat is used to retrieve information from Check_MK's livestatus and send it to Elasticsearch / Logstash / etc.

WARNING: This repository is still under active development, and there is no guarantee that it will run (or build) properly!

## Getting Started with Cmkbeat

### Requirements

* [Golang](https://golang.org/dl/) 1.7
* [Glide](https://glide.sh) 0.12

### Build

CMKBeat uses Glide for dependency management. To install glide, see https://github.com/Masterminds/glide

or (in most cases) run "go get github.com/Masterminds/glide.

To install all dependencies, simply run

"glide up"

in the cmkbeat directory, and to build the cmk binary and default configuration files, either

"make"
-or-
"go build"

will do the trick. Once you have built the cmkbeat binary, you can simply copy the executable and configuration files into your desired directories, and there is a systemv style init script in the 'services directory.
You can also just run

"make install"

to install everything to the default locations.

To start cmkbeat manually, run

./cmkbeat -c /path/to/cmkbeat.yml 
