# terraform

```sh
gpg --list-keys # get your public key id
pass init C25565E4701F4ED36A0711AA114F3606EFD923BB # id of your public GPG key
pass insert DO_TOKEN
pass ls
pass show
```

* Create Terraform plan.

```bash
export DO_TOKEN=$(pass DO_TOKEN) || export DO_TOKEN=(pass DO_TOKEN)
#terraform plan -out franklin.plan -var="do_token=${DO_TOKEN}" # BASH
terraform plan -out franklin.plan -var="do_token=$DO_TOKEN" # FISH
terraform show -json franklin.plan > tfplan.json
```

