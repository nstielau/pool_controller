#!/bin/sh
./echoserver.py & 
python ./pool_controller.py run
