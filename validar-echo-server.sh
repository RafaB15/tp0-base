#!/bin/bash

mensaje="Probando el echo server"

respuesta=$(docker run --rm --network=tp0_testing_net busybox \
  sh -c "echo '$mensaje' | nc server 12345" 2>/dev/null)

if [ $? -eq 0 ] && [ "$respuesta" = "$mensaje" ]; then
  echo "action: test_echo_server | result: success"
else
  echo "action: test_echo_server | result: fail"
fi
