<h3 align="center">
	<img src="https://raw.githubusercontent.com/catppuccin/catppuccin/main/assets/logos/exports/1544x1544_circle.png" width="100" alt="Logo"/><br/>
	Gith (Catppuccin Theme)
</h3>

<p align="center">
	<a href="https://github.com/kurtschambach/gith/releases/latest"><img src="https://img.shields.io/github/v/release/kurtschambach/gith?colorA=363a4f&colorB=b7bdf8&style=for-the-badge"></a>
	<a href="https://github.com/kurtschambach/gith/issues"><img src="https://img.shields.io/github/issues/kurtschambach/gith?colorA=363a4f&colorB=f5a97f&style=for-the-badge"></a>
	<a href="https://github.com/kurtschambach/gith/actions/workflows/lint.yaml"><img src="https://img.shields.io/github/check-runs/kurtschambach/gith/main?colorA=363a4f&colorB=a6da95&style=for-the-badge"></a>
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

For the terminal customization / starship config, check out my [ubuntu customization blog article](https://a3chron.vercel.app/blog/ubuntu-setup).  
This is the full setup, for only starship scroll down to the starship section.

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

- [ ] add config file (accent & flavor)

- [ ] add init repo (&add remote?) with presets for some frameworks, gitignores etc

- [ ] get latest tag in utils version, in add tag show current latest tag

- [ ] simplify level logic to lewer risk of wrong level outputs -> with this add level var to show current active branch with accent color

- [ ] add loading states when fetching sth

- [ ] prettier version (& help) output

- [ ] list branches

<br />

<p align="center">
 <a href="https://github.com/kurtschambach/gith/LICENSE"><img src="https://img.shields.io/github/license/kurtschambach/gith?colorA=363a4f&colorB=b7bdf8&style=for-the-badge"></a>
</p>
