**ABSOLUTELY NO WARRANTY! USE AT YOUR OWN RISK!**

# Hörmann BiSecur Gateway CLI client and GoLang SDK

Goal is to create a fully fledged replacement client after Hörmann stated end of life of their cloud and android application.

If all goes fine, later this repository will provide you both a GoLang SDK and a CLI client you can use to
- open and close your door
- manage users
- etc.

![gateway image](gateway.webp)

## TODOs
* [ ] Create json output for machines. Improve documentation accordingly.
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
      --mac string        MAC address of the Hörmann BiSecur gateway
      --password string   Valid password belongs to the given username
      --port int           (default 4000)
      --token uint32      Valid authentication token
      --username string   Valid username

Use "halsecur [command] --help" for more information about a command.
```

### Ping
```bash
$ ./halsecur ping --host 192.168.3.232 --mac 54:10:EC:85:28:BB --count 3 --delay 1000
INFO[2024-01-14T10:18:54+01:00] Response 1 of 3 received in 73 ms
INFO[2024-01-14T10:18:55+01:00] Response 2 of 3 received in 75 ms
INFO[2024-01-14T10:18:56+01:00] Response 3 of 3 received in 76 ms
```

### Get device name
```bash
$ ./dist/halsecur get-name
DEBU[2023-12-26T21:02:50+01:00] Connecting to 192.168.3.232:4000             
DEBU[2023-12-26T21:02:50+01:00] Request: SrcMAC=0x000000000009, DstMAC=0x5410EC8528BB, BodyLength=0x0, packet=[Tag=0x1, Token=0x7974DB57, CommandID=0x26 (0x26), payload=[], Checksum=0x0, isResponse=false], Checksum=0x0, isResponse: false 
DEBU[2023-12-26T21:02:50+01:00] Request bytes: 3030303030303030303030393534313045433835323842423030303930313739373444423537323634464346 
DEBU[2023-12-26T21:02:50+01:00] Length of received bytes: 74                 
DEBU[2023-12-26T21:02:50+01:00] Response bytes: 5410EC8528BB0000000000060018017974DB57A64269536563757220476174657761797D11 
DEBU[2023-12-26T21:02:50+01:00] Received TC: SrcMAC=0x5410EC8528BB, DstMAC=0x000000000006, BodyLength=0x18, packet=[Tag=0x1, Token=0x7974DB57, CommandID=0x26 (0xA6), payload=[GetNameResponse: BiSecur Gateway], Checksum=0x7D, isResponse=true], Checksum=0x11, isResponse: true 
INFO[2023-12-26T21:02:50+01:00] Received name: BiSecur Gateway
```

### Login
```bash
$ ./dist/halsecur login --host 192.168.3.232 --mac 54:10:EC:85:28:BB --password Gabor123456789. --username app --debug
DEBU[2023-12-26T21:00:19+01:00] Connecting to 192.168.3.232:4000             
DEBU[2023-12-26T21:00:19+01:00] Request: SrcMAC=0x000000000009, DstMAC=0x5410EC8528BB, BodyLength=0x0, packet=[Tag=0x1, Token=0x0, CommandID=0x10 (0x10), payload=[appGabor123456789.], Checksum=0x0, isResponse=false], Checksum=0x0, isResponse: false 
DEBU[2023-12-26T21:00:19+01:00] Request bytes: 30303030303030303030303935343130454338353238424230303143303130303030303030303130303336313730373034373631363236463732333133323333333433353336333733383339324536373446 
DEBU[2023-12-26T21:00:19+01:00] Length of received bytes: 54                 
DEBU[2023-12-26T21:00:19+01:00] Response bytes: 5410EC8528BB000000000006000E010000000090017974DB57BFC8 
DEBU[2023-12-26T21:00:19+01:00] Received TC: SrcMAC=0x5410EC8528BB, DstMAC=0x000000000006, BodyLength=0xE, packet=[Tag=0x1, Token=0x0, CommandID=0x10 (0x90), payload=[SenderID: 0x1, Token: 0x7974DB57], Checksum=0xBF, isResponse=true], Checksum=0xC8, isResponse: true 
INFO[2023-12-26T21:00:19+01:00] Token: 0x7974DB57
```

### Get users
```bash
$ ./dist/halsecur users list
DEBU[2023-12-26T21:01:38+01:00] Connecting to 192.168.3.232:4000             
DEBU[2023-12-26T21:01:38+01:00] Request: SrcMAC=0x000000000009, DstMAC=0x5410EC8528BB, BodyLength=0x0, packet=[Tag=0x1, Token=0x7974DB57, CommandID=0x6 (0x6), payload=[Jcmp: {"CMD":"GET_USERS"}], Checksum=0x0, isResponse=false], Checksum=0x0, isResponse: false 
DEBU[2023-12-26T21:01:38+01:00] Request bytes: 30303030303030303030303935343130454338353238424230303143303137393734444235373036374232323433344434343232334132323437343535343546353535333435353235333232374441314431 
DEBU[2023-12-26T21:01:38+01:00] Length of received bytes: 250                
DEBU[2023-12-26T21:01:38+01:00] Response bytes: 5410EC8528BB0000000000060070017974DB57865B7B226964223A302C226E616D65223A2261646D696E222C22697341646D696E223A747275652C2267726F757073223A5B5D7D2C7B226964223A312C226E616D65223A22617070222C22697341646D696E223A66616C73652C2267726F757073223A5B305D7D5D26EF 
DEBU[2023-12-26T21:01:38+01:00] Received TC: SrcMAC=0x5410EC8528BB, DstMAC=0x000000000006, BodyLength=0x70, packet=[Tag=0x1, Token=0x7974DB57, CommandID=0x6 (0x86), payload=[Jcmp: [{"id":0,"name":"admin","isAdmin":true,"groups":[]},{"id":1,"name":"app","isAdmin":false,"groups":[0]}]], Checksum=0x26, isResponse=true], Checksum=0xEF, isResponse: true 
INFO[2023-12-26T21:01:38+01:00] Users: [ID=0, Name="admin", IsAdmin=true, Groups:[]][ID=1, Name="app", IsAdmin=false, Groups:[0]] 
```

### Get groups
```bash
$ ./dist/halsecur groups list
DEBU[2023-12-26T21:02:20+01:00] Connecting to 192.168.3.232:4000             
DEBU[2023-12-26T21:02:20+01:00] Request: SrcMAC=0x000000000009, DstMAC=0x5410EC8528BB, BodyLength=0x0, packet=[Tag=0x1, Token=0x7974DB57, CommandID=0x6 (0x6), payload=[Jcmp: {"CMD":"GET_GROUPS"}], Checksum=0x0, isResponse=false], Checksum=0x0, isResponse: false 
DEBU[2023-12-26T21:02:20+01:00] Request bytes: 303030303030303030303039353431304543383532384242303031443031373937344442353730363742323234333444343432323341323234373435353435463437353234463535353035333232374446303446 
DEBU[2023-12-26T21:02:20+01:00] Length of received bytes: 152                
DEBU[2023-12-26T21:02:20+01:00] Response bytes: 5410EC8528BB000000000006003F017974DB57865B7B226964223A302C226E616D65223A22676172617A73222C22706F727473223A5B7B226964223A302C2274797065223A317D5D7D5DD1EF 
DEBU[2023-12-26T21:02:20+01:00] Received TC: SrcMAC=0x5410EC8528BB, DstMAC=0x000000000006, BodyLength=0x3F, packet=[Tag=0x1, Token=0x7974DB57, CommandID=0x6 (0x86), payload=[Jcmp: [{"id":0,"name":"garazs","ports":[{"id":0,"type":1}]}]], Checksum=0xD1, isResponse=true], Checksum=0xEF, isResponse: true 
INFO[2023-12-26T21:02:20+01:00] Groups: ID=0 Name="garazs" Ports=[ID=0 Type=IMPULS] 
```

### Get door status
```bash
$ ./dist/halsecur status --devicePort 0
DEBU[2023-12-26T21:03:45+01:00] Connecting to 192.168.3.232:4000             
DEBU[2023-12-26T21:03:45+01:00] Request: SrcMAC=0x000000000009, DstMAC=0x5410EC8528BB, BodyLength=0x0, packet=[Tag=0x1, Token=0x7974DB57, CommandID=0x70 (0x70), payload=[HmGetTransition], Checksum=0x0, isResponse=false], Checksum=0x0, isResponse: false 
DEBU[2023-12-26T21:03:45+01:00] Request bytes: 30303030303030303030303935343130454338353238424230303041303137393734444235373730303039413336 
DEBU[2023-12-26T21:03:46+01:00] Length of received bytes: 76                 
DEBU[2023-12-26T21:03:46+01:00] Response bytes: 5410EC8528BB00000000000600190100000000F000000000010102020000000000000000107B 
DEBU[2023-12-26T21:03:46+01:00] Received TC: SrcMAC=0x5410EC8528BB, DstMAC=0x000000000006, BodyLength=0x19, packet=[Tag=0x1, Token=0x0, CommandID=0x70 (0xF0), payload=[HmGetTransitionResponse[StateInPercent: 0, DesiredStateInPerced: 0, Error: false, AutoClose: false, DriveTime: 0, Gk: 257, Hcp: HCP[PositionOpen: false, PositionClose: true, OptionRelais: false, LightBarrier: false, Error: false, DrivingToClose: false, Driving: false, HalfOpened: false, ForecastLeadTime: false, Learned: true, NotReferenced: false], Exst: [0 0 0 0 0 0 0 0], Time: 2023-12-26 21:03:46.233468596 +0100 CET m=+0.797944431]], Checksum=0x10, isResponse=true], Checksum=0x7B, isResponse: true 
INFO[2023-12-26T21:03:46+01:00] Transition: HmGetTransitionResponse[StateInPercent: 0, DesiredStateInPerced: 0, Error: false, AutoClose: false, DriveTime: 0, Gk: 257, Hcp: HCP[PositionOpen: false, PositionClose: true, OptionRelais: false, LightBarrier: false, Error: false, DrivingToClose: false, Driving: false, HalfOpened: false, ForecastLeadTime: false, Learned: true, NotReferenced: false], Exst: [0 0 0 0 0 0 0 0], Time: 2023-12-26 21:03:46.233468596 +0100 CET m=+0.797944431]
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