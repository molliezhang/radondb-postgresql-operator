package api

/*
 Copyright 2018 - 2021 Qingcloud Data Solutions, Inc.
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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	msgs "github.com/qingcloud/postgres-operator/pkg/apiservermsgs"
	log "github.com/sirupsen/logrus"
)

// DeleteBackup  makes an API requests to delete a pgBackRest backup
func DeleteBackup(httpclient *http.Client, SessionCredentials *msgs.BasicAuthCredentials, request msgs.DeleteBackrestBackupRequest) (msgs.DeleteBackrestBackupResponse, error) {
	var response msgs.DeleteBackrestBackupResponse

	// explicitly set the client version here
	request.ClientVersion = msgs.PGO_VERSION

	log.Debugf("DeleteBackup called [%+v]", request)

	ctx := context.TODO()
	jsonValue, _ := json.Marshal(request)
	url := SessionCredentials.APIServerURL + "/backrest"

	action := "DELETE"
	req, err := http.NewRequestWithContext(ctx, action, url, bytes.NewBuffer(jsonValue))
	if err != nil {
		return response, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(SessionCredentials.Username, SessionCredentials.Password)

	resp, err := httpclient.Do(req)
	if err != nil {
		return response, err
	}

	defer resp.Body.Close()

	log.Debugf("%+v", resp)

	if err := StatusCheck(resp); err != nil {
		return response, err
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Print("Error: ")
		fmt.Println(err)
		return response, err
	}

	return response, nil
}

func ShowBackrest(httpclient *http.Client, arg, selector string, SessionCredentials *msgs.BasicAuthCredentials, ns string) (msgs.ShowBackrestResponse, error) {
	var response msgs.ShowBackrestResponse
	url := SessionCredentials.APIServerURL + "/backrest/" + arg + "?version=" + msgs.PGO_VERSION + "&selector=" + selector + "&namespace=" + ns

	log.Debugf("show backrest called [%s]", url)

	action := "GET"
	req, err := http.NewRequest(action, url, nil)
	if err != nil {
		return response, err
	}

	req.SetBasicAuth(SessionCredentials.Username, SessionCredentials.Password)
	resp, err := httpclient.Do(req)
	if err != nil {
		fmt.Println("Error: Do: ", err)
		return response, err
	}
	defer resp.Body.Close()
	log.Debugf("%v", resp)
	err = StatusCheck(resp)
	if err != nil {
		return response, err
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Debugf("%v", resp.Body)
		log.Debug(err)
		return response, err
	}

	return response, err
}

func CreateBackrestBackup(httpclient *http.Client, SessionCredentials *msgs.BasicAuthCredentials, request *msgs.CreateBackrestBackupRequest) (msgs.CreateBackrestBackupResponse, error) {
	var response msgs.CreateBackrestBackupResponse

	jsonValue, _ := json.Marshal(request)

	url := SessionCredentials.APIServerURL + "/backrestbackup"

	log.Debugf("create backrest backup called [%s]", url)

	action := "POST"
	req, err := http.NewRequest(action, url, bytes.NewBuffer(jsonValue))
	if err != nil {
		return response, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(SessionCredentials.Username, SessionCredentials.Password)

	resp, err := httpclient.Do(req)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	log.Debugf("%v", resp)
	err = StatusCheck(resp)
	if err != nil {
		return response, err
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Printf("%v\n", resp.Body)
		fmt.Println("Error: ", err)
		log.Println(err)
		return response, err
	}

	return response, err
}
