#!/bin/sh

set -e

kind=${1-}
case "$kind" in
patch | minor | major) ;;
*)
	printf 'usage: %s patch|minor|major\n' "$0" >&2
	exit 1
	;;
esac

top=$(git rev-parse --show-toplevel 2>/dev/null) || {
	printf '%s\n' 'error: not a git repository' >&2
	exit 1
}
cd "$top" || exit 1

git fetch --tags --quiet 2>/dev/null || :

base=$(git tag --list --sort=-version:refname 2>/dev/null | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$' | head -n 1)

if [ -z "$base" ]; then
	case "$kind" in
	major) next=v1.0.0 ;;
	minor) next=v0.1.0 ;;
	patch) next=v0.0.1 ;;
	esac
else
	v=${base#v}
	maj=$(printf '%s\n' "$v" | cut -d. -f1)
	min=$(printf '%s\n' "$v" | cut -d. -f2)
	pat=$(printf '%s\n' "$v" | cut -d. -f3)
	case "$kind" in
	patch) pat=$((pat + 1)) ;;
	minor)
		min=$((min + 1))
		pat=0
		;;
	major)
		maj=$((maj + 1))
		min=0
		pat=0
		;;
	esac
	next="v${maj}.${min}.${pat}"
fi

if git rev-parse -q --verify "refs/tags/$next" >/dev/null 2>&1; then
	printf '%s\n' "error: tag $next already exists" >&2
	exit 1
fi

printf 'Create, push and release tag %s? [y/n] ' "$next"
read -r confirm
case "$confirm" in
y) ;;
*)
	printf '%s\n' 'Aborted.' >&2
	exit 1
	;;
esac

git tag "$next"
git push origin "$next"
printf 'Pushed %s\n' "$next"
