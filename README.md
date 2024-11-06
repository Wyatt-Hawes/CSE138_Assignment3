# CSE138_Assignment3

# Setup

- Install `deno` for local testing at [https://docs.deno.com/runtime/getting_started/installation/](https://docs.deno.com/runtime/getting_started/installation/)

- Download the `Deno` vscode extension

- type `> Deno: Enable` in the top search bar to enable deno package importing into vs-code

- Now you should be able to run deno normally with
  `deno run --allow-net server.ts`
  - you can use `deno run --allow-net --watch server.ts` to have the server auto-refresh with your code changes

# Docker

- Run in the docker container with:

  - `docker build -t dino .`

  - `docker run --rm -p 8090:8090 dino`
