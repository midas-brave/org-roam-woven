;;; init.el --- Description -*- lexical-binding: t; -*-
;;;
(require 'package)
(setq package-enable-at-startup nil)
(setq gnutls-algorithm-priority "NORMAL:-VERS-TLS1.3")
(setq package-archives
      '(("MELPA"        . "https://melpa.org/packages/")))

(package-initialize)
(setq use-package-always-ensure t)

(use-package org-roam
  :ensure t)
