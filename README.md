# lsc-chocos

This is a wrapping/simulation module for IOT data collection @ LSC
currently under construction.


First you will need a local version of Mainflux, an opensource IOT platform, on <https://localhost>.

To start simply run

`> go run main.go`

This will simply create 1 thing and 1 channel on the server,

Local go routines are also created for generating and uploading messages.