package internal

// Strings
const (
	StrLatest = "latest"
)

// Errors
const (
	ErrMaxBackupsExceeded         = "maximum backup threshold reached"
	ErrInvalidBackups             = "max backups must be greater than zero"
	ErrOperationAlreadyInProgress = "operation already in progress"
	ErrNoitaRunning               = "noita.exe cannot be running"
	ErrErrorGettingBackups        = "error getting backups"
	ErrFailureDeletingBackups     = "failure deleting backups"
	ErrFailedToLaunch             = "failed to launch noita"
	ErrSourceDestination          = "cannot copy source to destination"
	ErrProcessingSave00           = "error processing save00"
	ErrBackupNotFound             = "backup %s not found in backup directory"
	ErrCannotCreateDestination    = "cannot create destination path"
	ErrDuringRestore              = "during restore"
	ErrDuringBackup               = "during backup"
	ErrFailedGettingBackupDirs    = "failed to get backup dirs"
	ErrNoBackupDirs               = "no backup dirs found, cannot restore"
	ErrRestoringToSave00          = "error restoring backup file to save00"
	ErrCopyingToSave00            = "error copying latest backup %s to save00: %v"
	ErrLaunchingExplorer          = "error launching explorer"
	ErrLaunchingNoita             = "error launching noita"
	ErrNumBackups                 = "number of backups to keep must be between 1 and 64"
	ErrGettingSourcePath          = "error getting source path"
	ErrGettingDestinationPath     = "error getting destination path"
	ErrSourcePathNotExist         = "source path does not exist"
	ErrDestinationPathNotExist    = "destination path does not exist"
	ErrSteamPathNotExist          = "steam path does not exist"
	ErrClosingFile                = "error closing file"
	ErrFailedCreateDir            = "failed to create directory: %s, error: %v"
	ErrRunningSteam               = "error running steam"
)

// Info
const (
	InfoNumberOfBackups   = "number of backups"
	InfoRemovingBackup    = "removing backup folder"
	InfoTimestamp         = "timestamp"
	InfoSource            = "source"
	InfoDestination       = "destination"
	InfoTotalTime         = "total time"
	InfoTotalDirCopied    = "total dirs copied"
	InfoTotalFileCopied   = "total files copied"
	InfoCreatingSave00    = "creating save00 directory"
	InfoCopyBackup        = "copying latest backup %s to save00"
	InfoSuccessfulRestore = "successfully restored backup"
	InfoDeletingSave00Bak = "deleting save00.bak folder"
	InfoRename            = "renaming save00 to save00.bak"
	InfoDebugLogSet       = "debug log set to"
	InfoAutoLaunchSet     = "auto-launch set to"
	InfoStartingRestore   = "starting restore"
	InfoStartingBackup    = "starting backup"
)

// Viper
const (
	ViperAutoLaunch      = "auto-launch"
	ViperNumBackups      = "num-backups"
	ViperSourcePath      = "source-path"
	ViperDestinationPath = "destination-path"
	ViperSteamPath       = "steam-path"
)

// Buttons
const (
	BtnLaunch  = "Launch Noita"
	BtnBackup  = "Backup Noita"
	BtnRestore = "Restore Noita"
	BtnExplore = "Explore Backups"
	BtnQuit    = "Quit"
)

// Checkbox
const (
	ChkAutoLaunch = "Auto Launch"
	ChkDebugLog   = "Debug Log"
)

// Slider
const (
	SldNumBackupsToKeep = "Number backups to keep"
)
