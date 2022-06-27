package proto

import (
	"context"
	"strings"

	"github.com/mediocregopher/radix/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"laxo.vn/laxo/laxo"
)

var errInvalidToken = status.Errorf(codes.Unauthenticated, "invalid token")

type key int

const (
	keyUID key = iota
)

type GRPCMiddleware struct {
	server *laxo.Server
}

func NewGRPCMiddleware(server *laxo.Server) GRPCMiddleware {
	return GRPCMiddleware{
		server: server,
	}
}

func (g *GRPCMiddleware) StreamAuthFunc(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
	}

	cookie, cookiePresent := md["cookie"]

	if !cookiePresent {
		return nil, errInvalidToken
	}

	if len(cookie) == 0 {
		return nil, errInvalidToken
	}

	arr := strings.Split(cookie[0], ";")

	var token string
	for _, str := range arr {
		str = strings.TrimSpace(str)
		if strings.HasPrefix(str, g.server.Config.AuthCookieName) {
			token = str[len(g.server.Config.AuthCookieName)+1:]
		}
	}

	if token == "" {
		return nil, errInvalidToken
	}

	var uID string
	redisCtx := context.Background()
	err := g.server.RedisClient.Do(redisCtx, radix.Cmd(&uID, "GET", token))

	if err != nil {
		g.server.Logger.Errorw("Error in grpc auth interceptor function (Redis)",
			"error", err,
		)
		return nil, errInvalidToken
	}

	if uID == "" {
		return nil, errInvalidToken
	}

	newCtx := context.WithValue(ctx, keyUID, uID)
	return newCtx, nil
}
