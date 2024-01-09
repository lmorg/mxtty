<h1>mxtty</h1>

- [Multimedia Terminal Emulator](#multimedia-terminal-emulator)
- [How It Works](#how-it-works)
- [Whats Left To Do](#whats-left-to-do)
  - [Compatibility Guide](#compatibility-guide)
    - [xterm compatibility](#xterm-compatibility)
    - [extended features](#extended-features)
    - [unique features](#unique-features)
    - [application customisation](#application-customisation)
    - [supported platforms](#supported-platforms)
- [How To Support](#how-to-support)

## Multimedia Terminal Emulator

The aim of this project is to provide an easy to use terminal emulator that
supports inlining multimedia widgets using native code as opposed to web
technologies like Electron.

Currently the project is _very_ alpha.

The idea behind this terminal emulator is that is can be used by any $SHELL,
however hooks will be built into [Murex](https://github.com/lmorg/murex) so
the terminal will be instantly usable even before wider support across other
shells and command line applications is adopted.

At its heart, `mxtty` is a regular terminal emulator. Like Kitty, iTerm2, and
PuTTY (to name a few). But where `mxtty` differs is that it also supports
inlining rich content. Some terminal emulators support inlining images. Others
might also allow videos. But none bar some edge case Electron terminals offer
collapsible trees for JSON printouts. Easy to navigate directory views. Nor any
other interactive elements that we have come to expect on modern user
interfaces.

The few terminal emulators that do attempt to offer this usually fail to be
good, or even just compatible, with all the CLI tools that we've come to depend
on.

`mxtty` aims to do _both well_. Even if you never want for any interactive
widgets, `mxtty` will be a good terminal emulator. And for those who want a
little more GUI in their CLI, `mxtty` will be a great modern user interface.

## How It Works

`mxtty` uses SDL ([Simple DirectMedia Layer](https://en.wikipedia.org/wiki/Simple_DirectMedia_Layer))
which is a simple hardware-assisted multimedia library. This enables the
terminal emulator to be both performant and also cross-platform. Essentially
providing some of the conveniences that people have come to love from tools
like Electron while still offering the benefits of native code.

The multimedia and interactive components will be passed from the controlling
terminal applications via ANSI escape sequences. Before groan, yes I agree that
in-band escape sequences are a lousy way of encoding meta-information. However
to succeed at being a good terminal emulator, it needs to support some historic
design decisions no matter how archaic they might seem today. This allows
`mxtty` to work with existing terminal applications _and_ for third parties to
easily add support for their applications to render rich content in `mxtty`
without breaking compatibility for legacy terminal emulators.

## Whats Left To Do

In short, pretty much everything. Most of what has been detailed above is still
only aspirational.

### Compatibility Guide

#### xterm compatibility

- C0 codes
  - [x] common: can run most CLI applications
  - [ ] broad: can run more older or more CLI applications
  - [ ] xterm compatibility
- C1 codes
  - [x] common: can run most CLI applications
  - [ ] broad: can run older or more CLI applications
  - [ ] xterm compatible
- CSI codes
  - [ ] common: can run most CLI applications
  - [ ] broad: can run more older or more CLI applications
  - [ ] xterm compatible
- SGR codes
  - [x] common: can run most CLI applications
  - [x] broad: can run more older or more CLI applications
  - [ ] xterm compatible
- OSC codes
  - [x] common: can run most CLI applications
  - [x] broad: can run more older or more CLI applications
  - [ ] xterm compatible
- DCS codes
  - [ ] common: can run most CLI applications
  - [ ] broad: can run more older or more CLI applications
  - [ ] xterm compatible
- PM codes (out of scope)
  - [x] common: can run most CLI applications
  - [x] broad: can run more older or more CLI applications
  - [x] xterm compatible
- [x] runs `tmux` glitch-free
- [ ] runs `vim` glitch-free
- mouse support
  - [ ] common: can run most CLI applications
  - [ ] broad: can run more older or more CLI applications
  - [ ] xterm compatible
- [ ] alt character sets
- [ ] wide characters
- [ ] resize support
- [ ] scrollback history

#### extended features

- inlining images
  - [ ] own API?
  - [ ] iterm2 compatible
  - [ ] sixel compatible
- [ ] hyperlink support
- [ ] extended SGR (eg "standardised" features not found in xterm)

#### unique features

- code folding
  - [ ] alpha: available but expect changes to the API
  - [ ] stable: available to use in Murex
- table sorting
  - [ ] alpha: available but expect changes to the API
  - [ ] stable: available to use in Murex

#### application customisation

- [ ] default typeface
- [ ] default colour scheme
- [ ] default bell sound
- [ ] default term size
- [ ] default command / shell

#### supported platforms

Support for the following platforms is planned:

- Linux
  - [x] ArchLinux
  - [ ] Ubuntu
  - [ ] Rocky
- BSD
  - [ ] FreeBSD
  - [ ] NetBSD
  - [ ] OpenBSD
  - [ ] DragonflyBSD
- [x] macOS
- [ ] Windows


## How To Support

Regardless of your time and skill set, there are multiple ways you can support
this project:

- **contributing code**: This could be bug fixes, new features, or even just
  correcting any typos.
- **testing**: There is a plethora of different software that needs to run
  inside a terminal emulator and a multitude of distinct platforms that this
  could run on. Any support testing `mxtty` would be greatly appreciated.
- **documentation**: This is possibly the hardest part of any project to get
  right. Eventually documentation for this will follow the same structure as
  [Murex Rocks](https://murex.rocks) (albeit its own website) however, for now,
  any documentation written in markdown is better than none.
- **architecture discussions**: I'm always open to discussing code theory. And
  if it results in building a better terminal emulator, then that is a
  worthwhile discussion to have.
- **porting escape codes to other applications**: Currently [Murex](https://github.com/lmorg/murex)
  is the pioneer for supporting `mxtty`-specific ANSI escape codes. However it
  would be good to see some of these extensions expanded out further. Maybe
  even to a point where this terminal emulator isn't required any more than a
  place to beta test future proposed escape sequences.

