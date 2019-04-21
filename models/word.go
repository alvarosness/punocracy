package models

import (
	"errors"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Word is the core of our project

// Specifies the structure of words stored in the Words entity
type WordRow struct {
	WordID         int    `db:"wordID"`
	Word           string `db:"word"`
	HomophoneGroup int    `db:"homophoneGroup"`
}

type Word struct {
	Base
}

func NewWord(db *sqlx.DB) *Word {
	word := &Word{}
	word.db = db
	word.table = "Words_T"
	word.hasID = true

	return word
}

//Query the firstLetter for a lit of strucs
//input: a single "rune" representing the first letter
//this is case insensitive
//output: a list of type WordRow
func (w *Word) QueryAlph(tx *sqlx.Tx, firstLetter rune) ([]WordRow, error) {
	words := []WordRow{} //this creates a nil slice of type wordRow, that can be appended to

	queryString := string(firstLetter) + "%"

	erri := w.db.Select(&words, "SELECT * FROM Words_T WHERE word LIKE ? ORDER BY word;", queryString)

	if erri != nil {
		return words, erri
	}
	if len(words) == 0 {
		err1 := errors.New("empty list")
		return words, err1
	}

	return words, nil
}

/*
   Pass in a word string to generate a list of homophones in alphabetical order
   Not including the word tested.
*/
func (w *Word) QueryHlistString(tx *sqlx.Tx, inputWord string) ([]WordRow, error) {

	words := []WordRow{} //this creates a nil slice of type WordRow, that can be appended to

	erri := w.db.Select(&words, "SELECT * FROM Words_T WHERE homophoneGroup = (SELECT homophoneGroup FROM Words_T WHERE word LIKE ? ) AND word NOT LIKE ? ORDER BY word;", inputWord, inputWord)

	if erri != nil {
		return words, erri
	}

	if len(words) == 0 {
		err1 := errors.New("empty list")
		return words, err1
	}

	return words, nil
}

/*
given a list of words, return the associated ID
SELECT wordID FROM Words_T WHERE word IN ('asdf, meme, lul')
*/
func (w *Word) GetWordIDList(tx *sqlx.Tx, wordSlice []string) ([]int, error) {
	questionMarks := []string{}
	values := make([]interface{}, 0)

	//for every wordSlice unit
	for _, v := range wordSlice {
		questionMarks = append(questionMarks, "?")
		values = append(values, v)
	}

	query := fmt.Sprintf("SELECT wordID FROM Words_T WHERE word IN ( %v )", strings.Join(questionMarks, ","))

	rows, err := w.db.Queryx(query, values...)

	if err != nil {
		return nil, err
	}

	idList := []int{}
	for rows.Next() {
		var idVal int
		err = rows.Scan(&idVal)
		if err != nil {
			return nil, err
		}
		idList = append(idList, idVal)
	}

	if len(idList) == 0 {
		return nil, errors.New("list is empty.")
	}

	return idList, nil
}

/*
rand list of words in words table
*/
