# Test platformer

![demo](demo.gif)

## Purpose

This is an example program that demonstrate usage of
the coroutine library [carrot](https://github.com/nvlled/carrot).
The relevant files can be found in the dino\_\* directories,
in particular the [dino_coroutine/dino.go](dino_coroutine/dino.go).
See the Notes section of the this readme for more details.

Note: this whole project was written
while figuring how to use the [ebiten game engine](https://ebitengine.org/).

## How to run

1. `git clone https://blah/todo`
2. `go run .`
   or `./dev` to enable auto-reload

## Controls

- **left/right arrow keys** - move left and right
- **up/down arrow keys** - move up and down (only when flying)
- **space key** - jump while on ground or air, hold to jump higher

## Instructions

- **aerial jump** - Press space key again while in midair to jump further
- **flying** - To fly, hold left or right until the dino is running
  really fast, then do a triple jump. To stop flying, hold down key
  then press space key.
- **jump charge** - Some silly feature, for purposes of using coroutines
  with more complicated state transitions. To activate:

  1. press and hold space and up arrow until dino starts spinning
  2. release up key, then keep pressing the space key until dino stops spinning.
  3. keep pressing up, left, right or down arrows
  4. wait
  5. hold up, right, left or right, then press space to dash
  6. press space again to return

## Notes and discussion

- The platformer has scuffed, janky home-made physics and collistion detection.
  I could polish it more but it's not really the point of the demo.

- With coroutines, there are less methods to define,
  thus less boilerplate. Of course, methods can
  still be optionally used on coroutines,
  whereas on the other
  implementations, it can not be avoided.
  The added boilerplate makes it harder to iterate
  experimental or exploratory solutions (subjectively speaking).
  If the resulting coroutine code turned out to be messy, it
  can always be refactored later. Then again, complex stateful gameplay
  code are rarely clean (citation needed?).

- Consider the ControllerCoroutine function:
  it's about 500 lines of loosely organized
  imperative code (with goto you say!?!).
  But, it's an exploratory code that was barfed
  out with less effort without needing to define a bunch of
  classes, structs, methods or enums. Instead, most of
  the logic code is in one place, free from boilerplate.

- Also note the local, inter-frame states are possible
  with coroutines without needing to define it outside
  the coroutine scope. With the non-coroutine implementations,
  state would have to be moved to outer scopes,
  breaking encapsulation in some cases, or worse
  cluttering the parent state.

- Coroutines can also be used on sprite animation. Instead of
  just providing an array of tileIDs, frame delay and other
  state can be arbritrarily adjusted.

- There is less or worse locality with
  non-coroutine solutions. Related code are
  broken apart and further away, making it
  harder to read or change. To be fair, goto sort
  makes it harder too, depending who's reading it.

- With coroutines, I get a scripting functionality
  without needing to embed a scripting engine.
