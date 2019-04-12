FROM archlinux/base
USER root
RUN pacman -Syyu --noconfirm
RUN yes | pacman -Sy go
RUN yes | pacman -Sy npm
RUN yes | pacman -Sy dep
RUN yes | pacman -Sy git
RUN yes | pacman -Sy docker
