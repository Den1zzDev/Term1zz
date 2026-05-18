# ──────────────────────────────────────────────
# ENVIRONMENT VARIABLES
# ──────────────────────────────────────────────
set -gx VISUAL   zeditor
set -gx EDITOR   zeditor
set -gx TERM     xterm-256color

# bat → Catppuccin Mocha theme everywhere
set -gx BAT_THEME "Catppuccin Mocha"

# fzf → match Catppuccin Mocha palette + bat preview
set -gx FZF_DEFAULT_OPTS "\
  --height=50% --layout=reverse --border=rounded --info=inline \
  --prompt='  ' --pointer='▶' --marker='✓' \
  --color=bg+:#313244,bg:#1e1e2e,spinner:#f5e0dc,hl:#f38ba8 \
  --color=fg:#cdd6f4,header:#f38ba8,info:#cba6f7,pointer:#f5e0dc \
  --color=marker:#f5e0dc,fg+:#cdd6f4,prompt:#cba6f7,hl+:#f38ba8 \
  --preview='bat --style=numbers --color=always {}' \
  --preview-window='right:50%:wrap'"
set -gx FZF_DEFAULT_COMMAND "fd --type f --hidden --follow --exclude .git"
set -gx FZF_CTRL_T_COMMAND  "$FZF_DEFAULT_COMMAND"
set -gx FZF_ALT_C_COMMAND   "fd --type d --hidden --follow --exclude .git"

# less → colours
set -gx LESS "-R --use-color"
set -gx MANPAGER "sh -c 'col -bx | bat -l man -p'"

# ──────────────────────────────────────────────
# FISH SYNTAX-HIGHLIGHT COLORS  (Catppuccin Mocha)
# ──────────────────────────────────────────────
set -U fish_color_normal         cdd6f4          # text
set -U fish_color_command        cba6f7          # commands  → mauve
set -U fish_color_keyword        cba6f7          # keywords  → mauve
set -U fish_color_quote          a6e3a1          # strings   → green
set -U fish_color_redirection    f5c2e7          # redirects → pink
set -U fish_color_end            fab387          # semicolon → peach
set -U fish_color_error          f38ba8          # errors    → red
set -U fish_color_param          cdd6f4          # params    → text
set -U fish_color_comment        9399b2          # comments  → overlay2
set -U fish_color_operator       89dceb          # operators → sky
set -U fish_color_escape         f5c2e7          # escapes   → pink
set -U fish_color_autosuggestion 6c7086          # ghost text → overlay0
set -U fish_color_valid_path     --underline
set -U fish_color_match          --background=313244
set -U fish_color_search_match   --background=313244
set -U fish_pager_color_prefix   cba6f7 --bold --underline
set -U fish_pager_color_completion cdd6f4
set -U fish_pager_color_description 6c7086 --italics
set -U fish_pager_color_selected_background --background=313244

