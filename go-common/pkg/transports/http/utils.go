package http

import (
	"path/filepath"
	"strings"
)

func formatCertFile(o *SSLOptions) {
	if !strings.HasPrefix(o.CertFile, "configs") {
		o.CertFile = filepath.Join("configs", o.CertFile)
	}
	if !strings.HasPrefix(o.KeyFile, "configs") {
		o.KeyFile = filepath.Join("configs", o.KeyFile)
	}
	if !strings.HasPrefix(o.TrustedCaFile, "configs") {
		o.TrustedCaFile = filepath.Join("configs", o.TrustedCaFile)
	}
	o.CertFile = filepath.Clean(o.CertFile)
	o.KeyFile = filepath.Clean(o.KeyFile)
	o.TrustedCaFile = filepath.Clean(o.TrustedCaFile)

	o.KeyFileBytes = RSAReadKeyFromFile(o.KeyFile)
	o.CertFileBytes = RSAReadKeyFromFile(o.CertFile)
	o.TrustedCaBytes = RSAReadKeyFromFile(o.TrustedCaFile)
}
