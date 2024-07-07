package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isTrue := form.Valid()
	if !isTrue {
		t.Errorf("Form valid failed")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Errorf("form shows valid when required fields are missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	r, _ = http.NewRequest("POST", "/whatever", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Errorf("form does not have required fields when it does")
	}
}

func TestForm_MinLength(t *testing.T) {
	postedData := url.Values{}
	f := New(postedData)

	f.MinLength("whatever", 10)
	if f.Valid() {
		t.Errorf("form shows min length for non-existent field")
	}

	postedData = url.Values{}
	postedData.Add("a", "James")
	f = New(postedData)

	f.MinLength("a", 100)
	if f.Valid() {
		t.Errorf("min length showed true when it shouldn't be")
	}

	postedData = url.Values{}
	postedData.Add("b", "abc123")

	f = New(postedData)

	f.MinLength("b", 1)
	if !f.Valid() {
		t.Errorf("min length showed false when it shouldn't be")
	}
}

func TestForm_Has(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	has := form.Has("whatever")
	if has {
		t.Errorf("form shows it has required fields when it does not")
	}

	postedData = url.Values{}
	postedData.Add("a", "a")

	form = New(postedData)

	has = form.Has("a")
	if !has {
		t.Errorf("form shows it does not have required fields when it does not")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	form.IsEmail("something")
	if form.Valid() {
		t.Errorf("showing true for a non-existent email")
	}

	isError := form.Errors.Get("something")
	if isError == "" {
		t.Errorf("should have gotten an error, but did not get one")
	}

	postedData = url.Values{}
	postedData.Add("a", "abc")

	form = New(postedData)

	form.IsEmail("a")
	if form.Valid() {
		t.Errorf("showing true for an invalid email")
	}

	postedData = url.Values{}
	postedData.Add("b", "something@example.com")

	form = New(postedData)

	form.IsEmail("b")
	if !form.Valid() {
		t.Errorf("showing false for a valid email")
	}

	isError = form.Errors.Get("b")
	if isError != "" {
		t.Errorf("should have not gotten an error, but got one")
	}
}