# ──────────────────────────────────────────────
# INTERACTIVE SESSION
# ──────────────────────────────────────────────
if status is-interactive

    # ── Greeting ─────────────────────────────
    set -g fish_greeting ""
    clear
    if command -q bfetch
        bfetch
    else if command -q fastfetch
        fastfetch
    end

    # ── Prompt: Starship ──────────────────────
    # If starship is available it takes over; falls back to fish default.
    if command -q starship
        starship init fish | source
    end

    # ── Smart cd: zoxide ─────────────────────
    if command -q zoxide
        zoxide init --cmd cd fish | source
    end

    # ── Fuzzy finder: fzf keybindings ────────
    # Ctrl-R → fzf history  |  Ctrl-T → fzf files  |  Alt-C → fzf dirs
    if command -q fzf
        fzf --fish | source
    end

    # ── Atuin: shell history sync (optional) ─
    if command -q atuin
        atuin init fish --disable-up-arrow | source
    end

    # =========================================
    # ALIASES — Modern CLI replacements
    # =========================================

    # ls → eza
    if command -q eza
        alias ls  "eza --icons --color=always --group-directories-first"
        alias ll  "eza -lah --icons --color=always --group-directories-first --git"
        alias la  "eza -a   --icons --color=always --group-directories-first"
        alias lt  "eza --tree --level=2 --icons --color=always"
        alias ltt "eza --tree --level=3 --icons --color=always"
        alias lg  "eza -lah --icons --color=always --git --git-repos"
    end

    # cat → bat
    if command -q bat
        alias cat  "bat"
        alias less "bat --paging=always"
    end

    # du → dust  (if installed)
    if command -q dust
        alias du "dust"
    end

    # top → btop  (if installed)
    if command -q btop
        alias top "btop"
    end

    # find → fd  (if installed)
    if command -q fd
        alias find "fd"
    end

    # grep → ripgrep  (if installed)
    if command -q rg
        alias grep "rg"
    end

    # nano → micro  (if installed)
    if command -q micro
        alias nano "micro"
    end

    # =========================================
    # uutils coreutils (Dynamic Mapping)
    # =========================================
    set -g uutils_exclude ls cat du top find grep

    # Dynamically alias EVERY uu-* binary (safely)
    for uu_bin in /usr/bin/uu-*
        set -l base_cmd (string replace 'uu-' '' (basename $uu_bin))

        # NEVER override Fish internal builtins (test, echo, true, etc.)
        if contains $base_cmd (builtin -n)
            continue
        end

        if not contains $base_cmd $uutils_exclude
            if contains $base_cmd cp mv rm
                alias $base_cmd "$uu_bin -i"
            else if test "$base_cmd" = "mkdir"
                alias $base_cmd "$uu_bin -pv"
            else
                alias $base_cmd "$uu_bin"
            end
        end
    end

    # =========================================
    # Sudo muscle-memory wrapper (routes to run0)
    # =========================================
    function sudo --description "Type sudo, but execute run0 safely with uutils"
        set -l run0_opts --background=""

        if test (count $argv) -eq 0
            command run0 $run0_opts
            return
        end

        set -l cmd $argv[1]

        # If passing flags, pass them directly to run0
        if string match -q -- "-*" "$cmd"
            command run0 $run0_opts $argv
            return
        end

        set -e argv[1] # Shift the arguments

        # ---------------------------------------------------------
        # Dynamic Alias/Function Resolution
        # ---------------------------------------------------------
        # If the command is a Fish alias (which is a function with --wraps),
        # extract the underlying binary and prepend its flags to argv.
        if functions -q "$cmd"
            set -l func_def (functions "$cmd")
            set -l wrapped (string match -r -- "--wraps='([^']+)'" "$func_def")
            if test -z "$wrapped"
                set wrapped (string match -r -- "--wraps=([^\s]+)" "$func_def")
            end

            if test -n "$wrapped"
                set -l split_target (string split " " -- $wrapped[2])
                set cmd $split_target[1]
                if test (count $split_target) -gt 1
                    set argv $split_target[2..-1] $argv
                end
            end
        end
        # ---------------------------------------------------------

        if test -x "/usr/bin/uu-$cmd"; and not contains "$cmd" $uutils_exclude
            if contains "$cmd" cp mv rm
                command run0 $run0_opts "uu-$cmd" -i $argv
            else if test "$cmd" = "mkdir"
                command run0 $run0_opts "uu-$cmd" -pv $argv
            else
                command run0 $run0_opts "uu-$cmd" $argv
            end
        else
            command run0 $run0_opts $cmd $argv
        end
    end

    # =========================================
    # NAVIGATION
    # =========================================
    alias ..    "cd .."
    alias ...   "cd ../.."
    alias ....  "cd ../../.."
    alias ..... "cd ../../../.."
    # Note: '~' and '-' can't be alias/function names in fish.
    # Use 'cd ~' and 'cd -' directly, or bind them as abbreviations:
    abbr -- cdh "cd ~"            # go home
    abbr -- cdb "cd -"            # jump back to last dir

    # =========================================
    # ABBREVIATIONS  (expand on <Space>)
    # =========================================

    # — Git —
    abbr g      "git"
    abbr gi     "git init"
    abbr gs     "git status -sb"
    abbr ga     "git add"
    abbr gaa    "git add -A"
    abbr gc     "git commit -m"
    abbr gca    "git commit --amend --no-edit"
    abbr gco    "git checkout"
    abbr gcb    "git checkout -b"
    abbr gp     "git push"
    abbr gpf    "git push --force-with-lease"
    abbr gpu    "git push -u origin HEAD"
    abbr pl     "git pull"
    abbr gpl    "git pull --rebase"
    abbr gf     "git fetch --all --prune"
    abbr glog   "git log --oneline --graph --decorate --all"
    abbr gd     "git diff"
    abbr gds    "git diff --staged"
    abbr gb     "git branch"
    abbr gba    "git branch -a"
    abbr gbd    "git branch -d"
    abbr gst    "git stash"
    abbr gstp   "git stash pop"
    abbr grb    "git rebase"
    abbr grbi   "git rebase -i"
    abbr gcp    "git cherry-pick"
    abbr gtag   "git tag"
    abbr greset "git reset --hard HEAD"
    abbr gclean "git clean -fd"

    # — GitHub CLI —
    abbr ghpr   "gh pr create"
    abbr ghprl  "gh pr list"
    abbr ghprv  "gh pr view --web"

    # — Package managers —
    abbr ni   "npm install"
    abbr nid  "npm install --save-dev"
    abbr nig  "npm install -g"
    abbr nr   "npm run"
    abbr nrd  "npm run dev"
    abbr nrb  "npm run build"
    abbr nrt  "npm run test"
    abbr nrc  "npm run check"
    abbr nrp  "npm run preview"

    abbr pi   "pnpm install"
    abbr pa   "pnpm add"
    abbr pad  "pnpm add -D"
    abbr pr   "pnpm run"
    abbr prd  "pnpm run dev"
    abbr prb  "pnpm run build"

    abbr yi   "yarn install"
    abbr ya   "yarn add"
    abbr yr   "yarn run"

    # — Python —
    abbr py    "python3"
    abbr pip   "pip3"
    abbr venv  "python3 -m venv .venv && source .venv/bin/activate.fish"
    abbr va    "source .venv/bin/activate.fish"
    abbr vd    "deactivate"

    # — Docker —
    abbr d      "docker"
    abbr dc     "docker compose"
    abbr dcu    "docker compose up -d"
    abbr dcd    "docker compose down"
    abbr dcl    "docker compose logs -f"
    abbr dps    "docker ps"
    abbr dpsa   "docker ps -a"
    abbr dimg   "docker images"
    abbr dex    "docker exec -it"
    abbr drm    "docker rm"
    abbr drmi   "docker rmi"
    abbr dprune "docker system prune -af --volumes"

    # — Kubernetes —
    abbr k      "kubectl"
    abbr kgp    "kubectl get pods"
    abbr kgs    "kubectl get svc"
    abbr kgn    "kubectl get nodes"
    abbr kd     "kubectl describe"
    abbr kl     "kubectl logs -f"
    abbr kex    "kubectl exec -it"
    abbr kctx   "kubectx"
    abbr kns    "kubens"
    abbr kaf    "kubectl apply -f"
    abbr kdf    "kubectl delete -f"

    # — Misc shortcuts —
    abbr e      "$EDITOR ."
    abbr q      "exit"
    abbr cl     "clear"
    abbr reload "source ~/.config/fish/config.fish"
    abbr fishrc "$EDITOR ~/.config/fish/config.fish"
    abbr path   "printf '%s\n' \$PATH"
    abbr ports  "ss -tulnp"
    abbr myip   "curl -s ifconfig.me"
    abbr weather "curl -s wttr.in"
    abbr speedtest "curl -s https://raw.githubusercontent.com/sivel/speedtest-cli/master/speedtest.py | python3 - --secure"

    # =========================================
    # USEFUL FUNCTIONS (inline)
    # =========================================

    # mkcd: make a directory and cd into it
    function mkcd --description "mkdir + cd"
        mkdir -p $argv[1] && cd $argv[1]
    end

    # up N: go up N directories
    function up --description "Go up N directories"
        set -l n 1
        if test (count $argv) -gt 0
            set n $argv[1]
        end
        set -l path ""
        for i in (seq 1 $n)
            set path "../$path"
        end
        cd $path
    end

    # extract: universal archive extractor
    function extract --description "Extract any archive"
        if test -f $argv[1]
            switch $argv[1]
                case "*.tar.bz2";  tar xjf  $argv[1]
                case "*.tar.gz";   tar xzf  $argv[1]
                case "*.tar.xz";   tar xJf  $argv[1]
                case "*.tar.zst";  tar --zstd -xf $argv[1]
                case "*.bz2";      bunzip2  $argv[1]
                case "*.gz";       gunzip   $argv[1]
                case "*.tar";      tar xf   $argv[1]
                case "*.tbz2";     tar xjf  $argv[1]
                case "*.tgz";      tar xzf  $argv[1]
                case "*.zip";      unzip    $argv[1]
                case "*.7z";       7z x     $argv[1]
                case "*.rar";      unrar x  $argv[1]
                case "*.xz";       xz -d    $argv[1]
                case "*.zst";      zstd -d  $argv[1]
                case "*";          echo "'$argv[1]' — unknown archive format"
            end
        else
            echo "'$argv[1]' is not a valid file"
        end
    end

    # fcd: fzf-powered interactive directory jump (fallback when zoxide unavailable)
    function fcd --description "Fuzzy cd"
        set dir (fd --type d --hidden --exclude .git | fzf +m)
        and cd $dir
    end

    # fkill: interactively kill a process via fzf
    function fkill --description "Fuzzy kill process"
        set pid (ps -eo pid,comm,args | fzf --header="Select process to kill" | awk '{print $1}')
        and kill -9 $pid
        and echo "Killed PID $pid"
    end

    # gitroot: cd to repo root
    function gitroot --description "cd to git root"
        cd (git rev-parse --show-toplevel)
    end

    # sshfzf: pick from known SSH hosts via fzf
    function sshfzf --description "Fuzzy SSH connect"
        set host (grep -E "^Host " ~/.ssh/config 2>/dev/null | awk '{print $2}' | fzf --prompt=" SSH host: ")
        and ssh $host
    end

    # twitch: watch a Twitch stream via streamlink + mpv
    function twitch --description "Watch a Twitch stream in mpv"
        streamlink --player "mpv" https://www.twitch.tv/$argv[1] best
    end

end
# ╚══════════════════════════════════════════════════════════════╝

fish_add_path "$HOME/.local/bin"
if status is-interactive
    if not pgrep -xu $USER ssh-agent >/dev/null
        eval (ssh-agent -c)
        set -Ux SSH_AUTH_SOCK $SSH_AUTH_SOCK
        set -Ux SSH_AGENT_PID $SSH_AGENT_PID
    end
end
