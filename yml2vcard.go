package main

import (
	"bytes"
	"code.google.com/p/go-uuid/uuid"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

type Employee struct {
	Id           int      `yaml:"id"`
	Name         string   `yaml:"name"`
	Phone        string   `yaml:"phone"`
	ExtraPhones  []string `yaml:"extra_phones"`
	WorkMail     string   `yaml:"work_mail"`
	PersonalMail string   `yaml:"personal_mail"`
	ExtraMails   []string `yaml:"extra_mails"`
}

func main() {
	employeesfile := flag.String("file", "", "path to employees.yaml")
	org := flag.String("org", "", "optional organization name")
	suffix := flag.String("suffix", "", "optional suffix will be added to names")
	prefix := flag.String("prefix", "", "optional prefix will be added to names")
	flag.Parse()

	if *employeesfile == "" {
		flag.PrintDefaults()
		return
	}

	emps, err := parseEmployeesYml(*employeesfile)
	if err != nil {
		log.Println(err)
		return
	}
	for _, emp := range emps {
		fmt.Print(vcardify(emp, *org, *prefix, *suffix))
	}
}

func parseEmployeesYml(file string) ([]Employee, error) {
	var emps []Employee

	var contents, err = ioutil.ReadFile(file)
	if err != nil {
		return emps, err
	}

	err = yaml.Unmarshal(contents, &emps)
	if err != nil {
		return emps, err
	}

	return emps, nil
}

func vcardify(emp Employee, org, prefix, suffix string) string {
	var buf bytes.Buffer

	buf.WriteString("BEGIN:VCARD\n")
	buf.WriteString("VERSION:3.0\n")
	buf.WriteString(fmt.Sprintf(
		"REV:%s\n", time.Now().UTC().Format("2006-01-02T15:04:05Z")))

	buf.WriteString("PRODID:-//Topface//yml2vcard//EN\n")

	names := strings.Split(emp.Name, " ")
	var firstname, lastname string
	switch len(names) {
	case 1:
		firstname, lastname = names[0], ""
	case 2:
		firstname, lastname = names[0], names[1]
	case 3:
		firstname, lastname = names[0], names[2]
	default:
		firstname, lastname = emp.Name, ""
	}

	buf.WriteString(fmt.Sprintf(
		"N:%s;%s;;%s;%s\n", lastname, firstname, prefix, suffix))

	buf.WriteString("FN:")
	if prefix != "" {
		buf.WriteString(prefix + " ")
	}
	buf.WriteString(emp.Name)
	if suffix != "" {
		buf.WriteString(" " + suffix)
	}
	buf.WriteString("\n")

	if org != "" {
		buf.WriteString("ORG:")
		buf.WriteString(org)
		buf.WriteString("\n")
	}

	if emp.Phone != "" {
		buf.WriteString(fmt.Sprintf(
			"TEL:%s\n", emp.Phone))
	}

	if len(emp.ExtraPhones) > 0 {
		for _, p := range emp.ExtraPhones {
			buf.WriteString(fmt.Sprintf(
				"TEL:%s\n", p))
		}
	}

	if emp.WorkMail != "" {
		buf.WriteString(fmt.Sprintf(
			"EMAIL;TYPE=work:%s\n", emp.WorkMail))
	}

	if emp.PersonalMail != "" {
		buf.WriteString(fmt.Sprintf(
			"EMAIL;TYPE=home:%s\n", emp.PersonalMail))
	}

	if len(emp.ExtraMails) > 0 {
		for _, m := range emp.ExtraMails {
			buf.WriteString(fmt.Sprintf(
				"EMAIL:%s\n", m))
		}
	}

	uuid := uuid.NewMD5(uuid.NIL, []byte(emp.Name))
	buf.WriteString(fmt.Sprintf(
		"X-RADICALE-NAME:%s.vcf\n", uuid.String()))
	buf.WriteString(fmt.Sprintf(
		"UID:%s\n", uuid.String()))

	buf.WriteString("END:VCARD\n")

	return buf.String()
}
