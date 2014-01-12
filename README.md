channels
========

[![Build Status](https://travis-ci.org/eapache/channels.png?branch=master)](https://travis-ci.org/eapache/channels)

A collection of helper functions and special types for working with and
extending golang's existing channels.

See https://godoc.org/github.com/eapache/channels for documentation.

*Note:* Several types in this package provide so-called "infinite" buffers. Be
very careful using these, as no buffer is truly infinite. If such a buffer
grows too large your program will run out of memory and crash. Caveat emptor.
