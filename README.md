# Cmkbeat

Welcome to Cmkbeat.

CMKBeat is used to retrieve information from Check_MK's livestatus and send it to Elasticsearch / Logstash / etc.

## Getting Started with Cmkbeat

### Requirements

* [Golang](https://golang.org/dl/) 1.7
* [Glide](https://glide.sh) 0.12

### Build

CMKBeat uses Glide for dependency management. To install glide, see https://github.com/Masterminds/glide

or (in most cases) run "go get github.com/Masterminds/glide".

To install all of the dependencies and build the binary, run

"make all"

in the cmkbeat directory.  Once built, you can simply copy the executable and configuration files into your
desired directories, or run

"make install"

to install everything to the default locations. If installing manually, there is a sysv init script in the
system/ directory.

### Run

To start cmkbeat manually, run

./cmkbeat -c /path/to/cmkbeat.yml 
