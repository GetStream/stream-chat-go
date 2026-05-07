package stream_chat

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	webhookTestAPIKey    = "tkey1"
	webhookTestAPISecret = "tsec2"
	webhookTestFixture   = `{"type":"message.new","message":{"text":"the quick brown fox"}}`
)

func newWebhookTestClient(t *testing.T) *Client {
	t.Helper()
	c, err := NewClient(webhookTestAPIKey, webhookTestAPISecret)
	require.NoError(t, err, "new client")
	return c
}

func hmacHex(t *testing.T, secret, body []byte) []byte {
	t.Helper()
	mac := hmac.New(sha256.New, secret)
	_, err := mac.Write(body)
	require.NoError(t, err)
	return []byte(hex.EncodeToString(mac.Sum(nil)))
}

func gzipBytes(t *testing.T, src []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err := zw.Write(src)
	require.NoError(t, err)
	require.NoError(t, zw.Close())
	return buf.Bytes()
}

func base64Bytes(t *testing.T, src []byte) []byte {
	t.Helper()
	return []byte(base64.StdEncoding.EncodeToString(src))
}

func TestVerifyWebhook_BackwardCompatibility(t *testing.T) {
	c := newWebhookTestClient(t)
	body := []byte(webhookTestFixture)
	sig := hmacHex(t, []byte(webhookTestAPISecret), body)

	require.True(t, c.VerifyWebhook(body, sig), "valid signature must verify")
	require.False(t, c.VerifyWebhook(body, []byte("not-a-valid-hex-signature")), "invalid signature must not verify")
	require.False(t, c.VerifyWebhook([]byte("tampered"), sig), "tampered body must not verify")
}

