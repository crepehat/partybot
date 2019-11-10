import React, { Component } from "react";

import Block from "./Block"

import { URL } from "./Constants"


class Grid extends Component {
  constructor(props) {
    super(props)
    this.state = {
      grid:[]
    }
  }

  // ws = new WebSocket("ws://" + window.location.host + "/api/socket")
  ws = new WebSocket("ws://localhost:8080/api/socket")

  handleWs = (ws) => {
    ws.onmessage = evt => {
      // on receiving a message, add it to the list of messages
      console.log(evt)
      // var data = JSON.parse(evt.data)
      // this.setState({data:data})
    }
    ws.onopen = () => {
      // on connecting, do nothing but log it to the console
      console.log("connected")
      this.sendMessage()
    }

  }

  getGrid = () => {
    return fetch(URL+"grid")
    .then(resp => resp.json())
    .then(grid => {
      console.log(grid)
      return this.setState({grid: grid})}
      )
    .then(console.log(this.state.grid))
    .catch(err => console.log(err))
  }

  constructGrid = (grid) => {
    var htmlGrid = ""
    grid.forEach(line => {
      htmlGrid+="<tr>"
      line.forEach(block => {
        htmlGrid+="<td><Block name='"+block+"'/></td>"
      });
      htmlGrid+="</tr>"

    });
    return htmlGrid
  }

  componentDidMount() {
    this.getGrid()
    this.handleWs(this.ws)
    // .catch(err => console.log(err))
  }

  sendMessage = () => {
    console.log("messaging")
    this.ws.send("swagger")
  }

  render() {

    var grid = this.constructGrid(this.state.grid)
    // this.sendMessage()
    return (
      <table class="table">
        <tbody>
         {this.state.grid.map(line => {
           return <tr>
             {line.map(block => <td><Block name={block}/></td>)}
           </tr>
         })}
        </tbody>
      </table>
    );
  }
}

export default Grid;