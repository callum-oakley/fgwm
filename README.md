# fgwm (floating grid window manager)

*fgwm* (pronounced "fugwum") sits somewhere between a floating and a tiling
window manager. Movement and resizing is done manually, as in a floating WM,
but every window always snaps perfectly to a (coarse) grid. Movements and
resizes are always done in multiples of the width and height of a cell in this
grid. This allows for the pixel perfect and efficient layouts of a tiling WM,
with the flexibility and aesthetics of a floating WM.

## Install

*fgwm* relies on [wmutils][0] ([core][1] and [opt][2]) being installed and
available in your path.

### From source

You'll need a [go environment][3] and [dep][4].

    $ go get -u github.com/hot-leaf-juice/fgwm
    $ cd $GOPATH/src/github.com/hot-leaf-juice/fgwm
    $ dep ensure
    $ go install

## Run

The *fgwm* binary starts as a daemon when run with no commands, or acts as a
client when provided with a command. Start the daemon from `.xinitrc` or
similar:

    exec fgwm

and then then issue commands like

    $ fgwm throw left
    $ fgwm move 2 0

You'll want to bind these commands with a hotkey daemon such as [sxhkd][5].
See [here][6] for an example sxhkdrc. See [here][7] for a description of
available commands.

## Configure

By default *fgwm* looks in `$HOME/.config/fgwm/config.toml` for a configuration
file. See [here][8] for a documented example config. Pass an alternative path
with the `--config` or `-c` flags.

    exec fgwm -c CONFIG_PATH

## Commands

The following commands act on the currently focussed window unless otherwise
stated.

    fgwm center

Center the window on the screen.

    fgwm focus next

Focus the next window in the stack. See [Focus][9].

    fgwm focus prev

Focus the previous window in the stack. See [Focus][9].

    fgwm fullscreen

Toggle the current window in or out of full screen (restores the pre-fullscreen
position and size).

    fgwm grow x y

Resize the window so that the top and bottom edges move away from the center by
`x` cells each, and the left and right by `y` cells each. `x` or `y` negative
causes the window to shrink. e.g. `fgwm grow 2 0` makes the window two cells
wider in each direction, four cells wider over all.

    fgwm kill

Close the window in the current view (see [Views][10]).

    fgwm move x y

Move the window by `x` cells to the left, and `y` cells down. Negative
arguments reverse the direction. e.g. `fgwm move -2 0` moves the window two
cells to the left.

    fgwm snap

Force the window to the grid if something has put it out of alignment.

    fgwm spread direction

Where `direction` is one of `up`, `down`, `left`, or `right`. Resise the window
so that the side indicated by the direction moves all the way to the edge of
the screen.

    fgwm teleport a b c d

Move and resize the window so that it occupies the rectangle with top left
corner at `(a, b)`, and bottom left at, `(c, d)`. The top left corner is `(0,
0)`. e.g. `fgwm teleport 6 0 18 24` resizes the window to take up a third of
the screen, and places it centrally, in a `24x24` grid.

    fgwm throw direction

Where `direction` is one of `up`, `down`, `left`, or `right`. Moves the window
in the given direction all the way to the edge of the screen.

    fgwm view-include n

Includes the window in the view `n`. See [Views][10].

    fgwm view-set n

Sets the current view to `n`. See [Views][10].

    fgwm help

List available commands.

## Focus

*fgwm* maintains a stack of the *most recently used* windows, and `fgwm focus
next` focusses the next window in the stack (i.e. *the* most recently used
window). A window isn't considered *used* until it has been focussed for a set
amount of time (500ms by default, configurable with the `focus_timeout_ms`
option), so calling `fgwm focus next` again before that timeout moves to the
next window in the stack, and so on. When a window has been focussed for long
enough to be considered *used* it is moved to the top of the stack, and your
position in the stack is reset.

This allows a single call of `fgwm focus next` to swap between your two most
recently used windows, while multiple calls can be used to traverse the window
stack as far back as you like. `fgwm focus prev` moves back up the stack, and
is for situations where you accidentally focus past the window you wanted and
don't have to loop all the way back around.

## Views

*Views* are similar to *desktops*, *workspaces*, or *groups* found in most
window managers. They have the following properties:

- A view is denoted by a single integer.
- Exactly one view is active at a time, the active view can be set to `n` with
  `fgwm view-set n`.
- Windows can belong to any number of views, and have an independent position
  and size in each. To include a window in view `n`, use `fgwm view-include n`.

When a window is created, it gets included in the currently active view. When
including a window in another view, it initially has the same position and
location in both, but henceforth can be moved and resized indepently. Killing a
window which belongs to multiple views with `fgwm kill` only removes it from
the current view.

The initial view can be changed with the `initial_view` option (`1` by default).

[0]: https://github.com/wmutils
[1]: https://github.com/wmutils/core
[2]: https://github.com/wmutils/opt
[3]: https://golang.org/doc/install
[4]: https://github.com/golang/dep#setup
[5]: https://github.com/baskerville/sxhkd
[6]: https://github.com/hot-leaf-juice/dots/blob/master/.config/sxhkd/sxhkdrc
[7]: #commands
[8]: https://github.com/hot-leaf-juice/fgwm/blob/master/config.example.toml
[9]: #focus
[10]: #views
