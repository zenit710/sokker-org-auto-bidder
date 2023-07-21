package repository

import "errors"

var ErrCanNotCreateDbSchema = errors.New("could not create schema for sqlite db")
