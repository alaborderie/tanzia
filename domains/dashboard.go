package domains

import (
	"html/template"
	"net/http"
	"tanzia/helpers"
)

type DashboardData struct {
	Persons        []Person
	Bills          []Bill
	Provisions     []Provision
	TotalTanzias int
	Balance        float64
}

func getDashboardData(userId string) DashboardData {
	db, err := helpers.GetConnectionManager().GetConnection("postgres")
	if err != nil {
		panic(err)
	}

	personRows, _ := db.Query("SELECT name, tanzia FROM persons WHERE userId = $1", userId)
	billRows, _ := db.Query("SELECT label, amount FROM bills WHERE userId = $1", userId)
	provisionRows, _ := db.Query("SELECT label, amount FROM provisions WHERE userId = $1", userId)

	var Persons []Person
	var Bills []Bill
	var Provisions []Provision
	var Balance float64
	TotalTanzias := 0

	for personRows.Next() {
		var person Person
		err := personRows.Scan(&person.Name, &person.Tanzia)
		if err != nil {
			panic(err)
		}
		Persons = append(Persons, person)
		TotalTanzias += person.Tanzia
	}

	for billRows.Next() {
		var bill Bill
		err := billRows.Scan(&bill.Label, &bill.Amount)
		if err != nil {
			panic(err)
		}
		Balance -= bill.Amount
		Bills = append(Bills, bill)
	}

	for provisionRows.Next() {
		var provision Provision
		err := provisionRows.Scan(&provision.Label, &provision.Amount)
		if err != nil {
			panic(err)
		}
		Balance += provision.Amount
		Provisions = append(Provisions, provision)
	}

	return DashboardData{
		Persons,
		Bills,
		Provisions,
		TotalTanzias,
		Balance,
	}
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	userId, ok := GetAuthenticatedUserId(w, r)

	if !ok {
		http.Redirect(w, r, "/logout", http.StatusFound)
		return
	}

	data := getDashboardData(userId)

	t, _ := template.ParseFiles("templates/dashboard.html")
	err := t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
