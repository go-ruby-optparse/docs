# go-ruby-optparse documentation

**Ruby's `OptionParser` argv-parsing engine in pure Go — MRI-compatible, no cgo.**

`go-ruby-optparse/optparse` is a faithful, pure-Go (zero cgo) reimplementation of Ruby's OptionParser,
matching reference Ruby (MRI) byte-for-byte. The module path is
`github.com/go-ruby-optparse/optparse`.

It was **extracted from rbgo's prelude/internals into a reusable standalone
library**: the module is standalone and importable by any Go program, and it is
the backend bound into [go-embedded-ruby](https://github.com/go-embedded-ruby/ruby)
by `rbgo` as a native module — just like
[go-ruby-regexp](https://github.com/go-ruby-regexp) and
[go-ruby-erb](https://github.com/go-ruby-erb). The dependency runs the other
way: this library has **no dependency on the Ruby runtime**.

!!! success "Status: engine complete — MRI byte-exact"
    Faithful port of MRI's `lib/optparse.rb` matching engine: **long / short / abbreviated / bundled** options, the **`--[no-]opt`** form, **type coercion**, the **`parse!` / `order!` / `permute!` / `getopts`** modes, and the **full error taxonomy**. The per-option callbacks stay in the host. Validated by a **differential oracle** against the system `ruby` — residual argv and error classes compared — at 100% coverage, `gofmt` + `go vet` clean, CI green across the six 64-bit Go targets and three OSes.

## Quick taste

```go
p := optparse.New()
p.On(optparse.Spec{Long: "verbose", Short: "v"})
p.On(optparse.Spec{Long: "level", Arg: optparse.Required, Type: optparse.Int})

m, rest, err := p.ParseBang([]string{"-v", "--lev=3", "file"})
// m matches -v and the abbreviated --level=3 (coerced to int 3)
// rest == []string{"file"}; the host runs each Spec's callback
```

## Repositories

| Repo | What it is |
| --- | --- |
| [`optparse`](https://github.com/go-ruby-optparse/optparse) | the library — Ruby's OptionParser in pure Go |
| [`docs`](https://github.com/go-ruby-optparse/docs) | this documentation site (MkDocs Material, versioned with mike) |
| [`go-ruby-optparse.github.io`](https://github.com/go-ruby-optparse/go-ruby-optparse.github.io) | the organization landing page (Hugo) |
| [`brand`](https://github.com/go-ruby-optparse/brand) | logo and brand assets |

## Principles

- **Pure Go, `CGO_ENABLED=0`** — trivial cross-compilation, a single static
  binary, no C toolchain.
- **MRI byte-exact.** Output matches reference Ruby exactly, not approximately,
  validated by a differential oracle against the `ruby` binary.
- **Standalone & reusable.** Extracted from rbgo's internals; no dependency on
  the Ruby runtime — the dependency runs the other way.
- **100% test coverage** is the target, enforced as a CI gate, across 6 arches
  and 3 OSes.

## Where to go next

- [Why pure Go](why.md) — why this slice of Ruby is deterministic enough to live
  as a standalone, interpreter-independent Go library.
- [Usage & API](api.md) — the public surface and worked examples.
- [Roadmap](roadmap.md) — what is done and what is downstream by design.

Source lives at [github.com/go-ruby-optparse/optparse](https://github.com/go-ruby-optparse/optparse).
