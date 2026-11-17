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

func hmacHex(t *testing.T, secret, body []byte) string {
	t.Helper()
	mac := hmac.New(sha256.New, secret)
	_, err := mac.Write(body)
	require.NoError(t, err)
	return hex.EncodeToString(mac.Sum(nil))
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

func base64String(t *testing.T, src []byte) string {
	t.Helper()
	return base64.StdEncoding.EncodeToString(src)
}

func TestVerifyWebhook_BackwardCompatibility(t *testing.T) {
	c := newWebhookTestClient(t)
	body := []byte(webhookTestFixture)
	sig := []byte(hmacHex(t, []byte(webhookTestAPISecret), body))

	require.True(t, c.VerifyWebhook(body, sig), "valid signature must verify")
	require.False(t, c.VerifyWebhook(body, []byte("not-a-valid-hex-signature")))
	require.False(t, c.VerifyWebhook([]byte("tampered"), sig))
}

func TestGunzipPayload(t *testing.T) {
	body := []byte(webhookTestFixture)

	t.Run("passthrough plain bytes", func(t *testing.T) {
		got, err := GunzipPayload(body)
		require.NoError(t, err)
		require.Equal(t, body, got)
	})

	t.Run("inflates gzip bytes", func(t *testing.T) {
		got, err := GunzipPayload(gzipBytes(t, body))
		require.NoError(t, err)
		require.Equal(t, body, got)
	})

	t.Run("empty input returns empty", func(t *testing.T) {
		got, err := GunzipPayload([]byte{})
		require.NoError(t, err)
		require.Equal(t, []byte{}, got)
	})

	t.Run("short input below magic length", func(t *testing.T) {
		got, err := GunzipPayload([]byte("ab"))
		require.NoError(t, err)
		require.Equal(t, []byte("ab"), got)
	})

	t.Run("gunzipPayload returns ErrInvalidWebhook on corrupt gzip", func(t *testing.T) {
		bad := append(append([]byte{}, gzipMagic...), 0, 0, 0)
		got, err := GunzipPayload(bad)
		require.Error(t, err)
		require.Nil(t, got)
		require.True(t, errors.Is(err, ErrInvalidWebhook))
		assert.Contains(t, err.Error(), "gzip decompression failed")
	})

	t.Run("decompresses helloworld fixture", func(t *testing.T) {
		compressed, err := base64.StdEncoding.DecodeString("H4sIAGrYAWoAA8tIzcnJL88vykkBAK0g6/kKAAAA")
		require.NoError(t, err)
		got, err := GunzipPayload(compressed)
		require.NoError(t, err)
		require.Equal(t, []byte("helloworld"), got)
	})
}

func TestGunzipPayloadLimit(t *testing.T) {
	limit := MaxWebhookPayloadBytes
	require.Equal(t, int64(10*1024*1024), limit)

	t.Run("inflates body at limit", func(t *testing.T) {
		got, err := GunzipPayload(gzipBytes(t, make([]byte, limit)))
		require.NoError(t, err)
		require.Len(t, got, int(limit))
	})

	t.Run("rejects gzip bomb before signature check", func(t *testing.T) {
		bomb := gzipBytes(t, make([]byte, limit+1))
		want := "gzip decompression failed: payload exceeds 10485760 bytes: invalid webhook"

		got, err := VerifyAndParseWebhook(bomb, "deadbeef", webhookTestAPISecret)
		require.Nil(t, got)
		require.ErrorIs(t, err, ErrInvalidWebhook)
		require.EqualError(t, err, want)

		got, err = ParseSqs(base64String(t, bomb))
		require.Nil(t, got)
		require.ErrorIs(t, err, ErrInvalidWebhook)
		require.EqualError(t, err, want)
	})

	t.Run("does not limit uncompressed body", func(t *testing.T) {
		body := make([]byte, limit+1)
		got, err := GunzipPayload(body)
		require.NoError(t, err)
		require.Len(t, got, len(body))
	})

	t.Run("inflates multiple members", func(t *testing.T) {
		got, err := GunzipPayload(append(gzipBytes(t, []byte("ab")), gzipBytes(t, []byte("cd"))...))
		require.NoError(t, err)
		require.Equal(t, []byte("abcd"), got)
	})

	t.Run("corrupt deflate data still errors", func(t *testing.T) {
		bad := gzipBytes(t, []byte(webhookTestFixture))
		copy(bad[10:14], []byte{0xff, 0xff, 0xff, 0xff})
		_, err := GunzipPayload(bad)
		require.ErrorIs(t, err, ErrInvalidWebhook)
		require.EqualError(t, err, "gzip decompression failed: invalid webhook")
	})

	t.Run("honours custom limit", func(t *testing.T) {
		t.Cleanup(func() { MaxWebhookPayloadBytes = limit })
		body := []byte(webhookTestFixture)

		MaxWebhookPayloadBytes = int64(len(body) - 1)
		_, err := GunzipPayload(gzipBytes(t, body))
		require.ErrorIs(t, err, ErrInvalidWebhook)
		assert.Contains(t, err.Error(), "payload exceeds")

		MaxWebhookPayloadBytes = int64(len(body))
		got, err := GunzipPayload(gzipBytes(t, body))
		require.NoError(t, err)
		require.Equal(t, body, got)

		MaxWebhookPayloadBytes = 1
		got, err = GunzipPayload(body)
		require.NoError(t, err)
		require.Equal(t, body, got)
	})
}

func TestDecodeSqsPayload(t *testing.T) {
	body := []byte(webhookTestFixture)

	t.Run("base64 only - no compression", func(t *testing.T) {
		got, err := DecodeSqsPayload(base64String(t, body))
		require.NoError(t, err)
		require.Equal(t, body, got)
	})

	t.Run("base64 plus gzip", func(t *testing.T) {
		got, err := DecodeSqsPayload(base64String(t, gzipBytes(t, body)))
		require.NoError(t, err)
		require.Equal(t, body, got)
	})

	t.Run("decodeSqsPayload returns ErrInvalidWebhook on invalid base64", func(t *testing.T) {
		got, err := DecodeSqsPayload("!!!not-base64!!!")
		require.Error(t, err)
		require.Nil(t, got)
		require.True(t, errors.Is(err, ErrInvalidWebhook))
		assert.Contains(t, err.Error(), "invalid base64 encoding")
	})

	t.Run("helloworld base64 fixture", func(t *testing.T) {
		got, err := DecodeSqsPayload("aGVsbG93b3JsZA==")
		require.NoError(t, err)
		require.Equal(t, []byte("helloworld"), got)
	})

	t.Run("helloworld base64+gzip fixture", func(t *testing.T) {
		got, err := DecodeSqsPayload("H4sIAGrYAWoAA8tIzcnJL88vykkBAK0g6/kKAAAA")
		require.NoError(t, err)
		require.Equal(t, []byte("helloworld"), got)
	})
}

func TestDecodeSnsPayload(t *testing.T) {
	body := []byte(webhookTestFixture)
	wrapped := base64String(t, gzipBytes(t, body))

	t.Run("pre-extracted message backward compat", func(t *testing.T) {
		sns, err := DecodeSnsPayload(wrapped)
		require.NoError(t, err)
		sqs, err := DecodeSqsPayload(wrapped)
		require.NoError(t, err)
		require.Equal(t, sqs, sns)
		require.Equal(t, body, sns)
	})

	t.Run("full SNS HTTP notification envelope", func(t *testing.T) {
		envelope := snsEnvelope(t, wrapped)
		got, err := DecodeSnsPayload(envelope)
		require.NoError(t, err)
		require.Equal(t, body, got)
	})

	t.Run("envelope with whitespace prefix", func(t *testing.T) {
		envelope := "  \n" + snsEnvelope(t, wrapped)
		got, err := DecodeSnsPayload(envelope)
		require.NoError(t, err)
		require.Equal(t, body, got)
	})
}

// snsEnvelope returns a realistic SNS HTTP POST notification body that
// wraps payload as the Message field. Matches the documented SNS schema.
func snsEnvelope(t *testing.T, payload string) string {
	t.Helper()
	env := map[string]any{
		"Type":             "Notification",
		"MessageId":        "22b80b92-fdea-4c2c-8f9d-bdfb0c7bf324",
		"TopicArn":         "arn:aws:sns:us-east-1:123456789012:stream-webhooks",
		"Message":          payload,
		"Timestamp":        "2026-05-11T10:00:00.000Z",
		"SignatureVersion": "1",
		"MessageAttributes": map[string]any{
			"X-Signature": map[string]string{
				"Type":  "String",
				"Value": "<signature placeholder>",
			},
		},
	}
	out, err := json.Marshal(env)
	require.NoError(t, err)
	return string(out)
}

func TestVerifySignature(t *testing.T) {
	body := []byte(webhookTestFixture)
	sig := hmacHex(t, []byte(webhookTestAPISecret), body)

	t.Run("valid signature returns nil", func(t *testing.T) {
		require.NoError(t, VerifySignature(body, sig, webhookTestAPISecret))
	})

	t.Run("verifySignature returns ErrInvalidWebhook on mismatch", func(t *testing.T) {
		err := VerifySignature(body, "0000000000000000000000000000000000000000000000000000000000000000", webhookTestAPISecret)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrInvalidWebhook))
		assert.Contains(t, err.Error(), "signature mismatch")
	})

	t.Run("wrong secret rejected", func(t *testing.T) {
		err := VerifySignature(body, sig, "different-secret")
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrInvalidWebhook))
		assert.Contains(t, err.Error(), "signature mismatch")
	})

	t.Run("signature over compressed bytes rejected", func(t *testing.T) {
		compressed := gzipBytes(t, body)
		sigOverCompressed := hmacHex(t, []byte(webhookTestAPISecret), compressed)
		err := VerifySignature(body, sigOverCompressed, webhookTestAPISecret)
		require.Error(t, err, "signature must be computed over uncompressed bytes")
		require.True(t, errors.Is(err, ErrInvalidWebhook))
	})
}

