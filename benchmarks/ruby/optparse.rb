# frozen_string_literal: true
# Copyright (c) the go-ruby-optparse authors
# SPDX-License-Identifier: BSD-3-Clause
#
# Reference OptionParser workload, mirroring benchmarks/go/main.go op-for-op over
# an identical, representative argv. Run normally it reports ns/op per op through
# the shared harness; run with CHECK=1 it prints one "CHECK\t<label>\t<value>"
# line per op so the Go output can be proven identical to MRI (the oracle) before
# any timing is trusted.
require "optparse"
require_relative "_harness"

# Byte-for-byte the same argv the Go driver builds: short and long flags, a
# "--opt VALUE" separated form, an "--opt=VALUE" attached form, an Integer and a
# Float coercion, a negated "--no-color", and two interleaved operands.
ARGV_INPUT = [
  "-v",
  "--name", "alice",
  "input1.txt",
  "-c", "42",
  "--rate=3.14",
  "--no-color",
  "-d",
  "--output", "out.log",
  "input2.txt",
].freeze

# HITS collects each option's coerced value in the order the switches fire (argv
# order), exactly like the Go driver's ordered Match slice.
HITS = []

# build_parser registers the same eight options, in the same order, as the Go
# driver's build(). Each block pushes the value optparse hands it into HITS.
def build_parser
  OptionParser.new do |o|
    o.banner = "Usage: tool [options] files..."
    o.on("-v", "--verbose", "run verbosely") { |x| HITS << x }
    o.on("-q", "--quiet", "suppress output") { |x| HITS << x }
    o.on("--name NAME", "the actor name") { |x| HITS << x }
    o.on("-c", "--count N", Integer, "iteration count") { |x| HITS << x }
    o.on("-r", "--rate X", Float, "sampling rate") { |x| HITS << x }
    o.on("--output FILE", "write output here") { |x| HITS << x }
    o.on("--[no-]color", "colorize output") { |x| HITS << x }
    o.on("-d", "--debug", "enable debugging") { |x| HITS << x }
  end
end

# checksum folds a parse result into one integer identically to the Go driver:
# each value (bool / Integer / Float / String) in argv order, then the match and
# operand counts, then the operand byte lengths.
def checksum(hits, rest)
  acc = 0
  hits.each do |v|
    c =
      case v
      when true    then 3
      when false   then 5
      when Integer then v
      when Float   then (v * 1000).to_i
      when String  then v.bytesize
      when Array   then v.sum(&:bytesize)
      else 0
      end
    acc = acc * 33 + c
  end
  acc = acc * 131 + hits.length
  acc = acc * 131 + rest.length
  rest.each { |r| acc = acc * 33 + r.bytesize }
  acc
end

# build: construct the parser alone (the construction-heavy path). Checksum = 8.
def op_build
  $sink = build_parser
  8
end

# BUILT is the parser reused by op_parse so construction stays outside timing.
BUILT = build_parser

# parse: parse! against the pre-built parser, isolating the parsing engine.
def op_parse
  HITS.clear
  rest = BUILT.parse!(ARGV_INPUT.dup)
  checksum(HITS, rest)
end

# build-parse: the full lifecycle a real program runs — build then parse!.
def op_build_parse
  p = build_parser
  HITS.clear
  rest = p.parse!(ARGV_INPUT.dup)
  8 * 1_000_003 + checksum(HITS, rest)
end

OPS = [
  ["build",       method(:op_build)],
  ["parse",       method(:op_parse)],
  ["build-parse", method(:op_build_parse)],
].freeze

if ENV["CHECK"] && !ENV["CHECK"].empty?
  OPS.each { |label, m| printf("CHECK\t%s\t%d\n", label, m.call) }
else
  INNER = 500
  OPS.each { |label, m| bench(label, INNER) { m.call } }
end
