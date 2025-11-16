# Kerberos

If you log in on a host that has a properly installed `/etc/krb5.keytab`
you will get a similar effect to requesting a new ticket, or typing `kinit -R`
to refresh.

This is an example of how the keytab should look on an example machine:

```sh
root@thelio:~# klist -ke /etc/krb5.keytab 
Keytab name: FILE:/etc/krb5.keytab
KVNO Principal
---- --------------------------------------------------------------------------
   2 host/thelio.lab.bitsmasher.net@LAB.BITSMASHER.NET (aes256-cts-hmac-sha1-96) 
   2 host/thelio.lab.bitsmasher.net@LAB.BITSMASHER.NET (aes128-cts-hmac-sha1-96) 
   2 ldap/thelio.lab.bitsmasher.net@LAB.BITSMASHER.NET (aes256-cts-hmac-sha1-96) 
   2 ldap/thelio.lab.bitsmasher.net@LAB.BITSMASHER.NET (aes128-cts-hmac-sha1-96) 
   2 nfs/thelio.lab.bitsmasher.net@LAB.BITSMASHER.NET (aes256-cts-hmac-sha1-96) 
   2 nfs/thelio.lab.bitsmasher.net@LAB.BITSMASHER.NET (aes128-cts-hmac-sha1-96)
```

## USER Principal

* Add the user to the KDC:

```sh
kinit -f root/admin
klist
kadmin -p  root/admin@LAB.BITSMASHER.NET
addprinc -randkey sly@LAB.BITSMASHER.NET
cpw sly # Reset Passwd for a User
klist -ke /etc/krb5.keytab
```

## HOST Principal

* Find the correct hostname for each principal:

```sh
getent hosts $(hostname) | awk '{print $1; exit}' | xargs getent hosts | awk '{print $2}'
```

* Add the NFS principal for the hosts to the KDC:

```sh
kinit -f root/admin
klist
kadmin -p  root/admin@LAB.BITSMASHER.NET
addprinc -randkey nfs/snowy.lab.bitsmasher.net@LAB.BITSMASHER.NET
ktadd nfs/snowy.lab.bitsmasher.net@LAB.BITSMASHER.NET host/snowy.lab.bitsmasher.net@LAB.BITSMASHER.NET ldap/snowy.lab.bitsmasher.net@LAB.BITSMASHER.NET
```

* Validate

```sh
klist -ke /etc/krb5.keytab # linnux
/usr/local/heimdal/bin/klist #openbsd
```

## KDC Files

```sh
/etc/krb5.keytab
/usr/share/krb5-kdc
/etc/krb5kdc/kdc.conf
/var/lib/krb5kdc
```

Check the files with the `krb_client.sh` script.

## KDC - Create a kadmind Keytab

* [Create a kadmind Keytab](https://web.mit.edu/kerberos/krb5-1.5/krb5-1.5.4/doc/krb5-install/Create-a-kadmind-Keytab-_0028optional_0029.html)

The kadmind keytab is the key that the legacy admininstration daemons kadmind4 and v5passwdd
will use to decrypt administrators' or clients' Kerberos tickets to determine whether or not
they should have access to the database. You need to create the kadmin keytab with entries for
the principals `kadmin/admin` and `kadmin/changepw`.

## SSSD integration with Active Directory

The System Security Services Daemon is software originally developed for the Linux
operating system that provides a set of daemons to manage access to remote\
directory services and authentication mechanisms.

## SSH

```sh
for key in ~/.ssh/id_*; do ssh-keygen -l -f "${key}"; done | uniq
ssh-keygen -t ed25519 -C "fdiaz@paloaltonetworks.com" -f ~/.ssh/id_ed25519_work -o -a 100
```

### Ticket Forwarding

* set `forwardable = True` in `/etc/krb5.conf`
* set `GSSAPIDelegateCredentials=yes` and `GSSAPIAuthentication=yes` in `~/.ssh/config`

### Mac

If you’re using macOS Sierra 10.12.2 or later, to load the keys automatically and store
the passphrases in the Keychain, you need to configure your ~/.ssh/config file:

```sh
Host *
  AddKeysToAgent yes
  UseKeychain yes
  IdentityFile ~/.ssh/id_ed25519
  IdentityFile ~/.ssh/id_rsa # Keep any old key files if you want
```

Once the SSH config file is updated, add the private-key to the SSH agent:

```sh
ssh-add -K ~/.ssh/id_ed25519
```

### OpenBSD SSH setup

TBD

## User Setup

```sh
gpg --list-keys
```
