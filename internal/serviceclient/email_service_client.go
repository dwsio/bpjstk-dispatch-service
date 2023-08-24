package serviceclient

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"

	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/model"
	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/pkg/constants"
)

type SendEmailEnvelopeReq struct {
	XMLName  xml.Name         `xml:"x:Envelope"`
	XmlnsX   string           `xml:"xmlns:x,attr"`
	XmlnsBpj string           `xml:"xmlns:bpj,attr"`
	Header   string           `xml:"x:Header"`
	Body     SendEmailBodyReq `xml:"x:Body"`
}

type SendEmailBodyReq struct {
	SendEmail SendEmailReq `xml:"bpj:sendEmail"`
}

type SendEmailReq struct {
	Cfg        string `xml:"bpj:cfg"`
	From       string `xml:"bpj:from"`
	To         string `xml:"bpj:to"`
	Cc         string `xml:"bpj:cc"`
	Bcc        string `xml:"bpj:bcc"`
	Subject    string `xml:"bpj:subject"`
	Body       string `xml:"bpj:body"`
	IsHTML     string `xml:"bpj:isHTML"`
	BodyHTML   string `xml:"bpj:bodyHTML"`
	IsAttach   string `xml:"bpj:isAttach"`
	Attach     string `xml:"bpj:attach"`
	AttachName string `xml:"bpj:attachName"`
	Avl        string `xml:"bpj:avl"`
}

func (sc *ServiceClient) SendEmail(ctx context.Context, replyCfg string, email *model.Email, traceID string) (string, string, error) {
	_, span := sc.tracer.Start(ctx, "ServiceClient.SendEmail")
	defer span.End()

	client := sc.NewHttpClient()

	xmlReq := getEmailRequestXml(replyCfg, sc.cfg.ServiceClient.EmailService.From, email)

	httpReqBody := strings.NewReader(xmlReq)

	httpReq, err := http.NewRequest(constants.METHOD_POST, sc.cfg.ServiceClient.EmailService.Url, httpReqBody)
	if err != nil {
		sc.logRestMessage(traceID, sc.cfg.ServiceClient.EmailService.Url, constants.METHOD_POST, xmlReq, err, "Create new request", "Failed to create request")
		return "", "", err
	}

	httpReq.Header.Add("Content-Type", "text/xml; charset=utf-8")
	httpReq.Header.Add("SOAPAction", "")

	httpRes, err := client.Do(httpReq)
	if err != nil {
		sc.logRestMessage(traceID, sc.cfg.ServiceClient.EmailService.Url, constants.METHOD_POST, xmlReq, err, "Get http result", "Failed to get result")
		return "", "", err
	}
	defer httpRes.Body.Close()

	httpResBody, err := io.ReadAll(httpRes.Body)
	if err != nil {
		sc.logRestMessage(traceID, sc.cfg.ServiceClient.EmailService.Url, constants.METHOD_POST, xmlReq, err, "Get result body", "Failed to get result body")
		return "", "", err
	}

	if httpRes.StatusCode != http.StatusOK {
		sc.logRestMessage(traceID, sc.cfg.ServiceClient.EmailService.Url, constants.METHOD_POST, string(httpResBody), fmt.Errorf("status code %v", httpRes.StatusCode), "Check HTTP result code", "Failed to send email request")
		return "", "", err
	}

	resMsg := strings.Split(string(httpResBody), "<ax21:msg>")
	resMsg = strings.Split(resMsg[1], "</ax21:msg>")
	msg := resMsg[0]

	resKode := strings.Split(string(httpResBody), "<ax21:kode>")
	resKode = strings.Split(resKode[1], "</ax21:kode>")
	kode := resKode[0]

	sc.logRestMessage(traceID, sc.cfg.ServiceClient.EmailService.Url, constants.METHOD_POST, string(httpResBody), nil, "Success to send email request", "")

	return msg, kode, nil
}

func getEmailRequestXml(replyCfg string, from string, email *model.Email) string {
	isHtml := "F"
	if email.IsHTML {
		isHtml = "T"
	}

	isAttach := "F"
	if email.IsAttach {
		isAttach = "T"
	}

	payload := `
	<x:Envelope xmlns:x="http://schemas.xmlsoap.org/soap/envelope/" xmlns:bpj="http://bpjs.com">
    <x:Header/>
    <x:Body>
        <bpj:sendEmail>
            <bpj:cfg>` + replyCfg + `</bpj:cfg>
            <bpj:from>` + from + `</bpj:from>
            <bpj:to>` + strings.Join(email.RecipientTo, ",") + `</bpj:to>
            <bpj:cc></bpj:cc>
            <bpj:bcc></bpj:bcc>
            <bpj:subject>` + email.Subject + `</bpj:subject>
            <bpj:body>` + email.ContentText + `</bpj:body>
            <bpj:isHTML>` + isHtml + `</bpj:isHTML>
            <bpj:bodyHTML>` + email.ContentHTML + `</bpj:bodyHTML>
            <bpj:isAttach>` + isAttach + `</bpj:isAttach>
            <bpj:attach>` + strings.Join(email.Attachment, ",") + `</bpj:attach>
            <bpj:attachName>` + strings.Join(email.AttachName, ",") + `</bpj:attachName>
            <bpj:avl></bpj:avl>
        </bpj:sendEmail>
    </x:Body>
	</x:Envelope>
	`

	return payload
}
