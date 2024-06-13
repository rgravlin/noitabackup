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

# How To Use
## First Use (Backup)
1. Download the latest release from this [GitHub Repo Releases](https://github.com/rgravlin/noitabackup/releases)
1. Create the default backup directory `%USERPROFILE%\NoitaBackup`
1. Launch `noitabackup.exe` in either GUI mode for ease of use
1. You must quit Noita `noita.exe` to execute `launch`, `backup`, and `restore` functions 
1. When Noita is not running, you must first execute a `backup`
  * This will copy `%USERPROFILE%\..\LocalLow\Nolla_Games_Noita\save00` to `%USERPROFILE%\NoitaBackup`
  * The timestamp will look like `2024-06-12-17-49-18` (these are your backups)

## Restore
1. Whenever you want to restore the _LATEST_ backup, quit Noita, and click `restore`
  * This will, assuming your base directory is `%USERPROFILE%\..\LocalLow\Nolla_Games_Noita\`:
    * Delete `%BASE%\save00.bak` (emergency backup that is rotated on every restore)
    * Rename `%BASE%\save00` to `%BASE%\save00.bak`
    * Copy the _LATEST_ backup to `%BASE%\save00`
    * Launch Noita if you have auto-launch enabled

## Advanced Use
### PowerShell - Alter default configuration

```commandline
PS C:\> $env:CONFIG_NOITA_DST_PATH = 'C:\Users\demo\backuptest'
PS C:\> .\noitabackup.exe backup
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
### Command Line - Alter default configuration
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