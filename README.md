# 🍯 Runny: for running things

Makefiles are for boomers. The future is Runny.

## Features

* ❤️ Simple YAML syntax inspired by Github Actions
* 🪄 Full schema validaton == autocomplete in your favourite code editor
* 🧱 Build workflows through composition with `needs`
* 🏃‍♂️ Skip the steps you don't need to run with `if`

## Installation

```command
brew install simonwhitaker/tap/runny
```

## Usage

Create a .runny.yaml:

```yaml
shell: /bin/bash
commands:
  install-uv:
    if: "! command -v uv"
    run: pip install uv
  pip-sync:
    needs:
      - install-uv
    run: uv pip sync requirements.txt
  pip-compile-and-sync:
    needs:
      - install-uv
    run: |
      uv pip compile requirements.in -o requirements.txt
      uv pip sync requirements.txt

```

Then run commands with runny:

```command
runny pip-compile-and-sync
```
