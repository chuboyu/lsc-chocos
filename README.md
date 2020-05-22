# lsc-chocos

This is a wrapping/simulation module for IOT data collection @ LSC
currently under construction.

## Quick Start

First you will need a running Mainflux service, an opensource IOT platform.

You need to modify the `config files` under `configs` to add the host to `BaseUrl`. 

Running this program will simply create 1 thing and 1 channel on the server.

Local go routines are also created, and they start generating and uploading messages.

The process will be hanging with no output.

### Command Line

You will need to specify the server certificate file path in the config files.

After that simply run

```
    > go run main.go <config_file_path> -f <confile_file_path> -cacert <crt_file_path>
```

There is a config file ready in `configs/config.json`

### Docker

It will use `configs/config.json` and `ssl/mainflux-server.crt`, so you may need to modify/replace them accordingly.

Please see the Dockerfile for more information.

To start simply run

```
> docker build . -t <image_name>
> docker run <image_name>
```