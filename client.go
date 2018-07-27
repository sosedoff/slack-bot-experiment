package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/nlopes/slack"
)

type Client struct {
	botId      string
	userPrefix string
	user       *slack.User
	api        *slack.Client
	rtm        *slack.RTM
	config     *Config
}

func NewClient(config *Config) *Client {
	api := slack.New(config.SlackToken)
	rtm := api.NewRTM()

	return &Client{
		api:    api,
		rtm:    rtm,
		config: config,
	}
}

func (c *Client) Start() error {
	resp, err := c.api.AuthTest()
	if err != nil {
		return err
	}

	user, err := c.api.GetUserInfo(resp.UserID)
	if err != nil {
		return err
	}
	c.user = user
	c.botId = user.ID
	c.userPrefix = fmt.Sprintf("<@%s>", c.botId)

	go c.rtm.ManageConnection()
	return c.handleEvents()
}

func (c *Client) handleEvents() error {
	for msg := range c.rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			go c.handleMessage(ev)
		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())
		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
		}
	}

	return nil
}

func (c *Client) handleMessage(event *slack.MessageEvent) {
	// Skip it's own and other bot's messages
	if event.BotID != "" || event.BotID == c.botId {
		return
	}
	// Skip messages without mention of the bot
	if !strings.HasPrefix(event.Text, c.userPrefix) {
		return
	}

	// Remove bot mention and trim whitespace
	text := strings.TrimSpace(strings.Replace(event.Text, c.userPrefix, "", 1))

	// Determine handler for the message
	handler, args := c.config.FindHandler(text)
	if handler == nil {
		fmt.Println("No handler for message:", text)
		return
	}

	// Shell out to whatever script the user provided
	output, err := c.runHandler(handler, text, args)
	if err != nil {
		fmt.Errorf("Run error: %v %s\n", err, output)
		c.write(event.Channel, "ERROR:"+err.Error())
		return
	}

	c.write(event.Channel, string(output))
}

func (c *Client) runHandler(handler *Handler, text string, args []string) ([]byte, error) {
	str := handler.Script

	// Replace argument references in the script
	// example: 'ping -c $1' turns into 'ping -c google.com'
	for idx, val := range args {
		key := fmt.Sprintf("$%v", idx+1)
		str = strings.Replace(str, key, val, -1)
	}

	stdout := bytes.NewBuffer(nil)

	cmd := exec.Command("bash", "-c", str)
	cmd.Stdout = stdout

	// Feed any user defined environment variables into process
	if len(handler.Env) > 0 {
		for k, v := range handler.Env {
			cmd.Env = append(cmd.Env, k+"="+v)
		}
	}

	if err := cmd.Start(); err != nil {
		return stdout.Bytes(), err
	}

	err := cmd.Wait()
	return stdout.Bytes(), err
}

func (c *Client) write(channel, message string) {
	c.rtm.SendMessage(&slack.OutgoingMessage{
		Type:    "message",
		Channel: channel,
		Text:    message,
	})
}
