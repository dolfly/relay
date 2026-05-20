package relay

import (
	"bytes"
	"io"
	"testing"
)

// Request / Response round-trip benchmarks

func BenchmarkRequestReadFrom(b *testing.B) {
	req := Request{
		Version: Version1,
		Cmd:     CmdConnect,
		Features: []Feature{
			&UserAuthFeature{Username: "admin", Password: "secret"},
			&AddrFeature{AType: AddrDomain, Host: "example.com", Port: 8080},
			&TunnelFeature{ID: NewTunnelID(bytes.Repeat([]byte{0xDE}, 16))},
			&NetworkFeature{Network: NetworkTCP},
		},
	}
	var buf bytes.Buffer
	if _, err := req.WriteTo(&buf); err != nil {
		b.Fatal(err)
	}
	data := buf.Bytes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var decoded Request
		decoded.ReadFrom(bytes.NewReader(data))
	}
}

func BenchmarkRequestWriteTo(b *testing.B) {
	req := Request{
		Version: Version1,
		Cmd:     CmdConnect,
		Features: []Feature{
			&UserAuthFeature{Username: "admin", Password: "secret"},
			&AddrFeature{AType: AddrDomain, Host: "example.com", Port: 8080},
			&TunnelFeature{ID: NewTunnelID(bytes.Repeat([]byte{0xDE}, 16))},
			&NetworkFeature{Network: NetworkTCP},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req.WriteTo(io.Discard)
	}
}

func BenchmarkRequestRoundTrip(b *testing.B) {
	req := Request{
		Version: Version1,
		Cmd:     CmdConnect,
		Features: []Feature{
			&UserAuthFeature{Username: "admin", Password: "secret"},
			&AddrFeature{AType: AddrDomain, Host: "example.com", Port: 8080},
			&TunnelFeature{ID: NewTunnelID(bytes.Repeat([]byte{0xDE}, 16))},
			&NetworkFeature{Network: NetworkTCP},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		req.WriteTo(&buf)
		var decoded Request
		decoded.ReadFrom(&buf)
	}
}

func BenchmarkResponseReadFrom(b *testing.B) {
	resp := Response{
		Version: Version1,
		Status:  StatusOK,
		Features: []Feature{
			&TunnelFeature{ID: NewConnectorID(bytes.Repeat([]byte{0xCD}, 16))},
			&NetworkFeature{Network: NetworkUDP},
		},
	}
	var buf bytes.Buffer
	if _, err := resp.WriteTo(&buf); err != nil {
		b.Fatal(err)
	}
	data := buf.Bytes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var decoded Response
		decoded.ReadFrom(bytes.NewReader(data))
	}
}

func BenchmarkResponseWriteTo(b *testing.B) {
	resp := Response{
		Version: Version1,
		Status:  StatusOK,
		Features: []Feature{
			&TunnelFeature{ID: NewConnectorID(bytes.Repeat([]byte{0xCD}, 16))},
			&NetworkFeature{Network: NetworkUDP},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp.WriteTo(io.Discard)
	}
}

// Feature encode/decode benchmarks

func BenchmarkUserAuthFeatureEncode(b *testing.B) {
	f := &UserAuthFeature{Username: "admin", Password: "secret"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Encode()
	}
}

func BenchmarkUserAuthFeatureDecode(b *testing.B) {
	f := &UserAuthFeature{Username: "admin", Password: "secret"}
	buf, _ := f.Encode()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var decoded UserAuthFeature
		decoded.Decode(buf)
	}
}

func BenchmarkAddrFeatureEncodeIPv4(b *testing.B) {
	f := &AddrFeature{AType: AddrIPv4, Host: "192.168.1.1", Port: 8080}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Encode()
	}
}

func BenchmarkAddrFeatureDecodeIPv4(b *testing.B) {
	f := &AddrFeature{AType: AddrIPv4, Host: "192.168.1.1", Port: 8080}
	buf, _ := f.Encode()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var decoded AddrFeature
		decoded.Decode(buf)
	}
}

