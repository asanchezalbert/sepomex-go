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
	Zone           string `gorm:"size:10;not null"`
	TypeID         uint   `gorm:"not null"`
	MunicipalityID uint   `gorm:"not null"`
	CityID         uint
	StateID        uint `gorm:"not null"`
}

func Clean(c *xlsx.Cell) string {
	return strings.Trim(c.String(), " ")
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

			code := Clean(row.Cells[7]) // c_estado

			if db.First(state, "code = ?", code).RecordNotFound() {
				state = &State{
					Name: name,
					Code: code,
				}

				db.Create(state)
			}

			// Create Municipality

			code = Clean(row.Cells[11]) // c_mnpio

			if db.First(municipality, "code = ?", code).RecordNotFound() {
				municipality = &Municipality{
					Name:    Clean(row.Cells[3]), // D_mnpio
					Code:    code,
					StateID: state.ID,
				}

				db.Create(municipality)
			}

			// Create City

			if Clean(row.Cells[5]) != "" {
				code = Clean(row.Cells[14]) // c_cve_ciudad

				if db.First(city, "code = ?", code).RecordNotFound() {
					city = &City{
						Name:    Clean(row.Cells[5]), // d_ciudad
						Code:    code,
						StateID: state.ID,
					}

					db.Create(city)
				}
			}

			// Create NeighborhoodType

			code = Clean(row.Cells[10]) // c_tipo_asenta

			if db.First(t, "code = ?", code).RecordNotFound() {
				t = &NeighborhoodType{
					Name: Clean(row.Cells[2]), // d_tipo_asenta
					Code: code,
				}

				db.Create(t)
			}

			// Create Neighborhood

			code = Clean(row.Cells[12]) // id_asenta_cpcons

			db.Create(&Neighborhood{
				Name:           Clean(row.Cells[1]), // d_asenta
				Code:           code,
				PostalCode:     Clean(row.Cells[0]),  // d_codigo
				Zone:           Clean(row.Cells[13]), // d_zona
				TypeID:         t.ID,
				MunicipalityID: municipality.ID,
				CityID:         city.ID,
				StateID:        state.ID,
			})
		}
	}
}
