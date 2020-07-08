package main

import (
	"bufio"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

type request struct {
	hostname      string
	cloudPlatform string
	instanceID    string
	project       string
	zone          string
	region        string
}

func readcsr(stdin *os.File) request {
	var line string
	var lines []string
	var csr string

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line = scanner.Text()
		lines = append(lines, line)
	}
	for _, v := range lines {
		csr = csr + fmt.Sprintf("%s\n", v)
	}

	return parseAttributes([]byte(csr))
}

func parseAttributes(puppetCSR []byte) request {
	var attributes request
	block, _ := pem.Decode(puppetCSR)
	csr, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}
	subjectCN := csr.Subject.String()
	attributes.hostname = trimLeftChars(subjectCN, 3)

	for _, extension := range csr.Extensions {
		oid := extension.Id.String()
		valueString := string(extension.Value)
		// Extension value needs to be trimmed because of ASN.1 encoding information
		// Read the Note from the link below
		// https://puppet.com/docs/puppet/latest/ssl_attributes_extensions.html#manually-checking-for-extensions-in-csrs-and-certificates
		value := trimLeftChars(valueString, 2)
		switch oid {
		case "1.3.6.1.4.1.34380.1.1.23":
			attributes.cloudPlatform = value
		case "1.3.6.1.4.1.34380.1.1.2":
			attributes.instanceID = value
		case "1.3.6.1.4.1.34380.1.1.7":
			attributes.project = value
		case "1.3.6.1.4.1.34380.1.1.20":
			attributes.zone = value
		case "1.3.6.1.4.1.34380.1.1.18":
			attributes.region = value
		}
	}

	return attributes
}

func trimLeftChars(s string, n int) string {
	m := 0
	for i := range s {
		if m >= n {
			return s[i:]
		}
		m++
	}
	return s[:0]
}
