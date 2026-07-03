// Copyright (c) the go-ruby-optparse authors
// SPDX-License-Identifier: BSD-3-Clause
//
// Library-level benchmark driver for the pure-Go go-ruby-optparse OptionParser.
// It exercises the same OptionParser lifecycle as ruby/optparse.rb — building a
// parser with several options (flags, "--opt VALUE", Integer/Float coercion, a
// negatable flag) and running parse! over an identical, representative argv — so
// the ns/op numbers compare this pure-Go library primitive against each Ruby
// runtime's own pure-Ruby optparse stdlib.
//
// With CHECK=1 it instead prints one "CHECK\t<label>\t<value>" line per op: an
// integer checksum of the op's result, used to prove the Go output is identical
// to MRI (the oracle) before any timing is trusted.
package main

import (
	"fmt"
	"math/big"
	"os"

	optparse "github.com/go-ruby-optparse/optparse"
)

// argv is a fixed, representative command line: short and long flags, a
// "--opt VALUE" separated form, an "--opt=VALUE" attached form, an Integer and a
// Float coercion, a negated "--no-color", and two interleaved operands. It is
// byte-for-byte the same argv the Ruby workload builds.
var argv = []string{
	"-v",
	"--name", "alice",
	"input1.txt",
	"-c", "42",
	"--rate=3.14",
	"--no-color",
	"-d",
	"--output", "out.log",
	"input2.txt",
}

// build constructs a fresh OptionParser and registers the eight options the
// Ruby workload registers, in the same order. It returns the built parser so a
// caller can either time construction alone or go on to parse.
func build() *optparse.Parser {
	p := optparse.New()
	p.Banner = "Usage: tool [options] files..."
	p.Define([]string{"-v", "--verbose", "run verbosely"}, optparse.CoerceNone, nil, nil)
	p.Define([]string{"-q", "--quiet", "suppress output"}, optparse.CoerceNone, nil, nil)
	p.Define([]string{"--name NAME", "the actor name"}, optparse.CoerceString, nil, nil)
	p.Define([]string{"-c", "--count N", "iteration count"}, optparse.CoerceInteger, nil, nil)
	p.Define([]string{"-r", "--rate X", "sampling rate"}, optparse.CoerceFloat, nil, nil)
	p.Define([]string{"--output FILE", "write output here"}, optparse.CoerceString, nil, nil)
	p.Define([]string{"--[no-]color", "colorize output"}, optparse.CoerceNone, nil, nil)
	p.Define([]string{"-d", "--debug", "enable debugging"}, optparse.CoerceNone, nil, nil)
	return p
}

// checksum folds a parse result into a single integer, identically to the Ruby
// driver, so the two can be proven to agree. It folds each match's coerced value
// (bool / integer / float / string) in argv order, then the operand and match
// counts and the operand byte lengths.
func checksum(matches []optparse.Match, rest []string) int {
	acc := 0
	for _, m := range matches {
		var c int
		switch v := m.Value.(type) {
		case bool:
			if v {
				c = 3
			} else {
				c = 5
			}
		case int64:
			c = int(v)
		case int:
			c = v
		case *big.Int:
			c = int(v.Int64())
		case float64:
			c = int(v * 1000)
		case string:
			c = len(v)
		case []string:
			for _, s := range v {
				c += len(s)
			}
		default:
			c = 0
		}
		acc = acc*33 + c
	}
	acc = acc*131 + len(matches)
	acc = acc*131 + len(rest)
	for _, r := range rest {
		acc = acc*33 + len(r)
	}
	return acc
}

// opBuild times parser construction alone (the construction-heavy path that
// dominates real optparse usage). Checksum = number of registered options.
func opBuild() int {
	p := build()
	sink = p
	return 8
}

// opParse times parse! against a parser built once outside the timed region,
// isolating the argv-parsing engine. Checksum = fold of the parse result.
func opParse() int {
	matches, rest, _ := parser.ParseBang(argv)
	return checksum(matches, rest)
}

// opBuildParse times the full lifecycle a program actually runs: build the
// parser, then parse!. Checksum = registered count folded with the parse result.
func opBuildParse() int {
	p := build()
	matches, rest, _ := p.ParseBang(argv)
	return 8*1_000_003 + checksum(matches, rest)
}

// parser is the pre-built parser used by opParse (construction excluded).
var parser = build()

// ops is the ordered op table shared by the timing and CHECK paths.
var ops = []struct {
	label string
	fn    func() int
}{
	{"build", opBuild},
	{"parse", opParse},
	{"build-parse", opBuildParse},
}

func main() {
	if os.Getenv("CHECK") != "" {
		for _, o := range ops {
			fmt.Printf("CHECK\t%s\t%d\n", o.label, o.fn())
		}
		return
	}
	const inner = 500
	for _, o := range ops {
		fn := o.fn
		bench(o.label, inner, func() { sink = fn() })
	}
}
