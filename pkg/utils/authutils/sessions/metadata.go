package sessions

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/ua-parser/uap-go/uaparser"
)

func getSessionMetadata(r *http.Request, userAgent string) (*SessionMetadata, error) {
	meta := &SessionMetadata{IP: getIP(r, true)}

	err := meta.parseLocation()
	if err != nil {
		return nil, err
	}
	meta.parseDevice(userAgent)

	return meta, nil
}

func getIP(r *http.Request, isDev bool) string {
	if isDev {
		return "109.207.190.162"
	}

	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}

	if strings.Contains(ip, ",") {
		parts := strings.Split(ip, ",")
		ip = strings.TrimSpace(parts[0])
	}
	return ip
}

func (m *SessionMetadata) parseLocation() error {
	info, err := ipinfo.GetIPInfo(net.IP(m.IP).To4())
	if err != nil {
		return err
	}

	m.Location.City = info.City

	if info.Country != "" {
		m.Location.Country = strings.ToLower(info.Country)
	} else {
		m.Location.Country = "Неизвестно"
	}

	if len(info.Location) > 0 {
		coords := strings.Split(info.Location, ",")
		if len(coords) == 2 {
			m.Location.Lat = parseFloat(coords[0])
			m.Location.Lng = parseFloat(coords[1])
		}
	}
	return nil
}

func (m *SessionMetadata) parseDevice(userAgent string) error {
	parser, err := uaparser.New()
	if err != nil {
		return err
	}

	client := parser.Parse(userAgent)

	m.Device.Browser = client.UserAgent.Family
	m.Device.OS = client.Os.Family
	m.Device.Type = client.Device.Family
	m.Device.Model = client.Device.Model

	return nil
}

func parseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}
