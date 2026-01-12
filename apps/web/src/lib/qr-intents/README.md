# QR Intents

This is a small library useful for standardizing and generation and parsing of QR strings. It is only applicable for SwampHacks platform and is used to abstract away action matching etc.

Example:

```ts
  // Receive checkin::userId+eventId
  const input = "checkin::9c2f4c3a-8e7b-4f2e-9a1d-3b6a1c5f9e42+3f6d1a92-1e5c-4c8a-bf2d-7a9e8c4b2d11"
  const (action, inputs) = parseQrIntent(input)
  switch (action) {
    // Do something based on them
  }

```

Here is an example of generating the QR string

```ts
  const userId = "3f6d1a92-1e5c-4c8a-bf2d-7a9e8c4b2d11"
  const eventId = "event1"

  const qrString = generateQrIntentString({Action.CheckIn, userId, eventId})

  // qrString = "checkin::3f6d1a92-1e5c-4c8a-bf2d-7a9e8c4b2d11+event1"
```
