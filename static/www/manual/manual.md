# Bitsmasher Lab Operations Manual

Plaintext Markdown export compiled from 22 operational chapters.



============================================================
## Ansible Lint
============================================================
## Ansible-Linting Infrastructure


All ansible roles in the lab collection should be linted on a regular basis to maintain configuration quality and catch issues early. The lab uses three tools: ansible-lint, yamllint, and ansible-test sanity.


### Tools and Configuration


Each linter is configured in the test/ directory:

- .ansible-lint.yml — ansible-lint configuration with warnings for deprecated modules and experimental rules disabled
- .yamllint / .yamllintrc — YAML linting extends default rules with lab-specific indentation and spacing constraints


### Coverage Status (44 Roles)


Based on the most recent audit of ansible/collections/ansible\_collections/lab/franklin/roles/ role:

description
    [Complete (30 roles)] apt-mirror, beagleboard, chonk, cluster, common, container-registry, ctfd, desktop, dhcp, dns, docker, documentation, golang, k3s-agent, k3s-server, kerberos, ldap, logging, minecraft, music, nfs, nix, ntp, openbsd, paloalto, prereq, pypi-internal, python, raspberrypi, ssh
    [Minimal/Functional (9 roles)] media, samba, shell, tls, website, apt-mirror-stub, jetson-nano, k8s (partial), dns-stub
    [Stubs (4 roles)] edge, extensions, latex, odroid
description

Roles with the smallest role definitions tend to have the fewest linting issues. The most problematic roles are usually those that include dynamic template generation or custom module invocation.


### Proposed Linting Pipeline


A new test/lint\_all.sh script should be created to run all three linters in sequence:

verbatim
# Example lint_all.sh structure:
#!/bin/bash
set -euo pipefail

cd "$(dirname "$0")/.."

echo "=== ansible-lint ==="
ansible-lint roles/ 2>&1 || echo "ansible-lint warnings found"

echo "=== yamllint ==="
yamllint -c test/.yamllint roles/ 2>&1 || echo "yamllint issues found"

echo "=== ansible-test sanity ==="
ansible-test sanity --docker default -v 2>&1 || echo "sanity check failures"
verbatim

The linting pipeline should be integrated into:

- GitHub Actions: Run on PR and nightly builds (similar to existing bandit/trivy/tfsec workflows)
- Hourly cron job on stargate: Optional for CI-style validation without external triggers


### Linting Exceptions by Role


Some roles require special handling during linting:

- tls: Includes shell commands with kss -kfe which may trigger warning 403 (deprecated module) — explicitly skip in config
- samba: Uses custom smb.conf.j2 templates that require validation via testparm post-converge
- k8s: Contains static manifests rather than role tasks — should be moved to a shared files directory or documented as a reference role


============================================================
## Ansible
============================================================
## Ansible Automation


### Current Landscape (August 2026)


As of August 2026, the Ansible ecosystem is split between two distribution paths:


- Ansible community package -- the traditional monolithic install (pip install ansible). Current version is 13.x (built on ansible-core 2.20). This includes 85+ collections with thousands of modules and plugins bundled together. Follows semantic versioning, releases two major versions per year plus monthly minor releases.
- ansible-core -- the minimal core engine. Maintains three versions at a time (current + two older). Versioned independently from the community package (2.19, 2.20, 2.21). Releases every four weeks for minor updates.


**EOL versions as of August 2026**: Ansible 10--12 and ansible-core 2.15--2.18 are EOL. Do not use these in new work. Only Ansible 13.x / ansible-core 2.20 is current; 14.0.0/2.21 is in development.


### Molecule Testing -- Deprecated (August 2026)


**Molecule has been removed from the active workflow as of August 2026.** It no longer serves as a testing target or standard. All role validation now uses:


- ansible-test --sanity -- built-in collection sanity checks, run natively on stargate.research.bitsmasher.net with zero external overhead
- ansible-test --integration -- integration tests targeting real hosts or containerized environments via ansible-test's native docker driver
- ansible-playbook --syntax-check -- YAML structure validation for all playbooks and role task files


Existing Molecule test harnesses (in ansible/collections/ansible\_collections/lab/franklin/roles/ role /molecule/) remain as historical artifacts but are no longer invoked or maintained. New roles must use native ansible-test exclusively.

lab/franklin Collection

The lab/franklin Ansible collection lives at:
lstlisting[style=mystyle]
/mnt/clusterfs2/workspace/lab-franklin/ansible/collections/ansible_collections/lab/franklin/
    lstlisting

It is the central automation artifact for the entire lab.bitsmasher.net infrastructure, managing provisioning, configuration, and testing across all hosts.


## Collection Manifest (galaxy.yml)


The collection follows standard Ansible Galaxy conventions with namespace lab and name franklin:


- Version: 1.0.0
- License: MIT
- Author: franklin <franklin@bitsmasher.net>
- Repository: https://github.com/devsecfranklin/lab-franklin
- Documentation: https://github.com/devsecfranklin/lab-franklin/docs


## Role Inventory


The collection contains roles organized under the roles/ directory. As of August 2026, the following roles are maintained with native ansible-test harnesses:

lstlisting[style=mystyle]
roles/
├── apt_mirror/          # APT mirror configuration
├── beagleboard/         # BeagleBoard device provisioning
├── common/              # Shared base configuration (bootstrap.sh, common utils)
├── container_registry/  # Docker/container registry setup
├── ctfd/                # CTFd platform deployment
├── desktop/             # Desktop environment provisioning
├── docker/              # Docker engine installation
├── documentation/       # Documentation generation roles
├── dhcp/                # DHCP server configuration (dhcpd.conf)
├── dns/                 # BIND/named DNS server
├── golang/              # Go language toolchain
├── jetson-nano/         # NVIDIA Jetson Nano provisioning
├── k3s_agent/           # k3s worker node
├── k3s_server/          # k3s control plane server
├── k8s/                 # Kubernetes generic setup
├── latex/               # LaTeX environment
├── logging/             # Centralized logging
├── music/               # Music service configuration
├── nfs/                 # NFS server/client
├── nix/                 # Nix package manager
├── openbsd/             # OpenBSD host bootstrap (blowfish target)
│   ├── tasks/main.yml   # Includes files.yml, timezone setup
│   └── files/           # bootstrap.sh, login.conf, hosts entries
├── paloalto/            # Palo Alto firewall config
├── prereq/              # Prerequisite packages
├── pypi_internal/       # Internal PyPI repository setup
│   ├── tasks/main.yml   # Symlink pypi_dir -> html_dir/pypi
│   └── vars/main.yml    # html_dir=/var/www/html, pypi_dir=/mnt/storage1/LAB/pypi
├── python/              # Python environment
├── raspberrypi/         # Raspberry Pi provisioning
├── samba/               # Samba file sharing
├── security/            # Hardening and security baselines
├── shell/               # Shell configuration
├── ssh/                 # SSH server hardening
├── tls/                 # TLS certificate management
├── chonk/               # chonk-specific configuration
├── cluster/             # Cluster-wide settings
    lstlisting


## Workspace Structure


