package test

import (
	"bytes"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gw123/GMQ/bootstarp"
	"github.com/gw123/GMQ/core"
	"github.com/gw123/GMQ/services"
	"strings"
	"testing"
)

type Goods struct {
	Id       int      `json:"id"`
	Title    string   `json:"title"`
	OldPrice int      `json:"old_price"`
	Price    int      `json:"price"`
	Desc     string   `json:"desc"`
	Tags     []string `json:"tags"`
	Images   []string `json:"images"`
	Banner   string   `json:"banner"`
}

func getEs() *elasticsearch.Client {
	bootstarp.SetConfigFile("./config.yml")
	config := bootstarp.GetConfig()
	app := core.NewApp(config)
	bootstarp.LoadServices(app)
	app.Start()
	esService := app.GetService("EsService").(*services.EsService)
	es := esService.GetEs()
	return es
}

func TestMapping(t *testing.T) {
	es := getEs()
	indexName := "es-service-test-goods2"
	mapping := `{
  "mappings": {
    "properties": {
      "price":    { "type": "integer" },  
      "title":  { "type": "keyword"  }, 
      "desc":  { "type": "keyword"  }, 
      "tag":   { "type": "text"  },     
      "image":   { "type": "text"  }    
    }
  }
}`
	res, err := es.Index(indexName,
		strings.NewReader(mapping),
		es.Index.WithPretty())

	if err != nil {
		t.Error(err)
	}
	defer res.Body.Close()
	t.Log(res)

}

func TestEsService(t *testing.T) {
	bootstarp.SetConfigFile("./config.yml")
	config := bootstarp.GetConfig()
	app := core.NewApp(config)
	bootstarp.LoadServices(app)
	app.Start()

	esService := app.GetService("EsService").(*services.EsService)

	es := esService.GetEs()
	//打印esinfo
	res, err := es.Info()
	if err != nil {
		t.Error(err)
	}
	defer res.Body.Close()
	//t.Log(res)

	//
	indexName := "es-service-test-goods"

	for i := 1; i < 100000; {
		good := Goods{
			Id:       i,
			Title:    "xyt1",
			OldPrice: i,
			Price:    i,
			Desc:     "desc",
			Tags:     []string{"t1", "t2"},
			Images:   []string{"/img/default.png"},
			Banner:   "海飞丝",
		}
		data, err := json.Marshal(good)
		if err != nil {
			t.Error(err)
		}

		res, err = es.Index(indexName,
			bytes.NewBuffer(data),
			es.Index.WithPretty())
		if err != nil {
			t.Error(err)
		}
		res.Body.Close()
	}

	//res, err = es.Index(
	//	indexName,                    // Index name
	//	strings.NewReader(`{"title" : "Test"}`), // Document body
	//	es.Index.WithDocumentID("1"),            // Document ID
	//	es.Index.WithRefresh("true"),            // Refresh
	//	es.Index.WithPretty(),
	//
	//)
	//if err != nil {
	//	t.Error(err)
	//}
	//defer res.Body.Close()
	//t.Log(res)
}
