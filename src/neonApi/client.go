package neonApi

import (
	"fmt"

	"github.com/imroc/req/v3"
)

type NeonApiClient struct {
	*req.Client
}

type NeonApiClientOptions struct {
	NumRetries int
}

type NeonApiErrorResponseBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type NeonApiError struct {
	Code     string
	Message  string
	Response *req.Response
}

type NeonApiRequestResult interface {
	NeonProject
}

func NewNeonApiClient(httpClient *req.Client, authToken string) NeonApiClient {

	httpClient.
		SetCommonHeader("Accept", "application/json").
		SetCommonBearerAuthToken(authToken).
		SetBaseURL("https://console.neon.tech/").
		// EnableDump at the request level in request middleware which dump content into
		// memory (not print to stdout), we can record dump content only when unexpected
		// exception occurs, it is helpful to troubleshoot problems in production.
		OnBeforeRequest(func(c *req.Client, r *req.Request) error {
			if r.RetryAttempt > 0 { // Ignore on retry.
				return nil
			}
			r.EnableDump()
			return nil
		}).
		SetCommonError(&NeonApiErrorResponseBody{}).
		OnAfterResponse(func(client *req.Client, resp *req.Response) error {
			if err, ok := resp.Error().(*NeonApiErrorResponseBody); ok {
				// Server returns an error message, convert it to human-readable go error.
				err := NeonApiError{
					Code:     err.Code,
					Message:  err.Message,
					Response: resp,
				}
				return err
			}
			// Corner case: neither an error response nor a success response,
			// dump content to help troubleshoot.
			if !resp.IsSuccess() {
				return fmt.Errorf("Neon API request failed. Bad response, raw dump:\n%s", resp.Dump())
			}
			return nil
		})

	return NeonApiClient{
		httpClient,
	}
}

func (e NeonApiError) Error() string {
	return fmt.Sprintf("Neon API request failed. request_url: %s status_code: %d message: %s code: %s", e.Response.Request.URL.String(), e.Response.StatusCode, e.Message, e.Code)
}

func (c *NeonApiClient) SetDebug(enable bool) *NeonApiClient {
	if enable {
		c.EnableDebugLog()
		c.EnableDumpAll()
	} else {
		c.DisableDebugLog()
		c.DisableDumpAll()
	}
	return c
}

// func NeonApiRequestCreate(c *NeonApiClient, endpointPath string, options NeonApiClientOptions) (Result, error) {

// 	url, err := url.JoinPath("https://console.neon.tech/api/v1/", endpointPath)
// 	if err != nil {
// 		return Result{}, errors.New(fmt.Sprintf("Could not join path. path: %s", endpointPath))
// 	}

// 	var result Result
// 	var apiErr NeonApiResponseError

// 	resp, err := c.httpClient.R().
// 		SetHeader("Accept", "application/json").
// 		SetBearerAuthToken(options.authToken).
// 		SetResult(&result).
// 		SetError(&apiErr).
// 		SetURL(url)

// 	if err != nil {
// 		return Result{}, errors.New(fmt.Sprintf("Neon API request failed. request_path: %s, status_code: %d, raw_response: %+v, error: %+v", resp.Request.URL.Path, resp.StatusCode, resp, err))
// 	}

// 	if apiErr.Message != "" {
// 		return Result{}, errors.New(fmt.Sprintf("Neon API request failed. request_path: %s, status_code: %d, error_code: %s, error_message: %s, raw_response: %+v", resp.Request.URL.Path, resp.StatusCode, apiErr.Code, apiErr.Message, resp))
// 	}

// 	return result, nil
// }

// func NeonApiRequestExecute()
