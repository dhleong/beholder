beholder [![Build Status](http://img.shields.io/travis/dhleong/beholder.svg?style=flat)](https://travis-ci.org/dhleong/beholder)
=========

*A CLI tool for D&D players*

## What?

Beholder provides an ultra-fast, "fuzzy" search interface to just about anything
in 5th edition D&D. It is meant as a quick reference for players, game masters,
and maybe even fans watching a [podcast][1]. It is not intended to be a replacement
for the Player's Handbook or Dungeon Master's Guide, and as such there's not really
a "browse" feature. Instead, it is intended to augment the official materials,
providing a quick and easy way to look up the specifics of a monster, spell, class
feature, or other "thing" in the D&D universe.

## How?

Precompiled builds are provided for the [latest release][2].

Beholder is written in [Go][3], so you can also `go get` it:

    go get -u github.com/dhleong/beholder

On macOS you can install it via [homebrew][4]:

    brew tap dhleong/beholder
    brew install beholder

[1]: https://www.youtube.com/watch?v=byva0hOj8CU
[2]: https://github.com/dhleong/beholder/releases/latest
[3]: https://golang.org/
[4]: https://brew.sh/
