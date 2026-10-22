package verifier

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"
)

type Verifier struct {
	client *http.Client
}

func New() *Verifier {
	return &Verifier{client: &http.Client{Timeout: 5 * time.Second}}
}

func (v *Verifier) Check(ctx context.Context, host, token string) (bool, error) {
	body, err := v.fetch(ctx, "http://"+host+"/.well-known/tenantly-verification")
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(body)) == token, nil
}

func (v *Verifier) fetch(ctx context.Context, target string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return nil, err
	}
	resp, err := v.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(io.LimitReader(resp.Body, 4096))
}
