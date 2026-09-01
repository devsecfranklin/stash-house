# SSH Role — lab/franklin Ansible Collection

## Purpose

Deploy, distribute SSH keys, and harden sshd across all lab.bitsmasher.net hosts. The role manages two primary identities:

- **openclaw@chonk** (`id_ed25519_openclaw`) — primary automation identity, deployed to `openclaw` user on all production hosts
- **franklin@bitsmasher.net** (`blowfish-key`) — operator admin key, deployed to `franklin` user on all hosts

## Prerequisites

- Control node runs Ansible 2.19+ (ansible-core)
- Control node's `~/.ssh/id_ed25519_openclaw.pub` and `id_ed25519.pub` are accessible
- Target hosts are reachable via network (static inventory or dynamic inventory)
- The `openclaw` user exists on all target hosts (created by this role if absent)

## File Structure

```
roles/ssh/
├── tasks/
│   └── main.yml        # Main task entry point (runs deploy and verify)
├── templates/
│   └── sshd_config.j2  # Hardened sshd config template
├── defaults/main.yml   # Default variables (keys, users, hardening options)
├── meta/main.yml       # Role metadata (dependency declaration)
└── README.md           # This file
```

## Key Inventory

| Key Label | Comment File | Purpose |
|-----------|-------------|---------|
| openclaw@chonk | `id_ed25519_openclaw.pub` | Primary automation identity — deployed to both `openclaw` and `franklin` users |
| blowfish-key | `id_ed25519.pub` | Operator admin key — deployed to `franklin` user |
| ansible@stargate | `id_ed25519_wonderland.pub` (wonderland only) | Ansible cross-host connectivity key for stargate→wonderland communication |

## Usage

### Deploy keys and harden sshd on a subset of hosts

```bash
ansible-playbook playbooks/purge_openclaw.yml \
  -i inventory/hosts \
  --limit "stargate.research.bitsmasher.net,wonderland.lab.bitsmasher.net"
```

### Syntax check before deploying

```bash
ansible-playbook playbooks/<role>_playbook.yml -i inventory/hosts --syntax-check
```

### Sanity test

```bash
cd collections/ansible_collections/lab/franklin
ansible-test sanity
```

## Hardening Defaults

The role enforces these sshd settings by default:

- `PasswordAuthentication no` (except Jetson subnet: temporarily `yes`)
- `PermitRootLogin prohibit-password`
- `PubkeyAcceptedAlgorithms +ssh-ed25519,sk-ssh-ed25519@openssh.com,rsa-sha2-256,rsa-sha2-512`
- `MaxAuthTries 3`
- `GSSAPIAuthentication no` (enabled only on hosts with Kerberos)
- `LogLevel VERBOSE`, `SyslogFacility AUTH`
- Modern KEX/ciphers/MACs: curve25519-chacha20-poly1305, chacha20-poly1305, AES-GCM

## Deployment Phases

1. **User creation** — Ensures `openclaw` and `franklin` users exist with correct shell/groups
2. **Key distribution** — Copies public keys to `~/.ssh/authorized_keys` for both users
3. **sshd hardening** — Deploys hardened template, restarts sshd service
4. **Verification** — Confirms deployed key fingerprint matches the control node

## Known Issues & Notes

- **OpenBSD hosts**: `AuthorizedKeysFile` directive may need adjustment (`authorized_keys2` variant)
- **Jetson nodes (node900–903)**: Password auth is temporary (`yes`) until keys are confirmed working. Remove `PasswordAuthentication yes` override after key deployment verification.
- **Stargate locked out event (2026-08-29)**: When purging the `openclaw` user, ensure a backup SSH key exists for the `franklin` user before running `pkill -u openclaw`. The purge playbook should always include franklin's key in authorized_keys before killing active sessions.

## Changelog

### 2026-08-29
- Added `blowfish-key` (blowfish@bitsmasher.net) to key inventory
- Added `ansible@stargate` key for stargate→wonderland Ansible connectivity
- Updated user table: `openclaw` purged from stargate; only present on wonderland
- Fixed Jetson password auth note (temporary fallback, not permanent)
- Added locked-out recovery warning for openclaw purge operations

### 2026-08-01 (initial role creation)
- Original sshd hardening template with ed25519 key deployment
- Basic two-user setup: openclaw + franklin
