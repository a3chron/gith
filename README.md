<h3 align="center">
	<img src="https://raw.githubusercontent.com/catppuccin/catppuccin/main/assets/logos/exports/1544x1544_circle.png" width="100" alt="Logo"/><br/>
	Gith (Catppuccin Theme)
</h3>

<p align="center">
	<a href="https://github.com/kurtschambach/gith/releases/latest"><img src="https://img.shields.io/github/v/release/kurtschambach/gith?colorA=363a4f&colorB=b7bdf8&style=for-the-badge"></a>
	<a href="https://github.com/kurtschambach/gith/issues"><img src="https://img.shields.io/github/issues/kurtschambach/gith?colorA=363a4f&colorB=f5a97f&style=for-the-badge"></a>
	<img src="https://img.shields.io/github/check-runs/kurtschambach/gith/main?colorA=363a4f&colorB=a6da95&style=for-the-badge">
</p>

A command-line git helper with catppuccin theme written in Go

> [!WARNING]
> Still in development

![](/assets/peek-usage-preview.gif)

<details>
<summary>More Images</summary>

![](/assets/preview-actions.png)
![](/assets/preview-tags.png)
![](/assets/preview-status.png)

</details>

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

<br />

<p align="center">
 <a href="https://github.com/kurtschambach/gith/LICENSE"><img src="https://img.shields.io/github/license/kurtschambach/gith?colorA=363a4f&colorB=b7bdf8&style=for-the-badge"></a>
</p>
