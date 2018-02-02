package main

import (
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/tealeg/xlsx"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// State represents a country state
type State struct {
	ID   uint
	Name string `gorm:"size:32;not null"`
	Code string `gorm:"size:2;not null"`
}

// Municipality represents a municipality in a state
type Municipality struct {
	ID      uint
	Name    string `gorm:"size:255;not null"`
	Code    string `gorm:"size:3;not null"`
	StateID uint   `gorm:"not null"`
}

// City represents a city in a state
type City struct {
	ID      uint
	Name    string `gorm:"size:255;not null"`
	Code    string `gorm:"size:2;not null"`
	StateID uint   `gorm:"not null"`
}

// NeighborhoodType represents a type of neighborhood
type NeighborhoodType struct {
	ID   uint
	Name string `gorm:"size:255;not null"`
	Code string `gorm:"size:2;not null"`
}

// Neighborhood represents a neighborhood in a municipality
type Neighborhood struct {
	ID             uint
	Name           string `gorm:"size:255;not null"`
	Code           string `gorm:"size:4;not null"`
	PostalCode     string `gorm:"size:5;not null"`
	Zone           string `gorm:"size:6;not null"`
	TypeID         uint   `gorm:"not null"`
	MunicipalityID uint   `gorm:"not null"`
	CityID         uint   `gorm:"not null"`
	StateID        uint   `gorm:"not null"`
}

func main() {
	f, err := xlsx.OpenFile("SEPOMEX.xlsx")

	if err != nil {
		panic(err)
	}

	db, err := gorm.Open("postgres", "user=postgres dbname=sepo password=postgres port=32768 sslmode=disable")

	if err != nil {
		panic(err)
	}

	db.SingularTable(true)

	db.AutoMigrate(&State{})
	db.AutoMigrate(&Municipality{})
	db.AutoMigrate(&City{})
	db.AutoMigrate(&NeighborhoodType{})
	db.AutoMigrate(&Neighborhood{})

	state := new(State)
	municipality := new(Municipality)
	city := new(City)
	t := new(NeighborhoodType)

	for _, sheet := range f.Sheets[1:] {
		name := strings.Replace(sheet.Name, "_", " ", -1)

		for _, row := range sheet.Rows[1:] {

			// Create State

			code := row.Cells[7].String() // c_estado

			if db.First(state, "code = ?", code).RecordNotFound() {
				state = &State{
					Name: name,
					Code: code,
				}

				db.Create(state)
			}

			// Create Municipality

			code = row.Cells[11].String() // c_mnpio

			if db.First(municipality, "code = ?", code).RecordNotFound() {
				municipality = &Municipality{
					Name:    row.Cells[3].String(), // D_mnpio
					Code:    code,
					StateID: state.ID,
				}

				db.Create(municipality)
			}

			// Create City

			code = row.Cells[14].String() // c_cve_ciudad

			if db.First(city, "code = ?", code).RecordNotFound() {
				city = &City{
					Name:    row.Cells[5].String(), // d_ciudad
					Code:    code,
					StateID: state.ID,
				}

				db.Create(city)
			}

			// Create NeighborhoodType

			code = row.Cells[10].String() // c_tipo_asenta

			if db.First(t, "code = ?", code).RecordNotFound() {
				t = &NeighborhoodType{
					Name: row.Cells[2].String(), // d_tipo_asenta
					Code: code,
				}

				db.Create(t)
			}

			// Create Neighborhood

			code = row.Cells[12].String() // id_asenta_cpcons

			db.Create(&Neighborhood{
				Name:           row.Cells[1].String(), // d_asenta
				Code:           code,
				PostalCode:     row.Cells[0].String(),  // d_codigo
				Zone:           row.Cells[13].String(), // d_zona
				TypeID:         t.ID,
				MunicipalityID: municipality.ID,
				CityID:         city.ID,
				StateID:        state.ID,
			})
		}

		break
	}
}
