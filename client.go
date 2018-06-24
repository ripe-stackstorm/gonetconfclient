package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	sshdriver "github.com/arsonistgopher/go-netconf/drivers/ssh"

	"golang.org/x/crypto/ssh"
)

const versionstring string = "0.2"

func main() {

	var cleanoutput = flag.Bool("cleanoutput", false, "Remove newline characters from XML")
	var version = flag.Bool("version", false, "Show version for code")
	var targethost = flag.String("targethost", "", "Target IPv4 or FQDN hostname of a NETCONF node")
	/*var transport = */ flag.String("transport", "ssh", "Transport mode (NOT IMPLEMENTED YET)")
	var envelope = flag.String("envelope", "", "XML Envelope")
	var envelopefile = flag.String("envelopefile", "", "XML envelope file path")
	var username = flag.String("username", "", "Username")
	var password = flag.String("password", "", "Password")
	var port = flag.Int("port", 830, "Port for accessing NETCONF")
	/*var sshkey = */ flag.String("sshkey", "", "SSHKey for accessing node (NOT IMPLEMENTED YET)")
	flag.Parse()

	// If version has been requested
	if *version {
		fmt.Printf("Version is: %s", versionstring)
		os.Exit(0)
	}

	nc := sshdriver.New()
	nc.Host = *targethost
	nc.Port = *port

	// Sort yourself out with SSH. Easiest to do that here.
	nc.SSHConfig = &ssh.ClientConfig{
		User:            *username,
		Auth:            []ssh.AuthMethod{ssh.Password(*password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	err := nc.Dial()

	if err != nil {
		log.Fatal(err)
	}

	// netconfcall content var
	netconfcall := ""
	// If the flag 'envelopefile' is set, then let's read the content into memory instead of using envelope
	if *envelopefile != "" {
		// Read content from disk
		bytes, err := ioutil.ReadFile(*envelopefile)
		netconfcall = string(bytes)
		if err != nil {
			log.Panic(err)
		}
		// Else, consume content of envelope
	} else {
		netconfcall = *envelope
	}

	reply, err := nc.SendRaw(netconfcall)
	if err != nil {
		panic(err)
	}

	if *cleanoutput == true {
		d1 := strings.Replace(reply.Data, "\n", "", -1)
		fmt.Print(d1 + "\n")
	}
	if *cleanoutput == false {
		d1 := strings.TrimLeft(reply.Data, "\n")
		fmt.Print(d1)
	}

}
