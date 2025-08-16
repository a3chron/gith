package internal

import (
	"github.com/a3chron/gith/internal/config"
	"github.com/charmbracelet/bubbles/spinner"
)

type Step int

const (
	StepLoad Step = iota
	StepAction

	StepBranchAction
	StepBranchSelect
	StepBranchCreate
	StepBranchInput

	StepCommitAction
	StepCommitSelectPrefix
	StepCommitInput

	StepTag
	StepTagSelect
	StepTagAdd
	StepTagInput

	StepRemote
	StepRemoteSelect
	StepRemoteNameInput
	StepRemoteUrlInput

	StepChanges

	StepOptions
	StepOptionsFlavorSelect
	StepOptionsAccentSelect
)

type ActionModel struct {
	Actions        []string
	SelectedAction string
}

type BranchModel struct {
	Actions        []string
	SelectedAction string
	Branches       []string
	SelectedBranch string
	Options        []string
	SelectedOption string
	Input          string
}

type CommitModel struct {
	Actions        []string
	SelectedAction string
	CommitPrefixes []string
	SelectedPrefix string
	CommitMessage  string
}

type TagModel struct {
	Actions        []string
	SelectedAction string
	Options        []string
	SelectedOption string
	AddOptions     []string
	SelectedAddTag string
	CurrentTag     string
	ManualInput    string
}

type RemoteModel struct {
	Actions        []string
	SelectedAction string
	Options        []string
	SelectedOption string
	NameInput      string
	UrlInput       string
}

type ConfigModel struct {
	Actions        []string
	SelectedAction string
	Flavors        []string
	SelectedFlavor string
	Accents        []string
	SelectedAccent string
}

type Model struct {
	CurrentStep   Step
	Loading       bool
	Selected      int
	ActionModel   ActionModel
	BranchModel   BranchModel
	CommitModel   CommitModel
	RemoteModel   RemoteModel
	TagModel      TagModel
	ConfigModel   ConfigModel
	CurrentConfig *config.Config
	Spinner       spinner.Model
	Level         int
	Output        []string
	Err           string
	Success       string
	StartAt       string
	StartAtLevel  int
}

type RepoUpdatedMsg struct{}
type RepoUpdateErrorMsg struct{}
