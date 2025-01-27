# PACT PoC

Evaluation of CDC contract testing using Pact.

### Dependencies

* Java 21
* Go 1.23
* pact-broker
* Docker


### Pact broker

In pact-broker root, exec `docker compose -f docker-compose-dev.yml up`

Visist `localhost:9292/`

### Consumer

Run Consumer (java/go) tests, publish results

### Provider

Run Provider (java/go) tests
Verify contract (pact)