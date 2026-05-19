# Copilot Instructions

## Go Style

- Do not use inline error assignment in `if` statements. Separate the assignment from the check:

```go
// Good
err := doSomething()
if err != nil {
    return err
}

// Bad
if err := doSomething(); err != nil {
    return err
}
```
