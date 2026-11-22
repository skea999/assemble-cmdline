# assemble-cmdline

Assembles the kernel command line from fragment files in a directory
(`/etc/kernel/cmdline.d` by default):

- files are read in **lexicographic order** (`00-base.conf` < `50-hardening.conf` < `99-debug.conf`)
- lines whose first non-blank character is `#` are dropped
- remaining lines are joined with single spaces — internal spacing is preserved

The result goes to **stdout**. Writing `/etc/kernel/cmdline` and rebuilding
the UKI is the caller's job (installer, kernel hook, apk trigger) — this
tool is pure: it never touches the filesystem.

## Usage

```
assemble-cmdline [dir]
```

Exit codes: `0` success, `1` I/O error, `2` usage error.

## Why

Package-shipped kernel cmdline fragments (e.g. a hardening package ships
its hardened parameters) must survive the whole chain:

1. the installer assembles `/etc/kernel/cmdline` from the fragments at install
2. the kernel hook rebuilds the UKI from the fragments at every kernel upgrade
3. an apk trigger reassembles + rebuilds the UKI when a fragment package upgrades

All three call this tool, so a fragment upgrade is picked up without a
post-install script.

## FHS

`/usr/sbin/assemble-cmdline` — system administration helper, executed by
root-owned scripts (kernel hook, apk trigger, installer). Per FHS 3.0,
`/usr/libexec` would be wrong: its binaries are not meant to be executed by
shell scripts, and this tool is.

## Build

```
go build -o /usr/sbin/assemble-cmdline .
```

Static binary (CGO disabled on target builds). No runtime dependencies.

## License

ISC — see `LICENSE`.
