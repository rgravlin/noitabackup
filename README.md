# Noita Backup and Restore (file) Utility

The goal of this project is to have an easy to use Noita backup and restore solution.  I work mostly
in Go, so that's why the language was chosen.  The UI was built with [GioUI](https://github.com/gioui/gio) and
the command line interface with [cobra](https://github.com/spf13/cobra) 

* Currently, Windows only!
* Basic GUI launcher features:
  * Backup and Restore (with auto-launch noita.exe)
  * Number of backups to keep
  * Open noita.exe
  * Explore Backups
* Supports custom source directory
* Supports custom destination directory
* Configured via environmental variables
* Command line interface

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
2024/06/12 07:14:16 timestamp: 2024-06-12-07-14-16
2024/06/12 07:14:16 source: C:\Users\demo\AppData\Roaming\..\LocalLow\Nolla_Games_Noita\save00
2024/06/12 07:14:16 destination: C:\Users\demo\backuptest\2024-06-12-07-14-16
2024/06/12 07:14:16 number of backups: 16
2024/06/12 07:14:16 maximum backup threshold reached
2024/06/12 07:14:16 removing backup folder: C:\Users\demo\backuptest\2024-06-11-07-33-30
2024/06/12 07:14:24 timestamp: 2024-06-12-07-14-24
2024/06/12 07:14:24 total time: 8.5259157s
2024/06/12 07:14:24 total dirs copied: 8
2024/06/12 07:14:24 total files copied: 10418
```
## Command Line - Alter default configuration
```commandline
C:\> set CONFIG_NOITA_DST_PATH='C:\Users\demo\backuptest'
C:\> noitabackup.exe backup
2024/06/12 07:14:16 timestamp: 2024-06-12-07-14-16
2024/06/12 07:14:16 source: C:\Users\demo\AppData\Roaming\..\LocalLow\Nolla_Games_Noita\save00
2024/06/12 07:14:16 destination: C:\Users\demo\NoitaBackups\2024-06-12-07-14-16
2024/06/12 07:14:16 number of backups: 16
2024/06/12 07:14:16 maximum backup threshold reached
2024/06/12 07:14:16 removing backup folder: C:\Users\demo\NoitaBackups\2024-06-11-07-33-30
2024/06/12 07:14:24 timestamp: 2024-06-12-07-14-24
2024/06/12 07:14:24 total time: 8.5259157s
2024/06/12 07:14:24 total dirs copied: 8
2024/06/12 07:14:24 total files copied: 10418
```