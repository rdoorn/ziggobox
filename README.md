Ziggo box library to interface with the ziggo connect box using http requests

basic usage added:

- login to connect box
- logout of connect box
- view global settings (access level)
- enabled previously added MAC to blocked list
- remove previously added MAC from blocked list

Usage:

```
z := New("http://url-to-connect-box") // url to ziggo box

// optionally enable debug logging (this will show the full request and reply including credentials!)
z.Debug(true)

// init sets up the initial sessiontoken
_, err := z.Init()

// do a proper login
err = z.Login("NULL", "password")
if err != nil {
  log.Fatal(err)
}

// get settings to see if we are logged in
res, err := z.GetGlobalSettings()
if err != nil {
  log.Fatal(err)
}
log.Printf("logged in: %t", res.AccessLevel == 1)

// deny/allow mac (only works on pre-configured macs, we don't add new macs or delete them)
err = z.AllowMac("00:00:00:00:01:02")
if err != nil {
  log.Fatal(err)
}

err = z.DenyMac("00:00:00:00:01:02")
if err != nil {
  log.Fatal(err)
}

// logout
z.Logout()
}
```

Restrictions:

This API uses the web interface, you can only log in to this web interface with 1 user at a time. As such, you need to be logged out of the web interface before you can use this api. and this api will error once you login, until your session expires, or you log out.
