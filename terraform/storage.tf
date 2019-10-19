resource "digitalocean_volume" "webserver_disk" {
  region      = "ams3"
  name        = "webserver_storage"
  size        = 10 
  description = "Extra space for your files"
}