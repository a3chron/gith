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
- [What is and what will be (features)](#what-is-and-what-will-be)
- [Contributing](#contributing)
- [Thanks](#thanks)

## Installation

### Go

If you already have Go installed, or plan to do so (easy to get new updates):
[Install via Go](https://gith.featurebase.app/help/articles/6375101-installation-via-go)

> [!TIP]
> If you installed gith via `go install`,
> you can run `gith update` to get the latest version of gith.

### Binaries

If you don't wish to use Go, you can install gith via Binaries:
[Install via Binaries](https://gith.featurebase.app/help/articles/2452108-installation-via-binaries)

> For updates, you have to download and install the new binaries manually again

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
gith tag
```

which will start interactive mode, but already at a point to select the tag, i.e. patch / minor / major / manual input.

For more info run `gith help` or check out the [help articles](https://gith.featurebase.app/help).

Gith tries to use intuitive, natural language commands,
combined with the usual git commands, for example `gith tag` or `gith update remote url`.

You can also get completions for fish, base or zsh: [Completions](https://gith.featurebase.app/help/articles/8096273)

## Customization

You can set your preferred flavor and accent in the Options.  
Just run `gith` and select "Options".

For more info check out the [help articles](https://gith.featurebase.app/help).

## What is and what will be

**You can view the Roadmap / upvote Feature Requests and Bugs at [Featurebase](https://gith.featurebase.app/).**

- [x] Branch

  - [x] Switch Branch _-- supports quick select --_

  - [x] List Branches _-- supports quick select --_

  - [x] Delete Branch _-- supports quick select --_

  - [x] Create Branch

- [x] Status

  - [x] View working tree status (Modified, Added, Deleted, Untracked files) _-- supports quick select --_

- [ ] Commit

  - [x] Undo Last Commit _-- supports quick select --_

  - [x] Commit staged changes _-- supports quick select --_

  - [x] Commit all changes _-- supports quick select --_

  - [ ] Amend last commit

- [x] Tag

  - [x] List Tags

  - [x] Remove Tag

  - [x] Push Tag _-- supports quick select --_

  - [x] Add Tag _-- supports quick select --_

- [ ] Remote

  - [x] List Remotes

  - [x] Add Remote _-- supports quick select --_

  - [x] Remove Remote

  - [ ] Update Remote url

  - [ ] Push with remotes

  - [ ] Pull with remotes

- [ ] Changes

  - [ ] View diff of changes

  - [ ] Stage individual files

  - [ ] Unstage individual files

- [x] Options

  - [x] Change UI flavor

  - [x] Change UI accent color

  - [x] Change fetch behaviour on Init

## Contributing

Contributions are welcome, please use [conventional commits](https://www.conventionalcommits.org/) for a constant commit message style.
If you reaaally struggle with conventinal commits, check out [Meteor](https://github.com/stefanlogue/meteor), or just use gith's built in commit feature.  
For feature requests or possible improvements please create an issue at [Featurebase](https://gith.featurebase.app/).

For local development, fork & clone the repo, then build using:

```bash
go build -o gith
```

Run with:

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
 <a href="https://github.com/a3chron/gith/blob/main/LICENSE"><img src="https://img.shields.io/github/license/a3chron/gith?colorA=363a4f&colorB=b7bdf8&style=for-the-badge"></a>
</p>
