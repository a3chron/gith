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

> **Why try gith?**  
> It brings little quality-of-life improvements that can make
> working with git feel smoother and sometimes even a bit quicker.
> For example, the “Add Tag” feature (also shown under “More Images” below the preview GIF)
> offers Patch / Minor / Major besides the Manual Input of tags.  
> Apart from that, it also just looks good ;)

![](/assets/peek-usage-preview.gif)

<details>
<summary>Catppuccin flavors preview</summary>

![](/assets/preview-flavors.webp)

</details>

<details>
<summary>More Images</summary>

![](/assets/preview-actions.png)
![](/assets/preview-status.png)
![](/assets/preview-add-tag.png)

</details>

For the terminal customization / starship config, check out my [ubuntu customization blog article](https://a3chron.vercel.app/blog/ubuntu-setup).  
This is the full setup, for only starship scroll down to the starship section.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Customization](#customization)
- [Usage in scripts](#usage-in-scripts)
- [Troubleshooting](#troubleshooting)
- [What is and what will be (features)](#what-is-and-what-will-be)
- [Completions](#completions)
- [Contributing](#contributing)
- [Thanks](#thanks)
- [TODOs](#todos)

## Installation

If you already have Go installed, or plan to do so (easy to get new updates):
[Install via Go](https://gith.featurebase.app/help/articles/6375101-installation-via-go)

Otherwise you can install gith via Binaries:
[Install via Binaries](https://gith.featurebase.app/help/articles/2452108-installation-via-binaries)

## Usage

After the installation finished, just run:

```bash
gith
```

Some commands like adding tags are also accessible via quick select,
check out the [What is and what will be (features)](#what-is-and-what-will-be) section,
to see which commands are supported.

As an example, quick select for "git tag <tag>":

```bash
gith add tag
```

which will start interactive mode, but already at a point to select the tag, i.e. patch / minor / major / manual input.

This approach tries to use intuitive, natural language commands, such as `gith add tag` or `gith update remote url`.

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
> you probably just discovered some other bug.  
> Consider opening an issue at [Featurebase](https://gith.featurebase.app/) or GitHub.

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

**You can view the Roadmap / upvote Feature Requests and Bugs at [Featurebase](https://gith.featurebase.app/).**

- [x] Branch

  - [x] Switch Branch

  - [x] List Branches

  - [x] Delete Branch

  - [x] Create Branch

- [x] Status

  - [x] View working tree status (Modified, Added, Deleted, Untracked files)

- [ ] Commit

  - [x] Undo Last Commit

  - [x] Commit staged changes

  - [x] Commit all changes

  - [ ] Amend last commit

- [x] Tag

  - [x] List Tags

  - [x] Remove Tag

  - [x] Push Tag _-- supports quick select --_

  - [x] Add Tag _-- supports quick select --_

- [ ] Remote

  - [x] List Remotes

  - [x] Add Remote

  - [x] Remove Remote

  - [ ] Update Remote url

  - [ ] Push to remote

  - [ ] Pull from remote

- [ ] Changes

  - [ ] View diff of changes

  - [ ] Stage individual files

  - [ ] Unstage individual files

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

Contributions are welcome, please use [conventional commits](https://www.conventionalcommits.org/) for constant commit message style.
If you reaaally struggle with conventinal commits, check out [Meteor](https://github.com/stefanlogue/meteor), or just use gith's built in commit feature.  
For feature requests or possible improvements please create an issue at [Featurebase](https://gith.featurebase.app/).

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
- [opencode](https://github.com/opencode-ai/opencode) for the design inspiration (copied their starting dialog 1:1)

## Support

If you want to support me, feel free to upvote gith on producthunt, and give some feedback :)
You can also upvote feature requests or bugs that you want fixed at [Featurebase](https://gith.featurebase.app/).

<a href="https://www.producthunt.com/products/gith?embed=true&utm_source=badge-featured&utm_medium=badge&utm_source=badge-gith&#0045;beta" target="_blank"><img src="https://api.producthunt.com/widgets/embed-image/v1/featured.svg?post_id=1004160&theme=neutral&t=1754987210952" alt="Gith&#0032;&#0040;beta&#0041; - A&#0032;Terminal&#0032;UI&#0032;for&#0032;git | Product Hunt" style="width: 250px; height: 54px;" width="250" height="54" /></a>

<br />

<p align="center">
 <a href="https://github.com/a3chron/gith/LICENSE"><img src="https://img.shields.io/github/license/a3chron/gith?colorA=363a4f&colorB=b7bdf8&style=for-the-badge"></a>
</p>
