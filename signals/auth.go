package signals

import "github.com/phonkee/patrol/models"

/* OnSuccessfulLoginSignalHandler
This signal is called on succesfull login
*/
type OnSuccessfulLoginSignalHandler interface {
	// when event request is received
	OnSuccessfulLogin(user *models.User)
}
