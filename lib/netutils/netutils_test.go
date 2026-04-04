package netutils

import (
	"net"
	"testing"
)

func TestBroadcastAddress(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		cidr string
		want string
	}{
		{
			name: "ipv4 /24",
			cidr: "10.0.0.0/24",
			want: "10.0.0.255",
		},
		{
			name: "ipv4 /30",
			cidr: "192.168.0.0/30",
			want: "192.168.0.3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, ipnet, err := net.ParseCIDR(tt.cidr)
			if err != nil {
				t.Fatal(err)
			}
			got := BroadcastAddress(ipnet)
			if got.String() != tt.want {
				t.Fatalf("BroadcastAddress(%s) = %s, want %s", tt.cidr, got, tt.want)
			}
		})
	}
}

func TestNextSubnet(t *testing.T) {
	t.Parallel()

	t.Run("advances /24", func(t *testing.T) {
		t.Parallel()
		_, a, _ := net.ParseCIDR("10.0.0.0/24")
		next := NextSubnet(a)
		if next == nil {
			t.Fatal("expected non-nil next subnet")
		}
		if next.String() != "10.0.1.0/24" {
			t.Fatalf("got %s, want 10.0.1.0/24", next.String())
		}
	})

	t.Run("overflow returns nil", func(t *testing.T) {
		t.Parallel()
		_, max, _ := net.ParseCIDR("255.255.255.255/32")
		if NextSubnet(max) != nil {
			t.Fatal("expected nil when incrementing past IPv4 end")
		}
	})
}

func TestIncrementIP(t *testing.T) {
	t.Parallel()

	t.Run("nil ip", func(t *testing.T) {
		t.Parallel()
		if IncrementIP(nil, 1) != nil {
			t.Fatal("expected nil")
		}
	})

	t.Run("zero increment normalizes", func(t *testing.T) {
		t.Parallel()
		ip := net.ParseIP("10.0.0.1")
		got := IncrementIP(ip, 0)
		if got == nil {
			t.Fatal("expected non-nil")
		}
		if !got.Equal(ip) {
			t.Fatalf("got %s, want %s", got, ip)
		}
	})

	t.Run("ipv4 plus one", func(t *testing.T) {
		t.Parallel()
		ip := net.ParseIP("10.0.0.1")
		got := IncrementIP(ip, 1)
		if got == nil || got.String() != "10.0.0.2" {
			t.Fatalf("got %v, want 10.0.0.2", got)
		}
	})

	t.Run("ipv4 overflow out of range", func(t *testing.T) {
		t.Parallel()
		ip := net.ParseIP("255.255.255.255")
		got := IncrementIP(ip, 1)
		if got != nil {
			t.Fatalf("expected nil, got %s", got)
		}
	})
}

func TestIncrementIP_IPv6(t *testing.T) {
	t.Parallel()
	ip := net.ParseIP("2001:db8::1")
	got := IncrementIP(ip, 1)
	want := net.ParseIP("2001:db8::2")
	if got == nil || !got.Equal(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestIncrementIP_negativeIPv4(t *testing.T) {
	t.Parallel()
	ip := net.ParseIP("10.0.0.10")
	got := IncrementIP(ip, -3)
	want := net.ParseIP("10.0.0.7")
	if got == nil || !got.Equal(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestIncrementIP_negativeIPv6(t *testing.T) {
	t.Parallel()
	ip := net.ParseIP("2001:db8::10")
	got := IncrementIP(ip, -1)
	want := net.ParseIP("2001:db8::f")
	if got == nil || !got.Equal(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}
