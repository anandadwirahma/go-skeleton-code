package httpgateway

import "context"

// Put performs an HTTP PUT against path.
//   - headers: per-request headers, merged on top of the client defaults.
//   - body: any value resty can JSON-marshal; pass nil for empty body.
//   - result: optional pointer into which a 2xx JSON body is unmarshalled.
//
// Returns the Response and an error; on non-2xx status codes the error is *HTTPError.
func (h *HTTPGateway) Put(
	ctx context.Context,
	path string,
	headers map[string]string,
	body any,
	result any,
) (*Response, error) {
	req := h.client.R().SetContext(ctx)

	if len(headers) > 0 {
		req.SetHeaders(headers)
	}
	if body != nil {
		req.SetBody(body)
	}
	if result != nil {
		req.SetResult(result)
	}

	resp, err := req.Put(path)
	if err != nil {
		return nil, err
	}
	return handleResponse(resp)
}
