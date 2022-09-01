# Python Base Image from https://hub.docker.com/r/arm32v7/python/
FROM arm32v7/python:2.7.13-jessie

ENV VERSION=1.0.0
RUN python --version


# Copy the Python Script to blink LED
COPY pool_controller.py ./
RUN mkdir -p gateway
COPY gateway/ ./gateway/
RUN mkdir -p gateway
COPY screenlogic/ ./screenlogic/
RUN ls -la ./
RUN ls -la ./screenlogic
RUN ls -la ./gateway

# Intall the rpi.gpio python module
RUN pip install --no-cache-dir rpi.gpio

# Trigger Python script
CMD ["python", "./pool_controller.py", "run"]
