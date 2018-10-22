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

package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"log"
	"os"
	"redhat.com/consulting/stager/types"
	"redhat.com/consulting/stager/util"
	"strings"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Args: cobra.NoArgs,
	Short: "Moves images from a registry to a staging directory",
	Long: `Moves images from a registry to a staging directory. Given a list of pull specs, The process
produces:

- a set of container images that are exported to a directory

- a file with a list of repo and image digests cooresponding to recently exported container images`,
	Run: run,
}



func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	getCmd.Flags().StringP("filename", "f", "", "File containing list of images to get")
	getCmd.Flags().StringP("docker-host", "d", "unix://var/run/docker.sock", "Path to Docker socket")
}

func run(command *cobra.Command, invocation [] string) {
	fmt.Println("get called")
	connectToDockerDaemon()
	generateSnapshots(command.Flag("filename").Value.String())
}

func connectToDockerDaemon() (*client.Client) {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		log.Println(err)
		panic(err)
	}
	//containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	//images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	//if err != nil {
	//	panic(err)
	//}

	return cli;
}

func generateSnapshots(filename string) ([]types.SnapshotItem, error) {
	f, err := os.Open(filename)
	util.Check(err)
	scanner := bufio.NewScanner(f)
	snapshots := []types.SnapshotItem{}
	for scanner.Scan() {
		pull_spec, err := componentizePullSpec(scanner.Text())
		util.Check(err)
		fmt.Printf("Read in registry: %s, namespace: %s, image: %s, tag: %s\n", pull_spec["registry"], pull_spec["namespace"], pull_spec["image"], pull_spec["tag"])
		item, err := types.NewSnapshotItem("https", pull_spec["registry"], pull_spec["namespace"], pull_spec["image"], pull_spec["tag"], "")
		util.Check(err)
		snapshots = append(snapshots, item)
	}
	util.Check(scanner.Err())
	return snapshots, nil
}

func componentizePullSpec(pull_spec string) (map[string]string, error) {
	split_pull_spec := strings.Split(pull_spec, "/")
	if len(split_pull_spec) != 3 {
		return nil, errors.New(fmt.Sprintf("Expected 3 components in pull spec, but found %d\n", len(split_pull_spec)))
	}
	imagetag := strings.Split(split_pull_spec[2], ":")
	if len(imagetag) != 2 {
		return nil, errors.New(fmt.Sprintf("Expected both an image and a tag in the pull spec, but found %s\n", split_pull_spec[2]))
	}

	elements := []string{"registry", "namespace", "image", "tag"}
	values := append(split_pull_spec[0:2], imagetag...)
	components := make(map[string]string)
	for i := 0; i < len(elements); i++ {
		components[elements[i]] = values[i]
	}
	return components, nil
}
