#!/bin/sh

go build -o thingies github.com/alexd765/thingies/backend
zip thingies.zip thingies
