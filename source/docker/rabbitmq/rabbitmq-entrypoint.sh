#!/bin/sh
export RABBITMQ_DEFAULT_PASS=$(cat /run/secrets/rabbitmq_pass)
export RABBITMQ_ERLANG_COOKIE=$(cat /run/secrets/rabbitmq_erlang_cookie)
exec docker-entrypoint.sh "$@"
