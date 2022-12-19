# GO Healthchecker

## Description

This simple golang application acts as a one-to-many healthcheck probe. It accepts `GET` and `HEAD` HTTP methods on `/healthcheck` address.

## Configuration

Set below environment variables

Variable Name|Description
|---|---|
MONITORED_URLS|semicolon separated list of urls to monitor. Make sure that the URL will not return 301 redirect as this will fail healthcheck. Will only accept 200
SERVER_PORT|Port on which server is listening for incoming connections. DEFAULT: 8080
