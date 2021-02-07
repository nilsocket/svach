# svach

Sanitize file name

## Library

[Documentation](https://pkg.go.dev/github.com/nilsocket/svach)

## CMD - Installation

```sh
go get github.com/nilsocket/svach/cmd/svach
```

## CMD - Usage

```sh
svach            # print intended file name changes (Clean)
svach -r         # print intended file name changes recursively (Clean)

svach -c         # change intended file names (Clean)
svach -c -r      # change intended file names recursively (Clean)

svach -n         # print intended file name changes (Name)
svach -c -n -r   # change intended file names recursively (Name)
```

## Difference between Name and Clean

### Name

Creates a valid file name for all operating systems.

```sh
❯ svach -n 'Hello___World!!!!!/\\'
Hello___World!!!!!
```

### Clean

Creates a valid file name and removes all control characters, repeated separators (`_`, `-`. `+`, `\`, `!`, ` `).
Different kinds of space are replaced with normal space character.

```sh
❯ svach 'Hello___World!!!!!/\\'
Hello_World!
```
