# GOST Relay Protocol

[![Go Reference](https://pkg.go.dev/badge/github.com/go-gost/relay.svg)](https://pkg.go.dev/github.com/go-gost/relay)

Zero-dependency Go implementation of the GOST relay wire protocol (v1). Used by [go-gost/gost](https://github.com/go-gost/gost) for relay tunnel communication.

## Install

```sh
go get github.com/go-gost/relay
```

## Protocol

Fixed 4-byte header, version always `0x01`.

```
Request:   | VER(1) | CMD/FLAGS(1) | FEALEN(2) | FEATURES(VAR) |
Response:  | VER(1) | STATUS(1)    | FEALEN(2) | FEATURES(VAR) |
```

### Commands

| Constant | Value | Purpose |
|----------|-------|---------|
| `CmdConnect` | `0x01` | TCP forward |
| `CmdBind` | `0x02` | Bind listener |
| `CmdAssociate` | `0x03` | UDP associate |

### Features

Each feature has a 3-byte header: `TYPE(1) + LEN(2) + DATA(VAR)`.

| Type | Payload |
|------|---------|
| `FeatureUserAuth` (0x01) | Username/password for authentication |
| `FeatureAddr` (0x02) | Target address (IPv4, IPv6, or domain + port) |
| `FeatureTunnel` (0x03) | 20-byte tunnel/connector ID |
| `FeatureNetwork` (0x04) | 2-byte network type (TCP, UDP, IP, Unix, Serial) |

## Usage

```go
package main

import (
    "bytes"
    "net"

    "github.com/go-gost/relay"
)

func main() {
    // Build a request
    req := &relay.Request{
        Version: relay.Version1,
        Cmd:     relay.CmdConnect,
        Features: []relay.Feature{
            &relay.UserAuthFeature{Username: "user", Password: "pass"},
            &relay.AddrFeature{AType: relay.AddrDomain, Host: "example.com", Port: 443},
            &relay.NetworkFeature{Network: relay.NetworkTCP},
        },
    }

    // Write to a connection
    conn, _ := net.Dial("tcp", "relay-server:8421")
    req.WriteTo(conn)

    // Read the response
    var resp relay.Response
    resp.ReadFrom(conn)
}
```

### AddrFeature from string

```go
var f relay.AddrFeature
f.ParseFrom("example.com:8080") // auto-detects IPv4/IPv6/domain
```

### TunnelID / ConnectorID

Both are 20-byte value types with UUID-format string output:

```go
tid := relay.NewTunnelID(uuidBytes)         // standard
tid := relay.NewPrivateTunnelID(uuidBytes)  // sets private flag (0x80)
tid.String() // "abababab-abab-abab-abab-abababababab"

cid := relay.NewConnectorID(uuidBytes)
cid := relay.NewUDPConnectorID(uuidBytes)   // sets UDP flag (0x01)
```

All setter methods (`SetPrivate`, `SetUDP`, `SetWeight`) return a new value — they do not mutate the receiver.

## License

MIT
