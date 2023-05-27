package pr

import (
	"fmt"
	"testing"

	"github.com/cli/go-gh/v2/pkg/api"
)

/*
type mockClient struct {
}

func (mc *mockClient) Get() {

}
*/
func TestSynced(t *testing.T) {
	client, err := api.DefaultRESTClient()
	if err != nil {
		t.Fatal(err)
	}

	// org := WordPressOrg
	repo := "gutenberg"
	id := "50731"

	pr, err := GetPr(client, repo, id)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Checking %s\n", pr.Url)

	/*
		for _, s := range synced {
			fmt.Println(s.StatusesUrl)
		}
	*/
	t.Fatal("not implemented")
}

func assertNoError(t testing.TB, got error) {
	t.Helper()
	if got != nil {
		t.Fatal("got an error but didn't expect one")
	}
}
