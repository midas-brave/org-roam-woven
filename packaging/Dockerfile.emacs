FROM alpine:latest
LABEL org.opencontainers.image.authors="master@southfox.me"

# RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && apk update
RUN apk add --no-cache sqlite emacs-nox && mkdir ~/.emacs.d/
COPY ./init.el /root/.emacs.d/init.el

RUN emacs --daemon && emacsclient --eval '(kill-emacs)'
