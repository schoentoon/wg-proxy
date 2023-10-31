# wg-proxy

This tool is intended for running a simple wireguard proxy in a rootless podman/kubernetes environment.
Because of this unlike most wireguard related OCI images it will use [the userspace implementation of wireguard](https://git.zx2c4.com/wireguard-go) instead of the kernel one.
To use this, first of all prepare a config file. Have a look at [the config example](./config.example.yml) and it should be fairly straight forward.
Afterwards simply run the OCI image using this config.

## Podman example

```bash
podman run -it --rm -v ./config.yml:/config.yml ghcr.io/schoentoon/wg-proxy:latest
```
