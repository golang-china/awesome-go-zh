# 🎮 awesomenes

A NES emulator written in Go.

<p align="center">
  <img src="https://i.imgur.com/z8xYcxV.png" alt="dk"  width="400px"/>
  <img src="https://i.imgur.com/ahSN16z.png" alt="smb" width="400px"/>
</p>
<p align="center">
  <img src="https://i.imgur.com/XX03vOV.png" alt="dk"  width="400px"/>
</p>

# Getting and running

`awesomenes` uses `sld2` for rendering and input processing. It may be necessary to install it beforehand. On macOS, using homebrew:

```
$ brew install sdl2
```

Other systems/package managers should provide similar `sld2`/`libsdl2` packages. Then use the `go get`:

```
$ go get github.com/rbaron/awesomenes
$ awesomenes MY_ROM.nes
```

# Status

Games that use the [mapper](http://wiki.nesdev.com/w/index.php/Mapper) 0 (NROM) mostly work, although without audio so far. Games that use mapper 4 (mmc3) should work with some eventual glitches. 

Tested games:

- Donkey Kong (NROM)
- Super Mario Bros. (NROM)
- Super Mario Bros. 2 (mmc3)
- Super Mario Bros. 3 (mmc3, with some glitches)

# Controller inputs

## Keyboard (controller 1)

```
Arrow keys  -> NES arrows
A           -> NES A
S           -> NES B
Enter       -> NES start
Right shift -> NES select
```

## Nintendo Switch Joycon (controller 1)

```
Directional -> NES arrows
Down arrow  -> NES A
Right arrow -> NES B
SL          -> NES select
SR          -> NES start
```

# Roadmap

✅ CPU emulation

✅ Video support (picture processing unit - PPU)

✅ Keyboard input

✅ Mapper 0

✅ Joystick input (tested with Nintendo Switch Joycon)

✅ Mapper 4 (...kinda)

➖ More mappers

➖ Save state

➖ Audio support (audio processing unit - APU)


# Resources

All the information used to build this emulator was found on the awesome [nesdev wiki](https://wiki.nesdev.com).
