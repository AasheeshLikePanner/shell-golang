# GoShell 🐚  
A small shell written in Go — because why not?

This started as a weekend project to figure out how shells actually work. I ended up building something that can run real commands, has some built-in features, and is honestly kinda fun to use.

### Why I built it  
I was curious about how shells handle input, spawn processes, and interact with the system. Writing it in Go made things super clean, and I learned a lot along the way. It’s not meant to replace your terminal — just something cool to hack on and learn from.

### How to run it  
If you’ve got Go installed, just do:

```bash
go run main.go
```

Boom. You're in the shell.

### What it supports  
I added support for a bunch of commands — some built-in, some passed to the system:

- `exit`, `ping`, `pwd`, `ls`, `cd`, `echo`
- `mkdir`, `touch`, `rm`, `cat`, `clear`, `date`
- `curl`, `ip`, `vim` 

It’s got enough to feel usable, but not overloaded with stuff.

### Things I might add later  
- [ ] Piping (`|`) and redirection  
- [ ] Command history (with arrow keys)  
- [X] Tab completion  
- [x] Maybe themes or aliases if I get fancy
