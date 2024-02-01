**ABSOLUTELY NO WARRANTY! USE AT YOUR OWN RISK!**

# Hörmann BiSecur Gateway CLI client and GoLang SDK

Goal is to create a fully fledged replacement client after Hörmann stated end of life of their cloud and android application.

If all goes fine, later this repository will provide you both a GoLang SDK and a CLI client you can use to
- open and close your door
- manage users
- etc.

![gateway image](gateway.webp)

## TODOs
* [x] Create json output for machines. Improve documentation accordingly.
* [x] Add retries. `status` command produces `PORT_ERROR` quite frequently while second try works fine.
* [x] Improve token handling. Token is stored in `config.yaml` but it seems to be invalidated after a while. It should be renewed on demand.
* [x] Create GitHub pipeline for releases
* [x] Create new gateway user
* [x] Delete gateway user
* [x] Change password of a gateway user
* [ ] Assign new door to the gateway
* [ ] Delete assigned door from the gateway

## Usage
```
$ ./halsecur 
Application to manage your Hörmann BiSecur gateway without the central cloud directly on your LAN.

Usage:
  halsecur [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  discover    Discover Hörmann BiSecur gateways on the local network
  get-name    Queries the name of the Hörmann BiSecur gateway
  groups      Manages doors defined in your Hörmann BiSecur gateway.
  help        Help about any command
  login       
  logout      
  ping        Check if your Hörmann BiSecur gateway is reachable or not.
  set-state   Open or close a door connected to your Hörmann BiSecur gateway.
  status      Queries the status (open/closed/etc) of your door.
  users       Manages users defined in your Hörmann BiSecur gateway.
  version     Print version information

Flags:
      --autologin         login automatically on demand (default true)
      --debug             debug log level (default true)
  -h, --help              help for halsecur
      --host string       IP or host name or the Hörmann BiSecure gateway
      --json              use json logging format instead of human readable
      --mac string        MAC address of the Hörmann BiSecur gateway
      --password string   Valid password belongs to the given username
      --port int           (default 4000)
      --token uint32      Valid authentication token
      --username string   Valid username

Use "halsecur [command] --help" for more information about a command.
```

### Ping
```bash
$ dist/halsecur ping --host 192.168.3.232 --mac 54:10:EC:85:28:BB --count 3 --delay 1000
INFO[2024-01-31T21:12:46+01:00] Response 1 of 3 received in 67 ms
INFO[2024-01-31T21:12:47+01:00] Response 2 of 3 received in 64 ms
INFO[2024-01-31T21:12:48+01:00] Response 3 of 3 received in 63 ms
```

### Get device name
```bash
$ dist/halsecur get-name
INFO[2024-01-31T21:08:47+01:00] Received name: BiSecur Gateway
```

### Login
```bash
$ ./dist/halsecur login --host 192.168.3.232 --mac 54:10:EC:85:28:BB --password Gabor123456789. --username app
INFO[2024-01-31T21:09:40+01:00] Token: 0x3AC29326
INFO[2024-01-31T21:09:40+01:00] Success
```

### Get users
```bash
$ ./dist/halsecur users list
INFO[2024-01-31T21:56:53+01:00] [{"id":0,"name":"admin","isAdmin":true,"Groups":[]},{"id":1,"name":"app","isAdmin":false,"Groups":[0]}]
```

### Get groups
```bash
$ ./dist/halsecur groups list
INFO[2024-02-01T17:35:32+01:00] [{"id":0,"name":"garazs","ports":[{"typeName":"IMPULS","id":0,"type":1}]}] 
```

### Get door status
```bash
$ ./dist/halsecur status --devicePort 0
INFO[2024-02-01T17:34:22+01:00] Token expired. Logging in...                 
INFO[2024-02-01T17:34:22+01:00] Token: 0xA0A67B43                            
INFO[2024-02-01T17:34:24+01:00] Transition: {"StateInPercent":0,"DesiredStateInPercent":0,"Error":false,"AutoClose":false,"DriveTime":0,"Gk":257,"Hcp":{"PositionOpen":false,"PositionClose":true,"OptionRelais":false,"LightBarrier":false,"Error":false,"DrivingToClose":false,"Driving":false,"HalfOpened":false,"ForecastLeadTime":false,"Learned":true,"NotReferenced":false},"Exst":"AAAAAAAAAAA=","Time":"2024-02-01T17:34:24.794359108+01:00"} 

```

### Open/close door
Door is fully closed. Start to open it:

