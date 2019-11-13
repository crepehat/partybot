import React, { useMemo, useState, useEffect, useContext } from "react";
import "./App.css";
import "bootstrap/dist/css/bootstrap.min.css";
import useWebSocket from "react-use-websocket";

const SocketContext = React.createContext({});

// Wrapper, initialises things
const App = () => {
  // Initialise websocket
  const STATIC_OPTIONS = useMemo(
    () => ({
      onOpen: () => console.log("opened"),
      shouldReconnect: closeEvent => true
    }),
    []
  );
  const [sendMessage, lastMessage, readyState] = useWebSocket(
    "ws://localhost:8080/api/socket",
    STATIC_OPTIONS
  );

  return (
    <SocketContext.Provider value={[sendMessage, lastMessage, readyState]}>
      <div>
        <Grid />
      </div>
    </SocketContext.Provider>
  );
};

// Grid, gets blocks and puts them in a table
const Grid = () => {
  return (
    <div>
      {ids.map(id => (
        <Block id={id} />
      ))}
    </div>
  );
};

// Block is one of the cells in the table. Receives the update payload and checks if it needs updating
const Block = ({ id }) => {
  const [, lastMessage] = useContext(SocketContext);
  const [light, setLight] = useState("");

  useEffect(() => {
    if (lastMessage !== null) {
      var values = JSON.parse(lastMessage.data);
      if (id in values) {
        setLight(values[id].light_magnitude);
      }
    }
  }, [lastMessage, id]);

  return (
    <div key={id}>
      {id}:{light}
    </div>
  );
};

const ids = [
  "0d",
  "1d",
  "2d",
  "3d",
  "4d",
  "5d",
  "6d",
  "7d",
  "8d",
  "9d",
  "ad",
  "bd",
  "cd",
  "dd",
  "0c",
  "1c",
  "2c",
  "3c",
  "4c",
  "5c",
  "6c",
  "7c",
  "8c",
  "9c",
  "ac",
  "bc",
  "cc",
  "dc",
  "0b",
  "1b",
  "2b",
  "3b",
  "4b",
  "5b",
  "6b",
  "7b",
  "8b",
  "9b",
  "ab",
  "bb",
  "cb",
  "db",
  "0a",
  "1a",
  "2a",
  "3a",
  "4a",
  "5a",
  "6a",
  "7a",
  "8a",
  "9a",
  "aa",
  "ba",
  "ca",
  "da",
  "09",
  "19",
  "29",
  "39",
  "49",
  "59",
  "69",
  "79",
  "89",
  "99",
  "a9",
  "b9",
  "c9",
  "d9",
  "08",
  "18",
  "28",
  "38",
  "48",
  "58",
  "68",
  "78",
  "88",
  "98",
  "a8",
  "b8",
  "c8",
  "d8",
  "07",
  "17",
  "27",
  "37",
  "47",
  "57",
  "67",
  "77",
  "87",
  "97",
  "a7",
  "b7",
  "c7",
  "d7",
  "06",
  "16",
  "26",
  "36",
  "46",
  "56",
  "66",
  "76",
  "86",
  "96",
  "a6",
  "b6",
  "c6",
  "d6",
  "05",
  "15",
  "25",
  "35",
  "45",
  "55",
  "65",
  "75",
  "85",
  "95",
  "a5",
  "b5",
  "c5",
  "d5",
  "04",
  "14",
  "24",
  "34",
  "44",
  "54",
  "64",
  "74",
  "84",
  "94",
  "a4",
  "b4",
  "c4",
  "d4",
  "03",
  "13",
  "23",
  "33",
  "43",
  "53",
  "63",
  "73",
  "83",
  "93",
  "a3",
  "b3",
  "c3",
  "d3",
  "02",
  "12",
  "22",
  "32",
  "42",
  "52",
  "62",
  "72",
  "82",
  "92",
  "a2",
  "b2",
  "c2",
  "d2",
  "01",
  "11",
  "21",
  "31",
  "41",
  "51",
  "61",
  "71",
  "81",
  "91",
  "a1",
  "b1",
  "c1",
  "d1",
  "00",
  "10",
  "20",
  "30",
  "40",
  "50",
  "60",
  "70",
  "80",
  "90",
  "a0",
  "b0",
  "c0",
  "d0"
];

// const

// class App extends Component {
//   constructor(props) {
//     super(props);
//     this.state = {
//       grid:{},
//       // update:{},
//     };
//   }

//   // ws = new WebSocket("ws://" + window.location.host + "/api/socket")
//   ws = new WebSocket("ws://localhost:8080/api/socket")

//   handleWs = (ws) => {
//     ws.onmessage = evt => {
//       // on receiving a message, add it to the list of messages
//       var updates = evt.data.split("\n")
//       // console.log(updates)
//       updates.forEach(update => {
//         var u = JSON.parse(update)
//         var newGrid = this.state.grid
//         var name = this.getBlockName(u.x,u.y)
//         newGrid[name] = update
//         this.setState({grid:newGrid})
//         console.log(u)
//       })
//     }
//     // ws.onopen = () => {
//     //   // on connecting, do nothing but log it to the console
//     //   console.log("connected")
//     // }
//   }

//   // i hate javascript
//   getBlockName = (x,y) => {
//     return ("00" + x).substr(-2,2)+("00" + y).substr(-2,2)
//   }

//   keyHandling = (e) => {
//     fetch(URL+"snake", {
//       body: JSON.stringify({"direction":e.key}),
//       method: 'POST',
//     })
//     .catch(err => console.log(err))
//   }

//   componentDidMount = () => {
//     // Add Event Listener on component mount
//     window.addEventListener("keyup", this.keyHandling);

//     this.handleWs(this.ws);

//   }

//   componentWillUnmount = () => {
//     // Remove event listener on component unmount
//     window.removeEventListener("keyup", this.keyHandling);
//   }

//   sendCommand = (command) => {
//     fetch(URL+"sequence/start", {
//       body: JSON.stringify({"name":command, "cycle_seconds":1}),
//       method: 'POST',
//     })
//     .catch(err => console.log(err))
//   }

//   unsendCommand = () => {
//     fetch(URL+"sequence/stop", {
//       method: 'GET',
//     })
//     .catch(err => console.log(err))
//   }

//   render() {
//     return (
//       <div className="container-fluid">
//         <div className="row">
//           <div className="col-1">
//             <div className="btn-group-vertical">
//               <button type="button" className="btn btn-primary" onClick={() => this.sendCommand("wave")}>Wave</button>
//               <button type="button" className="btn btn-primary" onClick={() => this.sendCommand("mexican_wave")}>Smooth Wave</button>
//               <button type="button" className="btn btn-primary" onClick={() => this.sendCommand("alt_wave")}>Alternating Wave</button>
//               <button type="button" className="btn btn-primary" onClick={() => this.sendCommand("alt_mexican_wave")}>Alternating Smooth Wave</button>
//               <button type="button" className="btn btn-primary" onClick={() => this.sendCommand("random_snake")}>Snake Random</button>
//               <button type="button" className="btn btn-primary" onClick={() => this.sendCommand("snake")}>Snake</button>
//             </div>
//             <div className="btn-group-vertical">
//               <button type="button" className="btn btn-danger" onClick={() => this.unsendCommand()}>Stop</button>
//              </div>
//           </div>
//           <div className="col">
//             <GridContext.Provider value={this.grid}>
//               <Grid/>
//             </GridContext.Provider>
//           </div>
//         </div>
//       </div>
//     );
//   }
// }

export default App;
