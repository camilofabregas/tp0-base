#!/bin/bash

PUERTO=$(awk -F' = ' '/SERVER_PORT/ {print $2}' server/config.ini | tr -d '[:space:]')
ADDRESS=$(awk -F' = ' '/SERVER_IP/ {print $2}' server/config.ini | tr -d '[:space:]')

# Mensaje de prueba
TEST_MSG="hola-echo-server"

RESULTADO=$(docker run --rm --network tp0_testing_net --entrypoint sh busybox -c "echo '$TEST_MSG' | nc $ADDRESS $PUERTO")

if [ "$RESULTADO" = "$TEST_MSG" ]; then
    echo "action: test_echo_server | result: success"
else
    echo "action: test_echo_server | result: fail"
fi