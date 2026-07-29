package main

import (
	"flag"
	"fmt"
	"log"
)

const bashCompletion = `_codedock_tunnel_completion() {
    local cur prev opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    opts="login open create list inspect start stop revoke http tcp health completion version"
    COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
    return 0
}
complete -F _codedock_tunnel_completion codedock-tunnel
`

const zshCompletion = `#compdef codedock-tunnel
_codedock_tunnel() {
    local -a commands
    commands=(
        'login:Log in to Codedock Tunnel'
        'open:Open a tunnel'
        'create:Create a new tunnel'
        'list:List active tunnels'
        'inspect:Inspect tunnel details'
        'start:Start a tunnel'
        'stop:Stop a tunnel'
        'revoke:Revoke a tunnel'
        'http:Open an HTTP tunnel'
        'tcp:Open a TCP tunnel'
        'health:Check relay and API health'
        'completion:Generate shell completion'
        'version:Print version'
    )
    _describe -t commands 'command' commands
}
_codedock_tunnel "$@"
`

const fishCompletion = `complete -c codedock-tunnel -n "__fish_use_subcommand" -a "login open create list inspect start stop revoke http tcp health completion version"
`

func runCompletion(args []string) {
	flags := flag.NewFlagSet("completion", flag.ExitOnError)
	shell := flags.String("shell", "bash", "shell type: bash, zsh, fish")
	_ = flags.Parse(args)

	switch *shell {
	case "bash":
		fmt.Print(bashCompletion)
	case "zsh":
		fmt.Print(zshCompletion)
	case "fish":
		fmt.Print(fishCompletion)
	default:
		log.Fatalf("unsupported shell %q; supported: bash, zsh, fish", *shell)
	}
}
