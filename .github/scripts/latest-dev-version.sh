#!/bin/sh
# POSIX sh: print dev version (latest semver-like tag + "-dev", or 1.0.0-dev).
# Semver-like: optional v, MAJOR.MINOR.PATCH, optional prerelease (- or + …).

set -e

top=$(git rev-parse --show-toplevel 2>/dev/null) || {
	printf '%s\n' '1.0.0-dev'
	exit 0
}
cd "$top" || {
	printf '%s\n' '1.0.0-dev'
	exit 0
}

git fetch --tags --quiet 2>/dev/null || :

tag=$(git tag --list --sort=-version:refname 2>/dev/null | grep -E '^v?[0-9]+\.[0-9]+\.[0-9]+([-+].*)?$' | head -n 1)

if [ -z "$tag" ]; then
	printf '%s\n' '1.0.0-dev'
	exit 0
fi

ver=${tag#v}
printf '%s\n' "${ver}-dev"
