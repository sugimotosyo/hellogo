package main

import (
	"fmt"
	"net"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	ifaces, _ := net.Interfaces()
	ips := ""
	// handle err
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			ips = ips + fmt.Sprint(ip) + ":"
			// process IP address
		}
	}

	fmt.Fprintf(w, "from pipeline Hello, World %s", ips)
}

func main() {
	http.HandleFunc("/", handler) // ハンドラを登録してウェブページを表示させる
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("error ======")
		return
	}

}
