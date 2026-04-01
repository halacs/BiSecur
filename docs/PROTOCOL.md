# BiSecur Gateway Wire Protocol

Wire protocol documentation for the Hörmann BiSecur gateway, based on the [Kotlin SDK](https://github.com/bisdk/sdk), the decompiled BiSecur.swf (ActionScript source via [JPEXS ffdec](https://github.com/niclasr/JPEXS-decompiler)), and testing against a single gateway (firmware EE001425-14, MCP V3.0). Some commands and payload formats have only been read from decompiled source and not verified on hardware. Corrections and additions are welcome.

## Transport Overview

```
TCP port 4000
  └─ TMCP frame (hex-encoded ASCII)
      └─ MCP frame (binary, encoded as hex within TMCP)
          └─ Command + Payload
              └─ JMCP (JSON sub-protocol, via CMD 0x06)
```

The gateway also listens on:
- **UDP 4001**: Discovery requests (broadcast `<Discover target="Gateway" />`)
- **UDP 4002**: Discovery responses (XML with `swVersion`, `hwVersion`, `mac`, `protocol`)

## TMCP Frame Format

TMCP frames are hex-encoded ASCII strings sent over a raw TCP socket. No framing delimiters; the receiver must parse the MCP LENGTH field to determine frame boundaries.

```
[SRC_MAC:12 hex][DST_MAC:12 hex][MCP_DATA:N hex][TMCP_CKSUM:2 hex]
```

| Field | Size | Description |
|-------|------|-------------|
| SRC_MAC | 12 hex chars | Sender MAC address |
| DST_MAC | 12 hex chars | Gateway MAC address |
| MCP_DATA | variable | Hex-encoded MCP binary frame |
| TMCP_CKSUM | 2 hex chars | Checksum of all preceding chars |

**TMCP Checksum:** sum of ASCII character codes of the entire hex string (SRC + DST + MCP_DATA), masked to 8 bits:

```
checksum = sum(ord(c) for c in hex_string) & 0xFF
```

## MCP Frame Format

The MCP frame is binary, transmitted as hex within the TMCP envelope.

```
[LENGTH:2][TAG:1][TOKEN:4][COMMAND:1][PAYLOAD:N][MCP_CKSUM:1]
```

| Field | Size (bytes) | Encoding | Description |
|-------|-------------|----------|-------------|
| LENGTH | 2 | uint16 big-endian | `9 + len(PAYLOAD)` |
| TAG | 1 | uint8 | Sequence number. `0` before auth. Used for request-response matching |
| TOKEN | 4 | uint32 big-endian | Session token from LOGIN. `0` before auth |
| COMMAND | 1 | uint8 | Command ID. Responses have bit 7 set (`cmd | 0x80`) |
| PAYLOAD | variable | raw bytes | Command-specific data |
| MCP_CKSUM | 1 | uint8 | Checksum |

**MCP Checksum:** LENGTH value (as uint16) plus sum of all remaining bytes (TAG through PAYLOAD), masked to 8 bits:

```
length_val = uint16_big_endian(data[0:2])
checksum = (length_val + sum(data[2:])) & 0xFF
```

The checksum byte itself is NOT included in the sum.

**Response identification:** Response command byte has bit 7 set: `response_cmd = request_cmd | 0x80`. Error responses use `CMD_ERROR (0x01)` with the error code in the payload.

## Command Reference

### Authentication

| Cmd | Hex | Name | Request | Response |
|-----|-----|------|---------|----------|
| 16 | 0x10 | LOGIN | `[uname_len:1][uname:N][pass:M]` | `[user_id:1][token:4]` |
| 17 | 0x11 | LOGOUT | empty | empty (no response expected) |
| 35 | 0x23 | CHANGE_PASSWD | `[new_pass_utf8:N]` | empty |
| 69 | 0x45 | CHANGE_PASSWORD_OF_USER | `[user_id:1][new_pass:N]` | empty |

**LOGIN payload:** Username is length-prefixed (`[len:1][username:N]`), password follows with no length prefix. The remaining bytes after the username are the password.

### Device Info

| Cmd | Hex | Name | Request | Response |
|-----|-----|------|---------|----------|
| 0 | 0x00 | PING | empty | empty |
| 2 | 0x02 | GET_MAC | empty | `[mac:6]` |
| 7 | 0x07 | GET_GW_VERSION | empty | `[version_str:N]` (UTF-8) |
| 38 | 0x26 | GET_NAME | empty | `[name_utf8:N]` |

### User Management

| Cmd | Hex | Name | Request | Response |
|-----|-----|------|---------|----------|
| 32 | 0x20 | GET_USER_IDS | empty | `[uid0:1][uid1:1]...` |
| 33 | 0x21 | GET_USER_NAME | `[user_id:1]` | `[name_utf8:N]` |
| 34 | 0x22 | ADD_USER | `[uname_len:1][uname:N][pass:M]` | `[user_id:1]` |
| 36 | 0x24 | REMOVE_USER | `[user_id:1]` | `[user_id:1]` |
| 37 | 0x25 | SET_USER_RIGHTS | `[user_id:1][port_id:1]` | echo |
| 40 | 0x28 | GET_USER_RIGHTS | `[user_id:1]` | `[port0:1][port1:1]...` |

**Security warning:** ADD_USER (0x22) appears to require no authentication based on the decompiled source. If true, anyone on the network could create a gateway user. See also the [SEC Consult analysis](https://sec-consult.com/blog/detail/hoermann-opening-doors-for-everyone/).

### Port Management

Ports represent radio channels. Each port corresponds to a paired device's radio code.

| Cmd | Hex | Name | Request | Response | Notes |
|-----|-----|------|---------|----------|-------|
| 48 | 0x30 | GET_PORTS | empty | `[port0:1]...` | List all port IDs |
| 49 | 0x31 | GET_TYPE | `[port_id:1]` | `[port_id:1][type:1]` | See Port Types |
| 54 | 0x36 | SET_TYPE | `[port_id:1][type:1]` | echo | UI metadata only |
| 41 | 0x29 | ADD_PORT | empty | `[new_port_id:1]` | Gateway LISTENS for radio signal |
| 65 | 0x41 | INHERIT_PORT | empty | `[new_port_id:1]` | Gateway TRANSMITS its radio code |
| 66 | 0x42 | REMOVE_PORT | `[port_id:1]` | `[removed_port_id:1]` | |

**ADD_PORT vs INHERIT_PORT:**
- **ADD_PORT (0x29):** The gateway enters receive mode and listens for a hand remote's radio signal. In my testing, the hand remote needed to be within ~10cm of the gateway for a successful clone. Timeout appears to be ~40 seconds.
- **INHERIT_PORT (0x41):** The gateway transmits its own radio code. The motor must be in learn mode (press P button on the motor). In my testing, the gateway needed to be within ~1 meter of the motor for pairing to succeed, though normal operation works at greater range.

In my testing with a ProMatic 3, INHERIT_PORT ports supported position feedback via HM_GET_TRANSITION, while ADD_PORT ports did not (always returned PORT_ERROR). This may vary with other motor models.

### Group Management

Groups organize ports into logical units (e.g., "Garage" containing a door impulse port).

| Cmd | Hex | Name | Request | Response |
|-----|-----|------|---------|----------|
| 42 | 0x2A | ADD_GROUP | empty | `[group_id:1]` |
| 43 | 0x2B | REMOVE_GROUP | `[group_id:1]` | empty |
| 44 | 0x2C | SET_GROUP_NAME | `[group_id:1][name_utf8:N]` | echo |
| 46 | 0x2E | SET_GROUPED_PORTS | `[group_id:1][port0:1]...` | echo |
| 47 | 0x2F | GET_GROUPED_PORTS | `[group_id:1]` | `[port0:1]...` |

### Value Store

The gateway has 32 addressable value slots (0-31). Addresses 0-15 store group type for group N (see Group Types). Addresses 16-31 store the "requestable port" for group N-16 (the port ID used for HM_GET_TRANSITION feedback).

| Cmd | Hex | Name | Request | Response |
|-----|-----|------|---------|----------|
| 3 | 0x03 | SET_VALUE | `[address:1][value:1]` | echo |
| 4 | 0x04 | GET_VALUE | `[address:1]` | `[address:1][value:1]` |

### Control

| Cmd | Hex | Name | Request | Response | Notes |
|-----|-----|------|---------|----------|-------|
| 51 | 0x33 | SET_STATE | `[port_id:1][state:1]` | varies | `state=0xFF` = impulse |
| 112 | 0x70 | HM_GET_TRANSITION | `[port_id:1]` | 16 bytes | Position feedback |

**SET_STATE notes:**
- `state=0xFF` triggers an impulse (open/close toggle). This is the only value I have observed to work
- Other state values (0x00, 0x06, 0x07) caused timeouts in our testing
- On INHERIT_PORT ports, the response is an unsolicited HM_GET_TRANSITION frame (not a SET_STATE ACK)
- On ADD_PORT ports, returns PORT_ERROR (10) but the door still moves

### HM_GET_TRANSITION Response (16 bytes)

| Offset | Size | Field | Description |
|--------|------|-------|-------------|
| 0 | 1 | position | Door position: raw 0-200 (0=closed, 200=open, divide by 2 for %) |
| 1 | 1 | desired | Target position (raw 0-200) |
| 2 | 1 | flags | Bit 7: Error, Bit 6: AutoClose |
| 3 | 1 | drv_time | Drive time remaining (counts down during motion, seconds) |
| 4-5 | 2 | gk | Device key (uint16 big-endian) |
| 6-7 | 2 | hcp | HCP bitfield (uint16 **little-endian**, see below) |
| 8-15 | 8 | exst | Extended status bytes |

**HCP bitfield (uint16 little-endian):**

| Bit | Flag | Description |
|-----|------|-------------|
| 0 | PositionOpen | Door is at fully open position |
| 1 | PositionClose | Door is at fully closed position |
| 2 | OptionRelais | Relay option active |
| 3 | LightBarrier | Light barrier triggered |
| 4 | Error | Motor error state |
| 5 | DrivingToClose | Motor is moving toward closed |
| 6 | Driving | Motor is currently running |
| 7 | HalfOpened | Door is at partial/venting position |
| 8 | ForecastLeadTime | Lead time forecast |
| 9 | Learned | Port has been learned/paired |
| 10 | NotReferenced | Port not referenced |

**Note:** The HCP bitfield uses **little-endian** byte order, unlike the rest of the MCP frame which is big-endian. The HCP data likely originates from the motor's Hörmann Communication Protocol layer rather than the gateway's MCP framing.

**Observed behavior (tested with ProMatic 3, firmware EE001425-14):**
- Only returned data on INHERIT_PORT (0x41) ports. ADD_PORT (0x29) ports returned PORT_ERROR (10)
- Transient PORT_ERROR responses occurred on roughly 20% of queries, even on supported ports
- Polling faster than every 2 seconds caused the gateway to reset the TCP connection

### WiFi

| Cmd | Hex | Name | Request | Response |
|-----|-----|------|---------|----------|
| 81 | 0x51 | SCAN_WIFI | empty | - |
| 82 | 0x52 | WIFI_FOUND | - | - |
| 83 | 0x53 | GET_WIFI_STATE | empty | `[state:1]` |

WiFi states: 0=CONNECTED, 1=NOT_CONNECTED, 64=BUSY, 128=AP_NOT_FOUND, 129=SECURITY_MISMATCH, 130=AUTHENTICATION_FAILURE, 131=CONNECTION_FAILURE

### JMCP (JSON Sub-Protocol)

Sent via CMD_JMCP (0x06). The payload is a UTF-8 JSON string. Note the inconsistent key casing in the firmware:

```json
{"cmd": "GET_USERS"}      // lowercase key
{"CMD": "GET_GROUPS"}     // uppercase key
{"CMD": "GET_VALUES"}     // uppercase key
```

Both casings work. The `GET_GROUPS` command supports an optional `FORUSER` filter:
```json
{"CMD": "GET_GROUPS", "FORUSER": 1}
```

## Error Codes

| Code | Name | Description |
|------|------|-------------|
| 0 | COMMAND_NOT_FOUND | Unknown command ID |
| 1 | INVALID_PROTOCOL | Protocol version mismatch |
| 2 | LOGIN_FAILED | Bad credentials |
| 3 | INVALID_TOKEN | Session expired or invalid |
| 4 | USER_ALREADY_EXISTS | Username taken |
| 5 | NO_EMPTY_USER_SLOT | Max users reached |
| 6 | INVALID_PASSWORD | Password validation failed |
| 7 | INVALID_USERNAME | Username validation failed |
| 8 | USER_NOT_FOUND | User ID doesn't exist |
| 9 | PORT_NOT_FOUND | Port ID doesn't exist |
| 10 | PORT_ERROR | Motor didn't ACK (normal for non-feedback ports) |
| 11 | GATEWAY_BUSY | Gateway processing another request |
| 12 | PERMISSION_DENIED | User lacks access, or token expired |
| 13 | NO_EMPTY_GROUP_SLOT | Max groups reached |
| 14 | GROUP_NOT_FOUND | Group ID doesn't exist |
| 15 | INVALID_PAYLOAD | Malformed request payload |
| 16 | OUT_OF_RANGE | Value outside allowed range |
| 17 | ADD_PORT_ERROR | ADD_PORT failed (no signal received) |
| 18 | NO_EMPTY_PORT_SLOT | Max ports reached (16 max) |
| 19 | ADAPTER_BUSY | Radio adapter in use |

## Type Enumerations

### Port Types

| Value | Name | Description |
|-------|------|-------------|
| 0 | NONE | Unconfigured |
| 1 | IMPULS | Momentary impulse (open/close toggle) |
| 2 | AUTO_CLOSE | Auto-close after open |
| 3 | ON_OFF | Toggle on/off |
| 4 | UP | Move up only |
| 5 | DOWN | Move down only |
| 6 | HALF | Half-open position |
| 7 | WALK | Pedestrian opening |
| 8 | LIGHT | Light control |
| 9 | ON | Switch on |
| 10 | OFF | Switch off |
| 11 | LOCK | Lock |
| 12 | UNLOCK | Unlock |
| 13 | OPEN_DOOR | Open door |
| 14 | LIFT | Raise |
| 15 | SINK | Lower |

### Group Types (stored in value slots 0-15)

| Value | Name |
|-------|------|
| 1 | SECTIONAL_DOOR |
| 2 | HORIZ_SECTIONAL_DOOR |
| 3 | SWING_GATE_SINGLE |
| 4 | SWING_GATE_DOUBLE |
| 5 | SLIDING_GATE |
| 6 | DOOR |
| 7 | LIGHT |
| 8 | OTHER |
| 9 | JACK |
| 10 | SMARTKEY |
| 15 | BARRIER |

## Session Management

### Token Lifetime

Session tokens expire after approximately 30 minutes. After expiration, commands return PERMISSION_DENIED (12). Re-login with a fresh connection to obtain a new token.

### Stale Session Drain

The gateway has limited session slots. If a connection is closed without sending LOGOUT, the session persists for ~30 minutes. On the next LOGIN, the gateway sends a LOGOUT frame (cmd 0x11, tag 0xFF) for each stale session before the actual LOGIN response.

The receiver must buffer incoming data and parse by MCP LENGTH to extract individual frames. Multiple LOGOUT frames may arrive before (or coalesced with) the LOGIN response.

**Always send LOGOUT before closing the TCP connection** to prevent stale session accumulation. With enough stale sessions, the flood of LOGOUTs can overwhelm the gateway and cause a TCP reset.

### Rate Limiting

The gateway drops connections if commands are sent too rapidly:
- Minimum ~100ms delay between any commands
- HM_GET_TRANSITION: minimum ~2 seconds between polls
- GET_VALUE: minimum ~100ms between calls

## Frame Construction Example

LOGIN frame for user `admin` with password `Password123!`:

```
Payload bytes: 05 61 64 6D 69 6E 50 61 73 73 77 6F 72 64 31 32 33 21
               ^^ ^^^^^^^^^^^^^^^^^ ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
               |  "admin" (5 bytes)  "Password123!" (12 bytes)
               username length

MCP frame (before checksum):
  00 1B  00  00000000  10  056164...3321
  ^^^^   ^^  ^^^^^^^^  ^^  ^^^^^^^^^^^^
  LEN=27 TAG TOKEN=0   CMD PAYLOAD (LOGIN)

TMCP frame:
  FFFFFFFFFFFF 801F12141DA5 001B00000000001005616...332143 DE
  ^^^^^^^^^^^^ ^^^^^^^^^^^^ ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ ^^
  SRC_MAC      DST_MAC      MCP_HEX (incl MCP checksum 0x43)  TMCP_CKSUM
```
