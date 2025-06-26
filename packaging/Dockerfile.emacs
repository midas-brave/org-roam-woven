FROM alpine:latest
LABEL org.opencontainers.image.authors="master@southfox.me"

ARG CN_MIRROR
RUN if [ "$CN_MIRROR" = "true" ]; then \
    sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
    apk update && \
    apk add --no-cache sqlite emacs-nox && \
    mkdir ~/.emacs.d/; \
    else \
    apk add --no-cache sqlite emacs-nox && \
    mkdir ~/.emacs.d/; \
    fi
COPY ./init.el /root/.emacs.d/init.el

RUN emacs --daemon && emacsclient --eval '(kill-emacs)'
