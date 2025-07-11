package pkg

import (
	"github.com/go-ini/ini"
	"github.com/pufferffish/wireproxy"

	"encoding/base64"
	"encoding/hex"
	"errors"
	"net"
	"os"
	"strings"

	"net/netip"
)

func parseString(section *ini.Section, keyName string) (string, error) {
	key := section.Key(strings.ToLower(keyName))
	if key == nil {
		return "", errors.New(keyName + " should not be empty")
	}
	value := key.String()
	if strings.HasPrefix(value, "$") {
		if strings.HasPrefix(value, "$$") {
			return strings.Replace(value, "$$", "$", 1), nil
		}
		var ok bool
		value, ok = os.LookupEnv(strings.TrimPrefix(value, "$"))
		if !ok {
			return "", errors.New(keyName + " references unset environment variable " + key.String())
		}
		return value, nil
	}
	return key.String(), nil
}

func parsePort(section *ini.Section, keyName string) (int, error) {
	key := section.Key(keyName)
	if key == nil {
		return 0, errors.New(keyName + " should not be empty")
	}

	port, err := key.Int()
	if err != nil {
		return 0, err
	}

	if port < 0 || port >= 65536 {
		return 0, errors.New("port should be >= 0 and < 65536")
	}

	return port, nil
}

func parseTCPAddr(section *ini.Section, keyName string) (*net.TCPAddr, error) {
	addrStr, err := parseString(section, keyName)
	if err != nil {
		return nil, err
	}
	return net.ResolveTCPAddr("tcp", addrStr)
}

func parseBase64KeyToHex(section *ini.Section, keyName string) (string, error) {
	key, err := parseString(section, keyName)
	if err != nil {
		return "", err
	}
	result, err := encodeBase64ToHex(key)
	if err != nil {
		return result, err
	}

	return result, nil
}

func encodeBase64ToHex(key string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", errors.New("invalid base64 string: " + key)
	}
	if len(decoded) != 32 {
		return "", errors.New("key should be 32 bytes: " + key)
	}
	return hex.EncodeToString(decoded), nil
}

func parseNetIP(section *ini.Section, keyName string) ([]netip.Addr, error) {
	key, err := parseString(section, keyName)
	if err != nil {
		if strings.Contains(err.Error(), "should not be empty") {
			return []netip.Addr{}, nil
		}
		return nil, err
	}

	keys := strings.Split(key, ",")
	var ips = make([]netip.Addr, 0, len(keys))
	for _, str := range keys {
		str = strings.TrimSpace(str)
		if len(str) == 0 {
			continue
		}
		ip, err := netip.ParseAddr(str)
		if err != nil {
			return nil, err
		}
		ips = append(ips, ip)
	}
	return ips, nil
}

func parseCIDRNetIP(section *ini.Section, keyName string) ([]netip.Addr, error) {
	key, err := parseString(section, keyName)
	if err != nil {
		if strings.Contains(err.Error(), "should not be empty") {
			return []netip.Addr{}, nil
		}
		return nil, err
	}

	keys := strings.Split(key, ",")
	var ips = make([]netip.Addr, 0, len(keys))
	for _, str := range keys {
		str = strings.TrimSpace(str)
		if len(str) == 0 {
			continue
		}

		if addr, err := netip.ParseAddr(str); err == nil {
			ips = append(ips, addr)
		} else {
			prefix, err := netip.ParsePrefix(str)
			if err != nil {
				return nil, err
			}

			addr := prefix.Addr()
			ips = append(ips, addr)
		}
	}
	return ips, nil
}

func parseAllowedIPs(section *ini.Section) ([]netip.Prefix, error) {
	key, err := parseString(section, "AllowedIPs")
	if err != nil {
		if strings.Contains(err.Error(), "should not be empty") {
			return []netip.Prefix{}, nil
		}
		return nil, err
	}

	keys := strings.Split(key, ",")
	var ips = make([]netip.Prefix, 0, len(keys))
	for _, str := range keys {
		str = strings.TrimSpace(str)
		if len(str) == 0 {
			continue
		}
		prefix, err := netip.ParsePrefix(str)
		if err != nil {
			return nil, err
		}

		ips = append(ips, prefix)
	}
	return ips, nil
}

