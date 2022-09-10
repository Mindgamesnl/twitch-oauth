package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/twitch"
	"log"
	"net/http"
	"os/exec"
	"runtime"
)

var (
	redirectURL    = "http://localhost:7001/redirect"
	oauth2Config   *oauth2.Config
	resultWait     = make(chan string)
	failurePointer error

	finishedHtml = "<style>*{\n    transition: all 0.6s;\n}\n\nhtml {\n    height: 100%;\n}\n\nbody{\n    font-family: 'Lato', sans-serif;\n    color: #888;\n    margin: 0;\n}\n\n#main{\n    display: table;\n    width: 100%;\n    height: 100vh;\n    text-align: center;\n}\n\n.fof{\n\t  display: table-cell;\n\t  vertical-align: middle;\n}\n\n.fof h1{\n\t  font-size: 50px;\n\t  display: inline-block;\n\t  padding-right: 12px;\n\t  animation: type .5s alternate infinite;\n}\n\n@keyframes type{\n\t  from{box-shadow: inset -3px 0px 0px #888;}\n\t  to{box-shadow: inset -3px 0px 0px transparent;}\n}</style><div id=\"main\">\n    \t<div class=\"fof\">\n        \t\t<h1>Login handled - please close this page</h1>\n    \t</div>\n</div>"
)

func openbrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err

}

func HandleOAuth2Callback(w http.ResponseWriter, r *http.Request) (err error) {
	token, err := oauth2Config.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, finishedHtml)

	resultWait <- token.AccessToken
	return
}

// HumanReadableError represents error information
// that can be fed back to a human user.
//
// This prevents internal state that might be sensitive
// being leaked to the outside world.
type HumanReadableError interface {
	HumanError() string
	HTTPCode() int
}

// HumanReadableWrapper implements HumanReadableError
type HumanReadableWrapper struct {
	ToHuman string
	Code    int
	error
}

func (h HumanReadableWrapper) HumanError() string { return h.ToHuman }
func (h HumanReadableWrapper) HTTPCode() int      { return h.Code }

// AnnotateError wraps an error with a message that is intended for a human end-user to read,
// plus an associated HTTP error code.
func AnnotateError(err error, annotation string, code int) error {
	if err == nil {
		return nil
	}
	return HumanReadableWrapper{ToHuman: annotation, error: err}
}

type Handler func(http.ResponseWriter, *http.Request) error

func HandleOauth(clientID string, clientSecret string, scopes []string) (*TwitchUser, error) {
	// Gob encoding for gorilla/sessions
	gob.Register(&oauth2.Token{})

	oauth2Config = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       scopes,
		Endpoint:     twitch.Endpoint,
		RedirectURL:  redirectURL,
	}

	var middleware = func(h Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) (err error) {
			// parse POST body, limit request size
			if err = r.ParseForm(); err != nil {
				return AnnotateError(err, "Something went wrong! Please try again.", http.StatusBadRequest)
			}

			return h(w, r)
		}
	}

	var errorHandling = func(handler func(w http.ResponseWriter, r *http.Request) error) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := handler(w, r); err != nil {
				var errorString = "Something went wrong! Please try again."
				var errorCode = 500

				if v, ok := err.(HumanReadableError); ok {
					errorString, errorCode = v.HumanError(), v.HTTPCode()
				}

				log.Println(err)
				w.Write([]byte(errorString))
				w.WriteHeader(errorCode)
				return
			}
		})
	}

	srv := &http.Server{Addr: ":7001"}

	var handleFunc = func(path string, handler Handler) {
		http.Handle(path, errorHandling(middleware(handler)))
	}
	handleFunc("/redirect", HandleOAuth2Callback)

	openErr := openbrowser(oauth2Config.AuthCodeURL("emlinksetup"))

	if openErr != nil {
		failurePointer = openErr
	}

	go func() {
		httpError := srv.ListenAndServe()
		if httpError != nil {
			failurePointer = httpError
		}
	}()

	token := ""
	if &failurePointer != nil {
		token = <-resultWait
	}

	err := srv.Shutdown(context.Background())
	if err != nil {
		return nil, err
	}

	if len(token) == 0 {
		eS := failurePointer
		failurePointer = nil
		return nil, eS
	}

	user, errU := GetUser(token, clientID, clientSecret)

	if errU != nil {
		return nil, errU
	}

	return user, nil
}
