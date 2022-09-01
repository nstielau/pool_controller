# Python Base Image from https://hub.docker.com/r/arm32v7/python/
#FROM balenalib/raspberry-pi-python:3.10-build
FROM balenalib/raspberrypi4-64-python:3.10-bookworm-build


ENV VERSION=1.0.0
RUN python --version

# Copy the Python files
COPY pool_controller.py ./
COPY echoserver.py ./
COPY run.sh ./
RUN mkdir -p gateway
COPY gateway/ ./gateway/
RUN mkdir -p gateway
COPY screenlogic/ ./screenlogic/

# Intall the rpi.gpio python module
RUN pip install --upgrade pip
RUN pip install --no-cache-dir rpi.gpio
RUN pip install --no-cache-dir bottle
RUN pip install --no-cache-dir ask-sdk
RUN pip install --no-cache-dir ask-sdk-webservice-support

# Trigger Python script
CMD ["bash", "./run.sh"]
