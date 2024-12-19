_output_block_begin() {
    printf "\033_begin;output-block\033\\"
}

_output_block_end() {
    printf "\033_end;output-block;{\"ExitNum\":$?}\033\\"
}

preexec_functions+=(_output_block_begin)
precmd_functions+=(_output_block_end)
