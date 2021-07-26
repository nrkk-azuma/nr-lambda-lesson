package main

import (
        "fmt"
        "context"
        "encoding/json"
        "os"
        "strconv"
        "net/http"
        "strings"

        "github.com/aws/aws-sdk-go/aws"
        "github.com/aws/aws-sdk-go/aws/session"
        "github.com/aws/aws-sdk-go/service/lambda"
	newrelic "github.com/newrelic/go-agent/v3/newrelic"
	"github.com/newrelic/go-agent/v3/integrations/nrlambda"
)

type getItemsRequest struct {
    Message    string      `json:"message"`
    Executor   string      `json:"executor"`
    Headers    [][]string `json:"headers"`
}

type getItemsResponseError struct {
    Message string `json:"message"`
}

type getItemsResponseBody struct {
    Result string                 `json:"result"`
    Error  getItemsResponseError  `json:"error"`
}

type getItemsResponseHeaders struct {
    ContentType string `json:"Content-Type"`
}

type getItemsResponse struct {
    StatusCode int                     `json:"statusCode"`
    Headers    getItemsResponseHeaders `json:"headers"`
    Body       getItemsResponseBody    `json:"body"`
}

func handler(ctx context.Context) {
        sess := session.Must(session.NewSessionWithOptions(session.Options{
                SharedConfigState: session.SharedConfigEnable,
        }))
        client := lambda.New(sess, &aws.Config{Region: aws.String("ap-northeast-1")})
 	// The nrlambda handler instrumentation will add the transaction to the
	// context.  Access it using newrelic.FromContext to add additional
	// instrumentation.
        hdrs := http.Header{}
        headers := [][]string{}
	if txn := newrelic.FromContext(ctx); nil != txn {
	        txn.InsertDistributedTraceHeaders(hdrs)
    for key, element := range hdrs {
        headers = append(headers, []string{strings.ToLower(key), element[0]})
    }
        fmt.Println(headers);
		txn.AddAttribute("userLevel", "gold")
		txn.Application().RecordCustomEvent("MyEvent", map[string]interface{}{
			"zip": "zap",
		})
        }
        request := getItemsRequest{"Hello from Go Executor", "GoExecutor", headers}
        payload, err := json.Marshal(request)
        result, err := client.Invoke(&lambda.InvokeInput{FunctionName: aws.String("PythonTest"), Payload: payload})
        if err != nil {
                fmt.Println("Error calling MyGetItemsFunction")
                os.Exit(0)
        }

    var resp getItemsResponse
    err = json.Unmarshal(result.Payload, &resp)
    fmt.Print(resp);
    if err != nil {
        fmt.Println("Error unmarshalling MyGetItemsFunction response")
        os.Exit(0)
    }

    // If the status code is NOT 200, the call failed
    if resp.StatusCode != 200 {
        fmt.Println("Error getting items, StatusCode: " + strconv.Itoa(resp.StatusCode))
        os.Exit(0)
    }

    // If the result is failure, we got an error
    if resp.Body.Result == "failure" {
        fmt.Println("Failed to get items")
        os.Exit(0)
    }
	fmt.Println("hello world")
}

func main() {
        app, err := newrelic.NewApplication(nrlambda.ConfigOption())
	if nil != err {
		fmt.Println("error creating app (invalid config):", err)
	}
	nrlambda.Start(handler, app)
}
