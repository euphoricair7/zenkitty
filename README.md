# zenkitty

A terminal timer where a kitty wants you to focus.

zenkitty is a cat that sits with you. not a gadget. sit as long as you like, pause whenever, leave early, or stay past the bell. no streaks. no shame. no alarm.

## sit

```bash
go run ./cmd/zenkitty
```

or install it:

```bash
go install github.com/euphoricair7/zenkitty/cmd/zenkitty@latest
```

needs a truecolor terminal if you can. still readable without one.

## keys

- `enter` / `s` / `space` — start sitting
- `space` — pause, or come back
- `d` — how long (or leave it open)
- `+` / `-` — a little more, a little less
- `n` — new sit
- `?` — a little note
- `q` — bye

if you pick a time, it is a soft until. the kitty gets sleepy near the end, then offers to stay a little longer. overtime is still a sit.

## how long

- a little while (~15)
- a good stretch (~25)
- a deep dive (~50)
- until we're done (open)
- or type minutes

default is open. the clock counts up, and nobody rings.
