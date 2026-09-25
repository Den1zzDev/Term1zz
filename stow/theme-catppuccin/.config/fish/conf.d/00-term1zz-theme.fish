# Catppuccin Mocha theme preset for fish
if status is-interactive
    if test -f ~/.config/fish/themes/catppuccin-mocha.theme
        fish_config theme choose "catppuccin-mocha" 2>/dev/null || true
    end
end
