#!/bin/bash
# Update GEOIP-DB for lookups
set -euo pipefail
IFS=$'\n\t'
mkdir -p "/tmp/geoip"
wget -q "http://geolite.maxmind.com/download/geoip/database/GeoLite2-Country.tar.gz" -O /tmp/geoip/country.mmdb.tgz
wget -q "http://geolite.maxmind.com/download/geoip/database/GeoLite2-ASN.tar.gz" -O /tmp/geoip/asn.mmdb.tgz
tar -zxf /tmp/geoip/country.mmdb.tgz -C /tmp/geoip --strip-components=1
tar -zxf /tmp/geoip/asn.mmdb.tgz -C /tmp/geoip --strip-components=1

mv /tmp/geoip/GeoLite2-Country.mmdb /home/dnsleak/country.mmdb
mv /tmp/geoip/GeoLite2-ASN.mmdb /home/dnsleak/asn.mmdb
