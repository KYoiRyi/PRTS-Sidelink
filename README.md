# Session Server (Fork)

This repository is a fork of the original session server branch. 

## Overview

It serves as a standalone session and matchmaking server designed to handle multiplayer lobby management, real-time battle synchronization, and duel mechanics. 

Key features include:
- TCP-based real-time state synchronization.
- Configurable multiplayer waiting times and match constraints.
- Built-in simulated bot system to automatically fill empty player slots during matchmaking shortages.
- Independent from the main API server, focusing solely on in-session logic.

## Usage

1. Ensure Go 1.24+ is installed.
2. Configure the server parameters in `configs/config.yaml`.
3. Build the server using `go build ./...` or the provided build scripts.
4. Run the executable to start handling session requests.
