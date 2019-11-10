import React, { Component } from "react";
import "./Block.css"

class Block extends Component {
  constructor(props) {
    super(props);
    this.state = { 
      data:{},
      name:'' 
    };
  }

  render() {
    return (
      <div>
        <div style={{backgroundColor: 'rgba(255, 0, 0,'+this.state.data.light_magnitude+')'}}>
        {this.state.name}
        </div>
      </div>
    );

  }
}

export default Block;