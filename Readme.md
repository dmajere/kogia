# Kogia

Kogia is a simple init for docker containers written in Go.

## Installation

You can use the provided Dockerfile which will compile a binary for you or build
yourself.

```bash
docker build -t kogia https://github.com/dmajere/kogia.git
docker run --name kogia kogia
docker cp kogia:/kogia ./
```
