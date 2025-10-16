package cookies

import (
	"fmt"
	"net/http"
	"time"
	"log"
)

type Cookie struct {
	Name  string
	Value string

	Path       string    // optional
	Domain     string    // optional
	Expires    time.Time // optional
	RawExpires string    // for reading cookies only

	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// MaxAge>0 means Max-Age attribute present and given in seconds
	MaxAge   int
	Secure   bool
	HttpOnly bool
	SameSite http.SameSite
	Raw      string
	Unparsed []string // Raw text of unparsed attribute-value pairs
}

func setCookie(w http.ResponseWriter, req *http.Request) { // set cookie for storing token
	cookie := http.Cookie{}
	cookie.Name = "struggleBussToken"
	cookie.Domain = "games.bitsmasher.net"
	cookie.Value = "ro8BS6Hiivgzy8Xuu09JDjlNLnSLldY5"
	cookie.MaxAge = 3600 * 24 * 365
	cookie.Expires = time.Now().Add(365 * 24 * time.Hour)
	cookie.Secure = true // TLS only
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteNoneMode
	cookie.Path = "/"
	log.Println("setting client cookie")
	http.SetCookie(w, &cookie)

	// cookie2 := http.Cookie{}
	// cookie2.Name = "page"
	// cookie2.Value = "GoLinuxCloud"
	// cookie2.Expires = time.Now().Add(365 * 24 * time.Hour)
	// cookie2.Secure = false
	// cookie2.HttpOnly = true
	// cookie2.Path = "/"
	// http.SetCookie(w, &cookie2)

	// fmt.Fprintf(w, "This is cookies!\n")
}

func getCookie(w http.ResponseWriter, req *http.Request) bool {
	var returnStr string

	log.Println("getting client cookie")

	for _, cookie := range req.Cookies() {
		returnStr = returnStr + cookie.Name + ":" + cookie.Value + "\n"
		log.Println(returnStr + cookie.Name + ":" + cookie.Value + "\n")
	}

	if returnStr != "" {
		//fmt.Fprintf(w, returnStr)
		fmt.Fprintf(w, returnStr)
		return true
	} else {
		return false
	}
}
