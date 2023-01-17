<img align="right" src="https://github.com/jafarlihi/file-hosting/blob/master/rssnix-logo.png?raw=true">

## Table of Contents

* [Demonstration](#demonstration)
* [Installation](#installation)
* [Packages](#packages)
* [Flags](#flags)
* [Config](#config)

## Demonstration

![Demo](https://github.com/jafarlihi/file-hosting/blob/3613fb4d60f0fc64ce77d8c56279bcd6bdf769a3/rssnix-demo2.gif?raw=true)

## Installation

You need to have Go version 1.19 or higher installed.

`git clone https://github.com/jafarlihi/rssnix --depth=1 && cd rssnix && go install`

## Packages

<a href="https://aur.archlinux.org/packages/rssnix">Arch Linux AUR Package (build, newest version)</a> <img src="https://img.shields.io/aur/version/rssnix?color=green" alt="AUR"> 

<a href="https://aur.archlinux.org/packages/rssnix-bin">Arch Linux AUR Package (binary, newest version)</a> <img src="https://img.shields.io/aur/version/rssnix-bin?color=green" alt="AUR"> 

<a href="https://aur.archlinux.org/packages/rssnix-git">Arch Linux AUR Package (build from git)</a> <img src="https://img.shields.io/aur/version/rssnix-git?color=green" alt="AUR"> 

## Flags

`config`
- Opens config file with `$EDITOR`

`update [feed name]`
- If the [feed name] argument is given and is a space-delimited list of feeds, then these feeds are updated
- If no [feed name] argument is given then all feeds are updated

`open [feed name]`
- If the [feed name] argument is given then the said feed's directory is opened with the configured viewer
- If no [feed name] argument is given then the root feed's directory is opened with the configured viewer

`add [feed name] [feed url]`
- Adds a new feed to the config file
- Example: `rssnix add CNN-Tech http://rss.cnn.com/rss/edition_technology.rss`

`import [OPML URL or file path]`
- Imports feeds from OPML file
- Example: `rssnix import feeds.opml`

`refetch [feed name]`
- delete and refetch given feed(s) or all feeds if no argument is given

`version`
- Prints the rssnix version

## Config

The config file is expected to be at `~/.config/rssnix/config.ini`.

Sample config file:

```
[settings]
viewer = vim
# The viewer option specifies the program that will be used to open the feed directory.
# By default, the value is set to vim, but you can change it to any other text editor of your choice.

feed_directory = ~/rssnix
# The feed_directory option specifies the location of the feed files on the file system.

[feeds]
CNN-Tech = http://rss.cnn.com/rss/edition_technology.rss
HackerNews = https://news.ycombinator.com/rss
```
(Tip: `ranger` is another great candidate for `viewer`)
