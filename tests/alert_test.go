package test


import (
	"time"
	"testing"
	"github.com/Bl4omArchie/ePrint-DB/src/utils"
)


func TestAlertChannel(t *testing.T) {
	ac := utils.CreateAlertChannel()
	go utils.ListenerAlertChannel(ac)
	utils.SendAlert(0xc2, "Can't download PDF n°497", ac)

	time.Sleep(2 * time.Second)
	utils.CloseChannel(ac)
}