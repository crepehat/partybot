import React, { Component } from "react";
import { URL } from "./Constants"
import "./Block.css"

class Block extends Component {
  constructor(props) {
    super(props);
    this.state = { data:{} };
  }

  ws = new WebSocket("ws://"+URL+"block/"+this.props.name)

  handleWs = (ws) => {
    ws.onopen = () => {
      // on connecting, do nothing but log it to the console
      console.log('connected')
    }

    ws.onmessage = evt => {
      // on receiving a message, add it to the list of messages
      var data = JSON.parse(evt.data)
      this.setState({data:data})
    }
  }

  componentDidMount() {
    this.handleWs(this.ws)

  }

  sendMessage = () => {
    this.ws.send("swagger")
  }

  render() {
    return (
      <div>
        <div style={{backgroundColor: 'rgba(255, 0, 0,'+this.state.data.light_magnitude+')'}}>{this.props.name}</div>
      </div>
    );
  }
}

export default Block;