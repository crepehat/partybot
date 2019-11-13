import React, { Component } from "react";
import "./Block.css";
import { GridContext } from "./Contexts";
// In component
import useWebSocket from "react-use-websocket";

// In component function
const [
  sendMessage,
  lastMessage,
  readyState
] = useWebSocket("wss://echo.websocket.org", { onOpen: console.log });

class Block extends Component {
  static contextType = GridContext;

  componentDidMount = () => {
    const grid = this.context;
    // console.log(grid)
    // console.log(grid[this.props.name])
  };

  render() {
    // console.log("rerendering" + this.props.data.name)
    return (
      <div>
        <div>{this.props.name}</div>
      </div>
    );
  }
}

export default Block;
// {/* <div style={{backgroundColor: 'rgba(255, 0, 0,'+this.props.data.light_magnitude+')'}}> */}
// {/* </div> */}
