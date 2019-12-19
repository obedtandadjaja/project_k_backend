# Services

Services directory is reserved for code that can potentially be decoupled from the monolith as its own stand-alone microservice.

In order to make this future migration easy, follow these guidelines:
- Methods in services should not be directly called by api. It should be called through a client
- Models owned by these services should not have an ORM relationship with the rest of the entities
  - No FK
  - No join table
- Services should accept a request struct as input and return a response struct as output
  - In this way, we can easily transition to marshall/unmarshall to JSON object
