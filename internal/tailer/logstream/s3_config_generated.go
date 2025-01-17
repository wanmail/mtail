// Code generated by github.com/wanmail/mtail/cmd/config-gen; DO NOT EDIT.
package logstream

import (
	"net/url"
	"github.com/aws/aws-sdk-go-v2/aws"
	"strconv"
)

func parseS3Config(u *url.URL, config *aws.Config) error {

	if Region := u.Query().Get("Region"); Region != "" {
		config.Region = Region
	}

  // Credentials is not supported

  // BearerAuthTokenProvider is not supported

  // HTTPClient is not supported

  // EndpointResolver is not supported

  // EndpointResolverWithOptions is not supported

	if RetryMaxAttempts := u.Query().Get("RetryMaxAttempts"); RetryMaxAttempts != "" {
		i, err := strconv.Atoi(RetryMaxAttempts)
		if err != nil {
			return err
		}
		config.RetryMaxAttempts = i
	}

	if RetryMode := u.Query().Get("RetryMode"); RetryMode != "" {
		config.RetryMode = aws.RetryMode(RetryMode)
	}

  // Retryer is not supported

  // Logger is not supported

  // ClientLogMode is not supported

	if DefaultsMode := u.Query().Get("DefaultsMode"); DefaultsMode != "" {
		config.DefaultsMode = aws.DefaultsMode(DefaultsMode)
	}

  // RuntimeEnvironment is not supported

	if AppID := u.Query().Get("AppID"); AppID != "" {
		config.AppID = AppID
	}

  // BaseEndpoint is not supported

	if DisableRequestCompression := u.Query().Get("DisableRequestCompression"); DisableRequestCompression != "" {
		b, err := strconv.ParseBool(DisableRequestCompression)
		if err != nil {
			return err
		}
		config.DisableRequestCompression = b
	}

	if RequestMinCompressSizeBytes := u.Query().Get("RequestMinCompressSizeBytes"); RequestMinCompressSizeBytes != "" {
		i, err := strconv.Atoi(RequestMinCompressSizeBytes)
		if err != nil {
			return err
		}
		config.RequestMinCompressSizeBytes = int64(i)
	}

	if AccountIDEndpointMode := u.Query().Get("AccountIDEndpointMode"); AccountIDEndpointMode != "" {
		config.AccountIDEndpointMode = aws.AccountIDEndpointMode(AccountIDEndpointMode)
	}

  return nil
}
