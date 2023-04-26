#!/usr/bin/env bash

set -e

echo "== Recompiling latest yarn =="

ysc compile -n game -o internal/gamedata/yarn/bin \
  internal/gamedata/yarn/**.yarn

printf "== Running game ==\n\n"

go run cmd/game/main.go