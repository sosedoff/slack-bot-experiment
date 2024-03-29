package main

import (
	"fmt"
	"os"
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
	rtm.SetDebug(true)

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

	// Allow direct messages or public messages with the bot prefix
	if !(strings.HasPrefix(event.Channel, "D") || strings.HasPrefix(event.Text, c.userPrefix)) {
		return
	}

	// Remove bot mention and trim whitespace
	text := strings.TrimSpace(strings.Replace(event.Text, c.userPrefix, "", 1))
	text = cleanupText(text)

	// Determine handler for the message
	handler, args := c.config.FindHandler(text)
	if handler == nil {
		fmt.Println("No handler for message:", text)
		c.write(event.Channel, "Sorry, don't know how to handle that.")
		return
	}

	// Shell out to whatever script the user provided
	output := make(chan string)

	go func() {
		for line := range output {
			c.write(event.Channel, line)
		}
	}()

	err := c.runHandler(handler, text, args, output)
	close(output)
	if err != nil {
		c.write(event.Channel, "ERROR:"+err.Error())
		return
	}
}

func (c *Client) runHandler(handler *Handler, text string, args []string, output chan string) error {
	str := handler.Script

	// Replace argument references in the script
	// example: 'ping -c $1' turns into 'ping -c google.com'
	for idx, val := range args {
		key := fmt.Sprintf("$%v", idx+1)
		str = strings.Replace(str, key, val, -1)
	}

	cmd := exec.Command("bash", "-c", str)
	cmd.Stdout = stdoutWriter{Lines: output}
	cmd.Stderr = os.Stderr

	// Feed any user defined environment variables into process
	if len(handler.Env) > 0 {
		env := os.Environ()
		for k, v := range handler.Env {
			env = append(env, k+"="+v)
		}
		cmd.Env = env
	}

	return cmd.Run()
}

func (c *Client) write(channel, message string) {
	c.rtm.SendMessage(&slack.OutgoingMessage{
		Type:    "message",
		Channel: channel,
		Text:    message,
	})
}
