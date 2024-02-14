package grpcApi

import (
	"context"
	metadata2 "google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	userAgentHeader            = "user-agent"
	xForwardForHeader          = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	metadata := &Metadata{}
	if md, ok := metadata2.FromIncomingContext(ctx); ok {

		if userAgents := md.Get(grpcGatewayUserAgentHeader); len(userAgents) > 0 {
			metadata.UserAgent = userAgents[0]
		}

		if userAgents := md.Get(userAgentHeader); len(userAgents) > 0 {
			metadata.UserAgent = userAgents[0]
		}

		if clientIPs := md.Get(xForwardForHeader); len(clientIPs) > 0 {
			metadata.ClientIP = clientIPs[0]
		}
	}
	if p, ok := peer.FromContext(ctx); ok {
		metadata.ClientIP = p.Addr.String()
	}
	return metadata
}
