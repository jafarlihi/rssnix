package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

const (
	bash = `#! /bin/bash

: ${PROG:=$(basename ${BASH_SOURCE})}

_cli_bash_autocomplete() {
  if [[ "${COMP_WORDS[0]}" != "source" ]]; then
    local cur opts base
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    if [[ "$cur" == "-"* ]]; then
      opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} ${cur} --generate-bash-completion )
    else
      opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} --generate-bash-completion )
    fi
    COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
    return 0
  fi
}

complete -o bashdefault -o default -o nospace -F _cli_bash_autocomplete $PROG
unset PROG`
	zsh = `#compdef $PROG

_cli_zsh_autocomplete() {
  local -a opts
  local cur
  cur=${words[-1]}
  if [[ "$cur" == "-"* ]]; then
    opts=("${(@f)$(${words[@]:0:#words[@]-1} ${cur} --generate-bash-completion)}")
  else
    opts=("${(@f)$(${words[@]:0:#words[@]-1} --generate-bash-completion)}")
  fi

  if [[ "${opts[1]}" != "" ]]; then
    _describe 'values' opts
  else
    _files
  fi
}

compdef _cli_zsh_autocomplete $PROG`
	zshRc = `
PROG=rssnix
_CLI_ZSH_AUTOCOMPLETE_HACK=1
source  ~/.rssnixrc`
)

func SetupAutocomplete(shell string) error {
	switch shell {
	case "bash":
		return setupBash()
	case "zsh":
		return setupZsh()
	}
	return nil
}

func setupBash() error {
	err := ioutil.WriteFile("/etc/bash_completion.d/rssnix", []byte(bash), 0644)
	fmt.Println("Please either restart your shell or run:")
	fmt.Println("\tsource ~/.bashrc")
	return err
}

func setupZsh() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	rssnixrcPath := path.Join(home, ".rssnixrc")
	err = ioutil.WriteFile(rssnixrcPath, []byte(zsh), 0644)
	if err != nil {
		return err
	}

	zshPath := path.Join(home, ".zshrc")
	f, err := os.OpenFile(zshPath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(zshRc); err != nil {
		return err
	}

	fmt.Println("Please either restart your shell or run:")
	fmt.Println("\tsource ~/.zshrc")
	return nil
}
