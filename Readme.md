# Kogia

   Small and simple init system for docker containers

## Installation

You can use the provided Dockerfile which will compile a binary for you or build
yourself.

```bash
docker build -t kogia https://github.com/dmajere/kogia.git
docker run --name kogia kogia
docker cp kogia:/kogia ./
```

## USAGE:
   kogia [global options] command [command options] [arguments...]


### COMMANDS:
   help, h  Shows a list of commands or help for one command

### GLOBAL OPTIONS:
```
   --skip-preinit, -S           Do not execute preinit scripts [$KOGIA_SKIP_PREINIT]
   --preinit, -s "/etc/preinit.d"   Path to preinit scripts [$KOGIA_PREINIT_DIR]
   --skip-postinit, -P          Do not execute postinit scripts [$KOGIA_SKIP_POSTINIT]
   --postinit, -p "/etc/postinit.d" Path to postinit scripts [$KOGIA_POSTINIT_DIR]
   --skip-env, -E           Do not load additional env from files [$KOGIA_SKIP_ENV]
   --env, -e "/etc/env"         Path to additional env file [$KOGIA_ENV]
   --level, -l "warning"        Verbosity level [$KOGIA_LEVEL]
   --help, -h               show help
   --version, -v            print the version
```
