package core



func StartEngine() {
	logChannel := CreateLogChannel()
	go ListenerLogFile(logChannel)

	eprint := InitEprint(logChannel)
	DownloadEprint(eprint, logChannel)

	close(logChannel.logChannel)
}