package otcdoc

import (
	"bytes"
	"errors"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type DocumentationPage struct {
	Namespace string
	Metrics   []MetricDocumentationEntry
}

type MetricDocumentationEntry struct {
	MetricId   string
	MetricName string
	Unit       string
}

type RawMetricDocumentationEntry struct {
	MetricId        *goquery.Selection
	Parameter       *goquery.Selection
	Description     *goquery.Selection
	Unit            *goquery.Selection
	ValueRange      *goquery.Selection
	MonitoredObject *goquery.Selection
	MonitorinPeriod *goquery.Selection
}

type DocumentationSource struct {
	Namespace    string
	Url          string
	SubComponent string
}

var DocumentationSources = []DocumentationSource{
	{
		Namespace: "ecs",
		Url:       "https://docs.otc.t-systems.com/elastic-cloud-server/umn/monitoring/basic_ecs_metrics.html",
	},
	{
		Namespace: "bms",
		Url:       "https://docs.otc.t-systems.com/bare-metal-server/umn/server_monitoring/monitored_metrics_with_agent_installed.html"},
	{
		Namespace: "as",
		Url:       "https://docs.otc.t-systems.com/usermanual/as/as_06_0105.html",
	},
	{
		Namespace: "evs",
		Url:       "https://docs.otc.t-systems.com/en-us/usermanual/evs/evs_01_0044.html",
	},
	{
		Namespace: "sfs",
		Url:       "https://docs.otc.t-systems.com/en-us/usermanual/sfs/sfs_01_0047.html",
	},
	{
		Namespace: "efs",
		Url:       "https://docs.otc.t-systems.com/en-us/usermanual/sfs/sfs_01_0048.html",
	},
	{
		Namespace: "cbr",
		Url:       "https://docs.otc.t-systems.com/en-us/usermanual/cbr/cbr_03_0114.html",
	},
	{
		Namespace: "vpc",
		Url:       "https://docs.otc.t-systems.com/usermanual/vpc/vpc010012.html",
	},
	{
		Namespace: "elb",
		Url:       "https://docs.otc.t-systems.com/usermanual/elb/elb_ug_jk_0001.html",
	},
	{
		Namespace: "nat",
		Url:       "https://docs.otc.t-systems.com/usermanual/nat/nat_ces_0002.html",
	},
	{
		Namespace: "waf",
		Url:       "https://docs.otc.t-systems.com/usermanual/waf/waf_01_0092.html",
	},
	{
		Namespace: "dms",
		Url:       "https://docs.otc.t-systems.com/en-us/usermanual/dms/dms-ug-180413002.html",
	},
	{
		Namespace: "dcs",
		Url:       "https://docs.otc.t-systems.com/usermanual/dcs/dcs-ug-0326019.html",
	},
	{
		Namespace:    "rds",
		SubComponent: "MySql",
		Url:          "https://docs.otc.t-systems.com/usermanual/rds/rds_06_0001.html",
	},
	{
		Namespace:    "rds",
		SubComponent: "Postgres",
		Url:          "https://docs.otc.t-systems.com/usermanual/rds/rds_pg_06_0001.html",
	},
	{
		Namespace:    "rds",
		SubComponent: "SqlServer",
		Url:          "https://docs.otc.t-systems.com/usermanual/rds/rds_sqlserver_06_0001.html",
	},
	{
		Namespace: "dds",
		Url:       "https://docs.otc.t-systems.com/usermanual/dds/dds_03_0026.html",
	},
	{
		Namespace: "nosql",
		Url:       "https://docs.otc.t-systems.com/usermanual/nosql/nosql_03_0011.html",
	},
	{
		Namespace: "gaussdb",
		Url:       "https://docs.otc.t-systems.com/usermanual/gaussdb/gaussdb_03_0085.html",
	},
	{
		Namespace: "gaussdbv5",
		Url:       "https://docs.otc.t-systems.com/usermanual/opengauss/opengauss_01_0071.html",
	},
	{
		Namespace: "dws",
		Url:       "https://docs.otc.t-systems.com/usermanual/dws/dws_01_0022.html",
	},
	{
		Namespace: "css",
		Url:       "https://docs.otc.t-systems.com/usermanual/css/css_01_0042.html",
	},
}

func ParseDocumentationPageFromHtmlBytes(htmlBytes []byte, namespace string) (DocumentationPage, error) {
	metricRows := make([]MetricDocumentationEntry, 0)

	reader := bytes.NewReader(htmlBytes)

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return DocumentationPage{}, err
	}

	errs := make([]error, 0)

	doc.Find("table").Each(func(i int, sTable *goquery.Selection) {

		if checkIfTableHasIncorrectStructure(sTable) {
			return
		}

		sTable.Find("tbody tr").Each(
			func(i int, sTableRow *goquery.Selection) {
				rm, err := loadRawDataFromTableRow(sTableRow)

				if err != nil {
					errs = append(errs, err)
				} else {
					m := parseDocumentationRow(rm)
					metricRows = append(metricRows, m)
				}

			},
		)
	})

	return DocumentationPage{
		Namespace: namespace,
		Metrics:   metricRows,
	}, errors.Join(errs...)
}

func parseDocumentationRow(rm RawMetricDocumentationEntry) MetricDocumentationEntry {
	m := MetricDocumentationEntry{
		MetricId:   rm.MetricId.Text(),
		MetricName: rm.Parameter.Text(),
		Unit:       parseUnitFromDescription(rm.Description),
	}
	return m
}

func parseUnitFromDescription(sDescription *goquery.Selection) string {
	unit := ""
	sDescription.Find("p").Each(func(i int, s *goquery.Selection) {
		if strings.HasPrefix(s.Text(), "Unit:") {
			unit, _ = strings.CutPrefix(s.Text(), "Unit:")
			unit = strings.TrimSpace(unit)
		}
	})

	return unit
}

func loadRawDataFromTableRow(sTableRow *goquery.Selection) (RawMetricDocumentationEntry, error) {
	m := RawMetricDocumentationEntry{}
	var err error = nil
	sTableRow.Find("td").Each(func(i int, sTableData *goquery.Selection) {
		switch i {
		case 0:
			m.MetricId = sTableData
		case 1:
			m.Parameter = sTableData
		case 2:
			m.Description = sTableData
		case 3:
			m.ValueRange = sTableData
		case 4:
			m.MonitoredObject = sTableData
		case 5:
			m.MonitorinPeriod = sTableData
		default:
			err = errors.New("invalid amout of columns, the input data or the parsing do not work as expected")
		}
	})

	return m, err
}

func checkIfTableHasIncorrectStructure(sTable *goquery.Selection) bool {
	columnCounter := 0
	sTable.Find("thead tr").First().Find("th").Each(func(i int, sTableHeader *goquery.Selection) {
		columnCounter++
	},
	)

	return columnCounter != 6
}
