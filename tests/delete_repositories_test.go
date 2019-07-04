package tests

import (
	"testing"
	"time"
)

func TestDeleteRepositories(t *testing.T) {
	takeAction("/test_files/create_repository.yml", t)

	time.Sleep(time.Duration(int64(3)) * time.Second)

	takeAction("/test_files/delete_repositories.yml", t)
}
