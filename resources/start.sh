#!/bin/bash

set -e

pppd call mint &

trap : TERM INT; sleep infinity & wait

