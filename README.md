**ABSOLUTELY NO WARRANTY! USE AT YOUR OWN RISK!**

# Hörmann BiSecur Gateway Protocol GoLang SDK

Goal is to create a fully fledged replacement client after Hörmann stated end of life of their cloud and android application.

If all goes fine, later this repository will provide you both a GoLang SDK and a CLI client you can use to
- open and close your door
- manage users
- etc.

## Usage
```
Application to manage your Hörmann BiSecur gateway without the central cloud directly on your LAN.

Usage:
  halsecur [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  discover    Discover Hörmann BiSecur gateways on the local network
  get-name    Queries the name of the Hörmann BiSecur gateway
  groups      Manages users defined in your Hörmann BiSecur gateway.
  help        Help about any command
  login       
  logout      
  ping        Check if your Hörmann BiSecur gateway is reachable or not.
  set-state   Open or close a door connected to your Hörmann BiSecur gateway.
  status      Queries the status (open/closed/etc) of your door.
  users       Manages users defined in your Hörmann BiSecur gateway.

Flags:
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