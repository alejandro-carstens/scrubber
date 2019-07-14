package tests

import (
	"testing"
	"time"
)

func TestDeleteRepositories(t *testing.T) {
	takeAction("/testfiles/create_repository.yml", t)

	time.Sleep(time.Duration(int64(3)) * time.Second)

	takeAction("/testfiles/delete_repositories.yml", t)
}
