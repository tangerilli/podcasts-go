podcasts
========

A port of https://github.com/tangerilli/podcasts to Go.

Usage
-----

Compile, then run the podcasts binary.  This will start an HTTP server, by 
default on port 4000 (can be changed with the --port argument).  Takes two
arguments, the local directory path to search for files that look like movies,
and the URL that the application will be available at (e.g. if the app is being
proxied through apache/nginx/etc.., it will be something other than localhost:4000).

Command line flags can also be used to modify the URL that the video files are
actually available for download from, and various bits of podcast metadata (run
with --help for more details).

Once it's running, add the appropriate URL as a podcast in iTunes.  Once it has
synced to your iPhone/iPad, iTunes is no longer necessary, and you can download
videos directly from the Podcasts app on your iDevice.