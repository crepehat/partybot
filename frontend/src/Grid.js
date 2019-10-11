import React, { Component } from "react";
import "./App.css"
import 'bootstrap/dist/css/bootstrap.min.css';

import Block from "./Block"

import { URL } from "./Constants"


class App extends Component {
  constructor(props) {
    super(props)
    this.state = {
      grid:[]
    }
  }

  getGrid = () => {
    fetch(URL+"grid")
    .then(resp => resp.json())
    .then(grid => this.setState({grid: grid}))
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
  }

  render() {

    var grid = this.constructGrid(this.state.grid)
    
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

export default App;