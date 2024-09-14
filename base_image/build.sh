#!/bin/bash

cd proxy && go build -o port_proxy && cd ..
docker build -t piko-proxy .