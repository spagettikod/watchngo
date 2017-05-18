# Go web service example
I use this watchngo setup when developing web services in Go.

## Usage
watchngo goweb goweb/run.sh

Make sure you don't run the script in the goweb directory since
temporary files are created in your current directory. These
files would trigger watchngo into an intfinite loop.