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

type SmsEnvelopeReq struct {
	XMLName  xml.Name   `xml:"x:Envelope"`
	XmlnsX   string     `xml:"xmlns:x,attr"`
	XmlnsBpj string     `xml:"xmlns:bpj,attr"`
	Header   string     `xml:"x:Header"`
	Body     SmsBodyReq `xml:"x:Body"`
}

type SmsBodyReq struct {
	SendSms SendSms `xml:"bpj:sendSMS"`
}

type SendSms struct {
	Username string `xml:"bpj:username"`
	Password string `xml:"bpj:password"`
	Msisdn   string `xml:"bpj:msisdn"`
	Txt      string `xml:"bpj:txt"`
	Avl      string `xml:"bpj:avl"`
}

func (sc *ServiceClient) SendSms(ctx context.Context, sms *model.Sms, traceID string) (string, string, error) {
	_, span := sc.tracer.Start(ctx, "ServiceClient.SendSMS")
	defer span.End()

	client := sc.NewHttpClient()

	xmlReq := getSmsRequestXml(sc.cfg.ServiceClient.SmsService.Username, sc.cfg.ServiceClient.SmsService.Password, sms.RecipientPhoneNumber, sms.Content)

	httpReqBody := strings.NewReader(xmlReq)

	httpReq, err := http.NewRequest(constants.METHOD_POST, sc.cfg.ServiceClient.SmsService.Url, httpReqBody)
	if err != nil {
		sc.logRestMessage(traceID, sc.cfg.ServiceClient.SmsService.Url, constants.METHOD_POST, xmlReq, err, "Create new request", "Failed to create request")
		return "", "", err
	}

	httpReq.Header.Add("Content-Type", "text/xml; charset=utf-8")
	httpReq.Header.Add("SOAPAction", "")

	httpRes, err := client.Do(httpReq)
	if err != nil {
		sc.logRestMessage(traceID, sc.cfg.ServiceClient.SmsService.Url, constants.METHOD_POST, xmlReq, err, "Get http result", "Failed to get result")
		return "", "", err
	}
	defer httpRes.Body.Close()

	httpResBody, err := io.ReadAll(httpRes.Body)
	if err != nil {
		sc.logRestMessage(traceID, sc.cfg.ServiceClient.SmsService.Url, constants.METHOD_POST, xmlReq, err, "Get result body", "Failed to get result body")
		return "", "", err
	}

	if httpRes.StatusCode != http.StatusOK {
		sc.logRestMessage(traceID, sc.cfg.ServiceClient.SmsService.Url, constants.METHOD_POST, string(httpResBody), fmt.Errorf("status code %v", httpRes.StatusCode), "Check HTTP result code", "Failed to send sms request")
		return "", "", err
	}

	resMsg := strings.Split(string(httpResBody), "<ax21:msg>")
	resMsg = strings.Split(resMsg[1], "</ax21:msg>")
	msg := resMsg[0]

	resKode := strings.Split(string(httpResBody), "<ax21:kode>")
	resKode = strings.Split(resKode[1], "</ax21:kode>")
	kode := resKode[0]

	sc.logRestMessage(traceID, sc.cfg.ServiceClient.EmailService.Url, constants.METHOD_POST, string(httpResBody), nil, "Success to send sms request", "")

	return msg, kode, nil
}

func getSmsRequestXml(username string, password string, msisdn string, txt string) string {
	payload := `
	<x:Envelope xmlns:x="http://schemas.xmlsoap.org/soap/envelope/" xmlns:bpj="http://bpjs.com">
    <x:Header/>
    <x:Body>
        <bpj:sendSMS>
            <bpj:username>` + username + `</bpj:username>
            <bpj:password>` + password + `</bpj:password>
            <bpj:msisdn>` + msisdn + `</bpj:msisdn>
            <bpj:txt>` + txt + `</bpj:txt>
            <bpj:avl></bpj:avl>
        </bpj:sendSMS>
    </x:Body>
  </x:Envelope>
  `
	return payload
}
