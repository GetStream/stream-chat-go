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
	Message  Message  `json:"message"`
	Reaction Reaction `json:"reaction"`
}

type reactionRequest struct {
	Reaction *Reaction `json:"reaction"`
}

// SendReaction sends a reaction about a message
//
// message: pointer to the message struct
// reaction: the reaction object, ie {type: 'love'}
// userID: the ID of the user that created the reaction
func (ch *Channel) SendReaction(msg *Message, reaction *Reaction, userID string) error {
	switch {
	case msg == nil:
		return errors.New("message is nil")
	case reaction == nil:
		return errors.New("reaction is nil")
	case msg.ID == "":
		return errors.New("message ID must be not empty")
	case userID == "":
		return errors.New("user ID must be not empty")
	}

	var resp reactionResponse

	reaction.UserID = userID

	p := path.Join("messages", url.PathEscape(msg.ID), "reaction")

	req := reactionRequest{Reaction: reaction}
	err := ch.client.makeRequest(http.MethodPost, p, nil, req, &resp)

	*msg = resp.Message
	*reaction = resp.Reaction

	return err
}

// DeleteReaction removes a reaction by user and type
//
// message:  pointer to the message from which we remove the reaction. Message will be updated from response body
// reaction_type: the type of reaction that should be removed
// userID: the id of the user
func (ch *Channel) DeleteReaction(message *Message, reactionType string, userID string) error {
	switch {
	case message == nil:
		return errors.New("message is nil")
	case reactionType == "":
		return errors.New("reaction type must be not empty")
	case message.ID == "":
		return errors.New("message ID must be not empty")
	case userID == "":
		return errors.New("user ID must be not empty")
	}

	p := path.Join("messages", url.PathEscape(message.ID), "reaction", url.PathEscape(reactionType))

	params := map[string][]string{
		"user_id": {userID},
	}

	var resp reactionResponse

	err := ch.client.makeRequest(http.MethodDelete, p, params, nil, &resp)
	if err != nil {
		return err
	}

	*message = resp.Message

	return nil
}