```bash
$ ./dist/halsecur set-state --devicePort 0
DEBU[2023-12-26T21:04:33+01:00] Connecting to 192.168.3.232:4000             
DEBU[2023-12-26T21:04:33+01:00] Request: SrcMAC=0x000000000009, DstMAC=0x5410EC8528BB, BodyLength=0x0, packet=[Tag=0x1, Token=0x7974DB57, CommandID=0x33 (0x33), payload=[SetState], Checksum=0x0, isResponse=false], Checksum=0x0, isResponse: false 
DEBU[2023-12-26T21:04:33+01:00] Request bytes: 303030303030303030303039353431304543383532384242303030423031373937344442353733333030464635444331 
DEBU[2023-12-26T21:04:34+01:00] Length of received bytes: 76                 
DEBU[2023-12-26T21:04:34+01:00] Response bytes: 5410EC8528BB00000000000600190100000000F0000000040101020200000000000000001483 
DEBU[2023-12-26T21:04:34+01:00] Received TC: SrcMAC=0x5410EC8528BB, DstMAC=0x000000000006, BodyLength=0x19, packet=[Tag=0x1, Token=0x0, CommandID=0x70 (0xF0), payload=[HmGetTransitionResponse[StateInPercent: 0, DesiredStateInPerced: 0, Error: false, AutoClose: false, DriveTime: 4, Gk: 257, Hcp: HCP[PositionOpen: false, PositionClose: true, OptionRelais: false, LightBarrier: false, Error: false, DrivingToClose: false, Driving: false, HalfOpened: false, ForecastLeadTime: false, Learned: true, NotReferenced: false], Exst: [0 0 0 0 0 0 0 0], Time: 2023-12-26 21:04:34.436781605 +0100 CET m=+0.770579311]], Checksum=0x14, isResponse=true], Checksum=0x83, isResponse: true 
DEBU[2023-12-26T21:04:34+01:00] Set State response: SrcMAC=0x5410EC8528BB, DstMAC=0x000000000006, BodyLength=0x19, packet=[Tag=0x1, Token=0x0, CommandID=0x70 (0xF0), payload=[HmGetTransitionResponse[StateInPercent: 0, DesiredStateInPerced: 0, Error: false, AutoClose: false, DriveTime: 4, Gk: 257, Hcp: HCP[PositionOpen: false, PositionClose: true, OptionRelais: false, LightBarrier: false, Error: false, DrivingToClose: false, Driving: false, HalfOpened: false, ForecastLeadTime: false, Learned: true, NotReferenced: false], Exst: [0 0 0 0 0 0 0 0], Time: 2023-12-26 21:04:34.436781605 +0100 CET m=+0.770579311]], Checksum=0x14, isResponse=true], Checksum=0x83, isResponse: true 
INFO[2023-12-26T21:04:34+01:00] Done        
```

Door is opening. Stop it in a half-open position:

```Bash
$ ./dist/halsecur set-state --devicePort 0
DEBU[2023-12-26T21:04:43+01:00] Connecting to 192.168.3.232:4000             
DEBU[2023-12-26T21:04:43+01:00] Request: SrcMAC=0x000000000009, DstMAC=0x5410EC8528BB, BodyLength=0x0, packet=[Tag=0x1, Token=0x7974DB57, CommandID=0x33 (0x33), payload=[SetState], Checksum=0x0, isResponse=false], Checksum=0x0, isResponse: false 
DEBU[2023-12-26T21:04:43+01:00] Request bytes: 303030303030303030303039353431304543383532384242303030423031373937344442353733333030464635444331 
DEBU[2023-12-26T21:04:43+01:00] Length of received bytes: 76                 
DEBU[2023-12-26T21:04:43+01:00] Response bytes: 5410EC8528BB00000000000600190100000000F00000000401014C0200000000000000005EAD 
DEBU[2023-12-26T21:04:43+01:00] Received TC: SrcMAC=0x5410EC8528BB, DstMAC=0x000000000006, BodyLength=0x19, packet=[Tag=0x1, Token=0x0, CommandID=0x70 (0xF0), payload=[HmGetTransitionResponse[StateInPercent: 0, DesiredStateInPerced: 0, Error: false, AutoClose: false, DriveTime: 4, Gk: 257, Hcp: HCP[PositionOpen: false, PositionClose: false, OptionRelais: true, LightBarrier: true, Error: false, DrivingToClose: false, Driving: true, HalfOpened: false, ForecastLeadTime: false, Learned: true, NotReferenced: false], Exst: [0 0 0 0 0 0 0 0], Time: 2023-12-26 21:04:43.979174613 +0100 CET m=+0.807240492]], Checksum=0x5E, isResponse=true], Checksum=0xAD, isResponse: true 
DEBU[2023-12-26T21:04:43+01:00] Set State response: SrcMAC=0x5410EC8528BB, DstMAC=0x000000000006, BodyLength=0x19, packet=[Tag=0x1, Token=0x0, CommandID=0x70 (0xF0), payload=[HmGetTransitionResponse[StateInPercent: 0, DesiredStateInPerced: 0, Error: false, AutoClose: false, DriveTime: 4, Gk: 257, Hcp: HCP[PositionOpen: false, PositionClose: false, OptionRelais: true, LightBarrier: true, Error: false, DrivingToClose: false, Driving: true, HalfOpened: false, ForecastLeadTime: false, Learned: true, NotReferenced: false], Exst: [0 0 0 0 0 0 0 0], Time: 2023-12-26 21:04:43.979174613 +0100 CET m=+0.807240492]], Checksum=0x5E, isResponse=true], Checksum=0xAD, isResponse: true 
INFO[2023-12-26T21:04:43+01:00] Done
```

Door is half-opened. Close it back:

