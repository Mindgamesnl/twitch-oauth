
# go twitch-oauth

A one-line util to fully manage twitch oauth on MacOS, Linux and Windows.


## Demo

[todo - add video]


## FAQ

#### How does it work?

twitch-oauth spawns a local webserver listening for redirects with oauth acknowedgements, requests an auth redirect and then opens it in a browser. The webserver will stop when the auth process complets (or fails) and then finish the method call.

#### Why a blocking call instead of channels?

Easier error handling and code eye-candy

#### Do I need to register my own application?

Yes, you need to register your own application in the twitch developer dashboard with a redirect callback to `http://localhost:7001/redirect`



## Usage/Examples

```go
func main() {
    // start login
    user, err := twitchauth.HandleOauth(
        "clientID",
        "clientSecret",
        []string{"user:read:email"}
    )

    if err != nil {
        panic(err)
    }

    fmt.Println(
		"Welcome! Email: " + user.Email +
			", ID: " + user.ID +
			", Username: " + user.Name +
			", Display Name: " + user.ReadableName +
			", Description: " + user.Description +
			", Profile Image: " + user.ProfilePicture,
	)
}
```

