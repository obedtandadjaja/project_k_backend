package models

type ScannableObject interface {
	Scan(dest ...interface{}) error
}
