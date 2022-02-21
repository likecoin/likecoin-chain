# Making release

We have official release at Github release and docker image on Docker Hub.

For release at Github release, it is done via Github Action automatically once
developer create tag which match `v*.*.*` regex. For detail, one can check the
`.github/workflows` folder.

For docker images, we are not doing automatic build. When the dev team push
releasing tag to Github, they will run the following.

```sh
make clean
make docker-build
make docker-push
```

Community should be able to find the release at
https://hub.docker.com/r/likecoin/likecoin-chain

## Building release locally

If community member would like to build the binary locally, one can run

```sh
make clean build
```

You will find the platform specific binary at `build/liked`

## Cross build

In case you want to build a linux binary for your validator locally at your
Macbook, we provide a cross platform binary build command as follow.

```sh
make clean build-reproducible
```

You will find few binary at `artifacts` folders. The binary name should be self
explanatory for what platform it built form.
