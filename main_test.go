package main_test

import (
	"bytes"
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	main "github.com/bhmj/resizer"
	"github.com/bhmj/resizer/storage"
	"github.com/gavv/httpexpect"
)

var img40 = []byte{
	'\xff', '\xd8', '\xff', '\xe0', '\x00', '\x10', '\x4a', '\x46', '\x49', '\x46', '\x00', '\x01', '\x02', '\x00', '\x00', '\x64',
	'\x00', '\x64', '\x00', '\x00', '\xff', '\xec', '\x00', '\x11', '\x44', '\x75', '\x63', '\x6b', '\x79', '\x00', '\x01', '\x00',
	'\x04', '\x00', '\x00', '\x00', '\x3c', '\x00', '\x00', '\xff', '\xee', '\x00', '\x0e', '\x41', '\x64', '\x6f', '\x62', '\x65',
	'\x00', '\x64', '\xc0', '\x00', '\x00', '\x00', '\x01', '\xff', '\xdb', '\x00', '\x84', '\x00', '\x06', '\x04', '\x04', '\x04',
	'\x05', '\x04', '\x06', '\x05', '\x05', '\x06', '\x09', '\x06', '\x05', '\x06', '\x09', '\x0b', '\x08', '\x06', '\x06', '\x08',
	'\x0b', '\x0c', '\x0a', '\x0a', '\x0b', '\x0a', '\x0a', '\x0c', '\x10', '\x0c', '\x0c', '\x0c', '\x0c', '\x0c', '\x0c', '\x10',
	'\x0c', '\x0e', '\x0f', '\x10', '\x0f', '\x0e', '\x0c', '\x13', '\x13', '\x14', '\x14', '\x13', '\x13', '\x1c', '\x1b', '\x1b',
	'\x1b', '\x1c', '\x1f', '\x1f', '\x1f', '\x1f', '\x1f', '\x1f', '\x1f', '\x1f', '\x1f', '\x1f', '\x01', '\x07', '\x07', '\x07',
	'\x0d', '\x0c', '\x0d', '\x18', '\x10', '\x10', '\x18', '\x1a', '\x15', '\x11', '\x15', '\x1a', '\x1f', '\x1f', '\x1f', '\x1f',
	'\x1f', '\x1f', '\x1f', '\x1f', '\x1f', '\x1f', '\x1f', '\x1f', '\x1f', '\x1f', '\x1f', '\x1f', '\x1f', '\x1f', '\x1f', '\x1f',
	'\x1f', '\x1f', '\x1f', '\x1f', '\x1f', '\x1f', '\x1f', '\x1f', '\x1f', '\x1f', '\x1f', '\x1f', '\x1f', '\x1f', '\x1f', '\x1f',
	'\x1f', '\x1f', '\x1f', '\x1f', '\x1f', '\x1f', '\x1f', '\x1f', '\x1f', '\x1f', '\x1f', '\x1f', '\x1f', '\xff', '\xc0', '\x00',
	'\x11', '\x08', '\x00', '\x28', '\x00', '\x28', '\x03', '\x01', '\x11', '\x00', '\x02', '\x11', '\x01', '\x03', '\x11', '\x01',
	'\xff', '\xc4', '\x00', '\x82', '\x00', '\x00', '\x03', '\x01', '\x01', '\x01', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00',
	'\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x05', '\x07', '\x04', '\x06', '\x08', '\x01', '\x00', '\x03', '\x01', '\x01',
	'\x01', '\x01', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x04', '\x05', '\x03',
	'\x06', '\x01', '\x02', '\x10', '\x00', '\x02', '\x01', '\x03', '\x03', '\x03', '\x01', '\x07', '\x05', '\x00', '\x00', '\x00',
	'\x00', '\x00', '\x00', '\x00', '\x01', '\x02', '\x03', '\x11', '\x04', '\x05', '\x00', '\x21', '\x12', '\x31', '\x13', '\x06',
	'\x41', '\x51', '\x61', '\x71', '\x81', '\x22', '\x42', '\x92', '\xa1', '\x32', '\xa2', '\x23', '\x07', '\x11', '\x00', '\x02',
	'\x01', '\x03', '\x02', '\x04', '\x06', '\x03', '\x01', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x01',
	'\x02', '\x11', '\x03', '\x04', '\x51', '\x12', '\x21', '\x31', '\x13', '\x05', '\xf0', '\x41', '\x81', '\xb1', '\x22', '\x14',
	'\x61', '\xd1', '\x52', '\x15', '\xff', '\xda', '\x00', '\x0c', '\x03', '\x01', '\x00', '\x02', '\x11', '\x03', '\x11', '\x00',
	'\x3f', '\x00', '\xf5', '\x4e', '\x80', '\x11', '\x5c', '\xf9', '\x9e', '\x16', '\x3b', '\xd9', '\x31', '\xf6', '\xc6', '\x5b',
	'\xec', '\x8c', '\x4c', '\x55', '\xad', '\x2d', '\xa3', '\x67', '\x6a', '\x8a', '\xf2', '\xfa', '\x8f', '\x18', '\xc7', '\x1a',
	'\x1a', '\xfd', '\x5e', '\x87', '\xd4', '\x6b', '\xed', '\x5b', '\x74', '\xaf', '\x91', '\x8c', '\xb2', '\x22', '\x9d', '\x39',
	'\xcb', '\x44', '\x2b', '\x4f', '\x33', '\xce', '\xc9', '\x7e', '\xd6', '\xb1', '\x62', '\x15', '\xa4', '\x5f', '\xdd', '\x6e',
	'\x64', '\x90', '\x4c', '\x00', '\xa7', '\x5a', '\xc7', '\xc4', '\x75', '\xf8', '\x6b', '\x47', '\x66', '\x8a', '\xb5', '\x42',
	'\xeb', '\x35', '\xb9', '\xec', '\x50', '\x95', '\x7d', '\x3f', '\x63', '\x0b', '\x2f', '\x38', '\xc3', '\xdc', '\x64', '\xd3',
	'\x15', '\x32', '\xcb', '\x67', '\x91', '\x6e', '\xb6', '\xf3', '\xa8', '\x04', '\x1a', '\xf1', '\x02', '\xaa', '\x5b', '\x62',
	'\x76', '\x07', '\xa5', '\x48', '\x1d', '\x48', '\x1a', '\xf8', '\x76', '\xda', '\x55', '\xf2', '\x37', '\x86', '\x4c', '\x1c',
	'\xb6', '\xf2', '\x96', '\x87', '\x43', '\xac', '\xcd', '\xc3', '\x40', '\x12', '\x6c', '\xfe', '\x2b', '\x19', '\x6d', '\xe5',
	'\x39', '\x19', '\x12', '\xd8', '\x10', '\x2f', '\x20', '\xb9', '\x65', '\x0f', '\x24', '\x65', '\xe4', '\x8e', '\x38', '\xae',
	'\x13', '\x93', '\xc6', '\xca', '\xe4', '\x09', '\x58', '\x9e', '\x24', '\xd3', '\xd2', '\x94', '\xdb', '\x4f', '\x5b', '\xb7',
	'\xbe', '\xdd', '\x19', '\x0f', '\x27', '\x21', '\xd9', '\xc8', '\xaa', '\x5e', '\x5e', '\xe6', '\x4b', '\xdb', '\xcc', '\xb5',
	'\xf7', '\x95', '\xc1', '\x9d', '\xbd', '\x6b', '\x29', '\xad', '\x6d', '\x99', '\x1e', '\xcf', '\x1b', '\xdb', '\x90', '\x4b',
	'\x1d', '\xcc', '\x41', '\x95', '\x67', '\x32', '\x2c', '\xcb', '\x55', '\x08', '\xcd', '\x40', '\xd1', '\xd2', '\xa7', '\xe1',
	'\xac', '\x7a', '\x72', '\xdf', '\xb2', '\x9f', '\x1d', '\x46', '\x3e', '\xd4', '\x15', '\xbe', '\xb5', '\x7e', '\x6f', '\x86',
	'\xdf', '\x1e', '\xe3', '\xac', '\x2c', '\xa3', '\x21', '\xe6', '\x16', '\x77', '\x77', '\xaa', '\xb2', '\x4d', '\x2c', '\xae',
	'\x76', '\xaa', '\xa8', '\x65', '\x85', '\x9d', '\x68', '\x2b', '\xf6', '\xf6', '\x45', '\x2b', '\x5f', '\x6f', '\x5d', '\xf5',
	'\xb5', '\xf8', '\x6d', '\xb7', '\x44', '\x2d', '\x83', '\x79', '\xdd', '\xc8', '\x72', '\x97', '\x3a', '\x14', '\xad', '\x22',
	'\x5d', '\x0d', '\x00', '\x4a', '\x33', '\x37', '\x17', '\x53', '\xe6', '\xaf', '\xe5', '\xb9', '\xb7', '\x36', '\xb3', '\x34',
	'\x89', '\xce', '\x06', '\x21', '\x8a', '\x91', '\x04', '\x42', '\x9c', '\x94', '\x90', '\x76', '\xa6', '\xfa', '\xa1', '\x86',
	'\xde', '\xce', '\x3a', '\x9c', '\xef', '\x76', '\x4b', '\xab', '\xc3', '\xf9', '\x46', '\x3d', '\x34', '\x4c', '\x35', '\xe2',
	'\x27', '\xbb', '\x83', '\x31', '\x61', '\x2d', '\xa5', '\xb9', '\xba', '\xb8', '\x59', '\x5b', '\xb7', '\x00', '\x21', '\x79',
	'\x13', '\x04', '\xa0', '\x8e', '\x4c', '\x40', '\x5d', '\x89', '\xdf', '\xf4', '\xd2', '\xd9', '\x75', '\xd9', '\xc3', '\x52',
	'\x97', '\x6a', '\xa7', '\x5b', '\x8e', '\x8c', '\xac', '\x6a', '\x71', '\xd1', '\x86', '\x80', '\x26', '\x5e', '\x59', '\x1b',
	'\xa7', '\x92', '\x64', '\x39', '\x0a', '\x77', '\x1a', '\x29', '\x53', '\xde', '\x86', '\x08', '\xd0', '\x37', '\xe5', '\x1b',
	'\x0f', '\x96', '\xa8', '\xe2', '\xbf', '\x89', '\xce', '\x77', '\x58', '\xb5', '\x76', '\xba', '\xa1', '\x4e', '\x99', '\x26',
	'\x8d', '\xbc', '\x52', '\x17', '\x97', '\xc8', '\xec', '\x02', '\x8a', '\xf6', '\x9a', '\x49', '\xa4', '\xf7', '\x22', '\xc4',
	'\xe9', '\x5f', '\xce', '\x54', '\x1f', '\x3d', '\x2d', '\x94', '\xe9', '\x1a', '\x14', '\xbb', '\x5c', '\x1b', '\xbb', '\x5d',
	'\x11', '\x4d', '\xd4', '\xe3', '\xa3', '\x0d', '\x00', '\x73', '\xd9', '\xcf', '\x13', '\x39', '\x8c', '\xaa', '\x5d', '\x4d',
	'\x76', '\xd0', '\xdb', '\xc7', '\x08', '\x8d', '\x16', '\x25', '\x5e', '\xea', '\xb7', '\x22', '\xcd', '\x46', '\x7e', '\x69',
	'\xc5', '\xb6', '\xa8', '\x28', '\x4e', '\xdb', '\x11', '\xaf', '\x63', '\x29', '\x46', '\x55', '\x4c', '\xca', '\xf5', '\x98',
	'\x5c', '\x8d', '\x24', '\xaa', '\x2d', '\x1f', '\xe7', '\x6f', '\xdc', '\xdf', '\x25', '\xfd', '\x5e', '\xc1', '\x00', '\xee',
	'\x7e', '\x5c', '\xca', '\xff', '\x00', '\x0d', '\x33', '\xf6', '\xe5', '\xf8', '\x10', '\xff', '\x00', '\x26', '\xd6', '\xb2',
	'\xf1', '\xe8', '\x35', '\xc1', '\xf8', '\x9a', '\x61', '\xf2', '\x52', '\xdd', '\x41', '\x76', '\xf2', '\xc3', '\x2c', '\x5d',
	'\xb7', '\x8e', '\x55', '\x5e', '\x65', '\x83', '\x02', '\x09', '\x75', '\xe2', '\xbc', '\x57', '\x7e', '\x2a', '\x13', '\xee',
	'\x6d', '\xfa', '\x51', '\x69', '\x4a', '\x52', '\x95', '\x5b', '\x28', '\x59', '\xb3', '\x0b', '\x71', '\xdb', '\x15', '\x41',
	'\xf6', '\xbc', '\x34', '\x3f', '\xff', '\xd9'}

