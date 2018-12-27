package teamtailorgo

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/pkg/errors"
)

type JobApplication struct {
	ID             string
	Type           string `json: "type"`
	Created        string `json:"created-at"`
	CoverLetter    string `json:"cover-letter"`
	UpdatedAt      string `json:"updated-at"`
	RejectedAt     string `json:"rejected-at"`
	ReferringSite  string `json:"referring-site"`
	ReferringURL   string `json:"referring-url"`
	Sourced        bool   `json:"sourced"`
	ChangedStateAt string `json:"changed-stage-at"`
}

type JA struct {
	Data JAData `json:"data"`
}

type JAData struct {
	Type          string          `json:"type"`
	Attributes    JAAttributes    `json:"attributes"`
	Relationships JARelationships `json:"relationships"`
}

type JAAttributes struct {
	Sourced bool `json:"sourced"`
}

type JARelationships struct {
	Candidate JACandidate `json:"candidate"`
	Job       JAJob       `json:"job"`
}

type JACandidate struct {
	Data JACandidateData `json:"data"`
}

type JACandidateData struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type JAJob struct {
	Data JAJobData `json:"data"`
}

type JAJobData struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// CreateJobApplication
func (t TeamTailor) CreateJobApplication(idjob string, idcand string) (JobApplication, error) {

	cand := JACandidate{Data: JACandidateData{ID: idcand, Type: "candidates"}}
	job := JAJob{Data: JAJobData{ID: idjob, Type: "jobs"}}

	data := JA{Data: JAData{Type: "job-applications", Attributes: JAAttributes{Sourced: true}, Relationships: JARelationships{Candidate: cand, Job: job}}}

	var ja JobApplication
	json, err := json.Marshal(data)
	if err != nil {
		return ja, errors.Wrap(err, "marshalling failed")
	}

	postData := bytes.NewReader(json)

	req, _ := http.NewRequest("POST", baseURL+"job-applications", postData)
	req.Header.Set("Authorization", "Token token="+t.Token)
	req.Header.Set("X-Api-Version", apiVersion)
	req.Header.Set("Content-Type", contentType)

	resp, err := t.HTTPClient.Do(req)
	if err != nil {
		return ja, err
	}
	if resp.StatusCode != 201 {
		return js, errors.New("Failed to create job application")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ja, err
	}

	err = jsonapi.Unmarshal(body, &ja)
	if err != nil {
		return ja, err
	}

	defer resp.Body.Close()

	return ja, nil
}

// TODO: GetJobApplication

// TODO: GetJobApplicationsByCandidate

// TODO: GetJobApplicationsByJob

// TODO: GetJobApplicationStage

// JSONAPI Functions

func (ja *JobApplication) SetID(ID string) error {
	ja.ID = ID
	return nil
}

func (ja JobApplication) GetID() string {
	return ja.ID
}

func (ja JobApplication) SetToOneReferenceID(name, ID string) error {
	ja.ID = ID
	return nil
}

func (ja JobApplication) GetName() string {
	return "job-applications"
}