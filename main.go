package main

import (
	"io"
	"os"
	"log"
	"bytes"
	"errors"
	"context"
	"encoding/json"
	"html/template"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
)

type PageData struct {
	Title  string
	Result string
}

type DistributionSummaryData struct {
	DomainName string `json:"DomainName"`
	Origins    string `json:"Origins"`
	Status     string `json:"Status"`
}

type Response events.APIGatewayProxyResponse

const title string = "Sample AWS CloudFront Page"

var cloudfrontClient *cloudfront.Client

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	tmp := template.New("tmp")
	var dat PageData
	funcMap := template.FuncMap{
		"safehtml": func(text string) template.HTML { return template.HTML(text) },
	}
	buf := new(bytes.Buffer)
	fw := io.Writer(buf)
	dat.Title = title
	distribution, e := getDistribution(ctx)
	if e != nil {
		dat.Result = "ERROR"
	} else {
		dat.Result = distribution
	}
	tmp = template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/index.html", "templates/view.html", "templates/header.html"))
	if e := tmp.ExecuteTemplate(fw, "base", dat); e != nil {
		log.Fatal(e)
	}
	res := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            string(buf.Bytes()),
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
	}
	return res, nil
}

func getDistribution(ctx context.Context)(string, error) {
	if cloudfrontClient == nil {
		cloudfrontClient = cloudfront.NewFromConfig(getConfig(ctx))
	}
	input := &cloudfront.ListDistributionsInput{}

	result, err := cloudfrontClient.ListDistributions(ctx, input)
	if err != nil {
		log.Print(err)
		return "", err
	}

	distributionSummaries := result.DistributionList.Items
	if len(distributionSummaries) < 1 {
		return "", errors.New("No DistributionSummary")
	}
	originsString := ""
	for _, v := range distributionSummaries[0].Origins.Items {
		originsString = originsString + aws.ToString(v.DomainName) + ","
	}
	originsString = originsString[:len(originsString) - 1]
	resultJson, err := json.Marshal(DistributionSummaryData{
		DomainName: aws.ToString(distributionSummaries[0].DomainName),
		Origins: originsString,
		Status: aws.ToString(distributionSummaries[0].Status),
	})
	if err != nil {
		log.Print(err)
		return "", err
	}
	return string(resultJson), nil
}

func getConfig(ctx context.Context) aws.Config {
	var err error
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(os.Getenv("REGION")))
	if err != nil {
		log.Print(err)
	}
	return cfg
}

func main() {
	lambda.Start(HandleRequest)
}
