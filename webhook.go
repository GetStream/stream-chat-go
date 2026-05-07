package stream_chat

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
)

// ErrInvalidWebhookSignature is returned by VerifyAndDecodeWebhook when the
// provided signature does not match the HMAC computed over the uncompressed
// JSON body.
var ErrInvalidWebhookSignature = errors.New("invalid webhook signature")

// DecompressWebhookBody undoes the encoding wrappers Stream applies to
// outbound webhook / SQS / SNS payloads, returning the raw JSON bytes
// the server signed.
//
//   - payloadEncoding: optional transport wrapper applied last by the
//     server. Currently the only supported value is "base64" (used for
//     SQS / SNS firehose so the message stays valid UTF-8). Pass "" for
//     HTTP webhooks.
//   - contentEncoding: optional compression. Currently the only
//     supported value is "gzip". Pass "" when no compression is set.
//
// Decode order is the inverse of how the server built the message:
// base64 first, then gunzip. Both are case-insensitive and trimmed.
func (c *Client) DecompressWebhookBody(body []byte, contentEncoding, payloadEncoding string) ([]byte, error) {
	out := body

	if pe := strings.ToLower(strings.TrimSpace(payloadEncoding)); pe != "" {
		switch pe {
		case "base64", "b64":
			decoded, err := decodeBase64(out)
			if err != nil {
				return nil, fmt.Errorf("decode webhook payload_encoding=base64: %w", err)
			}
			out = decoded
		default:
			return nil, fmt.Errorf("unsupported webhook payload_encoding: %s. This SDK only supports base64.", payloadEncoding)
		}
	}

	if ce := strings.ToLower(strings.TrimSpace(contentEncoding)); ce != "" {
		switch ce {
		case "gzip":
			decompressed, err := gunzip(out)
			if err != nil {
				return nil, fmt.Errorf("decompress webhook Content-Encoding=gzip: %w", err)
			}
			out = decompressed
		default:
			return nil, fmt.Errorf(`unsupported webhook Content-Encoding: %s. This SDK only supports gzip; set webhook_compression_algorithm to "gzip" on the app config.`, contentEncoding)
		}
	}

	return out, nil
}

// VerifyAndDecodeWebhook decompresses (when needed), verifies the HMAC
// signature, and returns the uncompressed JSON bytes. The signature is
// always computed over the innermost (uncompressed, base64-decoded)
// JSON, so the verification rule is invariant across HTTP webhooks and
// SQS / SNS.
//
//   - body: raw HTTP request body / SQS Body / SNS Message bytes
//   - signature: value of the X-Signature header / message attribute
//   - contentEncoding: value of Content-Encoding header / attribute
//   - payloadEncoding: "base64" for SQS / SNS firehose, "" for HTTP webhooks
//
// Returns ErrInvalidWebhookSignature when the signature does not match.
func (c *Client) VerifyAndDecodeWebhook(body []byte, signature, contentEncoding, payloadEncoding string) ([]byte, error) {
	decoded, err := c.DecompressWebhookBody(body, contentEncoding, payloadEncoding)
	if err != nil {
		return nil, err
	}

	mac := hmac.New(sha256.New, c.apiSecret)
	_, _ = mac.Write(decoded)
	expected := []byte(hex.EncodeToString(mac.Sum(nil)))

	if !hmac.Equal([]byte(signature), expected) {
		return nil, ErrInvalidWebhookSignature
	}

	return decoded, nil
}

func decodeBase64(b []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(b))
}

func gunzip(b []byte) ([]byte, error) {
	zr, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	defer zr.Close()
	return io.ReadAll(zr)
}
