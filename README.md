# cisshgo
Simple, small, fast, concurrent SSH server to emulate network equipment (i.e. Cisco IOS) for testing purposes.

## Usage

1. Clone the repository and change into that directory (All dependencies are included in the `/vendor` folder, so no installation step is necessary.)
2. Execute `go run cis.go` as shown below:

```bash
$ go run cis.go 
2020/08/22 00:17:34 starting ssh server on port :10049
2020/08/22 00:17:34 starting ssh server on port :10023
2020/08/22 00:17:34 starting ssh server on port :10024
... <snip>
```

Alternatively you can compile and run in separate steps (useful for docker images, etc):

```bash
$ go build cisgo-ios cis.go
$ ./cisgo-ios
2020/09/02 15:46:31 starting cis.go ssh server on port :10008
2020/09/02 15:46:31 starting cis.go ssh server on port :10005
2020/09/02 15:46:31 starting cis.go ssh server on port :10000
2020/09/02 15:46:31 starting cis.go ssh server on port :10006
... <snip>
```

3. SSH into one of the open ports with `admin` as the password. By default, you can run "show version"
 or "show ip interface brief" or "show running-config". Other commands can be added by modifying the
 transcript_map.yaml file and supplying transcripts as needed.

Example output:

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

## Advanced Usage

There are several options available to control the behavior
 of cisshgo see the below output of `-help`:

```
  -listners int
    	How many listeners do you wish to spawn? (default 50)
  -startingPort int
    	What port do you want to start at? (default 10000)
  -transcriptMap string
    	What file contains the map of commands to transcipted output? (default "transcripts/transcript_map.yaml")
```

For example, if you only wish to lauch with a single SSH listner for a testing process,
 you could simply apply `-listners 1` to the run command:

```
go run cis.go -listners 1
2020/09/03 19:41:04 Starting cis.go ssh server on port :10000
```

## Expanding Platform Support

cisgo-ios is built modularly to support easy expansion or customization. Potential options for enhancement are outlined below.

### Customized Output in Command Transcripts

If you wish to modify elements of the transcript dynamically, for example the hostname,
 you can instantiate templateable sections into your transcript.

For example, in the packaged output of `show_version.txt` the hostname is listed as:

```
ROM: IOS-XE ROMMON
{{.Hostname}} uptime is 4 hours, 55 minutes
Uptime for this control processor is 4 hours, 56 minutes
```

Any value in the `fakedevices.FakeDevice` struct can be referenced in this way, today these are:

```
type FakeDevice struct {
	Vendor            string            // Vendor of this fake device
	Platform          string            // Platform of this fake device
	Hostname          string            // Hostname of the fake device
	Password          string            // Password of the fake device
	SupportedCommands SupportedCommands // What commands this fake device supports
	ContextSearch     map[string]string // The available CLI prompt/contexts on this fake device
	ContextHierarchy  map[string]string // The heiarchy of the available contexts
}
```

If you wish to template additional/different values, they will need to be added to the FakeDevice struct
 and then instantiated in the transcript with a reference to `{{.MyNewAttribute}}`.

### Adding Additional Command Transcripts

If you wish to add additional command transcripts, you simply need to include a plain text file in the appropriate
 `vendor/platform` folder, and create an entry in the `transcript_map.yaml` file under the appropriate vendor/platform:

```
---
platforms:
  - csr1000v:
      command_transcripts:
        "my new fancy command": "transcripts/cisco/csr1000v/my_new_fancy_command.txt"
```

On the next execution of cisgo-ios it will read this map and respond to `my new fancy command`

### Adding Additional "Cisco-style" Platforms

If you wish to add a completely new Cisco-style device, that is one with `configure terminal`
 leading to a `(config)#` mode for example, you can simply supply additional device types and transcripts
 in the transcript_map.yaml file.

This however does not work if a device follows a different interaction pattern than the Cisco standard,
 for example a Juniper or F5 device, for that see the following section.

### Adding Additional Non-"Cisco-style" Platforms

**NOTE** This feature is not fully implemented yet!

If you wish to add a platform that is _not_ the "Cisco-style" of interaction, for example a Juniper or F5 device,
 you will need to implement a new `handler` module for it under `ssh_server/handlers` and add it to the 
 device mapping in code in `cis.go` where it chooses the SSH listner and handler.

The "handler" controls the basics of how we will emulate the SSH session, and provides a list of
 `if...else if...else if...` options to roughly simulate the device experience. Because many network
  devices vary in their CLI and interactions, the conditional tree that each requires will vary.
  This is implemented via the "handler" functionality.

### Disclaimer
Cisco IOS is the property/trademark of Cisco.
