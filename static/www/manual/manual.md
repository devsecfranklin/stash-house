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


Based on the most recent audit of ansible/collections/ansible\_collections/lab/franklin/roles/:

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


### Molecule Testing -- Status and Deprecation


Molecule (the canonical integration testing framework for Ansible roles and playbooks) remains **active** as of August 2026 but its role has shifted:


- Molecule v26.4.x is the latest series, released June 2026. It still supports only the latest two major versions of Ansible (N/N-1).
- Molecule is being superseded as the preferred test strategy by:
- pytest-ansible -- pytest extension for testing Ansible module/plugin Python code directly.
- tox-ansible -- integration with tox to run tests across multiple Python interpreters and ansible-core versions simultaneously.
- Execution Environments (containers) as the primary test target, tested via ansible-navigator.

    
     The **community.molecule** collection still exists but new projects should prefer pytest-ansible + tox-ansible for role/test harness work. Molecule is not yet officially deprecated -- it is maintained on a best-effort basis while the community transitions.
    
     For existing molecule test harnesses (such as those in lab-franklin's Ansible collection at ansible/collections/ansible\_collections/lab/franklin/roles/), they will continue to work with ansible-core 2.20 and pytest-6.x but should be audited for:
    
- Podman driver compatibility (molecule continues using podman or docker as default driver)
- Python version requirements (molecule requires Python $$3.9)
- Plugin deprecations in newer collections


itemize

**Note**: For lab-franklin's molecule tests, keep testing with podman for now. The container-based approach is still valid; only the *preferred tooling around it* is shifting toward pytest-ansible + execution environments.

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


The collection contains roles organized under the roles/ directory. As of August 2026, the following roles are maintained with molecule test harnesses:

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

Roles currently **missing molecule tests**: edge, mail, odroid, website (incomplete roles). These have tasks/main.yml but no test harness.


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
- Makefile.am: Top-level targets include test (ansible-lint + molecule), security (venv bootstrap), dev (Python venv in BUILDDIR)
- bootstrap.sh: Cross-platform bootstrapper that detects OS (Debian/RedHat/OpenBSD/macOS/Linux) and installs platform-specific packages, then runs aclocal $$ autoreconf $$ automake $$ configure
- docs/manual/Makefile.am: Standalone autotools stub for LaTeX docs build (clean target removes *.aux files)


The network\_update.sh script was hardened on 2026-08-02 with: set -euo pipefail, auto-resolved paths, pinned K3s version validation, ANSIBLE_HOME environment checks, removed dangerous clush/apt-get blast radius, and dead code removal. It runs the main playbook at ansible/playbook.yml as its primary action.

Testing from Stargate


## Molecule Test Harnesses on stargate.lab.bitsmasher.net


The stargate host (10.10.16.66, Debian 12 bookworm) serves as the primary testing platform for the lab/franklin Ansible collection. Testing is run from /mnt/clusterfs2/workspace/lab-franklin/ansible/collections/ansible\_collections/lab/franklin/roles/ role /molecule/.


### Test Environment


Molecule tests on stargate use:

- Driver: podman (containers, not Docker -- avoids daemon dependency)
- Base image: Ubuntu 20.04 for most roles, with role-specific variations
- Execution user: openclaw user (has sudo -n, no password required)
- Python: $$ 3.9 (required by both ansible and molecule)


### Test Execution Flow


For a given role, the molecule test cycle is:

- dependency: Install collection dependencies from galaxy.yml
- syntax: Validate playbook YAML syntax
- create: Launch podman container (e.g., "server1") as the test target host
- prepare: Verify role prerequisites within the container
- converge: Run ansible-playbook against the container -- this is where the role's tasks are executed
- test/verify: Optional pytest sequence runs inside the container to validate convergence


### Example: DNS Role Testing


The DNS role molecule test confirmed:

- Podman creates an Ubuntu 20.04 container (hostname "server1")
- Bind9 packages install correctly
- /etc/bind directory is created with named.conf.options written
- Molecule platform name was changed from "molecule-ubuntu" to "server1" to match the role's hostname gate
- Group assignment updated to "dns\_servers" for proper role targeting


### Known Limitations


Some roles cannot be fully tested on molecule containers:


- nfs: Hostname gates ('snowy' in inventory\_hostname) and hardcoded mount paths depend on real NFS servers
- openbsd: The role targets OpenBSD hosts specifically; molecule tests run on Ubuntu containers by default
- jetson-nano: Hardware-specific (NVIDIA Jetson); molecule tests only verify file placement, not actual provisioning


### Running Tests Manually


To test a specific role from the stargate host:
lstlisting[style=mystyle]
cd /mnt/clusterfs2/workspace/lab-franklin/ansible/collections/ansible_collections/lab/franklin/roles/<role>/molecule
molecule converge --scenario-name default
    lstlisting

To verify convergence results:
lstlisting[style=mystyle]
molecule test --scenario-name default
    lstlisting

If the podman driver is unavailable, fall back to docker:
lstlisting[style=mystyle]
molecule converge --driver-name docker
    lstlisting


### Test Results Dashboard


Testing output logs are captured in the molecule scenario directory under logs/. The full lab-franklin test suite (all roles with molecule tests) is run via the top-level Makefile.am test target:
lstlisting[style=mystyle]
make test
# Runs ansible-lint --all-roles and molecule converge per role
    lstlisting


### pytest-ansible Transition Plan


As molecule is being superseded, new tests should use pytest-ansible where possible:

- Migrate the test directory structure from molecule/ to a tests/ directory with pytest configuration
- Use ansible-test for collection sanity tests (built into ansible-core)
- Keep molecule only where full environment isolation is needed (e.g., DNS zone file validation, NTP time sync testing)


### Pending Tasks and Todos


- TODO: Re-establish internal PyPI server on stargate. The pypi\_internal role is ready (role exists, molecule test added Aug 2026) but the actual PyPI service needs to be provisioned on the storage host. Requires: NFS mount of /mnt/storage1 available on stargate, web server (nginx/apache) installed and configured, pip packages populated into /mnt/storage1/LAB/pypi with PEP 503 structure. Run: ansible-playbook ... -l storage --tags pypi\_internal once the infrastructure is in place.


============================================================
## Cluster
============================================================
## Clustershell


all
head
compute
gpu


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


### BIND Overview


The lab uses BIND (named) as its authoritative DNS server. The DNS infrastructure is managed by the dns Ansible role in the lab-franklin collection, located at:
lstlisting[style=mystyle]
ansible/collections/ansible_collections/lab/franklin/roles/dns/
    lstlisting


### Zones Maintained


- Forward zone: db.home.lab -- maps hostnames to IPs for all lab infrastructure
- Reverse zones: PTR records for each subnet (10.10.x.0/24 ranges)
- SRV records: Kerberos service discovery (\_kerberos.\_tcp, \_kdc.tcp)


### DNS Server Location


The DNS server (ns1) is the single authoritative nameserver for the lab. It has been unreachable during recent infrastructure audits, which explains why many hosts cannot resolve hostnames -- ns1 *is* the DNS server and it's offline.

When ns1 is available, it resolves all internal hostnames. Key entries include:

- time.lab.bitsmasher.net $$ 10.10.12.2 (NTP server)
- odroid-c1.lab.bitsmasher.net $$ 10.10.12.254 (KDC)
- ldap.lab.bitsmasher.net $$ 10.10.13.1 (OpenLDAP)
- blowfish.lab.bitsmasher.net $$ 10.10.14.85 (reassigned from 10.10.12.15)


### Molecule Test Harness


The DNS role has a molecule test harness that confirms:

- Converge: BIND packages installed, /etc/bind directory created, named.conf.options written
- Test sequence: dependency $$ syntax $$ create (podman container) $$ prepare $$ converge
- Limitation: molecule test uses Ubuntu 20.04 container with hostname "server1"; the role has a hostname gate that targets the real DNS server's hostname


The DNS molecule test confirms correct package installation and configuration file placement but does not verify actual name resolution (no live DNS server in the container test environment).


### DNS Failure Impact


When ns1 is offline:

- Lab hosts fall back to /etc/hosts entries (which are authoritative but not dynamic)
- Kerberos KDC discovery via SRV records fails -- clients need manual kdc = entry in krb5.conf
- DHCP reservations on the Dream Machine remain functional (DHCP server is separate from DNS)


The DNS role's molecule test harness ensures that when the server comes back online, its configuration will install and write correctly.


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

chonk & 10.10.8.60 & franklin/openclaw & Gateway host, primary workstation
stargate & 10.10.16.66 & openclaw & Ansible collection workspace, build host
skynet & 10.10.16.10 & openclaw/franklin & Minecraft server, IP-direct login (hostname resolves oddly)
time & 10.10.12.2 & openclaw/franklin & Stratum-1 NTP server with GPS reference
wonderland & 178.62.60.55 & franklin/openclaw & Public cloud VM (Debian 12, 6.1 kernel)
music.lan & 192.168.86.38 & root & Desktop machine via skynet route
tabularx

*{Jetson Cluster ( 87 days uptime)}

All three Jetsons run Ubuntu 18.04 with ed25519 key auth (openclaw permanent, franklin has password 123 as fallback):

tabularx{}{l l X}
**Host** & **IP** & **Status** 

node900 & 10.10.12.90 & Online, key auth
node901 & 10.10.12.91 & Online, key auth
node903 & 10.10.12.93 & Online, key auth
node902 & 10.10.12.92 & Off / no route to host (likely powered down)
tabularx

*{KDC and Directory Services}

tabularx{}{l l X}
**Host** & **IP** & **Status** 

odroid-c1 / kdc1 & 10.10.12.254 & Online (~242d uptime), root key auth, KDC for lab.bitsmasher.net Kerberos realm
ldap/bbb1 & 10.10.13.1 & SSH reachable as root; slapd failed Dec 2025 (TLS cert issue)
tabularx


### Historical / Decommissioned Hosts


- snowy: Hard disk physically removed. /mnt/snowy is now just a mounted drive (not SSH-accessible). Should be cleaned from DNS and SSH config.
- blowfish: Previously at 10.10.12.15, reassigned to 10.10.14.85. Key not authorized yet -- unreachable from chonk but port 22 is open. Expected role: database host with OpenBSD.
- node902: Powered off; no route to host from any known peer.


### Planned Hardware Expansion


- blowfish as the primary database server (OpenBSD planned)
- head2 for cluster management (referenced in network\_update.sh but unreachable)
- k3s_server and k3s_agent roles for Kubernetes across node infrastructure


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

Key file: ~/.ssh/id\_ed25519\_openclaw (the standard openclaw key).


### Principal Conventions


Kerberos principals in the lab.bitsmasher.net realm should follow these conventions:

- Service principals: service/FQDN@lab.BITSMASHER.NET
- User principals: username@lab.BITSMASHER.NET
- Host principals: host/FQDN@lab.BITSMASHER.NET


### DNS/Kerberos Integration


Kerberos relies on DNS SRV records for KDC discovery:
lstlisting[style=mystyle]
_kerberos._tcp.lab.bitsmasher.net  IN SRV 0 0 88 odroid-c1.lab.bitsmasher.net.
_kerberos._udp.lab.bitsmasher.net  IN SRV 0 0 88 odroid-c1.lab.bitsmasher.net.
_kdc.tcp.lab.bitsmasher.net        IN SRV 0 0 88 odroid-c1.lab.bitsmasher.net.
    lstlisting

These records should be present in the BIND DNS zones on the lab's DNS server. If ns1 is offline (as has been the case), Kerberos clients may fail to discover the KDC automatically and will need manual configuration via kdc = odroid-c1.lab.bitsmasher.net in /etc/krb5.conf.


### Keytab Management


Keytabs for service principals should be:

- Generated on the KDC host (odroid-c1) using kadmin.local or kadmin -q
- Transferred securely to target hosts (via scp with SSH keys)
- Set to mode 600 and owned by the appropriate service user


The lab-franklin Ansible collection should have a role for managing keytabs across hosts, though this may need to be verified against current KDC state.


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


The NFS role at ansible/collections/ansible\_collections/lab/franklin/roles/nfs has a molecule test harness, but testing is complex due to hostname-dependent task gates. This chapter documents how to run tests and what to expect.


## Prerequisites


- podman available on the test host (chonk tested successfully)
- The lab.franklin collection with all dependencies installed
- A molecule environment with an Ubuntu-based container (the role expects Debian/Ubuntu apt packages)


## Running the Test


From the role directory:

verbatim
cd /mnt/clusterfs2/workspace/lab-franklin/ansible/collections/ansible_collections/lab/franklin/roles/nfs
molecule converge --scenario-name default
molecule verify --scenario-name default
verbatim


## What the Test Validates


- Package installation: nfs-kernel-server, nfs-common, rpcbind, nfs-client are installed
- Configuration file deployment: /etc/default/nfs-kernel-server, /etc/default/nfs-common, /etc/idmapd.conf are written with correct content
- Directory creation: Export mount points (/mnt/clusterfs, /mnt/storage\{1,2,3\}) are created
- Export template rendering: Jinja2 templates generate valid /etc/exports content


## Current Test Limitations


The role has two categories of gates that block full testing:

*{1. Server-side hostname gates (in tasks/main.yml)}

verbatim
when: 'snowy' in inventory_hostname
when: 'thelio' in inventory_hostname
verbatim

Molecule containers cannot satisfy these conditions, so nfs\_server\_*.yml task files never execute. The server-side converge phase is effectively skipped.

*{2. molecule environment variable guards}

Most mount and service-start tasks check:
verbatim
when: HOMELAB_MOLECULE_TEST is not defined
verbatim

This prevents actual mounts from happening during testing (which would require real NFS servers), but it also means the nfs\_client.yml tasks that matter most for client-side validation never run their key stages.


## How to Fix the Test Harness


To make molecule tests meaningful, the following changes are needed:


- Add test-specific inventory: In molecule/default/inventory.yml, define a hostgroup (e.g., nfs\_servers) that includes container hostnames matching the role's hostname gates. Example:
- Add override variables in molecule\_prepare.yml: Set nfs\_exports\_snowy and nfs\_exports\_thelio to container-safe values that use local paths, so export template rendering can be validated without real storage.
- Create a test scenario for client-only testing: A second molecule scenario with only nfs\_clients inventory that tests mount directory creation, fstab line insertion, and idmapd.conf deployment. Use state: present instead of state: mounted in the molecule overrides so it doesn't try to connect to a live NFS server.
- Replace hostname gates with role variables: Replace 'snowy' in inventory\_hostname with a variable like nfs\_server\_mode | default(false) so tests can enable server tasks without depending on real hostnames.
- Add idempotence tests: A second converge pass to verify that running the role again produces no changes (the force: no on idmapd.conf copy already helps here).


## Integration Test Targets


An integration test target exists at tests/integration/targets/nfs-server/ but it is minimal. For thorough testing, add:


- Export verification (check showmount -e)
- Mount/unmount cycle tests on client container
- Kerberos security validation (try mounting with and without sec=krb5i)
- Permission checks after mount (verify root\_squash behavior)


## Molecule Test Checklist for Future Work


Before declaring the NFS role test-ready:


- [ ] Server-side converge passes (with hostname override in molecule inventory)
- [ ] Client-side converge passes (mount directory creation, fstab entries)
- [ ] Verify phase checks exports file content and configuration files
- [ ] Idempotence test: second converge shows zero changed tasks
- [ ] Integration test validates live mount on container


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


### Password Policy

Passwords are used only as a temporary fallback. The Jetson nodes (node900--node903) currently use password ``123'' for the franklin user during provisioning; this should be replaced with key-based authentication when possible.


### User Accounts


center
tabular{lll}
  
  **User** & **Purpose** & **Access Scope** 
 
  franklin & Primary admin user & All hosts; sudo NOPASSWD on skynet 
  openclaw & Automation user & All production hosts via ed25519 key 
  root & Emergency access & KDC, ldap, music host only 
 
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
  DNS (BIND) & ns1.lab.bitsmasher.net & Offline & December 2025 
  OpenLDAP & ldap.lab.bitsmasher.net & Offline & December 2025 
  Minecraft & skynet.lab.bitsmasher.net & Online & Current 
 
tabular
center


### Disaster Recovery Priorities


- NTP (time host) -- time sync failure cascades to all Kerberos authentication
- DNS (BIND) -- resolution failure makes most hosts unreachable by name
- Kerberos KDC -- authentication cascade failure affects all service access
- LDAP -- directory services, offline since December 2025


## Container and Kubernetes Security


The k3s cluster is hardened with:

- Version pinning to 2.1.x
- Containerd runtime with NVIDIA GPU support
- TLS for all API communication (k3s default)
- Pod security policies via Kubernetes native admission controllers


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


### NFS Mounts


The lab relies on NFS for shared storage across hosts. Key mount points:


- /mnt/clusterfs2: Primary cluster workspace mounted on chonk and stargate (from the NFS server). This is where Ansible roles, molecule tests, and project files live.
- /mnt/snowy: Formerly a full NFS share; now just a mounted drive after the hard disk was physically removed. Should be audited for stale mount entries.


The NFS role in lab-franklin's Ansible collection handles:

- Server-side export configuration (/etc/exports)
- Client-side fstab management and mount points
- Mount path conventions per role (previously /mnt/clusterfs, /mnt/backup1, /mnt/storage{1,2,3})


Note: The NFS role has hostname gates ('snowy' in inventory\_hostname) that prevent testing on molecule containers. Hardcoded mount paths depend on real NFS servers being online. This is a known limitation of the test harness.


### Snapshot and Backup Strategy


The lab uses bare git repos for version control and backup:

- GPG-encrypted repositories synced across hosts for credential protection (Stash House project)
- Pass-based password store with GPG backend
- No LFS -- ABSOLUTE RULE: do not use git-lfs anywhere. Microsoft once held Minecraft chunk files hostage when a repo hit the 10GB GitHub limit with LFS enabled.


### Local Storage Considerations


- Jetson nodes (node90x): limited eMMC storage; rely on NFS for workspace data
- chonk: large local disk, primary build and gateway host
- stargate: Debian 12 box with full Ansible workspace mirror


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