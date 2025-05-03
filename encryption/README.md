```go
func main() {
    rsaGen, _ := KeyGeneratorFactory("RSA", 4096) // Pass bits explicitly
    ed25519Gen, _ := KeyGeneratorFactory("ED25519") // No arguments needed
}

```