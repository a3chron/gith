# gith

A command-line git helper

> [!WARNING]
> In development, would not recommend using yet

## Installation

Install directly using Go

```bash
go install github.com/kurtschambach/gith@latest
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
