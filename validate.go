package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

const autosignerConfPath string = "/etc/puppetlabs/puppet/autosigner_hostnames.conf"

func (r request) validate() bool {

	// Check for local file first
	if validateFile(r) {
		// log.Printf("Sucessfully validated %v using autosigner_hostnames.conf\n", r.hostname)
		return true
	}

	// Confirm AWS and GCP projects are valid
	var valid bool = false
	for _, validProject := range config.allowedProjects {
		if validProject == r.project {
			// Check cloud platform
			switch r.cloudPlatform {
			case "aws":
				if validateAWS(r) {
					valid = true
				}
			case "gcp":
				if validateGCP(r) {
					valid = true
				}
			}
			if valid {
				// log.Printf("Sucessfully validated %v with Instance ID: %v, Zone: %v, Project: %v, Cloud Platform: %v\n", r.hostname, r.instanceID, r.zone, r.project, r.cloudPlatform)
				return true
			}
		}
	}
	if valid == false {
		log.Printf("Could not validate %v. The project \"%v\" does not match an allowed project.\n", r.hostname, r.project)
	}

	return false
}

// TODO: Add autosigner_hostname.conf rfc1035 validation
func validateFile(r request) bool {
	file, err := os.Open(autosignerConfPath)
	if err != nil {
		log.Fatalf("Failed opening file: %s, with error:%s", autosignerConfPath, err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string
	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}
	file.Close()

	for _, eachLine := range txtlines {
		// if eachline == r.hostname {
		if glob(eachLine, r.hostname) {
			// log.Printf("Hostname matched in %v, line \"%v\"\n", autosignerConfPath, eachLine)
			return true
		}
	}

	return false
}

func validateGCP(r request) bool {
	ctx := context.Background()
	instanceID, _ := strconv.ParseUint(r.instanceID, 10, 64)

	computeService, err := compute.NewService(ctx)
	if err != nil {
		log.Fatalf("Failed to create service: %v\n", err)
	}
	var completeList []*compute.Instance
	list, err := computeService.Instances.List(r.project, r.zone).Do()
	if err != nil {
		log.Fatalf("Failed to get InstanceList: %v\n", err)
	}
	completeList = append(completeList, list.Items...)

	if list.NextPageToken != "" {
		pageToken := list.NextPageToken
		for pageToken != "" {
			nextList, err := computeService.Instances.List(r.project, r.zone).Do(googleapi.QueryParameter("pageToken", pageToken))
			if err != nil {
				log.Fatalf("Failed to get InstanceList: %v\n", err)
			}
			completeList = append(completeList, nextList.Items...)
			pageToken = nextList.NextPageToken
		}
	}

	for _, v := range completeList {
		if v.Id == instanceID {
			return true
		}
	}
	return false
}

func validateAWS(r request) bool {
	fmt.Println("Validate AWS ran")
	return false
}

func glob(pattern, subj string) bool {
	const GLOB string = "*"

	// Empty pattern can only match empty subject
	if pattern == "" {
		return subj == pattern
	}

	// If the pattern _is_ a glob, it matches everything
	if pattern == GLOB {
		return true
	}

	parts := strings.Split(pattern, GLOB)

	if len(parts) == 1 {
		// No globs in pattern, so test for equality
		return subj == pattern
	}

	leadingGlob := strings.HasPrefix(pattern, GLOB)
	trailingGlob := strings.HasSuffix(pattern, GLOB)
	end := len(parts) - 1

	// Go over the leading parts and ensure they match.
	for i := 0; i < end; i++ {
		idx := strings.Index(subj, parts[i])

		switch i {
		case 0:
			// Check the first section. Requires special handling.
			if !leadingGlob && idx != 0 {
				return false
			}
		default:
			// Check that the middle parts match.
			if idx < 0 {
				return false
			}
		}

		// Trim evaluated text from subj as we loop over the pattern.
		subj = subj[idx+len(parts[i]):]
	}

	// Reached the last section. Requires special handling.
	return trailingGlob || strings.HasSuffix(subj, parts[end])
}
