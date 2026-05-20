# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Test

```bash
go build ./...           # verify compilation (produces no binary — library-only)
go vet ./...             # static analysis
go test ./...            # run all tests
go test -v -run TestName ./...  # run a single test
go test -cover ./...     # with coverage (currently 88.2%)
```

There are intentionally zero external dependencies — stdlib only.

## Architecture

This is a zero-dependency Go wire-protocol library that implements the GOST relay protocol (v1). It is imported by `github.com/go-gost/x` for relay tunnel implementations but depends on nothing outside the standard library.

Two files define the entire public API:

| File | Purpose |
|------|---------|
| [relay.go](relay.go) | `Request` and `Response` types with `ReadFrom`/`WriteTo` methods, plus the `readFeatures` helper |
| [feature.go](feature.go) | `Feature` interface, all 4 feature types, `TunnelID`, `ConnectorID`, `NetworkID`, and their wire encoding |

### Protocol (4-byte fixed header)

```
Request:   | VER(1) | CMD/FLAGS(1) | FEALEN(2) | FEATURES(VAR) |
Response:  | VER(1) | STATUS(1)    | FEALEN(2) | FEATURES(VAR) |
```

Version is always `0x01`. Commands are `CmdConnect` (0x01), `CmdBind` (0x02), `CmdAssociate` (0x03). `FUDP` (0x80) is a deprecated flag — superseded by `NetworkFeature`.

### Features

Each feature has a 3-byte header: `TYPE(1) + LEN(2) + DATA(VAR)`.

| Type | Struct | Payload |
|------|--------|---------|
| `FeatureUserAuth` (0x01) | `UserAuthFeature` | `ULEN(1) + UNAME + PLEN(1) + PASSWD` |
| `FeatureAddr` (0x02) | `AddrFeature` | `ATYP(1) + ADDR(VAR) + PORT(2)` |
| `FeatureTunnel` (0x03) | `TunnelFeature` | 20-byte tunnel/connector ID |
| `FeatureNetwork` (0x04) | `NetworkFeature` | 2-byte network ID |

### TunnelID / ConnectorID (20-byte value types)

Both are `[20]byte` arrays with identical layout: `ID(16) + FLAG(1) + RSV(2) + WEIGHT(1)`.

- **`TunnelID`**: flag `0x80` = private. Methods: `IsPrivate()`, `SetPrivate(bool)`.
- **`ConnectorID`**: flag `0x01` = UDP. Methods: `IsUDP()`, `SetUDP(bool)`.
- Both share `SetWeight(uint8)` / `Weight() uint8` on byte 19.
- `String()` returns UUID-format hex of the 16-byte ID.
- `Equal()` compares only the 16-byte ID portion (ignores flags/rsv/weight).
- **Value semantics**: all setter methods return a new value — they do not mutate the receiver. `tid.SetPrivate(true)` returns a private copy; the original is unchanged.

### Network IDs

`NetworkTCP` (0x00), `NetworkUDP` (0x01), `NetworkIP` (0x02), `NetworkUnix` (0x10), `NetworkSerial` (0x11). Unknown IDs default to "tcp" in `String()`.

### Key patterns

- `Request` and `Response` implement both `io.ReaderFrom` and `io.WriterTo` for streaming wire I/O.
- `Feature` is an interface with `Type()`, `Encode()`, `Decode([]byte)`. `NewFeature(t FeatureType, data []byte)` is the factory — unknown types return an error.
- All encode/decode uses `encoding/binary.BigEndian`.

### Error sentinels

`ErrBadVersion`, `ErrShortBuffer`, `ErrBadAddrType`, `ErrBadTunnelID`, `ErrBadConnectorID`, `ErrBadNetworkID`.

## Relationship to the broader GOST project

This module is part of the [go-gost](https://github.com/go-gost) workspace. The root [CLAUDE.md](../CLAUDE.md) covers the full project architecture. The `relay/` module sits at the protocol layer — it's imported by `x/` (implementations) but imports nothing from the rest of the workspace.
