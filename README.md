# Registry Cache

This repo provides several way to implement a Registry Cache (or pull-through mirror).
The goal is to explore several implementations to easily deploy high-performant, and easy-to-use Registry mirrors.

## Dagger Module

Initialize a local Registry Cache (caching images from docker.io)

```console
dagger -m github.com/samalba/registry-cache/dagger call --sock /var/run/docker.sock init-registry-mirror
```

Initialize a local Registry Cache - use a specific directory on the host for storage:

```console
dagger -m github.com/samalba/registry-cache/dagger call --sock /var/run/docker.sock init-registry-mirror --storage-path "$PWD/storage"
```

Re-configure the dagger engine to leverage the local mirror:
```console
dagger -m github.com/samalba/registry-cache/dagger call --sock /var/run/docker.sock configure-dagger-engine
# Copy paste the export command from stdout
export _EXPERIMENTAL_DAGGER_RUNNER_HOST=docker-container://<NEW CONTAINER IMAGE>
# From now on, all dagger calls will use the registry mirror
```

Re-configure the dagger engine to point to a remote mirror (local network or internet):

```console
dagger -m github.com/samalba/registry-cache/dagger call --sock /var/run/docker.sock configure-dagger-engine --registry-mirror my.registry.domain.tld:5000
export _EXPERIMENTAL_DAGGER_RUNNER_HOST=docker-container://<NEW CONTAINER IMAGE>
# From now on, all dagger calls will use the registry mirror
```
