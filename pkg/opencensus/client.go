package opencensus

import (
	"context"
	"net/http"

	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/b3"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
)

type preOCTransport struct {
	base http.RoundTripper
}

type postOCTransport struct {
	base http.RoundTripper
}

type clientTagKey struct{}

type clientTag struct {
	operation  string
	attributes []trace.Attribute
}

func WrapTransport(transport http.RoundTripper) http.RoundTripper {
	return &preOCTransport{
		base: &ochttp.Transport{
			Base: &postOCTransport{
				base: transport,
			},
			Propagation: &b3.HTTPFormat{},
		},
	}
}

func (t *preOCTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if data := getClientTag(req.Context()); data != nil {
		operationKey, _ := tag.NewKey("operation")
		ctx, _ := tag.New(req.Context(), tag.Upsert(operationKey, data.operation))
		req = req.WithContext(ctx)
	}
	return t.base.RoundTrip(req)
}

func (t *postOCTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if data := getClientTag(req.Context()); data != nil {
		span := trace.FromContext(req.Context())
		attributes := append(data.attributes, trace.StringAttribute("operation", data.operation))
		if len(attributes) > 0 {
			span.AddAttributes(attributes...)
		}
	}
	return t.base.RoundTrip(req)
}

func getClientTag(ctx context.Context) *clientTag {
	if data, ok := ctx.Value(clientTagKey{}).(clientTag); ok {
		return &data
	}
	return nil
}
