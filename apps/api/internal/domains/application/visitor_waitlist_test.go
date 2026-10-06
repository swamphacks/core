package application

import "testing"

func TestVisitorSubmissionValidation(t *testing.T) {
	valid := ApplicationSubmissionFields{
		FirstName: "Test", LastName: "Visitor", Age: "18-22",
		Phone: "3525550100", PreferredEmail: "visitor@example.com",
		UniversityEmail: "visitor@ufl.edu", Country: "US",
		Linkedin: "https://www.linkedin.com/in/test", AgeCertification: true,
		School: "University of Florida", Level: "undergrad_three_plus_year",
		Year: "first_year", GraduationYear: "05/2027", Majors: "Computer Science",
		Experience: "first_time", ProjectExperience: "yes", ShirtSize: "M",
		Referral: "website", PictureConsent: "agree",
		InPersonAcknowledgement: "agree", AgreeToConduct: "agree",
		InfoShareAuthorization: "agree",
	}
	if err := validateVisitorSubmission(valid); err != nil {
		t.Fatalf("essays and optional MLH marketing consent should not be required: %v", err)
	}
	cases := []struct {
		name   string
		change func(*ApplicationSubmissionFields)
	}{
		{"missing name", func(d *ApplicationSubmissionFields) { d.FirstName = "" }},
		{"not EDU", func(d *ApplicationSubmissionFields) { d.UniversityEmail = "x@example.com" }},
		{"invalid phone", func(d *ApplicationSubmissionFields) { d.Phone = "abcdefghij" }},
		{"invalid month", func(d *ApplicationSubmissionFields) { d.GraduationYear = "13/2027" }},
		{"not age certified", func(d *ApplicationSubmissionFields) { d.AgeCertification = false }},
		{"photo declined", func(d *ApplicationSubmissionFields) { d.PictureConsent = "no" }},
		{"attendance declined", func(d *ApplicationSubmissionFields) { d.InPersonAcknowledgement = "no" }},
		{"conduct declined", func(d *ApplicationSubmissionFields) { d.AgreeToConduct = "no" }},
		{"sharing declined", func(d *ApplicationSubmissionFields) { d.InfoShareAuthorization = "no" }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data := valid
			tc.change(&data)
			if validateVisitorSubmission(data) == nil {
				t.Fatal("invalid registration accepted")
			}
		})
	}
}
