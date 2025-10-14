package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockFlex struct {
	items        []string
	removedItems []string
}

func (f *mockFlex) GetItemCount() int {
	return len(f.items)
}

func (f *mockFlex) RemoveItem(item interface{}) {
	itemString, boolean := item.(string)
	if boolean {
		for i := 0; i < len(f.items); i++ {
			if f.items[i] == itemString {
				f.items = append(f.items[:i], f.items[i+1:]...)
				break
			}
		}

		f.removedItems = append(f.removedItems, itemString)
	}
}

type mockGuildsTree struct {
	focus bool
}

func (g *mockGuildsTree) HasFocus() bool {
	return g.focus
}

// Mock application
type mockApp struct {
	flex       *mockFlex
	guildsTree *mockGuildsTree
	focusSet   string
	initCalled bool
}

func (a *mockApp) SetFocus(name interface{}) {
	switch name.(type) {
	case *mockFlex:
		a.focusSet = "dummy-flex"
	case *mockGuildsTree:
		a.focusSet = "fake-guild"
	}
}

func (a *mockApp) init() {
	a.initCalled = true
}

func (a *mockApp) toggleGuildsTree(removeItem string) {
	// The guilds tree is visible if the number of items is two.
	if a.flex.GetItemCount() == 2 {
		a.flex.RemoveItem(removeItem)
		if a.guildsTree.HasFocus() {
			a.SetFocus(a.flex)
		}
	} else {
		a.init()
		a.SetFocus(a.guildsTree)
	}
}

func Test_toggleGuildsTree_RemoveGuild(t *testing.T) {
	removeItem := "fake-guild"
	mockApplication := &mockApp{
		flex:       &mockFlex{items: []string{"dummy-flex", removeItem}},
		guildsTree: &mockGuildsTree{focus: true},
	}

	mockApplication.toggleGuildsTree(removeItem)
	assert.Contains(t, mockApplication.flex.removedItems, removeItem, "removed items must contain fake-guild")
	assert.Len(t, mockApplication.flex.items, 1, "expected 1 item after removal")
	assert.Equal(t, "dummy-flex", mockApplication.focusSet, "expected focus to be on dummy-flex")
}

func Test_toggleGuildsTree_AddGuild(t *testing.T) {
	addItem := "fake-guild"
	mockApplication := &mockApp{
		flex:       &mockFlex{items: []string{"dummy-flex"}},
		guildsTree: &mockGuildsTree{focus: true},
	}

	mockApplication.toggleGuildsTree(addItem)
	assert.True(t, mockApplication.initCalled, "init must be called when adding fake-guild")
	assert.Equal(t, "fake-guild", mockApplication.focusSet, "focus must be set to fake-guild")
}
