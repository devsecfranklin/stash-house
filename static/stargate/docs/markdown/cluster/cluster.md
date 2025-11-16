# cluster

Raspberry Pi cluster setup

## Requirements

`requirements.yml`

## Role Variables

A description of the settable variables for this role should go here, including any variables that are in defaults/main.yml, vars/main.yml, and any variables that can/should be set via parameters to the role. Any variables that are read from other roles and/or the global scope (ie. hostvars, group vars, etc.) should be mentioned here as well.

## Dependencies

A list of other roles hosted on Galaxy should go here, plus any details in regards to parameters that may need to be set for other roles, or variables that are used from other roles.

## Example Playbook

```sh
ansible-playbook /home/franklin/workspace/LAB/lab-franklin/ansible/collections/ansible_collections/lab/franklin/playbooks/cluster_raspi.yml -i /etc/ansible/hosts -b -e 'ansible_python_interpreter=/usr/bin/python3'
```

## License

BSD

## Author Information

An optional section for the role authors to include contact information, or a website (HTML is not allowed).

## SLURM

Reconfigure all Slurm daemons on all nodes. This should be done after changing the Slurm configuration file.

```sh
scontrol reconfig
scontrol show nodes
srun --mpi=list
```

```sh
sinfo –s # summary of cluster resources (-s --summarize)
sinfo -p <partition-name> -o %n,%C,%m,%z # compute info of nodes in a partition (-o --format)
sinfo -p Gpu -o %n,%C,%m,%G # GPUs information in Gpu partition (-p --partition)
sjstat –c # show computing resources per node
scontrol show partition <partition-name> # partition information
scontrol show node <node-name> # node information
sacctmgr show qos format=name,maxwall,maxsubmit # show quality of services
```

## CLUSH

<https://clustershell.readthedocs.io/en/latest/config.html>

### Checking status of the cluster

* clustat
* clustat -m  -> Display status of and exit
* clustat -s ->  Display status of and exit
* clustat -l -> Use long format for services
* cman_tool status -> Show local record of cluster status
* cman_tool nodes -> Show local record of cluster nodes
* cman_tool nodes -af
* ccs_tool lsnode -> List nodes
* ccs_tool lsfence ->  List fence devices
* group_tool ->  displays the status of fence, dlm and gfs groups
* group_tool ls ->  displays the list of groups and their membership

### Resource Group Control Commands

* clusvcadm -d -> Disable
* clusvcadm -e -> Enable
* clusvcadm -e -F -> Enable according to failover domain rules
* clusvcadm -e -m -> Enable on
* clusvcadm -r -m -> Relocate to member
* clusvcadm -R ->  Restart a group in place.
* clusvcadm -s -> Stop

### Resource Group Locking (for cluster Shutdown / Debugging)

* clusvcadm -l -> Lock local resource group manager. This prevents resource groups from starting on the local node.
* clusvcadm -S -> Show lock state
* clusvcadm -Z -> Freeze group in place
* clusvcadm -U -> Unfreeze/thaw group
* clusvcadm -u -> Unlock local resource group manager. This allows resource groups to start on the local node.
* clusvcadm -c -> Convalesce (repair, fix) resource group. Attempts to start failed, non-critical resources within a resource group.

### Cluster Command Quick Reference

* clustat -> Display the status of the cluster as viewed from the executing host
* clusvcadm -> Manage services across the cluster.
* clusvcadm -r -m -> Move a service to another node
* clusvcadm -d -> Stop a service
* clusvcadm -e -> Start a service

### Cluster configuration system

* ccs_tool ->  Online management of cluster configuration
* ccs_tool update /etc/cluster/cluster.conf ->  Update the cluster.conf file across the cluster
* cman_tool -> Manage cluster nodes and display the current state of the cluster
* cman_tool status
* fence_node  -> Eject a node from the cluster
* fence_tool dump -> Print fence debug messages

## Node Health Check

## R and Snow

```sh
sudo apt install dirmngr apt-transport-https ca-certificates software-properties-common -y
gpg --keyserver keyserver.ubuntu.com --recv-key '95C0FAF38DB3CCAD0C080A7BDC78B2DDEABC47B7'
gpg --armor --export '95C0FAF38DB3CCAD0C080A7BDC78B2DDEABC47B7' | gpg --dearmor | sudo tee /usr/share/keyrings/cran.gpg > /dev/null
echo "deb [signed-by=/usr/share/keyrings/cran.gpg] https://cloud.r-project.org/bin/linux/debian bookworm-cran40/" | sudo tee /etc/apt/sources.list.d/cran.list
sudo apt update
sudo apt install r-base r-base-dev r-recommended
```
