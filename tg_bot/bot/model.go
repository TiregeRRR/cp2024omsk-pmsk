package bot

const (
	StatusUploaded = iota
	StatusTranscription
	StatusNers
	StatusReport
	StatusDone
	StatusFailed
)

const (
	START                   = "/start"
	START_TEXT              = "Здравствуйте. Для начала работы загрузите аудиофайл."
	NO_AUDIO_ATTACHED       = "Ошибка. Загрузите аудиофайл."
	NOT_SUPPORTED_TYPE      = "Ошибка. Загрузите аудиофайл формата: mp3, ogg или wav."
	FAILED_TO_DOWNLOAD_FILE = "Ошибка. Не получилось загрузить файл. Повторите попытку."
)
