// Backend orchestrates all database actions
package backend

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	b64 "encoding/base64"

	"github.com/mszalewicz/frosk/helpers"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
)

var EmptyPassword = errors.New("No password given to insert.")
var EmptyUsername = errors.New("No username name given to insert.")
var EmptyServiceName = errors.New("No service name given to insert.")
var EmptyMasterPassord = errors.New("No master passwrod given to compare.")
var ServiceNameNotFound = errors.New("Provided service name is not present in database.")
var ServiceNameAlreadyTaken = errors.New("Provided service name is already present in database.")

type Backend struct {
	DB *sql.DB
}

type PasswordEntry struct {
	Username    string
	Password    string
	ServiceName string
}

// Returns GCM block cipher based on secret key
func InitGCM(secretKey []byte) (cipher.AEAD, error) {
	helpers.Assert(len(secretKey), 32) // secret key has to 32 byte long

	aes, err := aes.NewCipher(secretKey)

	if err != nil {
		errorWrapped := fmt.Errorf("Error during creation of AES cipher block: %w", err)
		slog.Error(errorWrapped.Error())
		return nil, errorWrapped
	}

	gcm, err := cipher.NewGCM(aes)

	if err != nil {
		errorWrapped := fmt.Errorf("Error during creation of AES GCM instance: %w", err)
		slog.Error(errorWrapped.Error())
		return nil, errorWrapped
	}

	return gcm, nil
}

func (backend *Backend) GetUserSecretKey(masterPasswordGUI string) ([]byte, error) {
	var (
		masterPasswordEncryptedBase64    string
		masterPasswordEncrypted          []byte
		userSecretKeyEncryptedBase64     string
		userSecretKeyEncrypted           []byte
		userSecretKey                    []byte
		initialVectorUserSecretKeyBase64 string
		initialVectorUserSecretKey       []byte
		saltBase64                       string
		salt                             []byte
	)

	if len(masterPasswordGUI) == 0 {
		return userSecretKey, EmptyMasterPassord
	}

	row := backend.DB.QueryRow("SELECT \"password\", secret_key, salt, initial_vector FROM master")
	err := row.Scan(&masterPasswordEncryptedBase64, &userSecretKeyEncryptedBase64, &saltBase64, &initialVectorUserSecretKeyBase64)

	if err != nil {
		errorWrapped := fmt.Errorf("Error during select query to master table: %w", err)
		slog.Error(errorWrapped.Error())
		return userSecretKey, err
	}

	masterPasswordEncrypted, errDecodingMasterPassword := b64.StdEncoding.DecodeString(masterPasswordEncryptedBase64)
	userSecretKeyEncrypted, errDecodingUserSecretKey := b64.StdEncoding.DecodeString(userSecretKeyEncryptedBase64)
	initialVectorUserSecretKey, errDecodingInitialVector := b64.StdEncoding.DecodeString(initialVectorUserSecretKeyBase64)
	salt, errDecodingSalt := b64.StdEncoding.DecodeString(saltBase64)

	if errDecodingMasterPassword != nil || errDecodingUserSecretKey != nil || errDecodingSalt != nil || errDecodingInitialVector != nil {
		errorWrapped := fmt.Errorf("Error during decoding base64 in | master password: %w | user secret key %w | salt: %w | initial vector: %w",
			errDecodingMasterPassword,
			errDecodingUserSecretKey,
			errDecodingSalt,
			errDecodingSalt)

		slog.Error(errorWrapped.Error())
		return userSecretKey, errorWrapped
	}

	// Authenticate
	err = bcrypt.CompareHashAndPassword(masterPasswordEncrypted, []byte(masterPasswordGUI))

	if err != nil {
		errorWrapped := fmt.Errorf("Master password from GUI input do not match databse signature: %w", err)
		slog.Error(errorWrapped.Error())
		return userSecretKey, errorWrapped
	}

	secretKey := pbkdf2.Key([]byte(masterPasswordGUI), salt, 4096, 32, sha256.New)
	gcmForUserSecretKey, err := InitGCM(secretKey)

	if err != nil {
		errorWrapped := fmt.Errorf("Error during gcm initialization: %w", err)
		slog.Error(errorWrapped.Error())
		return userSecretKey, errorWrapped
	}

	userSecretKey, err = gcmForUserSecretKey.Open(nil, initialVectorUserSecretKey, userSecretKeyEncrypted, nil)

	if err != nil {
		errorWrapped := fmt.Errorf("Error during user secret key decryption: %w", err)
		slog.Error(errorWrapped.Error())
		return userSecretKey, errorWrapped
	}

	return userSecretKey, nil
}

