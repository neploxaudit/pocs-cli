package config

import (
	"net/url"

	"github.com/neploxaudit/pocs-cli/pkg/pocs"
)

var (
	BaseURL     *url.URL
	NeploxToken string
)

var Client *pocs.Client
