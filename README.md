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

## PowerShell

```commandline
PS C:\> $env:CONFIG_NOITA_DST_PATH = 'C:\Your\Backup\Path'
PS C:\> $env:CONFIG_NOITA_SRC_PATH = 'D:\Your\Noita\Path\save00'
PS C:\> noitabackup
2023/09/19 08:12:01 Timestamp: 2023-09-19-08-11-59
2023/09/19 08:12:01 Total dirs copied: 8
2023/09/19 08:12:01 Total files copied: 2057
```
## Command Line
```commandline
PS C:\> $env:CONFIG_NOITA_DST_PATH = 'C:\Your\Backup\Path'
PS C:\> $env:CONFIG_NOITA_SRC_PATH = 'D:\Your\Noita\Path\save00'
PS C:\> noitabackup
2023/09/19 08:12:01 Timestamp: 2023-09-19-08-11-59
2023/09/19 08:12:01 Total dirs copied: 8
2023/09/19 08:12:01 Total files copied: 2057
```