// Opens connection to the sqlite local database
func Initialize(applicationDB string) (*Backend, error) {

	db, err := sql.Open("sqlite3", applicationDB)
	if err != nil {
		errWrapped := fmt.Errorf("Could not initializa db file: %w", err)
		slog.Error(errWrapped.Error())
		return nil, errWrapped
	}

	backend := Backend{DB: db}
	return &backend, nil
}

// Create db from schema
func (backend *Backend) CreateStructure() error {

	const create_passwords_table = `
		CREATE TABLE IF NOT EXISTS passwords (
	       id INTEGER PRIMARY KEY AUTOINCREMENT,
	       service_name TEXT UNIQUE NOT NULL,
		   username TEXT NOT NULL,
	       password TEXT NOT NULL,
		   initial_vector TEXT UNIQUE NOT NULL,
	       created_at TEXT NULL,
	       updated_at TEXT NULL
	   ) STRICT;
	`

	_, err := backend.DB.Exec(create_passwords_table)

	if err != nil {
		errWrapped := fmt.Errorf("Error during creating passwords table: %w", err)
		slog.Error(errWrapped.Error())
		return err
	}

	const create_master_table = `
		CREATE TABLE IF NOT EXISTS master (
		   	id INTEGER PRIMARY KEY AUTOINCREMENT,
		    password TEXT UNIQUE NOT NULL,
		    secret_key TEXT UNIQUE NOT NULL,
		    salt TEXT UNIQUE NOT NULL,
			initial_vector TEXT UNIQUE NOT NULL,
		    created_at TEXT NULL,
		    updated_at TEXT NULL
		) STRICT;
	`

	_, err = backend.DB.Exec(create_master_table)

	if err != nil {
		errWrapped := fmt.Errorf("Error during creating passwords table: %w", err)
		slog.Error(errWrapped.Error())
		return errWrapped
	}

	return nil
}

func (backend *Backend) CountMasterEntries() (int, error) {
	var numberOfEntriesInMaster int
	row := backend.DB.QueryRow("SELECT COUNT(*) FROM master")

	err := row.Scan(&numberOfEntriesInMaster)

	if err != nil {
		errWrapped := fmt.Errorf("Query counting number of entries in master table: %w", err)
		slog.Error(errWrapped.Error())
		return 0, errWrapped
	}

	return numberOfEntriesInMaster, nil
}

func (backend *Backend) CountServiceNameOccurences(serviceName string) (int, error) {
	var numberOfServiceNameOccurences int
	row := backend.DB.QueryRow("SELECT COUNT(service_name) FROM passwords WHERE service_name = ?", serviceName)

	err := row.Scan(&numberOfServiceNameOccurences)

	if err != nil {
		errWrapped := fmt.Errorf("Query counting number of entries in master table: %w", err)
		slog.Error(errWrapped.Error())
		return 0, errWrapped
	}

	return numberOfServiceNameOccurences, nil
}

