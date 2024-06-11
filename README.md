# Lines
> It's easy to colour in between the lines.

This project aims to create a golang API project that's ready to roll out of the box. It distills a few key opinions 
I've formed over my short time programming into a sort of opinionated framework for building a scalable API.

There are a few key ideas at play here:

- Use a statically typed, compiled language to keep code simple and make it less error prone. We get this for free 
with Golang.
- Store structured data in a schema's database (Postgres) with the ability to dip into semi-structured, unschema'd JSON 
if the right use-case comes up.
- Start with a monolith, but build out 'applications' within the monolith following a domain driven design philosophy. 
    - Each app has it's own store(s), domain layer and ingress / egress.
    - The idea here is that you should be able to peel an app away at any time and create a microservice.
    - Apps should be loosely coupled, with async, event driven communications being the default, however sync calls 
  should be used where it makes sense.
    - While in the monolith, inter-app calls can be via services or use-cases (to prevent the HTTP overhead), however 
  apps should never access another apps business logic.
- Within an app, we follow a clean / hexagonalish architecture where we have:
  - An ingress layer, used to define the interface between the app and the outside world, either through events or HTTP
  calls. Here we define our data transfer schemas and handle requests / responses.
  - A service layer, this is where our business logic lives. We ingest data, orchestrate our data stores to perform 
  some sort of action and send data back out as egress.
  - A store layer, which handles our database operations.

In addition to this, I'm going to solve some common problems and handle some of the 'gubbins' that goes into 
building a production grade app:
- Logging
- Auth, using JWT's for API calls.
- Observability, using Datadog and Sentry to help us to see what our app is doing.
- Containerisation, for easy deployment.
- An admin console, for administering the app.