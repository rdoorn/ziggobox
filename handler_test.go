package ziggobox

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	z := New("http://url-to-connect-box") // url to ziggo box

	// init sets up the initial sessiontoken
	_, err := z.Init()

	// do a proper login
	err = z.Login("NULL", "password")
	if err != nil {
		assert.Failf(t, "login failed: %s", err.Error())
		return
	}

	// get settings to see if we are logged in
	res, err := z.GetGlobalSettings()
	if err != nil {
		assert.Failf(t, "get global settings failed: %s", err.Error())
		return
	}
	log.Printf("logged in: %t", res.AccessLevel == 1)

	// deny/allow mac (only works on pre-configured macs, we don't add new macs or delete them)
	err = z.AllowMac("00:00:00:00:01:02")
	assert.Nil(t, err)
	err = z.DenyMac("00:00:00:00:01:02")
	assert.Nil(t, err)

	// logout
	err = z.Logout()
	assert.Nil(t, err)
}