The full project layout:
lstlisting[style=mystyle]
/workspace/lab-franklin/
├── Makefile.am                  # Top-level autotools target (Python, Docker, dev)
├── configure.ac                 # Autotools config (autoconf)
├── bootstrap.sh                 # Cross-platform bootstrapper
├── network_update.sh            # One-shot maintenance script (hardened 2026-08-02)
├── ansible/                     # ANSIBLE_HOME root
│   ├── playbook.yml             # Main playbook
│   ├── hosts                    # Inventory file
│   └── collections/             # Collections path
│       ├── ansible_collections/ # Galaxy namespace
│           └── lab/franklin/    # The lab/franklin collection
│               ├── galaxy.yml
│               ├── roles/       # All role modules (37 total)
│               └── docs/latex/  # LaTeX docs build (autotools)
├── container/                   # Docker/Podman configs
├── terraform/                   # Terraform infrastructure definitions
├── bin/                         # Shared utility scripts (common.sh)
└── docs/manual/                 # Lab manual documentation (this document)
    lstlisting


## Build System (Autotools)


The project uses GNU Autotools for the root-level build:


- configure.ac: Defines package (lab-franklin v0.2), libtool support, Python $$ 3.9 requirement, podman/ansible/gpg checks, git version tagging
- Makefile.am: Top-level targets include test (ansible-lint + ansible-test sanity/integration), security (venv bootstrap), dev (Python venv in BUILDDIR)
- bootstrap.sh: Cross-platform bootstrapper that detects OS (Debian/RedHat/OpenBSD/macOS/Linux) and installs platform-specific packages, then runs aclocal $$ autoreconf $$ automake $$ configure
- docs/manual/Makefile.am: Standalone autotools stub for LaTeX docs build (clean target removes *.aux files)


The network\_update.sh script was hardened on 2026-08-02 with: set -euo pipefail, auto-resolved paths, pinned K3s version validation, ANSIBLE_HOME environment checks, removed dangerous clush/apt-get blast radius, and dead code removal. It runs the main playbook at ansible/playbook.yml as its primary action.

Testing from Stargate


## Native ansible-test on stargate.research.bitsmasher.net


The stargate host (10.10.16.66, Debian 12 bookworm) serves as the exclusive testing platform for the lab/franklin Ansible collection. All test execution runs with zero external billing -- direct shell access to the gateway sandbox eliminates API token overhead.


### Testing Standard: ansible-test (August 2026)


Native ansible-test has replaced Molecule as the canonical testing framework:


- sanity checks: Built into ansible-core; validates module/plugin imports, YAML structure, and collection metadata without any container dependency
- integration tests: Run via ansible-test integration targeting real hosts or the native docker driver (preferred over podman for consistency)
- syntax validation: ansible-playbook --syntax-check validates all playbooks and role task files before execution


The legacy Molecule directory structure at each role's molecule/ subdirectory has been deprecated. It remains in-place as historical reference but is excluded from CI pipelines and the Makefile.am test target.


### Test Execution Flow (ansible-test)


For a given role, the native test cycle is:

- syntax-check: ansible-playbook --syntax-check playbook.yml validates YAML structure
- sanity: ansible-test sanity runs built-in linting against the collection -- no container, no overhead
- integration: ansible-test integration <role> runs target-specific tests using docker driver containers


### Running Tests from Stargate


From stargate's workspace:
lstlisting[style=mystyle]
cd ~/workspace/lab-franklin/ansible/collections/ansible_collections/lab/franklin
ansible-test sanity
ansible-test integration dns --python 3.12
ansible-playbook --syntax-check playbook.yml
    lstlisting


## Dynamic Scoping and Pathing Convention


Roles use dynamic file inclusion via templated variables to avoid hardcoded paths:


- Roles include task fragments using: {\_role\_name\_files} and {\_role\_name\_templates} for dynamically scoped file/paths
- Template references follow the pattern: \{include\_task "tasks/\_role\_name.yml"\ } scoped within dynamically included tasks


This convention eliminates static path dependencies across role invocations and supports multi-target deployment without playbook-level path overrides.


## Package Standards -- Deprecated Utilities


The following deprecated utilities have been removed from baseline role manifests:


- neofetch: Removed from common role; replaced by standard Linux tools (lsb\_release, uname -r)
- tripwire: Removed from security role; no longer maintained and superseded by native file integrity monitoring
- Other legacy packages audited during August 2026 cleanup -- verify role manifests before adding new dependencies


## Dynamic Probing Standard


Static inventory reachability assumptions have been replaced with non-interactive batch probes:

lstlisting[style=mystyle]
ssh -o BatchMode=yes -o ConnectTimeout=3 -i ~/.ssh/id_ed25519_openclaw franklin@<host> "hostname && whoami"
    lstlisting

This approach eliminates stale host entries, prevents password prompts in automation pipelines, and provides deterministic connectivity feedback with a 3-second timeout. Never assume reachability from static notes -- always probe on demand.


============================================================
## Cluster
============================================================
## Ansible Workspace


All Ansible source code, playbooks, roles, and testing are centralized on **stargate.research.bitsmasher.net**:


- Primary workspace: ~/workspace/lab-franklin/ansible/ on stargate
- Secondary (chonk): /mnt/backup1/workspace/lab-franklin/ansible/ (mirror for local builds)
- Interpreter: /home/franklin/.local/bin/ansible on stargate
- Testing standard: Native ansible-test (sanity + integration) -- Molecule deprecated as of August 2026


## Cluster Architecture


The lab uses cluster-wide automation via ansible-playbook targeting inventory groups defined in ansible/hosts. Cluster management is handled through role-based provisioning rather than a dedicated orchestration layer.


### k3s Clustershell (Historical)


Previous cluster shell groupings remain documented for reference:
verbatim
all     -- all managed nodes
head    -- control plane nodes
compute -- worker nodes
gpu     -- GPU-enabled Jetson nodes
    verbatim

These groupings are now managed through Ansible inventory patterns rather than clustershell.


### Centralized Ansible Workflow


- Role changes committed to stargate workspace
- Native ansible-test runs on stargate (zero external cost)
- Playbook execution targets remote hosts via SSH batch probes
- Changes verified via dynamic reachability: ssh -o BatchMode=yes -o ConnectTimeout=3 franklin@ host


============================================================
## Database
============================================================
## Blowfish -- Database Host


### Host Identity


- Hostname: blowfish.lab.bitsmasher.net (alias: blowfish)
- IP: 10.10.12.15
- MAC: 98:b7:85:21:ad:77 (from DHCP reservation)
- OS: OpenBSD (planned; current DNS/DHCP entry references it as an OpenBSD host)


### Role in the Lab


Blowfish is provisioned as a **database host** within the lab infrastructure. It was intended to serve as a dedicated database server, separate from the primary infrastructure hosts (chonk, stargate, skynet).

The OpenBSD role in the lab-franklin Ansible collection includes configuration references specific to blowfish:

- Login configuration uses blowfish password hashing algorithm (localcipher=blowfish,a) in /etc/login.conf.
- The bootstrap.sh for OpenBSD hosts includes a direct /etc/hosts entry for blowfish (10.10.12.15).
- DHCP reservation is configured in the lab's dhcpd.conf role.


### Current Status


