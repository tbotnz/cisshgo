package main

import (
	// "io"

	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

// ssh listernet
func sshListener(a int, done chan bool) {

	// SHOW_VERSION_PAGING_ENABLED := `Cisco IOS XE Software, Version 16.04.01
	// Cisco IOS Software [Everest], CSR1000V Software (X86_64_LINUX_IOSD-UNIVERSALK9-M), Version 16.4.1, RELEASE SOFTWARE (fc2)
	// Technical Support: http://www.cisco.com/techsupport
	// Copyright (c) 1986-2016 by Cisco Systems, Inc.
	// Compiled Sun 27-Nov-16 13:02 by mcpre
	// Cisco IOS-XE software, Copyright (c) 2005-2016 by cisco Systems, Inc.
	// All rights reserved.  Certain components of Cisco IOS-XE software are
	// licensed under the GNU General Public License ("GPL") Version 2.0.  The
	// software code licensed under GPL Version 2.0 is free software that comes
	// with ABSOLUTELY NO WARRANTY.  You can redistribute and/or modify such
	// GPL code under the terms of GPL Version 2.0.  For more details, see the
	// documentation or "License Notice" file accompanying the IOS-XE software,
	// or the applicable URL provided on the flyer accompanying the IOS-XE
	// software.
	// ROM: IOS-XE ROMMON
	// csr1000v uptime is 4 hours, 55 minutes
	// Uptime for this control processor is 4 hours, 57 minutes
	// System returned to ROM by reload
	// System image file is "bootflash:packages.conf"
	// Last reload reason: reload
	// This product contains cryptographic features and is subject to United
	// States and local country laws governing import, export, transfer and
	//  --More--
	// `

	supportedCommands := make(map[string]string)

	hostname := "test_device"

	supportedCommands["show version"] = `Cisco IOS XE Software, Version 16.04.01
Cisco IOS Software [Everest], CSR1000V Software (X86_64_LINUX_IOSD-UNIVERSALK9-M), Version 16.4.1, RELEASE SOFTWARE (fc2)
Technical Support: http://www.cisco.com/techsupport
Copyright (c) 1986-2016 by Cisco Systems, Inc.
Compiled Sun 27-Nov-16 13:02 by mcpre
Cisco IOS-XE software, Copyright (c) 2005-2016 by cisco Systems, Inc.
All rights reserved.  Certain components of Cisco IOS-XE software are
licensed under the GNU General Public License ("GPL") Version 2.0.  The
software code licensed under GPL Version 2.0 is free software that comes
with ABSOLUTELY NO WARRANTY.  You can redistribute and/or modify such
GPL code under the terms of GPL Version 2.0.  For more details, see the
documentation or "License Notice" file accompanying the IOS-XE software,
or the applicable URL provided on the flyer accompanying the IOS-XE
software.
ROM: IOS-XE ROMMON
csr1000v uptime is 4 hours, 55 minutes
Uptime for this control processor is 4 hours, 56 minutes
System returned to ROM by reload
System image file is "bootflash:packages.conf"
Last reload reason: reload
This product contains cryptographic features and is subject to United
States and local country laws governing import, export, transfer and
use. Delivery of Cisco cryptographic products does not imply
third-party authority to import, export, distribute or use encryption.
Importers, exporters, distributors and users are responsible for
compliance with U.S. and local country laws. By using this product you
agree to comply with applicable laws and regulations. If you are unable
to comply with U.S. and local laws, return this product immediately.
A summary of U.S. laws governing Cisco cryptographic products may be found at:
http://www.cisco.com/wwl/export/crypto/tool/stqrg.html
If you require further assistance please contact us by sending email to
export@cisco.com.
License Level: ax
License Type: Default. No valid license found.
Next reload license Level: ax
cisco CSR1000V (VXE) processor (revision VXE) with 2052375K/3075K bytes of memory.
Processor board ID 9FKLJWM5EB0
10 Gigabit Ethernet interfaces
32768K bytes of non-volatile configuration memory.
3985132K bytes of physical memory.
7774207K bytes of virtual hard disk at bootflash:.
0K bytes of  at webui:.
Configuration register is 0x2102`

	supportedCommands["show ip interface brief"] = `Interface                  IP-Address      OK? Method Status                Protocol
FastEthernet0/0            10.0.2.27       YES NVRAM  up                    up
Serial0/0                  unassigned      YES NVRAM  administratively down down
FastEthernet0/1            unassigned      YES NVRAM  administratively down down
Serial0/1                  unassigned      YES NVRAM  administratively down down
FastEthernet1/0            unassigned      YES NVRAM  administratively down down
FastEthernet2/0            unassigned      YES NVRAM  administratively down down
FastEthernet3/0            unassigned      YES unset  up                    down
FastEthernet3/1            unassigned      YES unset  up                    down
FastEthernet3/2            unassigned      YES unset  up                    down
FastEthernet3/3            unassigned      YES unset  up                    down
FastEthernet3/4            unassigned      YES unset  up                    down
FastEthernet3/5            unassigned      YES unset  up                    down
FastEthernet3/6            unassigned      YES unset  up                    down
FastEthernet3/7            unassigned      YES unset  up                    down
FastEthernet3/8            unassigned      YES unset  up                    down
FastEthernet3/9            unassigned      YES unset  up                    down
FastEthernet3/10           unassigned      YES unset  up                    down
FastEthernet3/11           unassigned      YES unset  up                    down
FastEthernet3/12           unassigned      YES unset  up                    down
FastEthernet3/13           unassigned      YES unset  up                    down
FastEthernet3/14           unassigned      YES unset  up                    down
FastEthernet3/15           unassigned      YES unset  up                    down
Vlan1                      unassigned      YES NVRAM  up                    down`

	ssh.Handle(func(s ssh.Session) {
		// io.WriteString(s, fmt.Sprintf(SHOW_VERSION_PAGING_DISABLED))
		term := terminal.NewTerminal(s, hostname+"#")
		for {
			line, err := term.ReadLine()
			if err != nil {
				break
			}
			response := line
			log.Println(line)
			if supportedCommands[response] != "" {
				term.Write(append([]byte(supportedCommands[response]), '\n'))
			} else if response == "" {
				term.Write(append([]byte(response)))
			} else if response == "exit" {
				break
			} else {
				term.Write(append([]byte("% Ambiguous command:  \""+response+"\""), '\n'))
			}
		}
		log.Println("terminal closed")
	})

	portString := strconv.Itoa(a)
	log.Printf("starting ssh server on port %s\n", portString)
	log.Fatal(ssh.ListenAndServe(fmt.Sprintf(":%s", portString), nil))

	done <- true
}

func main() {

	// Gather command line arguments
	listnersPtr := flag.Int("listners", 50, "How many listeners do you wish to spawn?")
	startingPortPtr := flag.Int("startingPort", 10000, "What port do you want to start at?")

	flag.Parse()
	listners := *startingPortPtr + *listnersPtr

	// Make a Channel for handling Goroutines, name of `done` expects a bool as return value
	done := make(chan bool, 1)

	// Iterate through the server ports and spawn a Goroutine for each
	for a := *startingPortPtr; a < listners; a++ {
		go sshListener(a, done)
	}

	// Recieve all the values from the channel (essentially wait on it to be empty)
	<-done
}
