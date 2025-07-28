# gith

A command-line git helper with catppuccin theme

> [!WARNING]
> Still in development, rather big changes possible

![](/assets/preview-actions.png)
![](/assets/preview-switch-branch.png)

## Installation

> [!NOTE]
> If you don't have Go installed, check out the [gith project page](https://a3chron.vercel.app/projects/gith),
> Go installation is included there.

Install directly using Go

```bash
go install github.com/kurtschambach/gith@latest
```

## Local Development

Build using:

```bash
go build -o gith
```

Then run with:

```bash
./gith
```

---

built using [BubbleTea](https://github.com/charmbracelet/bubbletea), design heavily based on [opencode](https://github.com/opencode-ai/opencode)

## TODO

```
◆ Select action
 │ Switch Branch
 │
 ■ Select branch
 │ feat/base
 │
 ╰─╌ Switch failed: error: Your local changes to the following files would be overwritten by checkout:
 gith
 main.go
 Please commit your changes or stash them before you switch branches.
 Aborting
```

- [ ] Make a errOut (Output) and errMsg (Switch failed) or errOut in out and err in err?

- [ ] add config file (accent)

- [x] add help command

- [x] add explicit no git repo error

- [x] add undo last commit: git reset --soft HEAD~1

- [ ] add init repo (&add remote?) with presets for some frameworks, gitignores etc

- [x] for gith --version check for new versions

- [ ] investigate version injection not working

- [ ] get latest tag in utils von version, in add tag show current latest tag

- [ ] in utils wirteErrOut & writeOut with currentStep to swotch over steps where to add out in out arrays

- [ ] add loading states when fetching sth

- [ ] improve error handling (e.g. if branches no branches / no tags when listing add info No tags found / for switch / delete branch check if no branches found first)
