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
To run the tests
```bash
encore test ./...
```

## API Documentation

### Postman Collection
You can find the Postman collection for this API [here](./pavebank.postman_collection.json).

### Endpoint for bill creation
```
POST /bill
```
Payload:
```
{
    "customerID": "123" (required),
    "currency": "USD" (optional)
}
```

Response
```
{
    "bill": {
        "id": "709091cc-8c61-4ebd-90e6-dd1ca47dcd20",
        "customerID": "123",
        "start_date": "2024-10-14T14:53:28.983302Z",
        "end_date": "2024-11-14T14:53:28.983302Z",
        "currentCharges": null,
        "currency": "USD",
        "status": "open",
        "totalCharges": 0
    }
}
```

### Endpoint for adding a line item to a bill
```
POST /lineitem
```
Payload:
```
{
    "billID": "709091cc-8c61-4ebd-90e6-dd1ca47dcd20" (required),
    "lineItem": {
            "ID": "04fab9a5-557f-4ecb-8f4b-807556a8b2c3",
            "Amount": 55.1,
            "Currency": "USD",
            "Timestamp": "2024-10-14T14:53:28.983302Z"       
        } (required)
}
```

Response
```
{
    "bill": {
        "id": "709091cc-8c61-4ebd-90e6-dd1ca47dcd20",
        "customerID": "123",
        "start_date": "2024-10-14T14:53:28.983302Z",
        "end_date": "2024-11-14T14:53:28.983302Z",
        "currentCharges": [
            {
                "id": "04fab9a5-557f-4ecb-8f4b-807556a8b2c3",
                "description": "",
                "amount": 55.1,
                "timestamp": "2024-10-14T14:53:28.983302Z",
                "currency": "USD",
                "metadata": ""
            }
        ],
        "currency": "USD",
        "status": "open",
        "totalCharges": 55.1
    }
}
```

### Endpoint for updating a bill
```
PATCH /bill/:billID
```
Payload:
```
{
    "status": "closed" (required)
}
```

Response
```
{
    "bill": {
        "id": "709091cc-8c61-4ebd-90e6-dd1ca47dcd20",
        "customerID": "123",
        "start_date": "2024-10-14T14:53:28.983302Z",
        "end_date": "2024-11-14T14:53:28.983302Z",
        "currentCharges": [
            {
                "id": "04fab9a5-557f-4ecb-8f4b-807556a8b2c3",
                "description": "",
                "amount": 55.1,
                "timestamp": "2024-10-14T14:53:28.983302Z",
                "currency": "USD",
                "metadata": ""
            }
        ],
        "currency": "USD",
        "status": "open",
        "totalCharges": 55.1
    }
}
```

### Endpoint for retrieving a bill
```
GET /bill/:billID
```

Response
```
{
    "bill": {
        "id": "709091cc-8c61-4ebd-90e6-dd1ca47dcd20",
        "customerID": "123",
        "start_date": "2024-10-14T14:53:28.983302Z",
        "end_date": "2024-11-14T14:53:28.983302Z",
        "currentCharges": [
            {
                "id": "04fab9a5-557f-4ecb-8f4b-807556a8b2c3",
                "description": "",
                "amount": 55.1,
                "timestamp": "2024-10-14T14:53:28.983302Z",
                "currency": "USD",
                "metadata": ""
            }
        ],
        "currency": "USD",
        "status": "open",
        "totalCharges": 55.1
    }
}
```

## Areas of improvement
Functional
- Add significantly more support and abstractions for currency manipulation. This include support for floating point numbers, more currency options, more precise currency conversion rates, real time currency conversion rates, value markup etc.

- Other billing related functionalities include support for recurring billing, proration, billing cycles, discount applications, late payment penalties etc.


Non-functional
- Authentication and Authorization
- Security headers
- Implement more robust error handling and logging. Error responses would need to be mapped to vague responses like Internal Server Error for public APIs while more specific error messages can be returned to internal services. Logging can include more metadata like context information.
- linting and formatting
