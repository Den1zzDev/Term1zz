# Tokyo Night theme preset for fish
if status is-interactive
    if test -f ~/.config/fish/themes/tokyonight.theme
        fish_config theme choose "tokyonight" 2>/dev/null || true
    end
end
