// Package fakedevices provides types and initialization for simulated
// network devices used by cisshgo.
package fakedevices

import (
	"bytes"
	"text/template"
)

// TranscriptReader parses a transcript file and populates any variables that may exist in it
func TranscriptReader(transcript string, fd *FakeDevice) (string, error) {

	// Setup a template with our transcript
	tmpl, err := template.New("fakeDeviceTemplate").Parse(transcript)
	if err != nil {
		return "", err
	}

	// Setup a bytes buffer to accept the rendered template
	var renderedTemplate bytes.Buffer

	// Render (Execute) the template with our input
	if err := tmpl.Execute(&renderedTemplate, fd); err != nil {
		return "", err
	}

	return renderedTemplate.String(), nil
}
