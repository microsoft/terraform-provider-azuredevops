# Schema

## Collection Type: Null vs Zero Value

For a collection type, e.g. List, Set, sometimes it is confusing to choose the right value in `Read()` when it is not returned by API, or returned an empty/zero value. If you simply convert a zero Go value, e.g. `[]string{}`, to the tf types, it will be `null`. This might cause *inconsistent plan* issue if the planned value is actually `[]`.

The principle followed by this provider is to set what is returned by the API. In other word:

- If the API doesn't return the field, or it returns a `null`, which will end up to be a `nil` pointer in the SDK, then set it as `null` in the tf type
- If the API returns a zero value, then set it as the zero value in the tf type

This is build on the hypothesis that the API is roundtrip consistent, meaning a response will return the same value that the request has set. However, it is not the case for ADO APIs. Instead, we have to enforce the schema to *make* it consistent. E.g.

- For API that always returns `[]`, no matter `[]`/`null` is specified, set a default value to the corresponding TF schema, to make sure TF will never send a `null` value in the request. This also requires a special handling in the `Read()`, to conver the null value to zero value before setting to the state.
- For API that always returns `null`, no matter `[]`/`null` is specified, set a zero length validation to the corresponding TF schema, to make sure TF will never send a `[]` value in the request.
