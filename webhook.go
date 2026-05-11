package stream_chat

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// ErrInvalidWebhookSignature is returned by VerifyAndParse* helpers when
// the supplied signature does not match the HMAC computed over the
// uncompressed JSON body.
var ErrInvalidWebhookSignature = errors.New("invalid webhook signature")

var gzipMagic = []byte{0x1f, 0x8b}

// GunzipPayload returns body unchanged unless the first two bytes are
// the gzip magic (1f 8b, per RFC 1952), in which case the gzip stream
// is inflated and the decompressed bytes are returned.
//
// Magic-byte detection lets the same handler stay correct when
// middleware auto-decompresses the request before your code sees it.
func GunzipPayload(body []byte) ([]byte, error) {
	if len(body) < 2 || !bytes.Equal(body[:2], gzipMagic) {
		return body, nil
	}
	zr, err := gzip.NewReader(bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("decompress gzip payload: %w", err)
	}
	out, err := io.ReadAll(zr)
	if err != nil {
		_ = zr.Close()
		return nil, fmt.Errorf("read gzip payload: %w", err)
	}
	if err := zr.Close(); err != nil {
		return nil, fmt.Errorf("finalize gzip payload: %w", err)
	}
	return out, nil
}

// DecodeSqsPayload reverses the SQS firehose envelope: the message Body
// is base64-decoded and, when the result begins with the gzip magic, it
// is gzip-decompressed. The same call works whether or not Stream is
// currently compressing payloads.
func DecodeSqsPayload(body string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(body)
	if err != nil {
		return nil, fmt.Errorf("base64-decode payload: %w", err)
	}
	return GunzipPayload(decoded)
}

// DecodeSnsPayload reverses an SNS HTTP notification envelope: when the
// input is a JSON envelope ({"Type":"Notification","Message":"..."}),
// the inner Message field is extracted and run through the SQS pipeline
// (base64-decode, then gzip-if-magic). When the input is not a JSON
// envelope it is treated as the already-extracted Message string, so
// existing call sites that pre-unwrap continue to work.
func DecodeSnsPayload(notificationBody string) ([]byte, error) {
	if msg, ok := extractSnsMessage(notificationBody); ok {
		return DecodeSqsPayload(msg)
	}
	return DecodeSqsPayload(notificationBody)
}

// extractSnsMessage returns the inner Message field from an SNS HTTP
// notification envelope. The ok result is false when input is not a JSON
// object with a string Message field.
func extractSnsMessage(notificationBody string) (string, bool) {
	trimmed := bytes.TrimLeft([]byte(notificationBody), " \t\r\n")
	if len(trimmed) == 0 || trimmed[0] != '{' {
		return "", false
	}
	var envelope struct {
		Message *string `json:"Message"`
	}
	if err := json.Unmarshal(trimmed, &envelope); err != nil {
		return "", false
	}
	if envelope.Message == nil {
		return "", false
	}
	return *envelope.Message, true
}

// VerifySignature returns true when signature equals the hex-encoded
// HMAC-SHA256 of body using secret as the key. The comparison is
// constant-time. The signature is always computed over the
// uncompressed JSON bytes, so callers that decoded a gzipped or
// base64-wrapped payload must pass the inflated bytes here.
func VerifySignature(body []byte, signature, secret string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write(body)
	expected := []byte(hex.EncodeToString(mac.Sum(nil)))
	return hmac.Equal(expected, []byte(signature))
}

// ParseEvent decodes the JSON-encoded webhook payload into a typed
// Event. Unknown event types still parse successfully because Event.Type
// is a string alias.
func ParseEvent(payload []byte) (*Event, error) {
	var ev Event
	if err := json.Unmarshal(payload, &ev); err != nil {
		return nil, fmt.Errorf("parse webhook event: %w", err)
	}
	return &ev, nil
}

func verifyAndParse(payload []byte, signature, secret string) (*Event, error) {
	if !VerifySignature(payload, signature, secret) {
		return nil, ErrInvalidWebhookSignature
	}
	return ParseEvent(payload)
}

// VerifyAndParseWebhook decompresses body when gzipped, verifies the
// HMAC signature against secret, and returns the parsed Event. Returns
// ErrInvalidWebhookSignature on mismatch and a wrapped error on any
// decode failure.
func VerifyAndParseWebhook(body []byte, signature, secret string) (*Event, error) {
	inflated, err := GunzipPayload(body)
	if err != nil {
		return nil, err
	}
	return verifyAndParse(inflated, signature, secret)
}

// VerifyAndParseSqs decodes the SQS message Body, verifies the
// X-Signature attribute against secret, and returns the parsed Event.
func VerifyAndParseSqs(messageBody, signature, secret string) (*Event, error) {
	inflated, err := DecodeSqsPayload(messageBody)
	if err != nil {
		return nil, err
	}
	return verifyAndParse(inflated, signature, secret)
}

// VerifyAndParseSns decodes the SNS notification Message, verifies the
// X-Signature attribute against secret, and returns the parsed Event.
func VerifyAndParseSns(message, signature, secret string) (*Event, error) {
	inflated, err := DecodeSnsPayload(message)
	if err != nil {
		return nil, err
	}
	return verifyAndParse(inflated, signature, secret)
}

// VerifyAndParseWebhook is the client-bound form of the package-level
// helper; it pulls the API secret from the receiver so call sites only
// supply the request body and signature.
func (c *Client) VerifyAndParseWebhook(body []byte, signature string) (*Event, error) {
	return VerifyAndParseWebhook(body, signature, string(c.apiSecret))
}

// VerifyAndParseSqs is the client-bound form of the package-level helper.
func (c *Client) VerifyAndParseSqs(messageBody, signature string) (*Event, error) {
	return VerifyAndParseSqs(messageBody, signature, string(c.apiSecret))
}

// VerifyAndParseSns is the client-bound form of the package-level helper.
func (c *Client) VerifyAndParseSns(message, signature string) (*Event, error) {
	return VerifyAndParseSns(message, signature, string(c.apiSecret))
}
