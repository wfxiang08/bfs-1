package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestHTTPAPI(t *testing.T) {
	// get volume
	// upload
	// uploads
	// delete
	// deletes
	var (
		client http.Client
		s      *Store
		z      *Zookeeper
		w      *multipart.Writer
		f      *os.File
		bw     io.Writer
		req    *http.Request
		resp   *http.Response
		body   []byte
		err    error
		buf    = &bytes.Buffer{}
		tr     = &testRet{}
		bfile  = "./test/block_api_1"
		ifile  = "./test/block_api_1.idx"
	)
	os.Remove(testConf.StoreIndex)
	os.Remove(bfile)
	os.Remove(ifile)
	defer os.Remove(testConf.StoreIndex)
	defer os.Remove(bfile)
	defer os.Remove(ifile)
	if z, err = NewZookeeper([]string{"localhost:2181"}, time.Second*1, "/rack/test-api/"); err != nil {
		t.Errorf("NewZookeeper() error(%v)", err)
		t.FailNow()
	}
	z.DelVolume(1)
	if s, err = NewStore(z, testConf); err != nil {
		t.Errorf("NewStore() error(%v)", err)
		t.FailNow()

	}
	defer s.Close()
	StartAdmin(s, "localhost:6064")
	time.Sleep(1 * time.Second)
	buf.Reset()
	buf.WriteString("n=1&bdir=./test/&idir=./test/")
	if resp, err = http.Post("http://localhost:6064/add_free_volume", "application/x-www-form-urlencoded", buf); err != nil {
		t.Errorf("http.Post() error(%v)", err)
		t.FailNow()
	}
	defer resp.Body.Close()
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		t.Errorf("ioutil.ReadAll() error(%v)", err)
		t.FailNow()
	}
	t.Logf("%s", body)
	if err = json.Unmarshal(body, tr); err != nil {
		t.Errorf("json.Unmarshal() error(%v)", err)
		t.FailNow()
	}
	if tr.Ret != 1 {
		t.FailNow()
	}
	buf.Reset()
	buf.WriteString("vid=1")
	if resp, err = http.Post("http://localhost:6064/add_volume", "application/x-www-form-urlencoded", buf); err != nil {
		t.Errorf("http.Post() error(%v)", err)
		t.FailNow()
	}
	defer resp.Body.Close()
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		t.Errorf("ioutil.ReadAll() error(%v)", err)
		t.FailNow()
	}
	t.Logf("%s", body)
	if err = json.Unmarshal(body, tr); err != nil {
		t.Errorf("json.Unmarshal() error(%v)", err)
		t.FailNow()
	}
	if tr.Ret != 1 {
		t.FailNow()
	}
	StartApi(s, "localhost:6062")
	time.Sleep(1 * time.Second)
	buf.Reset()
	w = multipart.NewWriter(buf)
	if bw, err = w.CreateFormFile("file", "./test/1.jpg"); err != nil {
		t.Errorf("w.CreateFormFile() error(%v)", err)
		t.FailNow()
	}
	if f, err = os.Open("./test/1.jpg"); err != nil {
		t.Errorf("os.Open() error(%v)", err)
		t.FailNow()
	}
	defer f.Close()
	if _, err = io.Copy(bw, f); err != nil {
		t.Errorf("io.Copy() error(%v)", err)
		t.FailNow()
	}
	if err = w.WriteField("vid", "1"); err != nil {
		t.Errorf("w.WriteField() error(%v)", err)
		t.FailNow()
	}
	if err = w.WriteField("key", "15"); err != nil {
		t.Errorf("w.WriteField() error(%v)", err)
		t.FailNow()
	}
	if err = w.WriteField("cookie", "15"); err != nil {
		t.Errorf("w.WriteField() error(%v)", err)
		t.FailNow()
	}
	w.Close()
	if req, err = http.NewRequest("POST", "http://localhost:6062/upload", buf); err != nil {
		t.Errorf("http..NewRequest() error(%v)", err)
		t.FailNow()
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	if resp, err = client.Do(req); err != nil {
		t.Errorf("client.Do() error(%v)", err)
		t.FailNow()
	}
	defer resp.Body.Close()
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		t.Errorf("ioutil.ReadAll() error(%v)", err)
		t.FailNow()
	}
	if err = json.Unmarshal(body, tr); err != nil {
		t.Errorf("json.Unmarshal() error(%v)", err)
		t.FailNow()
	}
	if tr.Ret != 1 {
		t.Errorf("ret: %d", tr.Ret)
		t.FailNow()
	}
	buf.Reset()
	w = multipart.NewWriter(buf)
	if err = w.WriteField("vid", "1"); err != nil {
		t.Errorf("w.WriteField() error(%v)", err)
		t.FailNow()
	}
	for i := 1; i < 10; i++ {
		if bw, err = w.CreateFormFile("file", "./test/"+strconv.Itoa(i)+".jpg"); err != nil {
			t.Errorf("w.CreateFormFile() error(%v)", err)
			t.FailNow()
		}
		if f, err = os.Open("./test/" + strconv.Itoa(i) + ".jpg"); err != nil {
			t.Errorf("os.Open() error(%v)", err)
			t.FailNow()
		}
		defer f.Close()
		if _, err = io.Copy(bw, f); err != nil {
			t.Errorf("io.Copy() error(%v)", err)
			t.FailNow()
		}
		if err = w.WriteField("keys", strconv.Itoa(20+i)); err != nil {
			t.Errorf("w.WriteField() error(%v)", err)
			t.FailNow()
		}
		if err = w.WriteField("cookies", strconv.Itoa(20+i)); err != nil {
			t.Errorf("w.WriteField() error(%v)", err)
			t.FailNow()
		}
	}
	w.Close()
	if req, err = http.NewRequest("POST", "http://localhost:6062/uploads", buf); err != nil {
		t.Errorf("http..NewRequest() error(%v)", err)
		t.FailNow()
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	if resp, err = client.Do(req); err != nil {
		t.Errorf("client.Do() error(%v)", err)
		t.FailNow()
	}
	defer resp.Body.Close()
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		t.Errorf("ioutil.ReadAll() error(%v)", err)
		t.FailNow()
	}
	if err = json.Unmarshal(body, tr); err != nil {
		t.Errorf("json.Unmarshal() error(%v)", err)
		t.FailNow()
	}
	if tr.Ret != 1 {
		t.Errorf("ret: %d", tr.Ret)
		t.FailNow()
	}
	buf.Reset()
	buf.WriteString("vid=1&key=21")
	if resp, err = http.Post("http://localhost:6062/del", "application/x-www-form-urlencoded", buf); err != nil {
		t.Errorf("http.Post() error(%v)", err)
		t.FailNow()
	}
	defer resp.Body.Close()
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		t.Errorf("ioutil.ReadAll() error(%v)", err)
		t.FailNow()
	}
	t.Logf("%s", body)
	if err = json.Unmarshal(body, tr); err != nil {
		t.Errorf("json.Unmarshal() error(%v)", err)
		t.FailNow()
	}
	if tr.Ret != 1 {
		t.FailNow()
	}
	buf.Reset()
	buf.WriteString("vid=1&keys=21&keys=22")
	if resp, err = http.Post("http://localhost:6062/dels", "application/x-www-form-urlencoded", buf); err != nil {
		t.Errorf("http.Post() error(%v)", err)
		t.FailNow()
	}
	defer resp.Body.Close()
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		t.Errorf("ioutil.ReadAll() error(%v)", err)
		t.FailNow()
	}
	t.Logf("%s", body)
	if err = json.Unmarshal(body, tr); err != nil {
		t.Errorf("json.Unmarshal() error(%v)", err)
		t.FailNow()
	}
	if tr.Ret != 1 {
		t.FailNow()
	}
}
