package core



func StartEngine() {
	errChannel := CreateErrorChannel()
	go ListenerLogFile(errChannel)

	eprint := InitEprint(errChannel)
	DownloadEprint(eprint, errChannel)

	close(errChannel.errChannel)
}