package httpgateway

import "context"

// Get performs an HTTP GET against path.
//   - headers: per-request headers, merged on top of the client defaults.
//   - queryParams: URL query string parameters.
//   - result: optional pointer into which a 2xx JSON body is unmarshalled.
//
// Returns the Response (always non-nil when the transport succeeded) and an
// error. On non-2xx status codes the error is *HTTPError.
func (h *HTTPGateway) Get(
	ctx context.Context,
	path string,
	headers map[string]string,
	queryParams map[string]string,
	result any,
) (*Response, error) {
	req := h.client.R().SetContext(ctx)

	if len(headers) > 0 {
		req.SetHeaders(headers)
	}
	if len(queryParams) > 0 {
		req.SetQueryParams(queryParams)
	}
	if result != nil {
		req.SetResult(result)
	}

	resp, err := req.Get(path)
	if err != nil {
		return nil, err
	}
	return handleResponse(resp)
}
