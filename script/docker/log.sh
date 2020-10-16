#!/bin/bash

set +e
set -o noglob
bold=; underline=; reset=; red=; green=; white=; tan=; blue=;
#echo "the value of TERM is ${TERM}"
if test ! dumb = "${TERM}" ; then
        bold=$(tput bold)
        underline=$(tput sgr 0 1)
        reset=$(tput sgr0)
        red=$(tput setaf 1)
        green=$(tput setaf 76)
        white=$(tput setaf 7)
fi
function underline() { printf "${underline}${bold}%s${reset}\n" "$@"
}
function h1() { printf "\n${bold}${blue}%s${reset}\n" "$@"
}
function h2() { printf "\n${bold}${white}%s${reset}\n" "$@"
}
function debug() { printf "${white}%s${reset}\n" "$@"
}
function info() { printf "${white}➜ %s${reset}\n" "$@"
}
function success() { printf "${green}✔ %s${reset}\n" "$@"
}
function error() { printf "${red}✖ %s${reset}\n" "$@"
}
function warn() { printf "${tan}➜ %s${reset}\n" "$@"
}
function bold() { printf "${bold}%s${reset}\n" "$@"
}
function note() { printf "\n${underline}${bold}${blue}Note:${reset} ${blue}%s${reset}\n" "$@"
}
function hlog() { echo -e "\033[32;32m $1 \033[0m \n"
}
set -e
set +o noglob
