#compdef runny

# Install as _runny
#
# e.g. ln -s $PWD/runny.zsh /usr/local/share/zsh/site-functions/_runny

_arguments "1: :($(runny | cut -f1 ))"
