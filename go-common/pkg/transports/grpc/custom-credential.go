package grpc

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/credentials"
)

type ExtraMeta struct {
	appId  string
	appKey string
	meta   map[string]string
}

type customCredential struct {
	meta *ExtraMeta
}

func newCustomCredential(meta *ExtraMeta) credentials.PerRPCCredentials {
	c := new(customCredential)
	c.meta = meta
	return c
}

func (c customCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	kv := map[string]string{
		"appId":  c.meta.appId,
		"appKey": c.meta.appKey,
	}
	if c.meta.meta != nil && len(c.meta.meta) > 0 {
		for k, v := range c.meta.meta {
			kv[k] = v
		}
	}
	return kv, nil
}

func (c customCredential) RequireTransportSecurity() bool {
	return false
}