```Bash
$ ./dist/halsecur set-state --devicePort 0
DEBU[2023-12-26T21:04:47+01:00] Connecting to 192.168.3.232:4000             
DEBU[2023-12-26T21:04:47+01:00] Request: SrcMAC=0x000000000009, DstMAC=0x5410EC8528BB, BodyLength=0x0, packet=[Tag=0x1, Token=0x7974DB57, CommandID=0x33 (0x33), payload=[SetState], Checksum=0x0, isResponse=false], Checksum=0x0, isResponse: false 
DEBU[2023-12-26T21:04:47+01:00] Request bytes: 303030303030303030303039353431304543383532384242303030423031373937344442353733333030464635444331 
DEBU[2023-12-26T21:04:48+01:00] Length of received bytes: 76                 
DEBU[2023-12-26T21:04:48+01:00] Response bytes: 5410EC8528BB00000000000600190100000000F00000000401010C0200000000000000001EA5 
DEBU[2023-12-26T21:04:48+01:00] Received TC: SrcMAC=0x5410EC8528BB, DstMAC=0x000000000006, BodyLength=0x19, packet=[Tag=0x1, Token=0x0, CommandID=0x70 (0xF0), payload=[HmGetTransitionResponse[StateInPercent: 0, DesiredStateInPerced: 0, Error: false, AutoClose: false, DriveTime: 4, Gk: 257, Hcp: HCP[PositionOpen: false, PositionClose: false, OptionRelais: true, LightBarrier: true, Error: false, DrivingToClose: false, Driving: false, HalfOpened: false, ForecastLeadTime: false, Learned: true, NotReferenced: false], Exst: [0 0 0 0 0 0 0 0], Time: 2023-12-26 21:04:48.329008833 +0100 CET m=+0.781378758]], Checksum=0x1E, isResponse=true], Checksum=0xA5, isResponse: true 
DEBU[2023-12-26T21:04:48+01:00] Set State response: SrcMAC=0x5410EC8528BB, DstMAC=0x000000000006, BodyLength=0x19, packet=[Tag=0x1, Token=0x0, CommandID=0x70 (0xF0), payload=[HmGetTransitionResponse[StateInPercent: 0, DesiredStateInPerced: 0, Error: false, AutoClose: false, DriveTime: 4, Gk: 257, Hcp: HCP[PositionOpen: false, PositionClose: false, OptionRelais: true, LightBarrier: true, Error: false, DrivingToClose: false, Driving: false, HalfOpened: false, ForecastLeadTime: false, Learned: true, NotReferenced: false], Exst: [0 0 0 0 0 0 0 0], Time: 2023-12-26 21:04:48.329008833 +0100 CET m=+0.781378758]], Checksum=0x1E, isResponse=true], Checksum=0xA5, isResponse: true 
INFO[2023-12-26T21:04:48+01:00] Done
```

## Acknowledgement

Thanks for [SEC Consult](https://sec-consult.com/blog/detail/hoermann-opening-doors-for-everyone/) for the initial analysis and documentation.

Based on the above study someone could create a [Kotlin SDK](https://github.com/bisdk/sdk) which also helped me a lot.

Taken into consideration that Hörmann will stop their cloud required for BiSecur Gateway usages from 2024, I asked Hörmann support to publish their protocol already leaked in above repositories, but they stated that this code is impossible to write.

## Disclaimer

**ABSOLUTELY NO WARRANTY! USE AT YOUR OWN RISK!**

This software is provided by Halacs "as is" and "with all faults." Halacs makes no representations or warranties of any kind concerning the safety, suitability, lack of viruses, inaccuracies, typographical errors, or other harmful components of this software. There are inherent dangers in the use of any software, and you are solely responsible for determining whether this software is compatible with your equipment and other software installed on your equipment. You are also solely responsible for the protection of your equipment and backup of your data, and Halacs will not be liable for any damages you may suffer in connection with using, modifying, or distributing this software.

**NO LIABILITY FOR CONSEQUENTIAL DAMAGES**

To the maximum extent permitted by applicable law, in no event shall Halacs be liable for any direct, indirect, punitive, incidental, special, consequential damages or any damages whatsoever including, without limitation, damages for loss of use, data, or profits arising out of or in any way connected with the use or performance of this software, with the delay or inability to use this software, or for any information obtained through this software.

**NO RESPONSIBILITY FOR THIRD-PARTY COMPONENTS**

This software may include third-party software components subject to their own licenses. Halacs disclaims any responsibility or liability related to the inclusion, functionality, or use of such third-party components.

**DISCLAIMER SPECIFIC TO BiSecur Gateway**

By using this software, you acknowledge that it was developed with the best intentions; however, it may lead to the BiSecur Gateway device becoming inoperable, or the associated physical door may open or remain open. Halacs explicitly disclaims any responsibility or liability for such events. It is your responsibility to ensure the proper functioning and security of the BiSecur Gateway device and the connected physical door. If you do not agree to these terms, you may not use, modify, or distribute this software.

## Links
- https://sec-consult.com/blog/detail/hoermann-opening-doors-for-everyone/
- https://github.com/bisdk/sdk