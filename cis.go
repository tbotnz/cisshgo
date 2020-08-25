package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

// ssh listernet
func sshListener(portNumber int, done chan bool) {

	basicCommands := make(map[string]string)
	specialCommands := make(map[string]string)
	contextSearch := make(map[string]string)
	contextHierarchy := make(map[string]string)

	contextSearch["conf t"] = "(config)#"
	contextSearch["configure terminal"] = "(config)#"
	contextSearch["configure t"] = "(config)#"
	contextSearch["enable"] = "#"
	contextSearch["en"] = "#"
	contextSearch["base"] = ">"

	contextHierarchy["(config)#"] = "#"
	contextHierarchy["#"] = ">"
	contextHierarchy[">"] = "exit"

	defaultHostname := "cisgo1000v"
	defaultContextState := ">"
	hostname := defaultHostname
	password := "admin"

	specialCommands["show version"] = `Cisco IOS XE Software, Version 16.04.01
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
%[1]v uptime is 4 hours, 55 minutes
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
Configuration register is 0x2102
`

	basicCommands["show ip interface brief"] = `Interface                  IP-Address      OK? Method Status                Protocol
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

	specialCommands["show running-config"] = `Building configuration...

Current configuration : 2114 bytes
!
version 12.4
service timestamps debug datetime msec
service timestamps log datetime msec
no service password-encryption
!
hostname %[1]s
!
boot-start-marker
boot-end-marker
!
!
no aaa new-model
memory-size iomem 5
no ip icmp rate-limit unreachable
ip cef
!
!
!
!
no ip domain lookup
ip domain name test
ip auth-proxy max-nodata-conns 3
ip admission max-nodata-conns 3
!
!
!
!
!
!
!
!
!
!
!
!
!
!
!
username admin privilege 15 secret 5 $1$M1ce$SKeVGg2lUCPrsLkJMIdWf.
!
!
ip tcp synwait-time 5
ip ssh version 2
ip scp server enable
!
!
!
!
!
interface FastEthernet0/0
 description netpalm
 ip address 10.0.2.27 255.255.255.0
 duplex auto
 speed auto
!
interface Serial0/0
 no ip address
 shutdown
 clock rate 2000000
!
interface FastEthernet0/1
 no ip address
 shutdown
 duplex auto
 speed auto
!
interface Serial0/1
 no ip address
 shutdown
 clock rate 2000000
!
interface FastEthernet1/0
 no ip address
 shutdown
 duplex auto
 speed auto
!
interface FastEthernet2/0
 no ip address
 shutdown
 duplex auto
 speed auto
!
interface FastEthernet3/0
!
interface FastEthernet3/1
!
interface FastEthernet3/2
!
interface FastEthernet3/3
!
interface FastEthernet3/4
!
interface FastEthernet3/5
!
interface FastEthernet3/6
!
interface FastEthernet3/7
!
interface FastEthernet3/8
!
interface FastEthernet3/9
!
interface FastEthernet3/10
!
interface FastEthernet3/11
!
interface FastEthernet3/12
!
interface FastEthernet3/13
!
interface FastEthernet3/14
!
interface FastEthernet3/15
!
interface Vlan1
 no ip address
!
ip forward-protocol nd
!
!
no ip http server
no ip http secure-server
!
ip access-list standard bob
ip access-list standard yip
!
snmp-server community test RO
snmp-server community location RO yip
snmp-server community contact RO bob
no cdp log mismatch duplex
!
!
!
control-plane
!
!
!
!
!
!
!
!
!
!
line con 0
 exec-timeout 0 0
 privilege level 15
 logging synchronous
line aux 0
 exec-timeout 0 0
 privilege level 15
 logging synchronous
line vty 0 4
 privilege level 15
 login local
 transport input ssh
line vty 5 15
 privilege level 15
 login local
 transport input ssh
!
!
end
`
	basicCommands["terminal length 0"] = " "
	basicCommands["terminal width 511"] = " "
	ssh.Handle(func(s ssh.Session) {

		// io.WriteString(s, fmt.Sprintf(SHOW_VERSION_PAGING_DISABLED))
		term := terminal.NewTerminal(s, hostname+contextSearch["base"])
		contextState := defaultContextState
		for {
			line, err := term.ReadLine()
			if err != nil {
				break
			}

			response := line
			log.Println(line)
			if response == "reset state" {
				log.Println("resetting internal state")
				contextState = defaultContextState
				hostname = defaultHostname
				term.SetPrompt(hostname + contextState)

			} else if basicCommands[response] != "" {
				// lookup supported commands for response
				term.Write(append([]byte(basicCommands[response]), '\n'))

			} else if response == "" {
				// return if nothing is entered
				term.Write(append([]byte(response)))

			} else if contextSearch[response] != "" {
				// switch contexts as needed
				term.SetPrompt(string(hostname + contextSearch[response]))
				contextState = contextSearch[response]

			} else if response == "exit" || response == "end" {
				// drop down configs if required
				if contextHierarchy[contextState] == "exit" {
					break

				} else {
					term.SetPrompt(string(hostname + contextHierarchy[contextState]))
					contextState = contextHierarchy[contextState]
				}

			} else if specialCommands[response] != "" {
				term.Write(append([]byte(fmt.Sprintf(specialCommands[response], hostname))))

			} else if contextState != ">" { // we're in config mode
				fields := strings.Fields(response)
				command := fields[0]
				if command == "hostname" {
					hostname = strings.Join(fields[1:], " ")
					log.Printf("Setting hostname to %s\n", hostname)
					term.SetPrompt(hostname + contextState)

				} else {
					term.Write(append([]byte("% Ambiguous command:  \""+response+"\""), '\n'))
				}

			} else {
				term.Write(append([]byte("% Ambiguous command:  \""+response+"\""), '\n'))
			}

		}
		log.Println("terminal closed")

	})

	portString := ":" + strconv.Itoa(portNumber)
	//prt :=  portString
	log.Printf("starting cis.go ssh server on port %s\n", portString)

	log.Fatal(ssh.ListenAndServe(portString, nil,
		ssh.PasswordAuth(func(ctx ssh.Context, pass string) bool {
			return pass == password
		}),
	))

	done <- true
}

func main() {
	done := make(chan bool, 1)
	for portNumber := 10000; portNumber < 10050; portNumber++ {
		go sshListener(portNumber, done)
	}
	<-done
}
