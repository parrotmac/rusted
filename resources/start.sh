#!/bin/bash

set -e

pppd call twilio &

trap : TERM INT; sleep infinity & wait

