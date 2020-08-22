# cisgo-ios
simple concurrent ssh server posing as cisco ios

## installation
install dependencies

```
go get github.com/gliderlabs/ssh
go get golang.org/x/crypto/ssh/terminal
```

## starting
```
go run cis.go 
2020/08/22 00:17:34 starting ssh server on port :10049
2020/08/22 00:17:34 starting ssh server on port :10023
2020/08/22 00:17:34 starting ssh server on port :10024
2020/08/22 00:17:34 starting ssh server on port :10000
2020/08/22 00:17:34 starting ssh server on port :10001
2020/08/22 00:17:34 starting ssh server on port :10025
2020/08/22 00:17:34 starting ssh server on port :10026
2020/08/22 00:17:34 starting ssh server on port :10027
```

## using
ssh into one of the open ports with ```admin``` as password and run "show version" or "show ip interface brief" or "show running-config"
```
test_device#show version
Cisco IOS XE Software, Version 16.04.01
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
```
