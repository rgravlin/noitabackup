## noitabackup completion powershell

Generate the autocompletion script for powershell

### Synopsis

Generate the autocompletion script for powershell.

To load completions in your current shell session:

	noitabackup completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.


```
noitabackup completion powershell [flags]
```

### Options

```
  -h, --help              help for powershell
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
      --auto-launch               Auto-launch Noita after backup/restore operation
      --config string             config file (default is $HOME/.noitabackup.yaml)
      --destination-path string   Define the destination backup path (default "C:\\Users\\Demo\\NoitaBackups")
      --num-backups int           Define the maximum number of backups to keep (default 16)
      --source-path string        Define the source Noita save00 path (default "C:\\Users\\Demo\\AppData\\Roaming\\..\\LocalLow\\Nolla_Games_Noita\\save00")
      --steam-path string         Define the path for your Steam executable (default "C:\\Program Files (x86)\\Steam\\steam.exe")
```

### SEE ALSO

* [noitabackup completion](noitabackup_completion.md)	 - Generate the autocompletion script for the specified shell

###### Auto generated by spf13/cobra on 26-Jun-2024
