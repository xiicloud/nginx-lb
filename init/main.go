package main

import (
	"encoding/json"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		startCmds()
		return
	}

	if os.Args[1] == "gen-config" {
		genConfig()
		reload()
	}

	if os.Args[1] == "xxx" {
		s := Service{
			Domain:       "pma.test.com",
			App:          "demo",
			Service:      "app",
			BackendPort:  80,
			FrontendPort: 80,
			SslCertificate: `-----BEGIN CERTIFICATE-----
MIICKTCCAZICCQCBO2ekFdKngDANBgkqhkiG9w0BAQsFADBZMQswCQYDVQQGEwJD
TjEQMA4GA1UECAwHQmVpamluZzEhMB8GA1UECgwYSW50ZXJuZXQgV2lkZ2l0cyBQ
dHkgTHRkMRUwEwYDVQQDDAxwbWEudGVzdC5jb20wHhcNMTYwNTAzMDQxOTUyWhcN
MTcwNTA0MDQxOTUyWjBZMQswCQYDVQQGEwJDTjEQMA4GA1UECAwHQmVpamluZzEh
MB8GA1UECgwYSW50ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMRUwEwYDVQQDDAxwbWEu
dGVzdC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJAoGBALo7afOGNwyXqN24
ViyG+HQG4In0O4SqrAQr6WUOg2g95rf7qSEG6U4o+wO5BM8nd6wzW0JnNqTxKXpf
nV2Ebub0uoITUHtbFMRSEYfLShQqbKBltnk09P9p4IcVhM5vUd+G9reGagaH84bt
P2bSZ5JmKBsiUf277b4ZlM/nIS4rAgMBAAEwDQYJKoZIhvcNAQELBQADgYEAOp9q
oNF/sGwEJzUKkrZ9jNfr9nvcVwsR9VajsB0cQW859IQ75r0P80NwPwJ4qbIMNsid
/1qwqzZlnYvOm01176DQTCgRC42r4vbLMKzNR4Xf6O25gwWaa2ZSSpLNQLxauKSg
2h/qgz7Rfn7rYMYZmVLNnnJqujr8GbZ4cKyuy1w=
-----END CERTIFICATE-----`,
			SslCertificateKey: `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQC6O2nzhjcMl6jduFYshvh0BuCJ9DuEqqwEK+llDoNoPea3+6kh
BulOKPsDuQTPJ3esM1tCZzak8Sl6X51dhG7m9LqCE1B7WxTEUhGHy0oUKmygZbZ5
NPT/aeCHFYTOb1Hfhva3hmoGh/OG7T9m0meSZigbIlH9u+2+GZTP5yEuKwIDAQAB
AoGAbNZeRFkzAOP9Z57cledHep+uSFF5Gz6Xi1SScWH7AEf0959XJ5sfbHNcx78w
hVR+hx/4fKVPdTQP1pncoRPNr6GcK6+9NhURiBy4oaWIAlswvuMeCFBJp4KBI3np
/hQIXKHZ4hNasD5SzHBo5bJOG6P5577KD4t39QFcBMGy8xECQQDz77IbLXOG+UWJ
Spw84HUJxzWAc7h9IgjkybJQIDoHvgmsYIjLFnLimwIbyvMlPBnyAopvZoO5tabX
WfGmQHJpAkEAw3Eqe7XG2j023SUDe8slqOIlnivdP95M2tACOKIvbJFB83zHDK2C
46cdyoqfI/bLimnxgPxYMq5CTr33I7RhcwJBAJTlU17ZcHILx5EU1KcoDuiICzU7
7XmcA7e7Ebds5F8DdZ4dUoI8UqXVHgVe7OlmdSPOvzdeaLs7kPpUMXdcUTkCQGG5
umZ1dGM37LETivRhlgkmW20FvfHrtD5NeG7dGh2NXI7lu5opQKOYsprOSdjv1ML3
Sp0WkPt2iw1Yi7U8wuUCQQCh3Z/svnkDSxBrexFUwt5RLhF5YBJUtIgeOBQ/wlDp
hQPBEwDsQoM9LxnEmVNyzeF8Yz9RJ6gANpnZYLszWrGk
-----END RSA PRIVATE KEY-----`,
		}
		json.NewEncoder(os.Stdout).Encode(s)
	}
}
