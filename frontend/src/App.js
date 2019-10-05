import React, { Component } from "react";
import "react-table/react-table.css";

const URL = "ws://localhost:8080/"

class App extends Component {

  ws1 = new WebSocket(URL+"1")
  ws2 = new WebSocket(URL+"2")

  handleWs = (ws) => {
    ws.onopen = () => {
      // on connecting, do nothing but log it to the console
      console.log('connected')
    }

    ws.onmessage = evt => {
      // on receiving a message, add it to the list of messages
      console.log(evt)
    }
  }

  componentDidMount() {
    console.log("start")
    this.handleWs(this.ws1)
    console.log("middle")
    this.handleWs(this.ws2)
  }

  sendMessage = () => {
    this.ws1.send("swagger")
    this.ws2.send("swigger")
  }

  render() {
    return (
      <div className="App">
        <button onClick={this.sendMessage}>Swag</button>
      </div>
    );
  }
}

export default App;