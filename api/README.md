# UatzAPI

### Send Message
```python
import requests
import json

url = "localhost:8080/send/message"

payload = json.dumps({
"device_id": 1,
"recipient_number": "5500000000000",
"message": "UatzApi is COOL!"
})
headers = {
'Content-Type': 'application/json'
}

response = requests.request("POST", url, headers=headers, data=payload)

print(response.text)

```

### Send Sticker
```python
import requests

url = "localhost:8080/send/sticker?device_id=1&recipient_number=5500000000000"

payload = {}
files=[
('sticker',('img.png',open('/Users/user/Downloads/img.png','rb'),'image/png'))
]
headers = {}

response = requests.request("POST", url, headers=headers, data=payload, files=files)

print(response.text)
```

### Device List
```python
import requests

url = "localhost:8080/device"

payload = {}
headers = {}

response = requests.request("GET", url, headers=headers, data=payload)

print(response.text)

```

### Device Connect
```python
import requests

url = "localhost:8080/connect"

payload = {}
headers = {}

response = requests.request("GET", url, headers=headers, data=payload)

print(response.text)
```

### Webhook Add
```python
import requests
import json

url = "localhost:8080/webhook"

payload = json.dumps({
"device_id": 3,
"webhook_url": "http://localhost:3030/whatsapp"
})
headers = {
'Content-Type': 'application/json'
}

response = requests.request("POST", url, headers=headers, data=payload)

print(response.text)
```

### Webhook List
```python
import requests

url = "localhost:8080/webhook"

payload = {}
headers = {}

response = requests.request("GET", url, headers=headers, data=payload)

print(response.text)
```

### Webhook Delete
```python
import requests
import json

url = "localhost:8080/webhook/1"

payload = json.dumps({
"device_id": "3",
"webhook_url": "http://localhost:3030/whatsapp"
})
headers = {
'Content-Type': 'application/json'
}

response = requests.request("DELETE", url, headers=headers, data=payload)

print(response.text)
```