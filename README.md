# Subdomain Enumeration Tool

The Subdomain Enumeration Tool is a command-line utility written in Go for discovering subdomains of a target domain. It supports concurrent subdomain resolution and includes the ability to read custom subdomains from a file.

## Table of Contents
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [Custom Subdomains](#custom-subdomains)
- [Contributing](#contributing)
- [License](#license)

## Features
- Concurrent enumeration of subdomains
- Predefined list of common subdomains to check
- Ability to read custom subdomains from a file
- Configurable concurrency level
- Timeout for resolving subdomains

## Prerequisites
- Go (Golang) installed on your system
- Internet connectivity to resolve subdomains

## Installation
1. Clone the repository:

2. Build the program:

    ```sh
    go build subenum.go
    ```

3. Run the program:

    ```sh
    ./subenum <target_domain> <concurrency_level>
    ```

## Usage
- `<target_domain>`: The domain for which you want to enumerate subdomains.
- `<concurrency_level>`: The number of concurrent workers to use for enumeration.

Example:

```sh
./subenum example.com 10
