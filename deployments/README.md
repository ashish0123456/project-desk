# Oracle VM Deployment

This deployment runs the application with Docker Compose on an Oracle Cloud VM. GitLab CI builds one image per service, pushes images to the GitLab Container Registry, and deploys the commit SHA image to the VM over SSH.

## Repository layout

- `deployments/docker/docker-compose.yml`: local development Compose definition.
- `deployments/docker/docker-compose.prod.yml`: production definition; pulls images instead of building them.
- `deployments/docker/.env.example`: template for secrets stored only on the Oracle VM.
- `deployments/scripts/deploy-compose.sh`: remote pull, restart, and cleanup script.

The root `docker-compose.yml` remains a convenience entry point for local development. The top-level `docker` directory has no tracked deployment role and can remain absent.

## Oracle VM setup

Create an Oracle Linux or Ubuntu VM with a public IPv4 address. Install Docker Engine and the Docker Compose plugin, then create the deployment directory:

```bash
sudo mkdir -p /opt/project-management
sudo chown "$USER":"$USER" /opt/project-management
cd /opt/project-management
```

Copy `deployments/docker/.env.example` to `/opt/project-management/.env` and replace these values:

```dotenv
POSTGRES_USER=pm_app
POSTGRES_PASSWORD=REPLACE_WITH_A_LONG_RANDOM_DATABASE_PASSWORD
AUTH_JWT_SECRET=REPLACE_WITH_A_LONG_RANDOM_JWT_SECRET
```

Generate secrets on the VM with:

```bash
openssl rand -base64 32
```

Do not commit the real `.env` file. The production Compose file exposes only the gateway port. Allow inbound TCP `80` and `443` in the Oracle security list and VM firewall. Keep SSH restricted to your own IP where possible. Add Caddy or NGINX in front of port `8080` for HTTPS before sharing the application publicly.

## GitLab CI/CD variables

Add these project variables under **Settings > CI/CD > Variables**. Mark them **Protected**; mark secret values **Masked** where GitLab permits it.

| Variable | Example | Required |
| --- | --- | --- |
| `DEPLOY_HOST` | `203.0.113.10` | Yes, Oracle VM public IP or DNS name |
| `DEPLOY_USER` | `ubuntu` or `opc` | Yes, VM SSH user |
| `DEPLOY_PATH` | `/opt/project-management` | Yes, must match the VM directory |
| `DEPLOY_SSH_PRIVATE_KEY` | contents of the deploy key | Yes, private key matching the VM's `authorized_keys` |
| `SSH_KNOWN_HOSTS` | output of `ssh-keyscan -H your-host` | Recommended |

`CI_REGISTRY`, `CI_REGISTRY_IMAGE`, `CI_REGISTRY_USER`, and `CI_REGISTRY_PASSWORD` are predefined by GitLab and do not need to be created manually.

For `DEPLOY_SSH_PRIVATE_KEY`, use a normal variable containing the complete private key, including the `BEGIN` and `END` lines. Never put that key in `.gitlab-ci.yml`.

For stronger SSH host verification, generate `SSH_KNOWN_HOSTS` from a trusted machine after confirming the fingerprint:

```bash
ssh-keyscan -H YOUR_ORACLE_HOSTNAME
```

## Deployment behavior

- Merge requests and branches run the test matrix.
- Every successful branch pipeline builds and pushes commit-tagged images.
- A successful pipeline on the default branch automatically deploys that commit to production.
- The deploy job uses the immutable `$CI_COMMIT_SHA` tag, so a rollback can deploy a previous commit tag.
- The VM keeps PostgreSQL data in Docker volumes across application updates.

## Manual VM commands

```bash
cd /opt/project-management
docker compose --env-file .env -f docker-compose.yml ps
docker compose --env-file .env -f docker-compose.yml logs -f api-gateway
docker compose --env-file .env -f docker-compose.yml down
```

Do not use `docker compose down -v` unless you intentionally want to delete all PostgreSQL data.
