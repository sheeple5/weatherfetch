# Weatherfetch
weatherfetch is a play off of the popular command "neofetch" (now known as "fastfetch" I believe). In trying to learn Golang, I created this project to pull 
weather information from an API and print it like neofetch. This currently does not require an API account or key as it utilizes the preview API from
[Visual Crossing](https://www.visualcrossing.com/weather-api/?gad_source=1&gad_campaignid=21410448611&gclid=Cj0KCQiAiqDJBhCXARIsABk2kSmj_X0ARYhU_7p1wvIFd8-z_nQrPQ0bQnIjyeFkJgAYNylRIIXqmsMaAnDMEALw_wcB).

## Building/Installing
Run the following command to build the binary:

`go build weatherfetch.go`

You can also use the makefile to build:

`make build`

To "install" (i.e. build and put it into a bin path), use the following to put it into your home's local bin:

`make install`

You should now be able to run the binary directly with `./weatherfetch` or use `weatherfetch` directly if "install" is chosen. 

## Usage
When executing the binary, you will be prompted to enter a location or address. The API service will do its best to identify the location and return the current
weather information in a pretty display. 

If desired, a config file can be written at `~/.weather_config` that contains your address. If this config file
is used, the address will be pulled directly from the file instead of prompting for user input. Use the following format for your configuration file:

`address = <chosen location>`
