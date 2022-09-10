
# go twitch-oauth

A one-line util to fully manage twitch oauth on MacOS, Linux and Windows.

## Demo

https://user-images.githubusercontent.com/10709682/189499070-ddc11e2b-31fa-46d0-b717-2fe42d48d9f8.mp4


## Usage/Examples

```go
func main() {
    // start login, block routine until the user finished or failed
    user, err := twitchauth.HandleOauth(
        "clientID",
        "clientSecret",
        []string{"user:read:email"}
    )

    // did the login fail, or user abort?
    if err != nil {
        panic(err)
    }

    // All good! print info
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


## FAQ

#### How does it work?

twitch-oauth spawns a local webserver listening for redirects with oauth acknowedgements, requests an auth redirect and then opens it in a browser. The webserver will stop when the auth process complets (or fails) and then finish the method call.

#### Why a blocking call instead of channels?

Easier error handling and code eye-candy

#### Do I need to register my own application?

Yes, you need to register your own application in the twitch developer dashboard with a redirect callback to `http://localhost:7001/redirect`
