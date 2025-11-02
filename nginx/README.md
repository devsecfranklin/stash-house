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

## TLS Cetrtificates

The `/etc/letsencrypt/live` directory contains your keys and certificates.

`[cert name]/privkey.pem`  : the private key for your certificate.
`[cert name]/fullchain.pem`: the certificate file used in most server software.
`[cert name]/chain.pem`    : used for OCSP stapling in Nginx >=1.3.7.
`[cert name]/cert.pem`     : will break many server configurations, and should not be used
                 without reading further documentation (see link below).

WARNING: DO NOT MOVE OR RENAME THESE FILES!
         Certbot expects these files to remain in this location in order
         to function properly!

We recommend not moving these files. For more information, see the Certbot
User Guide at https://certbot.eff.org/docs/using.html#where-are-my-certificates.
