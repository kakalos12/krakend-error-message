package krakend_error

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/martian"
	"github.com/google/martian/parse"
)

func init() {
	parse.Register("body.Modifier", modifierFromJSON)
}

// MarvelModifier contains the private and public Marvel API key
type BetProphetErrorModifier struct {
	StatusCode int
}

// MarvelModifierJSON to Unmarshal the JSON configuration
type BetProphetErrorModifierJSON struct {
	StatusCode int                  `json:"status_code"`
	Scope      []parse.ModifierType `json:"scope"`
}

type BetProphetErrorResponse struct {
	Error_Code int    `json:"error"`
	Error_Body string `json:"error_body"`
}

// ModifyRequest modifies the query string of the request with the given key and value.
func (m *BetProphetErrorModifier) ModifyResponse(res *http.Response) error {

	fmt.Println("Response Status code", res.StatusCode)

	if res.StatusCode >= 400 {
		// close the body first
		body, error := ioutil.ReadAll(res.Body)
		if error != nil {
			fmt.Println("Error reading response body")
			body = []byte{}
		}

		if m.StatusCode != 0 {
			fmt.Println("Setting status to ", m.StatusCode)
			res.StatusCode = m.StatusCode
		}

		res.Body.Close()

		response_obj := BetProphetErrorResponse{
			Error_Code: res.StatusCode,
			Error_Body: string(body),
		}

		r, _ := json.Marshal(response_obj)

		fmt.Println("Response ", string(r))

		res.Body = ioutil.NopCloser(bytes.NewReader(r))
	}
	return nil
}

// MarvelNewModifier returns a request modifier that will set the query string
// at key with the given value. If the query string key already exists all
// values will be overwritten.
func BetProphetNewErrorModifier() martian.ResponseModifier {
	return &BetProphetErrorModifier{}
}

func modifierFromJSON(b []byte) (*parse.Result, error) {
	msg := &BetProphetErrorModifierJSON{}

	if err := json.Unmarshal(b, msg); err != nil {
		fmt.Println("Error ", err)
		return nil, err
	}
	fmt.Println("message : ", msg)
	return parse.NewResult(BetProphetNewErrorModifier(), msg.Scope)
}
