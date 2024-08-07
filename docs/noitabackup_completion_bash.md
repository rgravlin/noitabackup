## noitabackup completion bash

Generate the autocompletion script for bash

### Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(noitabackup completion bash)

To load completions for every new session, execute once:

#### Linux:

	noitabackup completion bash > /etc/bash_completion.d/noitabackup

#### macOS:

	noitabackup completion bash > $(brew --prefix)/etc/bash_completion.d/noitabackup

You will need to start a new shell for this setup to take effect.


```
noitabackup completion bash
```

### Options

```
  -h, --help              help for bash
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
