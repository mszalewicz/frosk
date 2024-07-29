// Backend orchestrate all database actions
package backend

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	b64 "encoding/base64"

	"github.com/mszalewicz/frosk/helpers"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
)

type Backend struct {
	DB *sql.DB
}

// Returns GCM block cipher based on secret key
func InitGCM(secretKey []byte) (cipher.AEAD, error) {
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

func (backend *Backend) EncryptPasswordEntry(serviceName string, password string, masterPasswordGUI string) error {

	var (
		masterPasswordEncryptedBase64 string
		userSecretKeyEncryptedBase64  string
		initialVectorBase64           string
		saltBase64                    string
		masterPasswordEncrypted       []byte
		userSecretKeyEncrypted        []byte
		initialVector                 []byte
		salt                          []byte
	)

	// cleartext := []byte(password)

	row := backend.DB.QueryRow("SELECT password, secret_key, salt, initial_vector FROM master")
	err := row.Scan(&masterPasswordEncryptedBase64, &userSecretKeyEncryptedBase64, &saltBase64, &initialVectorBase64)

	if err != nil {
		errorWrapped := fmt.Errorf("Error during select query in master table: %w", err)
		slog.Error(errorWrapped.Error())
		return errorWrapped
	}

	masterPasswordEncrypted, errDecodingMasterPassword := b64.StdEncoding.DecodeString(masterPasswordEncryptedBase64)
	userSecretKeyEncrypted, errDecodingUserSecretKey := b64.StdEncoding.DecodeString(userSecretKeyEncryptedBase64)
	initialVector, errDecodingInitialVector := b64.StdEncoding.DecodeString(initialVectorBase64)
	salt, errDecodingSalt := b64.StdEncoding.DecodeString(saltBase64)

	if errDecodingMasterPassword != nil || errDecodingUserSecretKey != nil || errDecodingSalt != nil || errDecodingInitialVector != nil {
		errorWrapped := fmt.Errorf("Error during decoding base64 in | master password: %w | user secret key %w | salt: %w | initial vector: %w",
			errDecodingMasterPassword,
			errDecodingUserSecretKey,
			errDecodingSalt,
			errDecodingSalt)

		slog.Error(errorWrapped.Error())
		return errorWrapped
	}

	// Authenticate
	err = bcrypt.CompareHashAndPassword(masterPasswordEncrypted, []byte(masterPasswordGUI))

	if err != nil {
		errorWrapped := fmt.Errorf("Master password from GUI input do not match databse signature: %w", err)
		slog.Error(errorWrapped.Error())
		return errorWrapped
	}

	secretKey := pbkdf2.Key([]byte(masterPasswordGUI), salt, 4096, 32, sha256.New)

	gcmForUserSecretKey, err := InitGCM(secretKey)

	if err != nil {
		errorWrapped := fmt.Errorf("Error during initialization of GCM cipher block: %w", err)
		slog.Error(errorWrapped.Error())
		return errorWrapped
	}

	userSecretKey, err := gcmForUserSecretKey.Open(nil, initialVector, userSecretKeyEncrypted, nil)

	if err != nil {
		errorWrapped := fmt.Errorf("Error during decryption of user secret key (initialVector, userSecretKeyEncrypted) => (): %w", err)
		slog.Error(errorWrapped.Error())
		return errorWrapped
	}

	gcmPasswordEntry, err := InitGCM(userSecretKey)

	if err != nil {
		errorWrapped := fmt.Errorf("Error during initialization of GCM cipher block: %w", err)
		slog.Error(errorWrapped.Error())
		return errorWrapped
	}

	saltPasswordEntry := make([]byte, gcmPasswordEntry.NonceSize())
	_, err = rand.Read(saltPasswordEntry)

	if err != nil {
		errorWrapped := fmt.Errorf("Can't create random salt: %w", err)
		slog.Error(errorWrapped.Error())
		return errorWrapped
	}

	helpers.Assert(len(saltPasswordEntry), gcmPasswordEntry.NonceSize())

	passwordEncrypted := gcmPasswordEntry.Seal(nil, saltPasswordEntry, []byte(password), nil)

	saltPasswordEntryBase64 := b64.StdEncoding.EncodeToString(saltPasswordEntry)

	insertPasswordEntryQuery := `INSERT INTO passwords (service_name, password, initial_vector, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`

	return nil
}
