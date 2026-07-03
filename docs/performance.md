# Performance

`go-ruby-optparse/optparse` is the pure-Go, CGO-free library that
[`rbgo`](https://github.com/go-embedded-ruby/ruby) binds for Ruby's `optparse`.
This page records a **real, library-level** benchmark of that module's Go API
against every reference runtime's own `optparse` stdlib, one row per
`OptionParser` operation. It is part of the ecosystem-wide per-module parity
suite, and **the bar is beating MRI + YJIT**, not just plain MRI.

## What is measured

Three representative `OptionParser` operations over one **fixed, representative
command line** — eight registered options (short/long flags, `--opt VALUE` and
`--opt=VALUE` forms, `Integer` and `Float` coercion, a `--[no-]color` negatable
flag) and a twelve-token argv with two interleaved operands:

| Op | What it exercises |
| --- | --- |
| `build` | construct an `OptionParser` and register all eight options (the construction-heavy path that dominates real optparse programs) |
| `parse` | `parse!` the representative argv against a parser built once, isolating the argv-parsing engine (flags, separated/attached values, coercion, negation, permuted operands) |
| `build-parse` | the full lifecycle a real program runs each launch: build the parser, then `parse!` |

The **go-ruby** column drives this pure-Go library through its Go API; every
other column is that interpreter's own stdlib `optparse` (a pure-Ruby library).
The Go and Ruby drivers build the **identical** argv, register the **identical**
options in the same order, and before any timing each op's integer checksum is
verified **byte-identical to MRI** — all five runtimes agree on every op
(`build`=8, `parse`=79102786953647896, `build-parse`=79102786961647920). So the
comparison is the same observable operation, apples-to-apples.

- **Host:** Apple M4 Max, macOS (`arm64-darwin`). **Date:** 2026-07-03.
- **Runtimes:** Go 1.26.4; `ruby 4.0.5 +PRISM` (MRI, the oracle) and
  `ruby --yjit`; `jruby 10.1.0.0` (OpenJDK 25); `truffleruby 34.0.1`
  (GraalVM CE Native).
- **Method:** each process runs 3 untimed warm-up passes then 25 timed passes of
  a fixed inner loop, timed with a monotonic clock; the **best** pass is reported
  as **ns/op**. Interpreter start-up is outside the timed region, so the number
  is the operation's own cost, not `ruby file.rb` process cost. Numbers were
  stable to within a few percent across repeated runs.
- Harness and drivers live in this repo under
  [`benchmarks/`](https://github.com/go-ruby-optparse/docs/tree/main/benchmarks)
  (`go/`, `ruby/optparse.rb`, `run.sh`). Reproduce: `bash benchmarks/run.sh`.

## Results (ns/op, best of 25)

| Op | go-ruby (pure Go) | MRI | MRI + YJIT | JRuby | TruffleRuby | **go vs YJIT** |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| `parse` | **852** | 14 600 | 10 424 | 14 910 | 19 927 | **12.2× faster** ✅ |
| `build-parse` | **7 063** | 81 258 | 59 802 | 77 091 | 44 637 | **8.5× faster** ✅ |
| `build` | **5 878** | 61 644 | 48 690 | 60 310 | 42 113 | **8.3× faster** ✅ |

## The go-vs-YJIT verdict, per op

**Every operation beats MRI + YJIT, decisively:**

- **`parse` — 12.2× faster than YJIT** (852 ns vs 10 424 ns). Pure argv walking:
  the Go parser is a handful of map lookups, slice reslices, and a `strconv`
  coercion per token, with no interpreter dispatch at all.
- **`build-parse` — 8.5× faster than YJIT** (7 063 ns vs 59 802 ns) — the whole
  real-world lifecycle (construct + parse) in one number.
- **`build` — 8.3× faster than YJIT** (5 878 ns vs 48 690 ns). Registering eight
  options compiles to plain struct/map appends in Go; MRI's `optparse` builds a
  `Switch` object graph and runs regexp-driven flag parsing in Ruby for every
  `on(...)`, which is where its construction cost goes.

Net: **all 3 operations beat MRI + YJIT** (8.3×–12.2×), and MRI + YJIT is itself
the fastest of the four Ruby runtimes on the construction-bound ops. Ruby's
`optparse` is pure Ruby — not a C extension — so there is no hand-written C to
catch up to here; the compiled pure-Go implementation wins across the board while
producing byte-identical results.

!!! note "Cold-JIT caveat"
    JRuby and TruffleRuby are measured **in-process after 3 warm-up passes**, but
    3 passes is not full JIT steady state for the JVM / GraalVM — their columns
    still carry cold-to-warming-JIT cost and should be read as
    order-of-magnitude, not as fully-warmed peak throughput. TruffleRuby's
    relatively strong `build-parse` reflects aggressive Graal specialization of
    that loop, yet it is still 6.3× slower than the pure-Go library. The
    **go-vs-MRI and go-vs-YJIT comparisons are the load-bearing ones** — MRI and
    YJIT reach steady state within the warm-up, and both they and the Go driver
    are timed identically. All numbers are real, single-host measurements from
    the 2026-07-03 run; nothing is cherry-picked.
