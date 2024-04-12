package postgre

import(
	"database/sql"

	_ "github.com/lib/pq"
)

func main(){
	db, err := sql.Open("postgres", "user=postgres password=postgres dbname=postgres sslmode=disable")
	checkErr(err)
	


	defer db.close()
}