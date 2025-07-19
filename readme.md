# cs-schema-gen
This project aims to generate a clone of the csfloat.com item schema from CS2’s items_game.txt and related files, allowing you to generate your own market-accurate schema for analysis, tooling, and skin-based projects.

I’m quite new to Golang, so if anything is poorly written, feel free to open a PR! The codebase is okay and it works.

## Features
  ✅ Parses CS2 items_game.txt to structured schema data

  ✅ Outputs JSON just like https://csfloat.com/api/v1/schema

  ✅ Fast, no-dependency parsing using Go’s standard library

  ✅ Useful for CS2 skin projects, market analysis, or personal bots

## Prerequisites
- Go (Golang) installed
- CS2 files:
  - items_game.txt
  - Use Source2Viewer to extract the file from CS2
  - Make sure to extract the csgo_english file/or any language files

## Why bother?
Valve’s items_game.txt is difficult to read and consume directly for tooling, bots, or market graphs.
This project allows you to:
- Build your own CS2 schema for pricing tools
- Create personal trade or snipe bots
- Power inventory value calculators
- Build visual or CLI-based collection management tools

## Contributing
🚩 Open PRs if you see poorly written Go or have improvements.
🚩 If you find issues parsing a specific CS2 update or file version, please open an issue with:
  - CS2 version or date
  - File causing the error
  - The command you ran

## Disclaimer
This project is not affiliated with or endorsed by Valve. All CS2 assets and data belong to Valve Corporation.