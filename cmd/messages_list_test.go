package cmd

import (
	"fmt"
	"testing"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/stretchr/testify/assert"

	"github.com/ayn2op/discordo/internal/config"
	"github.com/ayn2op/tview"
)

func Test_extractURLs(t *testing.T) {
	textContent := "Testing links: https://dummy.com and also [Test](http://test.com). Also broken url: [Broken](broken)"
	urls := extractURLs(textContent)
	assert.Len(t, urls, 3, "There should be 3 URLs extracted")
	assert.NotContains(t, urls, "Testing links:")
	assert.Contains(t, urls, "https://dummy.com")
	assert.Contains(t, urls, "http://test.com")
	assert.Contains(t, urls, "broken")
}

var ml messagesList

func Test_setFetchingChunk(t *testing.T) {
	// Start fetch
	ml.setFetchingChunk(true, 0)
	assert.True(t, ml.fetchingMembers.value)
	assert.NotNil(t, ml.fetchingMembers.done)
	assert.Zero(t, ml.fetchingMembers.count)

	// Stop fetch
	ml.setFetchingChunk(false, 123)
	assert.False(t, ml.fetchingMembers.value)
	assert.Equal(t, uint(123), ml.fetchingMembers.count)
}

func Test_waitForChunkEvent(t *testing.T) {
	ml.setFetchingChunk(true, 0)

	go func() {
		time.Sleep(30 * time.Millisecond)
		ml.setFetchingChunk(false, 42)
	}()

	count := ml.waitForChunkEvent()
	assert.NotNil(t, ml.fetchingMembers.done, "Must not be nil after waitForChunkEvent")
	assert.Equal(t, uint(42), count, "waitForChunkEvent must return 42/should wait asynchronous call")
}

func Test_waitForChunkEvent_FalseValue(t *testing.T) {
	// Ensure fetchingMembers.value is false
	ml.fetchingMembers.value = false

	count := ml.waitForChunkEvent()
	if count != 0 {
		t.Errorf("Expected 0, got %d", count)
	}
}

func Test_newMessageList(t *testing.T) {
	ml := newTestMessagesList()
	assert.NotNil(t, ml)
	assert.NotNil(t, ml.TextView)
	assert.Equal(t, "Messages", ml.GetTitle())
	assert.NotNil(t, ml.GetInputCapture())
	assert.NotNil(t, ml.onHighlighted)
}

func Test_Reset(t *testing.T) {
	ml := newTestMessagesList()
	ml.SetTitle("Test title")
	ml.Highlight("msg")

	ml.reset()

	assert.Equal(t, discord.MessageID(0), ml.selectedMessageID)
	assert.Equal(t, "", ml.GetTitle())
	assert.Empty(t, ml.GetHighlights())
}

func Test_SetTitle(t *testing.T) {
	app = &application{
		messagesList: newTestMessagesList(),
	}
	ml := app.messagesList

	channel := discord.Channel{Name: "Dummy", Topic: "Test"}

	ml.setTitle(channel)
	assert.Equal(t, "#Dummy - Test", ml.GetTitle())

	channel.Topic = ""
	ml.setTitle(channel)
	assert.Equal(t, "#Dummy", ml.GetTitle())
}

func Test_FormatTimestamp(t *testing.T) {
	ml := newTestMessagesList()
	ml.cfg.Timestamps.Format = "15:04"
	thisTimestamp := "12:34"
	ts := discord.Timestamp(time.Date(2025, 10, 14, 12, 34, 0, 0, time.Local))
	assert.Equal(t, thisTimestamp, ml.formatTimestamp(ts))
}

func Test_DrawTimestamps(t *testing.T) {
	ml := newTestMessagesList()
	ml.cfg.Timestamps.Format = "15:04"
	ts := discord.Timestamp(time.Date(2025, 10, 14, 12, 34, 0, 0, time.Local))
	thisTimestamp := "12:34"

	ml.drawTimestamps(ts)
	text := ml.GetText(true)
	assert.Contains(t, text, thisTimestamp)
}

func Test_DrawAuthor(t *testing.T) {
	author := "dummy"
	ml := &messagesList{TextView: tview.NewTextView()}
	msg := discord.Message{
		Author:  discord.User{Username: author, ID: 1},
		GuildID: 0,
	}
	ml.drawAuthor(msg)

	text := ml.GetText(true)
	assert.Contains(t, text, author)
}

func Test_DrawContent(t *testing.T) {
	app = &application{
		messagesList: newTestMessagesList(),
		cfg:          &config.Config{Markdown: false},
	}
	ml := app.messagesList

	content := "Hello World"
	msg := discord.Message{Content: content}

	ml.drawContent(msg)
	assert.Contains(t, ml.GetText(true), content)
}

func Test_SelectedMsg_NoSelection(t *testing.T) {
	ml := newTestMessagesList()
	_, err := ml.selectedMsg()
	assert.Error(t, err)
}

func Test_OnHighlighted(t *testing.T) {
	ml := newTestMessagesList()
	ml.onHighlighted([]string{"13"}, nil, nil)
	assert.Equal(t, discord.MessageID(13), ml.selectedMessageID)

	// invalid input
	ml.selectedMessageID = 0
	ml.onHighlighted([]string{"not integer"}, nil, nil)
	assert.Equal(t, discord.MessageID(0), ml.selectedMessageID)
}

func newTestMessagesList() *messagesList {
	return newMessagesList(&config.Config{
		Timestamps: config.Timestamps{Enabled: false},
		Theme:      config.Theme{},
	})
}

func Test_DrawSnapshotContent(t *testing.T) {
	ml := newTestMessagesList()
	content := "Hello world"
	msg := discord.MessageSnapshotMessage{Content: content}

	ml.drawSnapshotContent(msg)
	assert.Contains(t, ml.GetText(true), content)
}

func Test_DrawDefaultMessage_WithAttachment(t *testing.T) {
	ml := newTestMessagesList()
	ml.cfg.ShowAttachmentLinks = true
	content := "hello world"
	testFile := "file.txt"
	testUrl := "http://example.com/file.txt"

	msg := discord.Message{
		Content: content,
		Attachments: []discord.Attachment{
			{Filename: testFile, URL: testUrl},
		},
	}

	ml.drawDefaultMessage(msg)
	text := ml.GetText(true)
	assert.Contains(t, text, content)
	assert.Contains(t, text, testFile)
	assert.Contains(t, text, testUrl)
}

func Test_DrawForwardedMessage(t *testing.T) {
	ml := newTestMessagesList()
	frw_indicator := "fwd: "
	content := "Forwarded Content"
	ml.cfg.Theme.MessagesList.ForwardedIndicator = frw_indicator

	msg := discord.Message{
		MessageSnapshots: []discord.MessageSnapshot{
			{Message: discord.MessageSnapshotMessage{Content: content}},
		},
	}

	ml.drawForwardedMessage(msg)
	text := ml.GetText(true)
	assert.Contains(t, text, content)
	assert.Contains(t, text, frw_indicator)
}

func Test_DrawPinnedMessage(t *testing.T) {
	author := "dummy"
	ml := newTestMessagesList()
	msg := discord.Message{
		Author: discord.User{Username: author},
	}

	ml.drawPinnedMessage(msg)
	text := ml.GetText(true)
	assert.Contains(t, text, fmt.Sprintf("%s pinned a message", author))
}
