package ctx

import "context"

// keyRequestID is the context key for request IDs.
const keyRequestID key = "requestID"

// key is a custom type for context keys in this package.
type key string

// RequestID extracts the request ID from the context.
func RequestID(ctx context.Context) string {
	requestID, _ := ctx.Value(keyRequestID).(string)
	return requestID
}

// SetRequestID sets the request ID in the context.
func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, keyRequestID, requestID)
}
