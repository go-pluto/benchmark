# IMAP Benchmark

A tool to generate IMAP traffic for [pluto](https://github.com/go-pluto/pluto), [Dovecot](https://www.dovecot.org), and other IMAP services (like [Gmail](https://www.google.com/gmail)).


## Traffic Generation

The major difference to previously introduced `imap-evaluation` is, that we now support **IMAP Sessions**. Sessions are sequences of IMAP commands that are executed consecutively. The commands are *more or less* reasonable.

For the moment we only focus on **state-changing** (i.e. write) commands like:
* CREATE
* DELETE
* APPEND
* STORE
* EXPUNGE


## Setup

To install `imap-benchmark`, please run

```
$ go get -u github.com/go-pluto/benchmark
```

Modify the config file `test-config.toml` and the user data base `userdb.passwd`.


## Usage

You can start benchmarking an IMAP service by running the `imap-benchmark.go` file.

```
$ go run imap-benchmark.go
```

Alternatively, you can provide paths for config file and userdb:

```
$ go run imap-benchmark.go --config /var/config.toml --userdb /var/private.passwd
```


## Logging

All response times are collected in a log file underneath the `results` folder.


## License

This project is [GPLv3](https://github.com/go-pluto/benchmark/blob/master/LICENSE) licensed.
