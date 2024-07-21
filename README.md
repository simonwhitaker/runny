# Runny: for running things

<img src="./assets/runny-logo.png" width="100" alt="Runny logo" />

Stop using Makefiles to run things. It's 2024. Use Runny.

## Installation

```command
brew install simonwhitaker/tap/runny
```

## Usage

Create a .runny.yaml:

```yaml
commands:
  install-uv:
    if: "! command -v uv"
    command: pip install uv
  pip-sync:
    needs: install-uv
    command: uv pip sync requirements.txt
  pip-compile-and-sync:
    needs: install-uv
    command: |
      uv pip compile requirements.in -o requirements.txt
      uv pip sync requirements.txt
```

Then run commands with runny:

```command
runny pip-compile-and-sync
```
