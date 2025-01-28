# HTMLTemplateCLI
Capable of inputting HTML Template (compatible with GoLang's html/template), an input file, manual input and outputting an HTML file.

## Usage

All parameters are required.

`./htcli --definition input/CoverLetterInputDefinition.json --gohtml input/CoverLetter.gohtml --input input/Me.json --output output/me_coverletter.html`

Inputs consist of:

| abc | a |
|-----| --- |
| --definition | Definition file provides information on all potential inputs |
| --gothml | Go html/template compliance HTML template file which references the defined fields |
| --input | Pre-created input values. Missing values will be prompted in command line |
| --output | Final HTML output file, having combined gohtml file, input file and manually entered CLI input |

## Inputs

### Definition

The format for `--definition` input.

```json
{
  "Name": "Technical Cover Letter",
  "Description": "Cover letter for a technical position",
  "Definitions": [
    {
      "Key": "Separator",
      "Prompt": "Common separator used for concatenating common details ",
      "Default": " | "
    },
    {
      "Key": "FullName",
      "Prompt": "First and Lastname of applicant",
      "Example": "First Lastname"
    },
    {
      "Key": "RoleNamesTop",
      "Prompt": "High level roles to highlight on header",
      "Example": [
        "AWS and Azure Cloud Architect",
        "Solution Architect",
        "Software Engineer"
      ],
      "Type": "StringList"
    }
  ]
}
```

### GoHTML Template

Details of the format can be found here: https://pkg.go.dev/html/template

### Input File

```json
{
  "Separator": " | ",
  "FullName": "Me",
  "ProvinceOrState": "MyProvince",
  "Country": "Canada",
  "PhoneNumber": "(123) 555-1234",
  "EmailAddress": "myemail@gmail.com",
  "LinkedInURL": "https://linkedin.com/in/Brad-Hannah",
  "GitHubURL": "https://github.com/BradHannah",
  "RoleNamesTop": [
    "AWS and Azure Cloud Architect",
    "Solution Architect",
    "Software Engineer"
  ],
  "TagLine": "Something Catchy That Will Capture Your Professional Brand"
}
```