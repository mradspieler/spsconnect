package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/mradspieler/tosi"
)

func main() {
	sps := flag.String("sps", "", "SPS host address")
	sport := flag.Int("sport", 102, "SPS port")
	host := flag.String("host", "", "hostname")
	lport := flag.Int("lport", -1, "local port")
	ltsap := flag.String("ltsap", "", "local TSAP")
	rtsap := flag.String("rtsap", "IPSTIPST", "remote TSAP")

	flag.Parse()

	spsAddr, err := rfc1006.ResolveTOSIAddr("tosi", *sps+":"+strconv.Itoa(*sport)+":"+*rtsap)
	if err != nil {
		errorf("1 - Resolve PLC address not possible: %v", err)
	}

	if *host == "" && *lport == -1 && *ltsap == "" {
		conn, err := rfc1006.DialTOSI("tosi", nil, spsAddr)
		if err != nil {
			errorf("2 - Connect to PLC not possible: %v", err)
		}
		defer conn.Close()
	}

	if *host != "" && *lport != -1 && *ltsap != "" {
		locAddr, err := rfc1006.ResolveTOSIAddr("tosi", *host+":"+strconv.Itoa(*lport)+":"+*ltsap)
		if err != nil {
			errorf("3 - Resolve PLC address not possible: %v", err)
		}
		conn, err := rfc1006.DialTOSI("tosi", locAddr, spsAddr)
		if err != nil {
			errorf("4 - Connect to PLC not possible: %v", err)
		}

		defer conn.Close()
	}

	if *host == "" && *lport == -1 && *ltsap != "" {
		locAddr := &rfc1006.TOSIAddr{}
		locAddr.TSel = []byte(*ltsap)
		conn, err := rfc1006.DialTOSI("tosi", locAddr, spsAddr)
		if err != nil {
			errorf("3 - Connect to PLC not possible: %v", err)
		}
		defer conn.Close()
	}

	fmt.Printf("\n=> Hooray! Connection to SPS successful\n")
}

func errorf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(2)
}
