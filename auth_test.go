package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

// newTestServer creates a server that runs with the given users in a temporary
// data folder. When no user is given, the default admin user is created.
func newTestServer(t *testing.T, users ...User) *httptest.Server {
	t.Helper()
	args := Args{"-data": t.TempDir()}
	config := &Config{Users: users}
	if len(users) == 0 {
		if _, err := EnsureDefaultAdmin(args, config); err != nil {
			t.Fatal(err)
		}
	}
	if err := WriteConfig(args, config); err != nil {
		t.Fatal(err)
	}
	tokens, err := initTokenAuth(args)
	if err != nil {
		t.Fatal(err)
	}
	dir, err := newDataDir(args.DataDir())
	if err != nil {
		t.Fatal(err)
	}
	s := &server{
		config:  config,
		cookies: initCookieStore(args),
		tokens:  tokens,
		dir:     dir,
	}
	router := mux.NewRouter()
	s.registerRoutes(router)
	server := httptest.NewServer(router)
	t.Cleanup(server.Close)
	return server
}

// TestAnonymousAccess checks that reading is open while writing needs a user.
func TestAnonymousAccess(t *testing.T) {
	server := newTestServer(t)
	client := server.Client()

	// reading is allowed without credentials
	assertStatus(t, get(t, client, server, "/resource/datastocks", ""),
		http.StatusOK)

	// writing requires credentials
	resp := postXML(t, client, server, "/resource/contacts", "", testContactXML)
	assertStatus(t, resp, http.StatusForbidden)
	assertText(t, resp, "Permission denied.")

	// invalid credentials are rejected on all routes
	assertStatus(t, get(t, client, server, "/resource/datastocks", "no.token"),
		http.StatusForbidden)
	assertStatus(t, postXML(t, client, server, "/resource/contacts",
		"no.token", testContactXML), http.StatusForbidden)
}

// TestTokenAuthentication checks the token based authentication.
func TestTokenAuthentication(t *testing.T) {
	server := newTestServer(t)
	client := server.Client()

	// request a token with the credentials of the default admin
	resp := postForm(t, client, server, "/resource/authenticate/getToken",
		url.Values{"username": {"admin"}, "password": {"default"}})
	assertStatus(t, resp, http.StatusOK)
	token := strings.TrimSpace(textOf(t, resp))
	if parts := strings.Split(token, "."); len(parts) != 3 {
		t.Fatalf("expected a token with 3 parts but got %q", token)
	}

	// the token can be used for reading and writing
	assertStatus(t, get(t, client, server, "/resource/datastocks", token),
		http.StatusOK)
	assertStatus(t, postXML(t, client, server, "/resource/contacts",
		token, testContactXML), http.StatusOK)

	// a wrong password does not give a token
	resp = postForm(t, client, server, "/resource/authenticate/getToken",
		url.Values{"username": {"admin"}, "password": {"wrong"}})
	assertStatus(t, resp, http.StatusUnauthorized)
	assertText(t, resp, "permission denied")
}

// TestSessionAuthentication checks the session cookie based authentication.
func TestSessionAuthentication(t *testing.T) {
	server := newTestServer(t)
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{Jar: jar}

	// login with the credentials of the default admin
	resp := postForm(t, client, server, "/resource/authenticate/login",
		url.Values{"username": {"admin"}, "password": {"default"}})
	assertStatus(t, resp, http.StatusOK)
	assertText(t, resp, "Login successful")

	// the session cookie allows writing
	assertStatus(t, postXML(t, client, server, "/resource/contacts",
		"", testContactXML), http.StatusOK)

	// and the status reports the logged in user
	resp = get(t, client, server, "/resource/authenticate/status", "")
	info := &struct {
		IsAuthenticated bool   `xml:"authenticated"`
		UserName        string `xml:"userName"`
	}{}
	if err := xml.Unmarshal([]byte(textOf(t, resp)), info); err != nil {
		t.Fatal(err)
	}
	if !info.IsAuthenticated || info.UserName != "admin" {
		t.Errorf("expected the admin user in the status but got %+v", info)
	}
}

// TestConfiguredToken checks tokens that are configured for a user.
func TestConfiguredToken(t *testing.T) {
	server := newTestServer(t, User{
		Name:     "alice",
		Password: "secret",
		Tokens:   []string{"alice-token"},
	})
	client := server.Client()

	assertStatus(t, get(t, client, server, "/resource/datastocks", "alice-token"),
		http.StatusOK)
	assertStatus(t, get(t, client, server, "/resource/datastocks", "bob-token"),
		http.StatusForbidden)
}

// TestTokenExpiry checks that generated tokens expire.
func TestTokenExpiry(t *testing.T) {
	auth, err := initTokenAuth(Args{"-data": t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	user := &User{Name: "admin"}

	token, err := auth.Generate(user, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if name, err := auth.Verify(token); err != nil || name != "admin" {
		t.Errorf("expected a valid token for admin but got %q, %v", name, err)
	}

	expired, err := auth.Generate(user, time.Now().Add(-2*tokenTTL))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := auth.Verify(expired); err == nil {
		t.Error("expected an expired token to be rejected")
	}

	if _, err := auth.Verify("not.a.token"); err == nil {
		t.Error("expected an invalid token to be rejected")
	}
}

// testContactXML is a minimal contact that can be posted to the server.
const testContactXML = `<contactDataSet>
  <contactInformation>
    <dataSetInformation>
      <UUID>11111111-1111-1111-1111-111111111111</UUID>
      <name>ACME</name>
    </dataSetInformation>
  </contactInformation>
  <administrativeInformation>
    <publicationAndOwnership>
      <dataSetVersion>01.00.000</dataSetVersion>
    </publicationAndOwnership>
  </administrativeInformation>
</contactDataSet>`

// get sends a GET request with an optional authentication token.
func get(t *testing.T, client *http.Client, server *httptest.Server,
	path string, token string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, server.URL+path, nil)
	if err != nil {
		t.Fatal(err)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return send(t, client, req)
}

// postForm sends a POST request with the given form values.
func postForm(t *testing.T, client *http.Client, server *httptest.Server,
	path string, form url.Values) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, server.URL+path,
		strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return send(t, client, req)
}

// postXML sends a POST request with the given XML data and an optional token.
func postXML(t *testing.T, client *http.Client, server *httptest.Server,
	path string, token string, data string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, server.URL+path,
		strings.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return send(t, client, req)
}

func send(t *testing.T, client *http.Client, req *http.Request) *http.Response {
	t.Helper()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { resp.Body.Close() })
	return resp
}

func assertStatus(t *testing.T, resp *http.Response, expected int) {
	t.Helper()
	if resp.StatusCode != expected {
		t.Errorf("expected status %d but got %d (%s)",
			expected, resp.StatusCode, textOf(t, resp))
	}
}

func assertText(t *testing.T, resp *http.Response, expected string) {
	t.Helper()
	text := textOf(t, resp)
	if !strings.Contains(text, expected) {
		t.Errorf("expected %q in %q", expected, text)
	}
}

func textOf(t *testing.T, resp *http.Response) string {
	t.Helper()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
