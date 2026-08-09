package proxy

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/StanleySun233/python-proxy/apps/node/api/internal/domain"
)

var errStreamFallbackUnavailable = errors.New("stream_fallback_unavailable")

type directPeerStreamOpener interface {
	OpenDirectStream(ctx context.Context, nextHop domain.Node, remaining []string, targetHost string, targetPort int) (net.Conn, error)
}

type fallbackStreamOpener interface {
	OpenStream(nextNodeID string, remaining []string, targetHost string, targetPort int) (net.Conn, error)
}

type directPeerAvailability interface {
	HasDirectPeer(peerNodeID string) bool
}

func openDirectFirstStream(ctx context.Context, direct directPeerStreamOpener, fallback fallbackStreamOpener, hop chainHop, targetHost string, targetPort int) (net.Conn, error) {
	var directErr error
	if direct != nil {
		conn, err := direct.OpenDirectStream(ctx, hop.node, append([]string(nil), hop.remainingHops...), targetHost, targetPort)
		if err == nil {
			return conn, nil
		}
		directErr = err
	}
	if fallback == nil {
		if directErr != nil {
			return nil, fmt.Errorf("direct_stream_failed=%v fallback_stream_failed=%w", directErr, errStreamFallbackUnavailable)
		}
		return nil, errStreamFallbackUnavailable
	}
	conn, err := fallback.OpenStream(hop.node.ID, append([]string(nil), hop.remainingHops...), targetHost, targetPort)
	if err != nil && directErr != nil {
		return nil, fmt.Errorf("direct_stream_failed=%v fallback_stream_failed=%w", directErr, err)
	}
	return conn, err
}

func openCandidateStream(ctx context.Context, direct directPeerStreamOpener, fallback fallbackStreamOpener, hops []chainHop, targetHost string, targetPort int) (net.Conn, chainHop, error) {
	var lastErr error
	for _, hop := range hops {
		conn, err := openDirectFirstStream(ctx, direct, fallback, hop, targetHost, targetPort)
		if err == nil {
			return conn, hop, nil
		}
		lastErr = err
	}
	if lastErr == nil {
		lastErr = errStreamFallbackUnavailable
	}
	return nil, chainHop{}, fmt.Errorf("chain_candidates_unavailable: %w", lastErr)
}

func (s *Server) fallbackStreamOpener() fallbackStreamOpener {
	if s.tunnelRegistry == nil {
		return nil
	}
	return s.tunnelRegistry
}
