# REST API Starter

This is a RESTful API Starter with a single Hello World API endpoint.

## Developing locally

When you have [installed Encore](https://encore.dev/docs/install), you can create a new Encore application and clone this example with this command.

```bash
encore app create my-app-name --example=hello-world
```

## Running locally
```bash
encore run
temporal server start-dev --namespace default
```

While `encore run` is running, open [http://localhost:9400/](http://localhost:9400/) to view Encore's [local developer dashboard](https://encore.dev/docs/observability/dev-dash).

## Using the API

To see that your app is running, you can ping the API.

```bash
curl http://localhost:4000/hello/World
```

## Areas of improvement
Functional
- Store bill and line items in the database to increase reliability and performance. This would also allow more sophisticated querying and reporting capabilities such as filtering, sorting, grouping, paging and aggregations. Currently, the only way to retrieve a bill is by querying the workflow history. The implementation of the database would need to be considered mindfully to avoid mixing business logic of bill manipulation with database query logic, bill manipulation should be done in the workflow and database query logic should be done in the service layer in order to enjoy the benefits of Temporal.

- Add significantly more support and abstractions for currency manipulation. This include support for floating point numbers, more currency options, more precise currency conversion rates, real time currency conversion rates, value markup etc.

- Other billing related functionalities include support for recurring billing, proration, billing cycles, discount applications, late payment penalties etc.


Non-functional
- Authentication and Authorization
- Security headers
- Implement more robust error handling and logging. Error responses would need to be mapped to vague responses like Internal Server Error for public APIs while more specific error messages can be returned to internal services. Logging can include more metadata like context information.
- linting and formatting
