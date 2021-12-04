<div align="center">
  <img src="doc/logo.png" width="250" />
  
  # WeaSeL
  
  A tool for running WSL distributions as easy as docker containers.
  <br/><br/>

</div>

## Description

WeaSeL tries to simplify the developer workflow by utilizing the Windows subsystem for Linux (WSL) for quickly spinning up a development environment containing your favourite tools.

The environment is defined like for any other docker container using a Dockerfile or also pre-existing images from the registry can be used.

The advantages compared to docker are that you get the simplicity of the docker workflow combined with higher performance as well as a better integration of Linux for Windows machines.

## Quick Start

WSL uses so-called distributions which are technically an archive containing the Linux filesystem. For building such an archive out of an existing docker image (in this case `debian:buster-slim`), you have to call:

```batch
weasel build --tag demo hub:debian:buster-slim
```

After executing this command it will be automatically stored and registered for the use with WeaSeL. For spinnung up the distribution as an instance and create an interactive session, call:

```batch
weasel run demo
```

If you just need to execute a command inside this distribution, it is also possible to append it to this command:

```batch
weasel run demo echo works!
```

## Installation

We directly provide the WeaSeL executable at the release page that must be downloaded. Then it can be placed at a location that is in the PATH variable.

There is even a simpler way to integrate WeaSeL in an exisiting repository. For that just download the wrapper batch script also available at the release section and place it at the root directory or the repository. When calling it, it automatically downloads the executable and executes it with the provided arguments.

> :warning: Note that there is only the executable for Windows because WSL is not available on other operating systems.
