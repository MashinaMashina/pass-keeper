package cliselect

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/manifoldco/promptui/list"
	"io"
	"os"
	"reflect"
)

const ModeSimple = 1
const ModeInteractive = 2

type Select struct {
	Label             interface{}
	Items             interface{}
	Size              int
	CursorPos         int
	IsVimMode         bool
	HideHelp          bool
	HideSelected      bool
	Templates         *promptui.SelectTemplates
	Keys              *promptui.SelectKeys
	Searcher          list.Searcher
	StartInSearchMode bool
	list              *list.List
	Pointer           promptui.Pointer
	Stdin             io.ReadCloser
	Stdout            io.WriteCloser

	Mode uint8
}

func (s *Select) Run() (int, string, error) {
	if s.Mode == ModeInteractive {
		sel := promptui.Select{
			Label:             s.Label,
			Items:             s.Items,
			Size:              s.Size,
			CursorPos:         s.CursorPos,
			IsVimMode:         s.IsVimMode,
			HideHelp:          s.HideHelp,
			HideSelected:      s.HideSelected,
			Templates:         s.Templates,
			Keys:              s.Keys,
			Searcher:          s.Searcher,
			StartInSearchMode: s.StartInSearchMode,
			Pointer:           s.Pointer,
			Stdin:             s.Stdin,
			Stdout:            s.Stdout,
		}

		return sel.Run()
	} else {
		return s.simpleRun()
	}
}

func (s *Select) simpleRun() (int, string, error) {
	if s.Stdout == nil {
		s.Stdout = os.Stdout
	}

	if s.Stdin == nil {
		s.Stdin = os.Stdin
	}

	if s.Items == nil || reflect.TypeOf(s.Items).Kind() != reflect.Slice {
		return 0, "", fmt.Errorf("items %v is not a slice", s.Items)
	}

	slice := reflect.ValueOf(s.Items)
	values := make([]*interface{}, slice.Len())

	for i := range values {
		item := slice.Index(i).Interface()
		values[i] = &item
	}

	if s.Label != "" {
		_, err := fmt.Fprintln(s.Stdout, s.Label)
		if err != nil {
			return 0, "", err
		}
	}

	for i, val := range values {
		fmt.Fprint(s.Stdout, i+1, ". ", *val, "\n")
	}

	_, err := fmt.Fprint(s.Stdout, "Enter number of the item: ")
	if err != nil {
		return 0, "", err
	}

	var res int
	fmt.Fscanln(s.Stdin, &res)

	if res == 0 || res > len(values) {
		return 0, "", nil
	}

	return res, fmt.Sprintf("%s", *values[res-1]), nil
}
