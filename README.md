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

## What is and what will work

- [ ] Branch

  - [x] Switch Branch (supports checking out remote branches locally)

  - [x] List Branches

  - [x] Delete Branch (with fallback to force delete)

  - [ ] Create Branch

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

- [ ] add loading states when fetching sth

- [ ] improve remove tag only show n latest tags wiht "pagination" instead of optional show all

- [ ] add confirmation for things like force branch delete

- [ ] adjust version info for small terminals

- [ ] "remove remote" add dimmed remote url in selection -> for all add Label with name & optional description

- [ ] clarify "out, err := git.GetRemotes(); m.Err = out" logic, siwth these two

- [ ] "Push Tag" Load All functionality

- [ ] double space before success end

- [ ] delete tag deletes localy not remote, for remote: git push origin --delete <tagname>

<br />

<p align="center">
 <a href="https://github.com/a3chron/gith/LICENSE"><img src="https://img.shields.io/github/license/a3chron/gith?colorA=363a4f&colorB=b7bdf8&style=for-the-badge"></a>
</p>
