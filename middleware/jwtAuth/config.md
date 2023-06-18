# Configuring the handler



```go

// To use the defaults...
app.Use(jwtchecker.New(jwtchecker.Config{}))

// We need to pass in a secret otherwise default secret is used
app.Use(jwtchecker.New(jwtchecker.Config{
    Secret: "my_own_custom" // in production please consider using env variables here
}))

// Although we are not required to pass a 
// function handler to the Filter property 
// of our Config, we can always do so if we 
// want to skip our middleware based on a custom condition. To demonstrate that let’s say we want to skip jwt token verification for any request with an X-Secret-Pass having a value of “c5cacd9002a3“. We can do so as shown below.
app.Use(jwtchecker.New(jwtchecker.Config{
    Filter: func (c *fiber.Ctx) bool {
    // if X-Secret-Pass header matches then we skip the jwt validation
    return c.Get("X-Secret-Pass") == "c5cacd9002a3"
    },
}))

```

