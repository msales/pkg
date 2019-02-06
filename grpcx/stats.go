package grpcx

import (
	"github.com/msales/pkg/v3/stats"
	"golang.org/x/net/context"
	grpcstats "google.golang.org/grpc/stats"
)

const (
	methodKey ctxKey = iota
	tagsKey
)

type ctxKey int

// Compile-time interface checks.
var (
	_ grpcstats.Handler = (*rpcStatsHandler)(nil)
	_ grpcstats.Handler = (*aggregateHandler)(nil)
	_ grpcstats.Handler = (*handler)(nil)
)

// WithHandlers returns an aggregated stats handler. All inner handlers are called in order.
func WithHandlers(handlers ...grpcstats.Handler) grpcstats.Handler {
	return &aggregateHandler{
		handlers: handlers,
	}
}

// WithRPCStats returns a handler that collects and tags RPC stats.
func WithRPCStats(stats stats.Stats, tagsFns ...TagsFunc) grpcstats.Handler {
	return &rpcStatsHandler{stats: stats, tagsFns: tagsFns}
}

// TagsFunc represents a function that returns a set of tags for the stats.
// The payload is the incoming message object.
type TagsFunc func(ctx context.Context, payload interface{}) stats.Tags

// rpcStatsHandler records stats of each RPC: message rate and request duration.
type rpcStatsHandler struct {
	handler

	stats   stats.Stats
	tagsFns []TagsFunc
}

// TagRPC can attach some information to the given context.
func (h *rpcStatsHandler) TagRPC(ctx context.Context, info *grpcstats.RPCTagInfo) context.Context {
	if h.tagsFns != nil {
		ctx = context.WithValue(ctx, tagsKey, &stats.Tags{})
	}
	return context.WithValue(ctx, methodKey, info.FullMethodName)
}

// HandleRPC processes the RPC stats.
func (h *rpcStatsHandler) HandleRPC(ctx context.Context, rpcStats grpcstats.RPCStats) {
	if in, ok := rpcStats.(*grpcstats.InPayload); ok && h.tagsFns != nil {
		tags := make(stats.Tags)

		for _, fn := range h.tagsFns {
			tags.Merge(fn(ctx, in.Payload))
		}

		ctxTags := ctx.Value(tagsKey).(*stats.Tags)
		*ctxTags = tags

		return
	}

	if end, ok := rpcStats.(*grpcstats.End); ok {
		tags := stats.Tags{}

		t, ok := ctx.Value(tagsKey).(*stats.Tags)
		if ok {
			tags = *t
		}

		tags.
			With("method", h.methodFromContext(ctx)).
			With("status", h.getStatus(end))

		h.stats.Inc("rpc.end", 1, 1, tags)
		h.stats.Timing("rpc.time", end.EndTime.Sub(end.BeginTime), 1, tags)
	}
}

// methodFromContext retrieves a full RPC method name from context.
func (h *rpcStatsHandler) methodFromContext(ctx context.Context) string {
	method := ctx.Value(methodKey)
	if method == nil {
		method = "unknown"
	}

	return method.(string)
}

// getStatus returns the status of the current RPC.
func (h *rpcStatsHandler) getStatus(end *grpcstats.End) string {
	status := "success"
	if end.Error != nil {
		status = "error"
	}

	return status
}

// aggregateHandler represents an aggregated stats handler.
type aggregateHandler struct {
	handlers []grpcstats.Handler
}

// TagRPC can attach some information to the given context.
func (a *aggregateHandler) TagRPC(ctx context.Context, info *grpcstats.RPCTagInfo) context.Context {
	a.withEachHandler(func(h grpcstats.Handler) {
		ctx = h.TagRPC(ctx, info)
	})

	return ctx
}

// HandleRPC processes the RPC stats.
func (a *aggregateHandler) HandleRPC(ctx context.Context, rpcStats grpcstats.RPCStats) {
	a.withEachHandler(func(h grpcstats.Handler) {
		h.HandleRPC(ctx, rpcStats)
	})
}

// TagConn can attach some information to the given context.
func (a *aggregateHandler) TagConn(ctx context.Context, connStats *grpcstats.ConnTagInfo) context.Context {
	a.withEachHandler(func(h grpcstats.Handler) {
		ctx = h.TagConn(ctx, connStats)
	})

	return ctx
}

// HandleConn processes the Conn stats.
func (a *aggregateHandler) HandleConn(ctx context.Context, connStats grpcstats.ConnStats) {
	a.withEachHandler(func(h grpcstats.Handler) {
		h.HandleConn(ctx, connStats)
	})
}

// withEachHandler executes a callback on each inner handler.
func (a *aggregateHandler) withEachHandler(fn func(grpcstats.Handler)) {
	for _, h := range a.handlers {
		fn(h)
	}
}

// handler represents a no-op stats handler. Can be used as a base for specialised handlers.
type handler struct{}

// TagRPC can attach some information to the given context.
func (*handler) TagRPC(ctx context.Context, _ *grpcstats.RPCTagInfo) context.Context {
	return ctx
}

// HandleRPC processes the RPC stats.
func (*handler) HandleRPC(context.Context, grpcstats.RPCStats) {}

// TagConn can attach some information to the given context.
func (*handler) TagConn(ctx context.Context, _ *grpcstats.ConnTagInfo) context.Context {
	return ctx
}

// HandleConn processes the Conn stats.
func (*handler) HandleConn(context.Context, grpcstats.ConnStats) {}
