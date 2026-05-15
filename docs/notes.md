User ➔ API ➔ PG

---

### PRODUCT
* id
* name
* price_in_cents
* quantity

➔

### ORDERS
* id
* customer_id
* created_at
* status

➔

### ORDERS ITEMS
* product_id
* order_id
* quantity
* price_in_cents

# API high level design
Prefix with /v1

**What we're going to build:**
* GET /health
* GET /products?name&limit=20&offset=0
* POST /orders

***

GET, POST, PATCH, PUT.

| | |
| :--- | :--- |
| DELETE | /orders/{id} |
| PATCH | /orders/{id} |

***

**PRACTICE EXERCISES**

**Exercise at the end:**
* GET /orders/{id} ➔ Fetch aggregated data using JOINS (dont forget to calcualte total price)
* POST /product ➔ Very similar to creating an order
* GET /products/{id}


# Placing an order

[User Icon] ➔ 

```json
{
  "customerId": 123,
  "items": [
    { "productId": 42, "quantity": 2}
  ]
}