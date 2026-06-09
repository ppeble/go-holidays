#!/usr/bin/env ruby
# frozen_string_literal: true

# parity/oracle.rb
#
# Ruby "oracle" for apples-to-apples parity testing against the Go port.
#
# It runs the installed `holidays` gem (CODE) but loads OUR v7.0.0 region YAML
# (DATA) from the `definitions/` submodule via Holidays.load_custom on startup.
# So the holiday RULES are identical to what the Go port compiles, and any
# difference in output is a real behavioural difference in the engines, not in
# the data.
#
# ============================================================================
# RUN CONTRACT (stable; the Go harness in a later bead depends on this)
# ============================================================================
#
# Transport: LINE-DELIMITED JSON ("JSON Lines" / NDJSON).
#   - Read: exactly one JSON object per input line on stdin.
#   - Write: exactly one JSON object per output line on stdout, in the same
#     order as the requests. stdout is line-buffered (sync = true), so a Go
#     harness can write a request line and read one response line back.
#   - Blank lines on stdin are ignored.
#   - One malformed/erroring request does NOT abort the process; it yields one
#     error response line, and processing continues with the next line.
#
# REQUEST shape (one JSON object per line):
#   { "func": "<name>", ... }
# Each request MAY carry an "id" passed through verbatim to the response, to
# let the caller correlate out-of-band. Every request carries a "func" naming
# one of the supported functions below, plus that function's arguments.
#
# Dates are always strings "YYYY-MM-DD" (UTC calendar dates, no time/zone).
# Regions are an array of strings, e.g. ["us"] or ["us_va","ca"]; a bare
# string is also accepted and wrapped into a one-element array.
# Options are booleans:
#   "informal": true  -> include informal holidays   (gem flag :informal)
#   "observed": true  -> use observed dates           (gem flag :observed)
#
# RESPONSE shape (one JSON object per line):
#   { "ok": true,  "id": <echoed-or-null>, "func": "<name>", "result": <...> }
#   { "ok": false, "id": <echoed-or-null>, "func": "<name>",
#     "error": "<class>: <message>" }
#
# A "holiday list" result is ALWAYS normalized to a sorted array of objects:
#   [ { "date": "YYYY-MM-DD", "name": "..." }, ... ]
# sorted by date ascending, then by name ascending. (regions and other gem
# fields are dropped: the Go port's public API exposes date+name.)
#
# ----------------------------------------------------------------------------
# SUPPORTED FUNCTIONS
# ----------------------------------------------------------------------------
#
# 1) on
#    req:  { "func":"on", "date":"2024-12-25", "regions":["us"],
#            "informal":false, "observed":false }
#    res:  result = holiday-list   (holidays on that single date)
#
# 2) between
#    req:  { "func":"between", "start":"2024-01-01", "end":"2024-12-31",
#            "regions":["us"], "informal":false, "observed":false }
#    res:  result = holiday-list   (holidays in [start, end] inclusive)
#
# 3) cache_between
#    req:  { "func":"cache_between", "start":"2024-01-01", "end":"2024-12-31",
#            "regions":["us"], "informal":false, "observed":false }
#    res:  result = holiday-list   (the holidays the gem cached for that range;
#                                   gem returns a hash, normalized to the list)
#
# 4) next_holidays
#    req:  { "func":"next_holidays", "count":3, "regions":["us"],
#            "from":"2024-12-01", "informal":false, "observed":false }
#    res:  result = holiday-list   ("from" optional, defaults to gem default)
#
# 5) year_holidays
#    req:  { "func":"year_holidays", "regions":["us"], "from":"2024-01-01",
#            "informal":false, "observed":false }
#    res:  result = holiday-list   (one calendar year starting at "from";
#                                   "from" optional, gem defaults to today)
#
# 6) any_holidays_during_work_week?
#    req:  { "func":"any_holidays_during_work_week?", "date":"2024-12-25",
#            "regions":["us"], "informal":false, "observed":false }
#    res:  result = true | false   (boolean, NOT a holiday list)
#
# 7) available_regions
#    req:  { "func":"available_regions" }
#    res:  result = ["ar","at",...]   (sorted array of region-code strings)
#
# 8) load_custom
#    req:  { "func":"load_custom", "files":["definitions/us.yaml"] }
#    res:  result = { "loaded": <int> }   (count of files ingested; lets a
#          caller load extra YAML at runtime. The standard region set is
#          already loaded at startup, so this is rarely needed.)
#
# ----------------------------------------------------------------------------
# STARTUP
# ----------------------------------------------------------------------------
# On boot, every definitions/*.yaml (excluding index.yaml, which is an
# aggregator and not a region) is ingested via Holidays.load_custom. The path
# to definitions/ is resolved relative to this file (../definitions), so the
# oracle can be launched from any working directory.
#
# Re-loading a region YAML that the gem already knew (its bundled defs are
# v6.0.0; ours are v7.0.0) does NOT duplicate holidays: load_custom merges by
# region key and overrides cleanly. Verified: loading us.yaml twice still
# yields exactly one "Christmas Day" on 2024-12-25, and US-2024 = 10 holidays.
# ============================================================================

