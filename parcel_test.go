package main

import (
	"database/sql"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())
	randRange  = rand.New(randSource)
)

func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func TestAddGetDelete(t *testing.T) {
	os.Remove("test_tracker.db")
	db, err := sql.Open("sqlite", "test_tracker.db")
	require.NoError(t, err)
	defer func() {
		db.Close()
		os.Remove("test_tracker.db")
	}()

	createTableSQL := `CREATE TABLE IF NOT EXISTS parcel (
		number INTEGER PRIMARY KEY AUTOINCREMENT,
		client INTEGER NOT NULL,
		status TEXT NOT NULL,
		address TEXT NOT NULL,
		created_at TEXT NOT NULL
	);`
	_, err = db.Exec(createTableSQL)
	require.NoError(t, err)

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.Greater(t, id, 0)

	result, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, parcel.Client, result.Client)
	require.Equal(t, parcel.Status, result.Status)
	require.Equal(t, parcel.Address, result.Address)
	require.Equal(t, parcel.CreatedAt, result.CreatedAt)

	err = store.Delete(id)
	require.NoError(t, err)

	_, err = store.Get(id)
	require.Error(t, err)
}

func TestSetAddress(t *testing.T) {
	os.Remove("test_tracker.db")
	db, err := sql.Open("sqlite", "test_tracker.db")
	require.NoError(t, err)
	defer func() {
		db.Close()
		os.Remove("test_tracker.db")
	}()

	createTableSQL := `CREATE TABLE IF NOT EXISTS parcel (
		number INTEGER PRIMARY KEY AUTOINCREMENT,
		client INTEGER NOT NULL,
		status TEXT NOT NULL,
		address TEXT NOT NULL,
		created_at TEXT NOT NULL
	);`
	_, err = db.Exec(createTableSQL)
	require.NoError(t, err)

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.Greater(t, id, 0)

	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	result, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, newAddress, result.Address)
}

func TestSetStatus(t *testing.T) {
	os.Remove("test_tracker.db")
	db, err := sql.Open("sqlite", "test_tracker.db")
	require.NoError(t, err)
	defer func() {
		db.Close()
		os.Remove("test_tracker.db")
	}()

	createTableSQL := `CREATE TABLE IF NOT EXISTS parcel (
		number INTEGER PRIMARY KEY AUTOINCREMENT,
		client INTEGER NOT NULL,
		status TEXT NOT NULL,
		address TEXT NOT NULL,
		created_at TEXT NOT NULL
	);`
	_, err = db.Exec(createTableSQL)
	require.NoError(t, err)

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.Greater(t, id, 0)

	err = store.SetStatus(id, ParcelStatusSent)
	require.NoError(t, err)

	result, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, ParcelStatusSent, result.Status)
}

func TestGetByClient(t *testing.T) {
	os.Remove("test_tracker.db")
	db, err := sql.Open("sqlite", "test_tracker.db")
	require.NoError(t, err)
	defer func() {
		db.Close()
		os.Remove("test_tracker.db")
	}()

	createTableSQL := `CREATE TABLE IF NOT EXISTS parcel (
		number INTEGER PRIMARY KEY AUTOINCREMENT,
		client INTEGER NOT NULL,
		status TEXT NOT NULL,
		address TEXT NOT NULL,
		created_at TEXT NOT NULL
	);`
	_, err = db.Exec(createTableSQL)
	require.NoError(t, err)

	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		require.Greater(t, id, 0)

		parcels[i].Number = id
		parcelMap[id] = parcels[i]
	}

	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	require.Len(t, storedParcels, 3)

	for _, parcel := range storedParcels {
		expected, exists := parcelMap[parcel.Number]
		require.True(t, exists)
		require.Equal(t, expected.Client, parcel.Client)
		require.Equal(t, expected.Status, parcel.Status)
		require.Equal(t, expected.Address, parcel.Address)
		require.Equal(t, expected.CreatedAt, parcel.CreatedAt)
	}
}