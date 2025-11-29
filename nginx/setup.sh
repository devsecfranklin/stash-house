sudo apt install -y nginx

sudo /usr/sbin/nginx -t


function test_nostr() {
  curl -k -L http://bitsmasher.net/.well-known/nostr.json
  ls -al /usr/share/nginx/html/.well-known/
}