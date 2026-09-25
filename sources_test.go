package main

import (
	"bytes"
	"encoding/xml"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

// sourceUUID is the UUID of the source that is uploaded in the tests.
const sourceUUID = "44444444-4444-4444-4444-444444444444"

// testSourceXML is a minimal source data set with a digital file reference and
// a version that needs to be normalized.
const testSourceXML = `<sourceDataSet>
  <sourceInformation>
    <dataSetInformation>
      <UUID>` + sourceUUID + `</UUID>
      <shortName>ACME source</shortName>
      <referenceToDigitalFile uri="docs/report%20final.pdf"/>
    </dataSetInformation>
  </sourceInformation>
  <administrativeInformation>
    <publicationAndOwnership>
      <dataSetVersion>1</dataSetVersion>
    </publicationAndOwnership>
  </administrativeInformation>
</sourceDataSet>`

// TestDigitalFile checks the version normalization and the digitalfile route.
func TestDigitalFile(t *testing.T) {
	server := newTestServer(t)
	client := login(t, server)

	uploadSource(t, client, server, "report final.pdf", "the pdf content")

	// the version of the uploaded source is normalized
	list := sourceList(t, client, server)
	if len(list) != 1 {
		t.Fatalf("expected 1 source but got %d", len(list))
	}
	if list[0].Version != "01.00.000" {
		t.Errorf("expected version 01.00.000 but got %q", list[0].Version)
	}

	// a non-normalized version can be used to read the data set
	assertStatus(t, get(t, client, server,
		"/resource/sources/"+sourceUUID+"?version=1", ""), http.StatusOK)

	// the first digital file of the source is returned
	resp := get(t, client, server,
		"/resource/sources/"+sourceUUID+"/digitalfile", "")
	assertStatus(t, resp, http.StatusOK)
	if ct := resp.Header.Get("Content-Type"); ct != "application/pdf" {
		t.Errorf("expected application/pdf but got %q", ct)
	}
	if body := textOf(t, resp); body != "the pdf content" {
		t.Errorf("expected the digital file content but got %q", body)
	}

	// an unknown file leads to a 404
	assertStatus(t, get(t, client, server,
		"/resource/sources/"+sourceUUID+"/no_such_file.pdf", ""),
		http.StatusNotFound)
}

// TestFileNameOf checks how file names are extracted from digital file
// references.
func TestFileNameOf(t *testing.T) {
	check := func(uri string, expected string) {
		if got := fileNameOf(uri); got != expected {
			t.Errorf("expected %q for %q but got %q", expected, uri, got)
		}
	}
	check("report.pdf", "report.pdf")
	check("docs/report%20final.pdf", "report final.pdf")
	check(`C:\docs\report.pdf`, "report.pdf")
	check("http://example.com/files/a.pdf", "a.pdf")
	check("", "")
}

// uploadSource posts a source data set with one external file.
func uploadSource(t *testing.T, client *http.Client, server *httptest.Server,
	fileName string, content string) {
	t.Helper()
	body := &bytes.Buffer{}
	form := multipart.NewWriter(body)
	if err := form.WriteField("file", testSourceXML); err != nil {
		t.Fatal(err)
	}
	if err := form.WriteField(fileName, content); err != nil {
		t.Fatal(err)
	}
	if err := form.Close(); err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest(http.MethodPost,
		server.URL+"/resource/sources/withBinaries", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", form.FormDataContentType())
	assertStatus(t, send(t, client, req), http.StatusOK)
}

// sourceList reads the sources of the root data stock.
func sourceList(t *testing.T, client *http.Client,
	server *httptest.Server) []SourceDescriptor {
	t.Helper()
	resp := get(t, client, server, "/resource/sources", "")
	assertStatus(t, resp, http.StatusOK)
	list := &DescriptorList{}
	if err := xml.Unmarshal([]byte(textOf(t, resp)), list); err != nil {
		t.Fatal(err)
	}
	return list.Sources
}