func TestParseEvent(t *testing.T) {
	t.Run("known event type", func(t *testing.T) {
		got, err := ParseEvent([]byte(webhookTestFixture))
		require.NoError(t, err)
		require.Equal(t, EventMessageNew, got.Type)
		require.NotNil(t, got.Message)
		require.Equal(t, "the quick brown fox", got.Message.Text)
	})

	t.Run("unknown event type still parses", func(t *testing.T) {
		got, err := ParseEvent([]byte(`{"type":"a.future.event","custom":42}`))
		require.NoError(t, err)
		require.Equal(t, EventType("a.future.event"), got.Type)
	})

	t.Run("parseEvent returns ErrInvalidWebhook on invalid JSON", func(t *testing.T) {
		got, err := ParseEvent([]byte("not json"))
		require.Error(t, err)
		require.Nil(t, got)
		require.True(t, errors.Is(err, ErrInvalidWebhook))
		assert.Contains(t, err.Error(), "invalid JSON payload")
	})
}

func TestVerifyAndParseWebhook(t *testing.T) {
	c := newWebhookTestClient(t)
	body := []byte(webhookTestFixture)
	sig := hmacHex(t, []byte(webhookTestAPISecret), body)

	t.Run("plain body via package", func(t *testing.T) {
		got, err := VerifyAndParseWebhook(body, sig, webhookTestAPISecret)
		require.NoError(t, err)
		require.Equal(t, EventMessageNew, got.Type)
	})

	t.Run("gzip body via package", func(t *testing.T) {
		got, err := VerifyAndParseWebhook(gzipBytes(t, body), sig, webhookTestAPISecret)
		require.NoError(t, err)
		require.Equal(t, EventMessageNew, got.Type)
	})

	t.Run("plain body via client", func(t *testing.T) {
		got, err := c.VerifyAndParseWebhook(body, sig)
		require.NoError(t, err)
		require.Equal(t, EventMessageNew, got.Type)
	})

	t.Run("gzip body via client", func(t *testing.T) {
		got, err := c.VerifyAndParseWebhook(gzipBytes(t, body), sig)
		require.NoError(t, err)
		require.Equal(t, EventMessageNew, got.Type)
	})

	t.Run("signature mismatch returns ErrInvalidWebhook", func(t *testing.T) {
		got, err := c.VerifyAndParseWebhook(body, strings.Repeat("0", 64))
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrInvalidWebhook))
		assert.Contains(t, err.Error(), "signature mismatch")
		require.Nil(t, got)
	})

	t.Run("signature over compressed bytes rejected", func(t *testing.T) {
		compressed := gzipBytes(t, body)
		sigOverCompressed := hmacHex(t, []byte(webhookTestAPISecret), compressed)
		got, err := c.VerifyAndParseWebhook(compressed, sigOverCompressed)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrInvalidWebhook))
		assert.Contains(t, err.Error(), "signature mismatch")
		require.Nil(t, got)
	})

	t.Run("propagates decompression error as ErrInvalidWebhook", func(t *testing.T) {
		bogus := append(append([]byte{}, gzipMagic...), []byte("garbage")...)
		got, err := c.VerifyAndParseWebhook(bogus, sig)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrInvalidWebhook))
		assert.Contains(t, err.Error(), "gzip decompression failed")
		require.Nil(t, got)
	})
}

