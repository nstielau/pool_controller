# Python Base Image from https://hub.docker.com/r/arm32v7/python/
FROM balenalib/raspberry-pi-python:3.10-build

ENV VERSION=1.0.0
RUN python --version

# Copy the Python files
COPY pool_controller.py ./
RUN mkdir -p gateway
COPY gateway/ ./gateway/
RUN mkdir -p gateway
COPY screenlogic/ ./screenlogic/

# Intall the rpi.gpio python module
RUN pip install --no-cache-dir rpi.gpio
RUN pip install --no-cache-dir bottle

# Trigger Python script
CMD ["python", "./pool_controller.py", "run"]
