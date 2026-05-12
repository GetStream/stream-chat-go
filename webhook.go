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

// ErrInvalidWebhook is the single sentinel that every webhook failure
// path wraps. Callers can branch on it with errors.Is and, when they
// need to distinguish the failure mode (signature mismatch, base64,
// gzip, JSON), match a substring of the error message.
var ErrInvalidWebhook = errors.New("invalid webhook")

var gzipMagic = []byte{0x1f, 0x8b}

// GunzipPayload returns body unchanged unless the first two bytes are
// the gzip magic (1f 8b, per RFC 1952), in which case the gzip stream
// is inflated and the decompressed bytes are returned.
//
// Magic-byte detection lets the same handler stay correct when
// middleware auto-decompresses the request before your code sees it.
//
// Any failure to inflate the gzip stream is wrapped in ErrInvalidWebhook
// with the prefix "gzip decompression failed".
func GunzipPayload(body []byte) ([]byte, error) {
	if len(body) < 2 || !bytes.Equal(body[:2], gzipMagic) {
		return body, nil
	}
	zr, err := gzip.NewReader(bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("gzip decompression failed: %v: %w", err, ErrInvalidWebhook)
	}
	out, err := io.ReadAll(zr)
	if err != nil {
		_ = zr.Close()
		return nil, fmt.Errorf("gzip decompression failed: %v: %w", err, ErrInvalidWebhook)
	}
	if err := zr.Close(); err != nil {
		return nil, fmt.Errorf("gzip decompression failed: %v: %w", err, ErrInvalidWebhook)
	}
	return out, nil
}

// DecodeSqsPayload reverses the SQS firehose envelope: the message Body
// is base64-decoded and, when the result begins with the gzip magic, it
// is gzip-decompressed. The same call works whether or not Stream is
// currently compressing payloads.
//
// A base64 failure is wrapped in ErrInvalidWebhook with the prefix
// "invalid base64 encoding"; gzip failures propagate from GunzipPayload.
func DecodeSqsPayload(body string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(body)
	if err != nil {
		return nil, fmt.Errorf("invalid base64 encoding: %v: %w", err, ErrInvalidWebhook)
	}
	return GunzipPayload(decoded)
}

// DecodeSnsPayload reverses an SNS HTTP notification envelope: when the
// input is a JSON envelope ({"Type":"Notification","Message":"..."}),
// the inner Message field is extracted and run through the SQS pipeline
// (base64-decode, then gzip-if-magic). When the input is not a JSON
// envelope it is treated as the already-extracted Message string, so
// existing call sites that pre-unwrap continue to work.
//
// Envelope parsing is lenient and never returns an error on its own;
// any error returned here originates from the SQS pipeline and already
// wraps ErrInvalidWebhook.
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

// VerifySignature returns nil when signature equals the hex-encoded
// HMAC-SHA256 of body using secret as the key, and an error wrapping
// ErrInvalidWebhook with the prefix "signature mismatch" otherwise. The
// comparison is constant-time. The signature is always computed over
// the uncompressed JSON bytes, so callers that decoded a gzipped or
// base64-wrapped payload must pass the inflated bytes here.
func VerifySignature(body []byte, signature, secret string) error {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write(body)
	expected := []byte(hex.EncodeToString(mac.Sum(nil)))
	if !hmac.Equal(expected, []byte(signature)) {
		return fmt.Errorf("signature mismatch: %w", ErrInvalidWebhook)
	}
	return nil
}

// ParseEvent decodes the JSON-encoded webhook payload into a typed
// Event. Unknown event types still parse successfully because Event.Type
// is a string alias. JSON decode failures are wrapped in
// ErrInvalidWebhook with the prefix "invalid JSON payload".
func ParseEvent(payload []byte) (*Event, error) {
	var ev Event
	if err := json.Unmarshal(payload, &ev); err != nil {
		return nil, fmt.Errorf("invalid JSON payload: %v: %w", err, ErrInvalidWebhook)
	}
	return &ev, nil
}

func verifyAndParse(payload []byte, signature, secret string) (*Event, error) {
	if err := VerifySignature(payload, signature, secret); err != nil {
		return nil, err
	}
	return ParseEvent(payload)
}

// VerifyAndParseWebhook decompresses body when gzipped, verifies the
// HMAC signature against secret, and returns the parsed Event. Every
// failure path wraps ErrInvalidWebhook, so callers can do
// errors.Is(err, ErrInvalidWebhook) for a unified check.
func VerifyAndParseWebhook(body []byte, signature, secret string) (*Event, error) {
	inflated, err := GunzipPayload(body)
	if err != nil {
		return nil, err
	}
	return verifyAndParse(inflated, signature, secret)
}

// VerifyAndParseSqs decodes the SQS message Body and returns the parsed
// Event. Signature verification is opt-in: pass both signature and
// secret as non-empty strings to run the full HMAC check, or pass both
// as empty strings to skip verification (the default for AWS-transport
// deliveries, where the queue itself is the authentication layer).
// Passing exactly one of the two is treated as a programmer error.
// Every failure path wraps ErrInvalidWebhook.
func VerifyAndParseSqs(messageBody, signature, secret string) (*Event, error) {
	inflated, err := DecodeSqsPayload(messageBody)
	if err != nil {
		return nil, err
	}
	if signature == "" && secret == "" {
		return ParseEvent(inflated)
	}
	if signature == "" || secret == "" {
		return nil, fmt.Errorf("signature and secret must both be provided to verify the SQS/SNS payload: %w", ErrInvalidWebhook)
	}
	return verifyAndParse(inflated, signature, secret)
}

// VerifyAndParseSns decodes the SNS notification Message and returns
// the parsed Event. Signature verification follows the same opt-in
// rules as VerifyAndParseSqs: both empty skips, both set runs the HMAC
// check, exactly one set returns an error wrapping ErrInvalidWebhook.
func VerifyAndParseSns(message, signature, secret string) (*Event, error) {
	inflated, err := DecodeSnsPayload(message)
	if err != nil {
		return nil, err
	}
	if signature == "" && secret == "" {
		return ParseEvent(inflated)
	}
	if signature == "" || secret == "" {
		return nil, fmt.Errorf("signature and secret must both be provided to verify the SQS/SNS payload: %w", ErrInvalidWebhook)
	}
	return verifyAndParse(inflated, signature, secret)
}

// VerifyAndParseWebhook is the client-bound form of the package-level
// helper; it pulls the API secret from the receiver so call sites only
// supply the request body and signature.
func (c *Client) VerifyAndParseWebhook(body []byte, signature string) (*Event, error) {
	return VerifyAndParseWebhook(body, signature, string(c.apiSecret))
}
