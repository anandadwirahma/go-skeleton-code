package httpgateway

import "context"

// Delete performs an HTTP DELETE against path.
//   - headers: per-request headers, merged on top of the client defaults.
//   - result: optional pointer into which a 2xx JSON body is unmarshalled
//     (some APIs return the deleted resource).
//
// Returns the Response and an error; on non-2xx status codes the error is *HTTPError.
func (h *HTTPGateway) Delete(
	ctx context.Context,
	path string,
	headers map[string]string,
	result any,
) (*Response, error) {
	req := h.client.R().SetContext(ctx)

	if len(headers) > 0 {
		req.SetHeaders(headers)
	}
	if result != nil {
		req.SetResult(result)
	}

	resp, err := req.Delete(path)
	if err != nil {
		return nil, err
	}
	return handleResponse(resp)
}
