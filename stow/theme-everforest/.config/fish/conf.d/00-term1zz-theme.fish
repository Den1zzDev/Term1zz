# Everforest Dark theme preset for fish
if status is-interactive
    if test -f ~/.config/fish/themes/everforest.theme
        fish_config theme choose "everforest" 2>/dev/null || true
    end
end
