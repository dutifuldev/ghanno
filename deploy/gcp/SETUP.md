# GCP Deploy Setup

The GCP deployment mounts the GitHub App private key from `/home/bob/prtags/secrets` into the `prtags` container.

Use these ownership and mode settings on the VM before restarting the service:

```bash
sudo chown 10001:10001 /home/bob/prtags/secrets
sudo chmod 0500 /home/bob/prtags/secrets
sudo chown 10001:10001 /home/bob/prtags/secrets/github-app.private-key.pem
sudo chmod 0400 /home/bob/prtags/secrets/github-app.private-key.pem
```

Keep the app env file pointing at that mounted key path:

```env
GITHUB_APP_PRIVATE_KEY_PATH=/home/bob/prtags/secrets/github-app.private-key.pem
```

HTTP requests and workers use separate database pools. Keep both bounded on the shared Cloud SQL instance:

```env
DB_NAME=ghreplica
PRTAGS_SCHEMA=prtags
GHREPLICA_SCHEMA=public
DB_MAX_OPEN_CONNS=10
DB_MAX_IDLE_CONNS=5
DB_WORKER_MAX_OPEN_CONNS=3
DB_WORKER_MAX_IDLE_CONNS=1
DB_CONN_MAX_IDLE_TIME=5m
DB_CONN_MAX_LIFETIME=30m
```
