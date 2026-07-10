#!/bin/bash

scriptPos=${0%/*}

COMPOSE_FILE=$scriptPos/../compose_env.yaml

# Detect available container runtime (docker or podman)
if docker --version &> /dev/null; then
  CONTAINER_RUNTIME="docker"
  COMPOSE_CMD="docker compose"
elif podman --version &> /dev/null; then
  CONTAINER_RUNTIME="podman"
  # Podman uses 'podman-compose' or 'podman' with compose
  COMPOSE_CMD="podman compose"
else
  echo "Error: Neither Docker nor Podman is installed. Please install one of them."
  exit 1
fi

function start() {
  echo "Starting ${CONTAINER_RUNTIME} Compose environment..."
  $COMPOSE_CMD -p dynsvg -f $COMPOSE_FILE up -d
}

function stop() {
  echo "Stopping ${CONTAINER_RUNTIME} Compose environment..."
  $COMPOSE_CMD -p dynsvg -f $COMPOSE_FILE down
}

function destroy() {
  echo "Destroying ${CONTAINER_RUNTIME} Compose environment..."
  $COMPOSE_CMD -p dynsvg -f $COMPOSE_FILE down -v
}

case "$1" in
  start)
    start
    ;;
  stop)
    stop
    ;;
  destroy)
    destroy
    ;;
  *)
    echo "Usage: $0 {start|stop|destroy}"
    exit 1
    ;;
esac

exit 0
