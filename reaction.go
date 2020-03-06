package stream

import (
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
	ExtraData map[string]interface{} `json:"-,extra"` //nolint: staticcheck
}

type reactionResponse struct {
	Message  *Message  `json:"message"`
	Reaction *Reaction `json:"reaction"`
}

type reactionRequest struct {
	Reaction *Reaction `json:"reaction"`
}

// SendReaction sends a reaction to message with given ID
func (ch *Channel) SendReaction(reaction *Reaction, messageID, userID string) (*Message, error) {
	switch {
	case reaction == nil:
		return nil, errors.New("reaction is nil")
	case messageID == "":
		return nil, errors.New("message ID must be not empty")
	case userID == "":
		return nil, errors.New("user ID must be not empty")
	}

	var resp reactionResponse

	reaction.UserID = userID

	p := path.Join("messages", url.PathEscape(messageID), "reaction")

	req := reactionRequest{Reaction: reaction}
	err := ch.client.makeRequest(http.MethodPost, p, nil, req, &resp)

	return resp.Message, err
}

// DeleteReaction removes a reaction from message with given ID
func (ch *Channel) DeleteReaction(messageID, reactionType, userID string) (*Message, error) {
	switch {
	case messageID == "":
		return nil, errors.New("message ID is empty")
	case reactionType == "":
		return nil, errors.New("reaction type is empty")
	case userID == "":
		return nil, ErrorMissingUserID
	}

	p := path.Join("messages", url.PathEscape(messageID), "reaction", url.PathEscape(reactionType))

	params := url.Values{}
	params.Set("user_id", userID)

	var resp reactionResponse

	err := ch.client.makeRequest(http.MethodDelete, p, params, nil, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Message == nil {
		return nil, errors.New("unexpected error: response message nil")
	}

	return resp.Message, nil
}

type reactionsResponse struct {
	Reactions []*Reaction `json:"reactions"`
}

// GetReactions returns list of the reactions for message with given ID.
// options: Pagination params, ie {"limit":{"10"}, "idlte": {"10"}}
func (ch *Channel) GetReactions(messageID string, options map[string][]string) ([]*Reaction, error) {
	if messageID == "" {
		return nil, errors.New("message ID is empty")
	}

	p := path.Join("messages", url.PathEscape(messageID), "reactions")

	var resp reactionsResponse

	err := ch.client.makeRequest(http.MethodGet, p, options, nil, &resp)

	return resp.Reactions, err
}
