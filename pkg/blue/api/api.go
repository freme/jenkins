package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

type (
	//https://mholt.github.io/json-to-go/
	RunsApi []struct {
		Links struct {
			Nodes struct {
				Class string `json:"_class"`
				Href  string `json:"href"`
			} `json:"nodes"`
		} `json:"_links"`
		ID     string `json:"id"`
		Name   string `json:"name"`
		Result string `json:"result"`
		State  string `json:"state"`
	}

	NodesApi []struct {
		Links struct {
			Steps struct {
				Class string `json:"_class"`
				Href  string `json:"href"`
			} `json:"steps"`
		} `json:"_links"`
		DisplayName string `json:"displayName"`
		ID          string `json:"id"`
		Result      string `json:"result"`
		State       string `json:"state"`
	}

	Node struct {
		Href        string
		DisplayName string
		Result      string
	}
	StepsApi []struct {
		Actions []struct {
			Links struct {
				Self struct {
					Class string `json:"_class"`
					Href  string `json:"href"`
				} `json:"self"`
			} `json:"_links"`
			URLName string `json:"urlName"`
		} `json:"actions"`
		DisplayName string `json:"displayName"`
		ID          string `json:"id"`
		Result      string `json:"result"`
		State       string `json:"state"`
	}
)

func GetJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func GetJsonFromFile(path string, target interface{}) error {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	json.Unmarshal(raw, &target)
	return err
}

func GetLogs(jenkins string, build string, pipeline string) error {
	runs := RunsApi{}
	url := jenkins + "/blue/rest/organizations/jenkins" + pipeline + "/runs/"
	state := ""
	result := ""
	GetJson(url, &runs)
	if build == "0" {
		// find first finished build
		for _, run := range runs {
			if run.State == "FINISHED" {
				build = run.ID
				url = jenkins + run.Links.Nodes.Href
				state = run.State
				result = run.Result
				break
			}
		}
	} else {
		buildExists := false
		for _, run := range runs {
			if run.ID == build {
				buildExists = true
				url = jenkins + run.Links.Nodes.Href
				state = run.State
				result = run.Result
				break
			}
		}
		if !buildExists {
			fmt.Printf("build %s not available via API\n", build)
		}
	}
	fmt.Printf("BuildURL: %s\n", url)
	fmt.Printf("BuildState: %s, BuildResult: %s\n", state, result)

	nodes := NodesApi{}
	failedNodes := []Node{}

	GetJson(url, &nodes)
	fmt.Printf("Found %d Nodes\n", len(nodes))
	failed := 0
	success := 0
	notBuilt := 0
	unknown := 0
	for _, node := range nodes {
		if node.Result == "FAILURE" {
			failedNodes = append(failedNodes, Node{Href: node.Links.Steps.Href, DisplayName: node.DisplayName, Result: node.Result})
			failed++
		} else if node.Result == "SUCCESS" {
			success++
		} else if node.Result == "NOT_BUILT" {
			notBuilt++
		} else {
			unknown++
			fmt.Printf("%s: %s\n", node.DisplayName, node.Result)
		}
	}
	fmt.Printf("SUCCESS: %d, FAILED: %d, NOT_BUILT: %d, UNKNOWN: %d\n", success, failed, notBuilt, unknown)
	for _, node := range failedNodes {
		fmt.Printf("Node: %s FAILED\n", node.DisplayName)
	}

	steps := StepsApi{}
	logUrls := []string{}

	for _, node := range failedNodes {
		url = jenkins + node.Href
		GetJson(url, &steps)
		for _, step := range steps {
			if step.Result == "FAILURE" {
				for _, action := range step.Actions {
					if action.URLName == "log" {
						logUrls = append(logUrls, action.Links.Self.Href)
					}
				}
			}
		}
	}
	fmt.Printf("\nLogs of failed steps:\n")
	for _, log := range logUrls {
		fmt.Println(jenkins + log)
	}

	return nil
}

// vim: ts=2 sw=2 et
