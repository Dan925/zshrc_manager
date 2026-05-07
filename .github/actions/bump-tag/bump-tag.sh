#!/usr/bin/env bash
set -euo pipefail

LATEST=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
MSG=$(git log -1 --pretty=%s)

if echo "$MSG" | grep -qE "^(chore|docs|ci):"; then
  echo "Skipping tag for commit: $MSG"
  echo "new_tag=" >> $GITHUB_OUTPUT
  exit 0
fi

VERSION=${LATEST#v}
MAJOR=$(echo $VERSION | cut -d. -f1)
MINOR=$(echo $VERSION | cut -d. -f2)
PATCH=$(echo $VERSION | cut -d. -f3)

if echo "$MSG" | grep -qE "^feat(\(.+\))?!:|BREAKING CHANGE"; then
  MAJOR=$((MAJOR + 1)); MINOR=0; PATCH=0
elif echo "$MSG" | grep -qE "^feat(\(.+\))?:"; then
  MINOR=$((MINOR + 1)); PATCH=0
else
  PATCH=$((PATCH + 1))
fi

NEW_TAG="v$MAJOR.$MINOR.$PATCH"
echo "Creating tag: $NEW_TAG"

git config user.name "github-actions"
git config user.email "github-actions@github.com"
git tag "$NEW_TAG"
git push origin "$NEW_TAG"

echo "new_tag=$NEW_TAG" >> $GITHUB_OUTPUT
