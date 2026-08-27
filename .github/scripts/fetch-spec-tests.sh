#!/usr/bin/env bash
#
# Lay out a consensus-specs release so the spec vector runner can read it.
#
#   .github/scripts/fetch-spec-tests.sh [version] [preset...]
#
# Downloads the vector tarball of every requested preset (default: mainnet)
# into .spec-tests/<version>:
#
#   .spec-tests/<version>/tests/<preset>/<fork>/ssz_static/...
#
# Without a version, or with "latest", the newest release that carries the
# requested vectors is used. Releases of the specs are pre-releases, so this is
# the newest release by date rather than what the API calls the latest one.
#
#   .github/scripts/fetch-spec-tests.sh                    # newest, mainnet
#   .github/scripts/fetch-spec-tests.sh latest minimal     # newest, minimal
#   .github/scripts/fetch-spec-tests.sh v1.7.0-alpha.14    # a specific release
#   .github/scripts/fetch-spec-tests.sh --latest-version   # only resolve, print, exit
#
# Run the vectors against the result with:
#
#   CONSENSUS_SPEC_TESTS_DIR=.spec-tests/<version> go test ./spec/...

set -euo pipefail

REPO="ethereum/consensus-specs"

# github_api fetches a repository API path, authenticated through gh when it is
# available so CI does not run into the anonymous rate limit.
github_api() {
    local path="$1"

    if command -v gh >/dev/null 2>&1 && gh auth status >/dev/null 2>&1; then
        gh api "repos/$REPO/$path"
    else
        curl -fsSL ${GITHUB_TOKEN:+-H "Authorization: Bearer $GITHUB_TOKEN"} \
            -H "Accept: application/vnd.github+json" \
            "https://api.github.com/repos/$REPO/$path"
    fi
}

# latest_version prints the newest release that has the vectors of the given
# preset attached. A release is cut before its assets finish uploading, so the
# newest release is not necessarily one that can be downloaded yet.
latest_version() {
    local preset="$1"

    github_api "releases?per_page=20" | python3 -c '
import json, sys

preset = sys.argv[1]
for release in json.load(sys.stdin):
    if release.get("draft"):
        continue
    if any(asset["name"] == preset + ".tar.gz" for asset in release.get("assets", [])):
        print(release["tag_name"])
        break
' "$preset"
}

if [ "${1:-}" = "--latest-version" ]; then
    version="$(latest_version "${2:-mainnet}")"
    if [ -z "$version" ]; then
        echo "could not resolve the newest $REPO release carrying vectors" >&2
        exit 1
    fi
    echo "$version"
    exit 0
fi

VERSION="${1:-latest}"
if [ $# -gt 0 ]; then shift; fi

PRESETS=("$@")
# The ssz_static vectors are read through the static SSZ path, whose sizes are
# the mainnet preset, so mainnet is the preset these tests need.
if [ ${#PRESETS[@]} -eq 0 ]; then
    PRESETS=("mainnet")
fi

if [ "$VERSION" = "latest" ]; then
    VERSION="$(latest_version "${PRESETS[0]}")"
    if [ -z "$VERSION" ]; then
        echo "could not resolve the newest $REPO release carrying vectors" >&2
        exit 1
    fi
    echo "==> newest release with ${PRESETS[0]} vectors: $VERSION"
fi

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
DEST="${SPEC_TESTS_DIR:-$ROOT/.spec-tests/$VERSION}"
mkdir -p "$DEST"

download() {
    local asset="$1" target="$2"
    local url="https://github.com/$REPO/releases/download/$VERSION/$asset"

    echo "==> $asset"
    if command -v gh >/dev/null 2>&1 && gh auth status >/dev/null 2>&1; then
        gh release download "$VERSION" -R "$REPO" -p "$asset" -O "$target" --clobber
    else
        curl -fSL --retry 3 -o "$target" "$url"
    fi
}

for preset in "${PRESETS[@]}"; do
    # The extracted tree carries the version in no way of its own, so the marker
    # records which release the directory was filled from.
    marker="$DEST/.$preset.stamp"
    if [ -f "$marker" ] && [ "$(cat "$marker")" = "$VERSION" ]; then
        echo "==> $preset vectors already present"
        continue
    fi

    archive="$DEST/$preset.tar.gz"
    download "$preset.tar.gz" "$archive"

    echo "==> extracting $preset"
    rm -rf "${DEST:?}/tests/$preset"
    tar xzf "$archive" -C "$DEST"
    rm -f "$archive"
    echo "$VERSION" > "$marker"
done

echo
echo "spec tests ready in $DEST ($VERSION)"
echo "run them with:"
echo "  CONSENSUS_SPEC_TESTS_DIR=$DEST go test ./spec/..."
