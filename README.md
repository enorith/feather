# Feather is a golang http client


## Usage

```golang
	c := feather.NewClient(feather.Options{BaseURI: "https://run.mocky.io/v3/", HttpErrors: true})

	pr, e := c.Get("ed68fdb5-e9e9-4846-bb0f-f208e6820039")
	if e != nil {
		t.Fatalf("request error %v", e)
	}

	pr.Then(func(r Result) {
		t.Logf("response: %v", r)
	})

	pr.Catch(func(err error) {
		if ue, ok := err.(feather.HttpError); ok {
			t.Logf("error: %v[%d]", err, ue.Response.StatusCode)
		}
	})
```