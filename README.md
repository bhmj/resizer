## JPEG image resizing service

### Run
`go build .`  
`./resizer`

### Usage

`http://{hostname}:8080/api/v1/resizer/?url={url}&width={width}&height={height}`  

`{url}` - image source  
`width` - new image width  
`height` - new image height  

Note: All parameters are required

Try it and see for yourself.