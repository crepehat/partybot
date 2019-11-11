import React, { Component } from "react";
import "./Block.css"
import { GridContext } from "./Contexts"

class Block extends Component {
  static contextType = GridContext

  componentDidMount = () => {
    const grid = this.context
    console.log(grid)
    // console.log(grid[this.props.name])
  }

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