func resolveIP(ip string) (*net.IPAddr, error) {
	return net.ResolveIPAddr("ip", ip)
}

func resolveIPPAndPort(addr string) (string, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return "", err
	}

	ip, err := resolveIP(host)
	if err != nil {
		return "", err
	}
	return net.JoinHostPort(ip.String(), port), nil
}

// ParseInterface parses the [Interface] section and extract the information into `device`
func ParseInterface(cfg *ini.File, device *wireproxy.DeviceConfig) error {
	sections, err := cfg.SectionsByName("Interface")
	if len(sections) != 1 || err != nil {
		return errors.New("one and only one [Interface] is expected")
	}
	section := sections[0]

	address, err := parseCIDRNetIP(section, "Address")
	if err != nil {
		return err
	}

	device.Endpoint = address

	privKey, err := parseBase64KeyToHex(section, "PrivateKey")
	if err != nil {
		return err
	}
	device.SecretKey = privKey

	dns, err := parseNetIP(section, "DNS")
	if err != nil {
		return err
	}
	device.DNS = dns

	if sectionKey, err := section.GetKey("MTU"); err == nil {
		value, err := sectionKey.Int()
		if err != nil {
			return err
		}
		device.MTU = value
	}

	if sectionKey, err := section.GetKey("ListenPort"); err == nil {
		value, err := sectionKey.Int()
		if err != nil {
			return err
		}
		device.ListenPort = &value
	}

	checkAlive, err := parseNetIP(section, "CheckAlive")
	if err != nil {
		return err
	}
	device.CheckAlive = checkAlive

	device.CheckAliveInterval = 5
	if sectionKey, err := section.GetKey("CheckAliveInterval"); err == nil {
		value, err := sectionKey.Int()
		if err != nil {
			return err
		}
		if len(checkAlive) == 0 {
			return errors.New("CheckAliveInterval is only valid when CheckAlive is set")
		}

		device.CheckAliveInterval = value
	}

	return nil
}

// ParsePeers parses the [Peer] section and extract the information into `peers`
func ParsePeers(cfg *ini.File, peers *[]wireproxy.PeerConfig) error {
	sections, err := cfg.SectionsByName("Peer")
	if len(sections) < 1 || err != nil {
		return errors.New("at least one [Peer] is expected")
	}

	for _, section := range sections {
		peer := wireproxy.PeerConfig{
			PreSharedKey: "0000000000000000000000000000000000000000000000000000000000000000",
			KeepAlive:    0,
		}

		decoded, err := parseBase64KeyToHex(section, "PublicKey")
		if err != nil {
			return err
		}
		peer.PublicKey = decoded

		if sectionKey, err := section.GetKey("PreSharedKey"); err == nil {
			value, err := encodeBase64ToHex(sectionKey.String())
			if err != nil {
				return err
			}
			peer.PreSharedKey = value
		}

		if sectionKey, err := section.GetKey("Endpoint"); err == nil {
			value := sectionKey.String()
			decoded, err = resolveIPPAndPort(strings.ToLower(value))
			if err != nil {
				return err
			}
			peer.Endpoint = &decoded
		}

		if sectionKey, err := section.GetKey("PersistentKeepalive"); err == nil {
			value, err := sectionKey.Int()
			if err != nil {
				return err
			}
			peer.KeepAlive = value
		}

		peer.AllowedIPs, err = parseAllowedIPs(section)
		if err != nil {
			return err
		}

		*peers = append(*peers, peer)
	}
	return nil
}

func parseTCPClientTunnelConfig(section *ini.Section) (wireproxy.RoutineSpawner, error) {
	config := &wireproxy.TCPClientTunnelConfig{}
	tcpAddr, err := parseTCPAddr(section, "BindAddress")
	if err != nil {
		return nil, err
	}
	config.BindAddress = tcpAddr

	targetSection, err := parseString(section, "Target")
	if err != nil {
		return nil, err
	}
	config.Target = targetSection

	return config, nil
}

