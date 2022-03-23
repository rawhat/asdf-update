# asdf-update

This super simple program will update all installed [asdf](https://asdf-vm.com/) plugins to their latest version.

It is essentially the same as you running:

```bash
asdf current $plugin
asdf latest $plugin
asdf install $plugin $latest
asdf global $plugin $latest
if [ -n "$current" ]; then
  asdf uninstall $plugin $current
fi
```

but for all of your installed plugins.  Probably could/should have done this in
`bash` but it's 2022.

## Supported Options

  * `-ignore` will skip updating the named plugin.  It can be repeated to skip
    multiple plugins
  * `-only` will update _only_ the named plugin.  It can be repeated to update
    some subset of plugins

Ignore takes precedence, and you can include both of them, technically.
¯\_(ツ)_/¯

Examples:

```bash
$ asdf-update -ignore bazel -ignore nodejs
Updating ...
Ignoring bazel
Updating ...
$ asdf-update -only golang
Updating golang
```
