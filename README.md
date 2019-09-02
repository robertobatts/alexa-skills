# golexa

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/cc64a91c4de74032bf559c98f6a78cdf)](https://app.codacy.com/app/battaroberto/golexa?utm_source=github.com&utm_medium=referral&utm_content=robertobatts/golexa&utm_campaign=Badge_Grade_Dashboard)
[![Codacy](https://img.shields.io/badge/Code%20Quality%20A-success.svg)](https://app.codacy.com/project/battaroberto/golexa/dashboard?bid=13997035) [![License](https://img.shields.io/badge/License-MPL%202.0-blue.svg)](https://www.mozilla.org/en-US/MPL/2.0/FAQ/)

## Usage

The Golexa struct is the initial interface point with the SDK.  Alexa must be
 initialized first.  The struct is defined as:

```Go
type Golexa struct {
	Triggerable    Triggerable
	AlexaRequest   *AlexaRequest
	TranslationMap map[string]map[string]string
}
```

Golexa embeds the Triggerable interface, whose methods must be implemented

```Go
type Triggerable interface {
	OnLaunch(ctx context.Context, req AlexaRequest, resp *AlexaResponse) error
	OnIntent(ctx context.Context, req AlexaRequest, resp *AlexaResponse) error
}
```

By calling 

```Go
var gxa = &Golexa{...}
gxa.LambdaStart()
```
the functions defined in the Triggerable interface are triggered.

## Example

You can see a complete example of how to use Golexa [here](https://github.com/robertobatts/golexa/blob/master/samples/scorekeeper/scorekeeper.go)
