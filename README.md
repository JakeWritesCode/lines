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
  - The overall idea is that we initialise a top level app, and that in turn cascades down with each app initialising
  it's own dependencies. This makes testing and spearation of concerns easy.

In addition to this, I'm going to solve some common problems and handle some of the 'gubbins' that goes into 
building a production grade app:
- Logging
- Auth, using JWT's for API calls.
- Observability, using Datadog and Sentry to help us to see what our app is doing.
- Containerisation, for easy deployment.
- An admin console, for administering the app.
- All the boring boilerplate for creating endpoints, serialisers etc that Django does for you.


# Getting Started
To get started, you'll need to have a few things installed on your machine:
- Docker
- Postgres
- Go


# Environment Variables
The following environment variables are required to run the app:
- `LOCAL_DEV` - Set to `true` if you're running the app locally, `false` otherwise.
- `TEST_RUNNER` - Set to `true` if you're running tests, `false` otherwise.
- `USE_SSL` - Set to `true` if you're using SSL, `false` otherwise.
- `SITE_DOMAIN` - The domain of the site.
- `LOG_LEVEL` - The log level of the app.
- `CORS_ORIGINS` - A comma separated list of origins that are allowed to make requests to the app.
- `SENTRY_DSN` - The DSN for Sentry.
- `HTTP_PORT` - The port the app will run on.
- `SECRET_KEY` - The secret key for the app.
- `TOKEN_EXPIRATION_TIME_MINUTES` - The time in minutes that a token will last for.
- `USER_POSTGRES_URL` - The URL for the user postgres database.