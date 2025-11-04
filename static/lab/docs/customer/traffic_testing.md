# Traffic Testing

```sh
tcpdump filter "host 10.0.120.1 and not port 22"
view-pcap mgmt-pcap mgmt.pcap
scp export mgmt-pcap from mgmt.pcap to user@scpserver:/tmp
tcpdump snaplen 0
```
