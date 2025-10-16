# nginx

- copy the file to `/etc/nginx/sites-available`
- make a link to `sites-enabled`

## service

Example:

```sh
cp nginx-lab.conf /etc/nginx/sites-available
cp lab.bitsmasher.net.website.service /etc/systemd/system/lab.bitsmasher.net.website.service
systemctl enable lab.bitsmasher.net.website
systemctl start lab.bitsmasher.net.website.service
```
