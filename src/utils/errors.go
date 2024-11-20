package utils


import "log"


const (
    Exit_listener_continue              = 0xc1
    Error_connection_db_continue        = 0xc2
    Error_deconnection_db_continue      = 0xc3
    Error_get_paper_data_continue       = 0xc4
    Error_reach_url_continue            = 0xc5
    Error_read_page_content             = 0xc6
    Error_inserting_document_continue   = 0xc7
    Error_downloading_document_continue = 0xc8
    Error_user_input_continue           = 0xc9
)


// check error and quit
func CheckErrorQuit(err error) {
    if err != nil {
        log.Fatalf("\033[31mError: %s\033[0m", err)
    }
}


// check error and raise flag
func CheckAlertError(err error, flag int, msg string, ac *AlertChannel) {
    if err != nil {
        SendAlert(flag, msg, ac)
    }
}