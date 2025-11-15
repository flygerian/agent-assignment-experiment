# Fun experiment

Hey! So I was thinking through the problem set again and I'm not sure I did what `would you do next` question justice. So I decided to have a bit of fun and flesh out the code with 2 things in mind.

1. Make it as multi-tenant setup with reasonable load (1m agents online across multiple accounts, so large some small)
2. Simulate load (100 new conversations / second - 8.64m conversations / day)
3. Write unit and integration tests to ensure correctness

I introduced some indexing for the agents accounts as a speed up tactic (difference of 7s / 100 conversations to 1ms / 100 conversations on my souped up custom built linux machine)

# Running the code

To run the load test with a system with goland on it simply run

```
go run cmd/main.go
```

# Running the tests

```
go test ./... -v
```

# Conslusion

The system right now will handle 100 conversation / sec with no issues.

~1ms on my souped up custom built linux machine with 100 conversations / second.
~10ms on my M1 pro macbook with 100 conversations / second.

with minimal memory and CPU overhead.

Obviously this is a simple implementation without the overhead of network calls or data storage, which from my tests will be the biggest bottle neck, but of those perform well then performance should be pretty good

The code if pretty easy to follow (I think) just start with the `cmd/main.go` file and follow the flow of the code.
