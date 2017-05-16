#!/bin/sh

#
# Script used with watchngo to recompile and restart
# a Go project.
#
# If script is run with an argument it runs a cleanup of
# temporary files (executable, PID file) and kills running
# executable.
#

OUTPUT=goweb_exe # name of temporary executable
PIDFILE=my.pid # file where PID of executable is stored
PROJ=github.com/spagettikod/watchngo/_examples/goweb # project to compile

set -e # fail fast and return exit code
MYPID=""
if [ -f "$PIDFILE" ]; then # read PID from file if file exists
    MYPID=$(cat "$PIDFILE")
fi
kill "$MYPID" || true # kill executable with found PID, supress exit code

# if there is no parameter given to this script compile, otherwise clean up.
if [ $# -eq 0 ]; then
    go build -o "$OUTPUT" "$PROJ" # build project into the custom executable
    nohup ./"$OUTPUT" &>/dev/null & # start executable
    echo $! > "$PIDFILE" # save executable PID to file
else
    rm "$PIDFILE" # clean up PID file
    rm "$OUTPUT" # clean up temporary executable
fi