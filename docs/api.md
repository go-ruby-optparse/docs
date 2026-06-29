# Usage & API

The public API lives at the module root (`github.com/go-ruby-optparse/optparse`). It is **Ruby-shaped but Go-idiomatic**: a `Parser` holds option `Spec`s and produces `Match`es; the per-option callbacks `OptionParser` would invoke are the host's job. The surface follows Go conventions — an explicit `error`, value types, no global state.

!!! success "Status: implemented"
    The library is built and importable as `github.com/go-ruby-optparse/optparse`, bound into
    `rbgo` as a native module; see [Roadmap](roadmap.md).

## Install

```sh
go get github.com/go-ruby-optparse/optparse
```

## Worked example

```go
p := optparse.New()
p.On(optparse.Spec{Long: "verbose", Short: "v"})
p.On(optparse.Spec{Long: "level", Arg: optparse.Required, Type: optparse.Int})

m, rest, err := p.ParseBang([]string{"-v", "--lev=3", "file"})
// m matches -v and the abbreviated --level=3 (coerced to int 3)
// rest == []string{"file"}; the host runs each Spec's callback
```

## Shape

```go
// Parser holds the option specs and parses an argv against them.
type Parser struct{ /* ... */ }

// Spec describes one option (long/short names, argument, type).
type Spec struct{ /* ... */ }

// Match is one resolved option occurrence (which Spec, its value).
type Match struct{ /* ... */ }

// ParseBang implements parse! — permute and consume the argv,
// returning the matches, the residual argv, and any ParseError.
func (p *Parser) ParseBang(argv []string) (matches []Match, rest []string, err error)
```

## MRI conformance

Correctness is defined by reference Ruby. A **differential oracle** runs a wide
corpus through both the system `ruby` and this library and compares the results
**byte-for-byte** — not approximated from memory. The oracle tests skip
themselves where `ruby` is not on `PATH` (e.g. the qemu arch lanes), so the
cross-arch builds still validate the library.

## Relationship to Ruby

`go-ruby-optparse/optparse` is **standalone and reusable**, and is the backend bound into
[go-embedded-ruby](https://github.com/go-embedded-ruby/ruby) by `rbgo` as a
native module — the same way [go-ruby-regexp](https://github.com/go-ruby-regexp)
and [go-ruby-erb](https://github.com/go-ruby-erb) are bound. The dependency runs
the other way: this library has no dependency on the Ruby runtime.
