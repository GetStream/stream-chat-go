package stream_chat

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
	ExtraData map[string]interface{} `json:"-,extra"`
}

type reactionResponse struct {
	Message  *Message  `json:"message"`
	Reaction *Reaction `json:"reaction"`
}

type reactionRequest struct {
	Reaction *Reaction `json:"reaction"`
}

// SendReaction sends a reaction about a message
// reaction: the reaction object, ie {type: 'love'}
// messageID: is of the message
// userID: the ID of the user that created the reaction
func (ch *Channel) SendReaction(reaction *Reaction, messageID string, userID string) (*Message, error) {
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

// DeleteReaction removes a reaction by user and type
//
// message:  pointer to the message from which we remove the reaction. Message will be updated from response body
// reaction_type: the type of reaction that should be removed
// userID: the id of the user
func (ch *Channel) DeleteReaction(messageID string, reactionType string, userID string) (*Message, error) {
	switch {
	case messageID == "":
		return nil, errors.New("message ID is empty")
	case reactionType == "":
		return nil, errors.New("reaction type is empty")
	case userID == "":
		return nil, errors.New("user ID is empty")
	}

	p := path.Join("messages", url.PathEscape(messageID), "reaction", url.PathEscape(reactionType))

	params := map[string][]string{
		"user_id": {userID},
	}

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
