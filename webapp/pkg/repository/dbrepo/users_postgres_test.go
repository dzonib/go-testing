package dbrepo

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"
	"webapp/pkg/data"
	"webapp/pkg/repository"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	host     = "localhost"
	user     = "postgres"
	password = "postgres"
	dbName   = "users_test"
	port     = "5435"
	dsn      = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5"
)

var resource *dockertest.Resource
var pool *dockertest.Pool
var testDB *sql.DB
var testRepo repository.DatabaseRepo

// TestMain gets executed before tests run
func TestMain(m *testing.M) {
	// connect to docker; fail if docker not running
	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker; is it running? %s", err)
	}

	pool = p

	// set up our docker options, specifying the image and so forth
	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14.5",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + dbName,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	// get a resource (docker image)
	resource, err = pool.RunWithOptions(&opts)
	if err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not start resource: %s", err)
	}

	// start the image and wait until it's ready
	if err := pool.Retry(func() error {
		var err error
		testDB, err = sql.Open("pgx", fmt.Sprintf(dsn, host, port, user, password, dbName))
		if err != nil {
			log.Println("Error:", err)
			return err
		}
		return testDB.Ping()
	}); err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not connect to database: %s", err)
	}

	// populate the database with empty tables
	err = createTables()
	if err != nil {
		log.Fatalf("error creating tables: %s", err)
	}

	// pool of database connections
	testRepo = &PostgresDBRepo{DB: testDB}

	// run tests
	code := m.Run()

	// clean up
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("could not purge resource: %s", err)
	}

	os.Exit(code)
}

func createTables() error {
	tableSQL, err := os.ReadFile("./testdata/users.sql")
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = testDB.Exec(string(tableSQL))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func Test_pingDB(t *testing.T) {
	err := testDB.Ping()
	if err != nil {
		t.Error("can't ping database")
	}
}

func TestPostgresDBRepoInsertUser(t *testing.T) {
	testUser := data.User{
		FirstName: "Admin",
		LastName:  "User",
		Email:     "admin@example.com",
		Password:  "secret",
		IsAdmin:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id, err := testRepo.InsertUser(testUser)

	if err != nil {
		t.Errorf("insert user returned an error %s", err)
	}

	if id != 1 {
		t.Errorf("insert user returns a wrong id; expected 1, but got %d", id)
	}
}

func TestPostgresDBRepoAllUsers(t *testing.T) {
	users, err := testRepo.AllUsers()

	if err != nil {
		t.Errorf("all users reports an error: %s", err)
	}

	if len(users) != 1 {
		t.Errorf("all users reports wrongs size; expected 1, but got %d", len(users))
	}

	jamesBond := data.User{
		FirstName: "James",
		LastName:  "Bond",
		Email:     "bond@example.com",
		Password:  "secret",
		IsAdmin:   0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, _ = testRepo.InsertUser(jamesBond)

	users, _ = testRepo.AllUsers()

	if len(users) != 2 {
		t.Errorf("all users reports wrongs size; expected 2, but got %d", len(users))
	}
}

func TestPostgresDBRepoGetUser(t *testing.T) {
	user, err := testRepo.GetUser(2)

	if err != nil {
		t.Error("expected to get james bond user, but got nothing")
	}

	if !strings.Contains("James", user.FirstName) {
		t.Errorf("Expected user name to be James, but got %s", user.FirstName)
	}

	_, err = testRepo.GetUser(34)

	if err == nil {
		t.Error("no error reported when getting a non existing user by id")
	}
}

func TestPostgresDBRepoGetUserByEmail(t *testing.T) {
	user, err := testRepo.GetUserByEmail("admin@example.com")

	if err != nil {
		t.Error("expected to get Admin user, but got nothing")
	}

	if user.ID != 1 {
		t.Errorf("Expected ID to be 1, but got %d", user.ID)
	}
}

func TestPostgresDBRepoUpdateUser(t *testing.T) {
	user, _ := testRepo.GetUser(2)

	const (
		newEmail = "newemail@example.com"
		newName  = "KINGKONG"
	)

	user.Email = newEmail
	user.FirstName = newName

	err := testRepo.UpdateUser(*user)

	if err != nil {
		t.Errorf("error updating user with id of %d: %s", 2, err)
	}

	user, _ = testRepo.GetUser(2)

	if user.Email != newEmail || user.FirstName != newName {
		t.Errorf("expected user email to be %s but got %s, and expected user name to be %s, but got %s", newEmail, user.Email, newName, user.FirstName)
	}

}
