# Golang Projects

Golang Projects:

- tutorial for basic REST API.
- API for library book model. Used what was learned in the tutoral, then added additional complexity/logic, e.g.: more complex model (manifested in multiple layers of nested JSON values); helper function to check that each new book is assigned a unique ID; some basic error handling, including a specific error check (in create/update) to ensure that a book that is "on shelf" will NOT also contain loan information (i.e., loan value must be null in this case, otherwise a specific error message is returned).
