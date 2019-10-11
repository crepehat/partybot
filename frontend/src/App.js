import React, { Component } from "react";
import "./App.css"
import 'bootstrap/dist/css/bootstrap.min.css';
import { URL } from "./Constants"
import Grid from "./Grid"

class App extends Component {

  keyHandling = (e) => {
    console.log(e.key)
    fetch("http://"+URL+"snake", {
      body: JSON.stringify({"direction":e.key}),
      method: 'POST',
    })
    .then(resp => resp.json()
    .then(data => console.log(data)))
    .catch(err => console.log(err))
  }
  
  componentDidMount = () => {
    // Add Event Listener on compenent mount
    window.addEventListener("keyup", this.keyHandling);
  }
  
  componentWillUnmount = () => {
    // Remove event listener on compenent unmount
    window.removeEventListener("keyup", this.keyHandling);
  }

  sendCommand = (command) => {
    fetch("http://"+URL+"sequence/start", {
      body: JSON.stringify({"name":command, "cycle_seconds":1}),
      method: 'POST',
    })
    .then(resp => resp.json()
    .then(data => console.log(data)))
    .catch(err => console.log(err))
  }

  unsendCommand = () => {
    fetch("http://"+URL+"sequence/stop", {
      method: 'GET',
    })
    .then(resp => resp.json()
    .then(data => console.log(data)))
    .catch(err => console.log(err))
  }

  render() {
    return (
      <div className="container-fluid">
        <div className="row">          
          <div className="col-1">
            <div class="btn-group-vertical">
              <button type="button" class="btn btn-primary" onClick={() => this.sendCommand("wave")}>Wave</button>
              <button type="button" class="btn btn-primary" onClick={() => this.sendCommand("mexican_wave")}>Smooth Wave</button>
              <button type="button" class="btn btn-primary" onClick={() => this.sendCommand("alt_wave")}>Alternating Wave</button>
              <button type="button" class="btn btn-primary" onClick={() => this.sendCommand("alt_mexican_wave")}>Alternating Smooth Wave</button>
              <button type="button" class="btn btn-primary" onClick={() => this.sendCommand("random_snake")}>Snake Random</button>
              <button type="button" class="btn btn-primary" onClick={() => this.sendCommand("snake")}>Snake</button>
            </div>
            <div class="btn-group-vertical">
              <button type="button" class="btn btn-danger" onClick={() => this.unsendCommand()}>Stop</button>
             </div>
          </div>
          <div className="col">
            <Grid/>
          </div>
        </div>
      </div>
    );
  }
}

export default App;