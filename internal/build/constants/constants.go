package constants

import tea "github.com/charmbracelet/bubbletea"

type BuildStage int

// FIXME: This is wrong. Need to find how to assign these types to the tea.msg type properly.
type ToRootModelMsg tea.Msg

type ConfirmDirectoryMsg tea.Msg

var ToRootModelCmd tea.Msg = tea.Msg("TRM")

const (
	ConfirmCurrentDirStage BuildStage = iota
	PickTargetDirStage
	CloneTemplateAppStage
)