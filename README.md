# 🍯 Runny: for running things

Makefiles are for boomers. The future is Runny.

## Features

* ❤️ Simple YAML syntax (inspired by Github Actions)
* 🪄 Full schema validaton == autocomplete in your favourite code editor
* 🧱 Build workflows through composition with `needs` and `then`
* 🏃‍♂️ Run steps conditionally with `if`
* 🙈 Hide helper commands from the command list with `internal`

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

Commands marked `internal: true` are hidden from the command list unless `--verbose` is given. This is useful for helper commands that are only used as dependencies:

```yaml
commands:
  install-uv:
    internal: true
    if: "! command -v uv"
    run: pip install uv
  pip-sync:
    needs:
      - install-uv
    run: uv pip sync requirements.txt
```

Then run commands with runny:

```command
runny pip-install ruff
```

## Examples

### Go

```yaml
commands:
  clean:
    run: |
      go clean ./...
      rm -rf dist
  install-goreleaser:
    if: "! command -v goreleaser"
    run: brew install goreleaser/tap/goreleaser
  release:
    needs:
      - clean
      - install-goreleaser
    run: |
      export GITHUB_TOKEN=$(gh auth token)
      goreleaser
  generate:
    run: go generate ./...
  test:
    run: go test ./...
  test-coverage:
    run: go test -coverprofile=c.out ./... && go tool cover -func="c.out"
  test-coverage-html:
    run: go test -coverprofile=c.out ./... && go tool cover -html="c.out"
```

### Python

```yaml
commands:
  update-requirements:
    run: pip freeze > requirements.txt
  pip-install:
    argnames:
      - packagespec
    run: pip install $packagespec
    then:
      - update-requirements
```

### Docker Compose

Docker Compose has good command-line completion already. But using runny, you can add entries for just the commands you use regularly, then get an uncluttered list of options when you tab-complete.

```yaml
commands:
  up:
    run: docker compose up -d
  down:
    run: docker compose down
  build-and-up:
    run: docker compose up --build -d
  logs:
    argnames:
      - service
    run: docker compose logs $service
  shell:
    argnames:
      - service
    run: docker compose exec $service sh
```