func parseSTDIOTunnelConfig(section *ini.Section) (wireproxy.RoutineSpawner, error) {
	config := &wireproxy.STDIOTunnelConfig{}
	targetSection, err := parseString(section, "Target")
	if err != nil {
		return nil, err
	}
	config.Target = targetSection

	return config, nil
}

func parseTCPServerTunnelConfig(section *ini.Section) (wireproxy.RoutineSpawner, error) {
	config := &wireproxy.TCPServerTunnelConfig{}

	listenPort, err := parsePort(section, "ListenPort")
	if err != nil {
		return nil, err
	}
	config.ListenPort = listenPort

	target, err := parseString(section, "Target")
	if err != nil {
		return nil, err
	}
	config.Target = target

	return config, nil
}

func parseSocks5Config(section *ini.Section) (wireproxy.RoutineSpawner, error) {
	config := &wireproxy.Socks5Config{}

	bindAddress, err := parseString(section, "BindAddress")
	if err != nil {
		return nil, err
	}
	config.BindAddress = bindAddress

	username, _ := parseString(section, "Username")
	config.Username = username

	password, _ := parseString(section, "Password")
	config.Password = password

	return config, nil
}

func parseHTTPConfig(section *ini.Section) (wireproxy.RoutineSpawner, error) {
	config := &wireproxy.HTTPConfig{}

	bindAddress, err := parseString(section, "BindAddress")
	if err != nil {
		return nil, err
	}
	config.BindAddress = bindAddress

	username, _ := parseString(section, "Username")
	config.Username = username

	password, _ := parseString(section, "Password")
	config.Password = password

	return config, nil
}

// Takes a function that parses an individual section into a config, and apply it on all
// specified sections
func parseRoutinesConfig(routines *[]wireproxy.RoutineSpawner, cfg *ini.File, sectionName string, f func(*ini.Section) (wireproxy.RoutineSpawner, error)) error {
	sections, err := cfg.SectionsByName(sectionName)
	if err != nil {
		return nil
	}

	for _, section := range sections {
		config, err := f(section)
		if err != nil {
			return err
		}

		*routines = append(*routines, config)
	}

	return nil
}

// ParseConfig takes the path of a configuration file and parses it into Configuration
func ParseConfig(path string) (*wireproxy.Configuration, error) {
	iniOpt := ini.LoadOptions{
		Insensitive:            true,
		AllowShadows:           true,
		AllowNonUniqueSections: true,
	}

	cfg, err := ini.LoadSources(iniOpt, path)
	if err != nil {
		return nil, err
	}

	device := &wireproxy.DeviceConfig{
		MTU: 1420,
	}

	root := cfg.Section("")
	wgConf, err := root.GetKey("WGConfig")
	wgCfg := cfg
	if err == nil {
		wgCfg, err = ini.LoadSources(iniOpt, wgConf.String())
		if err != nil {
			return nil, err
		}
	}

	err = ParseInterface(wgCfg, device)
	if err != nil {
		return nil, err
	}

	err = ParsePeers(wgCfg, &device.Peers)
	if err != nil {
		return nil, err
	}

	var routinesSpawners []wireproxy.RoutineSpawner

	err = parseRoutinesConfig(&routinesSpawners, cfg, "TCPClientTunnel", parseTCPClientTunnelConfig)
	if err != nil {
		return nil, err
	}

	err = parseRoutinesConfig(&routinesSpawners, cfg, "STDIOTunnel", parseSTDIOTunnelConfig)
	if err != nil {
		return nil, err
	}

	err = parseRoutinesConfig(&routinesSpawners, cfg, "TCPServerTunnel", parseTCPServerTunnelConfig)
	if err != nil {
		return nil, err
	}

	err = parseRoutinesConfig(&routinesSpawners, cfg, "Socks5", parseSocks5Config)
	if err != nil {
		return nil, err
	}

	err = parseRoutinesConfig(&routinesSpawners, cfg, "http", parseHTTPConfig)
	if err != nil {
		return nil, err
	}

	return &wireproxy.Configuration{
		Device:   device,
		Routines: routinesSpawners,
	}, nil
}
