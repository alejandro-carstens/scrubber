package tests

import (
	"testing"
	"time"
)

func TestDeleteRepositories(t *testing.T) {
	takeAction("/testdata/create_repository.yml", t)

	time.Sleep(time.Duration(int64(3)) * time.Second)

	takeAction("/testdata/delete_repositories.yml", t)
}
