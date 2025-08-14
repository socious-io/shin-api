# ShinID Verification Guide

This guide explains how to create verification links for your customers, redirect them to ShinID for verification, and retrieve their verification status.

---

## 1. Generate a Verification Link

Each verification record has a unique **Verification ID**.  
A typical verification link looks like this:

`https://app.shinid.com/connect/redirect/{verification_id}`

For example:  
`https://app.shinid.com/connect/redirect/af7e290e-5cc1-455b-a7d4-7ac1d81b4381`

---

## 2. Add a Customer Identifier (Optional)

You can attach a unique customer reference to the link by adding the `customer` query parameter. This helps you track which of your customers completed the verification.

Example:  
`https://app.shinid.com/connect/redirect/af7e290e-5cc1-455b-a7d4-7ac1d81b4381?customer=customer_unique_username`

When the customer clicks this link, they will be redirected to ShinID’s verification flow.

---

## 3. Customer Verification Process

- The customer is redirected to ShinID.  
- They complete identity verification using ShinID Wallet.  
- Once finished, their verification record will be available for retrieval.

---

## 4. Fetch Verification Status

After your customer has completed verification, you can fetch their status using the ShinID API.

Example request:

```
curl --location 'https://api.shinid.com/verifications/{verification_id}/individuals/{customer_unique_username}'
--header 'apikey: YOUR_API_KEY'
```


Replace:
- `{verification_id}` with the verification ID  
- `{customer_unique_username}` with the unique username you provided (if any)  
- `YOUR_API_KEY` with your ShinID API key  

Response for this API would be something like this:
```
{
  "id": "string",
  "user_id": "string",
  "recipient_id": "string",
  "verification_id": "string",
  "verification": {
    "id": "string",
    "name": "string",
    "description": "string",
    "schema": {
      "id": "string",
      "name": "string",
      "description": "string",
      "attributes": [
        {
          "id": "string",
          "name": "string",
          "description": "string",
          "type": "TEXT | NUMBER | DATETIME",
          "created_at": "datetime"
        }
      ]
    },
    "user": {
      "id": "string",
      "username": "string",
      "email": "string",
      "first_name": "string",
      "last_name": "string",
      "phone": "string",
      "status": "ACTIVE | INACTIVE"
    },
    "type": "SINGLE | MULTI",
    "created_at": "datetime",
    "updated_at": "datetime"
  },
  "body"?: {
    "type": "verification",
    "gender": "string",
    "country": "string",
    "id_number": "string",
    "last_name": "string",
    "first_name": "string",
    "issued_date": "datetime",
    "date_of_birth": "date",
    "document_type": "string",
    "document_number": "string"
  },
  "status": "CREATED | REQUESTED | VERIFIED | FAILED",
  "connection_id": "string",
  "connection_url": "string",
  "connection_at": "datetime",
  "verified_at": "datetime",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```