var img10 = []byte{
	'\xff', '\xd8', '\xff', '\xdb', '\x00', '\x84', '\x00', '\x08', '\x06', '\x06', '\x07', '\x06', '\x05', '\x08', '\x07', '\x07',
	'\x07', '\x09', '\x09', '\x08', '\x0a', '\x0c', '\x14', '\x0d', '\x0c', '\x0b', '\x0b', '\x0c', '\x19', '\x12', '\x13', '\x0f',
	'\x14', '\x1d', '\x1a', '\x1f', '\x1e', '\x1d', '\x1a', '\x1c', '\x1c', '\x20', '\x24', '\x2e', '\x27', '\x20', '\x22', '\x2c',
	'\x23', '\x1c', '\x1c', '\x28', '\x37', '\x29', '\x2c', '\x30', '\x31', '\x34', '\x34', '\x34', '\x1f', '\x27', '\x39', '\x3d',
	'\x38', '\x32', '\x3c', '\x2e', '\x33', '\x34', '\x32', '\x01', '\x09', '\x09', '\x09', '\x0c', '\x0b', '\x0c', '\x18', '\x0d',
	'\x0d', '\x18', '\x32', '\x21', '\x1c', '\x21', '\x32', '\x32', '\x32', '\x32', '\x32', '\x32', '\x32', '\x32', '\x32', '\x32',
	'\x32', '\x32', '\x32', '\x32', '\x32', '\x32', '\x32', '\x32', '\x32', '\x32', '\x32', '\x32', '\x32', '\x32', '\x32', '\x32',
	'\x32', '\x32', '\x32', '\x32', '\x32', '\x32', '\x32', '\x32', '\x32', '\x32', '\x32', '\x32', '\x32', '\x32', '\x32', '\x32',
	'\x32', '\x32', '\x32', '\x32', '\x32', '\x32', '\x32', '\x32', '\xff', '\xc0', '\x00', '\x11', '\x08', '\x00', '\x0a', '\x00',
	'\x0a', '\x03', '\x01', '\x22', '\x00', '\x02', '\x11', '\x01', '\x03', '\x11', '\x01', '\xff', '\xc4', '\x01', '\xa2', '\x00',
	'\x00', '\x01', '\x05', '\x01', '\x01', '\x01', '\x01', '\x01', '\x01', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00',
	'\x00', '\x01', '\x02', '\x03', '\x04', '\x05', '\x06', '\x07', '\x08', '\x09', '\x0a', '\x0b', '\x10', '\x00', '\x02', '\x01',
	'\x03', '\x03', '\x02', '\x04', '\x03', '\x05', '\x05', '\x04', '\x04', '\x00', '\x00', '\x01', '\x7d', '\x01', '\x02', '\x03',
	'\x00', '\x04', '\x11', '\x05', '\x12', '\x21', '\x31', '\x41', '\x06', '\x13', '\x51', '\x61', '\x07', '\x22', '\x71', '\x14',
	'\x32', '\x81', '\x91', '\xa1', '\x08', '\x23', '\x42', '\xb1', '\xc1', '\x15', '\x52', '\xd1', '\xf0', '\x24', '\x33', '\x62',
	'\x72', '\x82', '\x09', '\x0a', '\x16', '\x17', '\x18', '\x19', '\x1a', '\x25', '\x26', '\x27', '\x28', '\x29', '\x2a', '\x34',
	'\x35', '\x36', '\x37', '\x38', '\x39', '\x3a', '\x43', '\x44', '\x45', '\x46', '\x47', '\x48', '\x49', '\x4a', '\x53', '\x54',
	'\x55', '\x56', '\x57', '\x58', '\x59', '\x5a', '\x63', '\x64', '\x65', '\x66', '\x67', '\x68', '\x69', '\x6a', '\x73', '\x74',
	'\x75', '\x76', '\x77', '\x78', '\x79', '\x7a', '\x83', '\x84', '\x85', '\x86', '\x87', '\x88', '\x89', '\x8a', '\x92', '\x93',
	'\x94', '\x95', '\x96', '\x97', '\x98', '\x99', '\x9a', '\xa2', '\xa3', '\xa4', '\xa5', '\xa6', '\xa7', '\xa8', '\xa9', '\xaa',
	'\xb2', '\xb3', '\xb4', '\xb5', '\xb6', '\xb7', '\xb8', '\xb9', '\xba', '\xc2', '\xc3', '\xc4', '\xc5', '\xc6', '\xc7', '\xc8',
	'\xc9', '\xca', '\xd2', '\xd3', '\xd4', '\xd5', '\xd6', '\xd7', '\xd8', '\xd9', '\xda', '\xe1', '\xe2', '\xe3', '\xe4', '\xe5',
	'\xe6', '\xe7', '\xe8', '\xe9', '\xea', '\xf1', '\xf2', '\xf3', '\xf4', '\xf5', '\xf6', '\xf7', '\xf8', '\xf9', '\xfa', '\x01',
	'\x00', '\x03', '\x01', '\x01', '\x01', '\x01', '\x01', '\x01', '\x01', '\x01', '\x01', '\x00', '\x00', '\x00', '\x00', '\x00',
	'\x00', '\x01', '\x02', '\x03', '\x04', '\x05', '\x06', '\x07', '\x08', '\x09', '\x0a', '\x0b', '\x11', '\x00', '\x02', '\x01',
	'\x02', '\x04', '\x04', '\x03', '\x04', '\x07', '\x05', '\x04', '\x04', '\x00', '\x01', '\x02', '\x77', '\x00', '\x01', '\x02',
	'\x03', '\x11', '\x04', '\x05', '\x21', '\x31', '\x06', '\x12', '\x41', '\x51', '\x07', '\x61', '\x71', '\x13', '\x22', '\x32',
	'\x81', '\x08', '\x14', '\x42', '\x91', '\xa1', '\xb1', '\xc1', '\x09', '\x23', '\x33', '\x52', '\xf0', '\x15', '\x62', '\x72',
	'\xd1', '\x0a', '\x16', '\x24', '\x34', '\xe1', '\x25', '\xf1', '\x17', '\x18', '\x19', '\x1a', '\x26', '\x27', '\x28', '\x29',
	'\x2a', '\x35', '\x36', '\x37', '\x38', '\x39', '\x3a', '\x43', '\x44', '\x45', '\x46', '\x47', '\x48', '\x49', '\x4a', '\x53',
	'\x54', '\x55', '\x56', '\x57', '\x58', '\x59', '\x5a', '\x63', '\x64', '\x65', '\x66', '\x67', '\x68', '\x69', '\x6a', '\x73',
	'\x74', '\x75', '\x76', '\x77', '\x78', '\x79', '\x7a', '\x82', '\x83', '\x84', '\x85', '\x86', '\x87', '\x88', '\x89', '\x8a',
	'\x92', '\x93', '\x94', '\x95', '\x96', '\x97', '\x98', '\x99', '\x9a', '\xa2', '\xa3', '\xa4', '\xa5', '\xa6', '\xa7', '\xa8',
	'\xa9', '\xaa', '\xb2', '\xb3', '\xb4', '\xb5', '\xb6', '\xb7', '\xb8', '\xb9', '\xba', '\xc2', '\xc3', '\xc4', '\xc5', '\xc6',
	'\xc7', '\xc8', '\xc9', '\xca', '\xd2', '\xd3', '\xd4', '\xd5', '\xd6', '\xd7', '\xd8', '\xd9', '\xda', '\xe2', '\xe3', '\xe4',
	'\xe5', '\xe6', '\xe7', '\xe8', '\xe9', '\xea', '\xf2', '\xf3', '\xf4', '\xf5', '\xf6', '\xf7', '\xf8', '\xf9', '\xfa', '\xff',
	'\xda', '\x00', '\x0c', '\x03', '\x01', '\x00', '\x02', '\x11', '\x03', '\x11', '\x00', '\x3f', '\x00', '\xf5', '\x3d', '\x66',
	'\xf2', '\xea', '\x1d', '\x70', '\xc7', '\x1c', '\xf2', '\xaa', '\x7c', '\x98', '\x50', '\xd8', '\xea', '\x31', '\xc5', '\x75',
	'\x75', '\x5e', '\x5b', '\x78', '\x24', '\x9f', '\xcc', '\x78', '\x63', '\x67', '\x1d', '\x19', '\x94', '\x12', '\x3f', '\x1a',
	'\xb3', '\x55', '\x29', '\x26', '\x95', '\x91', '\x9c', '\x21', '\x28', '\xc9', '\xb6', '\xef', '\x73', '\xff', '\xd9'}

func testHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/image40.jpg":
		w.Header().Set("Content-Type", "image/jpeg")
		w.WriteHeader(http.StatusOK)
		w.Write(img40)
	case "/not_found":
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
	}
}

func createHttpTest(t *testing.T) *httpexpect.Expect {
	cache := storage.NewStore()

	mux := http.NewServeMux()
	src := httptest.NewServer(http.HandlerFunc(testHandler))
	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, src.Listener.Addr().String())
			},
		},
	}
	mux.HandleFunc("/api/v1/resizer/", main.MakeHandler(cli, cache))
	server := httptest.NewServer(mux)
	return httpexpect.New(t, server.URL)
}

func TestMain(t *testing.T) {
	e := createHttpTest(t)

	testCases := []struct {
		Url         string
		Width       string
		Height      string
		ExpectCode  int
		ContentType string
		Err         string
	}{
		{
			Url:         "http://foo.com/image40.jpg",
			Width:       "10",
			Height:      "10",
			ExpectCode:  http.StatusOK,
			ContentType: "image/jpeg",
			Err:         "",
		},
		{
			Url:         "http://foo.com/not_found",
			Width:       "10",
			Height:      "10",
			ExpectCode:  http.StatusNotFound,
			ContentType: "text/plain",
			Err:         "remote image not received",
		},
		{
			Url:         "",
			Width:       "10",
			Height:      "10",
			ExpectCode:  http.StatusBadRequest,
			ContentType: "text/plain",
			Err:         "no url specified",
		},
		{
			Url:         "http://foo.com/image40.jpg",
			Width:       "xxx",
			Height:      "10",
			ExpectCode:  http.StatusBadRequest,
			ContentType: "text/plain",
			Err:         "bad width or height",
		},
	}

	for _, tc := range testCases {
		resp := e.GET("/api/v1/resizer/").WithQuery("url", tc.Url).WithQuery("width", tc.Width).WithQuery("height", tc.Height).Expect()
		resp.Status(tc.ExpectCode)
		resp.ContentType(tc.ContentType)
		body := resp.Body().Raw()
		if tc.ExpectCode != http.StatusOK {
			// error message
			if body != tc.Err {
				t.Errorf("Error message mismatch\n:actual %s\nexpected: %s", body, tc.Err)
			}
		} else {
			// body
			raw := []byte(body)
			if bytes.Compare(raw, img10) != 0 {
				t.Errorf("Resized images mismatch:\nactual: %x\nexpected: %x", raw, img10)
			}
		}
	}
}
