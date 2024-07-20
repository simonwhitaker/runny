#!/usr/bin/env bash

complete -W "$(runny | cut -f1)" runny
