# Runny: for running things

<img src="./assets/runny-logo.png" width="100" alt="Runny logo" />

Stop using Makefiles to run things. It's 2024. Use Runny.

## Installation

```command
go install .
```

## Usage

Create a .runny.yaml:

```yaml
commands:
    git-root:
        command: git rev-parse --show-toplevel
    clean:
        command: rm -rf node_modules
    readme-stats:
        command: |
            if [[ -e README.md ]]; then
                wc README.md
            else
                echo "Couldn't find README.md"
            fi
```

Then run commands with runny:

```command
runny clean
```
