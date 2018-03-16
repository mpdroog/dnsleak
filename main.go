package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	ttl_map "github.com/leprosus/golang-ttl-map"
	"github.com/miekg/dns"
	geoip2 "github.com/oschwald/geoip2-golang"
	"log"
	"net"
	"net/http"
	"strings"
)

var (
	Verbose bool
	cache   ttl_map.Heap
	dbASN   *geoip2.Reader
	dbCity  *geoip2.Reader
)

type Handle struct {
}

func (h *Handle) ServeDNS(w dns.ResponseWriter, req *dns.Msg) {
	ip := ""
	if addr, ok := w.RemoteAddr().(*net.UDPAddr); ok {
		ip = addr.IP.String()
	}
	if addr, ok := w.RemoteAddr().(*net.TCPAddr); ok {
		ip = addr.IP.String()
	}
	if ip == "" {
		panic("IP not found?")
	}

	domain := req.Question[0].Name
	domain = domain[:len(domain)-1]
	if Verbose {
		fmt.Printf("Origin=%s\n", ip)
		fmt.Printf("Domain=%s\n", domain)
	}
	ips := cache.Get(domain)
	ips += "," + ip
	cache.Set(domain, ips, 300) // 5min
}

type Domains struct {
	Domain []string
}
type ResDomain struct {
	Domain string
	Origin string
}

type Response struct {
	ISP     string
	Country string
	IP      string
}

func doc(w http.ResponseWriter, r *http.Request) {
	var d Domains
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")

	if r.Method == "OPTIONS" {
		w.Write([]byte("CORS :)"))
		return
	}

	if e := json.NewDecoder(r.Body).Decode(&d); e != nil {
		log.Printf(e.Error())
		http.Error(w, "failed to decode input", 400)
		return
	}
	fmt.Printf("In=%+v\n", d)

	// Filter out duplicate IPs
	uniqips := make(map[string]int)
	for _, domain := range d.Domain {
		vals := cache.Get(domain)
		for _, ip := range strings.Split(vals, ",") {
			n, _ := uniqips[ip]
			uniqips[ip] = n + 1
		}
	}

	out := make(map[uint]Response)
	// Humanize
	for ipstr, _ := range uniqips {
		if len(ipstr) == 0 {
			break
		}
		// Convert IPs to company list
		ip := net.ParseIP(ipstr)
		country, e := dbCity.Country(ip)
		if e != nil {
			log.Printf(e.Error())
			http.Error(w, "failed parsing IPs", 400)
			return
		}

		isp, e := dbASN.ASN(ip)
		if e != nil {
			log.Printf(e.Error())
			http.Error(w, "failed parsing IPs", 400)
			return
		}

		if _, ok := out[isp.AutonomousSystemNumber]; !ok {
			out[isp.AutonomousSystemNumber] = Response{
				ISP:     isp.AutonomousSystemOrganization,
				Country: country.Country.IsoCode,
				IP:      ipstr,
			}
		}

	}

	buf := new(bytes.Buffer)
	if e := json.NewEncoder(buf).Encode(out); e != nil {
		log.Printf(e.Error())
		http.Error(w, "encoding failed", 400)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, e := w.Write(buf.Bytes()); e != nil {
		log.Printf(e.Error())
	}
}

func main() {
	var (
		dns_addr  string
		http_addr string
	)
	flag.BoolVar(&Verbose, "v", false, "Verbose-mode (log more)")
	flag.StringVar(&dns_addr, "d", "[::]:53", "DNS listen on (both tcp and udp)")
	flag.StringVar(&http_addr, "h", "[::]:80", "HTTP listen on")
	flag.Parse()

	handler := &Handle{}
	cache = ttl_map.New("/tmp/leak.tsv")

	var err error
	dbCity, err = geoip2.Open("city.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer dbCity.Close()

	dbASN, err = geoip2.Open("asn.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer dbASN.Close()

	go func() {
		if err := dns.ListenAndServe(dns_addr, "udp", handler); err != nil {
			panic(err)
		}
	}()
	go func() {
		if err := dns.ListenAndServe(dns_addr, "tcp", handler); err != nil {
			panic(err)
		}
	}()

	http.HandleFunc("/", doc)
	if e := http.ListenAndServe(http_addr, nil); e != nil {
		panic(e)
	}
}
