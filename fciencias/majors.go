// Package fciencias contains the functions used to parse the majors UNAM's
// Faculty of Science run.
package fciencias

import (
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pablotrinidad/courses-fciencias/crawler"
)

// Major offered in the faculty
type Major struct {
	ID            int
	Name          string
	AcademicPlans []AcademicPlan
}

// AcademicPlan that a major have
type AcademicPlan struct {
	ID         int
	ExternalID int
	Name       string
	Year       int
}

// GetMajors return all majors listed in the website.
func GetMajors() []Major {
	var majors []Major

	// Match content inside parenthesis
	planNameRegex := regexp.MustCompile(`\([^)]+\)`)
	planYearRegex := regexp.MustCompile(`(\d+)`)

	// Fetch main index document
	indexDoc := crawler.GetDocument("docencia/horarios/indice")

	planID := 0

	// Obtain careers and their academic plans
	indexDoc.Find("#info-contenido h2").Each(func(i int, s *goquery.Selection) {
		// Degree name
		var major Major
		major.Name = s.Text()
		major.ID = i + 1

		// Academic plans
		s.Next().Find("a").Each(func(j int, t *goquery.Selection) {

			// Crate new academic plan
			planID++
			var academicPlan AcademicPlan
			academicPlan.ID = planID

			// Parse academic plan name
			APName := planNameRegex.FindString(t.Text())
			APName = strings.Title(APName[1 : len(APName)-1])
			academicPlan.Name = APName

			// Parse plan year
			APYear := planYearRegex.FindString(APName)
			academicPlan.Year, _ = strconv.Atoi(APYear)

			// Obtain plan ID from URL
			ref, _ := t.Attr("href")
			planURL, err := url.Parse(ref)
			if err != nil {
				log.Fatal(err)
			}
			urlComponents := strings.Split(planURL.Path, "/")
			academicPlanID, err := strconv.Atoi(urlComponents[len(urlComponents)-1])
			academicPlan.ExternalID = academicPlanID
			if err != nil {
				log.Fatal(err)
			}
			major.AcademicPlans = append(major.AcademicPlans, academicPlan)
		})

		majors = append(majors, major)
	})

	return majors
}