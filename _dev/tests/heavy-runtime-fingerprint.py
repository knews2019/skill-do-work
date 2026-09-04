#!/usr/bin/env python3
"""Seal the effective module toolchain; refuse opaque runtime extensions.

The lane runner bounds this probe's whole process group. Browser lanes have no
fingerprint declaration: browser discovery, profiles, fonts, and shared runtime
assets are not a complete, stable closure we can currently seal.
"""
import hashlib
import json
import os
from pathlib import Path
import shutil
import shlex
import subprocess
import sys


def binary_seal(binary_path):
    resolved_path = Path(binary_path).resolve(strict=True)
    contents = resolved_path.read_bytes()
    # A wrapper can delegate to mutable files outside this repository. Its own
    # bytes alone cannot attest to the executable that actually runs.
    native_headers = (b"\x7fELF", b"\xcf\xfa\xed\xfe", b"\xfe\xed\xfa\xcf", b"\xce\xfa\xed\xfe", b"\xfe\xed\xfa\xce", b"\xca\xfe\xba\xbe", b"\xbe\xba\xfe\xca", b"MZ")
    if not contents.startswith(native_headers):
        raise ValueError("opaque executable wrapper")
    return {"path": str(resolved_path), "sha256": hashlib.sha256(contents).hexdigest()}


def runtime_seal(module_directory, tool_names):
    # These are extension points with inputs outside the declared repository:
    # startup files, dynamically loaded code, or arbitrary compiler overlays.
    # Nonempty extensions require execution, even if their path text matches.
    for variable_name in ("BASH_ENV", "ENV", "PYTHONPATH", "PYTHONHOME", "NODE_OPTIONS", "LD_PRELOAD", "LD_LIBRARY_PATH", "DYLD_INSERT_LIBRARIES", "DYLD_LIBRARY_PATH"):
        if os.environ.get(variable_name):
            raise ValueError("opaque runtime extension: " + variable_name)
    # The shipped manifest explicitly isolates system/global Git config for
    # both lane and probe. Other caller-supplied Git injection remains opaque.
    for variable_name in os.environ:
        if variable_name.startswith("GIT_CONFIG_") and variable_name not in ("GIT_CONFIG_NOSYSTEM", "GIT_CONFIG_GLOBAL"):
            raise ValueError("opaque Git configuration override")
    tools = {name: binary_seal(shutil.which(name) or name) for name in tool_names}
    tools["go"] = binary_seal(shutil.which("go") or "go")
    go_environment = json.loads(subprocess.check_output(["go", "env", "-json"], cwd=module_directory))
    # GOGCCFLAGS is derived from this toolchain/configuration and contains a
    # fresh go-build scratch directory on every query. Keep its determining
    # inputs below, rather than hashing an incidental random pathname.
    go_environment.pop("GOGCCFLAGS", None)
    # Effective settings include the GOENV file and module-selected toolchain,
    # unlike `go version` at the repository root. Flags/workspaces may reference
    # arbitrary external inputs; do not attempt to parse a purported closure.
    if go_environment.get("GOFLAGS") or go_environment.get("GOWORK") not in ("", "off"):
        raise ValueError("opaque Go flags or workspace")
    if go_environment.get("CGO_ENABLED") == "1":
        for compiler_name in ("CC", "CXX"):
            compiler_argv = shlex.split(go_environment[compiler_name])
            if len(compiler_argv) != 1:
                raise ValueError("opaque compiler invocation")
            tools[compiler_name] = binary_seal(shutil.which(compiler_argv[0]) or compiler_argv[0])
    # These are Git's scalar repository/identity settings. Any other setting
    # could introduce an external helper, hook, filter, or credential provider;
    # execute rather than claim a complete runtime closure for it.
    scalar_git_keys = {"user.name", "user.email", "init.defaultbranch", "safe.directory", "core.repositoryformatversion", "core.filemode", "core.bare", "core.logallrefupdates", "core.ignorecase", "core.precomposeunicode", "core.worktree", "extensions.worktreeconfig", "extensions.objectformat"}
    git_configuration = subprocess.check_output(["git", "config", "--null", "--list"], cwd=module_directory).decode()
    for setting in git_configuration.split("\0"):
        if not setting:
            continue
        key = setting.split("\n", 1)[0].lower()
        parts = key.split(".")
        remote_scalar = len(parts) >= 3 and parts[0] == "remote" and parts[-1] in ("url", "fetch", "pushurl")
        branch_scalar = len(parts) >= 3 and parts[0] == "branch" and parts[-1] in ("remote", "merge", "vscode-merge-base")
        if key not in scalar_git_keys and not remote_scalar and not branch_scalar:
            raise ValueError("opaque Git configuration")
    selected_go = Path(go_environment["GOROOT"]) / "bin" / "go"
    tools["module_go"] = binary_seal(selected_go)
    module_records = subprocess.check_output([str(selected_go), "list", "-m", "-json", "all"], cwd=module_directory).decode()
    # Local replace targets can change without any repository commit. The
    # module listing uses concatenated JSON, so decode each record explicitly.
    decoder = json.JSONDecoder()
    remaining = module_records.strip()
    while remaining:
        record, offset = decoder.raw_decode(remaining)
        replacement = record.get("Replace")
        if replacement and not replacement.get("Version"):
            raise ValueError("opaque local module replacement")
        remaining = remaining[offset:].lstrip()
    supplied_cli = os.environ.get("DO_WORK_TEST_DO_WORK_CLI_BINARY")
    if supplied_cli:
        tools["supplied_cli"] = binary_seal(supplied_cli)
    return {"tools": tools, "go_environment": go_environment, "module_records": module_records, "git_configuration": git_configuration, "platform": tuple(os.uname())}


if __name__ == "__main__":
    try:
        encoded = json.dumps(runtime_seal(sys.argv[1], sys.argv[2:]), sort_keys=True).encode()
        print(hashlib.sha256(encoded).hexdigest())
    except (OSError, ValueError, KeyError, subprocess.CalledProcessError) as error:
        print("fingerprint uncertain: " + str(error), file=sys.stderr)
        raise SystemExit(1)
