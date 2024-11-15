package mainBuildModel

import (
	buildConfig "ulld/cli/internal/build/config"
	"ulld/cli/internal/build/constants"
	keyMap "ulld/cli/internal/build/keymap"
	"ulld/cli/internal/build/ui/confirmdir"
	"ulld/cli/internal/build/ui/filepicker"
	"ulld/cli/internal/keymap"
	"ulld/cli/internal/signals"
	fs_utils "ulld/cli/internal/utils/fs"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/mitchellh/go-homedir"
)

type mainModel struct {
	stage           constants.BuildStage
	help            help.Model
	keys            keyMap.KeyMap
	confirmDirModel confirmdir.Model
	targetDirModel  filepicker.Model
	targetDir       string
	quitting        bool
}

func (m mainModel) Init() tea.Cmd {
	tea.SetWindowTitle("ULLD Build")
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.targetDirModel, cmd = m.targetDirModel.Update(msg)
		cmds = append(cmds, cmd)
		m.confirmDirModel, cmd = m.confirmDirModel.Update(msg)
		cmds = append(cmds, cmd)
	case signals.SetStageMsg:
		m.stage = msg.NewStage
	case signals.SetUseSelectedDirMsg:
		if msg.UseSelectedDir {
			m.stage = constants.CloneTemplateAppStage
		} else {
			m.stage = constants.PickTargetDirStage
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keymap.Keymap.Quit):
			m.quitting = true
			return m, tea.Quit
		}
	}
	switch m.stage {
	case m.confirmDirModel.Stage:
		m.confirmDirModel, cmd = m.confirmDirModel.Update(msg)
		cmds = append(cmds, cmd)
	case m.targetDirModel.Stage:
		m.targetDirModel, cmd = m.targetDirModel.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	var s string
	if m.quitting {
		return "\n  No worries.\n\n"
	}
	switch m.stage {
	case constants.ConfirmCurrentDirStage:
		return m.confirmDirModel.View()
	case constants.PickTargetDirStage:
		return m.targetDirModel.View()
	}
	return s
}

func InitialMainModel(cfg *buildConfig.BuildConfigOpts) *mainModel {
	homeDir, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
	val := mainModel{
		stage:           cfg.InitialStage,
		help:            help.New(),
		keys:            keyMap.DefaultKeymap,
		targetDirModel:  filepicker.NewModel(homeDir, fs_utils.DirOnlyDataType, "Where would you like to build ULLD?"),
		confirmDirModel: confirmdir.NewModel("Do you want to build ULLD in your current directory?"),
		targetDir:       cfg.TargetDir,
		quitting:        false,
	}

	return &val
}
