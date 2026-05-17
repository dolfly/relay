# github.com/go-gost/relay

Zero-dependency Go wire-protocol library implementing a relay feature-based tunneling protocol.

## Module

| Field | Value |
|-------|-------|
| Module path | `github.com/go-gost/relay` |
| Go version | `1.16` |
| Dependencies | None (stdlib only) |
| Package name | `relay` |
| License | MIT |

## Architecture

A simple Request/Response protocol with a 4-byte fixed header and variable-length features.

### Protocol

```
Request:   | VER(1) | CMD/FLAGS(1) | FEALEN(2) | FEATURES(VAR) |
Response:  | VER(1) | STATUS(1)    | FEALEN(2) | FEATURES(VAR) |
```

Version is always `0x01`.

### Commands

| Constant | Value | Purpose |
|----------|-------|---------|
| `CmdConnect` | `0x01` | TCP-like forward |
| `CmdBind` | `0x02` | Bind listener |
| `CmdAssociate` | `0x03` | UDP associate |

`FUDP` (`0x80`) is a deprecated flag — superseded by `NetworkFeature`.

### Features

Each feature has a 3-byte header: `TYPE(1) + LEN(2) + DATA(VAR)`.

| Type constant | Value | Struct | Wire payload |
|--------------|-------|--------|-------------|
| `FeatureUserAuth` | `0x01` | `UserAuthFeature` | `ULEN(1) + UNAME(VAR) + PLEN(1) + PASSWD(VAR)` |
| `FeatureAddr` | `0x02` | `AddrFeature` | `ATYP(1) + ADDR(VAR) + PORT(2)` |
| `FeatureTunnel` | `0x03` | `TunnelFeature` | `[20]byte` tunnel/connector ID |
| `FeatureNetwork` | `0x04` | `NetworkFeature` | `uint16` network ID |

### TunnelID / ConnectorID (20-byte layout)

```
| ID(16) | FLAG(1) | RSV(2) | WEIGHT(1) |
```

- `TunnelID`: flag `0x80` = private. Methods: `IsPrivate`, `SetPrivate`.
- `ConnectorID`: flag `0x01` = UDP. Methods: `IsUDP`, `SetUDP`.
- Both share `SetWeight` / `Weight` on byte 19.
- `String()` returns UUID-format hex of the 16-byte ID.

### Network IDs

| Constant | Value | String |
|----------|-------|--------|
| `NetworkTCP` | `0x00` | `tcp` |
| `NetworkUDP` | `0x01` | `udp` |
| `NetworkIP` | `0x02` | `ip` |
| `NetworkUnix` | `0x10` | `unix` |
| `NetworkSerial` | `0x11` | `serial` |

## Files

| File | Lines | Purpose |
|------|-------|---------|
| `relay.go` | 226 | `Request`, `Response`, `readFeatures` helper |
| `feature.go` | 536 | `Feature` interface, all 4 feature types, `TunnelID`, `ConnectorID`, `NetworkID` |
| `go.mod` | 3 | Module definition |
| `LICENSE` | 21 | MIT |

## Usage by go-gost

This library is imported by `github.com/go-gost/x` for relay tunnel implementations. The `x/` module imports `core/` (interfaces) and `relay/` (wire protocol) but `relay/` itself depends on neither.

## Building

```sh
cd relay && go build ./...     # library-only, compiles to nothing
cd relay && go vet ./...
```