func TestParseSqs(t *testing.T) {
	c := newWebhookTestClient(t)
	body := []byte(webhookTestFixture)

	t.Run("base64 only via client", func(t *testing.T) {
		got, err := c.ParseSqs(base64String(t, body))
		require.NoError(t, err)
		require.Equal(t, EventMessageNew, got.Type)
		require.NotNil(t, got.Message)
		require.Equal(t, "the quick brown fox", got.Message.Text)
	})

	t.Run("base64 plus gzip via client", func(t *testing.T) {
		wrapped := base64String(t, gzipBytes(t, body))
		got, err := c.ParseSqs(wrapped)
		require.NoError(t, err)
		require.Equal(t, EventMessageNew, got.Type)
	})

	t.Run("via package", func(t *testing.T) {
		wrapped := base64String(t, gzipBytes(t, body))
		got, err := ParseSqs(wrapped)
		require.NoError(t, err)
		require.Equal(t, EventMessageNew, got.Type)
	})

	t.Run("invalid base64 surfaced as ErrInvalidWebhook", func(t *testing.T) {
		got, err := c.ParseSqs("!!!not-base64!!!")
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrInvalidWebhook))
		assert.Contains(t, err.Error(), "invalid base64 encoding")
		require.Nil(t, got)
	})
}

func TestParseSns(t *testing.T) {
	c := newWebhookTestClient(t)
	body := []byte(webhookTestFixture)
	wrapped := base64String(t, gzipBytes(t, body))

	t.Run("pre-extracted message backward compat", func(t *testing.T) {
		got, err := c.ParseSns(wrapped)
		require.NoError(t, err)
		require.Equal(t, EventMessageNew, got.Type)

		pkg, err := ParseSns(wrapped)
		require.NoError(t, err)
		require.Equal(t, got.Type, pkg.Type)
	})

	t.Run("full SNS HTTP notification envelope", func(t *testing.T) {
		envelope := snsEnvelope(t, wrapped)
		got, err := c.ParseSns(envelope)
		require.NoError(t, err)
		require.Equal(t, EventMessageNew, got.Type)
		require.NotNil(t, got.Message)
		require.Equal(t, "the quick brown fox", got.Message.Text)
	})

	t.Run("invalid base64 surfaced as ErrInvalidWebhook", func(t *testing.T) {
		got, err := c.ParseSns("!!!not-base64!!!")
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrInvalidWebhook))
		assert.Contains(t, err.Error(), "invalid base64 encoding")
		require.Nil(t, got)
	})
}
