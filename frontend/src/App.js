import React, { Component } from "react";
import "./App.css"
import 'bootstrap/dist/css/bootstrap.min.css';
import { URL } from "./Constants"
import Grid from "./Grid"
import { GridContext } from "./Contexts"

class App extends Component {
  constructor(props) {
    super(props);
    this.state = { 
      grid:{},
      // update:{},
    };
  }

  // ws = new WebSocket("ws://" + window.location.host + "/api/socket")
  ws = new WebSocket("ws://localhost:8080/api/socket")

  handleWs = (ws) => {
    ws.onmessage = evt => {
      // on receiving a message, add it to the list of messages
      var updates = evt.data.split("\n")
      // console.log(updates)
      updates.forEach(update => {
        var u = JSON.parse(update)
        var newGrid = this.state.grid
        var name = this.getBlockName(u.x,u.y)
        newGrid[name] = update
        this.setState({grid:newGrid})
      })
    }
    // ws.onopen = () => {
    //   // on connecting, do nothing but log it to the console
    //   console.log("connected")
    // }
  }

  // i hate javascript
  getBlockName = (x,y) => {
    return ("00" + x).substr(-2,2)+("00" + y).substr(-2,2)
  }

  keyHandling = (e) => {
    fetch(URL+"snake", {
      body: JSON.stringify({"direction":e.key}),
      method: 'POST',
    })
    .catch(err => console.log(err))
  }
  
  componentDidMount = () => {
    // Add Event Listener on component mount
    window.addEventListener("keyup", this.keyHandling);

    this.handleWs(this.ws);

  }
  
  componentWillUnmount = () => {
    // Remove event listener on component unmount
    window.removeEventListener("keyup", this.keyHandling);
  }

  sendCommand = (command) => {
    fetch(URL+"sequence/start", {
      body: JSON.stringify({"name":command, "cycle_seconds":1}),
      method: 'POST',
    })
    .catch(err => console.log(err))
  }

  unsendCommand = () => {
    fetch(URL+"sequence/stop", {
      method: 'GET',
    })
    .catch(err => console.log(err))
  }

  render() {
    return (
      <div className="container-fluid">
        <div className="row">          
          <div className="col-1">
            <div className="btn-group-vertical">
              <button type="button" className="btn btn-primary" onClick={() => this.sendCommand("wave")}>Wave</button>
              <button type="button" className="btn btn-primary" onClick={() => this.sendCommand("mexican_wave")}>Smooth Wave</button>
              <button type="button" className="btn btn-primary" onClick={() => this.sendCommand("alt_wave")}>Alternating Wave</button>
              <button type="button" className="btn btn-primary" onClick={() => this.sendCommand("alt_mexican_wave")}>Alternating Smooth Wave</button>
              <button type="button" className="btn btn-primary" onClick={() => this.sendCommand("random_snake")}>Snake Random</button>
              <button type="button" className="btn btn-primary" onClick={() => this.sendCommand("snake")}>Snake</button>
            </div>
            <div className="btn-group-vertical">
              <button type="button" className="btn btn-danger" onClick={() => this.unsendCommand()}>Stop</button>
             </div>
          </div>
          <div className="col">
            <GridContext.Provider value={this.grid}>
              <Grid/>
            </GridContext.Provider>
          </div>
        </div>
      </div>
    );
  }
}

export default App;