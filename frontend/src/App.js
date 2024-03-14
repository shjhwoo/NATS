import './App.css';
import { connectNats } from './Nats';



function App() {
  connectNats()
  return (
    <div className="App">

    </div>
  );
}

export default App;
