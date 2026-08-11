#!/usr/bin/env bash
# Sourceable fixture helpers for shell behavior probes. This file runs no cases itself.

fixture_repo_init() {
  local repository_path="$1"
  git -c init.defaultBranch=main init -q "$repository_path"
  git -C "$repository_path" config user.name 'Fixture Runner'
  git -C "$repository_path" config user.email 'fixture@example.invalid'
}

fixture_repo_commit_all() {
  local repository_path="$1"
  local commit_message="$2"
  git -C "$repository_path" add -A
  git -C "$repository_path" commit -qm "$commit_message"
}

fixture_repo_clone_script() {
  local script_source="$1"
  local script_target="$2"
  mkdir -p "$(dirname "$script_target")"
  cp "$script_source" "$script_target"
  chmod +x "$script_target"
}
