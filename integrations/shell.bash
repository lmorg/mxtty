output_block_integration () {
    [ -n "$COMP_LINE" ] && return
    [ "$BASH_COMMAND" = "$PROMPT_COMMAND" ] && return

    printf "\033_begin;output-block\033\\"
}

PROMPT_COMMAND='printf "\033_end;output-block;{\"ExitNum\":$?}\033\\"'

trap 'output_block_integration' DEBUG
