import (
	"fmt"
	"net"
	"sync"
	"time"
)

// GetLocalIP returns the non-loopback local IP address of the host.
func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no local network IP address found")
}

// GetLocalSubnet returns the local IP and subnet mask.
func GetLocalSubnet() (*net.IPNet, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet, nil
			}
		}
	}

	return nil, fmt.Errorf("no local subnet found")
}

// ScanNetwork scans the local network for active hosts.
// It uses concurrent TCP connection attempts to common ports.
func ScanNetwork(timeout time.Duration) ([]string, error) {
	subnet, err := GetLocalSubnet()
	if err != nil {
		return nil, err
	}

	var activeHosts []string
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Generate all IPs in the subnet
	for ip := subnet.IP.Mask(subnet.Mask); subnet.Contains(ip); incIP(ip) {
		ipStr := ip.String()
		wg.Add(1)

		go func(host string) {
			defer wg.Done()
			if isHostActive(host, timeout) {
				mu.Lock()
				activeHosts = append(activeHosts, host)
				mu.Unlock()
			}
		}(ipStr)
	}

	wg.Wait()
	return activeHosts, nil
}

// isHostActive checks if a host is active by attempting connections to common ports.
func isHostActive(ip string, timeout time.Duration) bool {
	// Try common ports: HTTP, HTTPS, SSH, SMB
	ports := []string{"80", "443", "22", "445"}

	for _, port := range ports {
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, port), timeout)
		if err == nil {
			conn.Close()
			return true
		}
	}
	return false
}

// PingScan uses ICMP-style detection (connection-based for cross-platform compatibility).
func PingScan(timeout time.Duration) ([]string, error) {
	subnet, err := GetLocalSubnet()
	if err != nil {
		return nil, err
	}

	var activeHosts []string
	var mu sync.Mutex
	var wg sync.WaitGroup

	for ip := subnet.IP.Mask(subnet.Mask); subnet.Contains(ip); incIP(ip) {
		ipStr := ip.String()
		wg.Add(1)

		go func(host string) {
			defer wg.Done()
			// Try to connect to any port to see if host is up
			conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, "80"), timeout)
			if err == nil {
				conn.Close()
				mu.Lock()
				activeHosts = append(activeHosts, host)
				mu.Unlock()
				return
			}

			// Try alternative ports
			conn, err = net.DialTimeout("tcp", net.JoinHostPort(host, "443"), timeout)
			if err == nil {
				conn.Close()
				mu.Lock()
				activeHosts = append(activeHosts, host)
				mu.Unlock()
			}
		}(ipStr)
	}

	wg.Wait()
	return activeHosts, nil
}

// incIP increments an IP address.
func incIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

/*
EXAMPLE

func main() {
	// Get local IP
	localIP, err := GetLocalIP()
	if err != nil {
		fmt.Println("Error getting local IP:", err)
		return
	}
	fmt.Println("Your local IP:", localIP)

	// Get subnet info
	subnet, err := GetLocalSubnet()
	if err != nil {
		fmt.Println("Error getting subnet:", err)
		return
	}
	fmt.Printf("Scanning subnet: %s\n\n", subnet)

	fmt.Println("Scanning network for active hosts (this may take a minute)...")

	// Scan with 500ms timeout per host
	activeHosts, err := ScanNetwork(500 * time.Millisecond)
	if err != nil {
		fmt.Println("Error scanning network:", err)
		return
	}

	fmt.Printf("\nFound %d active host(s):\n", len(activeHosts))
	for _, host := range activeHosts {
		fmt.Printf("  - %s", host)
		if host == localIP {
			fmt.Print(" (this device)")
		}
		fmt.Println()
	}
}



*/
