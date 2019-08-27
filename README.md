## JPEG image resizing service

### Run
`dep ensure`
`go run .`  

### Usage

`http://{hostname}:8080/api/v1/resizer/?url={url}&width={width}&height={height}`  

##### Parameters
`{url}` - image source  
`{width}` - new image width  
`{height}` - new image height  

Note: All parameters are required

### Details

Source image cached for 1 hour.  
Resized image cached for 1 hour.  
Only JPEG format is supported.

### Hey

Try it and see for yourself.