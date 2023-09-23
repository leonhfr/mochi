package action

import "net/http"

var _ http.RoundTripper = &ThrottledTransport{}
