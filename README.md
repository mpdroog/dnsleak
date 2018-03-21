![SpyOFF](https://github.com/mpdroog/dnsleak/blob/master/web/spyoff.svg)

DNSleak
==================
Catch DNS-requests and keep in memory to offer
the origin IP and notice if your DNS-requests pass through
a faulty server.

Tool created for [SpyOFF](https://www.spyoff.com/dns-leak-test/?a_aid=11108&a_bid=02dc3d81)

Install
```
useradd -r dnsleak
mkdir -p /home/dnsleak
vi /etc/systemd/system/dnsleak.service
chmod 644 /etc/systemd/system/dnsleak.service
systemctl daemon-reload
systemctl enable dnsleak
systemctl start dnsleak
```

/etc/systemd/system/dnsleak.service
```
[Unit]
Description=DNS Leak tester by faking a DNS-server
After=network.target
Requires=network.target

[Service]
LimitNOFILE=8192
Type=notify

Restart=always
RestartSec=30
TimeoutStartSec=0

WorkingDirectory=/home/dnsleak
ExecStart=/home/dnsleak/dnsleak
User=dnsleak
Group=dnsleak

CapabilityBoundingSet=CAP_NET_BIND_SERVICE
AmbientCapabilities=CAP_NET_BIND_SERVICE
NoNewPrivileges=true

[Install]
WantedBy=multi-user.target
```
