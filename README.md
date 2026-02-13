# kra
# Development (for contributors / agents)

- Start here: `docs/START_HERE.md`
- Backlog: `docs/backlog/README.md`
- Specs: `docs/spec/`

## Quick Start

```sh
# 1) initialize root explicitly
kra init --root ~/kra

# 2) check current context
kra context current

# 3) create and use workspace
kra ws create --no-prompt TASK-1234
kra ws --id TASK-1234 --act go
```
