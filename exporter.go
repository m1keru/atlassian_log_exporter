package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	jira "github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"go.uber.org/zap"
)

var layout = "2006-01-02T15:04:05.000-0700"

type SavedState struct {
	Offset        int       `json:"offset"`
	LastEventDate time.Time `json:"last_event_date"`
}

func saveState(state SavedState, filename string) error {
	jsonData, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, jsonData, 0644)
}

func loadState(filename string) (SavedState, error) {
	var state SavedState
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return state, nil
	}

	jsonData, err := os.ReadFile(filename)
	if err != nil {
		return state, err
	}

	err = json.Unmarshal(jsonData, &state)
	return state, err
}

func main() {
	jiraAPIEmail := flag.String("jira_api_email", os.Getenv("JIRA_API_EMAIL"), "Jira API Email, can be set with JIRA_API_EMAIL env variable")
	jiraAPIToken := flag.String("jira_api_token", os.Getenv("JIRA_API_TOKEN"), "Jira API Token, can be set with JIRA_API_TOKEN env variable")
	jiraAPIEndpoint := flag.String("jira_api_endpoint", os.Getenv("JIRA_API_ENDPOINT"), "Jira API Endpoint, can be set with JIRA_API_ENDPOINT env variable")
	debug := flag.Bool("debug", false, "Enable debug mode")

	flag.Parse()

	if *jiraAPIEndpoint == "" || *jiraAPIEmail == "" || *jiraAPIToken == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	if *debug {
		logger.Log(zap.InfoLevel, "Debug mode enabled")
		config := zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		logger, _ = config.Build()
	}

	log := logger.Sugar()

	client, err := jira.New(nil, *jiraAPIEndpoint)
	if err != nil {
		log.Errorf("Error creating Jira client", err)
	}

	client.Auth.SetBasicAuth(*jiraAPIEmail, *jiraAPIToken)

	stateFilename := "jira_state.json"
	state, err := loadState(stateFilename)
	if err != nil {
		log.Errorf("Error loading state: %v. Starting from beginning.", err)
		state = SavedState{
			Offset:        0,
			LastEventDate: time.Now().AddDate(-1, 0, 0).UTC(),
		}
	}
	endTime := time.Now().UTC()
	startTime := state.LastEventDate

	log.Infof("Getting records from %s to %s", state.LastEventDate, endTime)

	for {
		log.Debugf("Getting records from %s to %s, offset %d\n", state.LastEventDate, endTime, state.Offset)
		records, response, err := client.Audit.Get(context.Background(), &models.AuditRecordGetOptions{
			From:   state.LastEventDate.UTC(),
			To:     endTime,
			Filter: "",
		}, state.Offset, 1000)

		if err != nil {
			log.Fatal("Failed to get records from API", err, response)
		}

		if len(records.Records) != 0 && records.Offset == 0 {
			state.LastEventDate, err = time.Parse(layout, records.Records[0].Created)
			if err != nil {
				log.Errorf("Error parsing time: %v", err)
			}
			state.LastEventDate = state.LastEventDate.Add(time.Second).UTC()
		}

		for _, record := range records.Records {
			created, err := time.Parse(layout, record.Created)
			if err != nil {
				log.Errorf("Error parsing time: %v", err)
			}
			if created.Before(endTime) && created.Before(startTime) {
				log.Debugf("Record created date is out of range: %s", record.Created)
				continue
			}

			var changeValues = []string{}
			for _, value := range record.ChangedValues {
				changeValues = append(changeValues,
					fmt.Sprintf("fieldName: %s, changedFrom: %s, changedTo: %s", value.FieldName, value.ChangedFrom, value.ChangedTo),
				)
			}

			log.Infof(
				"id:%d,created:%s,remoteAddress:%s,authorKey:%s,category:%s,changedValues:%v,description:%s,eventSource:%s,summary:%s",
				record.ID,
				record.Created,
				record.RemoteAddress,
				record.AuthorKey,
				record.Category,
				changeValues,
				record.Description,
				record.EventSource,
				record.Summary,
			)

		}

		state.Offset += len(records.Records)
		err = saveState(state, stateFilename)
		if err != nil {
			log.Errorf("Error saving state: %v", err)
		}

		if len(records.Records) < 1000 {
			state.Offset = 0
			err = saveState(state, stateFilename)
			if err != nil {
				log.Errorf("Error saving state: %v", err)
			}
			break
		}
	}
}
