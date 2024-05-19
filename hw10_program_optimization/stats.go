package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	domainStat := make(DomainStat)
	domain = "." + strings.ToLower(domain)

	bufReader := bufio.NewReader(r)
	decoder := json.NewDecoder(bufReader)

	for {
		var user User
		if err := decoder.Decode(&user); errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, fmt.Errorf("json decode error: %w", err)
		}

		email := strings.ToLower(user.Email)
		if strings.HasSuffix(email, domain) {
			atIndex := strings.LastIndex(email, "@")
			if atIndex != -1 {
				domainPart := email[atIndex+1:]
				domainStat[domainPart]++
			}
		}
	}

	return domainStat, nil
}