require "holidays"
require "json"
require "date"

DEFINITIONS_DIR = File.expand_path("../definitions", __dir__)

# Ingest our v7.0.0 YAML into the gem. Returns the count of files loaded.
def load_all_definitions
  files = Dir.glob(File.join(DEFINITIONS_DIR, "*.yaml"))
              .reject { |f| File.basename(f) == "index.yaml" }
              .sort
  files.each { |f| Holidays.load_custom(f) }
  files.size
end

# Normalize a "bare string or array" region argument to an array of symbols.
def regions_arg(req)
  raw = req["regions"]
  raw = [raw] unless raw.is_a?(Array)
  raw.compact.map(&:to_sym)
end

# Translate boolean option flags into the gem's trailing symbol flags.
def option_flags(req)
  flags = []
  flags << :informal if req["informal"]
  flags << :observed if req["observed"]
  flags
end

def parse_date(str)
  Date.strptime(str, "%Y-%m-%d")
end

# Normalize any gem holiday collection (array of hashes, or the hash that
# cache_between returns) into a sorted array of {date, name}.
def normalize_holidays(collection)
  list =
    case collection
    when Hash  then collection.values.flatten
    when Array then collection
    else Array(collection)
    end

  list
    .map { |h| { "date" => h[:date].strftime("%Y-%m-%d"), "name" => h[:name] } }
    .sort_by { |h| [h["date"], h["name"]] }
end

# ---- dispatch ----------------------------------------------------------------

def dispatch(req)
  case req.fetch("func")
  when "on"
    date = parse_date(req.fetch("date"))
    normalize_holidays(Holidays.on(date, *regions_arg(req), *option_flags(req)))

  when "between"
    s = parse_date(req.fetch("start"))
    e = parse_date(req.fetch("end"))
    normalize_holidays(Holidays.between(s, e, *regions_arg(req), *option_flags(req)))

  when "cache_between"
    s = parse_date(req.fetch("start"))
    e = parse_date(req.fetch("end"))
    normalize_holidays(Holidays.cache_between(s, e, *regions_arg(req), *option_flags(req)))

  when "next_holidays"
    count = Integer(req.fetch("count"))
    opts  = regions_arg(req) + option_flags(req)
    if req["from"]
      normalize_holidays(Holidays.next_holidays(count, opts, parse_date(req["from"])))
    else
      normalize_holidays(Holidays.next_holidays(count, opts))
    end

  when "year_holidays"
    opts = regions_arg(req) + option_flags(req)
    if req["from"]
      normalize_holidays(Holidays.year_holidays(opts, parse_date(req["from"])))
    else
      normalize_holidays(Holidays.year_holidays(opts))
    end

  when "any_holidays_during_work_week?"
    date = parse_date(req.fetch("date"))
    Holidays.any_holidays_during_work_week?(date, *regions_arg(req), *option_flags(req))

  when "available_regions"
    Holidays.available_regions.map(&:to_s).sort

  when "load_custom"
    files = Array(req.fetch("files"))
    files.each { |f| Holidays.load_custom(f) }
    { "loaded" => files.size }

  else
    raise ArgumentError, "unknown func: #{req["func"].inspect}"
  end
end

# ---- main loop ---------------------------------------------------------------

def main
  $stdout.sync = true
  load_all_definitions

  $stdin.each_line do |line|
    line = line.strip
    next if line.empty?

    req = nil
    begin
      req = JSON.parse(line)
      result = dispatch(req)
      puts JSON.generate("ok" => true, "id" => req["id"],
                         "func" => req["func"], "result" => result)
    rescue StandardError => e
      id   = req.is_a?(Hash) ? req["id"]   : nil
      func = req.is_a?(Hash) ? req["func"] : nil
      puts JSON.generate("ok" => false, "id" => id, "func" => func,
                         "error" => "#{e.class}: #{e.message}")
    end
  end
end

main if $PROGRAM_NAME == __FILE__
