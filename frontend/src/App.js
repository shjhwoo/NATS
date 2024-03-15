import React, { useState, useEffect } from 'react';

console.log("render?")
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onopen = () => {
  console.log('WebSocket connected');
};

ws.onclose = () => {
  console.log('WebSocket disconnected');
};

function App() {
  const [message, setMessage] = useState('');
  const [receivedMessage, setReceivedMessage] = useState([]);

  useEffect(() => {
    console.log(receivedMessage)
    if(ws){
      console.log("add new msg")
      ws.onmessage = (event) => {
        console.log([...receivedMessage, event.data])
        setReceivedMessage([...receivedMessage, event.data]);
      }
    }
  },[ws])

  const sendMessage = (e) => {
    e.preventDefault();
    if (ws && message.trim() !== '') {
      ws.send(`{"user":${ws.id},"text":"${message}"}`);
      setMessage('');
    }
  };

  return (
    <div>
      <h1>WebSocket Client</h1>
      <form>
        <input
          type="text"
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          onKeyUp={(e) => {
            if (message.trim() === ''){
              return
            }
            e.preventDefault();
            if (e.key === 'Enter'){
              ws.send(message);
              setMessage('');
            }
          }}
          placeholder="Type a message..."
        />
        <button onClick={sendMessage}>Send</button>
      </form>
      <div>
        <h2>Received Message:</h2>
        <ul>
          {receivedMessage.map((msg, index) => (
            <li key={index}>{msg}</li>
          ))}
        </ul>
      </div>
    </div>
  );
}

export default App;
