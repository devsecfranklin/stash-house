# README.md

## test web site locally

```bash
sudo docker pull nginx
docker build -t nginx-bitsmasher .
docker run --name docker-nginx-bitsmasher -p 8080:80 -d nginx-bitsmasher
```

Now navigate to http://http://0.0.0.0:8080/

## test twitter card

[twitter card validator](https://cards-dev.twitter.com/validator)
