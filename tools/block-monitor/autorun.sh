#!/bin/bash

while ! ./block-monitor -events-address0=localhost:31315 -events-address1=localhost:31316 -events-address2=localhost:31317 -events-address3=localhost:31318 -rest-address0=localhost:5000 -rest-address1=localhost:5001 -rest-address2=localhost:5002 -rest-address3=localhost:5003

do
  sleep 1
  # echo "Restarting program..."
done
