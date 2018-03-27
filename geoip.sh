#!/bin/bash
# Update GEOIP-DB for lookups

mkdir -p "/tmp/geoip"
wget -q "http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.mmdb.gz" -O /tmp/geoip/city.mmdb.gz
wget -q "http://geolite.maxmind.com/download/geoip/database/GeoLite2-ASN.tar.gz" -O /tmp/geoip/asn.mmdb.gz
gzip -d /tmp/geoip/city.mmdb.gz -f
gzip -d /tmp/geoip/asn.mmdb.gz -f
mv /tmp/geoip/* /home/dnsleak