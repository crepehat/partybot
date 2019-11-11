import React, { Component } from "react";

import Block from "./Block"
import { URL } from "./Constants"


class Grid extends Component {
  constructor(props) {
    super(props);
    this.state = { 
      grid:[],
      // update:{},
    };
  }

  getGrid = () => {
    return fetch(URL+"grid")
    .then(resp => resp.json())
    .then(grid => {
      grid.forEach(row => {
        row.forEach(block => {
          var name = this.getBlockName(block.x,block.y)
          console.log(name)
          this.setState({[name]:block})
        })
      });
      return this.setState({grid:grid})
    })
    .catch(err => console.log(err))
  }

  // i hate javascript
  getBlockName = (x,y) => {
    return ("00" + x).substr(-2,2)+("00" + y).substr(-2,2)
  }

  componentDidMount() {
    this.getGrid()

  }

  render() {

    // var grid = this.constructGrid(this.state.grid)
    return (
      <table className="table">
        <tbody>
         {this.state.grid.reverse().map((line,x) =>            
           <tr key={x}>
             {line.map((block,y) => {
              var name = this.getBlockName(block.x,block.y)
              return <td key={name}><Block key={name} name={name}/></td>
              })
            }
           </tr>
         )}
        </tbody>
      </table>
    );
  }
}

export default Grid;