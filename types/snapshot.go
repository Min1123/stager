// Copyright © 2018 Tristian Celestin <tristian@redhat.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	urllib "net/url"
	"path"
	"strings"
	"time"
	"redhat.com/consulting/stager/util"
)

type SnapshotItem struct {
	pull_spec string
	repo_digest string
	image_id string
}

func NewSnapshotItem(scheme string, registry string, namespace string, image string, tag string, oauth_token string) (SnapshotItem, error) {
	// Setup client
	endpoint := fmt.Sprintf("%s://%s/v2/%s/%s/manifests/%s", scheme, registry, namespace, image, tag)
	url, err := urllib.Parse(endpoint)
	util.Check(err)
	client := http.Client{}
	client.Timeout = time.Second * 10
	request, err := http.NewRequest("GET", url.String(), nil)
	util.Check(err)
	request.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	if oauth_token != "" {
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", oauth_token))
	}

	// Make request
	response, err := client.Do(request)
	util.Check(err)
	if response.StatusCode != 200 {
		return SnapshotItem{}, fmt.Errorf("Failed to construct SnapshotItem: couldn't get repo digest for %s", url.String())
	}

	// Get repo digest
	image_url, err := urllib.Parse(response.Header.Get("Location"))
	util.Check(err)
	repo_digest := path.Base(image_url.String())

	// Get image id
	result := make(map[string]interface{})
	err = json.NewDecoder(response.Body).Decode(&result)
	util.Check(err)
	config, _ := result["config"].(map[string]interface{})
	id, _ := config["digest"].(string)


	pull_spec := strings.Join([]string{registry, namespace, image}, "/")
	pull_spec = strings.Join([]string{pull_spec, tag}, ":")
	item := SnapshotItem{pull_spec: pull_spec, repo_digest: repo_digest, image_id: id}
	log.Printf("Constructed snapshot item: {pull_spec: %s, repo_digest: %s, image_id: %s}", pull_spec, repo_digest, id)

	return item, nil;
}