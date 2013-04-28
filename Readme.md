# CAPO

Capo is a go script that reads a `commands` file and executes every command in a separate go routine and displays the output in the same terminal session. It works pretty much like foreman and can also read foreman Procfiles for backward compatibility.

Capo doesn't try to do too much. There is no colored output except for the one provided by the commands executed. Standard output, input and error are not affected either. If any command fails, all the others are killed by sending a SIGKILL signal.

### Installation

For now, you'll need a go compiler. You can compile capo by running `go build capo.go` and then move the executable `capo` somewhere in your path.

### Usage

```
> capo # Reads a commands file and executes every command in a separate go routine

> capo file_name # Reads commands from file_name

> capo --foreman # Reads a Procfile in the format used by foreman
```

### Example commands file

```
> cat commands
ls
sleep 600
ruby my_script.rb
pwd
bundle exec rails s
```

### Bugs

This software is alpha. It has bugs.

Please report any bug/suggestion.

### License

MIT
