package irr

import (
	"fmt"
	"io"
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

	// 🔥 Lê até EOF
	response, err := io.ReadAll(conn)
	if err != nil {
		return "", err
	}

	return string(response), nil
}
