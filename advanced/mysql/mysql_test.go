package mysql

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"testing"
)

//
//func initDB(dns string) (*sql.DB, error) {
//	db, err := sql.Open("mysql", dns)
//	if err != nil {
//		return nil, err
//	}
//	db.SetConnMaxLifetime(time.Minute * 3)
//	db.SetMaxOpenConns(10)
//	db.SetMaxIdleConns(10)
//	return db, nil
//}
//
//var (
//	Username     = "root"
//	Password     = "1qaz@WSX"
//	Host         = "127.0.0.1"
//	Port         = 3306
//	Database     = "foo"
//	Timeout      = "10s"
//	ReadTimeout  = "10s"
//	WriteTimeout = "10s"
//)
//
//var dns = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True&timeout=%s&readTimeout=%s&writeTimeout=%s", Username,
//	Password,
//	Host,
//	Port,
//	Database,
//	Timeout,
//	ReadTimeout,
//	WriteTimeout)

func TestSelect(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	exceptedPosts := []*post{
		{
			1, "title1", "content1",
		},
	}

	rows := sqlmock.NewRows([]string{"id", "title", "content"}).
		AddRow(1, "title1", "content1")
	//AddRow(2, "title2", "content2")

	mock.ExpectQuery("^SELECT (.)+ FROM tb_post$").WillReturnRows(rows)

	posts, err := queryPost(db)
	if err != nil {
		t.Error(err)
	}

	assert.EqualValues(t, exceptedPosts, posts)
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Error(err)
	}
}

type post struct {
	ID    int
	Title string
	Body  string
}

func queryPost(db *sql.DB) ([]*post, error) {
	rows, err := db.Query("SELECT * FROM tb_post")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*post
	for rows.Next() {
		p := &post{}
		if err := rows.Scan(&p.ID, &p.Title, &p.Body); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	if rows.Err() != nil {
		return nil, err
	}
	return posts, nil
}

func recordStats(db *sql.DB, userID, productID int64) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	if _, err = tx.Exec("UPDATE products SET views = views + 1"); err != nil {
		return
	}
	if _, err = tx.Exec("INSERT INTO product_viewers (user_id, product_id) VALUES (?, ?)", userID, productID); err != nil {
		return
	}
	return
}

// a successful case
func TestShouldUpdateStats(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE products").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO product_viewers").WithArgs(2, 3).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// now we execute our method
	if err = recordStats(db, 2, 3); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
