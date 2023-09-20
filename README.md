# Noita Backup (file) Utility

* Currently, Windows only!
* Supports custom source directory
* Supports custom destination directory
* Configured via environmental variables

# Example Usage

If your Noita save files are in the default location:
`%USERPROFILE%\..\LocalLow\Nolla_Games_Noita\save00`

And you want to back up to `%USERPROFILE%\NoitaBackup`

You can simply download the latest release from the releases
page and run it.  I would recommend putting it in a Windows
path, so you can run it from anywhere.

If you'd like to configure the source or destination please
use the following:

## PowerShell - Alter default configuration

```commandline
PS C:\> $env:CONFIG_NOITA_DST_PATH = 'C:\Users\demo\backuptest'
PS C:\> .\noitabackup.exe
2023/09/20 08:55:39 Destination: C:\Users\demo\backuptest\2023-09-20-08-55-39
2023/09/20 08:55:39 Timestamp: 2023-09-20-08-55-39
2023/09/20 08:55:39 Total dirs copied: 8
2023/09/20 08:55:39 Total files copied: 2057
```
## Command Line - Alter default configuration
```commandline
C:\> set CONFIG_NOITA_DST_PATH='C:\Users\demo\backuptest'
C:\> noitabackup.exe
2023/09/20 08:58:06 Destination: C:\Users\demo\backuptest\2023-09-20-08-58-05
2023/09/20 08:58:06 Total dirs copied: 8
2023/09/20 08:58:06 Total files copied: 2078
```