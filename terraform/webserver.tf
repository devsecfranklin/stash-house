# This is so we can test puppet agent.
resource "digitalocean_droplet" "webserver" {
  ssh_keys           = [13805615]
  image              = "ubuntu-16-04-x64"
  #image              = "debian-9-x64"
  name               = "www"
  region             = "ams3"
  size               = "1gb"
  private_networking = true
  backups            = false
  ipv6               = false
  volume_ids = ["${digitalocean_volume.webserver_disk.id}"]

  connection {
    user = "root"
    type = "ssh"
    host = self.ipv4_address
    # eval `ssh-agent -s` ; ssh-add ~/.ssh/do_terraform_rsa
    # agent = false
    private_key = "${file(var.pvt_key)}"
    timeout = "2m"
  }

  provisioner "remote-exec" {
    inline = [
      "export PATH=$PATH:/usr/bin",
      "apt-get update",
      "apt-get install puppet-agent -y",
      "echo \"${digitalocean_droplet.webserver.ipv4_address} www www.bitsmasher.net\" >> /etc/hosts",
    ]
  }
}
