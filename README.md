# About
Basic utility to monitor a directory and run a command when changes
are detected.

I use it to recompile Go web services projects on changes. I also use it to precompile Handlebars templates when developing web applications.

It's tested on MacOS but should work on other OS's as well.

## Installation
Download and extract the latest version for MacOS with this command:

```curl -LO https://github.com/spagettikod/watchngo/releases/download/1.0.0/watchngo1.0.0-macos.tar.gz && tar -xvf watchngo1.0.0-macos.tar.gz && rm watchngo1.0.0-macos.tar.gz```

## Usage
```watchngo goweb goweb/run.sh```

The above example monitors the directory ```goweb```, if files or folders are modified, deleted or added to ```goweb``` the script ```run.sh``` will execute.

* watchngo checks for changes every second
* each check only triggers one run of the script eventhough more than one file or directory is modified
* stdout and stderr from the script will be passed on to stdout or stderr by watchngo
* use the ```-v``` flag while developing your script to see what watchngo is doing
* make sure you script doesn't change files within the directory you are watching.

See ```_examples``` folder for working examples recompiling and restarting a Go project when
files change.
