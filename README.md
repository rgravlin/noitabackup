# Noita Backup and Restore (file) Utility
![Noita Backup and Restore](images/noita-backup-restore-gui-default.png)

The goal of this project is to have an easy to use Noita backup and restore solution.  I work mostly
in Go, so that's why the language was chosen.  The UI was built with [GioUI](https://github.com/gioui/gio) and
the command line interface with [cobra](https://github.com/spf13/cobra) 

* Faster than Windows copy/paste!
* Currently, Windows only
* Command line interface
* Supports custom source directory
* Supports custom destination directory
* Configure via environmental variables and config file
* Protects against corrupted backups
* GUI launcher features:
  * Backup and Restore
  * Auto-Launch after backup/restore
  * Open Noita
  * Explore Backups
  * UI Debug Log

# How To Use

## First Use (Backup)
1. Download the latest release from this [GitHub Repo Releases](https://github.com/rgravlin/noitabackup/releases)
1. Create the default backup directory `%USERPROFILE%\NoitaBackup`
1. Launch `noitabackup.exe` in either GUI mode for ease of use
1. You must quit Noita `noita.exe` to execute `launch`, `backup`, and `restore` functions 
1. When Noita is not running, you must first execute a `backup`
  * This will copy `%APPDATA%\..\LocalLow\Nolla_Games_Noita\save00` to `%USERPROFILE%\NoitaBackup`
  * The timestamp will look like `2024-06-12-17-49-18` (these are your backups)

## Restore
1. Whenever you want to restore the _LATEST_ backup, quit Noita, and click `restore`
  * This will, assuming your base directory is `%APPDATA%\..\LocalLow\Nolla_Games_Noita\`:
    * Delete `%BASE%\save00.bak` (emergency backup that is rotated on every restore)
    * Rename `%BASE%\save00` to `%BASE%\save00.bak`
    * Copy the _LATEST_ backup to `%BASE%\save00`
    * Launch Noita if you have auto-launch enabled

## Advanced Use
### Configuration Parameters

The default configuration file is looked for in `$HOME/.noitabackup.yaml` and has the same configuration parameters as
the CLI application.

| Name               | Description                                            | Value                                            |
|--------------------|--------------------------------------------------------|--------------------------------------------------|
| `auto-launch`      | Auto-launch Noita after backup or restore              | `false`                                          |
| `num-backups`      | Total number of backups to keep                        | `16`                                             |
| `num-workers`      | Total number of Go routines to process copy operations | `4`                                              |
| `source-path`      | Source Noita save game path                            | `%APPDATA%\..\LocalLow\Nolla_Games_Noita\save00` |
| `destination-path` | Destination main backup path                           | `%USERPROFILE%\NoitaBackups`                     |
| `steam-path`       | Steam executable path                                  | `C:\Program Files (x86)\Steam\steam.exe`         |

### Configuration Example
```yaml
---
auto-launch: 'false'
num-backups: 16
num-workers: 4
source-path: C:\\Users\\Demo\\AppData\\LocalLow\\Nolla_Games_Noita\\save00
destination-path: C:\\Users\\Demo\\NoitaBackups
steam-path: C:\\Program Files (x86)\\Steam\\steam.exe
```

### PowerShell - Alter default configuration
```bash
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
```bash
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

## Screenshots
### Noita restore process
![Noita Restore Process](images/noita-backup-restore-gui-restore.png)
### Noita is running (red background)
![Noita UI Noita Running Red Background](images/noita-backup-autolaunch-noitarunning.png)
### Enable Debug log
![Noita UI Debug Log](images/noita-backup-restore-gui-debuglog-default.png)
