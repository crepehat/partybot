import React, { Component } from "react";
import "./Block.css"

const URL = "ws://localhost:8080/"

class Block extends Component {
  render() {
    return (
        <div class="box"><div class="inner">1</div></div>
    );
  }
}

export default Block;