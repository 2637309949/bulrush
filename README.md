# Bulrush Framework

![Bulrush flash](./assets/flash.jpg)


## Instruction
1. Install Bulrush
```shell
$ go get github.com/2637309949/juglans
```
2. Init a Bulrush Instance
```shell
import (
    "github.com/2637309949/bulrush"
)
// Use attachs a global Recovery middleware to the router
// Use attachs a Recovery middleware to the user router
// Use attachs a LoggerWithWriter middleware to the user router
app := bulrush.Default()
app.LoadConfig(CONFIGPATH)
app.Inject(&plugins.Middles{}, &routes.Routes{}, &models.Model{})
app.DebugPrintRouteFunc(func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
    fmt.Printf("%5v %9v\n", httpMethod, absolutePath)
})
app.Run()
```
or
```shell
import (
    "github.com/2637309949/bulrush"
)
// No middlewares has been attached
app := bulrush.New()
app.LoadConfig(CONFIGPATH)
app.Inject(&plugins.Middles{}, &routes.Routes{}, &models.Model{})
app.DebugPrintRouteFunc(func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
    fmt.Printf("%5v %9v\n", httpMethod, absolutePath)
})
app.Run()
```
3. For more details, Please reference to [bulrush_template](https://github.com/2637309949/bulrush_template). 
## MIT License

Copyright (c) 2016 Freax

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
