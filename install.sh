#!/bin/bash

make -C set/ install && \
  make -C deps/ install &&
  make -C parse/ install &&
  make -C build/ install &&
  make -C main/ install
