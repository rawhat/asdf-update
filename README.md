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
