package cmd

import (
	"testing"

	"github.com/diamondburned/ningen/v3"
	"github.com/stretchr/testify/assert"
)

func Test_unreadStyle(t *testing.T) {
	gt := &guildsTree{}

	styles := []struct {
		input    ningen.UnreadIndication
		wantBold bool
		wantDim  bool
		wantUL   bool
	}{
		{ningen.ChannelRead, false, true, false},
		{ningen.ChannelUnread, true, false, false},
		{ningen.ChannelMentioned, true, false, true},
	}

	for _, s := range styles {
		style := gt.unreadStyle(s.input)
		isBold := style.Bold(false) != style
		isDim := style.Dim(false) != style
		isUL := style.Underline(false) != style

		assert.Equal(t, s.wantBold, isBold, "Style must be bold")
		assert.Equal(t, s.wantDim, isDim, "Style must be dim")
		assert.Equal(t, s.wantUL, isUL, "Style must be underline")
	}
}
