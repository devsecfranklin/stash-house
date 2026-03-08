# 03 Safe Use

* Credentials

Store credentials in "pass" manager.

```sh
doas apk update && doas apk add pass
gpg --list-keys # get your public key id
pass init C25565E4701F4ED36A0711AA114F3606EFD923BB # id of your public GPG key
pass insert aws-access-key-id
pass ls
pass insert aws-secret-access-key
pass ls
pass show aws-secret-access-key
```

* Create Terraform plan.

```bash
pass show aws-secret-access-key
export AWS_ACCESS_KEY_ID=$(pass aws-access-key-id)
export AWS_SECRET_ACCESS_KEY=$(pass aws-secret-access-key)
terraform plan
terraform plan -out franklin.plan
terraform show -json franklin.plan > tfplan.json
```