package main

/*
   Listen on port 9191
   Reads toml file that has exporters with name and URL (http://192.168.30.23:9090/log)
   Has webserver in / listing all entries of the toml file and every entry is an a tag
   When a tag pressed, go to /<service> and show journalctl output like gokrazy does
   Runs in gokrazy so I can go to gokrazy:9191 and view all logs
*/
