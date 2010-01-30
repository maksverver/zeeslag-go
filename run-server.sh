#!/bin/sh

SOCK=server.sock
LOG=server.log
BIN=./server

HOST=
PORT=14000
ROOT=/player

if [ -e "${SOCK}" ]; then echo "${SOCK} exists!"; exit 1; fi
if [ ! -x "${BIN}" ]; then echo "${BIN} does is not executable!"; exit 1; fi
dtach -c "${SOCK}" /bin/sh -c "'${BIN}' -h '${HOST}' -p '${PORT}' -r '${ROOT}' | tee -a '${LOG}'"
