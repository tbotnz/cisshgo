package fakedevices

import (
	"bytes"
	"log"
	"text/template"
)

// TranscriptReader parses a transcript file and populates any variables that may exist in it
func TranscriptReader(transcript string, fakeDevice *FakeDevice) (string, error) {

	// Setup a template with our transcript
	tmpl, err := template.New("fakeDeviceTemplate").Parse(transcript)
	if err != nil {
		log.Fatal(err)
	}

	// Setup a bytes buffer to accept the rendered template
	var renderedTemplate bytes.Buffer

	// Render (Execute) the template with our input
	if err := tmpl.Execute(&renderedTemplate, fakeDevice); err != nil {
		log.Fatal(err)
	}

	return renderedTemplate.String(), nil
}
