package main

import (
	"fmt"
	"strings"

	"github.com/tealeg/xlsx"
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

	for _, sheet := range f.Sheets[1:] {
		state := strings.Replace(sheet.Name, "_", " ", -1)

		for _, row := range sheet.Rows[1:] {
			state := &State{
				Name: state,
				Code: row.Cells[7].String(), // c_estado
			}

			municipality := &Municipality{
				Name:    row.Cells[3].String(),  // D_mnpio
				Code:    row.Cells[11].String(), // c_mnpio
				StateID: state.ID,
			}

			city := &City{
				Name:    row.Cells[5].String(),  // d_ciudad
				Code:    row.Cells[14].String(), // c_cve_ciudad
				StateID: state.ID,
			}

			t := &NeighborhoodType{
				Name: row.Cells[2].String(),  // d_tipo_asenta
				Code: row.Cells[10].String(), // c_tipo_asenta
			}

			fmt.Println(&Neighborhood{
				Name:           row.Cells[1].String(),  // d_asenta
				Code:           row.Cells[12].String(), // id_asenta_cpcons
				PostalCode:     row.Cells[0].String(),  // d_codigo
				Zone:           row.Cells[13].String(), // d_zona
				TypeID:         t.ID,
				MunicipalityID: municipality.ID,
				CityID:         city.ID,
				StateID:        state.ID,
			})

			break
		}

		break
	}
}