func TestDecompressWebhookBody(t *testing.T) {
	c := newWebhookTestClient(t)
	body := []byte(webhookTestFixture)

	gzipped := gzipBytes(t, body)
	b64Plain := base64Bytes(t, body)
	b64Gzipped := base64Bytes(t, gzipped)

	tests := []struct {
		name            string
		body            []byte
		contentEncoding string
		payloadEncoding string
		want            []byte
	}{
		{
			name: "passthrough when both encodings empty",
			body: body,
			want: body,
		},
		{
			name:            "passthrough when both encodings whitespace",
			body:            body,
			contentEncoding: "  ",
			payloadEncoding: "\t",
			want:            body,
		},
		{
			name:            "gzip round-trip",
			body:            gzipped,
			contentEncoding: "gzip",
			want:            body,
		},
		{
			name:            "base64 round-trip without compression",
			body:            b64Plain,
			payloadEncoding: "base64",
			want:            body,
		},
		{
			name:            "base64 + gzip round-trip (SQS / SNS shape)",
			body:            b64Gzipped,
			contentEncoding: "gzip",
			payloadEncoding: "base64",
			want:            body,
		},
		{
			name:            "case-insensitive GZIP",
			body:            gzipped,
			contentEncoding: "GZIP",
			want:            body,
		},
		{
			name:            "case-insensitive BASE64",
			body:            b64Plain,
			payloadEncoding: "BASE64",
			want:            body,
		},
		{
			name:            "b64 alias",
			body:            b64Plain,
			payloadEncoding: "b64",
			want:            body,
		},
		{
			name:            "b64 alias + gzip",
			body:            b64Gzipped,
			contentEncoding: "gzip",
			payloadEncoding: "B64",
			want:            body,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.DecompressWebhookBody(tt.body, tt.contentEncoding, tt.payloadEncoding)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestDecompressWebhookBody_RejectsUnsupportedContentEncoding(t *testing.T) {
	c := newWebhookTestClient(t)
	body := []byte(webhookTestFixture)

	for _, ce := range []string{"br", "brotli", "zstd", "deflate", "compress", "lz4"} {
		ce := ce
		t.Run(ce, func(t *testing.T) {
			got, err := c.DecompressWebhookBody(body, ce, "")
			require.Error(t, err)
			require.Nil(t, got)
			msg := err.Error()
			assert.Contains(t, msg, "unsupported")
			assert.Contains(t, msg, "gzip")
		})
	}
}

func TestDecompressWebhookBody_RejectsUnsupportedPayloadEncoding(t *testing.T) {
	c := newWebhookTestClient(t)
	body := []byte(webhookTestFixture)

	for _, pe := range []string{"hex", "url", "binary"} {
		pe := pe
		t.Run(pe, func(t *testing.T) {
			got, err := c.DecompressWebhookBody(body, "", pe)
			require.Error(t, err)
			require.Nil(t, got)
			assert.Contains(t, err.Error(), "unsupported")
			assert.Contains(t, err.Error(), "payload_encoding")
		})
	}
}

func TestDecompressWebhookBody_InvalidGzipBytes(t *testing.T) {
	c := newWebhookTestClient(t)

	got, err := c.DecompressWebhookBody([]byte("not-actually-gzip"), "gzip", "")
	require.Error(t, err)
	require.Nil(t, got)
	assert.Contains(t, err.Error(), "gzip")
}

func TestDecompressWebhookBody_InvalidBase64Input(t *testing.T) {
	c := newWebhookTestClient(t)

	got, err := c.DecompressWebhookBody([]byte("!!!not base64!!!"), "", "base64")
	require.Error(t, err)
	require.Nil(t, got)
	assert.Contains(t, err.Error(), "payload_encoding")
}

func TestVerifyAndDecodeWebhook_Plain(t *testing.T) {
	c := newWebhookTestClient(t)
	body := []byte(webhookTestFixture)
	sig := hmacHex(t, []byte(webhookTestAPISecret), body)

	got, err := c.VerifyAndDecodeWebhook(body, string(sig), "", "")
	require.NoError(t, err)
	require.Equal(t, body, got)
}

func TestVerifyAndDecodeWebhook_Gzip(t *testing.T) {
	c := newWebhookTestClient(t)
	body := []byte(webhookTestFixture)
	gzipped := gzipBytes(t, body)
	sig := hmacHex(t, []byte(webhookTestAPISecret), body)

	got, err := c.VerifyAndDecodeWebhook(gzipped, string(sig), "gzip", "")
	require.NoError(t, err)
	require.Equal(t, body, got)
}

func TestVerifyAndDecodeWebhook_Base64Gzip(t *testing.T) {
	c := newWebhookTestClient(t)
	body := []byte(webhookTestFixture)
	wrapped := base64Bytes(t, gzipBytes(t, body))
	sig := hmacHex(t, []byte(webhookTestAPISecret), body)

	got, err := c.VerifyAndDecodeWebhook(wrapped, string(sig), "gzip", "base64")
	require.NoError(t, err)
	require.Equal(t, body, got)
}

func TestVerifyAndDecodeWebhook_SignatureMismatch(t *testing.T) {
	c := newWebhookTestClient(t)
	body := []byte(webhookTestFixture)
	wrongSig := hex.EncodeToString(make([]byte, sha256.Size))

	got, err := c.VerifyAndDecodeWebhook(body, wrongSig, "", "")
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrInvalidWebhookSignature))
	require.Nil(t, got)
}

func TestVerifyAndDecodeWebhook_RejectsSignatureOverCompressedBytes(t *testing.T) {
	c := newWebhookTestClient(t)
	body := []byte(webhookTestFixture)
	gzipped := gzipBytes(t, body)
	sigOverCompressed := hmacHex(t, []byte(webhookTestAPISecret), gzipped)

	got, err := c.VerifyAndDecodeWebhook(gzipped, string(sigOverCompressed), "gzip", "")
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrInvalidWebhookSignature))
	require.Nil(t, got)
}

func TestVerifyAndDecodeWebhook_RejectsSignatureOverWrappedBytes(t *testing.T) {
	c := newWebhookTestClient(t)
	body := []byte(webhookTestFixture)
	wrapped := base64Bytes(t, gzipBytes(t, body))
	sigOverWrapped := hmacHex(t, []byte(webhookTestAPISecret), wrapped)

	got, err := c.VerifyAndDecodeWebhook(wrapped, string(sigOverWrapped), "gzip", "base64")
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrInvalidWebhookSignature))
	require.Nil(t, got)
}

func TestVerifyAndDecodeWebhook_PropagatesDecompressionError(t *testing.T) {
	c := newWebhookTestClient(t)
	bogus := []byte("definitely-not-gzip-bytes")
	sig := hmacHex(t, []byte(webhookTestAPISecret), bogus)

	got, err := c.VerifyAndDecodeWebhook(bogus, string(sig), "gzip", "")
	require.Error(t, err)
	require.False(t, errors.Is(err, ErrInvalidWebhookSignature))
	require.Nil(t, got)
	assert.True(t, strings.Contains(err.Error(), "gzip"))
}
