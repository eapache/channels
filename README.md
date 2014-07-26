channels
========

[![Build Status](https://travis-ci.org/eapache/channels.svg?branch=master)](https://travis-ci.org/eapache/channels)
[![GoDoc](https://godoc.org/github.com/eapache/channels?status.png)](https://godoc.org/github.com/eapache/channels)

A collection of helper functions and special types for working with and
extending [golang](http://golang.org/)'s existing channels.

See https://godoc.org/github.com/eapache/channels for documentation.

Most of the channel types in this package are backed by a very fast queue
implementation that used to be built into this package but has now been
extracted into its own package at https://github.com/eapache/queue.

*Note:* Several types in this package provide so-called "infinite" buffers. Be
very careful using these, as no buffer is truly infinite. If such a buffer
grows too large your program will run out of memory and crash. Caveat emptor.
