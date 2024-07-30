# 🍯 Runny: for running things

Makefiles are for boomers. The future is Runny.

## Features

* ❤️ Simple YAML syntax (inspired by Github Actions)
* 🪄 Full schema validaton == autocomplete in your favourite code editor
* 🧱 Build workflows through composition with `needs` and `then`
* 🏃‍♂️ Run steps conditionally with `if`

## Installation

```command
brew install simonwhitaker/tap/runny
```

## Usage

Create a `.runny.yaml` file:

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
  pip-install:
    argnames:
      - packagespec
    run: echo $packagespec >> requirements.in
    then:
      - pip-compile-and-sync
```

Then run commands with runny:

```command
runny pip-install ruff
```

## Examples

Have a look in the [examples folder](./examples/) for examples of how you might use Runny with various languages and frameworks.
