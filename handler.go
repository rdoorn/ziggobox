package ziggobox

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Handler struct {
	sid          string
	sessionToken string
	baseURL      string
	debug        bool
}

func New(baseURL string) *Handler {
	return &Handler{
		baseURL:      baseURL,
		sessionToken: "",
		debug:        false,
	}
}

func (h *Handler) Debug(v bool) {
	h.debug = v
}

func (h *Handler) Login(username, password string) error {
	body, err := h.call("/xml/setter.xml", "POST", []string{
		fmt.Sprintf("token=%s", h.sessionToken),
		"fun=15",
		fmt.Sprintf("Username=%s", username),
		fmt.Sprintf("Password=%s", password),
	})
	if h.debug {
		log.Printf("response: %s error: %s\n", body, err)
	}
	if err != nil {
		return err
	}
	if body == "" {
		return fmt.Errorf("Body returned nil, this can happen if you are logged in to the web interface, and did not log out.")
	}
	// successful;SID=167772160
	if strings.HasPrefix("successful", body) {
		h.sid = strings.Split(body, "=")[1]
		if h.debug {
			log.Printf("new SID: %s\n", h.sid)
		}
	}

	return nil
}

func (h *Handler) Logout() error {
	_, err := h.call("/xml/setter.xml", "POST", []string{
		fmt.Sprintf("token=%s", h.sessionToken),
		"fun=16",
	})
	return err
}

func (h *Handler) call(path, method string, parameters []string) (string, error) {
	/*if h.sessionToken == "" && path != "/" {
		h.call("/", "GET", map[string]string{})
	}*/
	//data := url.Values{}
	data := strings.Join(parameters, "&")
	u, _ := url.ParseRequestURI(h.baseURL)
	u.Path = path
	if h.debug {
		log.Printf("doing call to %s data: %v\n", u, data)
	}
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	r, _ := http.NewRequest(method, u.String(), strings.NewReader(data))

	cookie := ""
	if h.sessionToken != "" {
		cookie += fmt.Sprintf("sessionToken=%s", h.sessionToken)
	}
	if h.sid != "" {
		cookie += fmt.Sprintf(";SID=%s", h.sid)
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Accept", "*/*")
	if cookie != "" {
		r.Header.Add("Cookie", cookie)
	}

	resp, err := client.Do(r)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		if k == "Set-Cookie" {
			vs := strings.Split(v[0], ";")
			for _, vss := range vs {
				if strings.HasPrefix(vss, "sessionToken") {
					h.sessionToken = strings.Split(vss, "=")[1]
					if h.debug {
						log.Printf("new session token: %s\n", h.sessionToken)
					}
				}
			}
		}
	}
	//read body to return
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		return string(bodyBytes), nil
	}
	return "", err
}

type GlobalSettings struct {
	AccessLevel int
}

// GetGlobalSettings returns the global settings
func (h *Handler) GetGlobalSettings() (*GlobalSettings, error) {
	body, err := h.call("/xml/getter.xml", "POST", []string{
		fmt.Sprintf("token=%s", h.sessionToken),
		"fun=1",
	})
	if err != nil {
		return nil, err
	}

	res := &GlobalSettings{}
	err = xml.Unmarshal([]byte(body), res)
	if err != nil {
		return nil, err
	}
	return res, nil
	// returns <?xml version="1.0" encoding="utf-8"?><GlobalSettings><AccessLevel>0</AccessLevel><CmProvisionMode>IPv4</CmProvisionMode><GwProvisionMode>IPv4</GwProvisionMode><GWOperMode>IPv4</GWOperMode><DsLite>0</DsLite><PortControl>0</PortControl><OperatorId>ZIGGO</OperatorId><AccessDenied>NONE</AccessDenied><LockedOut>Disable</LockedOut><CountryID>7</CountryID><title>Connect Box</title><Interface>1</Interface><operStatus>1</operStatus></GlobalSettings>
}

// Init sets the initial sessionToken
func (h *Handler) Init() (*GlobalSettings, error) {
	_, err := h.call("/common_page/login.html", "GET", []string{})
	return nil, err
}

// AllowMac disabled an existing blocked MAC
func (h *Handler) AllowMac(addr string) error {
	data := fmt.Sprintf("EN,devicename,%s,2,1;MODE=0,TIME=0;", addr)
	body, err := h.call("/xml/setter.xml", "POST", []string{
		fmt.Sprintf("token=%s", h.sessionToken),
		"fun=120",
		fmt.Sprintf("data=%s", url.QueryEscape(data)),
	})
	if err != nil {
		return err
	}
	if body == "" {
		return nil
	}
	return fmt.Errorf("Failed to allow MAC: %s", body)
}

// DenyMac enable an existing blocked MAC
func (h *Handler) DenyMac(addr string) error {
	data := fmt.Sprintf("EN,devicename,%s,1,1;MODE=0,TIME=0;", addr)
	body, err := h.call("/xml/setter.xml", "POST", []string{
		fmt.Sprintf("token=%s", h.sessionToken),
		"fun=120",
		fmt.Sprintf("data=%s", url.QueryEscape(data)),
	})
	if err != nil {
		return err
	}
	if body == "" {
		return nil
	}
	return fmt.Errorf("Failed to deny MAC: %s", body)
}
