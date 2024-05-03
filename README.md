# Pokedexcli

![](./out.gif)

## Overview
Here is a personal project built using [Go](https://go.dev/) and [Bubble Tea](https://github.com/charmbracelet/bubbletea/tree/master). It allows users to search for various pokemon and retrieve basic information about pokemon in the national pokedex. It's purpose is personal: an outlet for me to explore the Go standard library and some of the community tools like [Boltdb](https://github.com/boltdb/bolt) and the [Charm libraries](https://charm.sh/libs/). Hope you'll like what you see!

## Features

- [pokeapi.co REST api](https://pokeapi.co/) integration
- configurable caching layer
- pokemon search with a basic phonetic algorithm (soundex)
- ASCII rendering of png images
- responsive(ish) design using Bubble Tea messages and [Lip Gloss](https://github.com/charmbracelet/lipgloss)

## Prerequisites
Make sure you have the following installed before running the application:

- Go (version 1.21.5 or higher)
- Git

## Installation
Clone the repository:
```bash
git clone https://github.com/ethanefung/pokedexcli.git
```
Navigate to the project directory:
```bash
cd pokedexcli
```
Build the application:
```go
go build
```
## Usage
Run the application:
```bash
./pokedexcli
```

By default boltdb will be used for the cache. However, you can pass a flag to explore different caching strategies
```bash
./pokedexcli -cache filebased
```
or
```bash
./pokedexcli -cache inmemory
```

## Contributing
Because the purpose of this application is for personal edification, I don't plan on accepting pull requests in this project. I encourage forks of this project however! I will keep track of any issues, so if you do encounter a bug or would like to pitch an idea of how this application can improve be sure to let me know in an issue.

## License
This project is licensed under the MIT License. See the LICENSE file for details.
