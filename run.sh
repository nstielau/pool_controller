#!/bin/sh
python ./echoserver.py & 
python ./pool_controller.py run