Blowfish has been **unreachable** from chonk since the initial infrastructure audit. The IP 10.10.12.15 resolved from DNS (stargate's /etc/hosts), but authentication was denied -- no matching SSH key was authorized on that host.

Key references:

- network\_update.sh: blowfish is the default target for the openbsd\_bootstrap() function.
- OpenBSD role tasks: includes commented-out ansible setup command for blowfish (indicating the host was expected to be online).


The host may need network investigation -- it resolved from DNS, had a DHCP reservation, and MAC address on file, suggesting it was once part of the lab but has since gone offline or been decommissioned.


### Planned Configuration (OpenBSD)


When brought online, the expected setup for blowfish is:


- SSH key-based authentication via ~/.ssh/authorized\_keys
- OpenBSD's built-in pf firewall with custom rules
- Database software stack (likely PostgreSQL or MariaDB on OpenBSD)
- Ansible-managed configuration via the lab-franklin collection, using the openbsd role as the base profile.


============================================================
## Direnv
============================================================
## direnv


============================================================
## Disaster Recovery
============================================================
## Disaster Recovery


### Critical Dependencies Chain


The lab infrastructure has a cascading dependency chain. A failure in one layer propagates:


- Time (NTP) $$ if time drift exceeds 128 seconds, Kerberos tickets expire and TLS handshakes fail
- DNS (ns1/BIND) $$ without DNS resolution, all hostname-based services (Kerberos SRV records, NTP upstream, LDAP) break
- Kerberos (odroid-c1/KDC) $$ without the KDC, authentication fails for all Kerberos-authenticated services
- LDAP (bbb1/slapd) $$ without LDAP, directory lookups fail and user/service account information is unavailable


Recovery priority should follow this chain in reverse: fix time first, then DNS, then KDC, then LDAP.


### Failure Scenarios and Recovery


#### NTP Failure -- All Hosts Lose Time Sync


Symptoms: timedatectl shows "System clock synchronized: no", Kerberos tickets rejected.

Recovery:

- Verify GPS unit is locked on odroid-c1 (time host)
- Check ntpd/ntpsec service: systemctl status ntpsec
- If GPS is offline, enable holdover mode (internal oscillator continues at reduced accuracy)
- On affected clients, force resync: ntpdate -s time.lab.bitsmasher.net or wait for ntpd's step-sync to kick in
- Verify with ntpq -p on clients and ntptime on the server


#### DNS (ns1) Offline -- Cascade Failure


Symptoms: hostnames don't resolve, Kerberos SRV discovery fails.

Recovery:

- Check ns1 reachability: ping from any known-good host
- If ns1 is up but not answering: check BIND (named) service status
- If BIND has crashed due to config error: restore from last known-good named.conf
- As a workaround, add manual entries to /etc/hosts on affected hosts until DNS recovers
- Update krb5.conf with explicit kdc = odroid-c1.lab.bitsmasher.net


#### LDAP (bbb1) Down -- No Directory Service


Symptoms: slapd failed, status.sh reports both slapd and ldapsearch as failing.

Recovery:

- SSH to bbb1/lab.bitsmasher.net as root
- Check TLS certificate validity in the configured cert path
- Regenerate or replace expired certificates
- Restart slapd: systemctl restart slapd
- Verify with ldapsearch for franklin and sly DNs


#### KDC (odroid-c1) Offline -- No Authentication


Symptoms: kinit fails for all realms, service authentication denied.

Recovery:

- SSH to odroid-c1 as root (franklin user's key is not authorized)
- Check kdc process: kadmind/krb5kdc status
- Verify KDC database integrity: kdb5_util list
- Restart kerberos services if needed
- Regenerate host keytabs for affected services via kadmin.local


### Backup and Restoration


- Ansible playbooks (lab-franklin collection) are the single source of truth -- they can reconstruct any configured state
- Bare git repos with GPG-encrypted pass entries handle credential backup
- KDC database (\$KRB5\_KDB\_FILE) should be backed up regularly on odroid-c1
- DNS zone files should have local copies (in the DNS role's files/ directory)


### Prevention: Monitoring


The current monitoring approach is manual via status.sh. For improved resilience, consider:

- Automated ping/SSH checks on critical hosts in HEARTBEAT.md
- Nagios/Zabbix-style monitoring on a dedicated host
- Email/slack alerts for service down conditions


The lab-franklin Ansible collection should eventually include a monitoring role to automate this.


============================================================
## Dns
============================================================
## DNS Configuration


### BIND Overview -- Post-Migration (August 2026)


The DNS infrastructure has been migrated from the legacy server1 (10.10.12.12) to a new authoritative setup:


- DNS Master: node3 (10.10.12.3) -- promoted from client node to authoritative BIND9 master
- Legacy DNS host: server1 (10.10.12.12) has been decommissioned; all zone transfers and updates now target node3


The DNS role in lab-franklin's Ansible collection manages:

- Forward zone: db.home.lab (hostnames to IPs)
- Reverse zones: PTR records for each 10.10.x.0/24 subnet
- SRV records: Kerberos service discovery
- Zone file generation via Jinja2 templates on node3


### Zone and Host Realignment (August 2026)


The following entries have been updated in the forward zone (db.home.lab):

center
tabular{lll}

**Hostname** & **IP Address** & **Notes** 
 
stargate.research.bitsmasher.net & 10.10.16.66 & Ansible orchestration host, primary workspace 
chonk.lab.bitsmasher.net & 10.10.8.60 & Gateway host, primary compute + inference 
 
tabular
center

The corresponding reverse zones (PTR records) have been updated to match:

- 10.10.16.66 $$ stargate.research.bitsmasher.net
- 10.10.8.60 $$ chonk.lab.bitsmasher.net


### Kerberos SRV RFC Compliance


The following SRV records now point directly to FQDN A records rather than CNAMEs, per RFC 2782:

lstlisting[style=mystyle]
_kerberos-adm._tcp.lab.bitsmasher.net IN SRV 0 100 749 kdc1.lab.bitsmasher.net.
_kerberos._tcp.lab.bitsmasher.net     IN SRV 0 0 88  kdc1.lab.bitsmasher.net.
_kerberos._udp.lab.bitsmasher.net     IN SRV 0 0 88  kdc1.lab.bitsmasher.net.
_kdc._tcp.lab.bitsmasher.net          IN SRV 0 0 88  kdc1.lab.bitsmasher.net.
_kdc._udp.lab.bitsmasher.net          IN SRV 0 0 88  kdc1.lab.bitsmasher.net.
    lstlisting

**Note**: The _kerberos-adm.\_tcp record now points directly to the A record for kdc1.lab.bitsmasher.net rather than via a CNAME chain, eliminating resolution latency and failure modes from indirection.


### DNS Failure Impact


When the DNS master (node3) is offline:

- Lab hosts fall back to /etc/hosts entries (authoritative but static -- no dynamic updates)
- Kerberos KDC discovery via SRV records fails -- clients need manual kdc = kdc1.lab.bitsmasher.net in /etc/krb5.conf
- Dynamic zone updates cannot propagate; new hosts require manual /etc/hosts sync


### DNS Maintenance Procedures


Zone file updates are committed to the DNS Ansible role and deployed via:
lstlisting[style=mystyle]
ansible-playbook -l dns_servers -t dns deploy.yml
    lstlisting

Forward and reverse zone files should be audited quarterly for IP realignment accuracy.


============================================================
## Hardware
============================================================
## Hardware Inventory


### Desk Setup


- Franklin's desk: chonk (10.10.8.60) -- GAMING Windows PC, primary workstation, Gateway host running OpenClaw
- Sly's desk: GAMING Windows computer (Twitch: https://www.twitch.tv/s1y_b0rg) -- at main office near music.lan monitor
- Both sit at the main office where music.lan's monitor is displayed


### Server and Compute Hosts


*{Core Infrastructure (Debian 12)}

tabularx{}{l l l X}
**Host** & **IP** & **User** & **Role** 

chonk & 10.10.8.60 & franklin/openclaw & Gateway host, primary workstation + inference 
stargate & 10.10.16.66 & openclaw & Ansible workspace, build host, DNS master (node3), NFS server 
skynet & 10.10.16.10 & openclaw/franklin & Minecraft server, IP-direct login 
time & 10.10.12.2 & openclaw/franklin & Stratum-1 NTP server with GPS reference 
wonderland & 178.62.60.55 & franklin/openclaw & Public cloud VM (Debian 12, 6.1 kernel), web host 
music.lan & 192.168.86.38 & root & Desktop machine via skynet route 
node3 & 10.10.12.3 & openclaw & **Authoritative BIND9 DNS master** (post-server1 migration) 

tabularx

*{Jetson Cluster ( 87 days uptime)}

All Jetsons run Ubuntu 18.04 with ed25519 key auth (openclaw permanent, franklin has password 123 as fallback -- pending key-only migration):

tabularx{}{l l X}
**Host** & **IP** & **Status** 

node900 & 10.10.12.90 & Online, key auth; home dir from stargate:/mnt/clusterfs2 
node901 & 10.10.12.91 & Online, key auth; home dir from stargate:/mnt/clusterfs2 
node903 & 10.10.12.93 & Online, key auth; home dir from stargate:/mnt/clusterfs2 
node902 & 10.10.12.92 & Off / no route to host (likely powered down) 

tabularx

**NFS note**: Jetson home directories are NFS-mounted from stargate:/mnt/clusterfs2. StrictModes permissions apply (home 0750, .ssh 0700, authorized\_keys 0600).

*{KDC and Directory Services}

tabularx{}{l l X}
**Host** & **IP** & **Status** 

odroid-c1 / kdc1 & 10.10.12.254 & Online (~242d uptime), root key auth, KDC for lab.bitsmasher.net Kerberos realm 
ldap/bbb1 & 10.10.13.1 & SSH reachable as root; slapd failed Dec 2025 (TLS cert issue) 

tabularx


### Historical / Decommissioned Hosts


- server1 (10.10.12.12): Former DNS master -- decommissioned August 2026; role replaced by node3
- snowy: Hard disk physically removed; /mnt/snowy now just a raw mounted drive (not SSH-accessible)
- blowfish: Previously at 10.10.12.15, reassigned to 10.10.14.85. Key not authorized yet -- unreachable from chonk but port 22 is open. Expected role: database host with OpenBSD
- node902: Powered off; no route to host from any known peer


### Planned Hardware Expansion


- blowfish as the primary database server (OpenBSD planned)
- Full ed25519 key rollout for node900--node903 to remove password fallback
- k3s\_server and k3s\_agent roles for Kubernetes across node infrastructure


============================================================
## History
============================================================
History

Code Storage

main repos

workspace repo


============================================================
## Kerberos
============================================================
## Kerberos Authentication


### Realm Configuration


- Realm name: lab.bitsmasher.net
- KDC host: odroid-c1.lab.bitsmasher.net (10.10.12.254)
- KDC software: MIT Kerberos KDC running on odroid-c1
- Uptime: ~242 days (stable, rarely rebooted)


### Access and Key Management


The KDC only has SSH key auth configured for the root user. The franklin user has been rejected by the KDC's SSH configuration -- direct root access is required for any management tasks.

Key file: /.ssh/id\_ed25519\_openclaw (the standard openclaw key).


### Principal Conventions


Kerberos principals in the lab.bitsmasher.net realm should follow these conventions:

- Service principals: service/FQDN@lab.BITSMASHER.NET
- User principals: username@lab.BITSMASHER.NET
- Host principals: host/FQDN@lab.BITSMASHER.NET


### DNS/Kerberos Integration -- RFC Compliant (August 2026)


Kerberos relies on DNS SRV records for KDC discovery. The following records point directly to FQDN A records (no CNAME indirection):

lstlisting[style=mystyle]
_kerberos-adm._tcp.lab.bitsmasher.net IN SRV 0 100 749 kdc1.lab.bitsmasher.net.
_kerberos._tcp.lab.bitsmasher.net     IN SRV 0 0 88  kdc1.lab.bitsmasher.net.
_kerberos._udp.lab.bitsmasher.net     IN SRV 0 0 88  kdc1.lab.bitsmasher.net.
_kdc._tcp.lab.bitsmasher.net          IN SRV 0 0 88  kdc1.lab.bitsmasher.net.
_kdc._udp.lab.bitsmasher.net          IN SRV 0 0 88  kdc1.lab.bitsmasher.net.
    lstlisting

The _kerberos-adm.\_tcp record now points directly to the A record for kdc1.lab.bitsmasher.net, per RFC 2782 compliance.

These records are authoritative on node3 (10.10.12.3), the new DNS master following the server1 decommissioning. If DNS is unavailable, clients fall back to manual kdc = kdc1.lab.bitsmasher.net in /etc/krb5.conf.


### Keytab Management


Keytabs for service principals should be:

- Generated on the KDC host (odroid-c1) using kadmin.local or kadmin -q
- Transferred securely to target hosts (via scp with SSH keys)
- Set to mode 600 and owned by the appropriate service user


**Note**: Kerberos keytabs have not been deployed across the infrastructure. NFS exports with sec=krb5i are currently non-functional pending keytab deployment. The Ansible ssh role handles the StrictModes requirements for SSH-based keytab transfer.


============================================================
## Kubernetes
============================================================
## Kubernetes (k3s) Cluster


### Architecture


The lab uses k3s (Lightweight Kubernetes from Rancher) for container orchestration across the infrastructure:


- k3s\_server: Master node running k3s control plane (host: head2, referenced but unreachable in current audit)
- k3s\_agent: Worker nodes that register with the server and receive workloads
- GPU nodes: Specialized workers for NVIDIA container runtime and compute workloads


### Version Pinning


k3s version is pinned to the 2.1.x series (specifically 2.1.5 as of the latest network\_update.sh). This prevents surprise breakage from upstream k3s releases. The pinning policy follows:
lstlisting[style=mystyle]
K3S_VERSION="2.1.5"
curl -sfL https://get.k3s.io | sh -s - \
    --write-kubeconfig-mode 644
    lstlisting


### Container Runtime


k3s uses containerd as the default runtime, with NVIDIA Container Toolkit for GPU workloads:

- nvidia-container-toolkit: Installed on GPU nodes via the network\_update.sh automation
- containerd config: Located at /var/lib/rancher/k3s/agent/etc/containerd/config.toml
- NVIDIA runtime configured in containerd to expose GPU devices to containers


### Cluster Topology


tabularx{}{l l X}
**Node Role** & **Target Hosts** & **Notes** 

k3s\_server & head2 (planned) & Control plane; not yet online
k3s\_agent & node900, node901, node903 & Jetson nodes as workers
nvidia\_nodes & GPU-equipped hosts & Container Toolkit + k3s agent
tabularx


### kubeconfig Access


After installation, kubeconfig is written to /etc/rancher/k3s/k3s.yaml with mode 644. Cluster access requires:

- SSH key auth to the k3s\_server host
- Reading the kubeconfig file (or copying it securely)
- kubectl --kubeconfig /etc/rancher/k3s/k3s.yaml cluster-info to verify connectivity


The k3s_server role and k3s_agent role in lab-franklin handle the installation and configuration of these nodes via Ansible.


============================================================
## Ldap
============================================================
## LDAP Directory Services


### Server Configuration


- Hostname: ldap.lab.bitsmasher.net (alias: bbb1)
- IP: 10.10.13.1
- Software: OpenLDAP (slapd)
- Status: SSH reachable as root, but slapd service has been failed since December 2025


### Outage Details


The slapd failure manifests as:
lstlisting[style=mystyle]
TLS init def ctx failed (-1)
    lstlisting

This indicates a TLS certificate configuration issue -- likely an expired, missing, or misconfigured CA certificate. The server is listening on its network interface (SSH works for root), but the LDAP daemon itself has not recovered since December 2025.


### Monitoring Approach


Franklin maintains /home/franklin/status.sh which checks:

- slapd service status on bbb1/ldap
- ldapsearch queries for franklin and sly DN lookups -- both fail because the directory is down


The script provides a simple pass/fail indicator for LDAP health without requiring complex monitoring infrastructure.


### Recovery Steps


To restore the LDAP server:

- SSH to ldap.lab.bitsmasher.net as root (ed25519 key)
- Check slapd logs: journalctl -u slapd --since "2025-12-13"
- Verify TLS certificate chain in /etc/ssl or the configured cert path
- Update or regenerate the CA-signed certificate if expired
- Restart slapd: systemctl restart slapd
- Verify with ldapsearch for franklin and sly DNs


### Directory Structure (Expected)


The OpenLDAP directory should contain entries for:

- User entries for franklin, Sly (slyborg), and other lab members
- Group/organizational unit entries for lab access control
- Service accounts for LDAP-authorized services (Kerberos integration, web auth)


Once restored, the LDAP directory serves as the central authentication source for the lab, complementing Kerberos for identity management.


============================================================
## Minecraft
============================================================
## Minecraft Server


### Host and Location


The Minecraft server runs on **skynet.lab.bitsmasher.net** (10.10.16.10), one of the core lab infrastructure hosts alongside stargate and chonk. The Minecraft role is managed by the minecraft Ansible role in the lab-franklin collection.

Skynet's dual-role setup:

- Primary: Minecraft game server host
- Secondary: General-purpose infrastructure node (Ansible workspace, CI testing)


Note: skynet uses IP-based SSH login (10.10.16.10) rather than hostname -- the hostname resolves oddly in SSH config due to a double-domain suffix issue.


### World and Chunk Storage


Minecraft world data is stored on NFS-mounted storage (clusterfs2). The world files are version-controlled where practical, with the following constraint:

**NO GIT-LFS**: Chunk files and world data must not use git-lfs. The Minecraft repository hit GitHub's 10GB LFS limit in a prior incident where chunk files were held hostage. Use bare git repos for any versioned world data.


### Dreamland Connection


The *Dreamland Manual* referenced in the docs/manual directory covers the Dreamland Minecraft server -- this manual documents that specific server's configuration, players, and game rules. The physical infrastructure (skynet host) is shared between both operations.


### Server Management via Ansible


The minecraft role handles:

- Minecraft server installation and version updates
- World backup procedures (NFS snapshots or bare git repos)
- Player whitelist management
- Mod/plugin configuration if applicable


Molecule testing for the minecraft role confirms correct file placement and configuration generation, though live server connectivity is not verified by the test harness.


============================================================
## Network
============================================================
## Network Infrastructure


### Addressing Scheme


The lab uses private IPv4 addressing across multiple subnets, organized by function:

tabularx{}{l l X}
**Subnet** & **Allocation** & **Purpose** 

10.10.8.0/21 & chonk, chonk-wifi, music.lan (via skynet) & Main office / desk area
10.10.12.0/& time, node90x, odroid-c1, blowfish, ns1 & Core infrastructure
10.10.13.0/ & ldap/bbb1 & Directory services (moved to working hosts)
10.10.14.0/ & blowfish (reassigned) & New allocation for blowfish
10.10.15.0/ & femputer & Guest / temporary hosts
10.10.16.0/& stargate, skynet, edge-t & Lab backbone / staging
tabularx

Note: The full subnet masks vary by role and DHCP scope configuration on the Dream Machine (DHCP server).


### DHCP


The dhcp Ansible role manages the DHCP server configuration stored at roles/dhcp/files/dhcpd.conf. Key reserved hosts include:


- chonk: 10.10.8.60, MAC fc:9d:05:01:27:02 -- Franklin's desk
- chonk-wifi: 10.10.8.61 -- wireless client
- chonk-ten-gig-1: 10.10.8.?? -- dedicated ten-gig link
- blowfish: 10.10.8.15 / 10.10.12.15, MAC 98:b7:85:21:ad:77 -- database host
- femputer: 10.10.15.1, MAC d8:bb:c1:af:21:17 -- temporary/guest


### DNS


The DNS infrastructure is managed by the dns Ansible role (BIND/named). Zones maintained include:

- Forward zone: db.home.lab
- Reverse zones for each subnet
- A records and PTR records for all lab hosts


DNS role molecule testing confirms BIND packages install correctly and named.conf.options are written. The DNS server (ns1) has been offline during recent audits, which explains why many hosts in the 10.10.x.0/24 range cannot resolve -- ns1 is a single point of failure for lab DNS resolution.


### Firewall and Access Control


- The Dream Machine provides built-in DHCP and handles basic routing between subnets
- odroid-c1 serves as the KDC and likely has pf firewall rules managing cross-subnet Kerberos traffic
- OpenBSD hosts (blowfish planned) use pf for host-level firewalls
- stargate has configured SSH access from chonk via openclaw key pair


### Physical Topology


The lab spans multiple physical locations and subnets:

- Main office where Franklin and Sly have their desks monitors (music.lan)
- Server racks hosting core infrastructure (time, odroid-c1/KDC)
- Jetson cluster (node900--node903) for GPU compute
- Public cloud (wonderland/178.62.60.55) for external-facing services


Routing between the 192.168.86.x subnet (music) and the rest of the lab requires hopping through skynet due to the non-standard routing configuration.


============================================================
## Nfs Testing
============================================================
## Overview


The NFS role at ansible/collections/ansible\_collections/lab/franklin/roles/nfs has been migrated from Molecule to native ansible-test. The role validates:


- Server-side export configuration (/etc/exports)
- Client-side fstab management and mount points
- Kerberos security validation (sec=krb5i -- though keytabs remain undeployed; sec=krb5i is non-functional without keytab infrastructure)


## Testing with ansible-test (August 2026)


From stargate:
verbatim
cd ~/workspace/lab-franklin/ansible/collections/ansible_collections/lab/franklin
ansible-test sanity --python 3.12
ansible-test integration nfs --python 3.12
    verbatim

Molecule tests remain as historical artifacts only and are no longer invoked.


## Current Validation Focus


Testing now emphasizes:


- Export template rendering: Jinja2 templates generate valid /etc/exports content for active roles (chonk's /mnt/storage1, storage2, storage3 exports)
- Client-side mount validation: fstab line insertion and idmapd.conf deployment on client nodes
- Jetson home directory consistency: NFS-mounted homes from stargate:/mnt/clusterfs2 must maintain uid/gid alignment


## StrictModes Validation


NFS-mounted user homes require OpenSSH StrictModes compliance:

- Home directory 0750, .ssh 0700, authorized\_keys 0600 -- verified as part of integration tests
- Kerberos keytabs never deployed; sec=krb5i in exports vars is non-functional


## Integration Test Targets


For thorough testing, add:


- Export verification (check showmount -e)
- Mount/unmount cycle tests on client container
- StrictModes permission validation post-mount


## Molecule Test Checklist -- DEPRECATED


The following Molecule-specific items are no longer relevant as of August 2026:

- [DEPRECATED] Server-side converge with hostname override
- [DEPRECATED] Client-side converge (mount directory creation)
- [DEPRECATED] Verify phase checks exports file content
- [DEPRECATED] Idempotence test: second converge


All future work uses ansible-test exclusively.


============================================================
## Pypi
============================================================
## Internal PyPI Repository


The lab runs an internal PyPI package repository to provide Python packages to development machines without depending on external PyPI servers. This avoids network bottlenecks, supports offline development, and enables use of proprietary or lab-specific packages not available on PyPI.org.


### Architecture


The PyPI directory lives at /mnt/storage1/LAB/pypi. A symlink is created from the web server's document root (/var/www/html/pypi) pointing to this directory, making packages accessible via HTTP.

Access URL pattern:
lstlisting[style=mystyle]
http://<storage-host>/pypi/simple/
    lstlisting

This follows the PEP 503 simple repository API structure expected by pip.


### Configuration via Ansible


The pypi\_internal role in the lab/franklin collection automates this setup:


- Role path: roles/pypi\_internal/
- Variables:
- html\_dir: web root (/var/www/html)
- pypi\_dir: package directory (/mnt/storage1/LAB/pypi)

    
     **Task**: Creates a symbolic link from \$html\_dir/pypi $$ \$pypi\_dir
    
     **Playbook integration**: Included in the cluster/storage playbook alongside NFS, LDAP, Kerberos, DNS, and ntp roles
itemize


### Usage from Client Machines


Configure pip to use the internal PyPI server:
lstlisting[style=mystyle]
# In ~/.pip/pip.conf or /etc/pip.conf:
[global]
index-url = http://<storage-host>/pypi/simple/

# Optionally trust host if not using HTTPS
trusted-host = <storage-host>
    lstlisting

Or use pip's \-\-index\-url flag directly:
lstlisting[style=mystyle]
pip install --index-url http://<storage-host>/pypi/simple/ some-package
    lstlisting


### Package Maintenance


To add packages to the internal PyPI, place them in the pypi directory with PEP 503-compliant directory structure:

- Package directories (lowercase names)
- File links inside each package directory pointing to the actual wheel/sdist files
- The web server must have read access to /mnt/storage1/LAB/pypi


Packages can be built and added using standard pip tools:
lstlisting[style=mystyle]
pip install --no-index --find-links=/mnt/storage1/LAB/pypi/simple/ <package>
    lstlisting


============================================================
## Revision Control
============================================================
Code Storage

main repos

workspace repo


============================================================
## Security
============================================================
Security
ch:security

The lab.bitsmasher.net infrastructure relies on multiple security layers -- from SSH key-based authentication to automated vulnerability scanning. This chapter documents the security posture, procedures, and policies in place.


## Authentication and Access Control


### SSH Key Management

The lab uses ed25519 SSH keys for host authentication. All automation hosts share a common key pair:


- Private key: ~/.ssh/id\_ed25519\_openclaw (chonk)
- Public key: ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIJb9mV02PpxD8VpzYCnu7192dTwMnGWUc3qh55BIeElN openclaw@chonk
- Key users: ``openclaw'' user on all production hosts


### SSH Hardening and Drop-In Standards


All hosts use standardized hardening.conf drop-ins in /etc/ssh/sshd\_config.d/:


- PermitRootLogin prohibit-password -- root login allowed only via key, never password
- StrictModes yes -- OpenSSH enforces file permission checks on ~/.ssh and authorized\_keys
- Explicit user key sync -- ed25519 keys deployed to target hosts via ansible ssh role; no manual key management


### Password Policy -- Subnet Exceptions


Passwords are used only as a temporary fallback. The following hosts have documented exceptions:


- Jetson nodes (node900--node903): Temporary password authentication allowed pending ed25519 key rollout. Password is 123 for the franklin user. These nodes will be migrated to key-only auth when infrastructure permits.


All other hosts enforce key-only authentication. No new password-based access may be granted without operator approval.


### User Accounts


center
tabular{lll}
  
  **User** & **Purpose** & **Access Scope** 
 
  franklin & Primary admin user & All hosts; sudo NOPASSWD on skynet 
  openclaw & Automation user & All production hosts via ed25519 key 
  root & Emergency access & KDC, ldap, music host only; key auth only 
 
tabular
center


## Network Security


### Subnet Segmentation

The lab is divided into dedicated subnets providing implicit isolation:


- Main office: 10.10.8.0/21 -- core infrastructure (chonk, stargate)
- Backbone: 10.10.16.0/24 -- game server, DNS services
- Edge/specialized: 10.10.12.x -- Jetson nodes, NTP, KDC, edge devices
- Service zone: 10.10.13.0/24 -- LDAP directory services
- Storage zone: 10.10.14.0/24 -- Blowfish and shared storage hosts
- Non-standard: 192.168.86.x -- Music desktop (routed through skynet)


### Firewall Rules

The Ansible ``prereq'' role handles firewall exceptions for K3s API ports (6443, 2379--2381). Both UFW (Debian) and Firewalld (RedHat-family) are supported. Default policy is deny-all inbound with only explicitly allowed services accessible.


## Certificate and TLS Management

The ``tls'' role manages certificate provisioning across the lab:

- Certificates are stored on NFS (clusterfs2) for cross-host availability
- OpenLDAP server experienced a TLS failure in December 2025 due to an expired certificate -- this is the primary remaining TLS vulnerability
- All external-facing services use HTTPS with valid certificates where applicable


## Credential Management


### Pass + GPG (Stash House)

Credentials are stored using pass with GPG encryption. The Stash House project manages encrypted credential backups synced via bare git repositories:

- Credentials never exist in plaintext in any repository
- Git repos use bare repositories to store only the encrypted .gpg files
- No git-lfs is used for credential data (see LFS policy)


### Ansible Vault

Sensitive variables (firewall credentials, API keys) are stored using ansible-vault. Vault files should never be committed to version control. The vault passphrase is managed outside the repository.


## Backup and Disaster Recovery


### Data Backups

NFS-mounted clusterfs2 provides shared storage with redundant copies across hosts. DNS zones are version-controlled in BIND configuration files. Credential data uses GPG-encrypted pass entries synced via bare git repos.


### Service Availability Status


center
tabular{llll}
  
  **Service** & **Host** & **Status** & **Last Known Good** 
 
  NTP stratum-1 & time.lab.bitsmasher.net & Online & Current 
  KDC (Kerberos) & odroid-c1.lab.bitsmasher.net & Online & ~242 days uptime 
  DNS (BIND) & node3.lab.bitsmasher.net & Operational & August 2026 
  OpenLDAP & ldap.lab.bitsmasher.net & Offline & December 2025 
  Minecraft & skynet.lab.bitsmasher.net & Online & Current 
 
tabular
center


### Disaster Recovery Priorities


- NTP (time host) -- time sync failure cascades to all Kerberos authentication
- DNS (node3, BIND) -- resolution failure makes most hosts unreachable by name
- Kerberos KDC -- authentication cascade failure affects all service access
- LDAP -- directory services, offline since December 2025


## Container and Kubernetes Security


The k3s cluster is hardened with:
   Version pinning to 2.1.x
   Containerd runtime with NVIDIA GPU support
   TLS for all API communication (k3s default)
   Pod security policies via Kubernetes native admission controllers
itemize


## Automated Security Scanning

The .github directory contains the following automated workflows:

center
tabular{lcl}
  
  **Workflow** & **Tool** & **Purpose** 
 
  bandit.yml & Bandit & Python static analysis for security vulnerabilities 
  tfsec-sarif.yml & tfsec & Terraform IaC security scanning 
  trivy.yaml / trivy-cb.yml & Trivy & Container and filesystem vulnerability scanning 
  bash\_check.yaml & GNU Autotools & Build system validation (make security) 
  markdown.yml & markdownlint & Markdown style/lint checking 
 
tabular
center

All scanning workflows run on every pull request and main branch push. Results are uploaded as SARIF artifacts for review.


============================================================
## Storage
============================================================
## Storage Infrastructure


### Active Disk Topology (chonk -- August 2026)


All compiled documents and web artifacts target wonderland; chonk serves as local scratch space. Physical topology on chok:


- /mnt/backup1: Source directory for Ansible repos, workspace code, and build artifacts
- /mnt/clusterfs: Scratch disk for temporary compilation, LaTeX builds, and intermediate files
- /mnt/snowy: Bulk storage -- formerly a full NFS share; hard disk physically removed, now mounted as raw drive only


The following mount points have been purged:

- /mnt/storage1: Completely removed across all roles and /etc/exports on all hosts
- Legacy /mnt/storage2, /mnt/storage3: Purged from inventory references


### NFS Mounts -- Current State


The lab relies on NFS for shared storage. Active mounts as of August 2026:


- stargate:/mnt/clusterfs2: Primary cluster workspace; mounted on stargate and accessible from wonderland via SSH-sync pipelines
- Jetson home directories: Shared from stargate via stargate:/mnt/clusterfs2 -- provides consistent /home for node900--node903


### StrictModes and Permission Requirements


OpenSSH StrictModes enforces strict permission checks on home directories. Required permissions:


- Home directory: 0750 (rwxr-x---)
- /.ssh: 0700 (rwx------)
- /.ssh/authorized\_keys: 0600 (rw-------)


All Jetson nodes must comply with these permissions or SSH authentication fails silently. NFS-mounted homes require uid/gid consistency between stargate and each Jetson host.


### Snapshot and Backup Strategy


The lab uses bare git repos for version control and backup:

- GPG-encrypted repositories synced across hosts for credential protection (Stash House project)
- Pass-based password store with GPG backend
- No LFS -- ABSOLUTE RULE: do not use git-lfs anywhere. Microsoft once held Minecraft chunk files hostage when a repo hit the 10GB GitHub limit with LFS enabled.


### LFS Policy Reminder


**ABSOLUTE RULE: Do NOT use git-lfs anywhere in any lab-franklin repository.**

This applies across all repos: training-ai, minecraft, Ansible collections, everything. The incident with Minecraft chunk files being held hostage by GitHub when the repo hit 10GB with LFS is a real precedent -- never risk it.


============================================================
## Terraform Infra
============================================================
Terraform Infrastructure Management

Terraform configurations in this lab manage cloud infrastructure across two providers: DigitalOcean and Google Cloud Platform. All configurations are version-controlled and include automated security scanning via GitHub Actions.


## DigitalOcean Infrastructure


The terraform/digital\_ocean/ directory manages the external-facing web presence for bitsmasher.net.


### Provisioned Resources


- Droplet: 512MB Debian 12 in lon1 region (www.bitsmasher.net)
- Domain: bitsmasher.net with full DNS management
- DNS Records: A, MX, TXT (SPF/DMARC/DKIM), NS, and ACME challenge records


### Key DNS Records


center
tabular{lll}

Name & Type & Purpose 
 
www.bitsmasher.net & A & Web server host (178.62.60.55) 
@ & MX & mail.protonmail.ch 10 (ProtonMail routing) 
@ & TXT & SPF and DMARC policies 
protonmail.\_domainkey & CNAME & DKIM email verification 
\_acme-challenge.*.lab & TXT & Let's Encrypt certificate challenges 
 
tabular
center


### Configuration Notes


The droplet includes:

- Automated backups enabled (immutable against accidental deletion)
- SSH key injection via terraform data source
- No external monitoring (consider enabling DO agent for uptime tracking)


Provider authentication uses a DigitalOcean API token managed through pass/direnv integration. The token is passed at runtime, never committed to the repository.


## Google Cloud Platform Infrastructure


The terraform/google/ directory manages GCP resources for Stash House free-tier operations.


### Resources Managed


- VPC Network: stash-house-vpc dedicated network
- Compute Instance: e2-micro Debian 12 (Always Free Tier eligible) in us-central1-a
- Firewall Rule: allow-ssh opening port 22 to the VPC
- Billing Budget Alert: Tripwire alert at \$0.01 threshold for free-tier monitoring


The GCP configuration is intentionally minimal, designed around Google's Always Free Tier limits to maintain zero-cost infrastructure.


## Security Scanning


All Terraform configurations are automatically scanned by two GitHub Actions workflows:

- tfsec-sarif.yml — IaC security scanning (generates SARIF reports)
- tfsec-comment.yml — Posts security findings as PR comments


Additionally, the DigitalOcean directory includes existing README.md with doctl and pass integration instructions for manual operations.


## Testing


Terraform configurations are validated through:

- Go-based integration tests (terraform/test/main\_test.go)
- Plan dry-run validation via GitHub Actions (markdown.yml workflow)
- Manual state reconciliation checks (compare terraform state against live infrastructure)


State drift detection is recommended as a regular maintenance task — compare terraform state list output against actual cloud resource inventories quarterly.


============================================================
## Time
============================================================
## Network Time (NTP)


### Overview


Time synchronization is one of the most critical and least-visible aspects of network management. All lab services -- Kerberos authentication, TLS certificate validation, log correlation, distributed databases, and filesystem consistency -- depend on time accuracy. The lab uses a two-tier NTP architecture:


- Stratum 1 source: GPS-referenced clock on the time host (10.10.12.2) provides the lab's authoritative time reference.
- Stratum 2+ clients: All other hosts sync to time.lab.bitsmasher.net.


### The Time Host


- Hostname: time.lab.bitsmasher.net (alias: time)
- IP: 10.10.12.2
- User: openclaw
- OS: Debian 12 bookworm
- Role: Stratum-1 NTP server / reference clock source for the lab


The time host has a GPS unit connected (serial/tty interface) providing UTC discipline. Its ntpd is configured via ntpsec to serve the lab subnet and use the GPS PPS signal as the primary time reference.


### Configuration on Client Hosts


The NTP configuration is managed by Ansible via the ntp role in the lab-franklin collection. Here is what a typical client configuration looks like (as seen on stargate):

lstlisting[style=mystyle]
# /etc/ntpsec/ntp.conf -- managed by Ansible
driftfile /var/lib/ntp/ntp.drift
leapfile /usr/share/zoneinfo/leap-seconds.list

statsdir /var/log/ntpstats/

statistics loopstats peerstats clockstats
filegen loopstats file loopstats type day enable
filegen peerstats file peerstats type day enable
filegen clockstats file clockstats type day enable

server time.lab.bitsmasher.net
restrict 10.10.8.0/21
restrict -4 default
lstlisting

Key points:

- driftfile: Tracks oscillator drift between NTP restarts. Critical for stability on hosts that reboot frequently (Jetson devices, test VMs).
- statistics: Logs loop/peer/clock stats daily in /var/log/ntpstats/. Use these to verify sync health.
- server: Points to the lab's time host -- never pool.ntp.org for infrastructure hosts that need deterministic time sources.
- restrict:
- 10.10.8.0/21 grants query access to the main lab subnet (used by other internal clients or monitoring tools).
- -4 default blocks all IPv4 queries not explicitly allowed, providing a deny-by-default posture.

itemize


### Installation and Setup Steps


For a new host that needs NTP:


- Install ntpsec packages:
- Configure /etc/ntpsec/ntp.conf (either manually or via Ansible ntp role):
- Set server time.lab.bitsmasher.net as the upstream reference.
- Set timezone: timedatectl set-timezone America/Denver (or appropriate TZ).
- Create the log directory if it doesn't exist: mkdir -p /var/log/ntpstats.


     Start and enable the service:
    lstlisting[style=mystyle]
systemctl enable --now ntpsec.service
timedatectl  # verify "System clock synchronized: yes"
    lstlisting
enumerate


### Verification on Clients


Check synchronization status with these commands:

lstlisting[style=mystyle]
# Check if the system clock is synchronized
timedatectl | grep "System clock synchronized"

# Check NTP peer status (if ntpq is available)
ntpq -p

# View recent sync logs
journalctl -u ntpsec --since today
    lstlisting


### Notes for Deployment


- The Ansible ntp role handles all of the above automatically across lab hosts. It installs packages, copies configuration, sets timezone, and creates log directories.
- For Jetson devices (node900--node903), drift compensation is especially important -- they boot frequently and have less stable oscillators than desktop/server hardware.
- The GPS unit on the time host should be verified periodically: confirm it's locking to satellites, not using its internal oscillator as a holdover source. Check this by running ntptime on the time host itself.
- If the time host goes offline and all clients lose sync beyond their maximum correction threshold (typically 128 seconds), systems may need manual time correction before Kerberos/TLS will function again.


============================================================
## Troubleshooting
============================================================
## Troubleshooting Guide


### SSH Connection Failures


#### Host Key Mismatch


Symptom: @@@@@ WARNING: REMOTE HOST IDENTIFICATION HAS CHANGED! @@@@@

Fix:
lstlisting[style=mystyle]
ssh-keygen -R <hostname>
    lstlisting

This occurs when a host was re-imaged with the same IP. Check /etc/ssh/ssh\_config.d/ for stale entries.


#### Auth Denied (Port 22 Open, Key Rejected)


Symptom: TCP port 22 is reachable but SSH returns "Permission denied".

Common causes:

- No matching key in the remote host's authorized\_keys
- The key file permissions are too open on the target (must be mode 600 for .pub and 700 for ~/.ssh)
- Wrong user context (e.g., trying franklin@ when only root is configured)


Fix: Add the openclaw ed25519 key to the target host's authorized\_keys file. The pubkey is:
lstlisting[style=mystyle]
ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIJb9mV02PpxD8VpzYCnu7192dTwMnGWUc3qh55BIeElN openclaw@chonk
    lstlisting


#### No Route to Host


Symptom: SSH hangs or returns "No route to host". The host is likely powered off. Common in Jetson nodes (node902, etc.) that are manually managed.


#### GSSAPI Interference


Symptom: SSH times out waiting for GSSAPI authentication before falling back to key auth.

Fix: Add to ~/.ssh/config:
lstlisting[style=mystyle]
Host *
    GSSAPIAuthentication no
    IdentitiesOnly yes
    IdentityFile ~/.ssh/id_ed25519_openclaw
    lstlisting


### Service-Specific Diagnostics


#### Kerberos Authentication Failure


tabularx{}{l X}
**Symptom** & **Likely Cause** 

kinit fails with "clock skew" & NTP not synced -- check time.host first
kinit fails with "KDC has no support for encryption type" & Wrong krb5.conf realm or KDC unavailable
kadmin.local works but kadmin (remote) fails & Port 749 firewall rule on odroid-c1
Ticket expired after reboot & Time drift during host downtime
tabularx

Verify: klist, kinit username@LAB.BITSMASHER.NET, ktutil for keytabs.


#### NTP Issues


tabularx{}{l X}
**Symptom** & **Likely Cause** 

ntpd not running & ntpsec service not enabled
"no server suitable for synchronization" & No upstream reachable (time host offline)
Clock gradually drifting & Missing driftfile or Jetson oscillator drift
Time jumps by several seconds & Large correction applied; step-sync mode engaged
tabularx

Verify: ntpq -p, timedatectl, check /var/log/ntpstats/.


#### NFS Mount Issues


tabularx{}{l X}
**Symptom** & **Likely Cause** 

mount hangs & NFS server unreachable or export doesn't list client
"Stale file handle" & Server-side data was deleted/recreated
Permission denied on mount & /etc/exports restricts the client IP
tabularx

Verify: showmount -e <server>, check NFS server's exports, verify client IP matches the export rule.


### Common Ansible Troubleshooting


#### Module Not Found / Missing Collection


lstlisting[style=mystyle]
ansible-galaxy collection install lab.franklin
# or for individual modules:
ansible-galaxy collection install community.general
    lstlisting


#### Molecule Test Failures on New Hosts


If molecule tests fail on a new host, the most common causes are:

- Podman not installed or permissions wrong (openclaw user needs podman group)
- Docker driver configured instead of podman -- edit molecule.yml to set driver
- Python version mismatch (molecule requires $$3.9)
- Role hostname gates preventing test execution (e.g., the DNS role's hostname check for the real server vs. molecule container name)


### General Debugging Checklist


When a service is down, work through this checklist:


- Can you ping the host? ($$ network up?)
- Is SSH responding? ($$ host OS running?)
- Does the service process exist? (systemctl status or ps aux)
- Are logs showing errors? (journalctl -u <service> --since today)
- Can you connect to the service port from another lab host? (nc -zv <host> <port>)
- Is DNS resolving correctly? (dig or getent hosts)
- Is time synchronized? (timedatectl)


### Quick Reference Commands


lstlisting[style=mystyle]
# Check if a host is reachable from chonk
ping -c 1 <hostname>
ssh -o ConnectTimeout=5 openclaw@<host> "echo alive"

# Check service status remotely
ssh openclaw@<host> "systemctl status <service>"

# Verify time sync across hosts
for host in stargate skynet chonk; do
    ssh openclaw@$host "timedatectl | grep synchronized"
done

# Verify DNS resolution from a client
dig @time.lab.bitsmasher.net <hostname>

# Check Kerberos auth
kinit username@LAB.BITSMASHER.NET
klist
    lstlisting