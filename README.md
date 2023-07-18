# deployx

[![Go Reference](https://pkg.go.dev/badge/github.com/aaraney/deployx.svg)](https://pkg.go.dev/github.com/aaraney/deployx)

`deployx` is a Docker CLI plugin rewrite of `docker stack deploy` that is compliant
with [compose-spec](https://github.com/compose-spec/compose-spec).

## Features

- Compose file variable expansion (done by default)
- Support for one or more alternate `.env` files
- Support for container storage interface (CSI) [cluster volume](https://github.com/moby/moby/blob/master/docs/cluster_volumes.md) types

## Installing

### Prebuilt binaries

Interactive install script:

```bash
bash <(curl https://raw.githubusercontent.com/Rei-x/deployx/main/install.sh)
```

Or non-interactive install as docker-plugin (requires sudo):

```bash
curl -sL https://raw.githubusercontent.com/rei-x/deployx/main/install.sh | bash -s -- -y
```

Or non-interactive install just as a binary (./deployx):

```bash
curl -sL https://raw.githubusercontent.com/rei-x/deployx/main/install.sh | bash -s -- -n
```

### Building from source

Building from source requires go and make be installed. Install go using either your package manager
of choice (i.e. `brew`) or by following [these instructions](https://go.dev/doc/install).

```shell
# clone repo
git clone https://github.com/aaraney/deployx.git && cd deployx

make build && make install

# uninstall
# make uninstall
```

Note, the instructions above install `deployx` for a single user. To install globally run `make build`
and copy the `docker-deployx` binary from the `./build` directory to one of the following locations:

- `/usr/local/lib/docker/cli-plugins` OR `/usr/local/libexec/docker/cli-plugins`
- `/usr/lib/docker/cli-plugins` OR `/usr/libexec/docker/cli-plugins`

### Brew

Install using `homebrew`:

```shell
brew install aaraney/tap/deployx

# install as docker cli plugin. invoke using 'docker deployx'
ln -s $(which deployx) $HOME/.docker/cli-plugin/docker-deployx
```

### Dockerfile

The easiest way to get started it using a pre-built docker image and `deployx` in standalone mode.
The following snippet shows pulling and running `deployx`. Adjust the volume mount accordingly to
mount your compose and env files.

```shell
docker run -it --rm --volume $(pwd):/home --volume /var/run/docker.sock:/var/run/docker.sock aaraney/docker-deployx
docker-deployx --compose-file /home/<compose.yaml> mystack
```

## Usage

Use `deployx` just as you would `docker stack deploy` by instead calling `docker deployx`. Unlike
`docker stack deploy`, environment and `.env` variables are interpolated into compose files. So,
there is no longer a need to: `docker stack deploy -c <(docker-compose config) stack-name`.

```shell
$ docker deployx
Usage:  docker deployx [OPTIONS] STACK

Deploy a new stack or update an existing stack

Options:
  -c, --compose-file strings   Path to a Compose file, or "-" to read from stdin
      --env-file strings       Path to an alternative env file, or "-" to read from stdin
      --no-interpolate         Don't perform environment variable interpolation
      --prune                  Prune services that are no longer referenced
      --resolve-image string   Query the registry to resolve image digest and supported platforms ("always", "changed", "never") (default "always")
      --with-registry-auth     Send registry authentication details to Swarm agents
```