func BenchmarkAddrFeatureEncodeDomain(b *testing.B) {
	f := &AddrFeature{AType: AddrDomain, Host: "example.com", Port: 443}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Encode()
	}
}

func BenchmarkAddrFeatureDecodeDomain(b *testing.B) {
	f := &AddrFeature{AType: AddrDomain, Host: "example.com", Port: 443}
	buf, _ := f.Encode()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var decoded AddrFeature
		decoded.Decode(buf)
	}
}

func BenchmarkAddrFeatureEncodeIPv6(b *testing.B) {
	f := &AddrFeature{AType: AddrIPv6, Host: "::1", Port: 443}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Encode()
	}
}

func BenchmarkAddrFeatureDecodeIPv6(b *testing.B) {
	f := &AddrFeature{AType: AddrIPv6, Host: "::1", Port: 443}
	buf, _ := f.Encode()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var decoded AddrFeature
		decoded.Decode(buf)
	}
}

func BenchmarkAddrFeatureParseFrom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var f AddrFeature
		f.ParseFrom("example.com:8080")
	}
}

func BenchmarkTunnelFeatureEncode(b *testing.B) {
	f := &TunnelFeature{ID: NewTunnelID(bytes.Repeat([]byte{0xDE}, 16))}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Encode()
	}
}

func BenchmarkTunnelFeatureDecode(b *testing.B) {
	f := &TunnelFeature{ID: NewTunnelID(bytes.Repeat([]byte{0xDE}, 16))}
	buf, _ := f.Encode()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var decoded TunnelFeature
		decoded.Decode(buf)
	}
}

func BenchmarkNetworkFeatureEncode(b *testing.B) {
	f := &NetworkFeature{Network: NetworkTCP}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Encode()
	}
}

func BenchmarkNetworkFeatureDecode(b *testing.B) {
	f := &NetworkFeature{Network: NetworkTCP}
	buf, _ := f.Encode()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var decoded NetworkFeature
		decoded.Decode(buf)
	}
}

// Feature factory / reader benchmarks

func BenchmarkNewFeature(b *testing.B) {
	data := []byte{0x04, 'u', 's', 'e', 'r', 0x04, 'p', 'a', 's', 's'}
	for i := 0; i < b.N; i++ {
		NewFeature(FeatureUserAuth, data)
	}
}

func BenchmarkReadFeature(b *testing.B) {
	f := &UserAuthFeature{Username: "admin", Password: "secret"}
	fb, _ := f.Encode()

	header := make([]byte, featureHeaderLen+len(fb))
	header[0] = byte(FeatureUserAuth)
	header[1] = byte(len(fb) >> 8)
	header[2] = byte(len(fb))
	copy(header[3:], fb)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ReadFeature(bytes.NewReader(header))
	}
}

// TunnelID / ConnectorID benchmarks

func BenchmarkTunnelIDString(b *testing.B) {
	tid := NewTunnelID(bytes.Repeat([]byte{0xAB}, 16))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tid.String()
	}
}

func BenchmarkTunnelIDEqual(b *testing.B) {
	a := NewTunnelID(bytes.Repeat([]byte{0xAB}, 16))
	c := NewTunnelID(bytes.Repeat([]byte{0xCD}, 16))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Equal(c)
	}
}

func BenchmarkTunnelIDSetPrivate(b *testing.B) {
	tid := NewTunnelID(bytes.Repeat([]byte{0xAB}, 16))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tid.SetPrivate(true)
	}
}

func BenchmarkConnectorIDString(b *testing.B) {
	cid := NewConnectorID(bytes.Repeat([]byte{0xAB}, 16))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cid.String()
	}
}

func BenchmarkConnectorIDEqual(b *testing.B) {
	a := NewConnectorID(bytes.Repeat([]byte{0xAB}, 16))
	c := NewConnectorID(bytes.Repeat([]byte{0xCD}, 16))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Equal(c)
	}
}
