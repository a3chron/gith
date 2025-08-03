# Fish completion for gith
complete -c gith -f
complete -c gith -n "__fish_use_subcommand" -a "version config help" -d "Available commands"
complete -c gith -n "__fish_seen_subcommand_from version" -a "check" -d "Check for updates"
complete -c gith -n "__fish_seen_subcommand_from config" -a "show reset path update help" -d "Config commands"
complete -c gith -n "__fish_seen_subcommand_from config update" -l flavor -d "Catppuccin flavor" -a "latte frappe macchiato mocha"
complete -c gith -n "__fish_seen_subcommand_from config update" -l accent -d "Catppuccin accent" -a "rosewater flamingo pink mauve red maroon peach yellow green teal sky sapphire blue lavender gray"