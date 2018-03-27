![SpyOFF](https://github.com/mpdroog/dnsleak/blob/master/web/spyoff.svg)

DNSleak
==================
Small DNS-server that catches requests and offers origins through HTTP(s) API.

What is used?
- LetsEncrypt to offer easy HTTPS-requests
- CORS-headers are added to cross-domain do AJAX-requests
- ISP/Country lookup through Maxmind's GeoIP

How to use?
- Point A-record to this node i.e. ns-dnstest.spyoff.com
- Point NS-record to this node i.e. dnstest.spyoff.com
- POST https://ns-dnstest.spyoff.com/dns/leaktest
 IN: ```{domain: ["4eb4b123bbd72478a29bff21cd00f48722b704ce.dnstest.spyoff.com"]}```
 OUT: ```{"15169":{"ISP":"Google LLC","Country":"US","IP":"172.217.40.8"}```

Arguments
```
./dnsleak --help
Usage of ./dnsleak:
  -d string
    	DNS listen on (both tcp and udp) (default "[::]:53")
  -h string
    	HTTP listen on (default "[::]:80")
  -m string
    	HTTPS-domain (LetsEncrypt) (default "ns-dnstest.spyoff.com")
  -s string
    	HTTPS listen on (default "[::]:443")
  -v	Verbose-mode (log more)
```

Tool created for [SpyOFF](https://www.spyoff.com/dns-leak-test/?a_aid=11108&a_bid=02dc3d81)

Install
```
# User + systemd
useradd -r dnsleak
mkdir -p /home/dnsleak
vi /etc/systemd/system/dnsleak.service
# Systemd file below...

chmod 644 /etc/systemd/system/dnsleak.service
systemctl daemon-reload
systemctl enable dnsleak
systemctl start dnsleak

# MaxMind GeoIP
vi /etc/cron.d/dnsleak-geo
# @daily dnsleak /home/dnsleak/geoip.sh
mkdir -p /tmp/geoip
chown dnsleak:dnsleak -R /tmp/geoip
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