# hotdog

block grep command

## how to install
```
go get -u github.com/umaumax/hotdog/...
```

## how to use
```
# change separator
cat ~/dotfiles/.peco.zshrc |hotdog -first='^\s*#{30}' -middle='^\s*#' -last='^\s*#{30}' -separator="$(echo "\0x0a\x00")"
# split separator (multi lines)
cat ~/dotfiles/.peco.zshrc |hotdog -first='^\s*#{30}' -middle='^\s*#' -last='^\s*#{30}' | tr "\036" '\n'
```
