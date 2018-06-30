#!/bin/bash

while true
do
        container_id="$(docker run -d ubuntu:latest /bin/bash -e -c "echo -n \"foo\"")"
        docker logs -f "${container_id}"
done
