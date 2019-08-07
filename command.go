package stream_chat

import "encoding/json"

type Command struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Args        string `json:"args"`
	Set         string `json:"set"`
}

type Commands []Command

// response of commands can be one of []string, []Command
func (c *Commands) UnmarshalJSON(data []byte) error {
	var cmds []interface{}
	err := json.Unmarshal(data, &cmds)
	if err != nil {
		return err
	}

	for i := range cmds {
		var cmd Command
		switch t := cmds[i].(type) {

		case map[string]interface{}:
			cmd.Name = t["name"].(string)
			cmd.Set = t["set"].(string)
			cmd.Args = t["args"].(string)
			cmd.Description = t["description"].(string)

		case string:
			cmd.Name = t
		}
		*c = append(*c, cmd)
	}
	return nil
}

// commands are sent as []string
func (c Commands) MarshalJSON() ([]byte, error) {
	var cmds = make([]string, 0, len(c))
	for i := range c {
		cmds = append(cmds, c[i].Name)
	}
	if len(cmds) == 0 {
		cmds = []string{"all"}
	}

	return json.Marshal(cmds)
}
