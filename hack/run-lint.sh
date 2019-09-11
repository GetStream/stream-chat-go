#!/usr/bin/env bash

set -euo pipefail

[[ -n ${DEBUG:-} ]] && set -x

gopath="$(go env GOPATH)"

if ! [[ -x "$gopath/bin/golangci-lint" ]]; then
	echo >&2 'Installing golangci-lint'
	curl --silent --fail --location \
		https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b "$gopath/bin" v1.18.0
fi

# configured by .golangci.yml
"$gopath/bin/golangci-lint" run

install_impi() {
	impi_dir="$(mktemp -d)"
	trap 'rm -rf -- ${impi_dir}' EXIT

	GOPATH="${impi_dir}" \
		GOBIN="${gopath}/bin" \
		go get github.com/pavius/impi/cmd/impi
}

# install impi that ensures import grouping is done consistently
if ! [[ -x "${gopath}/bin/impi" ]]; then
	echo >&2 'Installing impi'
	install_impi
fi

"$gopath/bin/impi" \
	--local github.com/GetStream/stream-chat-go \
	--scheme stdThirdPartyLocal \
	--skip "stream_chat_easyjson.go" \
	./...

install_shfmt() {
	shfmt_dir="$(mktemp -d)"
	trap 'rm -rf -- ${shfmt_dir}' EXIT

	GOPATH="${shfmt_dir}" \
		GO111MODULE=off \
		GOBIN="${gopath}/bin" \
		go get mvdan.cc/sh/cmd/shfmt
}

# install shfmt that ensures consistent format in shell scripts
if ! [[ -x "${gopath}/bin/shfmt" ]]; then
	echo >&2 'Installing shfmt'
	install_shfmt
fi

SCRIPTDIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
shfmt_out="$($gopath/bin/shfmt -l -s ${SCRIPTDIR})"
if [[ -n ${shfmt_out} ]]; then
	echo >&2 "The following shell scripts need to be formatted, run: 'shfmt -w -s ${SCRIPTDIR}'"
	echo >&2 "${shfmt_out}"
	exit 1
fi
