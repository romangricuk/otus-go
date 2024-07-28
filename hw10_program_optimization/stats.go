package hw10programoptimization

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/ugorji/go/codec"
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
	handle := new(codec.JsonHandle)
	decoder := codec.NewDecoder(bufReader, handle)
	var user User
	var email string

	for {
		if err := decoder.Decode(&user); errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, fmt.Errorf("json decode error: %w", err)
		}

		email = strings.ToLower(user.Email)
		if strings.HasSuffix(email, domain) {
			atIndex := strings.LastIndex(email, "@")
			if atIndex != -1 {
				domainPart := email[atIndex+1:]
				domainStat[domainPart]++
			}
		}
		user = User{}
	}

	return domainStat, nil
}
