package scheduler

/*
 Copyright 2019 - 2021 Crunchy Data Solutions, Inc.
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/qingcloud/postgres-operator/internal/kubeapi"
	log "github.com/sirupsen/logrus"

	cv3 "github.com/robfig/cron/v3"
	v1 "k8s.io/api/core/v1"
)

func New(label, namespace string, client kubeapi.Interface) *Scheduler {
	clientset = client
	cronClient := cv3.New()
	_, _ = cronClient.AddFunc("* * * * *", phony)
	_, _ = cronClient.AddFunc("* * * * *", heartbeat)

	return &Scheduler{
		namespace:  namespace,
		label:      label,
		CronClient: cronClient,
		entries:    make(map[string]cv3.EntryID),
	}
}

func (s *Scheduler) AddSchedule(config *v1.ConfigMap) error {
	name := config.Name + config.Namespace
	if _, ok := s.entries[name]; ok {
		return nil
	}

	if len(config.Data) != 1 {
		return errors.New("Schedule configmaps should contain only one schedule")
	}

	var schedule ScheduleTemplate
	for _, data := range config.Data {
		if err := json.Unmarshal([]byte(data), &schedule); err != nil {
			return fmt.Errorf("Failed unmarhsaling configMap: %w", err)
		}
	}

	if err := validate(schedule); err != nil {
		return fmt.Errorf("Failed to validate schedule: %w", err)
	}

	id, err := s.schedule(schedule)
	if err != nil {
		return fmt.Errorf("Failed to schedule configmap: %w", err)
	}

	log.WithFields(log.Fields{
		"configMap":  string(config.Name),
		"type":       schedule.Type,
		"schedule":   schedule.Schedule,
		"namespace":  schedule.Namespace,
		"deployment": schedule.Deployment,
		"label":      schedule.Label,
		"container":  schedule.Container,
	}).Info("Added new schedule")

	s.entries[name] = id
	return nil
}

func (s *Scheduler) DeleteSchedule(config *v1.ConfigMap) {
	log.WithFields(log.Fields{
		"scheduleName": config.Name,
	}).Info("Removed schedule")

	name := config.Name + config.Namespace
	s.CronClient.Remove(s.entries[name])
	delete(s.entries, name)
}

func (s *Scheduler) schedule(st ScheduleTemplate) (cv3.EntryID, error) {
	var job cv3.Job

	switch st.Type {
	case "pgbackrest":
		job = st.NewBackRestSchedule()
	case "policy":
		job = st.NewPolicySchedule()
	default:
		var id cv3.EntryID
		return id, fmt.Errorf("schedule type not implemented yet")
	}
	return s.CronClient.AddJob(st.Schedule, job)
}

// phony implements a no-op schedule job to prevent a bug that runs newly
// scheduled jobs multiple times
func phony() {
	_ = time.Now()
}

// heartbeat modifies a sentinel file used as part of the liveness test
// for the scheduler
func heartbeat() {
	// #nosec: G303
	err := ioutil.WriteFile("/tmp/scheduler.hb", []byte(time.Now().String()), 0o600)
	if err != nil {
		log.Errorln("error writing heartbeat file: ", err)
	}
}
