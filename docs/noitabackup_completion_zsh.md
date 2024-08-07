## noitabackup completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(noitabackup completion zsh)

To load completions for every new session, execute once:

#### Linux:

	noitabackup completion zsh > "${fpath[1]}/_noitabackup"

#### macOS:

	noitabackup completion zsh > $(brew --prefix)/share/zsh/site-functions/_noitabackup

You will need to start a new shell for this setup to take effect.


```
noitabackup completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
      --auto-launch               auto-launch Noita after backup/restore operation
      --config string             config file (default is $HOME/.noitabackup.yaml)
      --destination-path string   destination backup path (default "C:\\Users\\Demo\\NoitaBackups")
      --num-backups int           maximum number of backups to keep (default 16)
      --num-workers int           total number of go routine workers (advanced usage) (default 4)
      --source-path string        source Noita save00 path (default "C:\\Users\\Demo\\AppData\\LocalLow\\Nolla_Games_Noita\\save00")
      --steam-path string         path for your Steam executable (default "C:\\Program Files (x86)\\Steam\\steam.exe")
```

### SEE ALSO

* [noitabackup completion](noitabackup_completion.md)	 - Generate the autocompletion script for the specified shell

###### Auto generated by spf13/cobra on 21-Jul-2024
