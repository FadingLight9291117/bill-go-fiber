# bill-go

# api

## bill

```ts
GET /list
Query {
    skip: int,
    limit: int,
}
return {
    code,
    data: {
        id: int,
        type: int,
        date: string,
        money: float,
        cls: string,
        label: string,
        options: string,
    }[]
    message,
}

POST /create
Body {
    type: int,
    date: string,
    money: float,
    cls: string,
    label: string,
    options: string,
}
return {
    code,
    data: {
        id: string,
    },
    message,
}

GET /search/:year/:month
Param {
    year: int,
    month: int,
}
Query {
    skip: int,
    limit: int,
}
return {
    code,
    data: {
        id: int,
        type: int,
        date: string,
        money: float,
        cls: string,
        label: string,
        options: string,
    }[]
    message,
}

GET /class
return {
    code,
    data: {
        consume: map<string, string[]>,
        income: string[]
    }
}
```
