package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	root.AddCommand(completionCmd)
}

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion <bash|zsh>",
	Short: "Generate bash/zsh completion scripts",
	Long: `To load completion run

For zsh:
source <(faas completion zsh)

If you would like to use alias:
alias f=faas
compdef _faas f

For bash:
source <(faas completion bash)

`,
	ValidArgs: []string{"bash", "zsh"},
	Args:      cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if len(args) < 1 {
			return errors.New("missing argument")
		}
		if args[0] == "bash" {
			err = root.GenBashCompletion(os.Stdout)
			return err
		}
		if args[0] == "zsh" {
			// manually edited script based on `root.GenZshCompletion(os.Stdout)`
			// unfortunately it doesn't support completion so well as for bash
			// some manual edits had to be done
			os.Stdout.WriteString(`
compdef _faas faas

function _faas {
  local -a commands

  _arguments -C \
    '--config[config file path]:file:_files' \
    '(-v --verbose)'{-v,--verbose}'[print verbose logs]' \
    "1: :->cmnds" \
    "*::arg:->args"

  case $state in
  cmnds)
    commands=(
      "completion:Generates bash/zsh completion scripts"
      "create:Create a Function"
      "delete:Delete deployed Function"
      "describe:Describe Function"
      "help:Help about any command"
      "list:Lists deployed Functions"
      "run:Run Function locally"
      "update:Update or create a deployed Function"
      "version:Print version"
    )
    _describe "command" commands
    ;;
  esac

  case "$words[1]" in
  completion)
    _faas_completion
    ;;
  create)
    _faas_create
    ;;
  delete)
    _faas_delete
    ;;
  describe)
    _faas_describe
    ;;
  help)
    _faas_help
    ;;
  list)
    _faas_list
    ;;
  run)
    _faas_run
    ;;
  update)
    _faas_update
    ;;
  version)
    _faas_version
    ;;
  esac
}

function _list_funs() {
    compadd $(faas list 2> /dev/null)
}

function _list_langs() {
		compadd node go quarkus
}

function _list_fmts() {
    compadd yaml xml json
}

function _list_regs() {
    local config="${HOME}/.docker/config.json"
    if command -v yq >/dev/null && test -f "$config";  then
		compadd $(jq -r ".auths | keys[] " "$config")
	fi
}

function _faas_completion {
  _arguments \
    '(-h --help)'{-h,--help}'[help for completion]' \
    '--config[config file path]:file:_files' \
    '(-v --verbose)'{-v,--verbose}'[print verbose logs]' \
    '1: :("bash" "zsh")'
}

function _faas_create {
  _arguments \
    '1:string:_list_langs' \
    '(-i --internal)'{-i,--internal}'[Create a cluster-local service without a publicly accessible route. $FAAS_INTERNAL]' \
    '(-l --local)'{-l,--local}'[create the service function locally only.]' \
    '(-n --name)'{-n,--name}'[optionally specify an explicit name for the serive, overriding path-derivation. $FAAS_NAME]:' \
    '(-s --namespace)'{-s,--namespace}'[namespace at image registry (usually username or org name). $FAAS_NAMESPACE]:' \
    '(-r --registry)'{-r,--registry}'[image registry (ex: quay.io). $FAAS_REGISTRY]:string:_list_regs' \
    '--config[config file path]:file:_files' \
    '(-v --verbose)'{-v,--verbose}'[print verbose logs]'
}

function _faas_delete {
  _arguments \
    '(-n --name)'{-n,--name}'[optionally specify an explicit name to remove, overriding path-derivation. $FAAS_NAME]:string:_list_funs' \
    '--config[config file path]:file:_files' \
    '(-v --verbose)'{-v,--verbose}'[print verbose logs]'
}

function _faas_describe {
  _arguments \
    '1:string:_list_funs' \
    '(-n --name)'{-n,--name}'[optionally specify an explicit name for the serive, overriding path-derivation. $FAAS_NAME]:string:_list_funs' \
    '(-o --output)'{-o,--output}'[optionally specify output format (yaml,xml,json).]:string:_list_fmts' \
    '--config[config file path]:file:_files' \
    '(-v --verbose)'{-v,--verbose}'[print verbose logs]'
}

function _faas_help {
  _arguments \
    '--config[config file path]:file:_files' \
    '(-v --verbose)'{-v,--verbose}'[print verbose logs]'
}

function _faas_list {
  _arguments \
    '--config[config file path]:file:_files' \
    '(-v --verbose)'{-v,--verbose}'[print verbose logs]'
}

function _faas_run {
  _arguments \
    '--config[config file path]:file:_files' \
    '(-v --verbose)'{-v,--verbose}'[print verbose logs]'
}

function _faas_update {
  _arguments \
    '--config[config file path]:file:_files' \
    '(-s --namespace)'{-s,--namespace}'[namespace at image registry (usually username or org name). $FAAS_NAMESPACE]:' \
    '(-r --registry)'{-r,--registry}'[image registry (ex: quay.io). $FAAS_REGISTRY]:string:_list_regs' \
    '(-v --verbose)'{-v,--verbose}'[print verbose logs]'
}

function _faas_version {
  _arguments \
    '--config[config file path]:file:_files' \
    '(-v --verbose)'{-v,--verbose}'[print verbose logs]'
}


`)
			return nil
		}
		return errors.New("unknown shell, only bash and zsh are supported")
	},
}
