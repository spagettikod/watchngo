# About
Basic utility to monitor a directory and run a command when changes
are detected.

I use it to recompile on changes when developing Go web services. Another
usage is to precompile Handlebars templates when developing web applications.

It's tested on MacOS but should work on other OS's as well.

## Installation
Download and extract the latest version for MacOS with this command:

```curl -LO https://github.com/spagettikod/watchngo/releases/download/1.0.0/watchngo1.0.0-macos.tar.gz && tar -xvf watchngo1.0.0-macos.tar.gz && rm watchngo1.0.0-macos.tar.gz```

## Usage
```watchngo goweb goweb/run.s```

See ```_examples```folder for working examples recompiling and restarting a Go code when
files change.
