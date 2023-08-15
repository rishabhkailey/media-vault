## Issues
### DNS issue when rebuilding container
#
```log
failed to solve: mcr.microsoft.com/vscode/devcontainers/go:0-1.18: failed to do request: Head "https://mcr.microsoft.com/v2/vscode/devcontainers/go/manifests/0-1.18": dial tcp: lookup mcr.microsoft.com on 192.168.0.1:53: read udp 172.17.0.2:42270->192.168.0.1:53: i/o timeout
```

this issue seems to be due to build kit
* https://github.com/docker/compose/issues/9550
* https://github.com/docker/buildx/issues/1459

#### Fix
#
temp workaround is to disable the buildkit

we can add below in ~/.bashrc 
```bash
export DOCKER_BUILDKIT=0
```

or we can add the above before the command printed by the error log
```bash
DOCKER_BUILDKIT=0 docker compose --project-name media-service_devcontainer -f /home/rishabhveer/not-work/media-service/.devcontainer/docker-compose.yml -f /home/rishabhveer/.config/Code/User/globalStorage/ms-vscode-remote.remote-containers/data/docker-compose/docker-compose.devcontainer.build-1688880841176.yml build
```

### File permissions
#
```log
checking context: can't stat '/home/rishabhveer/not-work/media-service/.devcontainer/certs/CAs'
```

## Fix
#
```bash
# CAs not required
sudo rmdir .devcontainer/certs/CAs/
```