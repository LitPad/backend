package routes

import "github.com/LitPad/backend/models/choices"

func IsAmongContractStatus(target string) bool {
    switch target {
    case string(choices.CTS_PENDING), string(choices.CTS_UPDATED), string(choices.CTS_APPROVED), string(choices.CTS_DECLINED):
        return true
    }
    return false
}