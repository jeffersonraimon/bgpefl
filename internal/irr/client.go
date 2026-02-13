package irr

import (
	"fmt"
	"net"
)

func FetchPrefixes(host string, asn uint32) (string, error) {

	conn, err := net.Dial("tcp", host+":43")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	query := fmt.Sprintf("-i origin AS%d\n", asn)

	_, err = conn.Write([]byte(query))
	if err != nil {
		return "", err
	}

	buf := make([]byte, 1024*1024)
	n, err := conn.Read(buf)
	if err != nil {
		return "", err
	}

	return string(buf[:n]), nil
}
