package group

import (
	"bufio"
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/sagernet/sing-box/adapter"
	"github.com/sagernet/sing-box/adapter/outbound"
	"github.com/sagernet/sing-box/common/interrupt"
	"github.com/sagernet/sing-box/common/urltest"
	"github.com/sagernet/sing-box/log"
	E "github.com/sagernet/sing/common/exceptions"
	M "github.com/sagernet/sing/common/metadata"
	N "github.com/sagernet/sing/common/network"

	"github.com/stretchr/testify/require"
)

type fakeOutbound struct {
	outbound.Adapter
	unavailable bool
}

func newFakeOutbound(tag string, network ...string) *fakeOutbound {
	if len(network) == 0 {
		network = []string{N.NetworkTCP, N.NetworkUDP}
	}
	return &fakeOutbound{Adapter: outbound.NewAdapter("fake", tag, network, nil)}
}

func (f *fakeOutbound) DialContext(ctx context.Context, network string, destination M.Socksaddr) (net.Conn, error) {
	if f.unavailable {
		return nil, E.New("unavailable")
	}
	clientConn, serverConn := net.Pipe()
	go func() {
		defer serverConn.Close()
		_, err := http.ReadRequest(bufio.NewReader(serverConn))
		if err != nil {
			return
		}
		serverConn.Write([]byte("HTTP/1.1 204 No Content\r\n\r\n"))
	}()
	return clientConn, nil
}

func (f *fakeOutbound) ListenPacket(ctx context.Context, destination M.Socksaddr) (net.PacketConn, error) {
	return nil, E.New("not implemented")
}

func newPinTestURLTest(outbounds ...adapter.Outbound) *URLTest {
	logger := log.NewNOPFactory().Logger()
	return &URLTest{
		logger: logger,
		group: &URLTestGroup{
			ctx:            context.Background(),
			logger:         logger,
			outbounds:      outbounds,
			link:           "http://www.example.com/generate_204",
			interval:       time.Minute,
			tolerance:      50,
			history:        urltest.NewHistoryStorage(),
			interruptGroup: interrupt.NewGroup(),
		},
	}
}

func storeDelay(history *urltest.HistoryStorage, tag string, delay uint16) {
	history.StoreURLTestHistory(tag, &adapter.URLTestHistory{Time: time.Now(), Delay: delay})
}

func TestURLTestPinnedSelect(t *testing.T) {
	t.Parallel()
	slow := newFakeOutbound("slow")
	fast := newFakeOutbound("fast")
	tcpOnly := newFakeOutbound("tcp-only", N.NetworkTCP)
	urlTest := newPinTestURLTest(slow, fast, tcpOnly)
	group := urlTest.group
	storeDelay(group.history, "slow", 500)
	storeDelay(group.history, "fast", 10)
	storeDelay(group.history, "tcp-only", 300)

	selected, _ := group.Select(N.NetworkTCP)
	require.Equal(t, "fast", selected.Tag())

	group.pinned.Store(slow)
	selected, _ = group.Select(N.NetworkTCP)
	require.Equal(t, "slow", selected.Tag())

	group.history.DeleteURLTestHistory("slow")
	selected, _ = group.Select(N.NetworkTCP)
	require.Equal(t, "fast", selected.Tag(), "unavailable pin must fall back to the fastest outbound")

	group.pinned.Store(tcpOnly)
	selected, _ = group.Select(N.NetworkTCP)
	require.Equal(t, "tcp-only", selected.Tag())
	selected, _ = group.Select(N.NetworkUDP)
	require.Equal(t, "fast", selected.Tag(), "pin must not apply to unsupported networks")
}

func TestURLTestPin(t *testing.T) {
	t.Parallel()
	pinned := newFakeOutbound("pinned")
	fast := newFakeOutbound("fast")
	urlTest := newPinTestURLTest(pinned, fast)
	group := urlTest.group
	storeDelay(group.history, "fast", 10)
	group.performUpdateCheck()
	require.Equal(t, "fast", urlTest.Selected(N.NetworkTCP).Tag())

	require.NoError(t, urlTest.Pin(context.Background(), "pinned"))
	require.NotNil(t, group.history.LoadURLTestHistory("pinned"), "untested pin must be tested first")
	require.Equal(t, pinned, urlTest.Pinned())
	require.Equal(t, "pinned", urlTest.Selected(N.NetworkTCP).Tag())

	urlTest.Unpin()
	require.Nil(t, urlTest.Pinned())
}

func TestURLTestPinRefused(t *testing.T) {
	t.Parallel()
	unavailable := newFakeOutbound("unavailable")
	unavailable.unavailable = true
	urlTest := newPinTestURLTest(unavailable, newFakeOutbound("fast"))

	require.Error(t, urlTest.Pin(context.Background(), "missing"))
	require.Error(t, urlTest.Pin(context.Background(), "unavailable"))
	require.Nil(t, urlTest.Pinned())
}
