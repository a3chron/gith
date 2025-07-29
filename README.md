<h3 align="center">
	<img src="https://raw.githubusercontent.com/catppuccin/catppuccin/main/assets/logos/exports/1544x1544_circle.png" width="100" alt="Logo"/><br/>
	Gith (Catppuccin Theme)
</h3>

<p align="center">
	<a href="https://github.com/a3chron/gith/releases/latest"><img src="https://img.shields.io/github/v/release/a3chron/gith?colorA=363a4f&colorB=b7bdf8&style=for-the-badge"></a>
	<a href="https://github.com/a3chron/gith/issues"><img src="https://img.shields.io/github/issues/a3chron/gith?colorA=363a4f&colorB=f5a97f&style=for-the-badge"></a>
	<a href="https://github.com/a3chron/gith/actions/workflows/lint.yaml"><img src="https://img.shields.io/github/check-runs/a3chron/gith/main?colorA=363a4f&colorB=a6da95&style=for-the-badge"></a>
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
go install github.com/a3chron/gith@latest
```

## What is and what will work

- [x] Branch

  - [x] Switch Branch (supports checking out remote branches locally)

  - [x] List Branches

  - [x] Delete Branch (with fallback to force delete)

  - [ ] Create Branch

- [x] Status

  - [x] View working tree status (Modified, Added, Deleted, Untracked files)

- [x] Commit

  - [x] Undo Last Commit

  - [ ] Commit staged changes

  - [ ] Commit all changes

  - [ ] Amend last commit

- [x] Tag

  - [x] List Tags (shows 10 latest)

  - [x] Remove Tag (from a list of all local tags)

  - [x] Push Tag (prompts to confirm pusing the latest tag)

  - [ ] Add Tag

- [x] Remote

  - [x] List Remotes

  - [ ] Add Remote

  - [ ] Remove Remote

  - [ ] Update Remote url

  - [ ] Push to remote

  - [ ] Pull from remote

- [ ] Changes

  - [ ] View diff of changes

  - [ ] Stage/unstage individual files

- [ ] Options

  - [ ] Change UI flavor

  - [ ] Change UI accent color

## Local Development

Build using:

```bash
go build -o gith
```

Then run with:

```bash
./gith
```

## Customization

I will add a config file & customization options soon.  
If you can't wait for the other catppuccin accents / flavors,
you can clone the repo, change everything you need in `/internal/ui/styles.go`.

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

- [ ] add init repo (&add remote?) with presets for some frameworks, gitignores etc -> move to seperate project gith-init

- [ ] get latest tag in utils version, in add tag show current latest tag

- [ ] add loading states when fetching sth

- [ ] prettier version (& help) output

- [ ] list branches

- [ ] fix double spaces between some things & others not

- [ ] remove tag only show n latest tags & last option "Show older tags" which shows n older tags

<br />

<p align="center">
 <a href="https://github.com/a3chron/gith/LICENSE"><img src="https://img.shields.io/github/license/a3chron/gith?colorA=363a4f&colorB=b7bdf8&style=for-the-badge"></a>
</p>
