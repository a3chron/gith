# gith

A command-line git helper

> [!WARNING]
> Still in development, rather big changes possible

![](/assets/preview.png)

## Installation

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

built using [BubbleTea](https://github.com/charmbracelet/bubbletea)

## TODO

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

-> Make a errOut (Output) and errMsg (Switch failed)

---

-> use the stepComplete for sth or remove

-> add config file (accent)

-> add help command, only show logo there

-> add undo last commit

-> add init repo (&add remote?) with presets for some frameworks, gitignores etc
