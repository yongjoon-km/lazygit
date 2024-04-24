package helpers

import (
	goContext "context"

	"github.com/jesseduffield/lazygit/pkg/gui/types"
)



type TextboxHelper struct {
	c *HelperCommon
}

func NewTextboxHelper(c *HelperCommon) *TextboxHelper {
	return &TextboxHelper{
		c: c,
	}
}

func (self *TextboxHelper) CreatePopupPanel(ctx goContext.Context, opts types.CreatePopupPanelOpts) error {
	self.c.Mutexes().PopupMutex.Lock()
	defer self.c.Mutexes().PopupMutex.Unlock()

	_, cancel := goContext.WithCancel(ctx)

	// we don't allow interruptions of non-loader popups in case we get stuck somehow
	// e.g. a credentials popup never gets its required user input so a process hangs
	// forever.
	// The proper solution is to have a queue of popup options
	currentPopupOpts := self.c.State().GetRepoState().GetCurrentPopupOpts()
	if currentPopupOpts != nil && !currentPopupOpts.HasLoader {
		self.c.Log.Error("ignoring create popup panel because a popup panel is already open")
		cancel()
		return nil
	}

	textboxView := self.c.Views().Textbox
	textboxView.Editable = opts.Editable

	textArea := textboxView.TextArea
	textArea.Clear()
	textArea.TypeString(opts.Prompt)

	self.c.State().GetRepoState().SetCurrentPopupOpts(&opts)

	return self.c.PushContext(self.c.Contexts().Textbox)
}
