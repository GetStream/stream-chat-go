package stream_chat

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"path"
)

type Reaction struct {
	MessageID string `json:"message_id"`
	UserID    string `json:"user_id"`
	Type      string `json:"type"`

	// any other fields the user wants to attach a reaction
	ExtraData map[string]interface{} `json:"-"`
}

type reactionForJSON Reaction

func (s *Reaction) UnmarshalJSON(data []byte) error {
	var s2 reactionForJSON
	if err := json.Unmarshal(data, &s2); err != nil {
		return err
	}
	*s = Reaction(s2)

	if err := json.Unmarshal(data, &s.ExtraData); err != nil {
		return err
	}

	removeFromMap(s.ExtraData, *s)
	flattenExtraData(s.ExtraData)
	return nil
}

func (s Reaction) MarshalJSON() ([]byte, error) {
	return addToMapAndMarshal(s.ExtraData, reactionForJSON(s))
}

type ReactionResponse struct {
	Message  *Message  `json:"message"`
	Reaction *Reaction `json:"reaction"`
	Response
}

type reactionRequest struct {
	Reaction *Reaction `json:"reaction"`
}

// SendReaction sends a reaction to message with given ID.
// Deprecated: SendReaction is deprecated, use client.SendReaction instead.
func (ch *Channel) SendReaction(ctx context.Context, reaction *Reaction, messageID, userID string) (*ReactionResponse, error) {
	return ch.client.SendReaction(ctx, reaction, messageID, userID)
}

// DeleteReaction removes a reaction from message with given ID.
// Deprecated: DeleteReaction is deprecated, use client.DeleteReaction instead.
func (ch *Channel) DeleteReaction(ctx context.Context, messageID, reactionType, userID string) (*ReactionResponse, error) {
	return ch.client.DeleteReaction(ctx, messageID, reactionType, userID)
}

// SendReaction sends a reaction to message with given ID.
func (c *Client) SendReaction(ctx context.Context, reaction *Reaction, messageID, userID string) (*ReactionResponse, error) {
	switch {
	case reaction == nil:
		return nil, errors.New("reaction is nil")
	case messageID == "":
		return nil, errors.New("message ID must be not empty")
	case userID == "":
		return nil, errors.New("user ID must be not empty")
	}

	reaction.UserID = userID
	p := path.Join("messages", url.PathEscape(messageID), "reaction")

	req := reactionRequest{Reaction: reaction}

	var resp ReactionResponse
	err := c.makeRequest(ctx, http.MethodPost, p, nil, req, &resp)
	return &resp, err
}

// DeleteReaction removes a reaction from message with given ID.
func (c *Client) DeleteReaction(ctx context.Context, messageID, reactionType, userID string) (*ReactionResponse, error) {
	switch {
	case messageID == "":
		return nil, errors.New("message ID is empty")
	case reactionType == "":
		return nil, errors.New("reaction type is empty")
	case userID == "":
		return nil, errors.New("user ID is empty")
	}

	p := path.Join("messages", url.PathEscape(messageID), "reaction", url.PathEscape(reactionType))

	params := url.Values{}
	params.Set("user_id", userID)

	var resp ReactionResponse
	err := c.makeRequest(ctx, http.MethodDelete, p, params, nil, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Message == nil {
		return nil, errors.New("unexpected error: response message nil")
	}

	return &resp, nil
}

type ReactionsResponse struct {
	Reactions []*Reaction `json:"reactions"`
	Response
}

// GetReactions returns list of the reactions for message with given ID.
// options: Pagination params, ie {"limit":{"10"}, "idlte": {"10"}}
func (c *Client) GetReactions(ctx context.Context, messageID string, options map[string][]string) (*ReactionsResponse, error) {
	if messageID == "" {
		return nil, errors.New("message ID is empty")
	}

	p := path.Join("messages", url.PathEscape(messageID), "reactions")

	var resp ReactionsResponse
	err := c.makeRequest(ctx, http.MethodGet, p, options, nil, &resp)
	return &resp, err
}
