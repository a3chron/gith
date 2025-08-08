<div align="center">
<table>
<tr>
<td>
<pre>
         ██┐         ██┐        
         └─┘   ██┐   ██│        
 ██████┐ ██┐ ██████┐ ██████┐    
 ██┌─██│ ██│ └─██┌─┘ ██┌─██│    
 ██│ ██│ ██│   ██│   ██│ ██│    
 ██████│ ██│   ████┐ ██│ ██│    
 └───██│ └─┘   └───┘ └─┘ └─┘    
     ██│  A TUI git helper      
 ██████│  with catppuccin theme,
 └─────┘  written in Go         
</pre>
</td>
<td>
<img src="https://raw.githubusercontent.com/catppuccin/catppuccin/main/assets/logos/exports/1544x1544_circle.png" width="100" alt="Logo"/>
</td>
</tr>
</table>
</div>

<p align="center">
	<a href="https://github.com/a3chron/gith/releases/latest"><img src="https://img.shields.io/github/v/release/a3chron/gith?colorA=363a4f&colorB=b7bdf8&style=for-the-badge"></a>
	<a href="https://github.com/a3chron/gith/issues"><img src="https://img.shields.io/github/issues/a3chron/gith?colorA=363a4f&colorB=f5a97f&style=for-the-badge"></a>
	<a href="https://github.com/a3chron/gith/actions/workflows/lint.yaml"><img src="https://img.shields.io/github/check-runs/a3chron/gith/main?colorA=363a4f&colorB=a6da95&style=for-the-badge"></a>
</p>

> [!WARNING]
> Still in development

![](/assets/peek-usage-preview.gif)

<details>
<summary>Catppuccin flavors preview</summary>

![](/assets/preview-flavors.webp)

</details>

<details>
<summary>More Images</summary>

![](/assets/preview-actions.png)
![](/assets/preview-tags.png)
![](/assets/preview-status.png)

</details>

For the terminal customization / starship config, check out my [ubuntu customization blog article](https://a3chron.vercel.app/blog/ubuntu-setup).  
This is the full setup, for only starship scroll down to the starship section.

## Table of Contents

- [Installation](#installation)
- [Customization](#customization)
- [Usage in scripts](#usage-in-scripts)
- [Troubleshooting](#troubleshooting)
- [What is and what will be (features)](#what-is-and-what-will-be)
- [Completions](#completions)
- [Contributing](#contributing)
- [Thanks](#thanks)
- [TODOs](#todos)

## Installation

> [!NOTE]
> If you don't have Go installed, check out the [gith project page](https://a3chron.vercel.app/projects/gith),
> Go installation is included there.

Install directly using Go

```bash
go install github.com/a3chron/gith@latest
```

or using one of the [pre-built releases](https://github.com/a3chron/gith/releases/latest).

After the installation finished, just run:

```bash
gith
```

## Customization

You can set your preferred flavor and accent in the Options.  
Just run `gith` and select "Options".

When running gith the first time, a config file storing your settings will be created at
`XDG_CONFIG_HOME/gith/config.json` if `XDG_CONFIG_HOME` is set,
otherwise at `~/.config/gith/config.json`.

You can also manually edit the config file, although editing with gith ensures that no invalid configurations are used.

**If you have anything you'd like to configure in the settings or options, don't hesitate to open an issue.**

## Usage in scripts

You are able to easily change flavor and accent of gith in scripts:

```bash
gith config update --flavor=Latte --accent=Red
```

Useful for example if you want a script to switch between
light and dark mode in all your catppuccin themed apps.

> [!NOTE]
> Run `gith config help` for more info

## Troubleshooting

**I update gith via `go install github.com/a3chron/gith@latest`, but nothing changes / version stays the same**

> [!NOTE]
> If you get any output when running `go install github.com/a3chron/gith@latest`,
> you probably don't have this issue.  
> Consider opening an issue.

Sometimes it takes some time for the go proxy server to recognize a new release,
so it is possible that the latest release for the proxy server is still the old one.

In that case, just request a lookup for the specific version, for example `v0.6.0` if this is the latest release:

```bash
go install github.com/a3chron/gith@v0.6.0
```

You should now see the output:

```bash
go: downloading github.com/a3chron/gith v0.6.0
```

> [!NOTE]
> You can check for the latest release on [github](https://github.com/a3chron/gith/releases/latest)
> or by simply running `gith version check`

## What is and what will be

- [x] Branch

  - [x] Switch Branch (supports checking out remote branches locally)

  - [x] List Branches

  - [x] Delete Branch (with fallback to force delete)

  - [x] Create Branch

- [x] Status

  - [x] View working tree status (Modified, Added, Deleted, Untracked files)

- [ ] Commit

  - [x] Undo Last Commit

  - [ ] Commit staged changes

  - [ ] Commit all changes

  - [ ] Amend last commit

- [x] Tag

  - [x] List Tags (shows 10 latest, sorted by semantic versioning)

  - [x] Remove Tag

  - [x] Push Tag (prompts to confirm pusing the latest tag)

  - [x] Add Tag

- [ ] Remote

  - [x] List Remotes

  - [ ] Add Remote

  - [x] Remove Remote

  - [ ] Update Remote url

  - [ ] Push to remote

  - [ ] Pull from remote

- [ ] Changes

  - [ ] View diff of changes

  - [ ] Stage/unstage individual files

- [x] Options

  - [x] Change UI flavor

  - [x] Change UI accent color

## Completions

For package installation the completions will be installed automatically.  
When installing via `go install` you can get completions for the few commands there are by running the following command:

**fish**

```bash
gith completion fish > ~/.config/fish/completions/gith.fish
```

**bash**

```bash
gith completion bash > ~/.local/share/bash-completion/completions/gith
```

**Zsh**

```bash
gith completion zsh > ~/.local/share/zsh/site-functions/_gith
```

## Contributing

Contributions are welcome, but please remember to use [conventional commits](https://www.conventionalcommits.org/) (at least when naming the PR).
For features or possible improvements please create an issue.

Build using:

```bash
go build -o gith
```

Then run with:

```bash
./gith
```

## Thanks

- [BubbleTea](https://github.com/charmbracelet/bubbletea) for making it possible to build this
- [opencode](https://github.com/opencode-ai/opencode) for the design inspiration (I copied their starting dialog 1:1)

## TODOs

This is mostly a list for me, but here so you can see what you don't need to create an issue for.
Also good starting point if you'd want to contribute, but please create an issue stating that you started working on it,
so we don't do the same thing twice.

- [ ] add loading states when fetching sth

- [ ] improve remove tag only show n latest tags wiht "pagination" instead of optional show all

- [ ] add confirmation for things like force branch delete

- [ ] "remove remote" add dimmed remote url in selection -> for all add Label with name & optional description

- [ ] clarify "out, err := git.GetRemotes(); m.Err = out" logic, siwth these two

- [ ] "Push Tag" Load All functionality

- [ ] bug: double space before success end (this could be a feature?)

- [ ] delete tag deletes localy not remote, for remote: git push origin --delete <tagname>

<br />

<p align="center">
 <a href="https://github.com/a3chron/gith/LICENSE"><img src="https://img.shields.io/github/license/a3chron/gith?colorA=363a4f&colorB=b7bdf8&style=for-the-badge"></a>
</p>
