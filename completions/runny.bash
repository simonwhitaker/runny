#!/usr/bin/env bash

# See https://opensource.com/article/18/3/creating-bash-completion-script
_runny_completions() {
    runny_commands=$(runny | cut -f1)
    COMPREPLY=($(compgen -W "${runny_commands}" "${COMP_WORDS[1]}"))
}

complete -F _runny_completions runny
