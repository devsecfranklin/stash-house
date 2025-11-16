# NTP

* [How to use ntp.org](https://www.pool.ntp.org/en/use.html)
* [Setting up an NTP server on a Raspberry Pi Zero W 2](https://gist.github.com/p0lr/ce2203cf6e23c2c95b556ea20f0144f7)

## configure

* Set the time zone: `sudo timedatectl set-timezone America/Denver`

```sh
sudo apt install ntp gpsd gpsd-clients # gpsd-clients has `xgps` used for testing
chmod a+rw /dev/ttyUSB0
sudo systemctl enable gpsd # enabled gpsd on boot
cgps -s # test gpsd
# stty -F /dev/ttyUSB0 speed 4800 && cat </dev/ttyUSB0 # test the gps
mkdir /var/log/ntpsec && chmod 755 /var/log/ntpsec
systemctl enable ntp
ntpmon 127.0.0.1
```

## Update System Time Sync

Modify the NTP line in `sudo nano /etc/systemd/timesyncd.conf` to the IP of the time-server
`time.lab.bitsmasher.net`

```sh
[Time]
NTP=127.0.0.1
```

## testing

* Be sure to use BASH not FISH or other shells while molecule testing

```sh
molecule lint
molecule converge
molecule test --scenario-name default # the create, converge and destroy steps will be run one after another.
```

## OpenBSD

* Configure it with `/etc/ntpd.conf`
* there is a `.j2` template config file for this

```sh
$ ntpctl -sa # check the time on openbsd

1/1 peers valid, constraint offset 0s, clock synced, stratum 2

peer
   wt tl st  next  poll          offset       delay      jitter
10.10.12.13 time.lab.bitsmasher.net
 *  1 10  1    7s   30s         4.084ms     0.514ms     0.064ms\
 ```