// Create all necessary crypto primitives and insert them with master password to db
func (backend *Backend) InitMaster(masterPassword string) error {
	// Flow:
	//     master password   -> encrypted as bcrypt || stored encrypted
	//     salt              -> used to derive secret key from master password, used in encrypting user secret key || created randomly || length = 16 || stored
	//     initial vector    -> used for storing user secret key || created randomly || length = gcm nonce size || stored
	//     user secret key   -> used to encrypt all user passwords || created randomly || length = 32 (maximal length, corresponding to AES-256) || stored encrypted
	// Info:
	//     master password secret key  -> derived from master password with PKBDF2, using salt
	//     user secret key             -> used in encryption of user stored passwords

	if len(masterPassword) == 0 {
		return EmptyMasterPassord
	}

	helpers.AssertBigger(len(masterPassword), 0)

	plaintext := []byte(masterPassword)

	// Increasing cost by 1, increases hashing execution time by 2x
	// i.e cost 14 will take 2x time than if cost was equal to 13
	bcryptHash, err := bcrypt.GenerateFromPassword(plaintext, 14)

	if err != nil {
		errWrapped := fmt.Errorf("Could not calculate bcrypt hash: %w", err)
		slog.Error(errWrapped.Error())
		return errWrapped
	}

	masterPasswordHashBase64 := b64.StdEncoding.EncodeToString(bcryptHash)

	// Length 16 is a compromise between security and execution time for PBKDF2
	salt := make([]byte, 16)
	_, err = rand.Read(salt)

	if err != nil {
		errorWrapped := fmt.Errorf("Error during randomizing salt: %w", err)
		slog.Error(errorWrapped.Error())
		return errorWrapped
	}

	saltBase64 := b64.StdEncoding.EncodeToString(salt)

	// Length 32 is maximal for AES secret key and corresponds to usage of AES-256
	secretKey := pbkdf2.Key(plaintext, salt, 4096, 32, sha256.New)

	helpers.Assert(len(secretKey), 32)

	gcm, err := InitGCM(secretKey)

	if err != nil {
		errorWrapped := fmt.Errorf("Error during initialization of GCM cipher block: %w", err)
		slog.Error(errorWrapped.Error())
		return errorWrapped
	}

	initialVector := make([]byte, gcm.NonceSize())
	_, err = rand.Read(initialVector)

	if err != nil {
		errorWrapped := fmt.Errorf("Error during randomizing initial vector: %w", err)
		slog.Error(errorWrapped.Error())
		return errorWrapped
	}

	initialVectorBase64 := b64.StdEncoding.EncodeToString(initialVector)

	userSecretKey := make([]byte, 32)
	_, err = rand.Read(userSecretKey)

	if err != nil {
		errorWrapped := fmt.Errorf("Error during randomizing userSecretKey: %w", err)
		slog.Error(errorWrapped.Error())
		return errorWrapped
	}

	helpers.Assert(len(userSecretKey), 32)

	userSecretKeyEncrypted := gcm.Seal(nil, initialVector, userSecretKey, nil)
	userSecretKeyEncryptedBase64 := b64.StdEncoding.EncodeToString(userSecretKeyEncrypted)

	now := helpers.TimeTo8601String(time.Now())

	queryResult, err := backend.DB.Exec(
		"INSERT INTO master (password, secret_key, salt, initial_vector, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		masterPasswordHashBase64, userSecretKeyEncryptedBase64, saltBase64, initialVectorBase64, now, now)

	if err != nil {
		err := fmt.Errorf("Error during insert into master execution: %w", err)
		slog.Error(err.Error())
		return err
	}

	rowsAffected, err := queryResult.RowsAffected()

	if rowsAffected != 1 {
		err := fmt.Errorf("Expected to insert exactly 1 row into master table. Inserted 0 / multiple rows")
		slog.Error(err.Error())
		return err
	}

	return nil
}

// Inserts encrypted password and username for given service name
func (backend *Backend) EncryptPasswordEntry(serviceName string, password string, username string, masterPasswordGUI string) error {

	if len(serviceName) == 0 {
		return EmptyServiceName
	}

	if len(password) == 0 {
		return EmptyPassword
	}

	if len(masterPasswordGUI) == 0 {
		return EmptyMasterPassord
	}

	if len(username) == 0 {
		return EmptyUsername
	}

	serviceNameOccurences, err := backend.CountServiceNameOccurences(serviceName)

	if err != nil {
		errorWrapped := fmt.Errorf("Problem quering count of service name occurences in passwords table: %v", err)
		slog.Error(errorWrapped.Error())
		return errorWrapped
	}

	if serviceNameOccurences != 0 {
		return ServiceNameAlreadyTaken
	}

	userSecretKey, err := backend.GetUserSecretKey(masterPasswordGUI)

	if err != nil {
		errorWrapped := fmt.Errorf("Error during decryption of user secret key: %w", err)
		slog.Error(errorWrapped.Error())
		return errorWrapped
	}

	gcmPasswordEntry, err := InitGCM(userSecretKey)

	if err != nil {
		errorWrapped := fmt.Errorf("Error during initialization of GCM cipher block: %w", err)
		slog.Error(errorWrapped.Error())
		return errorWrapped
	}

	initialVectorPasswordEntry := make([]byte, gcmPasswordEntry.NonceSize())
	_, err = rand.Read(initialVectorPasswordEntry)

	if err != nil {
		errorWrapped := fmt.Errorf("Can't create random salt: %w", err)
		slog.Error(errorWrapped.Error())
		return errorWrapped
	}

	helpers.Assert(len(initialVectorPasswordEntry), gcmPasswordEntry.NonceSize())

	initialVectorPasswordEntryBase64 := b64.StdEncoding.EncodeToString(initialVectorPasswordEntry)

	passwordEncrypted := gcmPasswordEntry.Seal(nil, initialVectorPasswordEntry, []byte(password), nil)
	passwordEncryptedBase64 := b64.StdEncoding.EncodeToString(passwordEncrypted)

	usernameEncrypted := gcmPasswordEntry.Seal(nil, initialVectorPasswordEntry, []byte(username), nil)
	usernameEncryptedBase64 := b64.StdEncoding.EncodeToString(usernameEncrypted)

	now := helpers.TimeTo8601String(time.Now())

	insertPasswordEntryQuery := `INSERT INTO passwords (service_name, username, password, initial_vector, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`

	_, err = backend.DB.Exec(insertPasswordEntryQuery, serviceName, usernameEncryptedBase64, passwordEncryptedBase64, initialVectorPasswordEntryBase64, now, now)

	if err != nil {
		errWrapped := fmt.Errorf("Error inserting password entry into passwords: %w", err)
		slog.Error(errWrapped.Error())
		return errWrapped
	}

	return nil
}

// Finds and decrypts passwrod and username for given service name
func (backend *Backend) DecryptPasswordEntry(serviceName string, masterPasswordGUI string) (PasswordEntry, error) {
	var (
		passwordEntry           PasswordEntry
		initialVectorBase64     string
		initialVector           []byte
		passwordEncryptedBase64 string
		passwordEncrypted       []byte
		usernameEncryptedBase64 string
		usernameEncrypted       []byte
		password                []byte
		username                []byte
	)

	row := backend.DB.QueryRow("SELECT service_name, username, \"password\", initial_vector FROM passwords WHERE	service_name = ?", serviceName)
	err := row.Scan(&passwordEntry.ServiceName, &usernameEncryptedBase64, &passwordEncryptedBase64, &initialVectorBase64)

	if err != nil {
		errorWrapped := fmt.Errorf("Error during count query on passwords table - looking for service name = %s: %w", serviceName, err)
		slog.Error(errorWrapped.Error())
		return passwordEntry, errorWrapped
	}

	initialVector, errDecodeInitialVectorBase64 := b64.StdEncoding.DecodeString(initialVectorBase64)
	passwordEncrypted, errDecodePasswordEncryptedBas64 := b64.StdEncoding.DecodeString(passwordEncryptedBase64)
	usernameEncrypted, errDecodeUsernameEncryptedBase64 := b64.StdEncoding.DecodeString(usernameEncryptedBase64)

	if errDecodeInitialVectorBase64 != nil || errDecodePasswordEncryptedBas64 != nil || errDecodeUsernameEncryptedBase64 != nil {
		errorWrapped := fmt.Errorf("Error during coversion from base 64 - initial vector: %w | password: %w | username: %w", errDecodeInitialVectorBase64, errDecodePasswordEncryptedBas64, errDecodeUsernameEncryptedBase64)
		slog.Error(errorWrapped.Error())
		return passwordEntry, errorWrapped
	}

	userSecretKey, err := backend.GetUserSecretKey(masterPasswordGUI)

	if err != nil {
		errorWrapped := fmt.Errorf("Error during decrytion of user secret key: %w", err)
		slog.Error(errorWrapped.Error())
		return passwordEntry, errorWrapped
	}

	gcm, err := InitGCM(userSecretKey)

	if err != nil {
		errorWrapped := fmt.Errorf("Error during creation of gcm cipher block: %w", err)
		slog.Error(errorWrapped.Error())
		return passwordEntry, errorWrapped
	}

	password, err = gcm.Open(nil, initialVector, passwordEncrypted, nil)

	if err != nil {
		errorWrapped := fmt.Errorf("Error during password decryption: %w", err)
		slog.Error(errorWrapped.Error())
		return passwordEntry, errorWrapped
	}

	username, err = gcm.Open(nil, initialVector, usernameEncrypted, nil)

	if err != nil {
		errorWrapped := fmt.Errorf("Error during username decryption: %w", err)
		slog.Error(errorWrapped.Error())
		return passwordEntry, errorWrapped
	}

	passwordEntry.Password = string(password)
	passwordEntry.Username = string(username)

	return passwordEntry, nil
}

func (backend *Backend) GetPasswordEntriesList() ([]string, error) {
	services := make([]string, 0)
	service := ""

	query := "SELECT service_name FROM passwords"
	rows, err := backend.DB.Query(query)
	if err != nil {
		errWrapped := fmt.Errorf("Error during getting service names for passwords: %w", err)
		slog.Error(errWrapped.Error())
		return nil, errWrapped
	}

	for rows.Next() {
		rows.Scan(&service)
		if len(service) > 0 {
			services = append(services, service)
		}
	}

	return services, nil
}
