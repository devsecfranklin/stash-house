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

```sh
sudo certbot certonly --nginx -d bitsmasher.net -d www.bitsmasher.net
```

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

### cert for internal

```sh
certbot -d time.lab.bitsmasher.net --manual  --preferred-challenges dns certonly
```

## CORS testing

- [cors test site](https://cors-test.codehappy.dev/?url=https%3A%2F%2Fbitsmasher.net%2F.well-known%2Fnostr.json%3Fname%3Dthedevilsvoice&origin=https%3A%2F%2Fcors-test.codehappy.dev%2F&method=get)

```sh
curl -H "Origin: https://example.com" -I https://bitsmasher.net/.well-known/nostr.json?name=thedevilsvoice
```
