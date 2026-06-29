# Roadmap

`go-ruby-optparse/optparse` is grown **test-first**, each capability differential-tested against MRI
rather than built in isolation. Ruby's optionparser — the
deterministic, interpreter-independent slice extracted from rbgo's internals — is
**complete**.

| Stage | What | Status |
| --- | --- | --- |
| Option specification | The `Spec` model for an option: long and short names, whether it takes a required/optional argument, the `--[no-]` negation form, and the coercion type — built the way `OptionParser#on` describes one. | **Done** |
| argv matching | Matching an argv against the spec: long `--opt`, short `-o`, **unambiguous abbreviations**, **bundled** short flags `-abc`, `--opt=val` and `-oval` argument attachment — exactly MRI's rules. | **Done** |
| Dispatch modes | `parse!` (permute, consume), `order!` (stop at the first non-option), `permute!` and the `getopts` convenience mode — each leaving the residual argv as MRI does. | **Done** |
| Type coercion | Coercion of option arguments to Integer, Float, and the other built-in OptionParser accept-types, raising MRI's `InvalidArgument` on a bad value. | **Done** |
| Error taxonomy | MRI's full `OptionParser::ParseError` family — `AmbiguousOption`, `InvalidOption`, `MissingArgument`, `NeedlessArgument`, `InvalidArgument` — surfaced as typed Go errors. | **Done** |
| Differential oracle & coverage | A wide argv/spec corpus parsed both here and by the system `ruby`, compared against MRI's residual argv and error classes; 100% coverage, gofmt + go vet clean, green across all six 64-bit Go arches and three OSes. | **Done** |

## Documented out-of-scope boundaries

These are **deliberate**, recorded so the module's surface is unambiguous:

- **No interpreter.** The library implements the deterministic algorithm; it
  never runs arbitrary Ruby. Anything that needs a live binding or evaluation is
  the consumer's job — that is why `rbgo` binds this module rather than the
  reverse.
- **Reference is reference Ruby (MRI).** Byte-for-byte conformance targets MRI's
  behaviour; differences across MRI releases are matched to the reference used by
  the differential oracle.
- **Standalone & reusable.** The module has no dependency on the Ruby runtime;
  the dependency runs the other way.

See [Usage & API](api.md) for the surface and [Why pure Go](why.md) for the
deterministic/interpreter split.
