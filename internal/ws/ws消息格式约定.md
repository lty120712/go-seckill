# ws消息格式约定

## 1.总体格式

## 2.聊天消息格式



## 3.聊天消息示例

```json
{
  "type": "chat", 
  "send_id": 3,    
  "data": {
    "sender_id": 3,      
    "receiver_id": 4,    
    "target_type": 0,    
    "content": [
      {
        "type": "text", 
        "content": "Hello, how are you?" 
      }
    ],
    "type": 0    
  }
}
```
```json
{
  "type": "chat",
  "send_id": 4,
  "data": {
    "sender_id": 4,
    "receiver_id": 3,
    "target_type": 0,
    "content": [
      {
        "type": "text",
        "content": "I'm fine, thanks."
      }
    ],
    "type": 0
  }
}
```

```json
{
  "type": "chat", 
  "send_id": 3,    
  "data": {
    "sender_id": 3,      
    "group_id": 4,    
    "target_type": 1,    
    "content": [
      {
        "type": "text", 
        "content": "Hello, how are you?" 
      }
    ],
    "type": 0    
  }
}

```

心跳检测

```json
{
  "type": "heartbeat",
  "send_id": 3
}
```